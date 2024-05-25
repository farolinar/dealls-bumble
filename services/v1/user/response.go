package userv1

import servicebase "github.com/farolinar/dealls-bumble/services/base"

type UserResponse struct {
	servicebase.ResponseBody
	Data *User `json:"data,omitempty"`
}

type UserAuthenticationResponse struct {
	servicebase.ResponseBody
	Data *UserAuthentication `json:"data,omitempty"`
}

type UserAuthentication struct {
	Token string `json:"token"`
}
