package racing_draft

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (
	ErrUnauthorized = errors.New("unauthorized to perform action")
)

func respToError(raw []byte) error {
	m := make(map[string]interface{})
	err := json.Unmarshal(raw, &m)
	if err != nil {
		return fmt.Errorf("server responded with poorly-understood response: %w\n%s", err, string(raw))
	}
	msg, ok := m["msg"]
	if !ok {
		return fmt.Errorf("invalid response from server - please contact an admin")
	}
	return fmt.Errorf("err - %s", msg)
}
