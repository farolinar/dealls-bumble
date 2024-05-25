package userv1

import (
	"fmt"

	"github.com/farolinar/dealls-bumble/internal/common/parser"
	servicebase "github.com/farolinar/dealls-bumble/services/base"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type UserCreatePayload struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Sex        Sex    `json:"sex"`
	Birthdate  string `json:"birthdate"`
	TimeLayout string `json:"-"`
}

func (p UserCreatePayload) NewLayoutDateOnly() UserCreatePayload {
	return UserCreatePayload{
		Name:       p.Name,
		Email:      p.Email,
		Username:   p.Username,
		Password:   p.Password,
		Sex:        p.Sex,
		Birthdate:  p.Birthdate,
		TimeLayout: parser.LayoutDateOnly,
	}
}

func (p UserCreatePayload) Validate() error {
	if !servicebase.MustAbove18Rule(p.Birthdate, p.TimeLayout) {
		return fmt.Errorf(fmt.Sprintf("birthdate: %s", MessageMustAbove18))
	}

	return validation.ValidateStruct(&p,
		validation.Field(&p.Name, validation.Required, validation.Length(3, 50)),
		validation.Field(&p.Email, validation.Required, is.Email),
		validation.Field(&p.Username, validation.Required, validation.Length(5, 30)),
		validation.Field(&p.Password, validation.Required, servicebase.PasswordValidationRule),
		validation.Field(&p.Sex, validation.Required, validation.In(SexList...)),
		validation.Field(&p.Birthdate, validation.Required, validation.Date(p.TimeLayout)),
	)
}
