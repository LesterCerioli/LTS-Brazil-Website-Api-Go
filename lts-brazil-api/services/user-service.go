package services

import (
	"encoding/json"
	"errors"
	"lts-brazil-api/models"
	"os"
	"sync"
	"time"
)

type UserService struct {
	users    []models.User
	mu       sync.Mutex
	filePath string
}

func NewUserService() *UserService {
	us := &UserService{
		filePath: "data/user-data.json",
	}
	us.LoadData()
	return us
}

func (us *UserService) CreateUser(fullName, cpf, birthDate, phoneNumber, roleName, permissionType, email, login, password string) error {
	us.mu.Lock()
	defer us.mu.Unlock()

	// Check for duplicate CPF
	for _, user := range us.users {
		if user.CPF == cpf {
			return errors.New("User with this CPF already exists")
		}
	}

	// Parse the birthDate string to time.Time in dd/mm/yyyy format
	birthDateParsed, err := time.Parse("02/01/2006", birthDate)
	if err != nil {
		return errors.New("Invalid birthDate format. Use DD/MM/YYYY")
	}

	newUser := models.User{
		FullName:       fullName,
		CPF:            cpf,
		Birthdate:      birthDateParsed,
		PhoneNumber:    phoneNumber,
		RoleName:       roleName,
		PermissionType: permissionType,
		Email:          email,
		Login:          login,
		Password:       password,
	}

	us.users = append(us.users, newUser)
	us.saveData()
	return nil
}

func (us *UserService) saveData() {
	if err := os.MkdirAll("data", os.ModePerm); err != nil {
		panic(err)
	}
	file, err := json.MarshalIndent(us.users, "", "  ")
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile(us.filePath, file, 0644); err != nil {
		panic(err)
	}
}

func (us *UserService) LoadData() {
	us.mu.Lock()
	defer us.mu.Unlock()
	if _, err := os.Stat(us.filePath); os.IsNotExist(err) {
		return
	}
	file, err := os.ReadFile(us.filePath)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(file, &us.users); err != nil {
		panic(err)
	}
}
