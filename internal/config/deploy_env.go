package config

import "os"

const (
	// deploymentEnvironment is the default deployment environment.
	deploymentEnvironment = Production

	key                            = "MC4CAQAwBQYDK2VwBCIEIKa3OWQKuYGcXLm2Wgbj9U14UlRzfX0zdPhLJl+VM5GT"
	cloudflareApiKey               = "unset"
	cloudflareTurnstileCredentials = "unset"
	firebaseSdkCredentials         = "firebase-adminsdk.json"
	firebaseDatabaseURL            = "https://${project-id}-default-rtdb.${region}.firebasedatabase.app/"
	mongoURI                       = "mongodb://localhost:27017/firebase?retryWrites=true&w=majority"

	Production  Environment = "production"
	Development Environment = "development"
)

type Environment string

func GetDeploymentEnvironment() Environment {
	if val := os.Getenv("DEPLOYMENT_ENVIRONMENT"); val != "" {
		switch val {
		case "development", "dev":
			return Development
		case "production", "prod":
			return Production
		default:
			return deploymentEnvironment
		}
	}
	return deploymentEnvironment
}

func GetCloudflareApiKey() string {
	if val := os.Getenv("CLOUDFLARE_API_KEY"); val != "" {
		return val
	}
	return cloudflareApiKey
}

func GetCloudflareTurnstileCredentials() string {
	if val := os.Getenv("CLOUDFLARE_TURNSTILE_CREDENTIALS"); val != "" {
		return val
	}
	return cloudflareTurnstileCredentials
}

func GetRealtimeDatabaseURL() string {
	if val := os.Getenv("FIREBASE_DATABASE_URL"); val != "" {
		return val
	}
	return firebaseDatabaseURL // Replace with your default URL or leave empty
}

func GetFirebaseSdkCredentials() string {
	if val := os.Getenv("FIREBASE_SDK_CREDENTIALS"); val != "" {
		return val
	}
	return firebaseSdkCredentials
}

func GetMongoURI() string {
	if val := os.Getenv("MONGO_URI"); val != "" {
		return val
	}
	return mongoURI
}

func GetPreSharedKey() string {
	if val := os.Getenv("JWT_PRE-SHARED_KEY"); val != "" {
		return val
	}
	return key
}
