package firebase

import (
	"context"
	"fmt"
	"sync"

	entity "github.com/weeback/grpc-project-template/internal/entity/firebase"

	"firebase.google.com/go/v4/auth"
	"firebase.google.com/go/v4/db"
	"google.golang.org/api/option"

	firebase "firebase.google.com/go/v4"
)

func Init(projectId string, databaseURL string, certificateJson []byte) entity.GoogleFirebaseAdmin {

	return &firebaseAdmin{
		projectId:   projectId,
		databaseURL: databaseURL,
		cert:        certificateJson,
	}
}

type firebaseAdmin struct {
	projectId        string
	databaseURL      string
	cert             []byte
	authOnce, dbOnce sync.Once
	auth             *auth.Client
	database         *db.Client
}

func (admin *firebaseAdmin) GetAuthClient(ctx context.Context) (*auth.Client, error) {
	var (
		opt     []option.ClientOption
		initErr error
	)
	admin.authOnce.Do(func() {
		if len(admin.cert) > 0 {
			opt = []option.ClientOption{option.WithCredentialsJSON(admin.cert)}
		}
		app, err := firebase.NewApp(ctx, &firebase.Config{ProjectID: admin.projectId}, opt...)
		if err != nil {
			initErr = err
			return
		}
		admin.auth, initErr = app.Auth(ctx)
	})
	if initErr != nil {
		//
		fmt.Printf("Failed to initialize Firebase app: %v\n", initErr)
		// reset the once to allow retrying
		admin.authOnce = sync.Once{} // Reset the once if there was an error
	}
	return admin.auth, initErr
}

func (admin *firebaseAdmin) GetDatabaseClient(ctx context.Context) (*db.Client, error) {
	var (
		opt     []option.ClientOption
		initErr error
	)
	admin.dbOnce.Do(func() {
		if len(admin.cert) > 0 {
			opt = []option.ClientOption{option.WithCredentialsJSON(admin.cert)}
		}
		app, err := firebase.NewApp(ctx, &firebase.Config{ProjectID: admin.projectId, DatabaseURL: admin.databaseURL}, opt...)
		if err != nil {

			initErr = err
			return
		}
		admin.database, initErr = app.Database(ctx)
	})
	if initErr != nil {
		//
		fmt.Printf("Failed to initialize Firebase app: %v\n", initErr)
		// reset the once to allow retrying
		admin.dbOnce = sync.Once{} // Reset the once if there was an error
	}
	return admin.database, initErr
}

func (admin *firebaseAdmin) GenerateCustomToken(ctx context.Context, sessionId string) (string, error) {
	client, err := admin.GetAuthClient(ctx)
	if err != nil {
		return "", err
	}
	// UID đại diện cho phiên, có thể thay đổi theo logics
	uid := fmt.Sprintf("session_%s", sessionId)

	// Thêm các custom claim nếu muốn
	claims := map[string]interface{}{
		"session": sessionId,
		"role":    "web_client", // tuỳ biến
	}

	token, err := client.CustomTokenWithClaims(ctx, uid, claims)
	if err != nil {
		return "", err
	}

	return token, nil
}
