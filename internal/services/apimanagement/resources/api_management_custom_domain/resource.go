// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_custom_domain

import (
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var apiManagementCustomDomainResourceName = "azurerm_api_management_custom_domain"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: apiManagementCustomDomainCreateUpdate,
		Read:   apiManagementCustomDomainRead,
		Update: apiManagementCustomDomainCreateUpdate,
		Delete: apiManagementCustomDomainDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.CustomDomainID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: apiManagementCustomDomainSchema(),
	}
}
