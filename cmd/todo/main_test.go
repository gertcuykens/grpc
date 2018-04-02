package main

import (
	"testing"

	pb "github.com/gertcuykens/grpc"
)

func TestTask(t *testing.T) {
	task := pb.Task{
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
