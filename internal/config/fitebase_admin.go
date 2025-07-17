package config

import (
	"encoding/json"
	"fmt"
	"os"
)

var (
	sharedFirebaseAdmin = OptionFirebaseAdmin{}
)

type ServiceAccount struct {
	Type                    string `json:"type"`
	ProjectId               string `json:"project_id"`
	PrivateKeyId            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientId                string `json:"client_id"`
	AuthUri                 string `json:"auth_uri"`
	TokenUri                string `json:"token_uri"`
	AuthProviderX509CertUrl string `json:"auth_provider_x509_cert_url"`
	ClientX509CertUrl       string `json:"client_x509_cert_url"`
	UniverseDomain          string `json:"universe_domain"`
}

type OptionFirebaseAdmin struct {
	ProjectId   string // Firebase project ID
	DatabaseURL string
	// Firebase service account certificate in JSON format
	CertificateJson []byte // Firebase service account certificate in JSON format
}

func LoadWithJsonFile(path string) (*OptionFirebaseAdmin, error) {
	var (
		certificateJson []byte
		temp            ServiceAccount
	)
	certificateJson, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(certificateJson, &temp); err != nil {
		return nil, fmt.Errorf("file data is not a valid JSON: %w", err)
	}
	sharedFirebaseAdmin = OptionFirebaseAdmin{
		ProjectId:       temp.ProjectId,
		DatabaseURL:     GetRealtimeDatabaseURL(), // Assuming DatabaseURL is the AuthUri, adjust as needed
		CertificateJson: certificateJson,
	}
	return &sharedFirebaseAdmin, nil
}

// GetOptionFirebaseAdmin returns the Firebase Admin options.
func GetOptionFirebaseAdmin() OptionFirebaseAdmin {
	return sharedFirebaseAdmin
}
