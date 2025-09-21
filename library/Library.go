package library

import (
	"errors"
	"fmt"
)

type LibraryOperations interface {
	AddBook(title string, author string) error
	GetBook(title string, author string) (Book, error)
	SetStorage(storage Storage)
	SetIDGenerator(func(string) int)
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
	newBook := NewBook(title, author, bookId)
	err := l.storage.add(newBook)
	if err != nil {
		return fmt.Errorf("error adding book %s", fullName)
	}
	l.titleIds[fullName] = bookId
	return nil
}

func (l *Library) GetBook(title string, author string) (Book, error) {
	fullName := getBookFullName(title, author)
	bookId, ok := l.titleIds[fullName]
	if !ok {
		return Book{}, errors.New("book not found")
	}
	book, err := l.storage.getById(bookId)
	if err != nil {
		return Book{}, err
	}
	return *book, nil
}

func (l *Library) SetStorage(storage Storage) {
	l.storage = storage
}

func (l *Library) SetIDGenerator(generator func(string) int) {
	l.idGenerator = generator
	currentBooks := l.storage.getAll()
	for _, book := range currentBooks {
		err := l.storage.remove(book.id)
		if err != nil {
			fmt.Println(err)
		}
	}
	for _, book := range currentBooks {
		fullName := getBookFullName(book.title, book.author)
		newBookId := l.idGenerator(fullName)
		book.id = newBookId
		err := l.storage.add(&book)
		if err != nil {
			fmt.Println(err)
			continue
		}
		l.titleIds[fullName] = book.id
	}
}
