package state

import "strings"

// IgnoreCase is a StateFunc from helper/schema that converts the
// supplied value to lower before saving to state for consistency.
func IgnoreCase(val interface{}) string {
	return strings.ToLower(val.(string))
}
