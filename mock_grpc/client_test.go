package mock_grpc_test

import (
	"testing"
	"time"

	pb "github.com/gertcuykens/grpc"
	pbmock "github.com/gertcuykens/grpc/mock_grpc"
	"github.com/golang/mock/gomock"
	"golang.org/x/net/context"
)

func TestAddTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTodoClient := pbmock.NewMockTodoClient(ctrl)
	req := &pb.Task{Text: "mock_test", Done: false}
	mockTodoClient.EXPECT().Add(
		gomock.Any(),
		&pbmock.Msg{M: req},
	).Return(&pb.Void{}, nil)

	test := func(t *testing.T, client pb.TodoClient) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		_, err := client.Add(ctx, &pb.Task{Text: "mock_test", Done: false})
		if err != nil {
			t.Error(err)
		}
	}

	test(t, mockTodoClient)
}

func TestList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTodoClient := pbmock.NewMockTodoClient(ctrl)
	req := &pb.Void{}
	mockTodoClient.EXPECT().List(
		gomock.Any(),
		&pbmock.Msg{M: req},
	).Return(&pb.TaskList{
		Tasks: []*pb.Task{&pb.Task{Text: "mock_test", Done: false}},
	}, nil)

	test := func(t *testing.T, client pb.TodoClient) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r, err := client.List(ctx, &pb.Void{})
		if err != nil {
			t.Error(err)
		}
		t.Log(r.Tasks)
	}

	test(t, mockTodoClient)
}
