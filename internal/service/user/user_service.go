package user

import (
	"errors"
	"fmt"
	"todo-list-api/internal/logger"
	"todo-list-api/internal/repository"
	"todo-list-api/internal/service/auth"
	"todo-list-api/internal/service/utils"
	"todo-list-api/lib/e"
	"todo-list-api/models"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo  repository.Repository
	jwtSecret string
}

func NewUserService(userRepo *repository.Repository, jwtSecret string) *UserService {
	return &UserService{
		userRepo:  *userRepo,
		jwtSecret: jwtSecret,
	}
}

func (us *UserService) SignUp(reqUser models.User) (string, error) {

	hashedPass, err := utils.HashPassword(reqUser.Password)
	if err != nil {
		return "", err
	}

	user := models.User{
		Name:     reqUser.Name,
		Email:    reqUser.Email,
		Password: hashedPass,
	}

	id, err := us.userRepo.User.CreateUser(&user)
	if err != nil {
		if errors.Is(err, e.UniqueViolationErr) {
			return "", fmt.Errorf("user with this email already exists %v", err)
		}
		return "", err
	}
	logger.Logger.Debug("user created with id", "id", id)

	token, err := auth.CreateToken(us.jwtSecret, int(id))
	if err != nil {
		return "", fmt.Errorf("error creating token: %v", err)
	}

	return token, nil
}

func (us *UserService) LogIn(reqUser models.User) (string, error) {
	dbUser, err := us.userRepo.User.GetUserByEmail(reqUser.Email)
	if err != nil {
		return "", err
	}
	if dbUser == nil {
		return "", e.InvalidCredentialsErr
	}
	err = utils.CompareHashPassword(reqUser.Password, dbUser.Password)
	if err != nil {
		if errors.Is(bcrypt.ErrMismatchedHashAndPassword, err) {
			return "", e.InvalidCredentialsErr
		}
		return "", err
	}

	token, err := auth.CreateToken(us.jwtSecret, dbUser.Id)
	if err != nil {
		return "", fmt.Errorf("error creating token: %v", err)
	}
	return token, nil

}
