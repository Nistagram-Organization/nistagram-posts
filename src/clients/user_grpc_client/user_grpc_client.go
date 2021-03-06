package user_grpc_client

import (
	"context"
	"github.com/Nistagram-Organization/nistagram-posts/src/dtos"
	"github.com/Nistagram-Organization/nistagram-shared/src/proto"
	"google.golang.org/grpc"
	"io"
)

type UserGrpcClient interface {
	GetUsername(dtos.GetUsernameRequest) (string, error)
	CheckPostIsInFavorites(dtos.CheckFavoritesRequest) (bool, error)
	CheckIfUserIsTaggable(dtos.CheckTaggableRequest) (bool, error)
	GetFollowingUsers(dtos.GetFollowingUsersRequest) ([]string, error)
	CheckIfUserIsBlocked(dtos.CheckIfUserIsBlockedRequest) (bool, error)
}

type userGrpcClient struct {
	address string
}

func NewUserGrpcClient(docker bool) UserGrpcClient {
	var address string
	if docker {
		address = "nistagram-users:8084"
	} else {
		address = "127.0.0.1:8084"
	}
	return &userGrpcClient{
		address: address,
	}
}

func (u *userGrpcClient) GetUsername(request dtos.GetUsernameRequest) (string, error) {
	conn, err := grpc.Dial(u.address, grpc.WithInsecure())
	if err != nil {
		return "", err
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := proto.NewUserServiceClient(conn)

	r, err := client.GetUsername(ctx,
		&proto.GetUsernameRequest{
			Email: request.Email,
		},
	)

	if err != nil {
		return "", err
	}

	return r.Username, nil
}

func (u *userGrpcClient) CheckPostIsInFavorites(request dtos.CheckFavoritesRequest) (bool, error) {
	conn, err := grpc.Dial(u.address, grpc.WithInsecure())
	if err != nil {
		return false, err
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := proto.NewUserServiceClient(conn)

	r, err := client.CheckIfPostIsInFavorites(ctx,
		&proto.CheckFavoritesRequest{
			Email:  request.Email,
			PostID: uint64(request.PostID),
		},
	)

	if err != nil {
		return false, err
	}

	return r.InFavorites, nil
}

func (u *userGrpcClient) CheckIfUserIsTaggable(request dtos.CheckTaggableRequest) (bool, error) {
	conn, err := grpc.Dial(u.address, grpc.WithInsecure())
	if err != nil {
		return false, err
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := proto.NewUserServiceClient(conn)

	r, err := client.CheckIfUserIsTaggable(ctx,
		&proto.CheckTaggableRequest{
			Username: request.Username,
		},
	)

	if err != nil {
		return false, err
	}

	return r.Taggable, nil
}

func (u *userGrpcClient) GetFollowingUsers(request dtos.GetFollowingUsersRequest) ([]string, error) {
	conn, err := grpc.Dial(u.address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := proto.NewUserServiceClient(conn)

	stream, err := client.GetFollowingUsers(ctx,
		&proto.GetFollowingUsersRequest{
			UserEmail: request.UserEmail,
		},
	)

	if err != nil {
		return nil, err
	}

	var posts []string
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		posts = append(posts, resp.User)
	}

	return posts, nil
}

func (u *userGrpcClient) CheckIfUserIsBlocked(request dtos.CheckIfUserIsBlockedRequest) (bool, error) {
	conn, err := grpc.Dial(u.address, grpc.WithInsecure())
	if err != nil {
		return false, err
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := proto.NewUserServiceClient(conn)

	r, err := client.CheckIfUserIsBlocked(ctx,
		&proto.CheckIfUserIsBlockedRequest{
			User: request.User,
			BlockedUser: request.BlockedUser,
		},
	)

	if err != nil {
		return false, err
	}

	return r.Blocked, nil
}