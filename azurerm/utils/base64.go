package utils

import "encoding/base64"

// Base64EncodeIfNot encodes data if the input isn't already encoded using
// base64.StdEncoding.EncodeToString. If the input is already base64 encoded,
// return the original input unchanged.
func Base64EncodeIfNot(data string) string {
	// Check whether the data is already Base64 encoded; don't double-encode
	if base64IsEncoded(data) {
		return data
	}
	// data has not been encoded encode and return
	return base64.StdEncoding.EncodeToString([]byte(data))
}

func base64IsEncoded(data string) bool {
	_, err := base64.StdEncoding.DecodeString(data)
	return err == nil
}
