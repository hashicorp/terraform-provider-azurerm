// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_openid_connect_provider

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/openidconnectprovider"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementOpenIDConnectProviderCreateUpdate,
		Read:   resourceApiManagementOpenIDConnectProviderRead,
		Update: resourceApiManagementOpenIDConnectProviderCreateUpdate,
		Delete: resourceApiManagementOpenIDConnectProviderDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := openidconnectprovider.ParseOpenidConnectProviderID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: apiManagementOpenidConnectProviderSchema(),
	}
}
