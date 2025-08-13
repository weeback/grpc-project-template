package hello

import (
	"context"

	common "github.com/weeback/grpc-project-template/pb/common"
	pb "github.com/weeback/grpc-project-template/pb/hello"
	"github.com/weeback/grpc-project-template/pkg/logger"
	"go.uber.org/zap"
)

type logging struct {
	next *controller
}

func (ins *logging) SayHello(ctx context.Context, request *pb.HelloRequest) (result *pb.HelloReply, err error) {
	defer func(entry *zap.Logger) {
		if err != nil {
			entry.Error("Failed to say hello",
				zap.String("name", request.GetName()),
				zap.Error(err))
		} else {
			entry.Debug("Successfully said hello",
				zap.String("result", result.String()))
		}
	}(logger.GetLoggerFromContext(ctx).With(zap.String(logger.KeyFunctionName, "SayHello")))
	// Call the next service
	return ins.next.SayHello(ctx, request)
}

func (ins *logging) UseStandardResponse(ctx context.Context, request *pb.PayloadRequest) (result *common.StandardResponse, err error) {
	defer func(entry *zap.Logger) {
		if err != nil {
			entry.Error("Failed to use standard response",
				zap.String("name", request.GetName()),
				zap.Error(err))
		} else {
			entry.Debug("Successfully used standard response",
				zap.String("result", result.String()))
		}
	}(logger.GetLoggerFromContext(ctx).With(zap.String(logger.KeyFunctionName, "UseStandardResponse")))
	// Call the next service
	return ins.next.UseStandardResponse(ctx, request)
}
