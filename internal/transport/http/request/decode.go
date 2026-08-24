package http_request

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var requestValidator = validator.New()

func DecodeAndValidateRequest(r *http.Request, dto any) error {
	if err := json.NewDecoder(r.Body).Decode(dto); err != nil {
		return fmt.Errorf("decode json: %w", err)
	}

	if err := requestValidator.Struct(dto); err != nil {
		return fmt.Errorf("request validation: %w", err)
	}

	return nil
}
