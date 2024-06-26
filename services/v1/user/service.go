package userv1

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/farolinar/dealls-bumble/config"
	"github.com/farolinar/dealls-bumble/internal/common/auth"
	"github.com/farolinar/dealls-bumble/internal/common/password"
	"github.com/farolinar/dealls-bumble/internal/common/uid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/zerolog/log"
)

type Service interface {
	Create(ctx context.Context, payload UserCreatePayload) (resp UserAuthentication, err error)
	Login(ctx context.Context, payload UserLoginPayload) (resp UserAuthentication, err error)
}

type userService struct {
	cfg        config.AppConfig
	repository Repository
}

func NewService(cfg config.AppConfig, repository Repository) Service {
	return &userService{cfg: cfg, repository: repository}
}

func (s *userService) Create(ctx context.Context, payload UserCreatePayload) (resp UserAuthentication, err error) {
	hashedPassword, err := password.Hash(s.cfg.App.BCryptSalt, payload.Password)
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
	accessToken, err := auth.CreateAccessToken(s.cfg, fmt.Sprint(user.UID))
	if err != nil {
		log.Debug().Msgf("error creating access token: %s", err.Error())
		return
	}
	resp.Token = accessToken

	// TODO: upload image

	return
}

func (s *userService) Login(ctx context.Context, payload UserLoginPayload) (resp UserAuthentication, err error) {
	user, err := s.repository.GetByUsername(ctx, payload.Username)
	if err != nil {
		log.Debug().Msgf("error getting user: %v", err)
		if err == sql.ErrNoRows {
			err = ErrNotFound
		}
		return
	}
	match, err := password.Matches(payload.Password, *user.HashedPassword)
	if err != nil {
		log.Debug().Msgf("error matching password: %v", err)
		return
	}
	if !match {
		err = ErrWrongPassword
		return
	}

	// create access token with signed jwt
	accessToken, err := auth.CreateAccessToken(s.cfg, fmt.Sprint(user.UID))
	if err != nil {
		log.Debug().Msgf("error creating access token: %v", err)
		return
	}

	resp.Token = accessToken
	return
}
