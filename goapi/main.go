package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Person struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
}

func main() {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://root:example@db:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	// Access the collection
	collection := client.Database("test").Collection("people")

	// Initialize the router
	router := mux.NewRouter()

	// Define the handler for the /people endpoint
	router.HandleFunc("/people", func(w http.ResponseWriter, r *http.Request) {
		// Define filter and find options
		filter := bson.D{{}}
		findOptions := options.Find()

		// Find documents
		cur, err := collection.Find(context.Background(), filter, findOptions)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer cur.Close(context.Background())

		// Loop through the cursor and add each document to the results slice
		var results []Person
		for cur.Next(context.Background()) {
			var person Person
			err := cur.Decode(&person)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			results = append(results, person)
		}

		// Encode the results slice as JSON and write it to the response
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(results)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", router))
}
