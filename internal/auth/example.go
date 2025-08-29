package auth

type Inter interface {
	AuthFunc(fullMethod string, bodyHash string, jwtStr string) error
}

func New() Inter {
	return &ins{}
}

type ins struct{}

func (i *ins) AuthFunc(fullMethod string, bodyHash string, jwtStr string) error {
	// Implement your authorization logic here
	// You can use the fullMethod, bodyHash, and jwtStr parameters as needed

	switch fullMethod {
	case "/hello.HelloService/SayHello":
		// Check if the request is authorized

		return nil
	default:
		return nil
	}
}
