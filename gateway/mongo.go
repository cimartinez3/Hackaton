package gateway

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"kushki/hackaton/schemas"
	"log"
)

type IMongo interface {
	GetDocument(email string) ([]schemas.TokenDB, error)
	PutDocument(request interface{}) bool
	UpdateDocument(id, cvv string) bool
}

type Mongo struct {
	conn *mongo.Collection
}

func NewMongoService() IMongo {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://root:root@localhost:27017/"))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database("hackaton").Collection("cards")
	return &Mongo{conn: collection}
}

func (m *Mongo) GetDocument(email string) ([]schemas.TokenDB, error) {
	items := make([]schemas.TokenDB, 0)
	cur, err := m.conn.Find(context.Background(), bson.D{{"email", email}})

	if err != nil {
		fmt.Println(err)
		return items, err
	}

	for cur.Next(context.Background()) {
		data := &schemas.TokenDB{}

		err := cur.Decode(data)

		if err != nil {
			fmt.Println(err)
			continue
		}

		items = append(items, *data)
	}
	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {
			return
		}
	}(cur, context.Background())

	return items, nil
}

func (m *Mongo) PutDocument(request interface{}) bool {
	res, err := m.conn.InsertOne(context.Background(), request)

	if err != nil {
		fmt.Println("Error ", err)
		return false
	}

	fmt.Println("SUCCESS", res.InsertedID)

	return true
}

func (m *Mongo) UpdateDocument(id, cvv string) bool {
	res, err := m.conn.UpdateOne(context.Background(), bson.M{"id": id}, bson.D{{"$set", bson.M{"cvv": cvv}}})

	if err != nil {
		log.Println("ERROR NO ITEMS UPDATED", err)
		return false
	}

	if res.ModifiedCount == 0 {
		log.Println("ZERO ITEMS UPDATED")
		return false
	}

	log.Println("UPDATED :D")
	return true
}
