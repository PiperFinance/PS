package configs

import (
	"context"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"sync"
	"time"
)

var mongoBackgroundContextOnce sync.Once

var mongoClient *mongo.Client

//mongoUrl : Sample url format is mongodb://[username:password@]host1[:port1][,...hostN[:portN]][/[defaultauthdb][?options]]
func mongoUrl() string {
	mongoURL, ok := os.LookupEnv("MONGO_URL")
	if !ok {
		log.Errorf("Missing MONGO_URL env, defaulting to mongodb://localhost:27017")
		mongoURL = "mongodb://localhost:27017"
	}
	return mongoURL
}

func newMongoClient() *mongo.Client {
	_mongoClient, err := mongo.NewClient(options.Client().ApplyURI(mongoUrl()))
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	return _mongoClient
}

//GetMongo Singleton Approach to get mongo connection (Defaults to context.Background)
func GetMongo() *mongo.Client {
	mongoBackgroundContextOnce.Do(func() {
		mongoClient := newMongoClient()
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		err := mongoClient.Connect(ctx)
		if err != nil {
			log.Fatal(err)
		}
	})
	return mongoClient
}
