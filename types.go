package main

import "time"

// User ...
type User struct {
	ExternalID  string    `json:"external_id" gorm:"uniqueIndex"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	DateOfBirth time.Time `json:"date_of_birth"`
}

type apiError struct {
	Error string `json:"error"`
}
