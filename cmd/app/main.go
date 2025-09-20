package main

import (
	"fmt"
	"library-app/library"
	"math/rand"
)

func main() {
	firstGenerator := func(name string) int {
		return rand.Int() % (1e9 + 7)
	}
	secondGenerator := func(name string) int {
		hash := 0
		var mod int = 1e9 + 7
		for _, char := range name {
			hash = (hash*2 + int(char)) % mod
		}
		return hash
	}
	sliceStorage := library.NewSliceStorage()
	bookLibrary := library.NewLibrary(sliceStorage, firstGenerator)
	books := []struct {
		name   string
		author string
	}{
		{"1", "John Doe"},
		{"2", "2"},
		{"3", "John Doe"},
	}
	for _, book := range books {
		err := bookLibrary.AddBook(book.name, book.author)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("Added book %s\n", book.name)
		}
	}
	addedBook, err := bookLibrary.GetBook(books[0].name, books[0].author)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(addedBook)
	}
	addedBook, err = bookLibrary.GetBook(books[1].name, books[1].author)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(addedBook)
	}
	bookLibrary.SetIDGenerator(secondGenerator)
	addedBook, err = bookLibrary.GetBook(books[2].name, books[2].author)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(addedBook)
	}
	mapStorage := library.NewMapStorage()
	bookLibrary.SetStorage(mapStorage)
	for _, book := range books {
		err = bookLibrary.AddBook(book.name, book.author)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("Added book %s\n", book.name)
		}
	}
	addedBook, err = bookLibrary.GetBook(books[2].name, books[2].author)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(addedBook)
	}
	addedBook, err = bookLibrary.GetBook(books[1].name, books[1].author)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(addedBook)
	}
}
