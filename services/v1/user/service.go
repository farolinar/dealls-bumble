package userv1

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/farolinar/dealls-bumble/internal/common/auth"
	"github.com/farolinar/dealls-bumble/internal/common/password"
	"github.com/farolinar/dealls-bumble/internal/common/uid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/zerolog/log"
)

type Service interface {
	Create(ctx context.Context, payload UserCreatePayload) (resp UserAuthentication, err error)
}

type userService struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &userService{repository: repository}
}

func (s *userService) Create(ctx context.Context, payload UserCreatePayload) (resp UserAuthentication, err error) {
	hashedPassword, err := password.Hash(payload.Password)
	if err != nil {
		log.Debug().Msgf("error hashing password: %s", err.Error())
		return
	}

	birthdateTime, err := time.Parse(payload.TimeLayout, payload.Birthdate)
	if err != nil {
		log.Debug().Msgf("error parsing birthdate: %s", err.Error())
		return
	}

	user := &User{
		UID:            uid.GenerateStringID(16),
		Name:           payload.Name,
		Email:          payload.Email,
		Username:       payload.Username,
		HashedPassword: &hashedPassword,
		Sex:            payload.Sex,
		Birthdate:      birthdateTime,
	}
	err = s.repository.Create(ctx, user)
	var pgErr *pgconn.PgError
	if err != nil {
		log.Debug().Msgf("error creating user: %s", err.Error())
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				err = ErrAlreadyExists
			default:
				return
			}
		}
		return
	}

	// create access token with signed jwt
	accessToken, err := auth.CreateAccessToken(fmt.Sprint(user.UID))
	if err != nil {
		log.Debug().Msgf("error creating access token: %s", err.Error())
		return
	}
	resp.Token = accessToken

	// TODO: upload image

	return
}
