package config

import "time"

var (
	// defaultGRPC is the default gRPC options.
	// It does not need to be changed (it is not recommended to change its value during the process).
	defaultGRPC = OptionGRPC{
		MaxRecvMsgSize:       50 << 20,          // Increased to 50MB
		MaxSendMsgSize:       50 << 20,          // Increased to 50MB
		ConnectionTimeout:    120 * time.Second, // 2 minutes connection timeout
		StreamIdleTimeout:    300 * time.Second, // 5 minutes stream idle timeout
		KeepaliveTime:        60 * time.Second,  // Send keepalive ping every 60 seconds
		KeepaliveTimeout:     20 * time.Second,  // Wait 20 seconds for keepalive ping response
		MinConnectionTimeout: 120 * time.Second, // Minimum connection timeout of 2 minutes
	}

	// sharedGRPC can be changed to store gRPC options for the entire application.
	// This value can be changed by functions from within this package.
	// Default value is the same as defaultGRPC.
	sharedGRPC = defaultGRPC
)

type OptionGRPC struct {
	MaxRecvMsgSize       int
	MaxSendMsgSize       int
	ConnectionTimeout    time.Duration // Connection timeout in seconds
	StreamIdleTimeout    time.Duration // Stream idle timeout in seconds
	KeepaliveTime        time.Duration // Send keepalive ping every X seconds
	KeepaliveTimeout     time.Duration // Wait X seconds for keepalive ping response
	MinConnectionTimeout time.Duration // Minimum connection timeout in seconds
}

func GetOptionGRPC() OptionGRPC {
	return sharedGRPC
}
