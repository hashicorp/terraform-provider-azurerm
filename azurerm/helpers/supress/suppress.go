package supress

import (
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"strings"
)

func CaseDifference(k, old, new string, d *schema.ResourceData) bool {
	return strings.ToLower(old) == strings.ToLower(new)
}

func Rfc3339Time(k, old, new string, d *schema.ResourceData) bool {
	ot, oerr := time.Parse(time.RFC3339, old)
	nt, nerr := time.Parse(time.RFC3339, new)

	if oerr != nil || nerr != nil {
		return false
	}

	return nt.Equal(ot)
}
