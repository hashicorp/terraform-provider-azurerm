package migration

import (
	"context"
	"fmt"
	"html"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2021-08-01/apimanagement"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
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

		client := meta.(*clients.Client).ApiManagement.PolicyClient
		resp, err := client.Get(ctx, id.ResourceGroup, id.ServiceName, apimanagement.PolicyExportFormatXML)
		if err != nil {
			return nil, fmt.Errorf("making Read request for API Management Policy (Resource Group %q / API Management Service %q / API %q): %+v", id.ResourceGroup, id.ServiceName, id.Name, err)
		}

		if properties := resp.PolicyContractProperties; properties != nil {
			// when you submit an `xml_link` to the API, the API downloads this link and stores it as `xml_content`
			// as such there is no way to set `xml_link` and we'll let Terraform handle it
			rawState["xml_content"] = html.UnescapeString(*properties.Value)
		}
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
