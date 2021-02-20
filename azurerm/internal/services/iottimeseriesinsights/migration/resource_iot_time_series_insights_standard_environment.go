package migration

import (
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func TimeSeriesInsightsAccessPolicyV0() schema.StateUpgrader {
	return schema.StateUpgrader{
		Type:    timeSeriesInsightsAccessPolicyV0StateMigration().CoreConfigSchema().ImpliedType(),
		Upgrade: timeSeriesInsightsAccessPolicyV0StateUpgradeV0ToV1,
		Version: 0,
	}
}

func timeSeriesInsightsAccessPolicyV0StateMigration() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"time_series_insights_environment_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"principal_object_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"roles": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func timeSeriesInsightsAccessPolicyV0StateUpgradeV0ToV1(rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	log.Println("[DEBUG] Migrating ResourceType from v0 to v1 format")
	oldId := rawState["id"].(string)
	newId := strings.Replace(oldId, "/accesspolicies/", "/accessPolicies/", 1)

	log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)

	rawState["id"] = newId

	return rawState, nil
}
