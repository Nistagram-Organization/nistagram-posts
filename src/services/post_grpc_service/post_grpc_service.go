package post_grpc_service

import (
	"context"
	"github.com/Nistagram-Organization/nistagram-posts/src/services/post"
	"github.com/Nistagram-Organization/nistagram-shared/src/proto"
)

type postGrpcService struct {
	proto.PostServiceServer
	postService post.PostService
}

func NewPostGrpcService(postService post.PostService) proto.PostServiceServer {
	return &postGrpcService{
		proto.UnimplementedPostServiceServer{},
		postService,
	}
}

func (s *postGrpcService) DecideOnPost(ctx context.Context, decideOnPostRequest *proto.DecideOnPostRequest) (*proto.DecideOnPostResponse, error) {
	id := uint(decideOnPostRequest.Post)
	deletePost := decideOnPostRequest.Delete

	if err := s.postService.DecideOnContent(id, deletePost); err != nil {
		return nil, err
	}

	response := proto.DecideOnPostResponse{Success: true}

	return &response, nil
}
