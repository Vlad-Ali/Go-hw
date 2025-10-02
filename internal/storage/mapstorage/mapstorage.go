package mapstorage

import (
	"fmt"
	"library-app/internal/book"
)

type MapStorage struct {
	books map[int]*book.Book
}

func NewMapStorage() *MapStorage {
	return &MapStorage{books: make(map[int]*book.Book)}
}

func (s *MapStorage) Add(book *book.Book) error {

	if _, ok := s.books[book.GetId()]; ok {
		return fmt.Errorf("book with id = %d already exists", book.GetId())
	}
	s.books[book.GetId()] = book
	return nil
}

func (s *MapStorage) Remove(id int) error {
	if _, ok := s.books[id]; !ok {
		return fmt.Errorf("book with id = %d not found", id)
	}
	delete(s.books, id)
	return nil
}

func (s *MapStorage) GetById(id int) (*book.Book, error) {
	foundBook, ok := s.books[id]

	if !ok {
		return nil, fmt.Errorf("foundBook with id = %d not found", id)
	}

	return foundBook, nil
}

func (s *MapStorage) GetAll() []book.Book {
	var books []book.Book
	for _, foundBook := range s.books {
		books = append(books, *foundBook)
	}
	return books
}
