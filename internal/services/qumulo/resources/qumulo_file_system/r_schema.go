// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package qumulo_file_system

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/qumulo/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of this Azure Native Qumulo Scalable File System resource.",
			ForceNew:     true,
			Required:     true,
			Type:         pluginsdk.TypeString,
			ValidateFunc: validate.FileSystemName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"admin_password": {
			Description:  "The initial administrator password of the Azure Native Qumulo Scalable File System.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			Sensitive:    true,
			ValidateFunc: validate.ValidatePasswordComplexity,
		},

		"email": {
			Description:  "The email address used for the Azure Native Qumulo Scalable File System.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsEmailAddress,
		},

		"storage_sku": {
			Description: "The storage SKU.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			// Use enum in SDK if swagger issue is resolved https://github.com/Azure/azure-rest-api-specs/issues/34017
			ValidateFunc: validation.StringInSlice([]string{
				"Cold_LRS",
				"Hot_LRS",
				"Hot_ZRS",
			}, false),
		},

		"subnet_id": {
			Description:  "The delegated subnet ID for VNet injection.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"zone": commonschema.ZoneSingleRequiredForceNew(),

		"offer_id": {
			Description:  "The marketplace offer ID.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      "qumulo-saas-mpp",
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"plan_id": {
			Description:  "The marketplace plan ID.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      "azure-native-qumulo-v3",
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"publisher_id": {
			Description:  "The marketplace publisher ID.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      "qumulo1584033880660",
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
