package mongodb

import (
	"context"
	"fmt"
	"github/weeback/grpc-project-template/internal/entity/db"
	"github/weeback/grpc-project-template/pkg/mongodb"
	"os"
)

type DB struct {
	ExampleDB db.ExampleDB
}

func NewMongoDB(ctx context.Context, withURI string) *DB {

	_, err := mongodb.NewConnection(ctx, withURI)
	if err != nil {
		fmt.Printf("mongodb.NewConnection err: %v\n", err)
		os.Exit(1)
	}
	// dbc := conn.Database()

	return &DB{
		// TODO: Add more repositories as needed
		// Example:
		// ExampleDB: NewExampleRepository(dbc),
	}
}
