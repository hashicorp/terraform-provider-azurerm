package schema

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/firewalls"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Plan struct {
	// TODO
}

// PlanSchema TODO returns the schema for a Plan Data block
func PlanSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				// TODO
			},
		},
	}
}

func FlattenPlanData(_ firewalls.PlanData) []Plan {
	return []Plan{}
}
