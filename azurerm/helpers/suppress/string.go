package suppress

import (
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func CaseDifference(_, old, new string, _ *schema.ResourceData) bool {
	return strings.ToLower(old) == strings.ToLower(new)
}

func IgnoreIfNotSet(_, old, new string, _ *schema.ResourceData) bool {
	retBool := false

	if new == "" {
		retBool = true
	} else {
		retBool = strings.EqualFold(old, new)
	}

	return retBool
}
