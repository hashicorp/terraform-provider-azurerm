package suppress

import (
	"github.com/hashicorp/terraform/helper/schema"
	"strings"
)

func CaseDifference(_, old, new string, _ *schema.ResourceData) bool {
	return strings.ToLower(old) == strings.ToLower(new)
}
