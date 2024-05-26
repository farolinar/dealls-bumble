package userv1

import (
	"errors"
	"net/http"

	"github.com/farolinar/dealls-bumble/config"
	"github.com/farolinar/dealls-bumble/internal/common/request"
	"github.com/farolinar/dealls-bumble/internal/common/response"
	servicebase "github.com/farolinar/dealls-bumble/services/base"
	"github.com/rs/zerolog/log"
)

var (
	MessageSuccess          = "Success"
	MessageInternalError    = "Internal server error"
	MessageFailedDecodeJSON = "Failed to decode JSON"
)

type Handler struct {
	cfg     config.AppConfig
	service Service
}

func NewHandler(cfg config.AppConfig, service Service) *Handler {
	return &Handler{cfg: cfg, service: service}
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	translateMessage(r)

	var payload UserCreatePayload
	var resp UserAuthenticationResponse

	err := request.DecodeJSON(w, r, &payload)
	if err != nil {
		err = response.JSON(w, http.StatusBadRequest, servicebase.ResponseBody{
			Message: MessageFailedDecodeJSON,
			Code:    servicebase.Code4XX,
		})
		if err != nil {
			log.Error().Msgf("error encoding response body: %v", err)
		}
		return
	}

	payload = payload.NewLayoutDateOnly()
	err = payload.Validate()
	if err != nil {
		err = response.JSON(w, http.StatusBadRequest, servicebase.ResponseBody{
			Message: err.Error(),
			Code:    servicebase.Code4XX,
		})
		if err != nil {
			log.Error().Msgf("error encoding response body: %v", err)
		}
		return
	}

	userResp, err := h.service.Create(r.Context(), payload)
	if errors.Is(err, ErrAlreadyExists) {
		err = response.JSON(w, http.StatusBadRequest, servicebase.ResponseBody{
			Message: err.Error(),
			Code:    servicebase.Code4XX,
		})
		if err != nil {
			log.Error().Msgf("error encoding response body: %v", err)
		}
		return
	}
	if err != nil {
		err = response.JSON(w, http.StatusInternalServerError, servicebase.ResponseBody{
			Message: MessageInternalError,
			Code:    servicebase.Code5XX,
		})
		if err != nil {
			log.Error().Msgf("error encoding response body: %v", err)
		}
		return
	}

	resp.Message = MessageSuccess
	resp.Code = servicebase.CodeSuccess
	resp.Data = &userResp
	err = response.JSON(w, http.StatusCreated, resp)
	if err != nil {
		log.Error().Msgf("error encoding response body: %v", err)
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	translateMessage(r)

	var payload UserLoginPayload
	var resp UserAuthenticationResponse

	err := request.DecodeJSON(w, r, &payload)
	if err != nil {
		err = response.JSON(w, http.StatusBadRequest, servicebase.ResponseBody{
			Message: MessageFailedDecodeJSON,
			Code:    servicebase.Code4XX,
		})
		if err != nil {
			log.Error().Msgf("error encoding response body: %v", err)
		}
		return
	}

	err = payload.Validate()
	if err != nil {
		err = response.JSON(w, http.StatusBadRequest, servicebase.ResponseBody{
			Message: err.Error(),
			Code:    servicebase.Code4XX,
		})
		if err != nil {
			log.Error().Msgf("error encoding response body: %v", err)
		}
		return
	}

	userResp, err := h.service.Login(r.Context(), payload)
	if errors.Is(err, ErrNotFound) {
		err = response.JSON(w, http.StatusBadRequest, servicebase.ResponseBody{
			Message: err.Error(),
			Code:    servicebase.Code4XX,
		})
		if err != nil {
			log.Error().Msgf("error encoding response body: %v", err)
		}
		return
	}
	if errors.Is(err, ErrWrongPassword) {
		err = response.JSON(w, http.StatusBadRequest, servicebase.ResponseBody{
			Message: err.Error(),
			Code:    servicebase.Code4XX,
		})
		if err != nil {
			log.Error().Msgf("error encoding response body: %v", err)
		}
		return
	}
	if err != nil {
		err = response.JSON(w, http.StatusInternalServerError, servicebase.ResponseBody{
			Message: MessageInternalError,
			Code:    servicebase.Code5XX,
		})
		if err != nil {
			log.Error().Msgf("error encoding response body: %v", err)
		}
		return
	}

	resp.Message = MessageSuccess
	resp.Code = servicebase.CodeSuccess
	resp.Data = &userResp
	err = response.JSON(w, http.StatusOK, resp)
	if err != nil {
		log.Error().Msgf("error encoding response body: %v", err)
	}
}

func translateMessage(r *http.Request) {
	lang := r.Header.Get("Accept-Language")
	servicebase.Translate(lang)
	Translate(lang)
}
