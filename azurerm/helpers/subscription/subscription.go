package subscription

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func SubscriptionSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"subscription_id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},

		"display_name": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"state": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"location_placement_id": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"quota_id": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"spending_limit": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}
