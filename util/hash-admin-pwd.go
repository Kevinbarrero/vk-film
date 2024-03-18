package util

func HashAdminPassword(password string) (string, error) {
	return HashPassword(password)
}
