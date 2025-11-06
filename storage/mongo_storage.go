package storage

import (
	"context"
	"time"

	"com.bonkelbansi/go-kanban/internals/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStorage struct {
	Collection *mongo.Collection
}

func NewMongoStorage(uri, dbName, collName string) (*MongoStorage, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		return nil, err
	}

	coll := client.Database(dbName).Collection(collName)
	return &MongoStorage{Collection: coll}, nil
}

func (m *MongoStorage) LoadTasks() ([]models.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cur, err := m.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var tasks []models.Task
	if err := cur.All(ctx, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (m *MongoStorage) SaveTasks(tasks []models.Task) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// alles raus, dann vollständig neu schreiben
	if _, err := m.Collection.DeleteMany(ctx, bson.M{}); err != nil {
		return err
	}

	if len(tasks) == 0 {
		return nil
	}

	var docs []interface{}
	for _, t := range tasks {
		docs = append(docs, t)
	}

	_, err := m.Collection.InsertMany(ctx, docs)
	return err
}
