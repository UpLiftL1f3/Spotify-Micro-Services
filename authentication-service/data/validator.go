package data

import (
	"fmt"
	"regexp"
	"strings"
)

func CreateUserValidator(u *User) error {
	err := FirstNameValidator(strings.TrimSpace(u.FirstName))
	if err != nil {
		return err
	}

	err = EmailValidator(u.Email)
	if err != nil {
		return err
	}

	err = PasswordValidator(strings.TrimSpace(u.Password))
	if err != nil {
		return err
	}

	return nil
}

func FirstNameValidator(firstName string) error {
	if firstName == "" {
		return fmt.Errorf("name is missing")
	}

	if len(firstName) < 3 {
		return fmt.Errorf("name has to be at least 3 characters")
	}

	if len(firstName) > 20 {
		return fmt.Errorf("name is too long! 20 characters is the max")
	}

	return nil
}

func EmailValidator(email string) error {
	if email == "" {
		return fmt.Errorf("email is missing")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	isValid := emailRegex.MatchString(email)
	if !isValid {
		return fmt.Errorf("invalid email")
	}

	return nil
}

func PasswordValidator(password string) error {
	if password == "" {
		return fmt.Errorf("password is missing")
	}

	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}

	// At least one uppercase letter
	if ok, _ := regexp.MatchString(`[A-Z]`, password); !ok {
		return fmt.Errorf("password must have at least one uppercase letter")
	}

	// At least one lowercase letter
	if ok, _ := regexp.MatchString(`[a-z]`, password); !ok {
		return fmt.Errorf("password must have at least one lowercase letter")
	}

	// At least one digit
	if ok, _ := regexp.MatchString(`\d`, password); !ok {
		return fmt.Errorf("password must have at least one digit")
	}

	// At least one special character
	if ok, _ := regexp.MatchString(`[@$!%*?&]`, password); !ok {
		return fmt.Errorf("password must have at least one special character")
	}

	return nil
}
