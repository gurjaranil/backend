package utill

import (
	"regexp"
	"strconv"
)

func IsContactNumberValid(contactNumber int) bool {
	// Check if contactNumber is empty
	contactNumber2 := strconv.Itoa(contactNumber)
	if len(contactNumber2) < 10 {
		return false
	}

	// Check if contactNumber contains only digits
	match, _ := regexp.MatchString("^[0-9]+$", contactNumber2)
	if !match {
		return false
	}

	// Check if contactNumber length is within a valid range (assuming max length is 15 digits)
	if len(contactNumber2) > 10 {
		return false
	}

	return true
}

func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
func IsNameValid(name string) bool {
	if len(name) == 0 {
		return false
	}
	if len(name) > 50 {
		return false
	}
	return true
}

func IsValidIsbn(isbn int) bool {
	// Check if contactNumber is empty
	contactNumber2 := strconv.Itoa(isbn)
	if len(contactNumber2) < 2 {
		return false
	}

	// Check if contactNumber length is within a valid range (assuming max length is 15 digits)
	if len(contactNumber2) > 10 {
		return false
	}

	return true
}

func IsPasswordValid(password string) bool {
	// Check if password length is between 8 and 15 characters
	if len(password) < 8 || len(password) > 15 {
		return false
	}

	// Check if password contains at least one uppercase letter
	if ok, _ := regexp.MatchString("[A-Z]+", password); !ok {
		return false
	}

	// Check if password contains at least one lowercase letter
	if ok, _ := regexp.MatchString("[a-z]+", password); !ok {
		return false
	}

	// Check if password contains at least one digit
	if ok, _ := regexp.MatchString("[0-9]+", password); !ok {
		return false
	}

	// Check if password contains at least one symbol (non-alphanumeric character)
	if ok, _ := regexp.MatchString(`[!@#$%^&*()-_=+{}\[\];:'",<.>/?\\|]+`, password); !ok {
		return false
	}

	return true // Password meets all criteria
}
