package mock_grpc_test

import (
	"testing"
	"time"

	gc "github.com/gertcuykens/grpc"
	gcmock "github.com/gertcuykens/grpc/mock_grpc"
	"github.com/golang/mock/gomock"
	"golang.org/x/net/context"
)

func TestAddTask2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTodoServer := gcmock.NewMockTodoServer(ctrl)
	req := &gc.Task{Text: "mock_test", Done: false}
	mockTodoServer.EXPECT().Add(
		gomock.Any(),
		&gcmock.Msg{M: req},
	).Return(&gc.Void{}, nil)

	test := func(t *testing.T, client gc.TodoServer) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		_, err := client.Add(ctx, &gc.Task{Text: "mock_test", Done: false})
		if err != nil {
			t.Error(err)
		}
	}

	test(t, mockTodoServer)
}

func TestList2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTodoServer := gcmock.NewMockTodoServer(ctrl)
	req := &gc.Void{}
	mockTodoServer.EXPECT().List(
		gomock.Any(),
		&gcmock.Msg{M: req},
	).Return(&gc.TaskList{
		Tasks: []*gc.Task{&gc.Task{Text: "mock_test", Done: false}},
	}, nil)

	test := func(t *testing.T, client gc.TodoServer) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r, err := client.List(ctx, &gc.Void{})
		if err != nil {
			t.Error(err)
		}
		t.Log(r.Tasks)
	}

	test(t, mockTodoServer)
}
