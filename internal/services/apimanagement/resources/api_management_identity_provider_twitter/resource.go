// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_identity_provider_twitter

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/identityprovider"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementIdentityProviderTwitterCreateUpdate,
		Read:   resourceApiManagementIdentityProviderTwitterRead,
		Update: resourceApiManagementIdentityProviderTwitterCreateUpdate,
		Delete: resourceApiManagementIdentityProviderTwitterDelete,

		Importer: helpers.IdentityProviderImportFunc(identityprovider.IdentityProviderTypeTwitter),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: apiManagementIdentityProviderTwitterSchema(),
	}
}
