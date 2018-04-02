package main

import "testing"

func TestSectionsCount(t *testing.T) {
	task := Task{
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
