package main

import (
	"testing"

	gc "github.com/gertcuykens/grpc"
)

func TestTask(t *testing.T) {
	task := gc.Task{
		Text: "test",
		Done: false,
	}
	err := add(&task)
	if err != nil {
		t.Error(err)
	}
	tasks, err := list()
	if err != nil {
		t.Error(err)
	}
	echo(tasks)
}
