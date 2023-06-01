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
	GetDocument(email string) bool
	PutDocument(request schemas.TokenRequest) bool
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

func (m *Mongo) GetDocument(email string) bool {
	items := make([]schemas.TokenRequest, 0)
	cur, err := m.conn.Find(context.Background(), bson.D{{"client.email", email}})

	if err != nil {
		fmt.Println(err)
		return false
	}

	for cur.Next(context.Background()) {
		data := &schemas.TokenRequest{}

		err := cur.Decode(data)

		if err != nil {
			fmt.Println(err)
			continue
		}

		items = append(items, *data)
	}

	for _, a := range items {
		fmt.Println(a)
	}

	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {

		}
	}(cur, context.Background())

	return true
}

func (m *Mongo) PutDocument(request schemas.TokenRequest) bool {
	res, err := m.conn.InsertOne(context.Background(), request)

	if err != nil {
		fmt.Println("Error ", err)
		return false
	}

	fmt.Println("SUCCESS", res.InsertedID)

	return true
}
