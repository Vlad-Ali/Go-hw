package slicestorage

import (
	"fmt"
	"library-app/internal/book"
)

type SliceStorage struct {
	books []*book.Book
}

func NewSliceStorage() *SliceStorage {
	return &SliceStorage{books: make([]*book.Book, 0)}
}

func (s *SliceStorage) Add(book *book.Book) error {
	for _, b := range s.books {
		if b.GetId() == book.GetId() {
			return fmt.Errorf("book with id = %d already exists", book.GetId())
		}
	}
	s.books = append(s.books, book)
	return nil
}

func (s *SliceStorage) Remove(id int) error {
	for i, b := range s.books {
		if b.GetId() == id {
			s.books = append(s.books[:i], s.books[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("book with id = %d not found", id)
}

func (s *SliceStorage) GetById(id int) (*book.Book, error) {
	for _, b := range s.books {
		if b.GetId() == id {
			return b, nil
		}
	}
	return nil, fmt.Errorf("book with id = %d not found", id)
}

func (s *SliceStorage) GetAll() []book.Book {
	books := make([]book.Book, len(s.books))
	for i, b := range s.books {
		books[i] = *b
	}
	return books
}
