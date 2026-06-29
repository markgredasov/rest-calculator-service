package utils_time

import (
	"encoding/json"
	"time"
)

type TimeOfDay string

func (t TimeOfDay) ToTime() (time.Time, error) {
	if t == "" {
		return time.Time{}, nil
	}
	return time.Parse("15:04", string(t))
}

func (t TimeOfDay) String() string {
	return string(t)
}

func (t TimeOfDay) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(t))
}

func (t *TimeOfDay) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	*t = TimeOfDay(str)
	return nil
}

func GetWeekday(date time.Time) int {
	day := int(date.Weekday())
	if day == 0 {
		return 7
	}
	return day
}
