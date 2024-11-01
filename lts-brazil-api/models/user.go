package models

import (
	"errors"
	"regexp"
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	RoleName   string    `json:"role_name"`
	Permission string    `json:"permission"` // "read" or "write"
}

type User struct {
	ID             uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	FullName       string    `json:"full_name" gorm:"not null"`
	CPF            string    `json:"cpf" gorm:"unique;not null"`
	Birthdate      time.Time `json:"birthdate"`
	PhoneNumber    string    `json:"phone_number"`
	RoleName       string    `json:"role"`
	PermissionType string    `json:"permission"`
	Email          string    `json:"email" gorm:"not null"`
	Login          string    `json:"login" gorm:"not null"`
	Password       string    `json:"password" gorm:"not null"`
}

func ValidateCPF(cpf string) error {

	re := regexp.MustCompile(`^\d{3}\.\d{3}\.\d{3}-\d{2}$`)
	if !re.MatchString(cpf) {
		return errors.New("CPF inválido. O formato deve ser xxx.xxx.xxx-xx")
	}

	// Remove non-numeric characters
	cpf = regexp.MustCompile(`[^0-9]`).ReplaceAllString(cpf, "")

	// Check for known invalid CPFs
	invalids := []string{
		"00000000000", "11111111111", "22222222222", "33333333333",
		"44444444444", "55555555555", "66666666666", "77777777777",
		"88888888888", "99999999999",
	}

	for _, invalid := range invalids {
		if cpf == invalid {
			return errors.New("CPF inválido")
		}
	}

	// Validate CPF algorithm
	sum := 0
	for i := 0; i < 9; i++ {
		sum += int(cpf[i]-'0') * (10 - i)
	}
	digit1 := 11 - (sum % 11)
	if digit1 >= 10 {
		digit1 = 0
	}

	sum = 0
	for i := 0; i < 10; i++ {
		sum += int(cpf[i]-'0') * (11 - i)
	}
	digit2 := 11 - (sum % 11)
	if digit2 >= 10 {
		digit2 = 0
	}

	if digit1 != int(cpf[9]-'0') || digit2 != int(cpf[10]-'0') {
		return errors.New("CPF inválido")
	}

	return nil
}

func ValidatePhoneNumber(phone string) error {
	re := regexp.MustCompile(`^\d{2}-\d{9}$`)
	if !re.MatchString(phone) {
		return errors.New("Número de telefone inválido. O formato deve ser dd-xxxxxxxxx")
	}
	return nil
}

func ValidateBirthdate(birthdate string) error {
	re := regexp.MustCompile(`^\d{2}/\d{2}/\d{4}$`)
	if !re.MatchString(birthdate) {
		return errors.New("Data de nascimento inválida. O formato deve ser DD/MM/AAAA")
	}

	_, err := time.Parse("02/01/2006", birthdate)
	if err != nil {
		return errors.New("Data de nascimento inválida")
	}
	return nil
}
