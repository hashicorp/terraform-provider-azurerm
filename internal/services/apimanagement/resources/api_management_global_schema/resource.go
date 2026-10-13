// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_global_schema

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementGlobalSchemaCreateUpdate,
		Read:   resourceApiManagementGlobalSchemaRead,
		Update: resourceApiManagementGlobalSchemaCreateUpdate,
		Delete: resourceApiManagementGlobalSchemaDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := schema.ParseSchemaID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: apiManagementGlobalSchemaSchema(),
	}
}
