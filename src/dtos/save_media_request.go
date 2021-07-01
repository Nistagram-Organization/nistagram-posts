package dtos

import (
	"github.com/Nistagram-Organization/nistagram-shared/src/proto"
)

type SaveMediaRequest struct {
	Image string
}

func (r *SaveMediaRequest) ToMediaMessage() *proto.MediaMessage {
	mediaMessage := proto.MediaMessage{
		ImageBase64: r.Image,
	}
	return &mediaMessage
}
