// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package compute_fleet

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurefleet/2024-11-01/fleets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/capacityreservationgroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/images"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleryimages"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-07-03/galleryimageversions"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/computefleet/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			// From Azure REST API, `name` can be up to 80 characters long. It must begin with a word character, and it must end with a word character or with '_'. The name may contain word characters or '.', '-', '_'. Instead, the validation is implemented according to Azure portal behavior
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile("^[a-zA-Z0-9][a-zA-Z0-9.-]{0,62}[a-zA-Z0-9]$|^[a-zA-Z0-9]$"),
				"The fleet name can only start or end with a number or a letter, and can contain only letters, numbers, periods (.), hyphens (-), up to 64 characters",
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"virtual_machine_profile": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"network_interface": helpers.NetworkInterfaceSchema(),

					"os_profile": helpers.OsProfileSchema(),

					"boot_diagnostic_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						ForceNew: true,
						Default:  false,
					},

					"boot_diagnostic_storage_account_endpoint": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.IsURLWithHTTPorHTTPS,
					},

					"capacity_reservation_group_id": commonschema.ResourceIDReferenceOptionalForceNew(&capacityreservationgroups.CapacityReservationGroupId{}),

					"data_disk": helpers.StorageProfileDataDiskSchema(),

					"encryption_at_host_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						ForceNew: true,
						Default:  false,
					},

					"extension": helpers.ExtensionSchema(),

					"extension_operations_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						ForceNew: true,
						Default:  true,
					},

					"extensions_time_budget_duration": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.ISO8601DurationBetween("PT15M", "PT2H"),
					},

					"gallery_application": helpers.GalleryApplicationSchema(),

					"license_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
						ValidateFunc: validation.StringInSlice([]string{
							"RHEL_BYOS",
							"SLES_BYOS",
							"Windows_Client",
							"Windows_Server",
						}, false),
					},

					"network_api_version": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
						// `Default` set according to portal behavior
						Default: string(fleets.NetworkApiVersionTwoZeroTwoZeroNegativeOneOneNegativeZeroOne),
						ValidateFunc: validation.StringMatch(
							regexp.MustCompile(`^\d{4}-\d{2}-\d{2}(-preview)?$`),
							"`network_api_version` must be in the format `YYYY-MM-DD` with an optional `-preview` suffix",
						),
					},

					"os_disk": helpers.StorageProfileOsDiskSchema(),

					"scheduled_event_os_image_timeout_duration": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
						ValidateFunc: validation.StringInSlice([]string{
							"PT15M",
						}, false),
					},

					"scheduled_event_termination_timeout_duration": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.ISO8601DurationBetween("PT5M", "PT15M"),
					},

					"secure_boot_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						ForceNew: true,
						Default:  false,
					},

					"source_image_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
						ValidateFunc: validation.Any(
							images.ValidateImageID,
							validation.AsGeneratedID(galleryimages.ParseGalleryImageIDInsensitively),
							validation.AsGeneratedID(galleryimageversions.ParseImageVersionIDInsensitively),
							validate.CommunityGalleryImageID,
							validate.CommunityGalleryImageVersionID,
							validate.SharedGalleryImageID,
							validate.SharedGalleryImageVersionID,
						),
						ExactlyOneOf: []string{
							"virtual_machine_profile.0.source_image_reference",
							"virtual_machine_profile.0.source_image_id",
						},
					},

					"source_image_reference": helpers.StorageProfileSourceImageReferenceSchema(),

					"user_data_base64": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
						ValidateFunc: validation.All(
							validation.StringIsBase64,
							validation.StringLenBetween(1, 349528),
						),
					},

					"vtpm_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						ForceNew: true,
						Default:  false,
					},
				},
			},
		},

		"virtual_machine_sizes_profile": {
			Type:     pluginsdk.TypeList,
			Required: true,
			// In Azure portal, maximum number of items is 15 while in Azure REST API, it is 10
			MaxItems: 10,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"rank": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntAtLeast(0),
					},
				},
			},
		},

		"additional_capabilities": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"hibernation_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						ForceNew: true,
						Default:  false,
						AtLeastOneOf: []string{
							"additional_capabilities.0.hibernation_enabled",
							"additional_capabilities.0.ultra_ssd_enabled",
						},
					},
					"ultra_ssd_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						ForceNew: true,
						Default:  false,
						AtLeastOneOf: []string{
							"additional_capabilities.0.hibernation_enabled",
							"additional_capabilities.0.ultra_ssd_enabled",
						},
					},
				},
			},
		},

		"compute_api_version": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// NOTE: O+C When creating new resource without providing this property, Azure will return a default value. As the default value is not documented anywhere, user experience will be poor to make this a `Required` property
			Computed: true,
			// This property is tested to be non-updatable
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^\d{4}-\d{2}-\d{2}(-preview)?$`),
				"`compute_api_version` must be in the format `YYYY-MM-DD` with an optional `-preview` suffix",
			),
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"on_demand_capacity": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			MaxItems:     1,
			AtLeastOneOf: []string{"on_demand_capacity", "spot_capacity"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"target_capacity": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(0, 10000),
					},

					"allocation_strategy": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(fleets.RegularPriorityAllocationStrategyLowestPrice),
						ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForRegularPriorityAllocationStrategy(), false),
					},

					"minimum_starting_capacity": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      0,
						ValidateFunc: validation.IntAtLeast(0),
					},
				},
			},
		},

		"plan": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						ForceNew:     true,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"product": {
						Type:         pluginsdk.TypeString,
						ForceNew:     true,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"publisher": {
						Type:         pluginsdk.TypeString,
						ForceNew:     true,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"promotion_code": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"platform_fault_domain_count": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      1,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(1, 2),
		},

		"spot_capacity": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			MaxItems:     1,
			AtLeastOneOf: []string{"on_demand_capacity", "spot_capacity"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					// `maintain_capacity_enabled` is `Required` as multiple properties depend on its value
					// `maintain_capacity_enabled` can be `true` using Azure REST API even if less than 3 `virtual_machine_sizes_profile` blocks are defined, which is different from Azure portal behavior
					// `virtual_machine_sizes_profile` can be updated in Azure portal even if `maintain_capacity_enabled` is `false`
					"maintain_capacity_enabled": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},

					"target_capacity": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(0, 10000),
					},

					"allocation_strategy": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(fleets.SpotAllocationStrategyPriceCapacityOptimized),
						ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForSpotAllocationStrategy(), false),
					},

					"eviction_policy": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(fleets.EvictionPolicyDelete),
						ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForEvictionPolicy(), false),
					},

					"maximum_hourly_price_per_virtual_machine_usd": {
						Type:         pluginsdk.TypeFloat,
						Optional:     true,
						Default:      -1,
						ValidateFunc: validate.SpotMaxPrice,
					},

					"minimum_capacity": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      0,
						ValidateFunc: validation.IntAtLeast(0),
					},
				},
			},
		},

		"zones": commonschema.ZonesMultipleOptionalForceNew(),

		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"unique_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var state ComputeFleetResourceModel
			if err := metadata.DecodeDiff(&state); err != nil {
				return fmt.Errorf("DecodeDiff: %+v", err)
			}

			if v := state.Identity; len(v) > 0 {
				// whilst the Swagger defines multiple at this time only UAI is supported
				if v[0].Type != identity.TypeUserAssigned {
					return fmt.Errorf("the type `%s` of `identity` is not supported currently", v[0].Type)
				}
			}

			if len(state.SpotCapacity) > 0 && len(state.OnDemandCapacity) > 0 {
				if state.SpotCapacity[0].TargetCapacity+state.OnDemandCapacity[0].TargetCapacity > 10000 {
					return errors.New("the sum of `spot_capacity.0.target_capacity` and `on_demand_capacity.0.target_capacity` must be between `0` and `10000`, inclusive")
				}
			}

			if len(state.SpotCapacity) > 0 {
				if state.SpotCapacity[0].MaintainCapacityEnabled {
					if state.SpotCapacity[0].MinimumCapacity > 0 {
						return errors.New("`spot_capacity.0.minimum_capacity` cannot be set if `spot_capacity.0.maintain_capacity_enabled` is set to `true`")
					}

					if state.SpotCapacity[0].EvictionPolicy != string(fleets.EvictionPolicyDelete) {
						return errors.New("`spot_capacity.0.eviction_policy` must be `Delete` if `spot_capacity.0.maintain_capacity_enabled` is `true`")
					}

					if len(state.VirtualMachineSizesProfile) < 3 {
						return errors.New("the number of `virtual_machine_sizes_profile` blocks must be at least 3 if `spot_capacity.0.maintain_capacity_enabled` is `true`")
					}

					if len(state.Zones) == 0 {
						return errors.New("enabling `spot_capacity.0.maintain_capacity_enabled` requires all qualified availability zones in the region to be supported")
					}
				}

				if state.SpotCapacity[0].MinimumCapacity > state.SpotCapacity[0].TargetCapacity {
					return errors.New("`spot_capacity.0.minimum_capacity` must be between `0` and `spot_capacity.0.target_capacity`, inclusive")
				}
			}

			// Relax `ForceNew` behavior of `spot_capacity` properties to allow updating resource from without to with `spot_capacity` properties
			// Apart from this, other changes are not allowed
			forceNewSpotCapacityProperties := []string{
				"maintain_capacity_enabled",
				"allocation_strategy",
				"eviction_policy",
				"maximum_hourly_price_per_virtual_machine_usd",
				"minimum_capacity",
			}
			oldSpotCapacity, _ := metadata.ResourceDiff.GetChange("spot_capacity")
			if len(oldSpotCapacity.([]any)) > 0 {
				if len(state.SpotCapacity) == 0 {
					metadata.ResourceDiff.ForceNew("spot_capacity")
				}

				for _, spotCapacityProperty := range forceNewSpotCapacityProperties {
					spotCapacityProperty = fmt.Sprintf("spot_capacity.0.%s", spotCapacityProperty)
					if metadata.ResourceDiff.HasChange(spotCapacityProperty) {
						metadata.ResourceDiff.ForceNew(spotCapacityProperty)
					}
				}
			}

			if len(state.OnDemandCapacity) > 0 {
				if state.OnDemandCapacity[0].MinimumStartingCapacity > state.OnDemandCapacity[0].TargetCapacity {
					return errors.New("`on_demand_capacity.0.minimum_starting_capacity` must be between `0` and `on_demand_capacity.0.target_capacity`, inclusive")
				}
			}

			// Relax `ForceNew` behavior of `on_demand_capacity` properties to allow updating resource from without to with `on_demand_capacity` properties
			// Apart from this, other changes are not allowed
			forceNewOnDemandCapacityProperties := []string{
				"allocation_strategy",
				"minimum_starting_capacity",
			}
			oldOnDemandCapacity, _ := metadata.ResourceDiff.GetChange("on_demand_capacity")
			if len(oldOnDemandCapacity.([]any)) > 0 {
				if len(state.OnDemandCapacity) == 0 {
					metadata.ResourceDiff.ForceNew("on_demand_capacity")
				}

				for _, onDemandCapacityProperty := range forceNewOnDemandCapacityProperties {
					onDemandCapacityProperty = fmt.Sprintf("on_demand_capacity.0.%s", onDemandCapacityProperty)
					if metadata.ResourceDiff.HasChange(onDemandCapacityProperty) {
						metadata.ResourceDiff.ForceNew(onDemandCapacityProperty)
					}
				}
			}

			if v := state.VirtualMachineProfile[0].DataDisks; len(v) > 0 {
				ultraSSDEnabled := false
				if ac := state.AdditionalCapabilities; len(ac) > 0 {
					ultraSSDEnabled = ac[0].UltraSsdEnabled
				}

				for i, dataDisk := range v {
					storageAccountType := dataDisk.StorageAccountType
					if !ultraSSDEnabled && storageAccountType == string(fleets.StorageAccountTypesUltraSSDLRS) {
						return errors.New("`UltraSSD_LRS` storage account type can be used only when `ultra_ssd_enabled` is `true`")
					}

					if dataDisk.CreateOption == string(fleets.DiskCreateOptionTypesEmpty) {
						if dataDisk.DiskSizeInGiB == 0 {
							return fmt.Errorf("`virtual_machine_profile.0.data_disk.%d.disk_size_in_gib` is required when `create_option` is `Empty`", i)
						}

						lunExist := metadata.ResourceDiff.GetRawConfig().AsValueMap()["virtual_machine_profile"].AsValueSlice()[0].AsValueMap()["data_disk"].AsValueSlice()[i].AsValueMap()["lun"]
						if lunExist.IsNull() {
							return fmt.Errorf("`virtual_machine_profile.0.data_disk.%d.lun` is required when `create_option` is `Empty`", i)
						}
					}

					if dataDisk.Caching != "" && dataDisk.DiskSizeInGiB > 4095 {
						return fmt.Errorf("`virtual_machine_profile.0.data_disk.%d.disk_size_in_gib` cannot be greater than 4095 GiB when `caching` is specified", i)
					}
				}
			}

			for i, virtualMachineSizesProfile := range state.VirtualMachineSizesProfile {
				if !metadata.ResourceDiff.GetRawConfig().AsValueMap()["virtual_machine_sizes_profile"].AsValueSlice()[i].AsValueMap()["rank"].IsNull() && (len(state.OnDemandCapacity) == 0 || state.OnDemandCapacity[0].AllocationStrategy != string(fleets.RegularPriorityAllocationStrategyPrioritized)) {
					return errors.New("`on_demand_capacity` should be specified with `allocation_strategy` equal to `Prioritized` when `virtual_machine_profile.0.virtual_machine_sizes_profile.#.rank` is specified")
				}

				if virtualMachineSizesProfile.Rank >= int64(len(state.VirtualMachineSizesProfile)) {
					return fmt.Errorf("`virtual_machine_profile.0.virtual_machine_sizes_profile.%d.rank` must be less than the number of `virtual_machine_sizes_profile` blocks", i)
				}
			}

			virtualMachineSizesProfileNames := make([]any, len(state.VirtualMachineSizesProfile))
			for i, v := range state.VirtualMachineSizesProfile {
				virtualMachineSizesProfileNames[i] = v.Name
			}

			virtualMachineSizesProfileNamesWithoutDuplicate := schema.NewSet(pluginsdk.HashString, virtualMachineSizesProfileNames)
			if len(state.VirtualMachineSizesProfile) > virtualMachineSizesProfileNamesWithoutDuplicate.Len() {
				return errors.New("`virtual_machine_profile.0.virtual_machine_sizes_profile` blocks with duplicated `name` are not allowed")
			}

			if err := validateSecuritySetting(state.VirtualMachineProfile); err != nil {
				return err
			}

			if err := validateWindowsSetting(state.VirtualMachineProfile, metadata.ResourceDiff); err != nil {
				return err
			}

			if err := validateLinuxSetting(state.VirtualMachineProfile, metadata.ResourceDiff); err != nil {
				return err
			}

			for i, galleryApplicationProfile := range state.VirtualMachineProfile[0].GalleryApplicationProfile {
				if galleryApplicationProfile.AutomaticUpgradeEnabled && galleryApplicationProfile.ConfigurationBlobUri != "" {
					return fmt.Errorf("when `virtual_machine_profile.0.gallery_application_profile.%d.automatic_upgrade_enabled` is `true`, `virtual_machine_profile.0.gallery_application_profile.%d.configuration_blob_uri` cannot be specified", i, i)
				}
			}

			for _, v := range state.VirtualMachineProfile[0].Extension {
				if v.ProtectedSettingsJson != "" && len(v.ProtectedSettingsFromKeyVault) > 0 {
					return errors.New("`protected_settings_from_key_vault` cannot be used with `protected_settings_json`")
				}
			}

			if len(state.Zones) > 0 && state.PlatformFaultDomainCount > 1 {
				return errors.New("specifying `zones` is not allowed when `platform_fault_domain_count` higher than 1")
			}

			return nil
		},
	}
}

