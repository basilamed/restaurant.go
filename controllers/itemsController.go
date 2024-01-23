package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"items/models"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var itemsCollection *mongo.Collection

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	connectionString := os.Getenv("MONGODB_CONNECTION_STRING")
	clientOptions := options.Client().ApplyURI(connectionString)

	client, error := mongo.Connect(context.TODO(), clientOptions)

	if error != nil {
		log.Fatal(error)
	}

	fmt.Println("Mongodb connection success")

	dbName := os.Getenv("DBNAME")
    colName := os.Getenv("COLNAME")

	itemsCollection = client.Database(dbName).Collection(colName)

	fmt.Println("Collection istance is ready")
}
//add new item
func AddItem(c *gin.Context) {
    var ctx, cancel = context.WithCancel(context.Background())
    defer cancel()

    var item models.Item

    if err := c.BindJSON(&item); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    _, insertErr := itemsCollection.InsertOne(ctx, item)

    if insertErr != nil {
		log.Println("Error inserting item:", insertErr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting item", "message": insertErr.Error()})
		return
	}

    c.JSON(http.StatusOK, item)
}
//get all items
func GetItems(c *gin.Context) {
	var ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	cursor, err := itemsCollection.Find(ctx, bson.M{})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while getting all items", "message": err.Error()})
		return
	}

	var items []models.Item = make([]models.Item, 0)

	for cursor.Next(ctx) {
		var item models.Item
		cursor.Decode(&item)
		items = append(items, item)
	}

	c.JSON(http.StatusOK, items)
}
//get item by id
func GetItemById(c *gin.Context) {
	var ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	id, _ := primitive.ObjectIDFromHex(c.Param("id"))

	var item models.Item

	err := itemsCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&item)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while getting item by id", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}
//update item by id
func UpdateItemById(c *gin.Context) {
    var ctx, cancel = context.WithCancel(context.Background())
    defer cancel()

    id, err := primitive.ObjectIDFromHex(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
        return
    }

    var item models.Item
    if err := c.BindJSON(&item); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    update := bson.M{
        "$set": item,
    }

    result, err := itemsCollection.UpdateOne(ctx, bson.M{"_id": id}, update)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while updating item by id", "message": err.Error()})
        return
    }

    if result.ModifiedCount == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
        return
    }

    c.JSON(http.StatusOK, item)
}
//delete item by id
func DeleteItemById(c *gin.Context) {
	var ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	id, _ := primitive.ObjectIDFromHex(c.Param("id"))

	var item models.Item

	err := itemsCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&item)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while getting item by id", "message": err.Error()})
		return
	}

	_, err = itemsCollection.DeleteOne(ctx, bson.M{"_id": id})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while deleting item by id", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}
//message consumer
func StartMessageConsumer() {
	fmt.Println("Message consumer started")

	connectionString := os.Getenv("RABBITMQ_CONNECTION_STRING")
	connection, err := amqp.Dial(connectionString)

	if err != nil {
		log.Fatalf("%s: %s", "Failed to connect to RabbitMQ", err)
	}

	defer connection.Close()

	channel, err := connection.Channel()

	if err != nil {
		log.Fatalf("%s: %s", "Failed to open a channel", err)
	}

	defer channel.Close()

	queue, err := channel.QueueDeclare(
        "items",
        false,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        log.Fatalf("Failed to declare a queue: %v", err)
    }

	messages, err := channel.Consume(
        queue.Name,
        "",
        true,
        false,
        false,
        false,
        nil,
    )

	if err != nil {
		log.Fatalf("%s: %s", "Failed to register a consumer", err)
	}

	for mesage := range messages {
        log.Printf("Received a message: %s", mesage.Body)

        var message models.Message
        if err := json.Unmarshal(mesage.Body, &message); err != nil {
            log.Printf("Error decoding message: %v", err)
            continue
        }
    }

	fmt.Println("Waiting for messages...")

}


