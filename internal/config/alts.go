package config

var sharedALTS = OptionALTS{}

type OptionALTS struct {
	// TargetServiceAccounts contains a list of expected target service
	// accounts.
	TargetServiceAccounts []string
	// HandshakerServiceAddress represents the ALTS handshaker gRPC service
	// address to connect to.
	HandshakerServiceAddress string

	// Add ALTS specific configuration fields if needed
}

func LoadALTS() (*OptionALTS, error) {
	sharedALTS = OptionALTS{
		TargetServiceAccounts:    GetTargetServiceAccounts(),
		HandshakerServiceAddress: GetHandshakerServiceAddress(),
	}
	return &sharedALTS, nil
}

// GetOptionALTS returns the ALTS options.
func GetOptionALTS() OptionALTS {
	return sharedALTS
}
