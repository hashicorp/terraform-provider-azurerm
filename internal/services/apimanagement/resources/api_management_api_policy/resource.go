// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_api_policy

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apipolicy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/resources/api_management_api_policy/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -properties "api_id:api_name,service_name:api_management_name,resource_group_name" -test-expect-non-empty

const azurermApiManagementApiPolicy = "azurerm_api_management_api_policy"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Manages an API Management API Policy.",
		Create:      resourceApiManagementAPIPolicyCreateUpdate,
		Read:        resourceApiManagementAPIPolicyRead,
		Update:      resourceApiManagementAPIPolicyCreateUpdate,
		Delete:      resourceApiManagementAPIPolicyDelete,
		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&apipolicy.ApiId{}),
		},

		Importer: pluginsdk.ImporterValidatingIdentity(&apipolicy.ApiId{}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		SchemaVersion: 2,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.ApiManagementApiPolicyV0ToV1{},
			1: migration.ApiManagementApiPolicyV1ToV2{},
		}),

		Schema: apiManagementApiPolicySchema(),
	}
}
