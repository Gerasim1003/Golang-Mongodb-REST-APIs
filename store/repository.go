package store

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

//Repository ..
type Repository struct {
	Client *mongo.Client
}

//SERVER ...
const SERVER string = "mongodb://localhost:8080"

//DATABASE ...
const DATABASE string = "civilact"

//COLLECTION ...
const COLLECTION string = "heroes"

//NewRepository ...
func NewRepository() Repository {
	log.Println("New repo")
	return Repository{
		Client: GetClient(),
	}
}

//GetClient ...
func GetClient() *mongo.Client {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:8080")

	client, err := mongo.NewClient(clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	return client
}

//ReturAllHeroes ..
func (r Repository) ReturAllHeroes(filter bson.M) []*Hero {
	var heroes []*Hero
	collection := r.Client.Database(DATABASE).Collection(COLLECTION)
	cur, err := collection.Find(context.TODO(), filter)

	if err != nil {
		log.Fatal("Error on Finding all the documents", err)
	}

	for cur.Next(context.TODO()) {
		var hero Hero
		err := cur.Decode(&hero)

		if err != nil {
			log.Fatal("Error on Decoding the document", err)
		}
		heroes = append(heroes, &hero)
	}

	return heroes

}

//ReturnOneHero ..
func (r Repository) ReturnOneHero(filter bson.M) Hero {
	var hero Hero
	collection := r.Client.Database(DATABASE).Collection(COLLECTION)
	documentReturned := collection.FindOne(context.TODO(), filter)
	documentReturned.Decode(&hero)
	return hero
}

//InsertNewHero ..
func (r Repository) InsertNewHero(hero Hero) interface{} {
	collection := r.Client.Database(DATABASE).Collection(COLLECTION)
	insertResult, err := collection.InsertOne(context.TODO(), hero)
	if err != nil {
		log.Fatal("")
	}

	return insertResult.InsertedID
}

//RemoveOneHero ..
func (r Repository) RemoveOneHero(filter bson.M) int64 {
	collection := r.Client.Database(DATABASE).Collection(COLLECTION)

	deleteRes, err := collection.DeleteOne(context.TODO(), filter)

	if err != nil {
		log.Fatal("Error on deleting one Hero", err)
	}

	return deleteRes.DeletedCount
}

//UpdateHero ..
func (r Repository) UpdateHero(updatedData bson.M, filter bson.M) int64 {
	collection := r.Client.Database(DATABASE).Collection(COLLECTION)
	atualizacao := bson.D{{Key: "$set", Value: updatedData}}
	updatedResult, err := collection.UpdateOne(context.TODO(), filter, atualizacao)

	if err != nil {
		log.Fatal("Error on updating one Hero", err)
	}

	return updatedResult.ModifiedCount
}
