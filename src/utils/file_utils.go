package utils

func IsFileImage(fileType string) bool {
	switch fileType {
	case "image/jpeg", "image/jpg", "image/gif", "image/png":
		return true
	default:
		return false
	}
}
