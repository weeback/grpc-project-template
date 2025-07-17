package hello

import (
	"fmt"

	pb "github/weeback/grpc-project-template/pb/hello"
)

func validateSayRequest(request *pb.HelloRequest) error {
	if request.GetName() == "" {
		return fmt.Errorf("name is required")
	}
	return nil
}
