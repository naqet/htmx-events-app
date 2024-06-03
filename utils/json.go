package utils

import (
	"encoding/json"
	"errors"
	"time"
)

type Time struct {
	time.Time
}

func (t *Time) UnmarshalJSON(b []byte) error {
	if len(b) < 1 {
		return errors.New("Invalid date")
	}

	val, err := time.Parse("2006-01-02T15:04", string(b[1:len(b)-1]))
	if err != nil {
		return err
	}
	*t = Time{val}
	return nil
}

type StringArr struct {
	Entries []string
}

func (h *StringArr) UnmarshalJSON(b []byte) error {
	var raw interface{}

	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	switch v := raw.(type) {
	case string:
		h.Entries = []string{v}
	case []interface{}:
		for _, item := range v {
			str, ok := item.(string)
			if !ok {
				return errors.New("Invalid entry in hosts")
			}

			h.Entries = append(h.Entries, str)
		}
	default:
		return errors.New("Invalid hosts")
	}

	return nil
}
