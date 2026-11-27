// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package app_service_connection

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicelinker/2024-04-01/servicelinker"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/serviceconnector/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"app_service_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.AsGeneratedID(commonids.ParseAppServiceIDInsensitively),
		},

		"target_resource_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: azure.ValidateResourceID,
		},

		"client_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(servicelinker.ClientTypeNone),
			ValidateFunc: validation.StringInSlice([]string{
				string(servicelinker.ClientTypeNone),
				string(servicelinker.ClientTypeDotnet),
				string(servicelinker.ClientTypeJava),
				string(servicelinker.ClientTypePython),
				string(servicelinker.ClientTypeGo),
				string(servicelinker.ClientTypePhp),
				string(servicelinker.ClientTypeRuby),
				string(servicelinker.ClientTypeDjango),
				string(servicelinker.ClientTypeNodejs),
				string(servicelinker.ClientTypeSpringBoot),
			}, false),
		},

		"secret_store": helpers.SecretStoreSchema(),

		"vnet_solution": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(servicelinker.PossibleValuesForVNetSolutionType(), false),
		},

		"authentication": helpers.AuthInfoSchema(),
	}
}

func (r Resource) Attributes() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{}
}
