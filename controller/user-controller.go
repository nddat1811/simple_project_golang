package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/nddat1811/simple_project_golang/dto"
	"github.com/nddat1811/simple_project_golang/entity"
	"github.com/nddat1811/simple_project_golang/helper"
	"github.com/nddat1811/simple_project_golang/service"
)

type UserController interface {
	Update(context *gin.Context)
	Profile(context *gin.Context)
	ProfileByName(context *gin.Context)
	GetAllUser(context *gin.Context)
}

type userController struct {
	userService service.UserService
	jwtService  service.JWTService
}

//NewUserController is creating a new instance of UserController
func NewUserController(userService service.UserService, jwtService service.JWTService) UserController {
	return &userController{
		userService: userService,
		jwtService:  jwtService,
	}
}

func (c *userController) Update(context *gin.Context) {
	var userUpdateDTO dto.UserUpdateDTO
	errDTo := context.ShouldBind(&userUpdateDTO)
	if errDTo != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTo.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	authHeader := context.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}
	claims := token.Claims.(jwt.MapClaims)
	id, err := strconv.ParseUint(fmt.Sprintf("%v", claims["user_id"]), 10, 64)
	if err != nil {
		panic(err.Error())
	}
	userUpdateDTO.ID = id
	u := c.userService.Update(userUpdateDTO)
	res := helper.BuildResponse(true, "OK!", u)
	context.JSON(http.StatusOK, res)
}

func (c *userController) Profile(context *gin.Context) {
	authHeader := context.GetHeader("Authorization")
	token, err := c.jwtService.ValidateToken(authHeader)
	if err != nil {
		panic(err.Error())
	}
	claims := token.Claims.(jwt.MapClaims)
	fmt.Println("SOS", claims)
	id := fmt.Sprintf("%v", claims["user+id"])
	fmt.Println("SOS", id)
	user := c.userService.Profile(id)

	res := helper.BuildResponse(true, "OK!", user)
	context.JSON(http.StatusOK, res)
}

func (c *userController) ProfileByName(context *gin.Context) {
	name := context.Param("name")
	if name == "" {
		res := helper.BuildErrorResponse("No param name was found", "Wrong", helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var user entity.User = c.userService.FindProfileByName(name)
	if (user == entity.User{}) {
		res := helper.BuildErrorResponse("Data not found", "No data with given id", helper.EmptyObj{})
		context.JSON(http.StatusNotFound, res)
	} else {
		res := helper.BuildResponse(true, "OK", user)
		context.JSON(http.StatusOK, res)
	}

}

func (c *userController) GetAllUser(context *gin.Context) {
	var users []entity.User = c.userService.GetAll()
	res := helper.BuildResponse(true, "OK", users)
	context.JSON(http.StatusOK, res)
}


// func (c *userController) RegisterUserRoutes(rg *gin.RouterGroup) {
// 	userRoutes := rg.Group("/user", middleware.AuthorizeJWT(jwtService))
// 	userRoutes.GET("/profile", c.Profile)
// 	userRoutes.PUT("/profile", c.Update)
// }
