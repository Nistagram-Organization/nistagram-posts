package comment

import (
	"github.com/Nistagram-Organization/nistagram-shared/src/model/comment"
	"github.com/Nistagram-Organization/nistagram-shared/src/utils/rest_error"
	"github.com/stretchr/testify/mock"
)

type CommentRepositoryMock struct {
	mock.Mock
}

func (c *CommentRepositoryMock) Create(comment *comment.Comment) rest_error.RestErr {
	args := c.Called(comment)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(rest_error.RestErr)
}

func (c *CommentRepositoryMock) GetComments(u uint) ([]comment.Comment, rest_error.RestErr) {
	panic("implement me")
}
