package media_grpc_client

import (
	"context"
	"github.com/Nistagram-Organization/nistagram-posts/src/dtos"
	"github.com/Nistagram-Organization/nistagram-shared/src/proto"
	"google.golang.org/grpc"
)

type MediaGrpcClient interface {
	SaveMedia(dtos.SaveMediaRequest) (*uint, error)
	GetMedia(dtos.GetMediaRequest) (string, error)
}

type mediaGrpcClient struct {
	address string
}

func NewMediaGrpcClient(docker bool) MediaGrpcClient {
	var address string
	if docker {
		address = "nistagram-media:8089"
	} else {
		address = "127.0.0.1:8089"
	}
	return &mediaGrpcClient{
		address: address,
	}
}

func (c *mediaGrpcClient) SaveMedia(request dtos.SaveMediaRequest) (*uint, error) {
	conn, err := grpc.Dial(c.address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := proto.NewMediaServiceClient(conn)

	r, err := client.SaveMedia(ctx,
		&proto.SaveMediaRequest{
			Image: request.ToMediaMessage(),
		},
	)

	if err != nil {
		return nil, err
	}

	var id *uint
	id = new(uint)
	*id = uint(r.Id)

	return id, nil
}

func (c *mediaGrpcClient) GetMedia(request dtos.GetMediaRequest) (string, error) {
	conn, err := grpc.Dial(c.address, grpc.WithInsecure())
	if err != nil {
		return "", err
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := proto.NewMediaServiceClient(conn)

	r, err := client.GetMedia(ctx,
		&proto.GetMediaRequest{
			Id: request.ID,
		},
	)

	if err != nil {
		return "", err
	}

	return r.Image.ImageBase64, nil
}
