package server

import (
	"fmt"
	"strings"
	"testing"
)

func TestTrimPrefix(t *testing.T) {

	name := "/xxx"
	if strings.HasPrefix(name, "/") {
		name = name[1:]
		fmt.Println(name)
	}

	fmt.Println(strings.TrimPrefix("xxx", "/"))
}
