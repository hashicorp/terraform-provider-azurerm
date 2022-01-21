package migration

import (
	"context"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = ApiManagementApiPolicyV0ToV1{}

type ApiManagementApiPolicyV0ToV1 struct{}

func (ApiManagementApiPolicyV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		apiMgmtId, err := parse.ApiManagementID(rawState["id"].(string))
		if err != nil {
			return rawState, nil
		}
		id := parse.NewPolicyID(apiMgmtId.SubscriptionId, apiMgmtId.ResourceGroup, apiMgmtId.ServiceName, "policy")
		rawState["id"] = id.ID()
		return rawState, nil
	}
}

func (ApiManagementApiPolicyV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"api_management_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"xml_content": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"xml_link": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
	}
}
