package utils

import (
	"encoding/json"
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
