// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_accelerator

import (
	appplatform_rm "github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/resources/spring_cloud_accelerator/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

func (s Resource) DeprecationMessage() string {
	return "Azure Spring Apps is now deprecated and will be retired on 2028-05-31 - as such the `azurerm_spring_cloud_accelerator` resource is deprecated and will be removed in a future major version of the AzureRM Provider. See https://aka.ms/asaretirement for more information."
}

var (
	_ sdk.Resource                                = Resource{}
	_ sdk.ResourceWithStateMigration              = Resource{}
	_ sdk.ResourceWithDeprecationAndNoReplacement = Resource{}
)

func (s Resource) ResourceType() string {
	return "azurerm_spring_cloud_accelerator"
}

func (s Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return appplatform_rm.ValidateApplicationAcceleratorID
}

func (s Resource) StateUpgraders() sdk.StateUpgradeData {
	return sdk.StateUpgradeData{
		SchemaVersion: 1,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.SpringCloudAcceleratorV0ToV1{},
		},
	}
}
