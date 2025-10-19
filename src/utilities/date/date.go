package date

import (
	"errors"
	"time"
)

func GenerateFutureDate(value int, unit string) (time.Time, error) {
	now := time.Now()

	switch unit {
	case "seconds":
		return now.Add(time.Duration(value) * time.Second), nil
	case "minutes":
		return now.Add(time.Duration(value) * time.Minute), nil
	case "hours":
		return now.Add(time.Duration(value) * time.Hour), nil
	case "days":
		return now.AddDate(0, 0, value), nil
	case "months":
		return now.AddDate(0, value, 0), nil
	case "years":
		return now.AddDate(value, 0, 0), nil
	default:
		return time.Time{}, errors.New("unidade de tempo inválida")
	}
}
