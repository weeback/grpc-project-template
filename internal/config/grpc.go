package config

var (
	// defaultGRPC is the default gRPC options.
	// It does not need to be changed (it is not recommended to change its value during the process).
	defaultGRPC = OptionGRPC{
		MaxRecvMsgSize: 10 << 20, // 1 << 10 = 1KB, 1 << 20 = 1MB, 10 << 20 = 10MB ...
		MaxSendMsgSize: 10 << 20, // 10MB
	}

	// sharedGRPC can be changed to store gRPC options for the entire application.
	// This value can be changed by functions from within this package.
	// Default value is the same as defaultGRPC.
	sharedGRPC = defaultGRPC
)

type OptionGRPC struct {
	MaxRecvMsgSize int
	MaxSendMsgSize int
}

func GetOptionGRPC() OptionGRPC {
	return sharedGRPC
}
