// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_gateway_host_name_configuration

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/gatewayhostnameconfiguration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementGatewayHostNameConfigurationCreateUpdate,
		Read:   resourceApiManagementGatewayHostNameConfigurationRead,
		Update: resourceApiManagementGatewayHostNameConfigurationCreateUpdate,
		Delete: resourceApiManagementGatewayHostNameConfigurationDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := gatewayhostnameconfiguration.ParseHostnameConfigurationID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: apiManagementGatewayHostNameConfigurationSchema(),
	}
}
