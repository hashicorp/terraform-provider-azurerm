// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package route_table

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/routetables"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity

var routeTableResourceName = "azurerm_route_table"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceRouteTableCreate,
		Read:   resourceRouteTableRead,
		Update: resourceRouteTableUpdate,
		Delete: resourceRouteTableDelete,

		Importer: pluginsdk.ImporterValidatingIdentity(&routetables.RouteTableId{}),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&routetables.RouteTableId{}),
		},

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: routeTableSchema(),
	}
}
