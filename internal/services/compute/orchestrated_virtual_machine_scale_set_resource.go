// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/capacityreservationgroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/images"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/proximityplacementgroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-07-01/virtualmachinescalesets"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	computeValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

const (
	SkuNameMix = "Mix"
)

func resourceOrchestratedVirtualMachineScaleSet() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceOrchestratedVirtualMachineScaleSetCreate,
		Read:   resourceOrchestratedVirtualMachineScaleSetRead,
		Update: resourceOrchestratedVirtualMachineScaleSetUpdate,
		Delete: resourceOrchestratedVirtualMachineScaleSetDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := commonids.ParseVirtualMachineScaleSetID(id)
			return err
		}, importOrchestratedVirtualMachineScaleSet),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		// The plan was to remove support the legacy Orchestrated Virtual Machine Scale Set in 3.0.
		// Turns out it's still in use
		// TODO: Revisit in 4.0
		// TODO: exposing requireGuestProvisionSignal once it's available
		// https://github.com/Azure/azure-rest-api-specs/pull/7246

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: computeValidate.VirtualMachineName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"network_interface": OrchestratedVirtualMachineScaleSetNetworkInterfaceSchema(),

			"os_disk": OrchestratedVirtualMachineScaleSetOSDiskSchema(),

			"instances": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 1000),
			},

			// For sku I will create a format like: tier_sku name.
			// NOTE: all of the exposed vm sku tier's are Standard so this will continue to be hardcoded
			// Examples: Standard_HC44rs_4, Standard_D48_v3_6, Standard_M64s_20, Standard_HB120-96rs_v3_8
			"sku_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: computeValidate.OrchestratedVirtualMachineScaleSetSku,
			},

			"sku_profile": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"allocation_strategy": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice(
								virtualmachinescalesets.PossibleValuesForAllocationStrategy(),
								false,
							),
						},

						"vm_sizes": {
							Type:     pluginsdk.TypeSet,
							Required: true,
							MinItems: 1,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},
			},

			"os_profile": OrchestratedVirtualMachineScaleSetOSProfileSchema(),

			// Optional
			// NOTE: The schema for the automatic instance repair has merged so they are
			// identical for both uniform and flex mode VMSS's
			"automatic_instance_repair": VirtualMachineScaleSetAutomaticRepairsPolicySchema(),

			"boot_diagnostics": bootDiagnosticsSchema(),

			"capacity_reservation_group_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: capacityreservationgroups.ValidateCapacityReservationGroupID,
				ConflictsWith: []string{
					"proximity_placement_group_id",
				},
			},

			"data_disk": OrchestratedVirtualMachineScaleSetDataDiskSchema(),

			// Optional
			"additional_capabilities": OrchestratedVirtualMachineScaleSetAdditionalCapabilitiesSchema(),

			"encryption_at_host_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"eviction_policy": {
				// only applicable when `priority` is set to `Spot`
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(virtualmachinescalesets.VirtualMachineEvictionPolicyTypesDeallocate),
					string(virtualmachinescalesets.VirtualMachineEvictionPolicyTypesDelete),
				}, false),
			},

			"extension_operations_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
			},

			// Due to bug in RP extensions cannot currently be supported in Terraform ETA for full support is mid Jan 2022
			"extension": OrchestratedVirtualMachineScaleSetExtensionsSchema(),

			"extensions_time_budget": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Default:      "PT1H30M",
				ValidateFunc: validate.ISO8601DurationBetween("PT15M", "PT2H"),
			},

			// whilst the Swagger defines multiple at this time only UAI is supported
			"identity": commonschema.UserAssignedIdentityOptional(),

			"license_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"None",
					"Windows_Client",
					"Windows_Server",
				}, false),
				DiffSuppressFunc: func(_, old, new string, _ *pluginsdk.ResourceData) bool {
					if old == "None" && new == "" || old == "" && new == "None" {
						return true
					}

					return false
				},
			},

			"max_bid_price": {
				Type:         pluginsdk.TypeFloat,
				Optional:     true,
				Default:      -1,
				ValidateFunc: computeValidate.SpotMaxPrice,
			},

			"plan": planSchema(),

			"platform_fault_domain_count": {
				Type:     pluginsdk.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"priority": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(virtualmachinescalesets.VirtualMachinePriorityTypesRegular),
				ValidateFunc: validation.StringInSlice([]string{
					string(virtualmachinescalesets.VirtualMachinePriorityTypesRegular),
					string(virtualmachinescalesets.VirtualMachinePriorityTypesSpot),
				}, false),
			},

			"proximity_placement_group_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: proximityplacementgroups.ValidateProximityPlacementGroupID,
				// the Compute API is broken and returns the Resource Group name in UPPERCASE :shrug:, github issue: https://github.com/Azure/azure-rest-api-specs/issues/10016
				DiffSuppressFunc: suppress.CaseDifference,
				ConflictsWith: []string{
					"capacity_reservation_group_id",
				},
			},

			"rolling_upgrade_policy": VirtualMachineScaleSetRollingUpgradePolicySchema(),

			// NOTE: single_placement_group is now supported in orchestrated VMSS
			// Since null is now a valid value for this field there is no default
			// for this bool
			"single_placement_group": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
				Optional: true,
			},

			"source_image_id": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.Any(
					images.ValidateImageID,
					computeValidate.SharedImageID,
					computeValidate.SharedImageVersionID,
					computeValidate.CommunityGalleryImageID,
					computeValidate.CommunityGalleryImageVersionID,
					computeValidate.SharedGalleryImageID,
					computeValidate.SharedGalleryImageVersionID,
				),
				ConflictsWith: []string{
					"source_image_reference",
				},
			},

			"source_image_reference": sourceImageReferenceSchemaOrchestratedVMSS(),

			"zone_balance": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},

			"termination_notification": OrchestratedVirtualMachineScaleSetTerminationNotificationSchema(),

			"zones": commonschema.ZonesMultipleOptional(),

			"tags": commonschema.Tags(),

			// Computed
			"unique_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"upgrade_mode": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(virtualmachinescalesets.UpgradeModeManual),
				ValidateFunc: validation.StringInSlice([]string{
					string(virtualmachinescalesets.UpgradeModeAutomatic),
					string(virtualmachinescalesets.UpgradeModeManual),
					string(virtualmachinescalesets.UpgradeModeRolling),
				}, false),
			},

			"user_data_base64": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsBase64,
			},

			"priority_mix": OrchestratedVirtualMachineScaleSetPriorityMixPolicySchema(),
		},

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			// Removing existing zones is currently not supported for Virtual Machine Scale Sets
			pluginsdk.ForceNewIfChange("zones", func(ctx context.Context, old, new, meta interface{}) bool {
				oldZones := zones.ExpandUntyped(old.(*schema.Set).List())
				newZones := zones.ExpandUntyped(new.(*schema.Set).List())

				for _, ov := range oldZones {
					found := false
					for _, nv := range newZones {
						if ov == nv {
							found = true
							break
						}
					}

					if !found {
						return true
					}
				}

				return false
			}),

			pluginsdk.CustomizeDiffShim(func(ctx context.Context, diff *pluginsdk.ResourceDiff, v interface{}) error {
				skuName, hasSkuName := diff.GetOk("sku_name")
				_, hasSkuProfile := diff.GetOk("sku_profile")

				if hasSkuProfile {
					if !hasSkuName || skuName != SkuNameMix {
						return fmt.Errorf("`sku_profile` can only be set when `sku_name` is set to `Mix`")
					}
				} else {
					if hasSkuName && skuName == SkuNameMix {
						return fmt.Errorf("`sku_profile` must be set when `sku_name` is set to `Mix`")
					}
				}

				upgradeMode := virtualmachinescalesets.UpgradeMode(diff.Get("upgrade_mode").(string))
				rollingUpgradePolicyRaw := diff.Get("rolling_upgrade_policy").([]interface{})

				if upgradeMode == virtualmachinescalesets.UpgradeModeManual && len(rollingUpgradePolicyRaw) > 0 {
					return fmt.Errorf("a `rolling_upgrade_policy` block cannot be specified when `upgrade_mode` is set to `%s`", string(upgradeMode))
				}

				if upgradeMode == virtualmachinescalesets.UpgradeModeRolling && len(rollingUpgradePolicyRaw) == 0 {
					return fmt.Errorf("a `rolling_upgrade_policy` block must be specified when `upgrade_mode` is set to `%s`", string(upgradeMode))
				}

				return nil
			}),
		),
	}
}

