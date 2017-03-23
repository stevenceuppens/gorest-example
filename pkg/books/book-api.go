package books

import "gopkg.in/mgo.v2/bson"

// BookAPI interface
type BookAPI interface {
	CreateOne(data NewBook) (*Book, error)
	FindAll() ([]*Book, error)
	FindOneByID(id bson.ObjectId) (*Book, error)
	DeleteOneByID(id bson.ObjectId) (int, error)
}
