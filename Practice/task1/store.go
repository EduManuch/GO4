package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type MyDB struct {
	client *mongo.Client
	db     *mongo.Database
	coll   *mongo.Collection
}

var myDB MyDB

func InitDB(ctx context.Context) {

	opts := options.Client().ApplyURI("mongodb://localhost:27017")
	var err error
	myDB.client, err = mongo.Connect(opts)
	if err != nil {
		log.Fatal(err)
	}

	if err := myDB.client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	myDB.db = myDB.client.Database("books_store")
	myDB.coll = myDB.db.Collection("books")
}

func (m *MyDB) CreateBook(book *Book) (*mongo.InsertOneResult, error) {
	newBook := bson.D{
		{Key: "title", Value: book.Title},
		{Key: "author", Value: book.Author},
		{Key: "year", Value: book.Year},
		{Key: "country", Value: book.Country},
		{Key: "pages", Value: book.Pages},
	}

	r, err := m.coll.InsertOne(ctx, newBook)
	return r, err
}

func (m *MyDB) FindOneBook(id string) (*Book, error) {

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	result := m.coll.FindOne(ctx, bson.M{"_id": bson.M{"$eq": objectID}})

	var r bson.M
	err = result.Decode(&r)
	if err != nil {
		return nil, err
	}

	s := func(i any) string { return fmt.Sprintf("%v", i) }

	year, err := strconv.Atoi(s(r["year"]))
	if err != nil {
		return nil, err
	}
	pages, err := strconv.Atoi(s(r["pages"]))
	if err != nil {
		return nil, err
	}

	book := &Book{
		Title:   s(r["title"]),
		Author:  s(r["author"]),
		Year:    year,
		Country: s(r["country"]),
		Pages:   pages,
	}
	return book, nil
}

func (m *MyDB) UpdateOneBook(id string, book *Book) (*Book, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	sMap := make(map[string]any)
	js, err := json.Marshal(*book)
	fmt.Println(string(js))
	if err != nil {
		return nil, err
	}
	json.Unmarshal(js, &sMap)

	for k, v := range sMap {
		if v != "" && v != 0 {
			strVal := fmt.Sprintf("%v", v)
			_, err := m.coll.UpdateOne(
				ctx,
				bson.M{"_id": objectID},
				bson.D{
					{Key: "$set", Value: bson.M{k: strVal}},
				},
			)
			if err != nil {
				return nil, err
			}
		}
	}

	updBook, err := m.FindOneBook(id)
	if err != nil {
		return nil, err
	}
	return updBook, nil
}

func (m *MyDB) DeleteOneBook(id string) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = m.coll.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}
	return nil
}
