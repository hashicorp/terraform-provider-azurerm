package schema

import (
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

// ignoreCaseDiffSuppressFunc is a DiffSuppressFunc from helper/schema that is
// used to ignore any case-changes in a return value.
func IgnoreCaseDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	return strings.ToLower(old) == strings.ToLower(new)
}

// ignoreCaseStateFunc is a StateFunc from helper/schema that converts the
// supplied value to lower before saving to state for consistency.
func IgnoreCaseStateFunc(val interface{}) string {
	return strings.ToLower(val.(string))
}
