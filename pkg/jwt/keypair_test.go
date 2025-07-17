package jwt

import "testing"

func Test_GenerateKeyPair(t *testing.T) {

	kp, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair err: %v", err)
	}

	key, err := kp.PKCS8PrivateKey()
	if err != nil {
		t.Fatalf("GenerateKeyPair err: %v", err)
	}

	t.Logf("base64: %v", key)

}
