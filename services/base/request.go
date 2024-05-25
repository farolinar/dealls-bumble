package servicebase

import (
	"regexp"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var (
	PasswordValidation = `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$`
	MinAge             = 18
)

var PasswordValidationRule = validation.NewStringRule(func(s string) bool {
	// Compile the regex
	r := regexp.MustCompile(PasswordValidation)

	return r.MatchString(s)
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
