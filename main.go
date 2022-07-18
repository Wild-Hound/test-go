package main

import (
	"context"
	// "errors"
	"fmt"
	"net/http"
	"os"

	// "strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type book struct{
	ID 	    primitive.ObjectID	`json:"_id,omitempty", bson:"_id",omitempty`
	Title 	string	`json:"title,omitempty", bson:"title,omitempty"`
	Author 	string	`json:"author,omitempty", bson:"author,omitempty"`
	Quantity int	`json:"quantity,omitempty", bson:"quantity,omitempty"`
	BookId  int 	`json:"book_id", bson:"book_id"`
}
var collection *mongo.Collection

var books = []book{
	{Title: "test", Author: "unknown", Quantity: 2},
	{Title: "test", Author: "unknown", Quantity: 2},
	{Title: "test", Author: "unknown", Quantity: 2},
}

func init(){
err := godotenv.Load(".env")
if err != nil{
	fmt.Println("error loading .env")
	return 
}
mongoUsername := os.Getenv("mongo_username")
mongoPassword := os.Getenv("mongo_password")
mongoConnectionLink := "mongodb+srv://"+mongoUsername+":"+mongoPassword+"@cluster0.2fhtx.mongodb.net/?retryWrites=true&w=majority"

fmt.Println("mongo connection link")
fmt.Println(mongoConnectionLink)

ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
client, mongoErr := mongo.Connect(ctx, options.Client().ApplyURI(mongoConnectionLink))
if(mongoErr != nil){
	fmt.Println("error connecting to mongo")
	fmt.Println(mongoErr)
	return
}
collection = client.Database(os.Getenv("mongo_db_name")).Collection(os.Getenv("mongo_collection_name"))
fmt.Println("mounted DB")
}

// func getBookById(id int) (*book, error){
// 	for i,b := range books{
// 		if(b.ID == id){
// 			return &books[i], nil
// 		}
// 	}

// 	return nil, errors.New("no book found")
// }

func getBooks(router *gin.Context){
	var mongoBooks []book

	ctx, _:= context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if(err != nil){
		router.IndentedJSON(http.StatusInternalServerError, "error while getting data")
		return
	}
	for cursor.Next(ctx){
		var mongoBook book
		cursor.Decode(&mongoBook)
		mongoBooks = append(mongoBooks, mongoBook)
	}

	router.IndentedJSON(http.StatusOK, mongoBooks)
}

func createBook(router *gin.Context){
	var newBook book
	if err := router.BindJSON(&newBook); err != nil{
		return
	}

	books = append(books, newBook)
	router.IndentedJSON(http.StatusCreated, newBook)
}

// func getSingleBook(router *gin.Context){
// 	id, convErr := strconv.Atoi(router.Param("id"))
// 	if(convErr != nil){
// 		router.IndentedJSON(http.StatusBadRequest, "Invalid id")
// 		return
// 	}
// 	book, err := getBookById(id)

// 	if(err != nil){
// 		router.IndentedJSON(http.StatusNotFound, "Book not found")
// 		return
// 	}

// 	router.IndentedJSON(http.StatusOK, book)
// }

// func checkoutBook(router *gin.Context){
// 	srtId, paramExists := router.GetQuery("id")
// 	id, err := strconv.Atoi(srtId)

// 	if(err != nil || !paramExists){
// 		router.IndentedJSON(http.StatusBadRequest,"Invalid id")
// 		return
// 	}

// 	book,bookErr := getBookById(id)
// 	if(bookErr != nil){
// 		router.IndentedJSON(http.StatusNotFound,"book not found")
// 		return
// 	}

// 	if(book.Quantity <= 0){
// 		router.IndentedJSON(http.StatusNotFound,"book is no longer available")
// 		return
// 	}

// 	book.Quantity -= 1
// 	router.IndentedJSON(http.StatusOK, book)
// }

func main(){
	router := gin.Default()

	router.GET("/books", getBooks)
	router.POST("/books", createBook)
	// router.GET("/book/:id", getSingleBook)
	// router.GET("/checkout", checkoutBook)
	router.Run("localhost:5000")
}