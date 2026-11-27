// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_gateway_route_config

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/resources/spring_cloud_gateway_route_config/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		DeprecationMessage: "Azure Spring Apps is now deprecated and will be retired on 2028-05-31 - as such the `azurerm_spring_cloud_gateway_route_config` resource is deprecated and will be removed in a future major version of the AzureRM Provider. See https://aka.ms/asaretirement for more information.",

		Create: resourceSpringCloudGatewayRouteConfigCreateUpdate,
		Read:   resourceSpringCloudGatewayRouteConfigRead,
		Update: resourceSpringCloudGatewayRouteConfigCreateUpdate,
		Delete: resourceSpringCloudGatewayRouteConfigDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.SpringCloudGatewayRouteConfigV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			// the ID is parsed insensitively to preserve the behaviour of the legacy parser this replaced -
			// we are ok with this remaining insensitive as Azure Spring Apps is deprecated and will be removed
			_, err := appplatform.ParseRouteConfigIDInsensitively(id)
			return err
		}),

		Schema: springCloudGatewayRouteConfigSchema(),
	}
}
