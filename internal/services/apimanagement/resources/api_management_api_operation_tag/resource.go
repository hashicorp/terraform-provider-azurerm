// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_api_operation_tag

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apioperationtag"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func RegisterResource() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Description: "Manages an API Management API Operation Tag.",
		Create:      resourceApiManagementApiOperationTagCreate,
		Read:        resourceApiManagementApiOperationTagRead,
		Delete:      resourceApiManagementApiOperationTagDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := apioperationtag.ParseOperationTagID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: apiManagementApiOperationTagSchema(),
	}

	if !features.SixPointOh() {
		resource.Update = resourceApiManagementApiOperationTagUpdate
		resource.Timeouts.Update = pluginsdk.DefaultTimeout(30 * time.Minute)

		resource.Schema["display_name"] = &pluginsdk.Schema{
			Description: "The display name of the API Management API Operation Tag.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			// NOTE: O+C when omitted the existing tag is assigned as is, so its display name is read back from Azure
			Computed:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Deprecated:   "`display_name` has been deprecated in favour of the `azurerm_api_management_tag` resource and will be removed in v6.0 of the AzureRM Provider",
		}
	}

	return resource
}
