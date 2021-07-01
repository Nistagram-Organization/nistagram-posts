package media_grpc_client

import (
	"context"
	"github.com/Nistagram-Organization/nistagram-posts/src/dtos"
	"github.com/Nistagram-Organization/nistagram-shared/src/proto"
	"google.golang.org/grpc"
)

type MediaGrpcClient interface {
	SaveMedia(dtos.SaveMediaRequest) (*uint, error)
}

type mediaGrpcClient struct {
}

func NewMediaGrpcClient() MediaGrpcClient {
	return &mediaGrpcClient{}
}

func (c *mediaGrpcClient) SaveMedia(request dtos.SaveMediaRequest) (*uint, error) {
	conn, err := grpc.Dial("127.0.0.1:8089", grpc.WithInsecure())
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
