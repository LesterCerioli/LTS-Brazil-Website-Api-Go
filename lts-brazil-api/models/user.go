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
}

func ValidateCPF(cpf string) error {

	re := regexp.MustCompile(`^\d{9}-\d{2}$`)
	if !re.MatchString(cpf) {
		return errors.New("CPF inválido. O formato deve ser xxxxxxxxx-xx")
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
