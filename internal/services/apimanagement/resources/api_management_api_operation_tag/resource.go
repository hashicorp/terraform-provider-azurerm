// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_api_operation_tag

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apioperationtag"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementApiOperationTagCreateUpdate,
		Read:   resourceApiManagementApiOperationTagRead,
		Update: resourceApiManagementApiOperationTagCreateUpdate,
		Delete: resourceApiManagementApiOperationTagDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := apioperationtag.ParseOperationTagID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: apiManagementApiOperationTagSchema(),
	}
}
