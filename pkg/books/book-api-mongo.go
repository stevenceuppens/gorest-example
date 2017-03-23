package books

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

/**
* BookAPIMongo doen not use Mutex
*
* Instead it uses session.Copy() from mgo library
* this creates (safely) a new session for each goroutine
 */

// BookAPIMongo ...
type BookAPIMongo struct {
	session *mgo.Session
}

// NewBookAPIMongo ...
func NewBookAPIMongo(session *mgo.Session) *BookAPIMongo {
	return &BookAPIMongo{
		session: session,
	}
}

// CreateOne ...
func (b *BookAPIMongo) CreateOne(data NewBook) (*Book, error) {
	// make sure new data is valid
	err := data.Validate()
	if err != nil {
		return nil, err
	}

	// copy NewBook into Book object and create id
	book := &Book{
		ID:      bson.NewObjectId(),
		NewBook: data,
	}

	// create a mgo session
	session := b.session.Copy()
	defer session.Close()

	// get collection and insert data
	collection := session.DB("").C("books")
	err = collection.Insert(book)
	if err != nil {
		return nil, err
	}

	return book, nil
}

// FindAll ...
func (b *BookAPIMongo) FindAll() ([]*Book, error) {
	session := b.session.Copy()
	defer session.Close()

	var books []*Book

	collection := session.DB("").C("books")
	err := collection.Find(bson.M{}).All(&books)
	if err != nil {
		return nil, err
	}

	return books, nil
}

// FindOneByID ...
func (b *BookAPIMongo) FindOneByID(id bson.ObjectId) (*Book, error) {
	session := b.session.Copy()
	defer session.Close()

	var book Book

	collection := session.DB("").C("books")
	err := collection.FindId(id).One(&book)
	if err != nil {
		return nil, err
	}

	return &book, nil
}

// DeleteOneByID ...
func (b *BookAPIMongo) DeleteOneByID(id bson.ObjectId) (int, error) {
	session := b.session.Copy()
	defer session.Close()

	collection := session.DB("").C("books")
	info, err := collection.RemoveAll(bson.M{"_id": id})
	if err != nil {
		return 0, err
	}

	return info.Matched, nil
}
