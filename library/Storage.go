package library

import "fmt"

type Storage interface {
	add(book *Book) error
	remove(id int) error
	getById(id int) (*Book, error)
	getAll() []Book
}

type SliceStorage struct {
	books []*Book
}

func NewSliceStorage() *SliceStorage {
	return &SliceStorage{books: make([]*Book, 0)}
}

func (s *SliceStorage) add(book *Book) error {
	for _, b := range s.books {
		if b.id == book.id {
			return fmt.Errorf("book with id = %d already exists", book.id)
		}
	}
	s.books = append(s.books, book)
	return nil
}

func (s *SliceStorage) remove(id int) error {
	for i, b := range s.books {
		if b.id == id {
			s.books = append(s.books[:i], s.books[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("book with id = %d not found", id)
}

func (s *SliceStorage) getById(id int) (*Book, error) {
	for _, b := range s.books {
		if b.id == id {
			return b, nil
		}
	}
	return nil, fmt.Errorf("book with id = %d not found", id)
}

func (s *SliceStorage) getAll() []Book {
	books := make([]Book, len(s.books))
	for i, b := range s.books {
		books[i] = *b
	}
	return books
}

type MapStorage struct {
	books map[int]*Book
}

func NewMapStorage() *MapStorage {
	return &MapStorage{books: make(map[int]*Book)}
}

func (s *MapStorage) add(book *Book) error {

	if _, ok := s.books[book.id]; ok {
		return fmt.Errorf("book with id = %d already exists", book.id)
	}
	s.books[book.id] = book
	return nil
}

func (s *MapStorage) remove(id int) error {
	if _, ok := s.books[id]; !ok {
		return fmt.Errorf("book with id = %d not found", id)
	}
	delete(s.books, id)
	return nil
}

func (s *MapStorage) getById(id int) (*Book, error) {
	book, ok := s.books[id]

	if !ok {
		return nil, fmt.Errorf("book with id = %d not found", id)
	}

	return book, nil
}

func (s *MapStorage) getAll() []Book {
	var books []Book
	for _, book := range s.books {
		books = append(books, *book)
	}
	return books
}
