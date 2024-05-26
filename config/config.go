package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"

	"github.com/go-playground/validator/v10"
	"github.com/iamolegga/enviper"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/spf13/viper"
	"gopkg.in/gorp.v2"
)

// Provider the config provider
type Provider interface {
	ConfigFileUsed() string
	Get(key string) interface{}
	GetBool(key string) bool
	GetDuration(key string) time.Duration
	GetFloat64(key string) float64
	GetInt(key string) int
	GetInt64(key string) int64
	GetSizeInBytes(key string) uint
	GetString(key string) string
	GetStringMap(key string) map[string]interface{}
	GetStringMapString(key string) map[string]string
	GetStringMapStringSlice(key string) map[string][]string
	GetStringSlice(key string) []string
	GetTime(key string) time.Time
	InConfig(key string) bool
	IsSet(key string) bool
}

var appConfig AppConfig

func GetConfig() AppConfig {
	return appConfig
}

func LoadEnvConfig() (conf AppConfig, err error) {
	v := viper.New()

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AllowEmptyEnv(true)
	v.AutomaticEnv()

	dirPath, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("error get working dir: %s", err))
	}

	err = godotenv.Load(fmt.Sprintf("%s/.env", dirPath))
	if err != nil {
		return
	}

	e := enviper.New(v)
	err = e.Unmarshal(&conf)
	if err != nil {
		return
	}

	validate := validator.New()
	err = validate.Struct(conf)
	if err != nil {
		return
	}

	appConfig = conf

	return
}

type AppConfig struct {
	App      App      `mapstructure:"app" validate:"required"`
	Postgres Postgres `mapstructure:"postgres" validate:"required"`
}

type App struct {
	Secret          string `mapstructure:"secret" validate:"required"`
	Host            string `mapstructure:"host" validate:"required"`
	Port            int    `mapstructure:"port" validate:"required"`
	Name            string `mapstructure:"name" validate:"required"`
	LogPretty       bool   `mapstructure:"log_pretty" validate:"required"`
	LogLevel        string `mapstructure:"log_level" validate:"required"`
	BCryptSalt      int    `mapstructure:"bcrypt_salt" validate:"required"`
	JWTSecret       string `mapstructure:"jwt_secret" validate:"required"`
	JWTHourDuration int    `mapstructure:"jwt_hour_duration" validate:"required"`
}

type Postgres struct {
	Host            string        `mapstructure:"host" validate:"required"`
	Port            int           `mapstructure:"port" validate:"required"`
	Username        string        `mapstructure:"username" validate:"required"`
	Password        string        `mapstructure:"password" validate:"required"`
	DbName          string        `mapstructure:"name" validate:"required"`
	Params          string        `mapstructure:"params" validate:"required"`
	ConnPoolSize    int           `mapstructure:"conn_pool_size" validate:"required"`
	ConnIdleMax     int           `mapstructure:"conn_idle_max" validate:"required"`
	ConnLifetimeMax time.Duration `mapstructure:"conn_lifetime_max" validate:"required"`
	Timeout         time.Duration `mapstructure:"timeout" validate:"required"`
}

type DbConnection struct{}

var gorpDb *gorp.DbMap

func (d *DbConnection) InitDbConnectionPool(db Postgres) {
	DB, err := sql.Open("pgx",
		fmt.Sprintf("host=%s "+
			"port=%d "+
			"user=%s "+
			"password=%s "+
			"dbname=%s "+
			"%s", db.Host, db.Port, db.Username, db.Password, db.DbName, db.Params))

	if err != nil {
		log.Fatalln(err)
	}

	gorpDb = &gorp.DbMap{Db: DB, Dialect: gorp.PostgresDialect{}}
}

func (d *DbConnection) GetDbConnectionPool() *gorp.DbMap {
	return gorpDb
}
