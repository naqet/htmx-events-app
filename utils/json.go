package utils

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"
)

const DATETIME_LOCAL = "2006-01-02T15:04"

type Time time.Time

func (t *Time) UnmarshalJSON(b []byte) error {
	if len(b) < 1 {
		return errors.New("Invalid date")
	}

	val, err := time.Parse(DATETIME_LOCAL, string(b[1:len(b)-1]))
	if err != nil {
		return err
	}
	*t = Time(val)
	return nil
}

type TimeArr []Time

func (t *TimeArr) UnmarshalJSON(b []byte) error {
	var raw interface{}

	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	switch v := raw.(type) {
	case string:
		var cTime Time
		err := json.Unmarshal([]byte(`"`+v+`"`), &cTime)

		if err != nil {
			return err
		}
		*t = []Time{cTime}
	case []interface{}:
		for _, item := range v {
			str, ok := item.(string)
			if !ok {
				return errors.New("Invalid entry in array")
			}

			var cTime Time
			err := json.Unmarshal([]byte(`"`+str+`"`), &cTime)

			if err != nil {
				return err
			}
			*t = append(*t, cTime)
		}
	default:
		return errors.New("Invalid time array")
	}

	return nil
}

type StringArr []string

func (h *StringArr) UnmarshalJSON(b []byte) error {
	var raw interface{}

	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	switch v := raw.(type) {
	case string:
		*h = []string{v}
	case []interface{}:
		for _, item := range v {
			str, ok := item.(string)
			if !ok {
				return errors.New("Invalid entry in array")
			}

			*h = append(*h, str)
		}
	default:
		return errors.New("Invalid string array")
	}

	return nil
}

type IntArr []int

func (h *IntArr) UnmarshalJSON(b []byte) error {
	var raw interface{}

	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	switch v := raw.(type) {
	case string:
		val, err := strconv.Atoi(v)

		if err != nil {
			return err
		}

		*h = []int{val}
	case []interface{}:
		for _, item := range v {
			strInt, ok := item.(string)

			if !ok {
				return errors.New("Invalid entry in array")
			}

			val, err := strconv.Atoi(strInt)

			if err != nil {
				return err
			}

			*h = append(*h, val)
		}
	default:
		return errors.New("Invalid int array")
	}

	return nil
}

type Float64Arr []float64

func (h *Float64Arr) UnmarshalJSON(b []byte) error {
	var raw interface{}

	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	switch v := raw.(type) {
	case string:
		val, err := strconv.ParseFloat(v, 64)

		if err != nil {
			return err
		}

		*h = []float64{float64(val)}
	case []interface{}:
		for _, item := range v {
			strInt, ok := item.(string)

			if !ok {
				return errors.New("Invalid entry in array")
			}

			val, err := strconv.ParseFloat(strInt, 64)

			if err != nil {
				return err
			}

			*h = append(*h, val)
		}
	default:
		return errors.New("Invalid float64 array")
	}

	return nil
}
