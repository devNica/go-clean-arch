package main

import (
	"github.com/devNica/go-clean-arch/config"
	"github.com/devNica/go-clean-arch/controller"
	"github.com/devNica/go-clean-arch/middleware"
	"github.com/devNica/go-clean-arch/repository"
	"github.com/devNica/go-clean-arch/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	db             *gorm.DB                  = config.SetupDatabaseConnection()
	userRepository repository.UserRepository = repository.NewUserRepository(db)
	bookRepository repository.BookRepository = repository.NewBookRepository(db)
	jwtService     service.JWTService        = service.NewJWTService()
	userService    service.UserService       = service.NewUserService(userRepository)
	authService    service.AuthService       = service.NewAuthService(userRepository)
	bookService    service.BookService       = service.NewBookService(bookRepository)
	authController controller.AuthController = controller.NewAuthController(authService, jwtService)
	userController controller.UserController = controller.NewUserController(userService, jwtService)
	bookController controller.BookController = controller.NewBookController(bookService, jwtService)
)

func main() {
	defer config.CloseDatabaseConnection(db)
	r := gin.Default()

	// middleware.AuthorizeJWT(jwtService)
	authRoutes := r.Group("api/auth")
	{
		authRoutes.POST("/signin", authController.Signin)
		authRoutes.POST("/signup", authController.Signup)
	}

	userRoutes := r.Group("api/users", middleware.AuthorizeJWT(jwtService))
	{
		userRoutes.GET("/profile", userController.Profile)
		userRoutes.PUT("/profile", userController.Update)

	}

	bookRoutes := r.Group("api/books", middleware.AuthorizeJWT(jwtService))
	{
		bookRoutes.GET("/", bookController.GetAllBook)
		bookRoutes.GET("/:id", bookController.FindByID)
		bookRoutes.POST("/", bookController.CreateBook)
		bookRoutes.PUT("/:id", bookController.UpdateBook)
		bookRoutes.DELETE("/:id", bookController.RemoveBook)

	}

	r.Run(":6700")
}
