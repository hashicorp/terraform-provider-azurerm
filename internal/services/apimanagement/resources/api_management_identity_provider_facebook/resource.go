// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_identity_provider_facebook

import (
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/helpers"

	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/identityprovider"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementIdentityProviderFacebookCreateUpdate,
		Read:   resourceApiManagementIdentityProviderFacebookRead,
		Update: resourceApiManagementIdentityProviderFacebookCreateUpdate,
		Delete: resourceApiManagementIdentityProviderFacebookDelete,

		Importer: helpers.IdentityProviderImportFunc(identityprovider.IdentityProviderTypeFacebook),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: apiManagementIdentityProviderFacebookSchema(),
	}
}
