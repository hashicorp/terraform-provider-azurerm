// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_api_operation_policy

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apioperationpolicy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/resources/api_management_api_operation_policy/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -properties "resource_group_name,service_name:api_management_name,api_id:api_name,operation_id" -test-name rawXml

const azureApiManagementApiOperationPolicyResourceName = "azurerm_api_management_api_operation_policy"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Manages an API Management API Operation Policy.",
		Create:      resourceApiManagementAPIOperationPolicyCreateUpdate,
		Read:        resourceApiManagementAPIOperationPolicyRead,
		Update:      resourceApiManagementAPIOperationPolicyCreateUpdate,
		Delete:      resourceApiManagementAPIOperationPolicyDelete,
		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&apioperationpolicy.OperationId{}),
		},

		Importer: pluginsdk.ImporterValidatingIdentity(&apioperationpolicy.OperationId{}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		SchemaVersion: 2,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.ApiManagementApiOperationPolicyV0ToV1{},
			1: migration.ApiManagementApiOperationPolicyV1ToV2{},
		}),

		Schema: apiManagementApiOperationPolicySchema(),
	}
}
