package mongodb

import (
	"context"
	"log"
	"net/url"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"mongodb",
)

type MongoDBConfig interface {
	GetMongoDBUser() string
	GetMongoDBPassword() string
	GetMongoDBAddress() string
	GetMongoDBName() string
	GetMongoAuthDBName() string
	GetMongoTLS() bool
	GetMongoTLSCa() string
}

func InitializeMongoDBConnection(config MongoDBConfig) *mongo.Database {
	u := &url.URL{
		Scheme: "mongodb",
		User:   url.UserPassword(config.GetMongoDBUser(), config.GetMongoDBPassword()),
		Host:   config.GetMongoDBAddress(),
		Path:   config.GetMongoDBName(),
	}

	q := url.Values{}

	if config.GetMongoAuthDBName() != "" {
		q.Set("authSource", config.GetMongoAuthDBName())
	}
	if config.GetMongoTLS() {
		q.Set("tls", "true")
		q.Set("tlsCAFile", config.GetMongoTLSCa())
	}
	u.RawQuery = q.Encode()
	uri := u.String()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("failed to ping MongoDB: %v", err)
	}

	log.Printf("✅ Connected to MongoDB at %s", config.GetMongoDBAddress())
	return client.Database(config.GetMongoDBName())
}
