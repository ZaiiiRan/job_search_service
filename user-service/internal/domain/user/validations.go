package user

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

var (
	phoneRegexp       = regexp.MustCompile(`^(?:\+7|7|8)\d{10}$|^\+375\d{9}$`)
	telegramRegexp    = regexp.MustCompile(`^@[A-Za-z][A-Za-z0-9_]{4,31}$`)
	emailRegexp       = regexp.MustCompile(`^[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,}$`)
	nameRegexp        = regexp.MustCompile(`^[\p{L}][\p{L}\s'\-]{0,63}$`)
	companyNameRegexp = regexp.MustCompile(`^[\p{L}\d\s.,&"'\-]{2,100}$`)
	dateLayout        = "02.01.2006"
)

func ValidatePhoneNumber(phoneNumber string) error {
	phoneNumber = strings.ReplaceAll(phoneNumber, " ", "")
	phoneNumber = strings.ReplaceAll(phoneNumber, "-", "")
	phoneNumber = strings.ReplaceAll(phoneNumber, "(", "")
	phoneNumber = strings.ReplaceAll(phoneNumber, ")", "")
	if !phoneRegexp.MatchString(phoneNumber) {
		return fmt.Errorf("invalid phone (RU/KZ +7|7|8 + 10 digits; BY +375 + 9 digits)")
	}
	return nil
}

func ValidateTelegram(telegram string) error {
	if !telegramRegexp.MatchString(telegram) {
		return fmt.Errorf("invalid username (5–32, starts with letter, letters/digits/_)")
	}
	return nil
}

func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email is required")
	}
	if !emailRegexp.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

func ValidateName(field, value string, required bool) error {
	value = strings.TrimSpace(value)
	if value == "" {
		if required {
			return fmt.Errorf("%s is required", field)
		}
		return nil
	}
	if !nameRegexp.MatchString(value) {
		return fmt.Errorf("%s must contain only letters, spaces, apostrophes or hyphens", field)
	}
	return nil
}

func ValidateBirthDate(birthDate string) error {
	if birthDate == "" {
		return fmt.Errorf("birth date is required")
	}

	dt, err := time.Parse(dateLayout, birthDate)
	if err != nil {
		return fmt.Errorf("birth date must be in format DD.MM.YYYY")
	}

	now := time.Now().UTC()
	if dt.After(now) {
		return fmt.Errorf("birth date cannot be in the future")
	}

	age := now.Year() - dt.Year()
	anniversary := time.Date(now.Year(), dt.Month(), dt.Day(), 0, 0, 0, 0, time.UTC)
	if now.Before(anniversary) {
		age--
	}

	if age < 14 || age > 100 {
		return fmt.Errorf("age must be between 14 and 100")
	}

	return nil
}

func ValidateCompanyName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("company_name is required")
	}

	if !companyNameRegexp.MatchString(name) {
		return fmt.Errorf("company_name contains invalid characters")
	}

	if len([]rune(name)) < 2 {
		return fmt.Errorf("company_name must be at least 2 characters")
	}
	if len([]rune(name)) > 100 {
		return fmt.Errorf("company_name must be no longer than 100 characters")
	}

	return nil
}
