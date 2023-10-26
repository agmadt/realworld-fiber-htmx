package internal

func ErrorMessage(fieldName string, fieldType string) string {

	message := ""

	switch fieldType {
	case "required":
		message = fieldName + " is required."
	case "email":
		message = fieldName + " must be an email."
	default:
		message = ""
	}

	return message
}
