package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nddat1811/simple_project_golang/dto"
	"github.com/nddat1811/simple_project_golang/entity"
	"github.com/nddat1811/simple_project_golang/helper"
	"github.com/nddat1811/simple_project_golang/service"
)

//AuthController interface is a contract what this controller can do
type AuthController interface {
	Login(ctx *gin.Context)
	Register(ctx *gin.Context)
}

type authController struct {
	//put service
	authService service.AuthService
	jwtService  service.JWTService
}

func NewAuthController(authService service.AuthService, jwtService service.JWTService) AuthController {
	return &authController{
		authService: authService,
		jwtService:  jwtService,
	}
}

// Login godoc
// @Summary Login
// @Schemes
// @Description Login directory
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {string} string "success"
// @Success 400 {string} string "error"
// @Success 404 {string} string "error"
// @Success 500 {string} string "error"
// @Router /auth/login [post]
func (c *authController) Login(ctx *gin.Context) {
	var loginDTO dto.LoginDTO
	errorDTO := ctx.ShouldBind(&loginDTO)
	if errorDTO != nil {
		response := helper.BuildErrorResponse("Failed to process request", errorDTO.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
	}
	authResult := c.authService.VerifyCredential(loginDTO.Email, loginDTO.Password)
	if v, ok := authResult.(entity.User); ok {
		generatedToken := c.jwtService.GenerateToken(strconv.FormatUint(v.ID, 10), v.Name, v.Email)
		v.Token = generatedToken
		response := helper.BuildResponse(true, "OK!", v)
		ctx.JSON(http.StatusOK, response)
		return
	}
	response := helper.BuildErrorResponse("Please check your credential", "Invalid Credential", helper.EmptyObj{})
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
}

// Register godoc
// @Summary Register
// @Schemes
// @Description Register directory
// @Tags Auth
// @Param name path string true "name"
// @Accept json
// @Produce json
// @Success 200 {string} string "success"
// @Success 400 {string} string "error"
// @Success 404 {string} string "error"
// @Success 500 {string} string "error"
// @Router /auth/register[post]
func (c *authController) Register(ctx *gin.Context) {
	var registerDTO dto.RegisterDTO
	errorDTO := ctx.ShouldBind(&registerDTO)
	fmt.Println("hi: ", registerDTO)
	if errorDTO != nil {
		response := helper.BuildErrorResponse("Failed to process request 1", errorDTO.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	if !c.authService.IsDuplicateEmail(registerDTO.Email) {
		response := helper.BuildErrorResponse("Failed to process request 2", "Duplicate email", helper.EmptyObj{})
		ctx.JSON(http.StatusConflict, response)
		return
	} else {
		createdUser := c.authService.CreateUser(registerDTO)
		token := c.jwtService.GenerateToken(strconv.FormatUint(createdUser.ID, 10), createdUser.Name, createdUser.Email)
		createdUser.Token = token
		response := helper.BuildResponse(true, "OK! register sucessful", createdUser)
		ctx.JSON(http.StatusCreated, response)
	}

}
