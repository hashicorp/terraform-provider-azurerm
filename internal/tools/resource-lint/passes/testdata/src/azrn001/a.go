package azrn001

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func validCases() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// Valid: Uses _percentage suffix
		"utilization_percentage": {
			Type:     schema.TypeInt,
			Optional: true,
		},
	}
}

func invalidCases() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// Invalid: Uses _in_percent instead of _percentage
		"cpu_in_percent": { // want `AZRN001`
			Type:     schema.TypeInt,
			Optional: true,
		},
	}
}