func resourceOrchestratedVirtualMachineScaleSetCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VirtualMachineScaleSetsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	isLegacy := true
	id := virtualmachinescalesets.NewVirtualMachineScaleSetID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		// Upgrading to the 2021-07-01 exposed a new expand parameter to the GET method
		existing, err := client.Get(ctx, id, virtualmachinescalesets.DefaultGetOperationOptions())
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_orchestrated_virtual_machine_scale_set", id.ID())
		}
	}

	t := d.Get("tags").(map[string]interface{})

	props := virtualmachinescalesets.VirtualMachineScaleSet{
		Location: location.Normalize(d.Get("location").(string)),
		Tags:     tags.Expand(t),
		Properties: &virtualmachinescalesets.VirtualMachineScaleSetProperties{
			PlatformFaultDomainCount: pointer.To(int64(d.Get("platform_fault_domain_count").(int))),
			// OrchestrationMode needs to be hardcoded to Uniform, for the
			// standard VMSS resource, since virtualMachineProfile is now supported
			// in both VMSS and Orchestrated VMSS...
			OrchestrationMode: pointer.To(virtualmachinescalesets.OrchestrationModeFlexible),
		},
	}

	// The RP now accepts true, false and null for single_placement_group value.
	// This is only valid for the Orchestrated VMSS Resource. If the
	// single_placement_group is null(e.g. not passed in the props) the RP will
	// automatically determine what values single_placement_group should be
	if !pluginsdk.IsExplicitlyNullInConfig(d, "single_placement_group") {
		props.Properties.SinglePlacementGroup = pointer.To(d.Get("single_placement_group").(bool))
	}

	zones := zones.ExpandUntyped(d.Get("zones").(*schema.Set).List())
	if len(zones) > 0 {
		props.Zones = &zones
	}

	upgradeMode := virtualmachinescalesets.UpgradeMode(d.Get("upgrade_mode").(string))
	rollingUpgradePolicy, err := ExpandVirtualMachineScaleSetRollingUpgradePolicy(d.Get("rolling_upgrade_policy").([]interface{}), len(zones) > 0, false)
	if err != nil {
		return fmt.Errorf("expanding `rolling_upgrade_policy`: %+v", err)
	}

	props.Properties.UpgradePolicy = &virtualmachinescalesets.UpgradePolicy{
		Mode:                 pointer.To(upgradeMode),
		RollingUpgradePolicy: rollingUpgradePolicy,
	}

	virtualMachineProfile := virtualmachinescalesets.VirtualMachineScaleSetVMProfile{
		StorageProfile: &virtualmachinescalesets.VirtualMachineScaleSetStorageProfile{},
	}

	networkProfile := &virtualmachinescalesets.VirtualMachineScaleSetNetworkProfile{
		// 2020-11-01 is the only valid value for this value and is only valid for VMSS in Orchestration Mode flex
		NetworkApiVersion: pointer.To(virtualmachinescalesets.NetworkApiVersionTwoZeroTwoZeroNegativeOneOneNegativeZeroOne),
	}

	if v, ok := d.GetOk("proximity_placement_group_id"); ok {
		props.Properties.ProximityPlacementGroup = &virtualmachinescalesets.SubResource{
			Id: pointer.To(v.(string)),
		}
	}

	// Not currently supported in OVMSS
	// healthProbeId := d.Get("health_probe_id").(string)
	// upgradeMode := virtualmachinescalesets.UpgradeMode(d.Get("upgrade_mode").(string))

	instances := d.Get("instances").(int)
	if v, ok := d.GetOk("sku_name"); ok {
		isLegacy = false
		sku, err := expandOrchestratedVirtualMachineScaleSetSku(v.(string), instances)
		if err != nil {
			return fmt.Errorf("expanding 'sku_name': %+v", err)
		}
		props.Sku = sku
	}

	if v, ok := d.GetOk("sku_profile"); ok {
		props.Properties.SkuProfile = expandOrchestratedVirtualMachineScaleSetSkuProfile(v.([]interface{}))
	}

	if v, ok := d.GetOk("capacity_reservation_group_id"); ok {
		if d.Get("single_placement_group").(bool) {
			return fmt.Errorf("`single_placement_group` must be set to `false` when `capacity_reservation_group_id` is specified")
		}

		virtualMachineProfile.CapacityReservation = &virtualmachinescalesets.CapacityReservationProfile{
			CapacityReservationGroup: &virtualmachinescalesets.SubResource{
				Id: pointer.To(v.(string)),
			},
		}
	}

	// hasHealthExtension is currently not needed but I added the plumming because we will need it
	// once upgrade policy is added to OVMSS
	hasHealthExtension := false

	if v, ok := d.GetOk("extension"); ok {
		var err error
		virtualMachineProfile.ExtensionProfile, hasHealthExtension, err = expandOrchestratedVirtualMachineScaleSetExtensions(v.(*pluginsdk.Set).List())
		if err != nil {
			return err
		}
	}

	if hasHealthExtension {
		log.Printf("[DEBUG] Orchestrated %s has a Health Extension defined", id)
	}

	// Virtual Machine Scale Set with Flexible Orchestration Mode and 'Rolling' upgradeMode must have Health Extension Present
	if upgradeMode == virtualmachinescalesets.UpgradeModeRolling && !hasHealthExtension {
		return fmt.Errorf("a health extension must be specified when `upgrade_mode` is set to `Rolling`")
	}

	if v, ok := d.GetOk("extensions_time_budget"); ok {
		if virtualMachineProfile.ExtensionProfile == nil {
			virtualMachineProfile.ExtensionProfile = &virtualmachinescalesets.VirtualMachineScaleSetExtensionProfile{}
		}
		virtualMachineProfile.ExtensionProfile.ExtensionsTimeBudget = pointer.To(v.(string))
	}

	sourceImageReferenceRaw := d.Get("source_image_reference").([]interface{})
	sourceImageId := d.Get("source_image_id").(string)
	if len(sourceImageReferenceRaw) != 0 || sourceImageId != "" {
		sourceImageReference := expandSourceImageReferenceVMSS(sourceImageReferenceRaw, sourceImageId)
		virtualMachineProfile.StorageProfile.ImageReference = sourceImageReference
	}

	if userData, ok := d.GetOk("user_data_base64"); ok {
		virtualMachineProfile.UserData = pointer.To(userData.(string))
	}

	osType := virtualmachinescalesets.OperatingSystemTypesWindows
	var winConfigRaw []interface{}
	var linConfigRaw []interface{}
	var vmssOsProfile *virtualmachinescalesets.VirtualMachineScaleSetOSProfile
	extensionOperationsEnabled := d.Get("extension_operations_enabled").(bool)
	osProfileRaw := d.Get("os_profile").([]interface{})

	if len(osProfileRaw) > 0 {
		osProfile := osProfileRaw[0].(map[string]interface{})
		winConfigRaw = osProfile["windows_configuration"].([]interface{})
		linConfigRaw = osProfile["linux_configuration"].([]interface{})
		customData := ""

		// Pass custom data if it is defined in the config file
		if v := osProfile["custom_data"]; v != nil {
			customData = v.(string)
		}

		if len(winConfigRaw) > 0 {
			winConfig := winConfigRaw[0].(map[string]interface{})
			provisionVMAgent := winConfig["provision_vm_agent"].(bool)
			patchAssessmentMode := winConfig["patch_assessment_mode"].(string)
			vmssOsProfile = expandOrchestratedVirtualMachineScaleSetOsProfileWithWindowsConfiguration(winConfig, customData)

			// if the Computer Prefix Name was not defined use the computer name
			if vmssOsProfile.ComputerNamePrefix == nil || len(*vmssOsProfile.ComputerNamePrefix) == 0 {
				// validate that the computer name is a valid Computer Prefix Name
				_, errs := computeValidate.WindowsComputerNamePrefix(id.VirtualMachineScaleSetName, "computer_name_prefix")
				if len(errs) > 0 {
					return fmt.Errorf("unable to assume default computer name prefix %s. Please adjust the 'name', or specify an explicit 'computer_name_prefix'", errs[0])
				}
				vmssOsProfile.ComputerNamePrefix = pointer.To(id.VirtualMachineScaleSetName)
			}

			if extensionOperationsEnabled && !provisionVMAgent {
				return fmt.Errorf("`extension_operations_enabled` cannot be set to `true` when `provision_vm_agent` is set to `false`")
			}

			if patchAssessmentMode == string(virtualmachinescalesets.WindowsPatchAssessmentModeAutomaticByPlatform) && !provisionVMAgent {
				return fmt.Errorf("when the 'patch_assessment_mode' field is set to %q the 'provision_vm_agent' must always be set to 'true'", virtualmachinescalesets.WindowsPatchAssessmentModeAutomaticByPlatform)
			}

			// Validate patch mode and hotpatching configuration
			isHotpatchEnabledImage := isValidHotPatchSourceImageReference(sourceImageReferenceRaw, sourceImageId)
			patchMode := winConfig["patch_mode"].(string)
			hotpatchingEnabled := winConfig["hotpatching_enabled"].(bool)

			if isHotpatchEnabledImage {
				// it is a hotpatching enabled image, validate hotpatching enabled settings
				if patchMode != string(virtualmachinescalesets.WindowsVMGuestPatchModeAutomaticByPlatform) {
					return fmt.Errorf("when referencing a hotpatching enabled image the 'patch_mode' field must always be set to %q", virtualmachinescalesets.WindowsVMGuestPatchModeAutomaticByPlatform)
				}

				if !provisionVMAgent {
					return fmt.Errorf("when referencing a hotpatching enabled image the 'provision_vm_agent' field must always be set to 'true'")
				}

				if !hasHealthExtension {
					return fmt.Errorf("when referencing a hotpatching enabled image the 'extension' field must always contain a 'application health extension'")
				}

				if !hotpatchingEnabled {
					return fmt.Errorf("when referencing a hotpatching enabled image the 'hotpatching_enabled' field must always be set to 'true'")
				}
			} else {
				// not a hotpatching enabled image verify Automatic VM Guest Patching settings
				if patchMode == string(virtualmachinescalesets.WindowsVMGuestPatchModeAutomaticByPlatform) {
					if !provisionVMAgent {
						return fmt.Errorf("when 'patch_mode' is set to %q then 'provision_vm_agent' must be set to 'true'", patchMode)
					}

					if !hasHealthExtension {
						return fmt.Errorf("when 'patch_mode' is set to %q then the 'extension' field must always contain a 'application health extension'", patchMode)
					}
				}

				if hotpatchingEnabled {
					return fmt.Errorf("'hotpatching_enabled' field is not supported unless you are using one of the following hotpatching enable images, '2022-datacenter-azure-edition', '2022-datacenter-azure-edition-core-smalldisk', '2022-datacenter-azure-edition-hotpatch', '2022-datacenter-azure-edition-hotpatch-smalldisk', '2025-datacenter-azure-edition', '2025-datacenter-azure-edition-smalldisk', '2025-datacenter-azure-edition-core' or '2025-datacenter-azure-edition-core-smalldisk'")
				}
			}
		}

		if len(linConfigRaw) > 0 {
			osType = virtualmachinescalesets.OperatingSystemTypesLinux
			linConfig := linConfigRaw[0].(map[string]interface{})
			provisionVMAgent := linConfig["provision_vm_agent"].(bool)
			patchAssessmentMode := linConfig["patch_assessment_mode"].(string)
			vmssOsProfile = expandOrchestratedVirtualMachineScaleSetOsProfileWithLinuxConfiguration(linConfig, customData)

			// if the Computer Prefix Name was not defined use the computer name
			if vmssOsProfile.ComputerNamePrefix == nil || len(*vmssOsProfile.ComputerNamePrefix) == 0 {
				// validate that the computer name is a valid Computer Prefix Name
				_, errs := computeValidate.LinuxComputerNamePrefix(id.VirtualMachineScaleSetName, "computer_name_prefix")
				if len(errs) > 0 {
					return fmt.Errorf("unable to assume default computer name prefix %s. Please adjust the 'name', or specify an explicit 'computer_name_prefix'", errs[0])
				}

				vmssOsProfile.ComputerNamePrefix = pointer.To(id.VirtualMachineScaleSetName)
			}

			if extensionOperationsEnabled && !provisionVMAgent {
				return fmt.Errorf("`extension_operations_enabled` cannot be set to `true` when `provision_vm_agent` is set to `false`")
			}

			if patchAssessmentMode == string(virtualmachinescalesets.LinuxPatchAssessmentModeAutomaticByPlatform) && !provisionVMAgent {
				return fmt.Errorf("when the 'patch_assessment_mode' field is set to %q the 'provision_vm_agent' must always be set to 'true'", virtualmachinescalesets.LinuxPatchAssessmentModeAutomaticByPlatform)
			}

			// Validate Automatic VM Guest Patching Settings
			patchMode := linConfig["patch_mode"].(string)

			if patchMode == string(virtualmachinescalesets.LinuxVMGuestPatchModeAutomaticByPlatform) {
				if !provisionVMAgent {
					return fmt.Errorf("when the 'patch_mode' field is set to %q the 'provision_vm_agent' field must always be set to 'true', got %q", patchMode, strconv.FormatBool(provisionVMAgent))
				}

				if !hasHealthExtension {
					return fmt.Errorf("when the 'patch_mode' field is set to %q the 'extension' field must contain at least one 'application health extension', got 0", patchMode)
				}
			}
		}

		virtualMachineProfile.OsProfile = vmssOsProfile
	}

	if virtualMachineProfile.OsProfile == nil {
		virtualMachineProfile.OsProfile = &virtualmachinescalesets.VirtualMachineScaleSetOSProfile{}
	}
	virtualMachineProfile.OsProfile.AllowExtensionOperations = pointer.To(extensionOperationsEnabled)

	if v, ok := d.GetOk("boot_diagnostics"); ok {
		virtualMachineProfile.DiagnosticsProfile = expandBootDiagnosticsVMSS(v.([]interface{}))
	}

	if v, ok := d.GetOk("priority"); ok {
		virtualMachineProfile.Priority = pointer.To(virtualmachinescalesets.VirtualMachinePriorityTypes(v.(string)))
	}

	if v, ok := d.GetOk("os_disk"); ok {
		virtualMachineProfile.StorageProfile.OsDisk = ExpandOrchestratedVirtualMachineScaleSetOSDisk(v.([]interface{}), osType)
	}

	additionalCapabilitiesRaw := d.Get("additional_capabilities").([]interface{})
	additionalCapabilities := ExpandOrchestratedVirtualMachineScaleSetAdditionalCapabilities(additionalCapabilitiesRaw)
	props.Properties.AdditionalCapabilities = additionalCapabilities

	if v, ok := d.GetOk("data_disk"); ok {
		ultraSSDEnabled := d.Get("additional_capabilities.0.ultra_ssd_enabled").(bool)
		dataDisks, err := ExpandOrchestratedVirtualMachineScaleSetDataDisk(v.([]interface{}), ultraSSDEnabled)
		if err != nil {
			return fmt.Errorf("expanding `data_disk`: %+v", err)
		}
		virtualMachineProfile.StorageProfile.DataDisks = dataDisks
	}

	if v, ok := d.GetOk("network_interface"); ok {
		networkInterfaces, err := ExpandOrchestratedVirtualMachineScaleSetNetworkInterface(v.([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `network_interface`: %+v", err)
		}

		networkProfile.NetworkInterfaceConfigurations = networkInterfaces
		virtualMachineProfile.NetworkProfile = networkProfile
	}

	if v, ok := d.Get("max_bid_price").(float64); ok && v > 0 {
		if *virtualMachineProfile.Priority != virtualmachinescalesets.VirtualMachinePriorityTypesSpot {
			return fmt.Errorf("`max_bid_price` can only be configured when `priority` is set to `Spot`")
		}

		virtualMachineProfile.BillingProfile = &virtualmachinescalesets.BillingProfile{
			MaxPrice: pointer.To(v),
		}
	}

	if v, ok := d.GetOk("encryption_at_host_enabled"); ok {
		virtualMachineProfile.SecurityProfile = &virtualmachinescalesets.SecurityProfile{
			EncryptionAtHost: pointer.To(v.(bool)),
		}
	}

	if v, ok := d.GetOk("eviction_policy"); ok {
		if *virtualMachineProfile.Priority != virtualmachinescalesets.VirtualMachinePriorityTypesSpot {
			return fmt.Errorf("an `eviction_policy` can only be specified when `priority` is set to `Spot`")
		}
		virtualMachineProfile.EvictionPolicy = pointer.To(virtualmachinescalesets.VirtualMachineEvictionPolicyTypes(v.(string)))
	} else if *virtualMachineProfile.Priority == virtualmachinescalesets.VirtualMachinePriorityTypesSpot {
		return fmt.Errorf("an `eviction_policy` must be specified when `priority` is set to `Spot`")
	}

	if v, ok := d.GetOk("license_type"); ok {
		virtualMachineProfile.LicenseType = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("termination_notification"); ok {
		virtualMachineProfile.ScheduledEventsProfile = ExpandOrchestratedVirtualMachineScaleSetScheduledEventsProfile(v.([]interface{}))
	}

	// Only inclued the virtual machine profile if this is not a legacy configuration
	if !isLegacy {
		if v, ok := d.GetOk("plan"); ok {
			props.Plan = expandPlanVMSS(v.([]interface{}))
		}

		if v, ok := d.GetOk("identity"); ok {
			identityExpanded, err := identity.ExpandSystemAndUserAssignedMap(v.([]interface{}))
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}
			props.Identity = identityExpanded
		}

		if v, ok := d.GetOk("automatic_instance_repair"); ok {
			if !hasHealthExtension {
				return fmt.Errorf("`automatic_instance_repair` can only be set if there is an application Health extension defined")
			}

			props.Properties.AutomaticRepairsPolicy = ExpandVirtualMachineScaleSetAutomaticRepairsPolicy(v.([]interface{}))
		}

		if v, ok := d.GetOk("zone_balance"); ok && v.(bool) {
			if props.Zones == nil || len(*props.Zones) == 0 {
				return fmt.Errorf("`zone_balance` can only be set to `true` when zones are specified")
			}

			props.Properties.ZoneBalance = pointer.To(v.(bool))
		}

		if v, ok := d.GetOk("priority_mix"); ok {
			if *virtualMachineProfile.Priority != virtualmachinescalesets.VirtualMachinePriorityTypesSpot {
				return fmt.Errorf("a `priority_mix` can only be specified when `priority` is set to `Spot`")
			}
			props.Properties.PriorityMixPolicy = ExpandOrchestratedVirtualMachineScaleSetPriorityMixPolicy(v.([]interface{}))
		}

		props.Properties.VirtualMachineProfile = &virtualMachineProfile
	}

	log.Printf("[DEBUG] Creating Orchestrated %s.", id)
	if err := client.CreateOrUpdateThenPoll(ctx, id, props, virtualmachinescalesets.DefaultCreateOrUpdateOperationOptions()); err != nil {
		return fmt.Errorf("creating Orchestrated %s: %+v", id, err)
	}

	log.Printf("[DEBUG] Orchestrated %s was created", id)
	log.Printf("[DEBUG] Retrieving Orchestrated %s.", id)

	d.SetId(id.ID())

	return resourceOrchestratedVirtualMachineScaleSetRead(d, meta)
}

func resourceOrchestratedVirtualMachineScaleSetUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VirtualMachineScaleSetsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualmachinescalesets.ParseVirtualMachineScaleSetID(d.Id())
	if err != nil {
		return err
	}

	isLegacy := true
	updateInstances := false
	isHotpatchEnabledImage := false
	linuxAutomaticVMGuestPatchingEnabled := false

	options := virtualmachinescalesets.DefaultGetOperationOptions()
	options.Expand = pointer.To(virtualmachinescalesets.ExpandTypesForGetVMScaleSetsUserData)
	existing, err := client.Get(ctx, *id, options)
	if err != nil {
		return fmt.Errorf("retrieving Orchestrated %s: %+v", id, err)
	}
	if existing.Model == nil {
		return fmt.Errorf("retrieving Orchestrated %s: `model` was nil", id)
	}
	if existing.Model.Sku != nil {
		isLegacy = false
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving Orchestrated %s: `properties` was nil", id)
	}

	if !isLegacy {
		if existing.Model.Properties.VirtualMachineProfile == nil {
			return fmt.Errorf("retrieving Orchestrated %s: `properties.virtualMachineProfile` was nil", id)
		}
		if existing.Model.Properties.VirtualMachineProfile.StorageProfile == nil {
			return fmt.Errorf("retrieving Orchestrated %s: `properties.virtualMachineProfile,storageProfile` was nil", id)
		}
	}

	updateProps := virtualmachinescalesets.VirtualMachineScaleSetUpdateProperties{}
	update := virtualmachinescalesets.VirtualMachineScaleSetUpdate{}
	osType := virtualmachinescalesets.OperatingSystemTypesWindows

	if !isLegacy {
		updateProps = virtualmachinescalesets.VirtualMachineScaleSetUpdateProperties{
			VirtualMachineProfile: &virtualmachinescalesets.VirtualMachineScaleSetUpdateVMProfile{
				// if an image reference has been configured previously (it has to be), we would better to include that in this
				// update request to avoid some circumstances that the API will complain ImageReference is null
				// issue tracking: https://github.com/Azure/azure-rest-api-specs/issues/10322
				StorageProfile: &virtualmachinescalesets.VirtualMachineScaleSetUpdateStorageProfile{
					ImageReference: existing.Model.Properties.VirtualMachineProfile.StorageProfile.ImageReference,
				},
			},
			// Currently not suppored in orchestrated VMSS
			// if an upgrade policy's been configured previously (which it will have) it must be threaded through
			// this doesn't matter for Manual - but breaks when updating anything on a Automatic and Rolling Mode Scale Set
			// UpgradePolicy: existing.Properties.UpgradePolicy,
		}

		priority := virtualmachinescalesets.VirtualMachinePriorityTypes(d.Get("priority").(string))

		if d.HasChange("single_placement_group") {
			// Since null is now a valid value for single_placement_group
			// make sure it is in the config file before you set the value
			// on the update props...
			if !pluginsdk.IsExplicitlyNullInConfig(d, "single_placement_group") {
				singlePlacementGroup := d.Get("single_placement_group").(bool)
				if singlePlacementGroup {
					return fmt.Errorf("'single_placement_group' can not be set to 'true' once it has been set to 'false'")
				}
				updateProps.SinglePlacementGroup = pointer.To(singlePlacementGroup)
			}
		}

		if d.HasChange("sku_profile") {
			updateInstances = true
			updateProps.SkuProfile = expandOrchestratedVirtualMachineScaleSetSkuProfile(d.Get("sku_profile").([]interface{}))
		}

		if d.HasChange("max_bid_price") {
			if priority != virtualmachinescalesets.VirtualMachinePriorityTypesSpot {
				return fmt.Errorf("`max_bid_price` can only be configured when `priority` is set to `Spot`")
			}

			updateProps.VirtualMachineProfile.BillingProfile = &virtualmachinescalesets.BillingProfile{
				MaxPrice: pointer.To(d.Get("max_bid_price").(float64)),
			}
		}

		osProfileRaw := d.Get("os_profile").([]interface{})
		vmssOsProfile := virtualmachinescalesets.VirtualMachineScaleSetUpdateOSProfile{}
		windowsConfig := virtualmachinescalesets.WindowsConfiguration{}
		windowsConfig.PatchSettings = &virtualmachinescalesets.PatchSettings{}
		linuxConfig := virtualmachinescalesets.LinuxConfiguration{}

		if len(osProfileRaw) > 0 {
			osProfile := osProfileRaw[0].(map[string]interface{})
			winConfigRaw := osProfile["windows_configuration"].([]interface{})
			linConfigRaw := osProfile["linux_configuration"].([]interface{})

			if d.HasChange("os_profile.0.custom_data") {
				updateInstances = true

				// customData can only be sent if it's a base64 encoded string,
				// so it's not possible to remove this without tainting the resource
				vmssOsProfile.CustomData = pointer.To(osProfile["custom_data"].(string))
			}

			if len(winConfigRaw) > 0 {
				winConfig := winConfigRaw[0].(map[string]interface{})
				provisionVMAgent := winConfig["provision_vm_agent"].(bool)
				patchAssessmentMode := winConfig["patch_assessment_mode"].(string)
				patchMode := winConfig["patch_mode"].(string)
				hotpatchingEnabled := winConfig["hotpatching_enabled"].(bool)

				// If the image allows hotpatching the patch mode can only ever be AutomaticByPlatform.
				sourceImageReferenceRaw := d.Get("source_image_reference").([]interface{})
				sourceImageId := d.Get("source_image_id").(string)
				isHotpatchEnabledImage = isValidHotPatchSourceImageReference(sourceImageReferenceRaw, sourceImageId)

				// PatchSettings is required by PATCH API when running Hotpatch-compatible images.
				if isHotpatchEnabledImage {
					windowsConfig.PatchSettings.AssessmentMode = pointer.To(virtualmachinescalesets.WindowsPatchAssessmentMode(patchAssessmentMode))
					windowsConfig.PatchSettings.PatchMode = pointer.To(virtualmachinescalesets.WindowsVMGuestPatchMode(patchMode))
					windowsConfig.PatchSettings.EnableHotpatching = pointer.To(hotpatchingEnabled)
				}

				if d.HasChange("os_profile.0.windows_configuration.0.enable_automatic_updates") ||
					d.HasChange("os_profile.0.windows_configuration.0.provision_vm_agent") ||
					d.HasChange("os_profile.0.windows_configuration.0.timezone") ||
					d.HasChange("os_profile.0.windows_configuration.0.secret") ||
					d.HasChange("os_profile.0.windows_configuration.0.winrm_listener") {
					updateInstances = true
				}

				if d.HasChange("os_profile.0.windows_configuration.0.enable_automatic_updates") {
					windowsConfig.EnableAutomaticUpdates = pointer.To(winConfig["enable_automatic_updates"].(bool))
				}

				if d.HasChange("os_profile.0.windows_configuration.0.provision_vm_agent") {
					if isHotpatchEnabledImage && !provisionVMAgent {
						return fmt.Errorf("when referencing a hotpatching enabled image the 'provision_vm_agent' field must always be set to 'true', got %q", strconv.FormatBool(provisionVMAgent))
					}
					windowsConfig.ProvisionVMAgent = pointer.To(provisionVMAgent)
				}

				if d.HasChange("os_profile.0.windows_configuration.0.patch_assessment_mode") {
					if !provisionVMAgent && (patchAssessmentMode == string(virtualmachinescalesets.WindowsPatchAssessmentModeAutomaticByPlatform)) {
						return fmt.Errorf("when the 'patch_assessment_mode' field is set to %q the 'provision_vm_agent' must always be set to 'true'", virtualmachinescalesets.WindowsPatchAssessmentModeAutomaticByPlatform)
					}
					windowsConfig.PatchSettings.AssessmentMode = pointer.To(virtualmachinescalesets.WindowsPatchAssessmentMode(patchAssessmentMode))
				}

				if d.HasChange("os_profile.0.windows_configuration.0.patch_mode") {
					if isHotpatchEnabledImage && (patchMode != string(virtualmachinescalesets.WindowsVMGuestPatchModeAutomaticByPlatform)) {
						return fmt.Errorf("when referencing a hotpatching enabled image the 'patch_mode' field must always be set to %q, got %q", virtualmachinescalesets.WindowsVMGuestPatchModeAutomaticByPlatform, patchMode)
					}
					windowsConfig.PatchSettings.PatchMode = pointer.To(virtualmachinescalesets.WindowsVMGuestPatchMode(patchMode))
				}

				// Disabling hotpatching is not supported in images that support hotpatching
				// so while the attribute is exposed in VMSS it is hardcoded inside the images that
				// support hotpatching to always be enabled and cannot be set to false, ever.
				if d.HasChange("os_profile.0.windows_configuration.0.hotpatching_enabled") {
					if isHotpatchEnabledImage && !hotpatchingEnabled {
						return fmt.Errorf("when referencing a hotpatching enabled image the 'hotpatching_enabled' field must always be set to 'true', got %q", strconv.FormatBool(hotpatchingEnabled))
					}
					windowsConfig.PatchSettings.EnableHotpatching = pointer.To(hotpatchingEnabled)
				}

				if d.HasChange("os_profile.0.windows_configuration.0.secret") {
					vmssOsProfile.Secrets = expandWindowsSecretsVMSS(winConfig["secret"].([]interface{}))
				}

				if d.HasChange("os_profile.0.windows_configuration.0.timezone") {
					windowsConfig.TimeZone = pointer.To(winConfig["timezone"].(string))
				}

				if d.HasChange("os_profile.0.windows_configuration.0.winrm_listener") {
					winRmListenersRaw := winConfig["winrm_listener"].(*pluginsdk.Set).List()
					vmssOsProfile.WindowsConfiguration.WinRM = expandWinRMListenerVMSS(winRmListenersRaw)
				}

				vmssOsProfile.WindowsConfiguration = &windowsConfig
			}

			if len(linConfigRaw) > 0 {
				osType = virtualmachinescalesets.OperatingSystemTypesLinux
				linConfig := linConfigRaw[0].(map[string]interface{})
				provisionVMAgent := linConfig["provision_vm_agent"].(bool)
				patchAssessmentMode := linConfig["patch_assessment_mode"].(string)
				patchMode := linConfig["patch_mode"].(string)

				if d.HasChange("os_profile.0.linux_configuration.0.provision_vm_agent") ||
					d.HasChange("os_profile.0.linux_configuration.0.disable_password_authentication") ||
					d.HasChange("os_profile.0.linux_configuration.0.admin_ssh_key") {
					updateInstances = true
				}

				if d.HasChange("os_profile.0.linux_configuration.0.provision_vm_agent") {
					linuxConfig.ProvisionVMAgent = pointer.To(provisionVMAgent)
				}

				if d.HasChange("os_profile.0.linux_configuration.0.disable_password_authentication") {
					linuxConfig.DisablePasswordAuthentication = pointer.To(linConfig["disable_password_authentication"].(bool))
				}

				if d.HasChange("os_profile.0.linux_configuration.0.admin_ssh_key") {
					sshPublicKeys := expandSSHKeysVMSS(linConfig["admin_ssh_key"].(*pluginsdk.Set).List())
					if linuxConfig.Ssh == nil {
						linuxConfig.Ssh = &virtualmachinescalesets.SshConfiguration{}
					}
					linuxConfig.Ssh.PublicKeys = &sshPublicKeys
				}

				if d.HasChange("os_profile.0.linux_configuration.0.patch_assessment_mode") {
					if !provisionVMAgent && (patchAssessmentMode == string(virtualmachinescalesets.LinuxPatchAssessmentModeAutomaticByPlatform)) {
						return fmt.Errorf("when the 'patch_assessment_mode' field is set to %q the 'provision_vm_agent' must always be set to 'true'", virtualmachinescalesets.LinuxPatchAssessmentModeAutomaticByPlatform)
					}

					if linuxConfig.PatchSettings == nil {
						linuxConfig.PatchSettings = &virtualmachinescalesets.LinuxPatchSettings{}
					}
					linuxConfig.PatchSettings.AssessmentMode = pointer.To(virtualmachinescalesets.LinuxPatchAssessmentMode(patchAssessmentMode))
				}

				if d.HasChange("os_profile.0.linux_configuration.0.patch_mode") {
					if patchMode == string(virtualmachinescalesets.LinuxPatchAssessmentModeAutomaticByPlatform) {
						if !provisionVMAgent {
							return fmt.Errorf("when the 'patch_mode' field is set to %q the 'provision_vm_agent' field must always be set to 'true', got %q", patchMode, strconv.FormatBool(provisionVMAgent))
						}

						linuxAutomaticVMGuestPatchingEnabled = true
					}

					if linuxConfig.PatchSettings == nil {
						linuxConfig.PatchSettings = &virtualmachinescalesets.LinuxPatchSettings{}
					}
					linuxConfig.PatchSettings.PatchMode = pointer.To(virtualmachinescalesets.LinuxVMGuestPatchMode(patchMode))
				}

				vmssOsProfile.LinuxConfiguration = &linuxConfig
			}

			updateProps.VirtualMachineProfile.OsProfile = &vmssOsProfile
		}

		if d.HasChange("data_disk") || d.HasChange("os_disk") || d.HasChange("source_image_id") || d.HasChange("source_image_reference") {
			updateInstances = true

			if updateProps.VirtualMachineProfile.StorageProfile == nil {
				updateProps.VirtualMachineProfile.StorageProfile = &virtualmachinescalesets.VirtualMachineScaleSetUpdateStorageProfile{}
			}

			if d.HasChange("data_disk") {
				ultraSSDEnabled := false // Currently not supported in orchestrated vmss
				dataDisks, err := ExpandOrchestratedVirtualMachineScaleSetDataDisk(d.Get("data_disk").([]interface{}), ultraSSDEnabled)
				if err != nil {
					return fmt.Errorf("expanding `data_disk`: %+v", err)
				}
				updateProps.VirtualMachineProfile.StorageProfile.DataDisks = dataDisks
			}

			if d.HasChange("os_disk") {
				osDiskRaw := d.Get("os_disk").([]interface{})
				updateProps.VirtualMachineProfile.StorageProfile.OsDisk = ExpandOrchestratedVirtualMachineScaleSetOSDiskUpdate(osDiskRaw)
			}

			if d.HasChange("source_image_id") || d.HasChange("source_image_reference") {
				sourceImageReferenceRaw := d.Get("source_image_reference").([]interface{})
				sourceImageId := d.Get("source_image_id").(string)

				if len(sourceImageReferenceRaw) != 0 || sourceImageId != "" {
					sourceImageReference := expandSourceImageReferenceVMSS(sourceImageReferenceRaw, sourceImageId)
					updateProps.VirtualMachineProfile.StorageProfile.ImageReference = sourceImageReference
				}

				// Must include all storage profile properties when updating disk image.  See: https://github.com/hashicorp/terraform-provider-azurerm/issues/8273
				updateProps.VirtualMachineProfile.StorageProfile.DataDisks = existing.Model.Properties.VirtualMachineProfile.StorageProfile.DataDisks
				updateProps.VirtualMachineProfile.StorageProfile.OsDisk = &virtualmachinescalesets.VirtualMachineScaleSetUpdateOSDisk{
					Caching:                 existing.Model.Properties.VirtualMachineProfile.StorageProfile.OsDisk.Caching,
					WriteAcceleratorEnabled: existing.Model.Properties.VirtualMachineProfile.StorageProfile.OsDisk.WriteAcceleratorEnabled,
					DiskSizeGB:              existing.Model.Properties.VirtualMachineProfile.StorageProfile.OsDisk.DiskSizeGB,
					Image:                   existing.Model.Properties.VirtualMachineProfile.StorageProfile.OsDisk.Image,
					VhdContainers:           existing.Model.Properties.VirtualMachineProfile.StorageProfile.OsDisk.VhdContainers,
					ManagedDisk:             existing.Model.Properties.VirtualMachineProfile.StorageProfile.OsDisk.ManagedDisk,
				}
			}
		}

		if d.HasChange("network_interface") {
			networkInterfacesRaw := d.Get("network_interface").([]interface{})
			networkInterfaces, err := ExpandOrchestratedVirtualMachineScaleSetNetworkInterfaceUpdate(networkInterfacesRaw)
			if err != nil {
				return fmt.Errorf("expanding `network_interface`: %+v", err)
			}

			updateProps.VirtualMachineProfile.NetworkProfile = &virtualmachinescalesets.VirtualMachineScaleSetUpdateNetworkProfile{
				NetworkInterfaceConfigurations: networkInterfaces,
				// 2020-11-01 is the only valid value for this value and is only valid for VMSS in Orchestration Mode flex
				NetworkApiVersion: pointer.To(virtualmachinescalesets.NetworkApiVersionTwoZeroTwoZeroNegativeOneOneNegativeZeroOne),
			}
		}

		if d.HasChange("boot_diagnostics") {
			updateInstances = true

			bootDiagnosticsRaw := d.Get("boot_diagnostics").([]interface{})
			updateProps.VirtualMachineProfile.DiagnosticsProfile = expandBootDiagnosticsVMSS(bootDiagnosticsRaw)
		}

		if d.HasChange("termination_notification") {
			notificationRaw := d.Get("termination_notification").([]interface{})
			updateProps.VirtualMachineProfile.ScheduledEventsProfile = ExpandOrchestratedVirtualMachineScaleSetScheduledEventsProfile(notificationRaw)
		}

		if d.HasChange("encryption_at_host_enabled") {
			updateProps.VirtualMachineProfile.SecurityProfile = &virtualmachinescalesets.SecurityProfile{
				EncryptionAtHost: pointer.To(d.Get("encryption_at_host_enabled").(bool)),
			}
		}

		if d.HasChange("license_type") {
			license := d.Get("license_type").(string)
			if license == "" {
				// Only for create no specification is possible in the API. API does not allow empty string in update.
				// So removing attribute license_type from Terraform configuration if it was set to value other than 'None' would lead to an endless loop in apply.
				// To allow updating in this case set value explicitly to 'None'.
				license = "None"
			}
			updateProps.VirtualMachineProfile.LicenseType = &license
		}

		if d.HasChange("automatic_instance_repair") {
			automaticRepairsPolicyRaw := d.Get("automatic_instance_repair").([]interface{})
			automaticRepairsPolicy := ExpandVirtualMachineScaleSetAutomaticRepairsPolicy(automaticRepairsPolicyRaw)

			if automaticRepairsPolicy != nil {
				// we need to know if the VMSS has a health extension or not
				hasHealthExtension := false

				if v, ok := d.GetOk("extension"); ok {
					var err error
					_, hasHealthExtension, err = expandOrchestratedVirtualMachineScaleSetExtensions(v.(*pluginsdk.Set).List())
					if err != nil {
						return err
					}
				}

				if !hasHealthExtension {
					return fmt.Errorf("`automatic_instance_repair` can only be set if there is an application Health extension defined")
				}
			}
			updateProps.AutomaticRepairsPolicy = automaticRepairsPolicy
		}

		if d.HasChange("identity") {
			identityExpanded, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]interface{}))
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			update.Identity = identityExpanded
		}

		if d.HasChange("plan") {
			planRaw := d.Get("plan").([]interface{})
			update.Plan = expandPlanVMSS(planRaw)
		}

		if d.HasChange("sku_name") || d.HasChange("instances") {
			// in-case ignore_changes is being used, since both fields are required
			// look up the current values and override them as needed
			sku := existing.Model.Sku
			instances := int(*sku.Capacity)
			skuName := d.Get("sku_name").(string)

			if d.HasChange("instances") {
				instances = d.Get("instances").(int)

				sku, err = expandOrchestratedVirtualMachineScaleSetSku(skuName, instances)
				if err != nil {
					return err
				}
			}

			if d.HasChange("sku_name") {
				updateInstances = true

				sku, err = expandOrchestratedVirtualMachineScaleSetSku(skuName, instances)
				if err != nil {
					return err
				}
			}

			update.Sku = sku
		}

		if d.HasChanges("extension", "extensions_time_budget") {
			updateInstances = true

			extensionProfile, hasHealthExtension, err := expandOrchestratedVirtualMachineScaleSetExtensions(d.Get("extension").(*pluginsdk.Set).List())
			if err != nil {
				return err
			}

			if isHotpatchEnabledImage && !hasHealthExtension {
				return fmt.Errorf("when referencing a hotpatching enabled image the 'extension' field must always contain a 'application health extension'")
			}

			if linuxAutomaticVMGuestPatchingEnabled && !hasHealthExtension {
				return fmt.Errorf("when the 'patch_mode' field is set to %q the 'extension' field must contain at least one 'application health extension', got 0", virtualmachinescalesets.LinuxPatchAssessmentModeAutomaticByPlatform)
			}

			updateProps.VirtualMachineProfile.ExtensionProfile = extensionProfile
			updateProps.VirtualMachineProfile.ExtensionProfile.ExtensionsTimeBudget = pointer.To(d.Get("extensions_time_budget").(string))
		}
	}

	if d.HasChange("zones") {
		update.Zones = pointer.To(zones.ExpandUntyped(d.Get("zones").(*schema.Set).List()))
	}

	if d.HasChange("rolling_upgrade_policy") {
		upgradePolicy := virtualmachinescalesets.UpgradePolicy{}

		if existing.Model.Properties.UpgradePolicy != nil {
			upgradePolicy = *existing.Model.Properties.UpgradePolicy
		}

		upgradePolicy.Mode = pointer.To(virtualmachinescalesets.UpgradeMode(d.Get("upgrade_mode").(string)))

		rollingRaw := d.Get("rolling_upgrade_policy").([]interface{})
		rollingUpgradePolicy, err := ExpandVirtualMachineScaleSetRollingUpgradePolicy(rollingRaw, len(zones.ExpandUntyped(d.Get("zones").(*schema.Set).List())) > 0, false)
		if err != nil {
			return fmt.Errorf("expanding `rolling_upgrade_policy`: %+v", err)
		}

		upgradePolicy.RollingUpgradePolicy = rollingUpgradePolicy
		updateProps.UpgradePolicy = &upgradePolicy
	}

	// Only two fields that can change in legacy mode
	if d.HasChange("proximity_placement_group_id") {
		if v, ok := d.GetOk("proximity_placement_group_id"); ok {
			updateInstances = true
			updateProps.ProximityPlacementGroup = &virtualmachinescalesets.SubResource{
				Id: pointer.To(v.(string)),
			}
		}
	}

	if d.HasChange("tags") {
		update.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if d.HasChange("user_data_base64") {
		updateInstances = true
		updateProps.VirtualMachineProfile.UserData = pointer.To(d.Get("user_data_base64").(string))
	}

	update.Properties = &updateProps

	if updateInstances {
		log.Printf("[DEBUG] Orchestrated %s - updateInstances is true", id)
	}

	// AutomaticOSUpgradeIsEnabled currently is not supported in orchestrated VMSS flex
	metaData := virtualMachineScaleSetUpdateMetaData{
		AutomaticOSUpgradeIsEnabled:  false,
		CanReimageOnManualUpgrade:    false,
		CanRollInstancesWhenRequired: false,
		UpdateInstances:              false,
		Client:                       meta.(*clients.Client).Compute,
		Existing:                     pointer.From(existing.Model),
		ID:                           id,
		OSType:                       osType,
	}

	if err := metaData.performUpdate(ctx, update); err != nil {
		return err
	}

	return resourceOrchestratedVirtualMachineScaleSetRead(d, meta)
}

