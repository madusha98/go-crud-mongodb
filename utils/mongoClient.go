package mongoclient

import (
	"context"
	"errors"
	"fmt"
	"go-crud-mongodb/models"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbName = "go-example-db"
	usersCollection = "users")

func getClient() *mongo.Client{
	

	// uncomment to run with .env file
	// err := godotenv.Load()
	// if err != nil {
	// 	fmt.Print("Error loading .env file")
	// }

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGODB_URL")))

	if err != nil {
        fmt.Println(err)
    }

	return client
}

func getCollection(collectionName string) *mongo.Collection {
	client := getClient()
	db := client.Database(dbName)
	collection := db.Collection(collectionName)

	return collection
}

func InsertUser(user models.User)  (*primitive.ObjectID, error) {

	collection := getCollection(usersCollection)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := collection.InsertOne(ctx, user)

	if err == context.DeadlineExceeded {
		fmt.Println("deadline exceeded")
		return nil, err
	}

    if err != nil {
        fmt.Printf("Unable to execute the query. %v", err)
		ctx.Done()
		return nil, err
    }

    fmt.Printf("Inserted a single record %v", res)

    return &user.ID, nil
}

func GetUser(id primitive.ObjectID) (*models.User, error) {

	collection := getCollection(usersCollection)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user := models.User{}
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)

	if  err != nil {
        fmt.Println(err)
        return nil, err
    }

	return &user, nil

}

func GetAllUsers() (*[]models.User, error) {
	collection := getCollection(usersCollection)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	users := []models.User{}
	cursor, err := collection.Find(ctx, bson.M{})

	if  err != nil {
        fmt.Println(err)
        return nil, err
    }

	for cursor.Next(ctx) {
		u := models.User{}
        cursor.Decode(&u)
        users = append(users, u)
    }

	return &users, nil
}

func UpdateUser(id primitive.ObjectID, user models.User) (*primitive.ObjectID, error) {
	collection := getCollection(usersCollection)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	updateObject := bson.M{}

	fmt.Println(user)
	if user.Name != "" {
		updateObject["name"] = user.Name
	}
	if user.Location != "" {
		updateObject["location"] = user.Location
	}
	if user.Age != 0 {
		updateObject["age"] = user.Age
	}

	updateResult, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{
		"$set": updateObject,
	}, )

	if updateResult.MatchedCount < 1{
		return nil, errors.New("record not found")
	}

	if updateResult.ModifiedCount < 1{
		return nil, errors.New("similar data exists in db")
		
	}

    if err != nil {
        fmt.Printf("Unable to execute the query. %v", err)
		return nil, err
    }

    return &id, nil
}

func DeleteUser(id primitive.ObjectID) (*primitive.ObjectID, error) {
	collection := getCollection(usersCollection)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	deleteResult, err := collection.DeleteOne(ctx, bson.M{"_id": id})

	if deleteResult.DeletedCount < 1{
		return nil, errors.New("record not found")
	}

	if  err != nil {
        fmt.Println(err)
        return nil, err
    }

	return &id, nil
}