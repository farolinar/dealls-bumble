package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/farolinar/dealls-bumble/config"
	"github.com/farolinar/dealls-bumble/config/postgres"
	"github.com/farolinar/dealls-bumble/internal/common/middleware"
	userv1 "github.com/farolinar/dealls-bumble/services/v1/user"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

func Initialize(cfg config.AppConfig) *mux.Router {

	postgresDB, _ := postgres.NewDBPostgreOptionBuilder().WithHost(cfg.Postgres.Host).
		WithPort(cfg.Postgres.Port).WithUsername(cfg.Postgres.Username).
		WithDBName(cfg.Postgres.DbName).Build()

	db, err := postgresDB.NewPostgreDatabase()

	if err != nil {
		log.Fatal().Msgf("Error connecting to database, will exit | %s", err.Error())
	}

	r := mux.NewRouter()
	r.Use(middleware.Logging)
	r.Use(middleware.PanicRecoverer)
	v1 := r.PathPrefix("/v1").Subrouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "Service ready")
	})

	// healthCheck := r.PathPrefix("/health-check").Subrouter()
	// healthCheck.HandleFunc("/db", readiness.DBReadinessHandler)

	// initialize user domain
	userRepository := userv1.NewRepository(db)
	userService := userv1.NewService(cfg, userRepository)
	userHandler := userv1.NewHandler(cfg, userService)

	ur := v1.PathPrefix("/user").Subrouter()
	ur.HandleFunc("/register", userHandler.CreateUser).Methods(http.MethodPost)
	ur.HandleFunc("/login", userHandler.Login).Methods(http.MethodPost)

	return r
}

func Serve() {
	cfg := config.GetConfig()

	r := Initialize(cfg)

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.App.Port),
		Handler: r,
	}

	go func() {
		log.Info().Msg(fmt.Sprintf("HTTP server listening on %s", httpServer.Addr))
		if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Error().Msg(fmt.Sprintf("HTTP server error: %v", err))
		}
		log.Info().Msg("Stopped serving new connections.")
	}()

	// Listen for the termination signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Block until termination signal received
	<-stop
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	log.Info().Msg(fmt.Sprintf("Shutting down HTTP server listening on %s", httpServer.Addr))
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Error().Msg(fmt.Sprintf("HTTP server shutdown error: %v", err))
	}
	log.Info().Msg("Shutdown complete.")
}
