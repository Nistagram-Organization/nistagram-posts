package dislike

import (
	"github.com/Nistagram-Organization/nistagram-shared/src/model/dislike"
	"github.com/Nistagram-Organization/nistagram-shared/src/utils/rest_error"
	"github.com/stretchr/testify/mock"
)

type DislikeRepositoryMock struct {
	mock.Mock
}

func (d *DislikeRepositoryMock) GetByUserAndPost(userEmail string, postId uint) (*dislike.Dislike, rest_error.RestErr) {
	args := d.Called(userEmail, postId)
	if args.Get(1) == nil {
		return args.Get(0).(*dislike.Dislike), nil
	}
	return nil, args.Get(1).(rest_error.RestErr)
}

func (d *DislikeRepositoryMock) Create(dislike *dislike.Dislike) rest_error.RestErr {
	args := d.Called(dislike)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(rest_error.RestErr)
}

func (d *DislikeRepositoryMock) Delete(d2 *dislike.Dislike) rest_error.RestErr {
	panic("implement me")
}

func (d *DislikeRepositoryMock) GetNumberOfDislikes(u uint) (int64, rest_error.RestErr) {
	panic("implement me")
}