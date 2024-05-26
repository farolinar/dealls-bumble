package auth

import (
	"time"

	"github.com/farolinar/dealls-bumble/config"
	"github.com/farolinar/dealls-bumble/internal/common/jwt"
)

func CreateAccessToken(cfg config.AppConfig, subject string) (string, error) {
	return jwt.Sign(time.Duration(cfg.App.JWTHourDuration)*time.Hour, cfg.App.Secret, subject)
}
