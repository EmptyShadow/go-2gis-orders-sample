package date

import (
	"time"
)

const DateFormat = "2006-01-02"

func FromString(str string) (time.Time, error) {
	return time.Parse(DateFormat, str)
}

func ToString(date time.Time) string {
	return date.Format(DateFormat)
}
