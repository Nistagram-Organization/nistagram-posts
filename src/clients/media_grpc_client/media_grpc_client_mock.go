package media_grpc_client

import (
	"github.com/Nistagram-Organization/nistagram-posts/src/dtos"
	"github.com/stretchr/testify/mock"
)

type MediaGrpcClientMock struct {
	mock.Mock
}

func (c *MediaGrpcClientMock) SaveMedia(request dtos.SaveMediaRequest) (*uint, error) {
	args := c.Called(request)
	if args.Get(1) != nil {
		return args.Get(0).(*uint), nil
	}
	return nil, args.Get(1).(error)
}

func (c *MediaGrpcClientMock) GetMedia(request dtos.GetMediaRequest) (string, error) {
	panic("implement me")
}
