package services

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"time"

	"lys-brazil-api/models"
)

type UserService struct {
	logFile string
}

type LogFileEntry struct {
	Date           string `json:"date"`
	StartHour      string `json:"start_hour"`
	HourAndRunning string `json:"hour_and_running"`
	DurationMS     int64  `json:"duration_ms"`
	Status         string `json:"status"`
}

func NewUserService() *UserService {
	return &UserService{logFile: "data/user-service.json"}
}

func (s *UserService) GetUserByCPFAndName(cpf, fullName string) (models.User, error) {

	data, err := ioutil.ReadFile("data/userData.json")
	if err != nil {
		return models.User{}, errors.New("could not read data file")
	}

	var users []models.User
	if err := json.Unmarshal(data, &users); err != nil {
		return models.User{}, errors.New("error deserializing user data")
	}

	for _, user := range users {
		if user.CPF == cpf && user.FullName == fullName {
			return user, nil // Return the found user
		}
	}

	return models.User{}, errors.New("user not found") // User not found
}

func (s *UserService) LogToJSON(entry LogFileEntry) {

	if err := os.MkdirAll("data", os.ModePerm); err != nil {
		log.Fatalf("Could not create data directory: %v", err)
	}

	file, err := os.OpenFile(s.logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Could not open log file: %v", err)
	}
	defer file.Close()

	var logs []LogFileEntry
	if data, err := ioutil.ReadFile(s.logFile); err == nil {
		json.Unmarshal(data, &logs)
	}

	logs = append(logs, entry)
	data, _ := json.MarshalIndent(logs, "", "  ")
	file.Write(data)
}

func (s *UserService) Run(cpf, fullName string) {
	for {
		startTime := time.Now()

		user, err := s.GetUserByCPFAndName(cpf, fullName)
		status := "success"
		if err != nil {
			status = "failure"
			log.Printf("Error retrieving user: %v", err)
		} else {
			log.Printf("Retrieved user: %+v", user)
		}

		durationMS := time.Since(startTime).Milliseconds()

		entry := LogFileEntry{
			Date:           time.Now().Format("2006-01-02"),
			StartHour:      startTime.Format("15:04:05"),
			HourAndRunning: time.Now().Format("15:04:05"),
			DurationMS:     durationMS,
			Status:         status,
		}

		s.LogToJSON(entry)

		if status == "success" {
			break
		}

		time.Sleep(5 * time.Second)
	}
}
