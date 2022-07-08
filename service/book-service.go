package service

import (
	"fmt"
	"log"

	"github.com/devNica/go-clean-arch/dto"
	"github.com/devNica/go-clean-arch/entity"
	"github.com/devNica/go-clean-arch/repository"
	"github.com/mashingan/smapping"
)

type BookService interface {
	CreateBook(book dto.BookCreateDTO) entity.Book
	UpdateBook(book dto.BookUpdateDTO) entity.Book
	RemoveBook(book entity.Book)
	AllBook() []entity.Book
	FindByID(bookID uint64) entity.Book
	IsAllowedToEdit(userID string, bookID uint64) bool
}

type bookService struct {
	bookRepository repository.BookRepository
}

func NewBookService(bookRep repository.BookRepository) BookService {
	return &bookService{
		bookRepository: bookRep,
	}
}

func (service *bookService) CreateBook(b dto.BookCreateDTO) entity.Book {
	book := entity.Book{}
	err := smapping.FillStruct(&book, smapping.MapFields(&b))
	if err != nil {
		log.Fatalf("Failed map %v: ", err)
	}
	res := service.bookRepository.InsertBook(book)
	return res
}

func (service *bookService) UpdateBook(b dto.BookUpdateDTO) entity.Book {
	book := entity.Book{}
	err := smapping.FillStruct(&book, smapping.MapFields(&b))
	if err != nil {
		log.Fatalf("Failed map %v: ", err)
	}
	res := service.bookRepository.UpdateBook(book)
	return res
}

func (service *bookService) RemoveBook(b entity.Book) {
	service.bookRepository.DeleteBook(b)
}

func (service *bookService) AllBook() []entity.Book {
	return service.bookRepository.AllBook()
}

func (service *bookService) FindByID(bookID uint64) entity.Book {
	return service.bookRepository.FindBookByID(bookID)
}

func (service *bookService) IsAllowedToEdit(userID string, bookID uint64) bool {
	b := service.bookRepository.FindBookByID(bookID)
	id := fmt.Sprintf("%v", b.UserID)
	return userID == id
}
