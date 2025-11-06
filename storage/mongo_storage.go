package storage

import (
	"context"
	"time"

	"com.bonkelbansi/go-kanban/internals/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoStore(uri string) (*MongoStore, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	db := client.Database("kanban")
	coll := db.Collection("tasks")

	return &MongoStore{client: client, coll: coll}, nil
}

func (m *MongoStore) LoadTasks() ([]models.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cur, err := m.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var ts []models.Task
	if err := cur.All(ctx, &ts); err != nil {
		return nil, err
	}
	return ts, nil
}

func (m *MongoStore) SaveTasks(ts []models.Task) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := m.coll.DeleteMany(ctx, bson.M{}); err != nil {
		return err
	}
	if len(ts) == 0 {
		return nil
	}
	docs := make([]interface{}, len(ts))
	for i, t := range ts {
		docs[i] = t
	}
	_, err := m.coll.InsertMany(ctx, docs)
	return err
}

func (m *MongoStore) ResetDemo(demo []models.Task) error {
	return m.SaveTasks(demo)
}

func (m *MongoStore) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = m.client.Disconnect(ctx)
}
