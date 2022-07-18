package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"errors"

	"strconv"
)

type book struct{
	ID 		 int	`json:"id"`
	Title 	string	`json:"title"`
	Author 	string	`json:"author"`
	Quantity int	`json:"quantity"`
}

var books = []book{
	{ID: 1, Title: "test", Author: "unknown", Quantity: 2},
	{ID: 2, Title: "test", Author: "unknown", Quantity: 2},
	{ID: 3, Title: "test", Author: "unknown", Quantity: 2},
}

func getBookById(id int) (*book, error){
	for i,b := range books{
		if(b.ID == id){
			return &books[i], nil
		}
	}

	return nil, errors.New("no book found")
}

func getBooks(router *gin.Context){
	router.IndentedJSON(http.StatusOK, books)
}

func createBook(router *gin.Context){
	var newBook book
	if err := router.BindJSON(&newBook); err != nil{
		return
	}

	books = append(books, newBook)
	router.IndentedJSON(http.StatusCreated, newBook)
}

func getSingleBook(router *gin.Context){
	id, convErr := strconv.Atoi(router.Param("id"))
	if(convErr != nil){
		router.IndentedJSON(http.StatusBadRequest, "Invalid id")
		return
	}
	book, err := getBookById(id)

	if(err != nil){
		router.IndentedJSON(http.StatusNotFound, "Book not found")
		return
	}

	router.IndentedJSON(http.StatusOK, book)
}

func main(){
	router := gin.Default()

	router.GET("/books", getBooks)
	router.POST("/books", createBook)
	router.GET("/book/:id", getSingleBook)
	router.Run("localhost:5000")
}