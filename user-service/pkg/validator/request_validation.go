import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func Validate(data interface{}) error {
	var errorMessages []string

	err := validate.Struct(data)
	if err != nil {

		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return err
		}

		for _, e := range validationErrors {
			switch e.Tag() {

			case "required":
				errorMessages = append(errorMessages,
					fmt.Sprintf("%s is required", e.Field()))

			case "email":
				errorMessages = append(errorMessages,
					fmt.Sprintf("%s is not a valid email", e.Field()))

			case "min":
				errorMessages = append(errorMessages,
					fmt.Sprintf("%s must be at least %s characters long",
						e.Field(), e.Param()))

			case "max":
				errorMessages = append(errorMessages,
					fmt.Sprintf("%s must be at most %s characters long",
						e.Field(), e.Param()))
			}
		}

		return errors.New("Validasi gagal: " + joinMessage(errorMessages))
	}

	return nil
}

func joinMessage(errorMessages []string) string {
	result := ""
	for i, message := renge errorMessages {
		if i > 0 {
			result +=","
		}
		result += message
	}
	return result
}