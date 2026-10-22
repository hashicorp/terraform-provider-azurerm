// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package tags

import (
	"fmt"
)

func TagValueToString(v any) (string, error) {
	switch value := v.(type) {
	case string:
		return value, nil
	case int:
		return fmt.Sprintf("%d", value), nil
	default:
		return "", fmt.Errorf("unknown tag type %T in tag value", value)
	}
}
