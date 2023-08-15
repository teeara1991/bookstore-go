package models

import (
	"github.com/jinzhu/gorm"
)

func NewBookModel(db *gorm.DB) *BookModel {
	return &BookModel{db: db}
}

type BookModel struct {
	db *gorm.DB
}

type Book struct {
	gorm.Model
	Name        string `json:"name"`
	Author      string `json:"author"`
	Publication string `json:"publication"`
}

func (bm *BookModel) CreateBook(b *Book) *Book {
	bm.db.NewRecord(b)
	bm.db.Create(&b)
	return b
}

func (bm *BookModel) GetAllBooks() []Book {
	var Books []Book
	bm.db.Find(&Books)
	return Books
}

func (bm *BookModel) GetBookById(Id int64) (*Book, *gorm.DB) {
	var getBook Book
	db := bm.db.Where("ID=?", Id).Find(&getBook)
	return &getBook, db
}

func (bm *BookModel) DeleteBook(ID int64) Book {
	var book Book
	bm.db.Where("ID=?", ID).Delete(book)
	return book
}
