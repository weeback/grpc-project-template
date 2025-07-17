package jwt

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
)

func GenerateKeyPair() (*KeyPair, error) {
	// Generate ECDSA key pair
	pub, key, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate key pair: %w", err)
	}
	return &KeyPair{PublicKey: pub, PrivateKey: key}, nil
}

func KeyPairFromSecret(secret string) (*KeyPair, error) {
	der, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return nil, fmt.Errorf("failed to decode secret: %w", err)
	}

	key, err := x509.ParsePKCS8PrivateKey(der)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}
	privateKey, ok := key.(ed25519.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("invalid private key")
	}
	publicKey, ok := privateKey.Public().(ed25519.PublicKey)
	if !ok {
		return nil, fmt.Errorf("invalid public key")
	}
	return &KeyPair{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}, nil
}

type KeyPair struct {
	PublicKey  ed25519.PublicKey
	PrivateKey ed25519.PrivateKey
}

func (p *KeyPair) PKIXPublicKey() (string, error) {
	bin, err := x509.MarshalPKIXPublicKey(p.PublicKey)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bin), nil
}

func (p *KeyPair) PKCS8PrivateKey() (string, error) {
	bin, err := x509.MarshalPKCS8PrivateKey(p.PrivateKey)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bin), nil
}

func (p *KeyPair) PEM() (public []byte, private []byte, err error) {
	// PKIX marshal public key
	pkix, err := x509.MarshalPKIXPublicKey(p.PublicKey)
	if err != nil {
		return nil, nil, err
	}
	// PKCS8 marshal private key
	bin, err := x509.MarshalPKCS8PrivateKey(p.PrivateKey)
	if err != nil {
		return nil, nil, err
	}
	// Encode result
	o1 := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pkix,
	})
	o2 := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: bin,
	})
	return o1, o2, nil
}
