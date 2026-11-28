// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

func ConvertHexToBase64Encoding(str string) (string, error) {
	data, err := hex.DecodeString(str)
	if err != nil {
		return "", fmt.Errorf("converting %q from Hex to Base64 Encoding: %+v", str, err)
	}

	return base64.StdEncoding.EncodeToString(data), nil
}

func ConvertBase64ToHexEncoding(str string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", fmt.Errorf("converting %q from Base64 to Hex Encoding: %+v", str, err)
	}

	return hex.EncodeToString(data), nil
}
