package data

import (
	"errors"
	"strconv"
	"time"
)

// Define an error that our UnmarshalJSON() method
// can return if we're unable to parse  or convert the JSON string successfully.
var ErrInvalidTimeFormat = errors.New("invalid Time format")

type CustomTime time.Time

func (ct CustomTime) MarshalJSON() ([]byte, error) {
	formattedTime := time.Time(ct).Format("2006-01-02 15:04:05")
	quotedJSONValue := strconv.Quote(formattedTime)
	return []byte(quotedJSONValue), nil
}

// Implement a UnmarshalJSON() method on the Runtime type so that it satisfies the json.Unmarshaler interface.
// IMPORTANT: Because UnmarshalJSON() needs to modify the receiver (our Runtime type),
// we must use a pointer receiver for this to work correctly.
// Otherwise, we will only be modifying a copy (which is then discarded when this method returns).
func (ct *CustomTime) UnmarshalJSON(jsonValue []byte) error {
	// We expect the incoming JSON value to be a string in the format "YYYY-MM-DD HH:MM:SS",
	// so we first remove the surrounding double-quotes from this string.
	unquotedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidTimeFormat
	}

	// Now, parse the unquoted JSON string into a time.Time value using a specified layout.
	// If the layout doesn't match the expected format, return the ErrInvalidTimeFormat error.
	const layout = "2006-01-02 15:04:05"
	parsedTime, err := time.Parse(layout, unquotedJSONValue)
	if err != nil {
		return ErrInvalidTimeFormat
	}

	// Convert the parsed time.Time value to the CustomTime type and assign it to the receiver.
	// Use the * operator to dereference the receiver (which is a pointer to CustomTime)
	// 		to set the underlying value of the pointer.
	*ct = CustomTime(parsedTime)
	return nil
}
