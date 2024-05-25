package postgres

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"database/sql"

	_ "github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

var dbConn *sql.DB

func GetDBConnection() (db *sql.DB) {
	if db == nil {
		log.Fatal().Msgf("database not initialized")
	}
	return dbConn
}

// DBPostgreOption options for postgre connection
type DBPostgreOption struct {
	Host            string
	Port            int
	Username        string
	Password        string
	DBName          string
	MaxPoolSize     int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// NewPostgreDatabase return gorp dbmap object with postgre options param
func (m DBPostgreOption) NewPostgreDatabaseWithDSN(dsn string) (db *sql.DB, err error) {
	db, err = sql.Open("pgx", dsn)
	if err != nil {
		return
	}

	db.SetConnMaxLifetime(m.ConnMaxLifetime)
	db.SetMaxIdleConns(m.MaxIdleConns)
	db.SetMaxOpenConns(m.MaxPoolSize)

	err = db.Ping()
	if err != nil {
		return
	}

	dbConn = db

	return
}

// NewPostgreDatabase return gorp dbmap object with postgre options param
func (m DBPostgreOption) NewPostgreDatabase() (db *sql.DB, err error) {
	db, err = sql.Open("pgx", fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable", m.Host, m.Port, m.Username, m.DBName, m.Password))
	if err != nil {
		return
	}

	db.SetConnMaxLifetime(m.ConnMaxLifetime)
	db.SetMaxIdleConns(m.MaxIdleConns)
	db.SetMaxOpenConns(m.MaxPoolSize)

	err = db.Ping()
	if err != nil {
		return
	}

	dbConn = db

	return
}

// DBPostgreOption builder pattern code
type DBPostgreOptionBuilder struct {
	dBPostgreOption *DBPostgreOption
}

func NewDBPostgreOptionBuilder() *DBPostgreOptionBuilder {
	postgresqlPort, _ := strconv.Atoi(os.Getenv("POSTGRESQL_PORT"))
	dBPostgreOption := &DBPostgreOption{
		Host:            os.Getenv("POSTGRES_HOST"),
		Port:            postgresqlPort,
		Username:        os.Getenv("POSTGRES_USERNAME"),
		Password:        os.Getenv("POSTGRES_PASSWORD"),
		MaxPoolSize:     10,
		MaxIdleConns:    2,
		ConnMaxLifetime: 30 * time.Second,
	}
	b := &DBPostgreOptionBuilder{dBPostgreOption: dBPostgreOption}
	return b
}

func (b *DBPostgreOptionBuilder) WithHost(host string) *DBPostgreOptionBuilder {
	b.dBPostgreOption.Host = host
	return b
}

func (b *DBPostgreOptionBuilder) WithPort(port int) *DBPostgreOptionBuilder {
	b.dBPostgreOption.Port = port
	return b
}

func (b *DBPostgreOptionBuilder) WithUsername(username string) *DBPostgreOptionBuilder {
	b.dBPostgreOption.Username = username
	return b
}

func (b *DBPostgreOptionBuilder) WithPassword(password string) *DBPostgreOptionBuilder {
	b.dBPostgreOption.Password = password
	return b
}

func (b *DBPostgreOptionBuilder) WithDBName(dBName string) *DBPostgreOptionBuilder {
	b.dBPostgreOption.DBName = dBName
	return b
}

func (b *DBPostgreOptionBuilder) WithMaxPoolSize(maxPoolSize int) *DBPostgreOptionBuilder {
	b.dBPostgreOption.MaxPoolSize = maxPoolSize
	return b
}

func (b *DBPostgreOptionBuilder) WithMaxIdleConns(maxIdleConns int) *DBPostgreOptionBuilder {
	b.dBPostgreOption.MaxIdleConns = maxIdleConns
	return b
}

func (b *DBPostgreOptionBuilder) WithConnMaxLifetime(connMaxLifetime time.Duration) *DBPostgreOptionBuilder {
	b.dBPostgreOption.ConnMaxLifetime = connMaxLifetime
	return b
}

func (b *DBPostgreOptionBuilder) Build() (*DBPostgreOption, error) {
	return b.dBPostgreOption, nil
}
