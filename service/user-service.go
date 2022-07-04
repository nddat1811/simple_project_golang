package service

import (
	"log"

	"github.com/mashingan/smapping"
	"github.com/nddat1811/simple_project_golang/dto"
	"github.com/nddat1811/simple_project_golang/entity"
	"github.com/nddat1811/simple_project_golang/repository"
)

type UserService interface {
	Update(user dto.UserUpdateDTO) entity.User
	Profile(userID string) entity.User
	FindProfileByName(userName string) entity.User
	GetAll() []entity.User
}

type userService struct {
	userRepository repository.UserRepository
}

//NewUserService creates a new instance of UserService
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepo,
	}
}

func (service *userService) Update(user dto.UserUpdateDTO) entity.User {
	userToUpdate := entity.User{}
	err := smapping.FillStruct(&userToUpdate, smapping.MapFields(&user))
	if err != nil {
		log.Fatalf("Failed map %v: ", err)
	}

	updatedUser := service.userRepository.UpdateUser(userToUpdate)
	return updatedUser
}

func (service *userService) Profile(userID string) entity.User {
	return service.userRepository.ProfileUSer(userID)
}


func (service *userService) FindProfileByName(userName string) entity.User {
	return service.userRepository.FindProfileByName(userName)
}

func (service *userService) GetAll() []entity.User {
	return service.userRepository.AllUser()
}