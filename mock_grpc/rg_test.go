package mock_grpc_test

import (
	"fmt"
	"testing"
	"time"

	pb "github.com/gertcuykens/grpc"
	pbmock "github.com/gertcuykens/grpc/mock_grpc"
	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
)

var msg = &pb.RouteNote{
	Location: &pb.Point{Latitude: 17, Longitude: 29},
	Message:  "Taxi-cab",
}

func TestRouteChat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock for the stream returned by RouteChat
	stream := pbmock.NewMockRouteGuide_RouteChatClient(ctrl)
	// set expectation on sending.
	stream.EXPECT().Send(
		gomock.Any(),
	).Return(nil)
	// Set expectation on receiving.
	stream.EXPECT().Recv().Return(msg, nil)
	stream.EXPECT().CloseSend().Return(nil)
	// Create mock for the client interface.
	rgclient := pbmock.NewMockRouteGuideClient(ctrl)
	// Set expectation on RouteChat
	rgclient.EXPECT().RouteChat(
		gomock.Any(),
	).Return(stream, nil)
	if err := testRouteChat(rgclient); err != nil {
		t.Fatalf("Test failed: %v", err)
	}
}

func testRouteChat(client pb.RouteGuideClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := client.RouteChat(ctx)
	if err != nil {
		return err
	}
	if err := stream.Send(msg); err != nil {
		return err
	}
	if err := stream.CloseSend(); err != nil {
		return err
	}
	got, err := stream.Recv()
	if err != nil {
		return err
	}
	if !proto.Equal(got, msg) {
		return fmt.Errorf("stream.Recv() = %v, want %v", got, msg)
	}
	return nil
}
