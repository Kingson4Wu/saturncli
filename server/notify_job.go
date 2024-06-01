package server

import (
	"errors"
	"github.com/Kingson4Wu/saturncli/base"
	"github.com/Kingson4Wu/saturncli/utils"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"sync"
)

type notifyJob struct {
	Name    string
	Handler any
}

//type JobHandler func(map[string]string, string) bool
//type StoppableJobHandler func(map[string]string, string, chan struct{}) bool

var (
	jobs                   = map[string]*notifyJob{}
	lock                   sync.Mutex
	stoppableJobRunningMap = make(map[string]*sync.Map)
	stoppableJobLock       sync.Mutex
)

func AddJob(name string, handler func(map[string]string, string) bool) error {
	return addJob(name, handler)
}

func AddStoppableJob(name string, handler func(map[string]string, string, chan struct{}) bool) error {
	stoppableJobLock.Lock()
	defer stoppableJobLock.Unlock()
	if _, ok := stoppableJobRunningMap[name]; !ok {
		stoppableJobRunningMap[name] = &sync.Map{}
	}
	return addJob(name, handler)
}

func addJob(name string, handler any) error {
	lock.Lock()
	defer lock.Unlock()
	if _, ok := jobs[name]; ok {
		return errors.New("the job is already exist")
	}
	jobs[name] = &notifyJob{Name: name, Handler: handler}
	return nil
}

func addStoppableJobToRunningMap(jobName, signature string, quit chan struct{}) {
	if runningMap, ok := stoppableJobRunningMap[jobName]; ok {
		_, _ = runningMap.LoadOrStore(signature, quit)
	}
}
func removeStoppableJobFromRunningMap(jobName, signature string) {
	if runningMap, ok := stoppableJobRunningMap[jobName]; ok {
		_, _ = runningMap.LoadAndDelete(signature)
	}
}
func stopTheSpecifiedStoppableJob(jobName, signature string) bool {
	if runningMap, ok := stoppableJobRunningMap[jobName]; ok {
		if v, ok := runningMap.Load(signature); ok {
			if quit, ok := v.(chan struct{}); ok {
				close(quit)
				return true
			}
		}
	}
	return false
}

func stopTheStoppableJob(jobName string) bool {
	if runningMap, ok := stoppableJobRunningMap[jobName]; ok {
		runningMap.Range(func(k, v interface{}) bool {
			if quit, ok := v.(chan struct{}); ok {
				close(quit)
			}
			return true
		})
		return true
	}
	return false
}

func init() {
	/*	AddNotifyJob("hello", func(args map[string]string) bool {
		//logger.Infof("hello, args : %s", args)
		return true
	})*/
}

type ser struct {
	logger   utils.Logger
	sockPath string
}

func NewServer(logger utils.Logger, sockPath string) *ser {
	return &ser{
		logger:   logger,
		sockPath: sockPath,
	}
}

func (s *ser) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	defer func() {
		if err := recover(); err != nil {
			stack := utils.Stack(3)
			s.logger.Errorf("saturn server panic, r: %s, err:%s, stack: %s", r.RequestURI, err, string(stack))
		}
	}()

	name := r.URL.Path
	name = strings.TrimPrefix(name, "/")

	if job, ok := jobs[name]; ok {
		if r.Header.Get(base.StopJobFlag) == "true" {
			s.stopJob(rw, r, job)
			return
		}
		s.runJob(rw, r, job)
		return
	}

	_, _ = rw.Write([]byte("not exist"))
	s.logger.Warnf("saturn server job not exist, name:%s", name)

}

func (s *ser) runJob(rw http.ResponseWriter, r *http.Request, job *notifyJob) {
	name := job.Name
	args := map[string]string{}
	for k, v := range r.URL.Query() {
		args[k] = v[0]
	}
	signature := r.Header.Get(base.RunSignature)
	if signature == "" {
		if v, err := uuid.NewUUID(); err == nil {
			signature = v.String()
		}
	}
	if signature == "" {
		signature = "cron"
	}
	var executeResult bool
	switch h := job.Handler.(type) {
	case func(map[string]string, string) bool:
		executeResult = h(args, signature)
	case func(map[string]string, string, chan struct{}) bool:
		quit := make(chan struct{})
		addStoppableJobToRunningMap(name, signature, quit)
		defer removeStoppableJobFromRunningMap(name, signature)
		executeResult = h(args, signature, quit)
		if isClosed(quit) {
			_, _ = rw.Write([]byte(base.INTERRUPT))
			s.logger.Warnf("saturn server job was interrupted, name:%s, args: %s, signature: %s", name, args, signature)
			return
		}
	default:
	}
	if executeResult {
		_, _ = rw.Write([]byte(base.SUCCESS))
		s.logger.Infof("saturn server job run success, name:%s, args: %s, signature: %s", name, args, signature)
	} else {
		_, _ = rw.Write([]byte(base.FAILURE))
		s.logger.Errorf("saturn server job run fail, name:%s, args: %s, signature: %s", name, args, signature)
	}
}

func isClosed(ch <-chan struct{}) bool {
	select {
	case <-ch:
		return true
	default:
		return false
	}
}

func (s *ser) stopJob(rw http.ResponseWriter, r *http.Request, job *notifyJob) {
	name := job.Name
	args := map[string]string{}
	for k, v := range r.URL.Query() {
		args[k] = v[0]
	}
	var executeResult bool
	signature := r.Header.Get(base.StopSignature)
	if signature != "" {
		//stop the task with specified signature.
		//s.logger.Warnf("saturn server job stop, name:%s, args: %s, signature: %s", name, args, signature)
		executeResult = stopTheSpecifiedStoppableJob(name, signature)
	} else {
		//stop all task with job name.
		//s.logger.Warnf("saturn server job stop, name:%s, args: %s", name, args)
		executeResult = stopTheStoppableJob(name)
	}

	if executeResult {
		_, _ = rw.Write([]byte(base.SUCCESS))
		s.logger.Infof("saturn server job stop success, name:%s, args: %s, signature: %s", name, args, signature)
	} else {
		_, _ = rw.Write([]byte(base.FAILURE))
		s.logger.Errorf("saturn server job stop failure, name:%s, args: %s, signature: %s", name, args, signature)
	}
}
