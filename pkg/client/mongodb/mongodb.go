package mongodb

import (
	"context"
	"fmt"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func NewClient(ctx context.Context,
	host, port, username, password, database, authDB string,
) (*mongo.Database, error) {
	err := sync.OnceValue(func() (err error) {
		var mongoDBURL string

		if username == "" && password == "" {
			mongoDBURL = fmt.Sprintf("mongodb://%s:%s", host, port)
		} else {
			mongoDBURL = fmt.Sprintf("mongodb://%s:%s@%s:%s", username, password, host, port)
		}

		clientOptions := options.Client().ApplyURI(mongoDBURL)
		client, err = mongo.Connect(ctx, clientOptions)
		return err
	})()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongoDB: %w", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping mongoDB: %w", err)
	}

	return client.Database(database), nil
}
