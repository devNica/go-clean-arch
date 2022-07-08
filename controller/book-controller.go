package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/devNica/go-clean/dto"
	"github.com/devNica/go-clean/entity"
	"github.com/devNica/go-clean/helper"
	"github.com/devNica/go-clean/service"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type BookController interface {
	CreateBook(ctx *gin.Context)
	UpdateBook(ctx *gin.Context)
	RemoveBook(ctx *gin.Context)
	GetAllBook(ctx *gin.Context)
	FindByID(ctx *gin.Context)
	getUserIDByToken(token string) string
}

type bookController struct {
	bookService service.BookService
	jwtService  service.JWTService
}

func NewBookController(bookService service.BookService, jwtService service.JWTService) BookController {
	return &bookController{
		bookService: bookService,
		jwtService:  jwtService,
	}
}

func (c *bookController) getUserIDByToken(token string) string {
	aToken, err := c.jwtService.ValidateToken(token)
	if err != nil {
		panic(err.Error())
	}
	claims := aToken.Claims.(jwt.MapClaims)
	id := fmt.Sprintf("%v", claims["user_id"])
	return id
}

func (c *bookController) CreateBook(ctx *gin.Context) {
	var bookCreateDTO dto.BookCreateDTO
	errDTO := ctx.ShouldBind(&bookCreateDTO)
	if errDTO != nil {
		res := helper.CreateErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	} else {
		authHeader := ctx.GetHeader("Authorization")
		userID := c.getUserIDByToken(authHeader)
		convertedUserID, err := strconv.ParseUint(userID, 10, 64)
		if err == nil {
			bookCreateDTO.UserID = convertedUserID
		}
		result := c.bookService.CreateBook(bookCreateDTO)
		response := helper.CreateResponse(true, "OK!", result)
		ctx.JSON(http.StatusOK, response)
	}

}

func (c *bookController) UpdateBook(ctx *gin.Context) {
	var bookUpdateDTO dto.BookUpdateDTO
	errDTO := ctx.ShouldBind(&bookUpdateDTO)
	if errDTO != nil {
		res := helper.CreateErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	authHeader := ctx.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}
	claims := token.Claims.(jwt.MapClaims)
	userID := fmt.Sprintf("%v", claims["user_id"])
	if c.bookService.IsAllowedToEdit(userID, bookUpdateDTO.ID) {
		id, errID := strconv.ParseUint(userID, 10, 64)
		if errID == nil {
			bookUpdateDTO.UserID = id
		}
		result := c.bookService.UpdateBook(bookUpdateDTO)
		response := helper.CreateResponse(true, "OK!", result)
		ctx.JSON(http.StatusOK, response)
	} else {
		response := helper.CreateErrorResponse("You dont have permissions", "You are not the owner", helper.EmptyObj{})
		ctx.JSON(http.StatusForbidden, response)
	}
}

func (c *bookController) RemoveBook(ctx *gin.Context) {
	var book entity.Book
	id, err := strconv.ParseUint(ctx.Param("id"), 0, 0)
	if err != nil {
		response := helper.CreateErrorResponse("Failed to get id", "No param id was found", helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
	}

	authHeader := ctx.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}
	claims := token.Claims.(jwt.MapClaims)
	userID := fmt.Sprintf("%v", claims["user_id"])

	book.ID = id
	if c.bookService.IsAllowedToEdit(userID, book.ID) {
		c.bookService.RemoveBook(book)
		res := helper.CreateResponse(true, "Deleted", helper.EmptyObj{})
		ctx.JSON(http.StatusOK, res)
	} else {
		response := helper.CreateErrorResponse("You dont have permissions", "You are not the owner", helper.EmptyObj{})
		ctx.JSON(http.StatusForbidden, response)
	}

}

func (c *bookController) GetAllBook(ctx *gin.Context) {
	var books []entity.Book = c.bookService.AllBook()
	res := helper.CreateResponse(true, "OK!", books)
	ctx.JSON(http.StatusOK, res)
}

func (c *bookController) FindByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 0, 0)
	if err != nil {
		res := helper.CreateErrorResponse("No param id was found", err.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var book entity.Book = c.bookService.FindByID(id)
	if (book == entity.Book{}) {
		res := helper.CreateErrorResponse("Data not found", "No data whit given id", helper.EmptyObj{})
		ctx.JSON(http.StatusNotFound, res)
		return
	} else {
		res := helper.CreateResponse(true, "OK!", book)
		ctx.JSON(http.StatusOK, res)
	}

}
