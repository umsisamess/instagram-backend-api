package main

import (
	"context"
	"fmt"
	"log"
	"encoding/json"
	"strings"
	"net/http"
	"time"
	
	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type users struct {
	ID			primitive.ObjectID		`json:"_id,omitempty" bson:"_id,omitempty"`
    Name 		string					`json:"name,omitempty" bson:"name,omitempty"`
    Email  		string					`json:"email,omitempty" bson:"email,omitempty"`
    Password 	string					`json:"password,omitempty" bson:"password,omitempty"`
}

type posts struct {
	ID					primitive.ObjectID		`json:"_id,omitempty" bson:"_id,omitempty"`
	Userid				primitive.ObjectID		`json:"userid,omitempty" bson:"userid,omitempty"`
    Caption 			string					`json:"caption,omitempty" bson:"name,omitempty"`
    Image_url  			string					`json:"image_url,omitempty" bson:"image_url,omitempty"`
    Posted_timestamp 	string					`json:"posted_timestamp,omitempty" bson:"posted_timestamp,omitempty"`
}


func createUser(w http.ResponseWriter,r *http.Request) {
	if r.URL.Path != "/users" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

// Connect to MongoDB
client, err := mongo.Connect(context.TODO(), clientOptions)

if err != nil{
    log.Fatal(err)
}

// Check the connection
err = client.Ping(context.TODO(), nil)

if err != nil {
    log.Fatal(err)
}

fmt.Println("Connected to MongoDB!")

	switch r.Method {
	case "GET":		
		 http.ServeFile(w, r, "form.html")
	case "POST":
		var user users
		_ = json.NewDecoder(r.Body).Decode(&user)
		collection := client.Database("backend-api").Collection("users")
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		result, _ := collection.InsertOne(ctx, user)
		json.NewEncoder(w).Encode(result)


		
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}

	err = client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}


func createPost(w http.ResponseWriter,r *http.Request) {
	if r.URL.Path != "/posts" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

// Connect to MongoDB
client, err := mongo.Connect(context.TODO(), clientOptions)

if err != nil{
    log.Fatal(err)
}

// Check the connection
err = client.Ping(context.TODO(), nil)

if err != nil {
    log.Fatal(err)
}

fmt.Println("Connected to MongoDB!")

	switch r.Method {
	case "POST":
		var post posts
		_ = json.NewDecoder(r.Body).Decode(&post)
		collection := client.Database("backend-api").Collection("posts")
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		result, _ := collection.InsertOne(ctx, post)
		json.NewEncoder(w).Encode(result)
		


		
	default:
		fmt.Fprintf(w, "Sorry, only POST methods are supported.")
	}

	err = client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}


func getUserorPostsbyId(w http.ResponseWriter,r *http.Request) {
	p := strings.Split(r.URL.Path, "/");
	id := p[2];
	cont := p[1]

	if cont=="users"{
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
		panic(err)
		}

		fmt.Printf(objID.Hex())

		clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

		// Connect to MongoDB
		client, err := mongo.Connect(context.TODO(), clientOptions)

		if err != nil{
			log.Fatal(err)
		}

		// Check the connection
		err = client.Ping(context.TODO(), nil)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Connected to MongoDB!")
		var user users
		collection := client.Database("backend-api").Collection("users")
		ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		er := collection.FindOne(ctx, user).Decode(&user)
		if er != nil {
			log.Fatal(er)
			return
		}
		json.NewEncoder(w).Encode(user.ID)
		json.NewEncoder(w).Encode(user.Name)
		json.NewEncoder(w).Encode(user.Email)

	}else if cont=="posts"{
		if p[2]=="users"{
			// useridd := p[3]
			// objIDD, err := primitive.ObjectIDFromHex(useridd)
			// if err != nil {
			// 	panic(err)
			// }
			clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

			// Connect to MongoDB
			client, err := mongo.Connect(context.TODO(), clientOptions)

			if err != nil{
				log.Fatal(err)
			}

			// Check the connection
			err = client.Ping(context.TODO(), nil)

			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("Connected to MongoDB!")

			var post posts
			collection := client.Database("backend-api").Collection("posts")
			ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
			er := collection.FindOne(ctx, post).Decode(&post)
			if er != nil {
				log.Fatal(er)
				return
			}
			json.NewEncoder(w).Encode(post)

		}else{
			objID, err := primitive.ObjectIDFromHex(id)
			if err != nil {
			panic(err)
			}

			fmt.Printf(objID.Hex())

			clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

			// Connect to MongoDB
			client, err := mongo.Connect(context.TODO(), clientOptions)

			if err != nil{
				log.Fatal(err)
			}

			// Check the connection
			err = client.Ping(context.TODO(), nil)

			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("Connected to MongoDB!")
			var post posts
			collection := client.Database("backend-api").Collection("posts")
			ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
			er := collection.FindOne(ctx, post).Decode(&post)
			if er != nil {
				log.Fatal(er)
				return
			}
			json.NewEncoder(w).Encode(post)
			
		}
		
	}

	
	


}

func main(){

	http.HandleFunc("/", getUserorPostsbyId)
	http.HandleFunc(`/users`, createUser)
	http.HandleFunc(`/posts`, createPost)
	

	fmt.Printf("Starting server for testing HTTP POST...\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}