func validateWindowsSetting(inputList []VirtualMachineProfileModel, d *schema.ResourceDiff) error {
	if len(inputList) == 0 || len(inputList[0].OsProfile) == 0 {
		return nil
	}

	input := &inputList[0]
	if v := input.OsProfile[0].WindowsConfiguration; len(v) > 0 {
		patchMode := v[0].PatchMode
		hotPatchingEnabled := v[0].HotPatchingEnabled
		provisionVirtualMachineAgentEnabled := v[0].ProvisionVirtualMachineAgentEnabled

		rebootSetting := v[0].PatchRebooting
		bypassPlatformSafetyChecksEnabledExist := d.GetRawConfig().AsValueMap()["virtual_machine_profile"].AsValueSlice()[0].AsValueMap()["os_profile"].AsValueSlice()[0].AsValueMap()["windows_configuration"].AsValueSlice()[0].AsValueMap()["bypass_platform_safety_checks_enabled"]
		if !bypassPlatformSafetyChecksEnabledExist.IsNull() || rebootSetting != "" {
			if patchMode != string(fleets.WindowsVMGuestPatchModeAutomaticByPlatform) {
				return errors.New("`bypass_platform_safety_checks_enabled` and `patch_rebooting` cannot be set if the `PatchMode` is not `AutomaticByPlatform`")
			}
		}

		if input.ExtensionOperationsEnabled && !provisionVirtualMachineAgentEnabled {
			return errors.New("`extension_operations_enabled` cannot be set to `true` when `provision_virtual_machine_agent_enabled` is set to `false`")
		}

		isHotPatchEnabledImage := isValidHotPatchSourceImageReference(input.SourceImageReference)
		hasHealthExtension := false
		if extensions := input.Extension; len(extensions) > 0 {
			for _, extension := range extensions {
				if extension.Type == "ApplicationHealthWindows" {
					hasHealthExtension = true
					break
				}
			}
		}

		if isHotPatchEnabledImage {
			// it is a hot patching enabled image, validate hot patching enabled settings
			if patchMode != string(fleets.WindowsVMGuestPatchModeAutomaticByPlatform) {
				return errors.New("when referencing a hot patching enabled image the `patch_mode` field must always be set to `AutomaticByPlatform`")
			}
			if !provisionVirtualMachineAgentEnabled {
				return errors.New("when referencing a hot patching enabled image the `provision_virtual_machine_agent_enabled` field must always be set to `true`")
			}
			if !hasHealthExtension {
				return errors.New("when referencing a hot patching enabled image the `extension` field must always contain a `application health extension`")
			}
			if !hotPatchingEnabled {
				return errors.New("when referencing a hot patching enabled image the `hot_patching_enabled` field must always be set to `true`")
			}
		} else {
			// not a hot patching enabled image verify Automatic VM Guest Patching settings
			if patchMode == string(fleets.WindowsVMGuestPatchModeAutomaticByPlatform) {
				if !provisionVirtualMachineAgentEnabled {
					return fmt.Errorf("when `patch_mode` is set to %q then `provision_virtual_machine_agent_enabled` must be set to `true`", patchMode)
				}
				if !hasHealthExtension {
					return fmt.Errorf("when `patch_mode` is set to %q then the `extension` field must always contain a `application health extension`", patchMode)
				}
			}

			if hotPatchingEnabled {
				return errors.New("`hot_patching_enabled` field is not supported unless you are using one of the following hot patching enable images, `2022-datacenter-azure-edition-core`, `2022-datacenter-azure-edition-core-smalldisk`, `2022-datacenter-azure-edition-hotpatch`, `2022-datacenter-azure-edition-hotpatch-smalldisk`, `2025-datacenter-azure-edition`, `2025-datacenter-azure-edition-smalldisk`, `2025-datacenter-azure-edition-core` or `2025-datacenter-azure-edition-core-smalldisk`")
			}
		}

		if len(v[0].WinRM) == 2 && v[0].WinRM[0].Protocol == v[0].WinRM[1].Protocol {
			return errors.New("`virtual_machine_profile.0.os_profile.0.windows_configuration.0.win_rm` cannot contain elements with same `protocol`")
		}

		for _, winRm := range v[0].WinRM {
			if winRm.Protocol == string(fleets.ProtocolTypesHTTP) && winRm.CertificateUrl != "" {
				return errors.New("`virtual_machine_profile.0.os_profile.0.windows_configuration.0.win_rm.#.certificate_url` cannot be set when `virtual_machine_profile.0.os_profile.0.windows_configuration.0.win_rm.#.protocol` is set to `HTTP`")
			}
		}
	}
	return nil
}

