// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"encoding/json"
	"fmt"
)

// NormalizeJson takes a JSON string (as an interface{}), parses it, and returns it
// as a compact, normalized string. If the input is not a string, it returns an error message.
func NormalizeJson(jsonString interface{}) string {
	if jsonString == nil || jsonString == "" {
		return ""
	}

	str, ok := jsonString.(string)
	if !ok {
		return fmt.Sprintf("Error parsing JSON: expected string, got %T", jsonString)
	}

	var j interface{}

	if err := json.Unmarshal([]byte(str), &j); err != nil {
		return fmt.Sprintf("Error parsing JSON: %+v", err)
	}
	b, _ := json.Marshal(j)
	return string(b)
}
