package books

import (
	"fmt"

	"gopkg.in/mgo.v2/bson"
)

// NewBook object (NO ID)
type NewBook struct {
	Title  string `bson:"title" json:"title"`
	Author string `bson:"author" json:"author"`
}

// Validate validates the NewBook struct
func (b *NewBook) Validate() error {
	if b.Title == "" {
		return fmt.Errorf("Book Title cannot be empty")
	}
	if b.Author == "" {
		return fmt.Errorf("Book Title cannot be empty")
	}

	return nil
}

// Book object, inherits NewBook
type Book struct {
	NewBook
	ID bson.ObjectId `bson:"_id" json:"id"`
}

// Validate validates the Book struct
func (b *Book) Validate() error {
	if !b.ID.Valid() {
		return fmt.Errorf("Book ID is not valid, %v", b.ID)
	}

	return b.NewBook.Validate()
}
