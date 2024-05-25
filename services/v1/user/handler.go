package userv1

import (
	"errors"
	"net/http"

	"github.com/farolinar/dealls-bumble/internal/common/request"
	"github.com/farolinar/dealls-bumble/internal/common/response"
	servicebase "github.com/farolinar/dealls-bumble/services/base"
)

var (
	MessageSuccess          = "Success"
	MessageInternalError    = "Internal server error"
	MessageFailedDecodeJSON = "Failed to decode JSON"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	translateMessage(r)

	var payload UserCreatePayload
	var resp UserAuthenticationResponse

	err := request.DecodeJSON(w, r, &payload)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, servicebase.ResponseBody{
			Message: MessageFailedDecodeJSON,
			Code:    servicebase.Code4XX,
		})
		return
	}

	payload = payload.NewLayoutDateOnly()
	err = payload.Validate()
	if err != nil {
		response.JSON(w, http.StatusBadRequest, servicebase.ResponseBody{
			Message: err.Error(),
			Code:    servicebase.Code4XX,
		})
		return
	}

	userResp, err := h.service.Create(r.Context(), payload)
	if errors.Is(err, ErrAlreadyExists) {
		response.JSON(w, http.StatusBadRequest, servicebase.ResponseBody{
			Message: err.Error(),
			Code:    servicebase.Code4XX,
		})
		return
	}
	if err != nil {
		response.JSON(w, http.StatusInternalServerError, servicebase.ResponseBody{
			Message: MessageInternalError,
			Code:    servicebase.Code5XX,
		})
		return
	}

	resp.Message = MessageSuccess
	resp.Code = servicebase.CodeSuccess
	resp.Data = &userResp
	response.JSON(w, http.StatusCreated, resp)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	translateMessage(r)

	var payload UserLoginPayload
	var resp UserAuthenticationResponse

	err := request.DecodeJSON(w, r, &payload)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, servicebase.ResponseBody{
			Message: MessageFailedDecodeJSON,
			Code:    servicebase.Code4XX,
		})
		return
	}

	err = payload.Validate()
	if err != nil {
		response.JSON(w, http.StatusBadRequest, servicebase.ResponseBody{
			Message: err.Error(),
			Code:    servicebase.Code4XX,
		})
		return
	}

	userResp, err := h.service.Login(r.Context(), payload)
	if errors.Is(err, ErrNotFound) {
		response.JSON(w, http.StatusBadRequest, servicebase.ResponseBody{
			Message: err.Error(),
			Code:    servicebase.Code4XX,
		})
		return
	}
	if errors.Is(err, ErrWrongPassword) {
		response.JSON(w, http.StatusBadRequest, servicebase.ResponseBody{
			Message: err.Error(),
			Code:    servicebase.Code4XX,
		})
		return
	}
	if err != nil {
		response.JSON(w, http.StatusInternalServerError, servicebase.ResponseBody{
			Message: MessageInternalError,
			Code:    servicebase.Code5XX,
		})
		return
	}

	resp.Message = MessageSuccess
	resp.Code = servicebase.CodeSuccess
	resp.Data = &userResp
	response.JSON(w, http.StatusOK, resp)
}

func translateMessage(r *http.Request) {
	lang := r.Header.Get("Accept-Language")
	servicebase.Translate(lang)
	Translate(lang)
}
