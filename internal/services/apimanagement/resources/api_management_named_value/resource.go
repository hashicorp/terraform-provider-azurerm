// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_named_value

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/namedvalue"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -properties "named_value_id:name,resource_group_name,service_name:api_management_name"

const azureApiManagementNamedValueResourceName = "azurerm_api_management_named_value"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Manages an API Management Named Value.",
		Create:      resourceApiManagementNamedValueCreateUpdate,
		Read:        resourceApiManagementNamedValueRead,
		Update:      resourceApiManagementNamedValueCreateUpdate,
		Delete:      resourceApiManagementNamedValueDelete,
		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&namedvalue.NamedValueId{}),
		},

		Importer: pluginsdk.ImporterValidatingIdentity(&namedvalue.NamedValueId{}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: apiManagementNamedValueSchema(),
	}
}
