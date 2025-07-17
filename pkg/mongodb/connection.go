package mongodb

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

const (
	connectTimeout = 10 * time.Second

	defaultURI = "mongodb://localhost:27017/dbnameDB?authSource=admin"
)

func getAppName() string {
	if name := os.Getenv("APP_NAME"); name != "" {
		return name
	}
	return "Go-Service"
}

func getDbName(uri string) (dbName string) {
	var queryStr string
	if idx := strings.LastIndex(uri, "://"); idx != -1 {
		dbName = uri[idx+3:]
	}
	if idx := strings.LastIndexByte(dbName, '/'); idx != -1 {
		dbName = dbName[idx+1:]
	} else {
		// must have a / separator between hosts and path
		return ""
	}
	if idx := strings.IndexRune(dbName, '?'); idx != -1 {
		queryStr = dbName[idx+1:]
		dbName = dbName[:idx]
	}
	if dbName == "" && queryStr != "" {
		q, _ := url.ParseQuery(queryStr)
		return q.Get("authSource")
	}
	return dbName
}

func NewConnection(ctx context.Context, withURI ...string) (*Connection, error) {

	opt := options.Client().
		ApplyURI(defaultURI).
		SetTimeout(connectTimeout).
		SetAppName(getAppName())

	if len(withURI) > 0 {
		opt = opt.ApplyURI(withURI[0])
	}

	replicaSet := opt.ReplicaSet != nil && *opt.ReplicaSet != ""
	if replicaSet {
		opt.SetReadPreference(readpref.SecondaryPreferred())
	}

	client, err := mongo.Connect(opt)
	if err != nil {
		return nil, err
	}

	var (
		primaryStatus, secondaryStatus error
	)
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		primaryStatus = err
		log.Printf("Primary ping failed: %v", err)
	}
	if err := client.Ping(ctx, readpref.Secondary()); err != nil {
		secondaryStatus = err
		log.Printf("Secondary ping failed: %v", err)
	}

	if primaryStatus != nil && secondaryStatus != nil {
		return nil, fmt.Errorf("unable to connect to primary or secondary, primary: %v, secondary: %v", primaryStatus, secondaryStatus)
	}

	return &Connection{
		client:     client,
		replicaSet: replicaSet,
		dbName:     getDbName(opt.GetURI()),
	}, nil
}

type Connection struct {
	client     *mongo.Client
	replicaSet bool
	dbName     string
}

func (c *Connection) WithDatabase(dbName string) error {
	if c.client != nil {
		c.dbName = dbName
		return nil
	}
	return fmt.Errorf("connection is unavailable")
}

func (c *Connection) Database() *mongo.Database {
	return c.client.Database(c.dbName)
}

func (c *Connection) Read(ctx context.Context, readFn func(*mongo.Database) error) error {
	// by default with readpref.Secondary
	if err := c.ReadSecondary(ctx, readFn); err != nil {
		log.Printf("Secondary node is unavailable, try reading with mode readpref.Primary")
		return c.ReadPrimary(ctx, readFn)
	}
	return nil
}

func (c *Connection) ReadPrimary(ctx context.Context, readFn func(*mongo.Database) error) error {
	if err := c.client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	// execute readFn by connect  connect readpref.Primary()
	return readFn(c.client.Database(c.dbName,
		options.Database().SetReadPreference(readpref.Primary())))
}

func (c *Connection) ReadSecondary(ctx context.Context, readFn func(*mongo.Database) error) error {
	if err := c.client.Ping(ctx, readpref.Secondary()); err != nil {
		return err
	}
	// execute readFn by connect readpref.Secondary()
	return readFn(c.client.Database(c.dbName,
		options.Database().SetReadPreference(readpref.Secondary())))
}

func (c *Connection) Write(ctx context.Context, writeFn func(*mongo.Database) error) error {
	// Is write operation, so we need to connect to primary node
	if err := c.client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	// execute readFn by connect readpref.Primary()
	return writeFn(c.client.Database(c.dbName,
		options.Database().SetReadPreference(readpref.Primary())))
}
