package userv1

import "errors"

var (
	ErrNotFound         = errors.New(MessageNotFound)
	ErrWrongPassword    = errors.New(MessageWrongPassword)
	ErrAlreadyExists    = errors.New(MessageAlreadyExists)
	ErrValidationFailed = errors.New(MessageValidationFailed)
)
