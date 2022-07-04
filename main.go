package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/nddat1811/simple_project_golang/config"
	"github.com/nddat1811/simple_project_golang/controller"
	"github.com/nddat1811/simple_project_golang/middleware"
	"github.com/nddat1811/simple_project_golang/repository"
	"github.com/nddat1811/simple_project_golang/service"
	"gorm.io/gorm"
)

var (
	r              *gin.Engine
	db             *gorm.DB                  = config.SetupDatabaseConnection()
	userRepository repository.UserRepository = repository.NewUserRepository(db)
	bookRepository repository.BookRepository = repository.NewBookRepository(db)
	jwtService     service.JWTService        = service.NewJWTService()
	userService    service.UserService       = service.NewUserService(userRepository)
	bookService    service.BookService       = service.NewBookService(bookRepository)
	authService    service.AuthService       = service.NewAuthService(userRepository)
	authController controller.AuthController = controller.NewAuthController(authService, jwtService)
	userController controller.UserController = controller.NewUserController(userService, jwtService)
	bookController controller.BookController = controller.NewBookController(bookService, jwtService)
)

// @title Swagger demo service API
// @version 1.0
// @description This is demo server
// @host localhost:9090
// @BasePath /api
// @securityDeifinitions.basic BasicAuth
// @securityDefinitions.apiKey ApiKeyAuth
// @in hearder
// @name Authorization
//docs "github.com/nddat1811/simple_project_golang/docs"
//swaggerFiles "github.com/swaggo/files"
//ginSwagger "github.com/swaggo/gin-swagger"

func main() {
	defer config.CloseDatabaseConnection(db)
	r := gin.Default()

	//docs.SwaggerInfo.BasePath = "/api"
	authRoutes := r.Group("api/auth")
	{
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/register", authController.Register)
	}

	userRoutes := r.Group("api/user", middleware.AuthorizeJWT(jwtService))
	{
		userRoutes.GET("/profile", userController.Profile)
		userRoutes.PUT("/profile", userController.Update)
		//userRoutes.GET("profile/:name", userController.ProfileByName)
	}

	//no authorize
	userRoute2s := r.Group("api/user")
	{

		userRoute2s.GET("profile/:name", userController.ProfileByName)
		userRoute2s.GET("/getall", userController.GetAllUser)
	}

	bookRoutes := r.Group("api/books", middleware.AuthorizeJWT(jwtService))
	{
		bookRoutes.GET("/", bookController.All)
		bookRoutes.POST("/new", bookController.Insert)
		bookRoutes.GET("/:id", bookController.FindByID)
		bookRoutes.PUT("/:id", bookController.Update)
		bookRoutes.DELETE("/:id", bookController.Delete)
	}

	//r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	log.Fatal(r.Run(":9090"))
}
