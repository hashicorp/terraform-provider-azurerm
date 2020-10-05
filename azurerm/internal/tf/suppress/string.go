package suppress

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func CaseDifference(_, old, new string, _ *schema.ResourceData) bool {
	return strings.EqualFold(old, new)
}
