package library

import (
	"errors"
	"fmt"
	"library-app/internal/book"
)

type Storage interface {
	Add(book *book.Book) error
	Remove(id int) error
	GetById(id int) (*book.Book, error)
	GetAll() []book.Book
}

type Library struct {
	storage     Storage
	titleIds    map[string]int
	idGenerator func(string) int
}

func NewLibrary(storage Storage, generator func(string) int) *Library {
	return &Library{storage: storage, titleIds: make(map[string]int), idGenerator: generator}
}

func getBookFullName(title string, author string) string {
	return fmt.Sprintf("%s %s", title, author)
}

func (l *Library) AddBook(title string, author string) error {
	fullName := getBookFullName(title, author)
	bookId := l.idGenerator(fullName)
	newBook := book.NewBook(title, author, bookId)
	err := l.storage.Add(newBook)
	if err != nil {
		return fmt.Errorf("error adding book %s", fullName)
	}
	l.titleIds[fullName] = bookId
	return nil
}

func (l *Library) GetBook(title string, author string) (book.Book, error) {
	fullName := getBookFullName(title, author)
	bookId, ok := l.titleIds[fullName]
	if !ok {
		return book.Book{}, errors.New("foundBook not found")
	}
	foundBook, err := l.storage.GetById(bookId)
	if err != nil {
		return book.Book{}, err
	}
	return *foundBook, nil
}

func (l *Library) SetStorage(storage Storage) {
	l.storage = storage
}

func (l *Library) SetIDGenerator(generator func(string) int) {
	l.idGenerator = generator
	currentBooks := l.storage.GetAll()
	for _, foundBook := range currentBooks {
		err := l.storage.Remove(foundBook.GetId())
		if err != nil {
			fmt.Println(err)
		}
	}
	for _, foundBook := range currentBooks {
		fullName := getBookFullName(foundBook.GetTitle(), foundBook.GetAuthor())
		newBookId := l.idGenerator(fullName)
		newBook := book.NewBook(foundBook.GetTitle(), foundBook.GetAuthor(), newBookId)
		err := l.storage.Add(newBook)
		if err != nil {
			fmt.Println(err)
			continue
		}
		l.titleIds[fullName] = newBookId
	}
}