func validateSecuritySetting(inputList []VirtualMachineProfileModel) error {
	if len(inputList) == 0 || len(inputList[0].OsProfile) == 0 {
		return nil
	}

	input := &inputList[0]
	if v := input.OsDisk; len(v) > 0 {
		secureBootEnabled := input.SecureBootEnabled
		vTpmEnabled := input.VTpmEnabled
		if v[0].SecurityEncryptionType != "" {
			if fleets.SecurityEncryptionTypesDiskWithVMGuestState == fleets.SecurityEncryptionTypes(v[0].SecurityEncryptionType) && (!secureBootEnabled || !vTpmEnabled) {
				return errors.New("`secure_boot_enabled` and `vtpm_enabled` must be set to `true` when `os_disk.0.security_encryption_type` is set to `DiskWithVMGuestState`")
			}
			if !vTpmEnabled {
				return errors.New("`vtpm_enabled` must be set to `true` when `os_disk.0.security_encryption_type` is set")
			}
		}
	}
	return nil
}

func validateLinuxSetting(inputList []VirtualMachineProfileModel, d *schema.ResourceDiff) error {
	if len(inputList) == 0 || len(inputList[0].OsProfile) == 0 {
		return nil
	}

	input := &inputList[0]
	if v := input.OsProfile[0].LinuxConfiguration; len(v) > 0 {
		patchMode := v[0].PatchMode
		provisionVirtualMachineAgentEnabled := v[0].ProvisionVirtualMachineAgentEnabled

		rebootSetting := v[0].PatchRebooting
		bypassPlatformSafetyChecksEnabledExist := d.GetRawConfig().AsValueMap()["virtual_machine_profile"].AsValueSlice()[0].AsValueMap()["os_profile"].AsValueSlice()[0].AsValueMap()["linux_configuration"].AsValueSlice()[0].AsValueMap()["bypass_platform_safety_checks_enabled"]
		if !bypassPlatformSafetyChecksEnabledExist.IsNull() || rebootSetting != "" {
			if patchMode != string(fleets.LinuxVMGuestPatchModeAutomaticByPlatform) {
				return errors.New("`bypass_platform_safety_checks_enabled` and `patch_rebooting` cannot be set if the `patch_mode` is not `AutomaticByPlatform`")
			}
		}

		if input.ExtensionOperationsEnabled && !provisionVirtualMachineAgentEnabled {
			return errors.New("`extension_operations_enabled` cannot be set to `true` when `provision_virtual_machine_agent_enabled` is set to `false`")
		}

		hasHealthExtension := false
		if extensions := input.Extension; len(extensions) > 0 {
			for _, extension := range extensions {
				if extension.Type == "ApplicationHealthLinux" {
					hasHealthExtension = true
					break
				}
			}
		}

		if patchMode == string(fleets.LinuxVMGuestPatchModeAutomaticByPlatform) {
			if !provisionVirtualMachineAgentEnabled {
				return fmt.Errorf("when the `patch_mode` field is set to %q the `provision_virtual_machine_agent_enabled` field must always be set to `true`, got %q", patchMode, strconv.FormatBool(provisionVirtualMachineAgentEnabled))
			}

			if !hasHealthExtension {
				return fmt.Errorf("when the `patch_mode` field is set to %q the `extension` field must contain at least one `application health extension`, got 0", patchMode)
			}
		}

		if v[0].AdminPassword == "" && v[0].PasswordAuthenticationEnabled {
			return errors.New("`admin_password` is required when `password_authentication_enabled` is enabled")
		}
	}
	return nil
}

func isValidHotPatchSourceImageReference(referenceInput []SourceImageReferenceModel) bool {
	if len(referenceInput) == 0 {
		return false
	}
	raw := referenceInput[0]
	pub := raw.Publisher
	offer := raw.Offer
	sku := raw.Sku

	supportedSkus := []string{
		"2022-datacenter-azure-edition-core",
		"2022-datacenter-azure-edition-core-smalldisk",
		"2022-datacenter-azure-edition-hotpatch",
		"2022-datacenter-azure-edition-hotpatch-smalldisk",
		"2025-datacenter-azure-edition",
		"2025-datacenter-azure-edition-smalldisk",
		"2025-datacenter-azure-edition-core",
		"2025-datacenter-azure-edition-core-smalldisk",
	}

	if pub == "MicrosoftWindowsServer" && offer == "WindowsServer" && slices.Contains(supportedSkus, sku) {
		return true
	}

	return false
}
