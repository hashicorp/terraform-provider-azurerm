// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_gateway_api

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/gatewayapi"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/resources/api_management_gateway_api/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementGatewayApiCreate,
		Read:   resourceApiManagementGatewayApiRead,
		Delete: resourceApiManagementGatewayApiDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := gatewayapi.ParseGatewayApiID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.ApiManagementGatewayApiV0ToV1{},
		}),

		Schema: apiManagementGatewayApiSchema(),
	}
}
