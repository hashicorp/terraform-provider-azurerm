package migration

import (
	"context"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = ApiManagementGatewayApiV0ToV1{}

type ApiManagementGatewayApiV0ToV1 struct{}

func (ApiManagementGatewayApiV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"api_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"gateway_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
	}
}

// UpgradeFunc this is a noop migration to account for a migration that was accidentally added from github.com/hashicorp/terraform-provider-azurerm/pull/22783/
// That migration didn't do anything for this resource so we'll just swap it for a no-op migration here
func (ApiManagementGatewayApiV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		return rawState, nil
	}
}
