package mock_grpc

import (
	"fmt"

	"github.com/gogo/protobuf/proto"
)

// Msg implements the gomock.Matcher interface
type Msg struct {
	M proto.Message
}

func (r *Msg) Matches(msg interface{}) bool {
	m, ok := msg.(proto.Message)
	if !ok {
		return false
	}
	return proto.Equal(m, r.M)
}

func (r *Msg) String() string {
	return fmt.Sprintf("is %s", r.M)
}
