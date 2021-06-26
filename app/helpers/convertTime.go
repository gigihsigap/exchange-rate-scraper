package helpers

import (
	"time"
)

// Parse YYYY-MM-DD string into time.Time format
func YYYYMMDDToISODate(date string) (time.Time, error) {
	layout := "2006-01-02T15:04:05-0700"
	date = date + "T00:00:00-0000"
	dateTime, err := time.Parse(layout, date)
	return dateTime, err
}

// Parse time.Time data into YYYY-MM-DD string format
// func ISODateToYYYYMMDD(date time.Time) string {
// 	return ""
// }
