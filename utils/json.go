package utils

import (
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
