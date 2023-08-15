// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"fmt"
	"html"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/policy"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = ApiManagementPolicyV0ToV1{}

type ApiManagementPolicyV0ToV1 struct{}

func (ApiManagementPolicyV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		apiMgmtId, err := policy.ParseServiceID(rawState["id"].(string))
		if err != nil {
			return rawState, nil
		}
		id := policy.NewServiceID(apiMgmtId.SubscriptionId, apiMgmtId.ResourceGroupName, apiMgmtId.ServiceName)
		rawState["id"] = id.ID()

		client := meta.(*clients.Client).ApiManagement.PolicyClient
		resp, err := client.Get(ctx, id, policy.GetOperationOptions{Format: pointer.To(policy.PolicyExportFormatXml)})
		if err != nil {
			return nil, fmt.Errorf("making Read request for API Management Policy (Resource Group %q / API Management Service %q / API %q): %+v", id.ResourceGroupName, id.ServiceName, "policy", err)
		}

		if model := resp.Model; model != nil {
			if props := model.Properties; props != nil {
				// when you submit an `xml_link` to the API, the API downloads this link and stores it as `xml_content`
				// as such there is no way to set `xml_link` and we'll let Terraform handle it
				rawState["xml_content"] = html.UnescapeString(props.Value)
			}
		}
		return rawState, nil
	}
}

func (ApiManagementPolicyV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return policySchemaForV0AndV1()
}

func policySchemaForV0AndV1() map[string]*pluginsdk.Schema {
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
