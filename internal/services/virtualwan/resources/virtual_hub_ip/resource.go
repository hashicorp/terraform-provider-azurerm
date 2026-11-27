// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_hub_ip

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -resource-name virtual_hub_ip -properties "name" -compare-values "subscription_id:virtual_hub_id,resource_group_name:virtual_hub_id,virtual_hub_name:virtual_hub_id"

const virtualHubResourceName = "azurerm_virtual_hub"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualHubIPCreate,
		Read:   resourceVirtualHubIPRead,
		Update: resourceVirtualHubIPUpdate,
		Delete: resourceVirtualHubIPDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingIdentity(&commonids.VirtualHubIPConfigurationId{}),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&commonids.VirtualHubIPConfigurationId{}),
		},

		Schema: virtualHubIpSchema(),
	}
}
