package service

import (
	"log"

	"github.com/mashingan/smapping"
	"github.com/nddat1811/simple_project_golang/dto"
	"github.com/nddat1811/simple_project_golang/entity"
	"github.com/nddat1811/simple_project_golang/repository"
	"golang.org/x/crypto/bcrypt"
)

//AuthService is a contract about sth that this service can do
type AuthService interface {
	VerifyCredential(email string, password string) interface{}
	CreateUser(user dto.RegisterDTO) entity.User
	FindByEmail(email string) entity.User
	IsDuplicateEmail(email string) bool
}

type authService struct {
	userReposity repository.UserRepository
}

//NewAuthService creates a new instance of AuthService
func NewAuthService(userRep repository.UserRepository) AuthService {
	return &authService{
		userReposity: userRep,
	}
}

func (service *authService) VerifyCredential(email string, password string) interface{} {
	res := service.userReposity.VerifyCredential(email, password)
	if v, ok := res.(entity.User); ok {
		comparedPassword := comparedPassword(v.Password, []byte(password))
		if v.Email == email && comparedPassword {
			return res
		}
		return false
	}
	return false
}

func (service *authService) CreateUser(user dto.RegisterDTO) entity.User {
	userToCreate := entity.User{}

	err := smapping.FillStruct(&userToCreate, smapping.MapFields(&user))
	if err != nil {
		log.Fatalf("Failed map %v", err)
	}
	res := service.userReposity.InsertUser(userToCreate)
	return res
}

func (service *authService) FindByEmail(email string) entity.User {
	return service.userReposity.FindByEmail(email)
}

func (service *authService) IsDuplicateEmail(email string) bool {
	res := service.userReposity.IsDuplicateEmail(email)
	return !(res.Error == nil)
}

func comparedPassword(hashedPwd string, plainPassword []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPassword)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