func resourceOrchestratedVirtualMachineScaleSetRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VirtualMachineScaleSetsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualmachinescalesets.ParseVirtualMachineScaleSetID(d.Id())
	if err != nil {
		return err
	}

	options := virtualmachinescalesets.DefaultGetOperationOptions()
	options.Expand = pointer.To(virtualmachinescalesets.ExpandTypesForGetVMScaleSetsUserData)
	resp, err := client.Get(ctx, *id, options)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] Orchestrated %s was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Orchestrated %s: %+v", id, err)
	}

	d.Set("name", id.VirtualMachineScaleSetName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))
		d.Set("zones", zones.FlattenUntyped(model.Zones))

		var skuName *string
		var instances int
		if model.Sku != nil {
			skuName, err = flattenOrchestratedVirtualMachineScaleSetSku(model.Sku)
			if err != nil || skuName == nil {
				return fmt.Errorf("setting `sku_name`: %+v", err)
			}

			if resp.Model.Sku.Capacity != nil {
				instances = int(*resp.Model.Sku.Capacity)
			}

			d.Set("sku_name", skuName)
			d.Set("instances", instances)
		}

		identityFlattened, err := flattenOrchestratedVirtualMachineScaleSetIdentity(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		if err := d.Set("identity", identityFlattened); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if err := d.Set("plan", flattenPlanVMSS(model.Plan)); err != nil {
			return fmt.Errorf("setting `plan`: %+v", err)
		}

		if props := model.Properties; props != nil {
			if err := d.Set("additional_capabilities", FlattenOrchestratedVirtualMachineScaleSetAdditionalCapabilities(props.AdditionalCapabilities)); err != nil {
				return fmt.Errorf("setting `additional_capabilities`: %+v", props.AdditionalCapabilities)
			}

			if err := d.Set("automatic_instance_repair", FlattenVirtualMachineScaleSetAutomaticRepairsPolicy(props.AutomaticRepairsPolicy)); err != nil {
				return fmt.Errorf("setting `automatic_instance_repair`: %+v", err)
			}

			d.Set("platform_fault_domain_count", props.PlatformFaultDomainCount)
			proximityPlacementGroupId := ""
			if props.ProximityPlacementGroup != nil && props.ProximityPlacementGroup.Id != nil {
				proximityPlacementGroupId = *props.ProximityPlacementGroup.Id
			}
			d.Set("proximity_placement_group_id", proximityPlacementGroupId)

			// only write state for single_placement_group if it is returned by the RP...
			if props.SinglePlacementGroup != nil {
				d.Set("single_placement_group", props.SinglePlacementGroup)
			}

			if err := d.Set("sku_profile", flattenOrchestratedVirtualMachineScaleSetSkuProfile(props.SkuProfile)); err != nil {
				return fmt.Errorf("setting `sku_profile`: %+v", err)
			}

			d.Set("unique_id", props.UniqueId)
			d.Set("zone_balance", props.ZoneBalance)

			extensionOperationsEnabled := true
			// if `VirtualMachineProfile` is nil, `UpgradeMode` will not exist in the response
			upgradeMode := string(virtualmachinescalesets.UpgradeModeManual)
			if profile := props.VirtualMachineProfile; profile != nil {
				if err := d.Set("boot_diagnostics", flattenBootDiagnosticsVMSS(profile.DiagnosticsProfile)); err != nil {
					return fmt.Errorf("setting `boot_diagnostics`: %+v", err)
				}

				capacityReservationGroupId := ""
				if profile.CapacityReservation != nil && profile.CapacityReservation.CapacityReservationGroup != nil && profile.CapacityReservation.CapacityReservationGroup.Id != nil {
					capacityReservationGroupId = *profile.CapacityReservation.CapacityReservationGroup.Id
				}
				d.Set("capacity_reservation_group_id", capacityReservationGroupId)

				// defaulted since BillingProfile isn't returned if it's unset
				maxBidPrice := float64(-1.0)
				if profile.BillingProfile != nil && profile.BillingProfile.MaxPrice != nil {
					maxBidPrice = *profile.BillingProfile.MaxPrice
				}
				d.Set("max_bid_price", maxBidPrice)

				d.Set("eviction_policy", pointer.From(profile.EvictionPolicy))
				d.Set("license_type", profile.LicenseType)

				// the service just return empty when this is not assigned when provisioned
				// See discussion on https://github.com/Azure/azure-rest-api-specs/issues/10971
				priority := virtualmachinescalesets.VirtualMachinePriorityTypesRegular
				if profile.Priority != nil {
					priority = pointer.From(profile.Priority)
				}
				d.Set("priority", priority)

				if storageProfile := profile.StorageProfile; storageProfile != nil {
					if err := d.Set("os_disk", FlattenOrchestratedVirtualMachineScaleSetOSDisk(storageProfile.OsDisk)); err != nil {
						return fmt.Errorf("setting `os_disk`: %+v", err)
					}

					if err := d.Set("data_disk", FlattenOrchestratedVirtualMachineScaleSetDataDisk(storageProfile.DataDisks)); err != nil {
						return fmt.Errorf("setting `data_disk`: %+v", err)
					}

					var storageImageId string
					if storageProfile.ImageReference != nil && storageProfile.ImageReference.Id != nil {
						storageImageId = *storageProfile.ImageReference.Id
					}
					if storageProfile.ImageReference != nil && storageProfile.ImageReference.CommunityGalleryImageId != nil {
						storageImageId = *storageProfile.ImageReference.CommunityGalleryImageId
					}
					if storageProfile.ImageReference != nil && storageProfile.ImageReference.SharedGalleryImageId != nil {
						storageImageId = *storageProfile.ImageReference.SharedGalleryImageId
					}
					d.Set("source_image_id", storageImageId)

					if err := d.Set("source_image_reference", flattenSourceImageReferenceVMSS(storageProfile.ImageReference, storageImageId != "")); err != nil {
						return fmt.Errorf("setting `source_image_reference`: %+v", err)
					}
				}

				if osProfile := profile.OsProfile; osProfile != nil {
					if err := d.Set("os_profile", FlattenOrchestratedVirtualMachineScaleSetOSProfile(osProfile, d)); err != nil {
						return fmt.Errorf("setting `os_profile`: %+v", err)
					}

					if osProfile.AllowExtensionOperations != nil {
						extensionOperationsEnabled = *osProfile.AllowExtensionOperations
					}
				}

				if nwProfile := profile.NetworkProfile; nwProfile != nil {
					flattenedNics := FlattenOrchestratedVirtualMachineScaleSetNetworkInterface(nwProfile.NetworkInterfaceConfigurations)
					if err := d.Set("network_interface", flattenedNics); err != nil {
						return fmt.Errorf("setting `network_interface`: %+v", err)
					}
				}

				if scheduleProfile := profile.ScheduledEventsProfile; scheduleProfile != nil {
					if err := d.Set("termination_notification", FlattenOrchestratedVirtualMachineScaleSetScheduledEventsProfile(scheduleProfile)); err != nil {
						return fmt.Errorf("setting `termination_notification`: %+v", err)
					}
				}

				extensionProfile, err := flattenOrchestratedVirtualMachineScaleSetExtensions(profile.ExtensionProfile, d)
				if err != nil {
					return fmt.Errorf("failed flattening `extension`: %+v", err)
				}
				d.Set("extension", extensionProfile)

				extensionsTimeBudget := "PT1H30M"
				if profile.ExtensionProfile != nil && profile.ExtensionProfile.ExtensionsTimeBudget != nil {
					extensionsTimeBudget = *profile.ExtensionProfile.ExtensionsTimeBudget
				}
				d.Set("extensions_time_budget", extensionsTimeBudget)

				encryptionAtHostEnabled := false
				if profile.SecurityProfile != nil && profile.SecurityProfile.EncryptionAtHost != nil {
					encryptionAtHostEnabled = *profile.SecurityProfile.EncryptionAtHost
				}
				d.Set("encryption_at_host_enabled", encryptionAtHostEnabled)
				d.Set("user_data_base64", profile.UserData)

				if policy := props.UpgradePolicy; policy != nil {
					upgradeMode = string(pointer.From(policy.Mode))
					flattenedRolling := FlattenVirtualMachineScaleSetRollingUpgradePolicy(policy.RollingUpgradePolicy)
					if err := d.Set("rolling_upgrade_policy", flattenedRolling); err != nil {
						return fmt.Errorf("setting `rolling_upgrade_policy`: %+v", err)
					}
				}
			}

			if priorityMixPolicy := props.PriorityMixPolicy; priorityMixPolicy != nil {
				if err := d.Set("priority_mix", FlattenOrchestratedVirtualMachineScaleSetPriorityMixPolicy(priorityMixPolicy)); err != nil {
					return fmt.Errorf("setting `priority_mix`: %+v", err)
				}
			}

			d.Set("extension_operations_enabled", extensionOperationsEnabled)
			d.Set("upgrade_mode", upgradeMode)
		}
		return tags.FlattenAndSet(d, model.Tags)
	}
	return nil
}

