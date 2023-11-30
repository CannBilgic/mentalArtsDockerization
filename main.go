package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	Username string `json:"username"`
	ID       int    `json:"id"`
}

var collection *mongo.Collection

func main() {
	r := gin.Default()
	baseURL := os.Getenv("baseURL")

	// MongoDB bağlantısı kurma
	clientOptions := options.Client().ApplyURI("mongodb://root:example@mongo:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// MongoDB koleksiyonunu seçme
	collection = client.Database("users").Collection("deneme")

	// MongoDB'ye bağlanma denemesi
	for {
		_, err := collection.Find(context.Background(), bson.M{})
		if err == nil {
			break
		}
		log.Println("MongoDB'ye bağlanılamadı. Tekrar denenecek.")
	}

	r.GET(baseURL+"/ping", handlePing)
	r.GET(baseURL+"/getUsers", getUsers)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func handlePing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func getUsers(c *gin.Context) {
	var users []User

	// MongoDB'den veri çekme
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching users"})
		return
	}
	defer cursor.Close(context.Background())

	// Belgeyi User tipine dönüştürme
	for cursor.Next(context.Background()) {
		var user User
		if err := cursor.Decode(&user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding user"})
			return
		}
		users = append(users, user)
	}

	c.JSON(http.StatusOK, users)
}
