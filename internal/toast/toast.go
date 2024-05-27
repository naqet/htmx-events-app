package toast

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	INFO    = "info"
	SUCCESS = "success"
	DANGER  = "danger"
)

func AddToast(w http.ResponseWriter, level, message string) error {
    if (level != INFO && level != SUCCESS && level != DANGER) {
        return fmt.Errorf("Level should be one of these: %s, %s, %s", INFO, SUCCESS, DANGER)
    }

	event := map[string]map[string]string{
		"showToast": {
			"level":   level,
			"message": message,
		},
	}

	header, err := json.Marshal(event)

	if err != nil {
		return err
	}

	w.Header().Add("HX-Trigger", string(header))
	return nil
}
