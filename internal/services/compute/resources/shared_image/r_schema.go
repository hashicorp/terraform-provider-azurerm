// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package shared_image

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleryimages"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-07-03/galleryimageversions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func sharedImageSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.SharedImageName,
		},

		"gallery_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.SharedImageGalleryName,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"architecture": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(galleryimages.ArchitectureXSixFour),
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(galleryimages.PossibleValuesForArchitecture(), false),
		},

		"os_type": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(galleryimages.PossibleValuesForOperatingSystemTypes(), false),
		},

		"disk_types_not_allowed": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
				ValidateFunc: validation.StringInSlice([]string{
					string(galleryimageversions.StorageAccountTypeStandardLRS),
					string(galleryimageversions.StorageAccountTypePremiumLRS),
				}, false),
			},
		},

		"end_of_life_date": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			DiffSuppressFunc: suppress.RFC3339Time,
			ValidateFunc:     validation.IsRFC3339Time,
		},

		"hyper_v_generation": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(galleryimages.HyperVGenerationVOne),
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(galleryimages.PossibleValuesForHyperVGeneration(), false),
		},

		"identifier": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"publisher": {
						Type:         pluginsdk.TypeString,
						ForceNew:     true,
						Required:     true,
						ValidateFunc: validate.SharedImageIdentifierAttribute(128),
					},
					"offer": {
						Type:         pluginsdk.TypeString,
						ForceNew:     true,
						Required:     true,
						ValidateFunc: validate.SharedImageIdentifierAttribute(64),
					},
					"sku": {
						Type:         pluginsdk.TypeString,
						ForceNew:     true,
						Required:     true,
						ValidateFunc: validate.SharedImageIdentifierAttribute(64),
					},
				},
			},
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"eula": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
		},

		"purchase_plan": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"publisher": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"product": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"privacy_statement_uri": {
			Type:     pluginsdk.TypeString,
			ForceNew: true,
			Optional: true,
		},

		"max_recommended_vcpu_count": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(1, 80),
		},

		"min_recommended_vcpu_count": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(1, 80),
		},

		"max_recommended_memory_in_gb": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(1, 640),
		},

		"min_recommended_memory_in_gb": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(1, 640),
		},

		"release_note_uri": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"specialized": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
		},

		"trusted_launch_supported": {
			Type:          pluginsdk.TypeBool,
			Optional:      true,
			ForceNew:      true,
			ConflictsWith: []string{"trusted_launch_enabled", "confidential_vm_supported", "confidential_vm_enabled"},
		},

		"trusted_launch_enabled": {
			Type:          pluginsdk.TypeBool,
			Optional:      true,
			ForceNew:      true,
			ConflictsWith: []string{"trusted_launch_supported", "confidential_vm_supported", "confidential_vm_enabled"},
		},

		"confidential_vm_supported": {
			Type:          pluginsdk.TypeBool,
			Optional:      true,
			ForceNew:      true,
			ConflictsWith: []string{"trusted_launch_supported", "trusted_launch_enabled", "confidential_vm_enabled"},
		},

		"confidential_vm_enabled": {
			Type:          pluginsdk.TypeBool,
			Optional:      true,
			ForceNew:      true,
			ConflictsWith: []string{"trusted_launch_supported", "trusted_launch_enabled", "confidential_vm_supported"},
		},

		"accelerated_network_support_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
		},

		"hibernation_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
		},

		"disk_controller_type_nvme_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
		},

		"tags": commonschema.Tags(),
	}
}