func resourceOrchestratedVirtualMachineScaleSetDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VirtualMachineScaleSetsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualmachinescalesets.ParseVirtualMachineScaleSetID(d.Id())
	if err != nil {
		return err
	}

	// Upgrading to the 2021-07-01 exposed a new expand parameter in the GET method
	resp, err := client.Get(ctx, *id, virtualmachinescalesets.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil
		}

		return fmt.Errorf("retrieving Orchestrated %s: %+v", id, err)
	}

	// Sometimes VMSS's aren't fully deleted when the `Delete` call returns - as such we'll try to scale the cluster
	// to 0 nodes first, then delete the cluster - which should ensure there's no Network Interfaces kicking around
	// and work around this Azure API bug:
	// Original Error: Code="InUseSubnetCannotBeDeleted" Message="Subnet internal is in use by
	// /{nicResourceID}/|providers|Microsoft.Compute|virtualMachineScaleSets|acctestvmss-190923101253410278|virtualMachines|0|networkInterfaces|example/ipConfigurations/internal and cannot be deleted.
	// In order to delete the subnet, delete all the resources within the subnet. See aka.ms/deletesubnet.
	if resp.Model.Sku != nil {
		resp.Model.Sku.Capacity = utils.Int64(int64(0))

		log.Printf("[DEBUG] Scaling instances to 0 prior to deletion - this helps avoids networking issues within Azure")
		update := virtualmachinescalesets.VirtualMachineScaleSetUpdate{
			Sku: resp.Model.Sku,
		}
		if err := client.UpdateThenPoll(ctx, *id, update, virtualmachinescalesets.DefaultUpdateOperationOptions()); err != nil {
			return fmt.Errorf("updating number of instances in Orchestrated %s to scale to 0: %+v", id, err)
		}
		log.Printf("[DEBUG] Scaled instances to 0 prior to deletion - this helps avoids networking issues within Azure")
	} else {
		log.Printf("[DEBUG] Unable to scale instances to `0` since the `sku` block is nil - trying to delete anyway")
	}

	log.Printf("[DEBUG] Deleting Orchestrated %s", id)
	// @ArcturusZhang (mimicking from windows_virtual_machine_pluginsdk.go): sending `nil` here omits this value from being sent
	// which matches the previous behaviour - we're only splitting this out so it's clear why
	if err = client.DeleteThenPoll(ctx, *id, virtualmachinescalesets.DefaultDeleteOperationOptions()); err != nil {
		return fmt.Errorf("deleting Orchestrated %s: %+v", id, err)
	}
	log.Printf("[DEBUG] Deleted Orchestrated %s", id)

	return nil
}

