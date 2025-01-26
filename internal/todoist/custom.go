package todoist

import (
	"encoding/json"
	"strings"
	"time"
)

var _ json.Unmarshaler = (*DateOnlyTime)(nil)

type DateOnlyTime struct {
	time.Time
}

func (t *DateOnlyTime) UnmarshalJSON(b []byte) (err error) {
	date, err := time.Parse(time.DateOnly, strings.Trim(string(b), `"`))
	if err != nil {
		return err
	}

	t.Time = date
	return
}
