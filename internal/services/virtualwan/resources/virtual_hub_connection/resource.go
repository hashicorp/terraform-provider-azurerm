// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_hub_connection

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -resource-name virtual_hub_connection -properties "name" -compare-values "subscription_id:virtual_hub_id,resource_group_name:virtual_hub_id,virtual_hub_name:virtual_hub_id"

const virtualHubResourceName = "azurerm_virtual_hub"

const azureVirtualHubConnectionResourceName = "azurerm_virtual_hub_connection"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualHubConnectionCreateOrUpdate,
		Read:   resourceVirtualHubConnectionRead,
		Update: resourceVirtualHubConnectionCreateOrUpdate,
		Delete: resourceVirtualHubConnectionDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingIdentity(&virtualwans.HubVirtualNetworkConnectionId{}),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&virtualwans.HubVirtualNetworkConnectionId{}),
		},

		Schema: virtualHubConnectionSchema(),
	}
}
