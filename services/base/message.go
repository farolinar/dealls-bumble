package servicebase

var (
	ID_LANG = "id"
)

var (
	MessageSuccess          = "Success"
	MessageInternalError    = "Internal server error"
	MessageFailedDecodeJSON = "Failed to decode JSON"
	MessagePasswordInvalid  = "Minimum eight characters, at least one uppercase letter, one lowercase letter, one number, and one special character"
)

func Translate(lang string) {
	switch lang {
	case ID_LANG:
		MessageSuccess = "Sukses"
		MessageInternalError = "Terjadi kegagalan pada server"
		MessageFailedDecodeJSON = "Minimum 8 karakter, satu huruf kapital, satu huruf kecil, satu angka, dan satu karakter spesial"
	}
}
