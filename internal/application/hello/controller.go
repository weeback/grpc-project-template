package hello

import (
	"context"

	"github.com/weeback/grpc-project-template/internal/entity/db"
	"github.com/weeback/grpc-project-template/internal/entity/hello"

	common "github.com/weeback/grpc-project-template/pb/common"
	pb "github.com/weeback/grpc-project-template/pb/hello"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewHelloServiceRepo(exDb db.ExampleDB) hello.Repository {
	return &logging{
		next: &controller{
			exDb: exDb,
		},
	}
}

type controller struct {
	// TODO: define some fields here
	exDb db.ExampleDB
}

func (ins *controller) SayHello(ctx context.Context, request *pb.HelloRequest) (*pb.HelloReply, error) {

	// Validate request
	if err := validateSayRequest(request); err != nil {
		return &pb.HelloReply{Message: err.Error()}, status.Errorf(codes.OK, "%#v", err)
	}
	// TODO implement me
	// ...

	return &pb.HelloReply{Message: "Hello " + request.GetName()}, nil
}

func (ins *controller) UseStandardResponse(ctx context.Context, request *pb.PayloadRequest) (*common.StandardResponse, error) {
	// TODO implement me
	// This is just a placeholder implementation
	return &common.StandardResponse{
		Code:    200,
		Message: "Hello " + request.GetName(),
	}, nil
}
