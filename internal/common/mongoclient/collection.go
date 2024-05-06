// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package mongoclient

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Collection struct {
	name       string
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
	ctx        context.Context
}

func (d *Database) Collection(name string) *Collection {
	collection := d.database.Collection(name)
	return &Collection{
		name:       name,
		client:     d.client,
		database:   d.database,
		collection: collection,
		ctx:        d.ctx,
	}
}

func (c *Collection) Name() string {
	return c.name
}

func (c *Collection) Client() *MongoClient {
	return c.Database().Client()
}

func (c *Collection) Database() *Database {
	return &Database{
		name:     c.database.Name(),
		client:   c.client,
		database: c.database,
		ctx:      c.ctx,
	}
}

func (c *Collection) WithContext(ctx context.Context) *Collection {
	c.ctx = ctx
	return c
}

func (c *Collection) Exists() (bool, error) {
	names, err := c.database.ListCollectionNames(c.ctx, bson.M{"name": c.name})
	if err != nil {
		return false, err
	}

	return len(names) > 0, nil
}

func (c *Collection) EnsureExistance() error {
	// Check if the collection exists
	exists, err := c.Exists()
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	// Create the collection
	if err := c.database.CreateCollection(c.ctx, c.name); err != nil {
		return err
	}

	return nil
}

func (c *Collection) Drop() error {
	return c.collection.Drop(c.ctx)
}

func (c *Collection) IsEmpty() (bool, error) {
	count, err := c.collection.CountDocuments(c.ctx, bson.M{})
	if err != nil {
		return false, err
	}

	return count == 0, nil
}

func (c *Collection) FindById(id string) (interface{}, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var document interface{}
	if err := c.collection.FindOne(c.ctx, bson.M{"_id": oid}).Decode(&document); err != nil {
		return nil, err
	}
	return document, nil
}

func (c *Collection) InsertOne(document interface{}) (string, error) {
	res, err := c.collection.InsertOne(c.ctx, document)
	if err != nil {
		return "", err
	}
	oid := res.InsertedID.(primitive.ObjectID)
	return oid.Hex(), nil
}

func (c *Collection) UpdateByID(id string, update interface{}) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = c.collection.UpdateOne(c.ctx, bson.M{"_id": oid}, update)
	return err
}

func (c *Collection) DeleteByID(id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = c.collection.DeleteOne(c.ctx, bson.M{"_id": oid})
	return err
}
