// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package route_server_bgp_connection

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -resource-name route_server_bgp_connection -properties "name" -compare-values "subscription_id:route_server_id,resource_group_name:route_server_id,hub_name:route_server_id"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceRouteServerBgpConnectionCreate,
		Read:   resourceRouteServerBgpConnectionRead,
		Delete: resourceRouteServerBgpConnectionDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingIdentity(&commonids.VirtualHubBGPConnectionId{}),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&commonids.VirtualHubBGPConnectionId{}),
		},

		Schema: routeServerBgpConnectionSchema(),
	}
}
