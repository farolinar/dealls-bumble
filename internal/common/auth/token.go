package auth

import (
	"time"

	"github.com/farolinar/dealls-bumble/config"
	"github.com/farolinar/dealls-bumble/internal/common/jwt"
)

func CreateAccessToken(subject string) (string, error) {
	cfg := config.GetConfig()
	return jwt.Sign(time.Duration(cfg.App.JWTHourDuration)*time.Hour, subject)
}
