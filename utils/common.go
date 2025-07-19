package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"

	"lambda-server/models"

)

func GenerateUserID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return fmt.Sprintf("user_%s", base64.URLEncoding.EncodeToString(bytes)[:22])
}

func GenerateJournalID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return fmt.Sprintf("journal_%s", base64.URLEncoding.EncodeToString(bytes)[:22])
}

// IsRunningLocally checks if the application is running locally
func IsRunningLocally() bool {
	return os.Getenv("AWS_LAMBDA_FUNCTION_NAME") == ""
}

// Helper functions for emergency contact validation
func isContactEmpty(contact models.Emergency) bool {
	return contact.Name == "" && contact.Email == "" && contact.Phone == ""
}

func isContactComplete(contact models.Emergency) bool {
	return contact.Name != "" && contact.Email != "" && contact.Phone != ""
}

// ValidateEmergencyContacts checks that the first contact is mandatory (all fields required),
// while the rest are optional (can be empty or complete, but not partial).
// Accepts a slice of contacts for flexibility.
func ValidateEmergencyContacts(contacts []models.Emergency) error {
	if len(contacts) == 0 {
		return fmt.Errorf("at least one emergency contact is required")
	}
	for i, contact := range contacts {
		if i == 0 {
			if !isContactComplete(contact) {
				return fmt.Errorf("first emergency contact is mandatory - all fields (name, email, phone) are required")
			}
		} else {
			if !isContactEmpty(contact) && !isContactComplete(contact) {
				return fmt.Errorf("all fields (name, email, phone) are required for contact %d when any field is provided", i+1)
			}
		}
	}
	return nil
}

// ValidateDOBFormat checks if dob is in dd/mm/yyyy format and is a valid date
func ValidateDOBFormat(dob string) error {
	if len(dob) != 10 {
		return fmt.Errorf("dob must be in dd/mm/yyyy format")
	}
	var day, month, year int
	// Only accept dd/mm/yyyy
	n, err := fmt.Sscanf(dob, "%02d/%02d/%04d", &day, &month, &year)
	if err != nil || n != 3 {
		return fmt.Errorf("dob must be in dd/mm/yyyy format")
	}
	if day < 1 || day > 31 || month < 1 || month > 12 || year < 1900 {
		return fmt.Errorf("dob contains invalid date values")
	}
	// Check for valid day in month
	if month == 2 {
		leap := (year%4 == 0 && year%100 != 0) || (year%400 == 0)
		maxDay := 28
		if leap {
			maxDay = 29
		}
		if day > maxDay {
			return fmt.Errorf("dob contains invalid day for February")
		}
	} else if month == 4 || month == 6 || month == 9 || month == 11 {
		if day > 30 {
			return fmt.Errorf("dob contains invalid day for the month")
		}
	}
	return nil
}

// GeneratePasswordResetToken generates a secure random token for password reset
func GeneratePasswordResetToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return base64.URLEncoding.EncodeToString(bytes)
}