func expandOrchestratedVirtualMachineScaleSetSkuProfile(input []interface{}) *virtualmachinescalesets.SkuProfile {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})
	vmSizesRaw := v["vm_sizes"].(*pluginsdk.Set).List()
	vmSizes := make([]virtualmachinescalesets.SkuProfileVMSize, 0)
	for _, vmSize := range vmSizesRaw {
		vmSizes = append(vmSizes, virtualmachinescalesets.SkuProfileVMSize{
			Name: pointer.To(vmSize.(string)),
		})
	}

	return &virtualmachinescalesets.SkuProfile{
		AllocationStrategy: pointer.To((virtualmachinescalesets.AllocationStrategy)(v["allocation_strategy"].(string))),
		VMSizes:            pointer.To(vmSizes),
	}
}

func flattenOrchestratedVirtualMachineScaleSetSkuProfile(input *virtualmachinescalesets.SkuProfile) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	vmSizes := make([]string, 0)
	if input.VMSizes != nil {
		for _, vmSize := range *input.VMSizes {
			vmSizes = append(vmSizes, *vmSize.Name)
		}
	}

	return []interface{}{
		map[string]interface{}{
			"allocation_strategy": string(pointer.From(input.AllocationStrategy)),
			"vm_sizes":            vmSizes,
		},
	}
}

