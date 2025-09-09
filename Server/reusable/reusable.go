package reusable

import (
	"encoding/json"
	"time"
)

// Convert []uint8 to time.Time
func Uint8ToTime(uintTime []uint8) (time.Time, error) {
	// Convert []uint8 to string
	dateStr := string(uintTime)
	//fmt.Println("Date String:", dateStr)

	// Parse the string as a time.Time object
	decodedTime, err := time.Parse("2006-01-02 15:04:05", dateStr)
	if err != nil {
		return time.Time{}, err
	}

	// Return the parsed time
	return decodedTime, nil
}

// convert json to string
func UnmarshalJSONtoStringSlice(arg []byte) ([]string, error) {
	var destination []string
	err := json.Unmarshal(arg, &destination)
	if err != nil {
		return nil, err
	}
	return destination, nil
}
