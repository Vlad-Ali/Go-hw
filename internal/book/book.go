package book

type Book struct {
	title  string
	author string
	id     int
}

func NewBook(title string, author string, id int) *Book {
	return &Book{title: title, author: author, id: id}
}

func (book *Book) GetTitle() string {
	return book.title
}

func (book *Book) GetAuthor() string {
	return book.author
}

func (book *Book) GetId() int { return book.id }

func (book *Book) SetId(id int) { book.id = id }
