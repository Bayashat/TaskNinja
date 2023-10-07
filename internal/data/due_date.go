package data

import (
	"strconv"
	"time"
)

type CustomTime time.Time

func (ct CustomTime) MarshalJSON() ([]byte, error) {
	formattedTime := time.Time(ct).Format("2006-01-02 15:04:05")
	quotedJSONValue := strconv.Quote(formattedTime)
	return []byte(quotedJSONValue), nil
}
