package like

import (
	"github.com/Nistagram-Organization/nistagram-shared/src/model/like"
	"github.com/Nistagram-Organization/nistagram-shared/src/utils/rest_error"
	"github.com/stretchr/testify/mock"
)

type LikeRepositoryMock struct {
	mock.Mock
}

func (l *LikeRepositoryMock) GetByUserAndPost(userEmail string, postId uint) (*like.Like, rest_error.RestErr) {
	args := l.Called(userEmail, postId)
	if args.Get(1) == nil {
		return args.Get(0).(*like.Like), nil
	}
	return nil, args.Get(1).(rest_error.RestErr)
}

func (l *LikeRepositoryMock) Create(like *like.Like) rest_error.RestErr {
	args := l.Called(like)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(rest_error.RestErr)
}

func (l *LikeRepositoryMock) Delete(l2 *like.Like) rest_error.RestErr {
	panic("implement me")
}

func (l *LikeRepositoryMock) GetNumberOfLikes(u uint) (int64, rest_error.RestErr) {
	panic("implement me")
}