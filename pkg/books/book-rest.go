package books

import (
	"net/http"

	"gopkg.in/mgo.v2/bson"

	restful "github.com/emicklei/go-restful"
	"github.com/stevenceuppens/gorest-example/pkg/utils"
)

// BookREST ...
type BookREST struct {
	api BookAPI
}

// NewBookREST creates a new BookREST object and passes the api object
func NewBookREST(api BookAPI) *BookREST {
	return &BookREST{
		api,
	}
}

// Register regitsers the REST endpoints for the Book resource
func (b *BookREST) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.
		Path("/api/v1/books").
		Doc("Manage Books").
		Consumes(restful.MIME_JSON, restful.MIME_XML).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	ws.Route(ws.POST("").To(b.createOne).
		Doc("Create a new book").
		Operation("createOne").
		Reads(NewBook{}).
		Writes(Book{}).
		Returns(400, "Invalid Request", nil))

	ws.Route(ws.GET("").To(b.findAll).
		Doc("Find all books").
		Operation("findAll").
		Writes([]Book{}))

	ws.Route(ws.GET("/{id}").To(b.findOneByID).
		Doc("Find a books by ID").
		Param(ws.PathParameter("id", "identifier of the book").DataType("string")).
		Operation("findOneByID").
		Writes(Book{}).
		Returns(400, "Invalid Request", nil).
		Returns(404, "Book not found", nil))

	ws.Route(ws.DELETE("/{id}").To(b.deleteOneByID).
		Doc("Find a books by ID").
		Param(ws.PathParameter("id", "identifier of the book").DataType("string")).
		Operation("deleteOneByID").
		Returns(204, "Book deleted", nil).
		Returns(400, "Invalid Request", nil).
		Returns(404, "Book not found", nil))

	container.Add(ws)
}

// POST /api/v1/books
func (b *BookREST) createOne(req *restful.Request, res *restful.Response) {
	// capture the payload of the HTTP Post request
	var newBook NewBook
	err := req.ReadEntity(&newBook)
	if err != nil {
		utils.ReplyError(http.StatusBadRequest, "Invalid payload", err, res)
		return
	}

	// call the api
	book, err := b.api.CreateOne(newBook)
	if err != nil {
		utils.ReplyError(http.StatusInternalServerError, "Oeps...", err, res)
		return
	}

	// handle response
	res.WriteHeaderAndEntity(http.StatusCreated, book)
}

// GET /api/v1/books
func (b *BookREST) findAll(req *restful.Request, res *restful.Response) {
	// call the api
	books, err := b.api.FindAll()
	if err != nil {
		utils.ReplyError(http.StatusInternalServerError, "Oeps...", err, res)
		return
	}
	if books == nil {
		// if no books are found, trick go to send empty array instead of null
		res.WriteHeaderAndEntity(http.StatusOK, []string{})
		return
	}

	// handle response
	res.WriteHeaderAndEntity(http.StatusOK, books)
}

// GET /api/v1/books/{id}
func (b *BookREST) findOneByID(req *restful.Request, res *restful.Response) {
	// capture the id
	id := req.PathParameter("id")
	if !bson.IsObjectIdHex(id) {
		utils.ReplyError(http.StatusBadRequest, "Invalid id", nil, res)
		return
	}
	bsonID := bson.ObjectIdHex(id)

	// call the api
	book, err := b.api.FindOneByID(bsonID)
	if err != nil {
		utils.ReplyError(http.StatusInternalServerError, "Oeps...", err, res)
		return
	}
	if book == nil {
		utils.ReplyError(http.StatusNotFound, "Book not found", nil, res)
		return
	}

	// handle response
	res.WriteHeaderAndEntity(http.StatusOK, book)
}

// DELETE /api/v1/books/{id}
func (b *BookREST) deleteOneByID(req *restful.Request, res *restful.Response) {
	// capture the id
	id := req.PathParameter("id")
	if !bson.IsObjectIdHex(id) {
		utils.ReplyError(http.StatusBadRequest, "Invalid id", nil, res)
		return
	}
	bsonID := bson.ObjectIdHex(id)

	// call the api
	count, err := b.api.DeleteOneByID(bsonID)
	if err != nil {
		utils.ReplyError(http.StatusInternalServerError, "Oeps...", err, res)
		return
	}
	if count == 0 {
		utils.ReplyError(http.StatusNotFound, "Book not found", nil, res)
		return
	}

	// handle response
	res.WriteHeader(http.StatusNoContent)
}
