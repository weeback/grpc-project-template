package errors

import (
	"fmt"

	commonpb "github.com/weeback/grpc-project-template/pb/common"
)

func Errorf(code int, format string, a ...any) Error {
	return Error{
		code:    ErrorCode(code),
		message: fmt.Sprintf(format, a...),
	}
}

func New(code ErrorCode, message string) Error {
	return Error{
		code:    code,
		message: message,
	}
}

type Error struct {
	code    ErrorCode
	message string
}

func (err Error) Error() string {
	return fmt.Sprintf("code (%d) - %s", err.code, err.message)
}

func (err Error) ToStandardResponse() *commonpb.StandardResponse {
	return &commonpb.StandardResponse{
		Code:    int32(err.code),
		Message: err.message,
	}
}

type ErrorCode int

// Common Error Code (0 - 99)
const (
	UnmarshalFailedCode ErrorCode = iota + 1 // 1
	ToMapFailedCode                          // 2
)

// Realtime Database Error Code (100 - 199)
const (
	RealtimeDBSetFailedCode    ErrorCode = iota + 100 // 100
	RealtimeDBUpdateFailedCode                        //101
)

var (
	ToMapFailed = Error{
		code:    ToMapFailedCode,
		message: "struct to map failed",
	}

	RealtimeDBSetFailed = Error{
		code:    RealtimeDBSetFailedCode,
		message: "realtimedb: set failed",
	}
	RealtimeDBUpdateFailed = Error{
		code:    RealtimeDBUpdateFailedCode,
		message: "realtimedb: update failed",
	}
)
