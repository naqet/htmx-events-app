package utils

import (
	"encoding/json"
	"htmx-events-app/internal/chttp"
	"io"
)

func GetDataFromBody(body io.Reader, dst any) error {
	content, err := io.ReadAll(body)

	if err != nil {
		return err
	}

	err = json.Unmarshal(content, dst)

	if err != nil {
		return chttp.BadRequestError()
	}

    return nil
}
