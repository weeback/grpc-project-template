package hello

import (
	"context"

	common "github.com/weeback/grpc-project-template/pb/common"
	pb "github.com/weeback/grpc-project-template/pb/hello"
)

type Repository interface {
	SayHello(ctx context.Context, request *pb.HelloRequest) (*pb.HelloReply, error)
	UseStandardResponse(ctx context.Context, payload *pb.PayloadRequest) (*common.StandardResponse, error)
}
