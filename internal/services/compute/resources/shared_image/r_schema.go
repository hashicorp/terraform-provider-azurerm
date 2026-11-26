// Copyright IBM Corp. 2014, 2026
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
			Description:  "The name of the Shared Image.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.SharedImageName,
		},

		"gallery_name": {
			Description:  "The name of the Shared Image Gallery in which this Shared Image should exist.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.SharedImageGalleryName,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"architecture": {
			Description:  "CPU architecture supported by an OS.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(galleryimages.ArchitectureXSixFour),
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(galleryimages.PossibleValuesForArchitecture(), false),
		},

		"os_type": {
			Description:  "The type of Operating System present in this Shared Image.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(galleryimages.PossibleValuesForOperatingSystemTypes(), false),
		},

		"disk_types_not_allowed": {
			Description: "One or more Disk Types not allowed for the Image.",
			Type:        pluginsdk.TypeSet,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
				ValidateFunc: validation.StringInSlice([]string{
					string(galleryimageversions.StorageAccountTypeStandardLRS),
					string(galleryimageversions.StorageAccountTypePremiumLRS),
				}, false),
			},
		},

		"end_of_life_date": {
			Description:      "The end of life date in RFC3339 format of the Image.",
			Type:             pluginsdk.TypeString,
			Optional:         true,
			DiffSuppressFunc: suppress.RFC3339Time,
			ValidateFunc:     validation.IsRFC3339Time,
		},

		"hyper_v_generation": {
			Description:  "The generation of HyperV that the Virtual Machine used to create the Shared Image is based on.",
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
						Description:  "The Publisher Name for this Gallery Image.",
						Type:         pluginsdk.TypeString,
						ForceNew:     true,
						Required:     true,
						ValidateFunc: validate.SharedImageIdentifierAttribute(128),
					},
					"offer": {
						Description:  "The Offer Name for this Shared Image.",
						Type:         pluginsdk.TypeString,
						ForceNew:     true,
						Required:     true,
						ValidateFunc: validate.SharedImageIdentifierAttribute(64),
					},
					"sku": {
						Description:  "The Name of the SKU for this Gallery Image.",
						Type:         pluginsdk.TypeString,
						ForceNew:     true,
						Required:     true,
						ValidateFunc: validate.SharedImageIdentifierAttribute(64),
					},
				},
			},
		},

		"description": {
			Description: "A description of this Shared Image.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"eula": {
			Description: "The End User Licence Agreement for the Shared Image.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ForceNew:    true,
		},

		"purchase_plan": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "The Purchase Plan Name for this Shared Image.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"publisher": {
						Description:  "The Purchase Plan Publisher for this Gallery Image.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"product": {
						Description:  "The Purchase Plan Product for this Gallery Image.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"privacy_statement_uri": {
			Description: "The URI containing the Privacy Statement associated with this Shared Image.",
			Type:        pluginsdk.TypeString,
			ForceNew:    true,
			Optional:    true,
		},

		"max_recommended_vcpu_count": {
			Description:  "Maximum count of vCPUs recommended for the Image.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(1, 80),
		},

		"min_recommended_vcpu_count": {
			Description:  "Minimum count of vCPUs recommended for the Image.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(1, 80),
		},

		"max_recommended_memory_in_gb": {
			Description:  "Maximum memory in GB recommended for the Image.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(1, 640),
		},

		"min_recommended_memory_in_gb": {
			Description:  "Minimum memory in GB recommended for the Image.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(1, 640),
		},

		"release_note_uri": {
			Description: "The URI containing the Release Notes associated with this Shared Image.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"specialized": {
			Description: "Whether the Operating System used inside this Image has been specialized (for example, `sysprep` on Windows has been run).",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
		},

		"trusted_launch_supported": {
			Description:   "Whether supports creation of both Trusted Launch Virtual Machines and Gen2 Virtual Machines with standard security created from the Shared Image.",
			Type:          pluginsdk.TypeBool,
			Optional:      true,
			ForceNew:      true,
			ConflictsWith: []string{"trusted_launch_enabled", "confidential_vm_supported", "confidential_vm_enabled"},
		},

		"trusted_launch_enabled": {
			Description:   "Whether Trusted Launch has to be enabled for the Virtual Machine created from the Shared Image.",
			Type:          pluginsdk.TypeBool,
			Optional:      true,
			ForceNew:      true,
			ConflictsWith: []string{"trusted_launch_supported", "confidential_vm_supported", "confidential_vm_enabled"},
		},

		"confidential_vm_supported": {
			Description:   "Whether supports creation of both Confidential Virtual Machines and Gen2 Virtual Machines with standard security from a compatible Gen2 OS disk VHD or Gen2 Managed image.",
			Type:          pluginsdk.TypeBool,
			Optional:      true,
			ForceNew:      true,
			ConflictsWith: []string{"trusted_launch_supported", "trusted_launch_enabled", "confidential_vm_enabled"},
		},

		"confidential_vm_enabled": {
			Description:   "Whether Confidential Virtual Machines enabled. It will enable all the features of trusted, with higher confidentiality features for isolate machines or encrypted data. Available for Gen2 machines.",
			Type:          pluginsdk.TypeBool,
			Optional:      true,
			ForceNew:      true,
			ConflictsWith: []string{"trusted_launch_supported", "trusted_launch_enabled", "confidential_vm_supported"},
		},

		"accelerated_network_support_enabled": {
			Description: "Whether the Shared Image supports Accelerated Network.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
		},

		"hibernation_enabled": {
			Description: "Whether the Shared Image supports hibernation.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
		},

		"disk_controller_type_nvme_enabled": {
			Description: "Whether the Shared Image supports NVMe disks.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
		},

		"tags": commonschema.Tags(),
	}
}
