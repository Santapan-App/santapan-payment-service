package repository

import (
	"encoding/base64"
	"strconv"
)

// DecodeCursor decodes a base64 encoded `id` (int64) back to its original value.
func DecodeCursor(encodedID string) (int64, error) {
	// Decode the base64 string
	byt, err := base64.StdEncoding.DecodeString(encodedID)
	if err != nil {
		return 0, err
	}

	// Convert the decoded bytes into an int64 string
	idStr := string(byt)

	// Parse the string to int64
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// EncodeCursor encodes an `id` (int64) into a base64 string.
func EncodeCursor(id int64) (string, error) {
	// Convert the int64 id to string
	idStr := strconv.FormatInt(id, 10)

	// Encode the string into base64
	encodedID := base64.StdEncoding.EncodeToString([]byte(idStr))

	return encodedID, nil
}
