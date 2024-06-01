package utils

import (
	"encoding/json"
	"fmt"
	"htmx-events-app/internal/chttp"
	"io"
	"net/http"
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

func GetEmailFromContext(r *http.Request) (string, error) {
	email, ok := r.Context().Value("email").(string)

	if !ok || email == "" {
		return email, fmt.Errorf("Senders email couldn't be obtained from the request")
	}
    return email, nil
}

func WriteJson(w http.ResponseWriter, value any) error {
	data, err := json.Marshal(value)

	if err != nil {
		return err
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
	return nil
}

func IsHtmxRequest(r *http.Request) bool {
	return r.Header.Get("HX-Request") != ""
}

func AddHtmxRedirect(w http.ResponseWriter, path string) {
	w.Header().Add("HX-Redirect", path)
}
