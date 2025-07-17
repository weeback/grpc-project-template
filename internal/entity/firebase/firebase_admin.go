package firebase

import (
	"context"

	"firebase.google.com/go/v4/db"
)

type GoogleFirebaseAdmin interface {
	GenerateCustomToken(ctx context.Context, sessionId string) (string, error)
	GetDatabaseClient(ctx context.Context) (*db.Client, error)
}
