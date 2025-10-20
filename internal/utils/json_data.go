package utils

import (
	"fmt"
	"strings"
	"time"
)

const format = "01-2006"

type JSONDate struct{ time.Time }

func (t JSONDate) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return nil, nil
	}
	return []byte(fmt.Sprintf(`"%s"`, t.Time.Format(format))), nil
}

func (t *JSONDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "null" {
		return nil
	}

	time, err := time.Parse(format, s)
	if err != nil {
		return err
	}

	t.Time = time
	return nil
}