func expandOrchestratedVirtualMachineScaleSetSku(input string, capacity int) (*virtualmachinescalesets.Sku, error) {
	skuParts := strings.Split(input, "_")

	if (input != SkuNameMix && len(skuParts) < 2) || strings.Contains(input, "__") || strings.Contains(input, " ") {
		return nil, fmt.Errorf("'sku_name'(%q) is not formatted properly", input)
	}

	sku := &virtualmachinescalesets.Sku{
		Name:     pointer.To(input),
		Capacity: utils.Int64(int64(capacity)),
	}

	if input != SkuNameMix {
		sku.Tier = pointer.To("Standard")
	}

	return sku, nil
}

func flattenOrchestratedVirtualMachineScaleSetSku(input *virtualmachinescalesets.Sku) (*string, error) {
	var skuName string
	if input != nil && input.Name != nil {
		if strings.HasPrefix(strings.ToLower(*input.Name), "standard") || *input.Name == SkuNameMix {
			skuName = *input.Name
		} else {
			skuName = fmt.Sprintf("Standard_%s", *input.Name)
		}

		return &skuName, nil
	}

	return nil, fmt.Errorf("sku struct 'name' is nil")
}

func expandOrchestratedVirtualMachineScaleSetPublicIPSku(input string) *virtualmachinescalesets.PublicIPAddressSku {
	skuParts := strings.Split(input, "_")

	if len(skuParts) < 2 || strings.Contains(input, "__") || strings.Contains(input, " ") {
		return &virtualmachinescalesets.PublicIPAddressSku{}
	}

	return &virtualmachinescalesets.PublicIPAddressSku{
		Name: pointer.To(virtualmachinescalesets.PublicIPAddressSkuName(skuParts[0])),
		Tier: pointer.To(virtualmachinescalesets.PublicIPAddressSkuTier(skuParts[1])),
	}
}

func flattenOrchestratedVirtualMachineScaleSetPublicIPSku(input *virtualmachinescalesets.PublicIPAddressSku) string {
	var skuName string
	if input != nil {
		name := string(pointer.From(input.Name))
		tier := string(pointer.From(input.Tier))
		if name != "" && tier != "" {
			skuName = fmt.Sprintf("%s_%s", name, tier)
		}
	}

	return skuName
}
