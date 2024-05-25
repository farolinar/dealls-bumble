package servicebase

import (
	"time"
	"unicode"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var (
	PasswordValidation = `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$`
	MinAge             = 18
)

var PasswordValidationRule = validation.NewStringRule(func(s string) bool {
	if len(s) < 8 {
		return false
	}

	var hasUpper, hasLower, hasSpecial bool
	for _, char := range s {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasSpecial
}, MessagePasswordInvalid)

var MustAbove18Rule = func(birthdate, layout string) bool {
	birthdateTime, err := time.Parse(layout, birthdate)
	if err != nil {
		return false
	}

	// Calculate the age.
	age := CalculateAge(birthdateTime)

	// Check if age is at least 18.
	return age >= MinAge
}

var CalculateAge = func(birthdate time.Time) int {
	now := time.Now()
	age := now.Year() - birthdate.Year()

	// Adjust age if the birthday hasn't occurred yet this year.
	if now.YearDay() < birthdate.YearDay() {
		age--
	}

	return age
}
