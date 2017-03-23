package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	restful "github.com/emicklei/go-restful"
	"github.com/stevenceuppens/gorest-example/pkg/books"
	"gopkg.in/mgo.v2"
)

// ENV Defaults
var (
	Port       = ":5000"
	MongoDBURL = "mongodb://127.0.0.8/gorest"
)

// Capture ENV
func init() {
	if p := os.Getenv("PORT"); p != "" {
		Port = ":" + p
	}
	if m := os.Getenv("MongoDBURL"); m != "" {
		MongoDBURL = m
	}
}

func main() {

	/**
	* REST endpoints use APIS !
	* APIS use Databases or other APIS!
	* APIS can be mocked via interface
	 */

	/**
	* First define API's
	 */

	// choose API
	var bookAPI books.BookAPI
	session, err := mgo.Dial(MongoDBURL)
	if err == nil {
		// Use Mongo backend
		bookAPI = books.NewBookAPIMongo(session)
		fmt.Println("Using Mongo Backend")
	} else {
		// if fails, use memory backend
		bookAPI = books.NewBookAPIMemory()
		fmt.Println("Using Memory Backend")
	}

	/**
	* Second setup REST (and provide dependent APIs)
	 */
	// regiter REST endpoints
	wsContainer := restful.NewContainer()
	bookREST := books.NewBookREST(bookAPI)
	bookREST.Register(wsContainer)

	// Launch Server
	log.Printf("start listening on localhost %v", Port)
	server := &http.Server{Addr: Port, Handler: wsContainer}
	log.Fatal(server.ListenAndServe())
}
