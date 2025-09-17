package server

import (
	"errors"
	"net/http"
	"strings"
	"sync"

	"github.com/Kingson4Wu/saturncli/base"
	"github.com/Kingson4Wu/saturncli/utils"
	"github.com/google/uuid"
)

// JobHandler processes a scheduled job and returns true on success.
type JobHandler func(map[string]string, string) bool

// StoppableJobHandler is invoked for cancellable jobs; implementations should
// watch the quit channel and stop work promptly when it is closed.
type StoppableJobHandler func(map[string]string, string, chan struct{}) bool

type notifyJob struct {
	name      string
	handler   JobHandler
	stoppable StoppableJobHandler
}

func (j *notifyJob) isStoppable() bool {
	return j != nil && j.stoppable != nil
}

// Registry maintains registered jobs and their active stoppable invocations.
type Registry struct {
	jobsMu    sync.RWMutex
	jobs      map[string]*notifyJob
	running   map[string]*sync.Map
	runningMu sync.RWMutex
}

// NewRegistry constructs an empty job registry for use with a Server.
func NewRegistry() *Registry {
	return &Registry{
		jobs:    make(map[string]*notifyJob),
		running: make(map[string]*sync.Map),
	}
}

var defaultRegistry = NewRegistry()

// AddJob registers a non-stoppable job in the package-level registry.
func AddJob(name string, handler JobHandler) error {
	return defaultRegistry.AddJob(name, handler)
}

// AddStoppableJob registers a stoppable job in the package-level registry.
func AddStoppableJob(name string, handler StoppableJobHandler) error {
	return defaultRegistry.AddStoppableJob(name, handler)
}

// AddJob registers a non-stoppable job against the receiver registry.
func (r *Registry) AddJob(name string, handler JobHandler) error {
	if handler == nil {
		return errors.New("handler is nil")
	}
	job := &notifyJob{name: name, handler: handler}
	return r.registerJob(job)
}

// AddStoppableJob registers a stoppable job against the receiver registry.
func (r *Registry) AddStoppableJob(name string, handler StoppableJobHandler) error {
	if handler == nil {
		return errors.New("handler is nil")
	}
	r.ensureRunningMap(name)
	job := &notifyJob{name: name, stoppable: handler}
	return r.registerJob(job)
}

func (r *Registry) registerJob(job *notifyJob) error {
	if job == nil || strings.TrimSpace(job.name) == "" {
		return errors.New("job name is empty")
	}
	r.jobsMu.Lock()
	defer r.jobsMu.Unlock()
	if _, ok := r.jobs[job.name]; ok {
		return errors.New("the job is already exist")
	}
	r.jobs[job.name] = job
	return nil
}

func (r *Registry) getJob(name string) (*notifyJob, bool) {
	r.jobsMu.RLock()
	defer r.jobsMu.RUnlock()
	job, ok := r.jobs[name]
	return job, ok
}

func (r *Registry) ensureRunningMap(name string) {
	r.runningMu.Lock()
	defer r.runningMu.Unlock()
	if _, ok := r.running[name]; !ok {
		r.running[name] = &sync.Map{}
	}
}

func (r *Registry) runningMap(name string) *sync.Map {
	r.runningMu.RLock()
	defer r.runningMu.RUnlock()
	return r.running[name]
}

func (r *Registry) trackStoppable(jobName, signature string, quit chan struct{}) {
	if signature == "" || quit == nil {
		return
	}
	if runningMap := r.runningMap(jobName); runningMap != nil {
		runningMap.Store(signature, quit)
	}
}

func (r *Registry) untrackStoppable(jobName, signature string) {
	if signature == "" {
		return
	}
	if runningMap := r.runningMap(jobName); runningMap != nil {
		runningMap.Delete(signature)
	}
}

func (r *Registry) stopSpecific(jobName, signature string) bool {
	if signature == "" {
		return false
	}
	if runningMap := r.runningMap(jobName); runningMap != nil {
		if value, ok := runningMap.LoadAndDelete(signature); ok {
			if quit, ok := value.(chan struct{}); ok {
				safeCloseQuit(quit)
			}
			return true
		}
	}
	return false
}

func (r *Registry) stopAll(jobName string) bool {
	if runningMap := r.runningMap(jobName); runningMap != nil {
		stopped := false
		runningMap.Range(func(key, value any) bool {
			runningMap.Delete(key)
			if quit, ok := value.(chan struct{}); ok {
				if safeCloseQuit(quit) {
					stopped = true
				}
			}
			return true
		})
		return stopped
	}
	return false
}

type ServerOption func(*ser)

// WithRegistry swaps the Server's backing registry, enabling isolated job sets.
func WithRegistry(registry *Registry) ServerOption {
	return func(s *ser) {
		if registry != nil {
			s.registry = registry
		}
	}
}

type ser struct {
	logger   utils.Logger
	sockPath string
	registry *Registry
}

func NewServer(logger utils.Logger, sockPath string, opts ...ServerOption) *ser {
	srv := &ser{
		logger:   logger,
		sockPath: sockPath,
		registry: defaultRegistry,
	}
	for _, opt := range opts {
		opt(srv)
	}
	return srv
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

	if job, ok := s.registry.getJob(name); ok {
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
	name := job.name
	args := map[string]string{}
	for k, v := range r.URL.Query() {
		if len(v) == 0 {
			continue
		}
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
	switch {
	case job.handler != nil:
		executeResult = job.handler(args, signature)
	case job.stoppable != nil:
		quit := make(chan struct{})
		s.registry.trackStoppable(name, signature, quit)
		defer s.registry.untrackStoppable(name, signature)
		executeResult = job.stoppable(args, signature, quit)
		if isClosed(quit) {
			_, _ = rw.Write([]byte(base.INTERRUPT))
			s.logger.Warnf("saturn server job was interrupted, name:%s, args: %s, signature: %s", name, args, signature)
			return
		}
	default:
		s.logger.Errorf("saturn server job handler missing, name:%s", name)
		_, _ = rw.Write([]byte(base.FAILURE))
		return
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
	jobName := ""
	if job != nil {
		jobName = job.name
	}
	if job == nil || !job.isStoppable() {
		_, _ = rw.Write([]byte(base.FAILURE))
		s.logger.Errorf("saturn server job stop failure, job is not stoppable, name:%s", jobName)
		return
	}
	name := job.name
	args := map[string]string{}
	for k, v := range r.URL.Query() {
		if len(v) == 0 {
			continue
		}
		args[k] = v[0]
	}
	var executeResult bool
	signature := r.Header.Get(base.StopSignature)
	if signature != "" {
		executeResult = s.registry.stopSpecific(name, signature)
	} else {
		executeResult = s.registry.stopAll(name)
	}

	if executeResult {
		_, _ = rw.Write([]byte(base.SUCCESS))
		s.logger.Infof("saturn server job stop success, name:%s, args: %s, signature: %s", name, args, signature)
	} else {
		_, _ = rw.Write([]byte(base.FAILURE))
		s.logger.Errorf("saturn server job stop failure, name:%s, args: %s, signature: %s", name, args, signature)
	}
}

func safeCloseQuit(quit chan struct{}) bool {
	if quit == nil {
		return false
	}
	select {
	case <-quit:
		return false
	default:
	}
	close(quit)
	return true
}
