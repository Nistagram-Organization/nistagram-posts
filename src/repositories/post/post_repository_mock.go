package post

import (
	"fmt"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/post"
	"github.com/Nistagram-Organization/nistagram-shared/src/utils/rest_error"
	"github.com/stretchr/testify/mock"
)

type PostRepositoryMock struct {
	mock.Mock
}

func (p *PostRepositoryMock) Create(postEntity *post.Post) rest_error.RestErr {
	args := p.Called(postEntity)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(rest_error.RestErr)
}

func (p *PostRepositoryMock) GetUsersPosts(userEmail string) ([]post.Post, rest_error.RestErr) {
	args := p.Called(userEmail)
	if args.Get(1) == nil {
		return args.Get(0).([]post.Post), nil
	}
	return nil, args.Get(1).(rest_error.RestErr)
}

func (p *PostRepositoryMock) Get(u uint) (*post.Post, rest_error.RestErr) {
	args := p.Called(u)
	fmt.Println(args.Get(1))
	if args.Get(1) == nil {
		return args.Get(0).(*post.Post), nil
	}
	return nil, args.Get(1).(rest_error.RestErr)
}

func (p *PostRepositoryMock) GetAll() []post.Post {
	panic("implement me")
}

func (p *PostRepositoryMock) Update(p2 *post.Post) rest_error.RestErr {
	panic("implement me")
}

func (p *PostRepositoryMock) GetInappropriateContent() []post.Post {
	panic("implement me")
}

func (p *PostRepositoryMock) Delete(p2 *post.Post) rest_error.RestErr {
	panic("implement me")
}

func (p *PostRepositoryMock) SearchByTag(s string) ([]post.Post, rest_error.RestErr) {
	panic("implement me")
}