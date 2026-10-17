// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_product_policy

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/productpolicy"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/resources/api_management_product_policy/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementProductPolicyCreateUpdate,
		Read:   resourceApiManagementProductPolicyRead,
		Update: resourceApiManagementProductPolicyCreateUpdate,
		Delete: resourceApiManagementProductPolicyDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := productpolicy.ParseProductID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		SchemaVersion: 2,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.ApiManagementProductPolicyV0ToV1{},
			1: migration.ApiManagementProductPolicyV1ToV2{},
		}),

		Schema: apiManagementProductPolicySchema(),
	}
}
