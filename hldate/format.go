package hldate

import "time"

func FormatDate(exp int64, format string) string {
	// Create a time.Time object from the POSIX timestamp in UTC location
	utcTime := time.Unix(exp, 0).UTC()
	// Format the time according to the desired layout YYYYMMDDHHMM
	formattedDate := utcTime.Format(format)
	return formattedDate
}
