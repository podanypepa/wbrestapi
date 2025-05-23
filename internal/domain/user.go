// Package domain ...
package domain

import "time"

// User struct
type User struct {
	ID          uint      `gorm:"primaryKey"`
	ExternalID  string    `json:"external_id" gorm:"uniqueIndex"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	DateOfBirth time.Time `json:"date_of_birth"`
}
