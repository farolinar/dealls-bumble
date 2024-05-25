package userv1

import servicebase "github.com/farolinar/dealls-bumble/services/base"

var (
	MessageNotFound         = "User not found"
	MessageWrongPassword    = "Wrong password"
	MessageAlreadyExists    = "User already exists"
	MessageValidationFailed = "Validation failed"
	MessageMustAbove18      = "Age must above 18"
)

func Translate(lang string) {
	switch lang {
	case servicebase.ID_LANG:
		MessageNotFound = "User tidak ditemukan"
		MessageWrongPassword = "Password salah"
		MessageAlreadyExists = "User sudah pernah dibuat"
		MessageValidationFailed = "Validasi gagal"
		MessageMustAbove18 = "Umur harus di atas 18 tahun"
	}
}
