package books

import (
	"sync"

	"gopkg.in/mgo.v2/bson"
)

/**
* BookAPIMemory uses Mutex
* http://www.alexedwards.net/blog/understanding-mutexes
*
* Every incomming call is handeled in a new goroutine,
* so sharing resources between endpoints can lead to race conditions
 */

// BookAPIMemory ...
type BookAPIMemory struct {
	sync.RWMutex
	books []*Book
}

// NewBookAPIMemory ...
func NewBookAPIMemory() *BookAPIMemory {
	return &BookAPIMemory{}
}

// CreateOne ...
func (b *BookAPIMemory) CreateOne(data NewBook) (*Book, error) {
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

	b.Lock()
	defer b.Unlock()

	b.books = append(b.books, book)

	return book, nil
}

// FindAll ...
func (b *BookAPIMemory) FindAll() ([]*Book, error) {
	b.RLock()
	defer b.RUnlock()

	return b.books, nil
}

// FindOneByID ...
func (b *BookAPIMemory) FindOneByID(id bson.ObjectId) (*Book, error) {
	b.RLock()
	defer b.RUnlock()

	for _, book := range b.books {
		if book.ID == id {
			return book, nil
		}
	}

	return nil, nil
}

// DeleteOneByID ...
func (b *BookAPIMemory) DeleteOneByID(id bson.ObjectId) (int, error) {
	b.Lock()
	defer b.Unlock()

	for index, book := range b.books {
		if book.ID == id {
			s := b.books
			s = append(s[:index], s[index+1:]...)
			b.books = s
			return 1, nil
		}
	}
	return 0, nil
}
