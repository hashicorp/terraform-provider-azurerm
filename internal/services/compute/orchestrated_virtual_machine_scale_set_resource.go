// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/capacityreservationgroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/images"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/proximityplacementgroups"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	computeValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/compute/2023-03-01/compute"
)

func resourceOrchestratedVirtualMachineScaleSet() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceOrchestratedVirtualMachineScaleSetCreate,
		Read:   resourceOrchestratedVirtualMachineScaleSetRead,
		Update: resourceOrchestratedVirtualMachineScaleSetUpdate,
		Delete: resourceOrchestratedVirtualMachineScaleSetDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := parse.VirtualMachineScaleSetID(id)
			return err
		}, importOrchestratedVirtualMachineScaleSet),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		// The plan was to remove support the the legacy Orchestrated Virtual Machine Scale Set in 3.0.
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
					string(compute.VirtualMachineEvictionPolicyTypesDeallocate),
					string(compute.VirtualMachineEvictionPolicyTypesDelete),
				}, false),
			},

			"extension_operations_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default: func() interface{} {
					if !features.FourPointOhBeta() {
						return nil
					}
					return true
				}(),
				Computed: !features.FourPointOhBeta(),
				ForceNew: true,
			},

			// Due to bug in RP extensions cannot curretntly be supported in Terraform ETA for full support is mid Jan 2022
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
				Default:  string(compute.VirtualMachinePriorityTypesRegular),
				ValidateFunc: validation.StringInSlice([]string{
					string(compute.VirtualMachinePriorityTypesRegular),
					string(compute.VirtualMachinePriorityTypesSpot),
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

			"zones": commonschema.ZonesMultipleOptionalForceNew(),

			"tags": tags.Schema(),

			// Computed
			"unique_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"user_data_base64": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsBase64,
			},

			"priority_mix": OrchestratedVirtualMachineScaleSetPriorityMixPolicySchema(),
		},
	}
}

func resourceOrchestratedVirtualMachineScaleSetCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VMScaleSetClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	isLegacy := true
	id := parse.NewVirtualMachineScaleSetID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		// Upgrading to the 2021-07-01 exposed a new expand parameter to the GET method
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name, compute.ExpandTypesForGetVMScaleSetsUserData)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_orchestrated_virtual_machine_scale_set", id.ID())
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	props := compute.VirtualMachineScaleSet{
		Location: utils.String(location),
		Tags:     tags.Expand(t),
		VirtualMachineScaleSetProperties: &compute.VirtualMachineScaleSetProperties{
			PlatformFaultDomainCount: utils.Int32(int32(d.Get("platform_fault_domain_count").(int))),
			// OrchestrationMode needs to be hardcoded to Uniform, for the
			// standard VMSS resource, since virtualMachineProfile is now supported
			// in both VMSS and Orchestrated VMSS...
			OrchestrationMode: compute.OrchestrationModeFlexible,
		},
	}

	// The RP now accepts true, false and null for single_placement_group value.
	// This is only valid for the Orchestrated VMSS Resource. If the
	// single_placement_group is null(e.g. not passed in the props) the RP will
	// automatically determine what values single_placement_group should be
	if !pluginsdk.IsExplicitlyNullInConfig(d, "single_placement_group") {
		props.VirtualMachineScaleSetProperties.SinglePlacementGroup = utils.Bool(d.Get("single_placement_group").(bool))
	}

	zones := zones.ExpandUntyped(d.Get("zones").(*schema.Set).List())
	if len(zones) > 0 {
		props.Zones = &zones
	}

	virtualMachineProfile := compute.VirtualMachineScaleSetVMProfile{
		StorageProfile: &compute.VirtualMachineScaleSetStorageProfile{},
	}

	networkProfile := &compute.VirtualMachineScaleSetNetworkProfile{
		// 2020-11-01 is the only valid value for this value and is only valid for VMSS in Orchestration Mode flex
		NetworkAPIVersion: compute.NetworkAPIVersionTwoZeroTwoZeroHyphenMinusOneOneHyphenMinusZeroOne,
	}

	if v, ok := d.GetOk("proximity_placement_group_id"); ok {
		props.VirtualMachineScaleSetProperties.ProximityPlacementGroup = &compute.SubResource{
			ID: utils.String(v.(string)),
		}
	}

	// Not currently supported in OVMSS
	// healthProbeId := d.Get("health_probe_id").(string)
	// upgradeMode := compute.UpgradeMode(d.Get("upgrade_mode").(string))

	instances := d.Get("instances").(int)
	if v, ok := d.GetOk("sku_name"); ok {
		isLegacy = false
		sku, err := expandOrchestratedVirtualMachineScaleSetSku(v.(string), instances)
		if err != nil {
			return fmt.Errorf("expanding 'sku_name': %+v", err)
		}
		props.Sku = sku
	}

	if v, ok := d.GetOk("capacity_reservation_group_id"); ok {
		if d.Get("single_placement_group").(bool) {
			return fmt.Errorf("`single_placement_group` must be set to `false` when `capacity_reservation_group_id` is specified")
		}

		virtualMachineProfile.CapacityReservation = &compute.CapacityReservationProfile{
			CapacityReservationGroup: &compute.SubResource{
				ID: utils.String(v.(string)),
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
		log.Printf("[DEBUG] Orchestrated Virtual Machine Scale Set %q (Resource Group %q) has a Health Extension defined", id.Name, id.ResourceGroup)
	}

	if v, ok := d.GetOk("extensions_time_budget"); ok {
		if virtualMachineProfile.ExtensionProfile == nil {
			virtualMachineProfile.ExtensionProfile = &compute.VirtualMachineScaleSetExtensionProfile{}
		}
		virtualMachineProfile.ExtensionProfile.ExtensionsTimeBudget = utils.String(v.(string))
	}

	sourceImageReferenceRaw := d.Get("source_image_reference").([]interface{})
	sourceImageId := d.Get("source_image_id").(string)
	if len(sourceImageReferenceRaw) != 0 || sourceImageId != "" {
		sourceImageReference := expandSourceImageReference(sourceImageReferenceRaw, sourceImageId)
		virtualMachineProfile.StorageProfile.ImageReference = sourceImageReference
	}

	if userData, ok := d.GetOk("user_data_base64"); ok {
		virtualMachineProfile.UserData = utils.String(userData.(string))
	}

	osType := compute.OperatingSystemTypesWindows
	var winConfigRaw []interface{}
	var linConfigRaw []interface{}
	var vmssOsProfile *compute.VirtualMachineScaleSetOSProfile
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
				_, errs := computeValidate.WindowsComputerNamePrefix(id.Name, "computer_name_prefix")
				if len(errs) > 0 {
					return fmt.Errorf("unable to assume default computer name prefix %s. Please adjust the 'name', or specify an explicit 'computer_name_prefix'", errs[0])
				}
				vmssOsProfile.ComputerNamePrefix = utils.String(id.Name)
			}

			if extensionOperationsEnabled && !provisionVMAgent {
				return fmt.Errorf("`extension_operations_enabled` cannot be set to `true` when `provision_vm_agent` is set to `false`")
			}

			if patchAssessmentMode == string(compute.WindowsPatchAssessmentModeAutomaticByPlatform) && !provisionVMAgent {
				return fmt.Errorf("when the 'patch_assessment_mode' field is set to %q the 'provision_vm_agent' must always be set to 'true'", compute.WindowsPatchAssessmentModeAutomaticByPlatform)
			}

			// Validate patch mode and hotpatching configuration
			isHotpatchEnabledImage := isValidHotPatchSourceImageReference(sourceImageReferenceRaw, sourceImageId)
			patchMode := winConfig["patch_mode"].(string)
			hotpatchingEnabled := winConfig["hotpatching_enabled"].(bool)

			if isHotpatchEnabledImage {
				// it is a hotpatching enabled image, validate hotpatching enabled settings
				if patchMode != string(compute.WindowsVMGuestPatchModeAutomaticByPlatform) {
					return fmt.Errorf("when referencing a hotpatching enabled image the 'patch_mode' field must always be set to %q", compute.WindowsVMGuestPatchModeAutomaticByPlatform)
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
				if patchMode == string(compute.WindowsVMGuestPatchModeAutomaticByPlatform) {
					if !provisionVMAgent {
						return fmt.Errorf("when 'patch_mode' is set to %q then 'provision_vm_agent' must be set to 'true'", patchMode)
					}

					if !hasHealthExtension {
						return fmt.Errorf("when 'patch_mode' is set to %q then the 'extension' field must always always contain a 'application health extension'", patchMode)
					}
				}

				if hotpatchingEnabled {
					return fmt.Errorf("'hotpatching_enabled' field is not supported unless you are using one of the following hotpatching enable images, '2022-datacenter-azure-edition' or '2022-datacenter-azure-edition-core-smalldisk'")
				}
			}
		}

		if len(linConfigRaw) > 0 {
			osType = compute.OperatingSystemTypesLinux
			linConfig := linConfigRaw[0].(map[string]interface{})
			provisionVMAgent := linConfig["provision_vm_agent"].(bool)
			patchAssessmentMode := linConfig["patch_assessment_mode"].(string)
			vmssOsProfile = expandOrchestratedVirtualMachineScaleSetOsProfileWithLinuxConfiguration(linConfig, customData)

			// if the Computer Prefix Name was not defined use the computer name
			if vmssOsProfile.ComputerNamePrefix == nil || len(*vmssOsProfile.ComputerNamePrefix) == 0 {
				// validate that the computer name is a valid Computer Prefix Name
				_, errs := computeValidate.LinuxComputerNamePrefix(id.Name, "computer_name_prefix")
				if len(errs) > 0 {
					return fmt.Errorf("unable to assume default computer name prefix %s. Please adjust the 'name', or specify an explicit 'computer_name_prefix'", errs[0])
				}

				vmssOsProfile.ComputerNamePrefix = utils.String(id.Name)
			}

			if extensionOperationsEnabled && !provisionVMAgent {
				return fmt.Errorf("`extension_operations_enabled` cannot be set to `true` when `provision_vm_agent` is set to `false`")
			}

			if patchAssessmentMode == string(compute.LinuxPatchAssessmentModeAutomaticByPlatform) && !provisionVMAgent {
				return fmt.Errorf("when the 'patch_assessment_mode' field is set to %q the 'provision_vm_agent' must always be set to 'true'", compute.LinuxPatchAssessmentModeAutomaticByPlatform)
			}

			// Validate Automatic VM Guest Patching Settings
			patchMode := linConfig["patch_mode"].(string)

			if patchMode == string(compute.LinuxVMGuestPatchModeAutomaticByPlatform) {
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

	if !features.FourPointOhBeta() {
		if !pluginsdk.IsExplicitlyNullInConfig(d, "extension_operations_enabled") {
			if virtualMachineProfile.OsProfile == nil {
				virtualMachineProfile.OsProfile = &compute.VirtualMachineScaleSetOSProfile{}
			}
			virtualMachineProfile.OsProfile.AllowExtensionOperations = utils.Bool(extensionOperationsEnabled)
		}
	} else {
		if virtualMachineProfile.OsProfile == nil {
			virtualMachineProfile.OsProfile = &compute.VirtualMachineScaleSetOSProfile{}
		}
		virtualMachineProfile.OsProfile.AllowExtensionOperations = utils.Bool(extensionOperationsEnabled)
	}

	if v, ok := d.GetOk("boot_diagnostics"); ok {
		virtualMachineProfile.DiagnosticsProfile = expandBootDiagnostics(v.([]interface{}))
	}

	if v, ok := d.GetOk("priority"); ok {
		virtualMachineProfile.Priority = compute.VirtualMachinePriorityTypes(v.(string))
	}

	if v, ok := d.GetOk("os_disk"); ok {
		virtualMachineProfile.StorageProfile.OsDisk = ExpandOrchestratedVirtualMachineScaleSetOSDisk(v.([]interface{}), osType)
	}

	additionalCapabilitiesRaw := d.Get("additional_capabilities").([]interface{})
	additionalCapabilities := ExpandOrchestratedVirtualMachineScaleSetAdditionalCapabilities(additionalCapabilitiesRaw)
	props.VirtualMachineScaleSetProperties.AdditionalCapabilities = additionalCapabilities

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
		if virtualMachineProfile.Priority != compute.VirtualMachinePriorityTypesSpot {
			return fmt.Errorf("`max_bid_price` can only be configured when `priority` is set to `Spot`")
		}

		virtualMachineProfile.BillingProfile = &compute.BillingProfile{
			MaxPrice: utils.Float(v),
		}
	}

	if v, ok := d.GetOk("encryption_at_host_enabled"); ok {
		virtualMachineProfile.SecurityProfile = &compute.SecurityProfile{
			EncryptionAtHost: utils.Bool(v.(bool)),
		}
	}

	if v, ok := d.GetOk("eviction_policy"); ok {
		if virtualMachineProfile.Priority != compute.VirtualMachinePriorityTypesSpot {
			return fmt.Errorf("an `eviction_policy` can only be specified when `priority` is set to `Spot`")
		}
		virtualMachineProfile.EvictionPolicy = compute.VirtualMachineEvictionPolicyTypes(v.(string))
	} else if virtualMachineProfile.Priority == compute.VirtualMachinePriorityTypesSpot {
		return fmt.Errorf("an `eviction_policy` must be specified when `priority` is set to `Spot`")
	}

	if v, ok := d.GetOk("license_type"); ok {
		virtualMachineProfile.LicenseType = utils.String(v.(string))
	}

	if v, ok := d.GetOk("termination_notification"); ok {
		virtualMachineProfile.ScheduledEventsProfile = ExpandOrchestratedVirtualMachineScaleSetScheduledEventsProfile(v.([]interface{}))
	}

	// Only inclued the virtual machine profile if this is not a legacy configuration
	if !isLegacy {
		if v, ok := d.GetOk("plan"); ok {
			props.Plan = expandPlan(v.([]interface{}))
		}

		if v, ok := d.GetOk("identity"); ok {
			identity, err := expandVirtualMachineScaleSetIdentity(v.([]interface{}))
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}
			props.Identity = identity
		}

		if v, ok := d.GetOk("automatic_instance_repair"); ok {
			props.VirtualMachineScaleSetProperties.AutomaticRepairsPolicy = ExpandVirtualMachineScaleSetAutomaticRepairsPolicy(v.([]interface{}))
		}

		if v, ok := d.GetOk("zone_balance"); ok && v.(bool) {
			if props.Zones == nil || len(*props.Zones) == 0 {
				return fmt.Errorf("`zone_balance` can only be set to `true` when zones are specified")
			}

			props.VirtualMachineScaleSetProperties.ZoneBalance = utils.Bool(v.(bool))
		}

		if v, ok := d.GetOk("priority_mix"); ok {
			if virtualMachineProfile.Priority != compute.VirtualMachinePriorityTypesSpot {
				return fmt.Errorf("a `priority_mix` can only be specified when `priority` is set to `Spot`")
			}
			props.VirtualMachineScaleSetProperties.PriorityMixPolicy = ExpandOrchestratedVirtualMachineScaleSetPriorityMixPolicy(v.([]interface{}))
		}

		props.VirtualMachineScaleSetProperties.VirtualMachineProfile = &virtualMachineProfile
	}

	log.Printf("[DEBUG] Creating Orchestrated %s.", id)
	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, props)
	if err != nil {
		return fmt.Errorf("creating Orchestrated %s: %+v", id, err)
	}

	log.Printf("[DEBUG] Waiting for %s to be created.", id)
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		// If it is a Retryable Error re-issue the PUT for 10 loops until either the error goes away or the limit has been reached
		if strings.Contains(err.Error(), "RetryableError") {
			log.Printf("[DEBUG] Retryable error hit for %s to be created..", id)
			errCount := 1

			for {
				log.Printf("[DEBUG] Retrying PUT %d for Orchestrated %s.", errCount, id)
				future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, props)
				if err != nil {
					return fmt.Errorf("creating Orchestrated %s after %d retries: %+v", id, errCount, err)
				}

				err = future.WaitForCompletionRef(ctx, client.Client)
				if err != nil && strings.Contains(err.Error(), "RetryableError") {
					if errCount == 10 {
						return fmt.Errorf("waiting for creation of Orchestrated %s after %d retries: %+v", id, err, errCount)
					}
					errCount++
				} else {
					if err != nil {
						// Hit an error while retying that is not retryable anymore...
						return fmt.Errorf("hit unretryable error waiting for creation of Orchestrated %s after %d retries: %+v", id, err, errCount)
					} else {
						// err is nil and finally succeeded continue with the rest of the create function...
						break
					}
				}
			}
		} else {
			// Not a retryable error...
			return fmt.Errorf("waiting for creation of Orchestrated %s: %+v", id, err)
		}
	}

	log.Printf("[DEBUG] Orchestrated %s was created", id)
	log.Printf("[DEBUG] Retrieving Orchestrated %s.", id)

	d.SetId(id.ID())

	return resourceOrchestratedVirtualMachineScaleSetRead(d, meta)
}

func resourceOrchestratedVirtualMachineScaleSetUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VMScaleSetClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualMachineScaleSetID(d.Id())
	if err != nil {
		return err
	}

	isLegacy := true
	updateInstances := false
	isHotpatchEnabledImage := false
	linuxAutomaticVMGuestPatchingEnabled := false

	// retrieve
	// Upgrading to the 2021-07-01 exposed a new expand parameter in the GET method
	existing, err := client.Get(ctx, id.ResourceGroup, id.Name, compute.ExpandTypesForGetVMScaleSetsUserData)
	if err != nil {
		return fmt.Errorf("retrieving Orchestrated Virtual Machine Scale Set %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	if existing.Sku != nil {
		isLegacy = false
	}
	if existing.VirtualMachineScaleSetProperties == nil {
		return fmt.Errorf("retrieving Orchestrated Virtual Machine Scale Set %q (Resource Group %q): `properties` was nil", id.Name, id.ResourceGroup)
	}

	if !isLegacy {
		if existing.VirtualMachineScaleSetProperties.VirtualMachineProfile == nil {
			return fmt.Errorf("retrieving Orchestrated Virtual Machine Scale Set %q (Resource Group %q): `properties.virtualMachineProfile` was nil", id.Name, id.ResourceGroup)
		}
		if existing.VirtualMachineScaleSetProperties.VirtualMachineProfile.StorageProfile == nil {
			return fmt.Errorf("retrieving Orchestrated Virtual Machine Scale Set %q (Resource Group %q): `properties.virtualMachineProfile,storageProfile` was nil", id.Name, id.ResourceGroup)
		}
	}

	updateProps := compute.VirtualMachineScaleSetUpdateProperties{}
	update := compute.VirtualMachineScaleSetUpdate{}
	osType := compute.OperatingSystemTypesWindows

	if !isLegacy {
		updateProps = compute.VirtualMachineScaleSetUpdateProperties{
			VirtualMachineProfile: &compute.VirtualMachineScaleSetUpdateVMProfile{
				// if an image reference has been configured previously (it has to be), we would better to include that in this
				// update request to avoid some circumstances that the API will complain ImageReference is null
				// issue tracking: https://github.com/Azure/azure-rest-api-specs/issues/10322
				StorageProfile: &compute.VirtualMachineScaleSetUpdateStorageProfile{
					ImageReference: existing.VirtualMachineScaleSetProperties.VirtualMachineProfile.StorageProfile.ImageReference,
				},
			},
			// Currently not suppored in orchestrated VMSS
			// if an upgrade policy's been configured previously (which it will have) it must be threaded through
			// this doesn't matter for Manual - but breaks when updating anything on a Automatic and Rolling Mode Scale Set
			// UpgradePolicy: existing.VirtualMachineScaleSetProperties.UpgradePolicy,
		}

		priority := compute.VirtualMachinePriorityTypes(d.Get("priority").(string))

		if d.HasChange("single_placement_group") {
			// Since null is now a valid value for single_placement_group
			// make sure it is in the config file before you set the value
			// on the update props...
			if !pluginsdk.IsExplicitlyNullInConfig(d, "single_placement_group") {
				singlePlacementGroup := d.Get("single_placement_group").(bool)
				if singlePlacementGroup {
					return fmt.Errorf("'single_placement_group' can not be set to 'true' once it has been set to 'false'")
				}
				updateProps.SinglePlacementGroup = utils.Bool(singlePlacementGroup)
			}
		}

		if d.HasChange("max_bid_price") {
			if priority != compute.VirtualMachinePriorityTypesSpot {
				return fmt.Errorf("`max_bid_price` can only be configured when `priority` is set to `Spot`")
			}

			updateProps.VirtualMachineProfile.BillingProfile = &compute.BillingProfile{
				MaxPrice: utils.Float(d.Get("max_bid_price").(float64)),
			}
		}

		osProfileRaw := d.Get("os_profile").([]interface{})
		vmssOsProfile := compute.VirtualMachineScaleSetUpdateOSProfile{}
		windowsConfig := compute.WindowsConfiguration{}
		windowsConfig.PatchSettings = &compute.PatchSettings{}
		linuxConfig := compute.LinuxConfiguration{}

		if len(osProfileRaw) > 0 {
			osProfile := osProfileRaw[0].(map[string]interface{})
			winConfigRaw := osProfile["windows_configuration"].([]interface{})
			linConfigRaw := osProfile["linux_configuration"].([]interface{})

			if d.HasChange("os_profile.0.custom_data") {
				updateInstances = true

				// customData can only be sent if it's a base64 encoded string,
				// so it's not possible to remove this without tainting the resource
				vmssOsProfile.CustomData = utils.String(osProfile["custom_data"].(string))
			}

			if len(winConfigRaw) > 0 {
				winConfig := winConfigRaw[0].(map[string]interface{})
				provisionVMAgent := winConfig["provision_vm_agent"].(bool)
				patchAssessmentMode := winConfig["patch_assessment_mode"].(string)
				patchMode := winConfig["patch_mode"].(string)

				// If the image allows hotpatching the patch mode can only ever be AutomaticByPlatform.
				sourceImageReferenceRaw := d.Get("source_image_reference").([]interface{})
				sourceImageId := d.Get("source_image_id").(string)
				isHotpatchEnabledImage = isValidHotPatchSourceImageReference(sourceImageReferenceRaw, sourceImageId)

				if d.HasChange("os_profile.0.windows_configuration.0.enable_automatic_updates") ||
					d.HasChange("os_profile.0.windows_configuration.0.provision_vm_agent") ||
					d.HasChange("os_profile.0.windows_configuration.0.timezone") ||
					d.HasChange("os_profile.0.windows_configuration.0.secret") ||
					d.HasChange("os_profile.0.windows_configuration.0.winrm_listener") {
					updateInstances = true
				}

				if d.HasChange("os_profile.0.windows_configuration.0.enable_automatic_updates") {
					windowsConfig.EnableAutomaticUpdates = utils.Bool(winConfig["enable_automatic_updates"].(bool))
				}

				if d.HasChange("os_profile.0.windows_configuration.0.provision_vm_agent") {
					if isHotpatchEnabledImage && !provisionVMAgent {
						return fmt.Errorf("when referencing a hotpatching enabled image the 'provision_vm_agent' field must always be set to 'true', got %q", strconv.FormatBool(provisionVMAgent))
					}
					windowsConfig.ProvisionVMAgent = utils.Bool(provisionVMAgent)
				}

				if d.HasChange("os_profile.0.windows_configuration.0.patch_assessment_mode") {
					if !provisionVMAgent && (patchAssessmentMode == string(compute.WindowsPatchAssessmentModeAutomaticByPlatform)) {
						return fmt.Errorf("when the 'patch_assessment_mode' field is set to %q the 'provision_vm_agent' must always be set to 'true'", compute.WindowsPatchAssessmentModeAutomaticByPlatform)
					}
					windowsConfig.PatchSettings.AssessmentMode = compute.WindowsPatchAssessmentMode(patchAssessmentMode)
				}

				if d.HasChange("os_profile.0.windows_configuration.0.patch_mode") {
					if isHotpatchEnabledImage && (patchMode != string(compute.WindowsVMGuestPatchModeAutomaticByPlatform)) {
						return fmt.Errorf("when referencing a hotpatching enabled image the 'patch_mode' field must always be set to %q, got %q", compute.WindowsVMGuestPatchModeAutomaticByPlatform, patchMode)
					}
					windowsConfig.PatchSettings.PatchMode = compute.WindowsVMGuestPatchMode(patchMode)
				}

				// Disabling hotpatching is not supported in images that support hotpatching
				// so while the attribute is exposed in VMSS it is hardcoded inside the images that
				// support hotpatching to always be enabled and cannot be set to false, ever.
				if d.HasChange("os_profile.0.windows_configuration.0.hotpatching_enabled") {
					hotpatchingEnabled := winConfig["hotpatching_enabled"].(bool)
					if isHotpatchEnabledImage && !hotpatchingEnabled {
						return fmt.Errorf("when referencing a hotpatching enabled image the 'hotpatching_enabled' field must always be set to 'true', got %q", strconv.FormatBool(hotpatchingEnabled))
					}
					windowsConfig.PatchSettings.EnableHotpatching = utils.Bool(hotpatchingEnabled)
				}

				if d.HasChange("os_profile.0.windows_configuration.0.secret") {
					vmssOsProfile.Secrets = expandWindowsSecrets(winConfig["secret"].([]interface{}))
				}

				if d.HasChange("os_profile.0.windows_configuration.0.timezone") {
					windowsConfig.TimeZone = utils.String(winConfig["timezone"].(string))
				}

				if d.HasChange("os_profile.0.windows_configuration.0.winrm_listener") {
					winRmListenersRaw := winConfig["winrm_listener"].(*pluginsdk.Set).List()
					vmssOsProfile.WindowsConfiguration.WinRM = expandWinRMListener(winRmListenersRaw)
				}

				vmssOsProfile.WindowsConfiguration = &windowsConfig
			}

			if len(linConfigRaw) > 0 {
				osType = compute.OperatingSystemTypesLinux
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
					linuxConfig.ProvisionVMAgent = utils.Bool(provisionVMAgent)
				}

				if d.HasChange("os_profile.0.linux_configuration.0.disable_password_authentication") {
					linuxConfig.DisablePasswordAuthentication = utils.Bool(linConfig["disable_password_authentication"].(bool))
				}

				if d.HasChange("os_profile.0.linux_configuration.0.admin_ssh_key") {
					sshPublicKeys := ExpandSSHKeys(linConfig["admin_ssh_key"].(*pluginsdk.Set).List())
					if linuxConfig.SSH == nil {
						linuxConfig.SSH = &compute.SSHConfiguration{}
					}
					linuxConfig.SSH.PublicKeys = &sshPublicKeys
				}

				if d.HasChange("os_profile.0.linux_configuration.0.patch_assessment_mode") {
					if !provisionVMAgent && (patchAssessmentMode == string(compute.LinuxPatchAssessmentModeAutomaticByPlatform)) {
						return fmt.Errorf("when the 'patch_assessment_mode' field is set to %q the 'provision_vm_agent' must always be set to 'true'", compute.LinuxPatchAssessmentModeAutomaticByPlatform)
					}

					if linuxConfig.PatchSettings == nil {
						linuxConfig.PatchSettings = &compute.LinuxPatchSettings{}
					}
					linuxConfig.PatchSettings.AssessmentMode = compute.LinuxPatchAssessmentMode(patchAssessmentMode)
				}

				if d.HasChange("os_profile.0.linux_configuration.0.patch_mode") {
					if patchMode == string(compute.LinuxPatchAssessmentModeAutomaticByPlatform) {
						if !provisionVMAgent {
							return fmt.Errorf("when the 'patch_mode' field is set to %q the 'provision_vm_agent' field must always be set to 'true', got %q", patchMode, strconv.FormatBool(provisionVMAgent))
						}

						linuxAutomaticVMGuestPatchingEnabled = true
					}

					if linuxConfig.PatchSettings == nil {
						linuxConfig.PatchSettings = &compute.LinuxPatchSettings{}
					}
					linuxConfig.PatchSettings.PatchMode = compute.LinuxVMGuestPatchMode(patchMode)
				}

				vmssOsProfile.LinuxConfiguration = &linuxConfig
			}

			updateProps.VirtualMachineProfile.OsProfile = &vmssOsProfile
		}

		if d.HasChange("data_disk") || d.HasChange("os_disk") || d.HasChange("source_image_id") || d.HasChange("source_image_reference") {
			updateInstances = true

			if updateProps.VirtualMachineProfile.StorageProfile == nil {
				updateProps.VirtualMachineProfile.StorageProfile = &compute.VirtualMachineScaleSetUpdateStorageProfile{}
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
					sourceImageReference := expandSourceImageReference(sourceImageReferenceRaw, sourceImageId)
					updateProps.VirtualMachineProfile.StorageProfile.ImageReference = sourceImageReference
				}

				// Must include all storage profile properties when updating disk image.  See: https://github.com/hashicorp/terraform-provider-azurerm/issues/8273
				updateProps.VirtualMachineProfile.StorageProfile.DataDisks = existing.VirtualMachineScaleSetProperties.VirtualMachineProfile.StorageProfile.DataDisks
				updateProps.VirtualMachineProfile.StorageProfile.OsDisk = &compute.VirtualMachineScaleSetUpdateOSDisk{
					Caching:                 existing.VirtualMachineScaleSetProperties.VirtualMachineProfile.StorageProfile.OsDisk.Caching,
					WriteAcceleratorEnabled: existing.VirtualMachineScaleSetProperties.VirtualMachineProfile.StorageProfile.OsDisk.WriteAcceleratorEnabled,
					DiskSizeGB:              existing.VirtualMachineScaleSetProperties.VirtualMachineProfile.StorageProfile.OsDisk.DiskSizeGB,
					Image:                   existing.VirtualMachineScaleSetProperties.VirtualMachineProfile.StorageProfile.OsDisk.Image,
					VhdContainers:           existing.VirtualMachineScaleSetProperties.VirtualMachineProfile.StorageProfile.OsDisk.VhdContainers,
					ManagedDisk:             existing.VirtualMachineScaleSetProperties.VirtualMachineProfile.StorageProfile.OsDisk.ManagedDisk,
				}
			}
		}

		if d.HasChange("network_interface") {
			networkInterfacesRaw := d.Get("network_interface").([]interface{})
			networkInterfaces, err := ExpandOrchestratedVirtualMachineScaleSetNetworkInterfaceUpdate(networkInterfacesRaw)
			if err != nil {
				return fmt.Errorf("expanding `network_interface`: %+v", err)
			}

			updateProps.VirtualMachineProfile.NetworkProfile = &compute.VirtualMachineScaleSetUpdateNetworkProfile{
				NetworkInterfaceConfigurations: networkInterfaces,
				// 2020-11-01 is the only valid value for this value and is only valid for VMSS in Orchestration Mode flex
				NetworkAPIVersion: compute.NetworkAPIVersionTwoZeroTwoZeroHyphenMinusOneOneHyphenMinusZeroOne,
			}
		}

		if d.HasChange("boot_diagnostics") {
			updateInstances = true

			bootDiagnosticsRaw := d.Get("boot_diagnostics").([]interface{})
			updateProps.VirtualMachineProfile.DiagnosticsProfile = expandBootDiagnostics(bootDiagnosticsRaw)
		}

		if d.HasChange("termination_notification") {
			notificationRaw := d.Get("termination_notification").([]interface{})
			updateProps.VirtualMachineProfile.ScheduledEventsProfile = ExpandOrchestratedVirtualMachineScaleSetScheduledEventsProfile(notificationRaw)
		}

		if d.HasChange("encryption_at_host_enabled") {
			updateProps.VirtualMachineProfile.SecurityProfile = &compute.SecurityProfile{
				EncryptionAtHost: utils.Bool(d.Get("encryption_at_host_enabled").(bool)),
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
			updateProps.AutomaticRepairsPolicy = automaticRepairsPolicy
		}

		if d.HasChange("identity") {
			identity, err := expandVirtualMachineScaleSetIdentity(d.Get("identity").([]interface{}))
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			update.Identity = identity
		}

		if d.HasChange("plan") {
			planRaw := d.Get("plan").([]interface{})
			update.Plan = expandPlan(planRaw)
		}

		if d.HasChange("sku_name") || d.HasChange("instances") {
			// in-case ignore_changes is being used, since both fields are required
			// look up the current values and override them as needed
			sku := existing.Sku
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
				return fmt.Errorf("when the 'patch_mode' field is set to %q the 'extension' field must contain at least one 'application health extension', got 0", compute.LinuxPatchAssessmentModeAutomaticByPlatform)
			}

			updateProps.VirtualMachineProfile.ExtensionProfile = extensionProfile
			updateProps.VirtualMachineProfile.ExtensionProfile.ExtensionsTimeBudget = utils.String(d.Get("extensions_time_budget").(string))
		}
	}

	// Only two fields that can change in legacy mode
	if d.HasChange("proximity_placement_group_id") {
		if v, ok := d.GetOk("proximity_placement_group_id"); ok {
			updateInstances = true
			updateProps.ProximityPlacementGroup = &compute.SubResource{
				ID: utils.String(v.(string)),
			}
		}
	}

	if d.HasChange("tags") {
		update.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if d.HasChange("user_data_base64") {
		updateInstances = true
		updateProps.VirtualMachineProfile.UserData = utils.String(d.Get("user_data_base64").(string))
	}

	update.VirtualMachineScaleSetUpdateProperties = &updateProps

	if updateInstances {
		log.Printf("[DEBUG] Orchestrated Virtual Machine Scale Set %q in Resource Group %q - updateInstances is true", id.Name, id.ResourceGroup)
	}

	// AutomaticOSUpgradeIsEnabled currently is not supported in orchestrated VMSS flex
	metaData := virtualMachineScaleSetUpdateMetaData{
		AutomaticOSUpgradeIsEnabled: false,
		// CanRollInstancesWhenRequired: meta.(*clients.Client).Features.VirtualMachineScaleSet.RollInstancesWhenRequired,
		// UpdateInstances:              updateInstances,
		CanRollInstancesWhenRequired: false,
		UpdateInstances:              false,
		Client:                       meta.(*clients.Client).Compute,
		Existing:                     existing,
		ID:                           id,
		OSType:                       osType,
	}

	if err := metaData.performUpdate(ctx, update); err != nil {
		return err
	}

	return resourceOrchestratedVirtualMachineScaleSetRead(d, meta)
}

func resourceOrchestratedVirtualMachineScaleSetRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VMScaleSetClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualMachineScaleSetID(d.Id())
	if err != nil {
		return err
	}

	// Upgrading to the 2021-07-01 exposed a new expand parameter in the GET method
	resp, err := client.Get(ctx, id.ResourceGroup, id.Name, compute.ExpandTypesForGetVMScaleSetsUserData)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Orchestrated Virtual Machine Scale Set %q was not found in Resource Group %q - removing from state!", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Orchestrated Virtual Machine Scale Set %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))
	d.Set("zones", zones.FlattenUntyped(resp.Zones))

	var skuName *string
	var instances int
	if resp.Sku != nil {
		skuName, err = flattenOrchestratedVirtualMachineScaleSetSku(resp.Sku)
		if err != nil || skuName == nil {
			return fmt.Errorf("setting `sku_name`: %+v", err)
		}

		if resp.Sku.Capacity != nil {
			instances = int(*resp.Sku.Capacity)
		}

		d.Set("sku_name", skuName)
		d.Set("instances", instances)
	}

	identity, err := flattenOrchestratedVirtualMachineScaleSetIdentity(resp.Identity)
	if err != nil {
		return fmt.Errorf("flattening `identity`: %+v", err)
	}
	if err := d.Set("identity", identity); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	if err := d.Set("plan", flattenPlan(resp.Plan)); err != nil {
		return fmt.Errorf("setting `plan`: %+v", err)
	}

	if resp.VirtualMachineScaleSetProperties == nil {
		return fmt.Errorf("retrieving Orchestrated Virtual Machine Scale Set %q (Resource Group %q): `properties` was nil", id.Name, id.ResourceGroup)
	}
	props := *resp.VirtualMachineScaleSetProperties

	if err := d.Set("additional_capabilities", FlattenOrchestratedVirtualMachineScaleSetAdditionalCapabilities(props.AdditionalCapabilities)); err != nil {
		return fmt.Errorf("setting `additional_capabilities`: %+v", props.AdditionalCapabilities)
	}

	if err := d.Set("automatic_instance_repair", FlattenVirtualMachineScaleSetAutomaticRepairsPolicy(props.AutomaticRepairsPolicy)); err != nil {
		return fmt.Errorf("setting `automatic_instance_repair`: %+v", err)
	}

	d.Set("platform_fault_domain_count", props.PlatformFaultDomainCount)
	proximityPlacementGroupId := ""
	if props.ProximityPlacementGroup != nil && props.ProximityPlacementGroup.ID != nil {
		proximityPlacementGroupId = *props.ProximityPlacementGroup.ID
	}
	d.Set("proximity_placement_group_id", proximityPlacementGroupId)

	// only write state for single_placement_group if it is returned by the RP...
	if props.SinglePlacementGroup != nil {
		d.Set("single_placement_group", props.SinglePlacementGroup)
	}
	d.Set("unique_id", props.UniqueID)
	d.Set("zone_balance", props.ZoneBalance)

	extensionOperationsEnabled := true
	if profile := props.VirtualMachineProfile; profile != nil {
		if err := d.Set("boot_diagnostics", flattenBootDiagnostics(profile.DiagnosticsProfile)); err != nil {
			return fmt.Errorf("setting `boot_diagnostics`: %+v", err)
		}

		capacityReservationGroupId := ""
		if profile.CapacityReservation != nil && profile.CapacityReservation.CapacityReservationGroup != nil && profile.CapacityReservation.CapacityReservationGroup.ID != nil {
			capacityReservationGroupId = *profile.CapacityReservation.CapacityReservationGroup.ID
		}
		d.Set("capacity_reservation_group_id", capacityReservationGroupId)

		// defaulted since BillingProfile isn't returned if it's unset
		maxBidPrice := float64(-1.0)
		if profile.BillingProfile != nil && profile.BillingProfile.MaxPrice != nil {
			maxBidPrice = *profile.BillingProfile.MaxPrice
		}
		d.Set("max_bid_price", maxBidPrice)

		d.Set("eviction_policy", string(profile.EvictionPolicy))
		d.Set("license_type", profile.LicenseType)

		// the service just return empty when this is not assigned when provisioned
		// See discussion on https://github.com/Azure/azure-rest-api-specs/issues/10971
		priority := compute.VirtualMachinePriorityTypesRegular
		if profile.Priority != "" {
			priority = profile.Priority
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
			if storageProfile.ImageReference != nil && storageProfile.ImageReference.ID != nil {
				storageImageId = *storageProfile.ImageReference.ID
			}
			if storageProfile.ImageReference != nil && storageProfile.ImageReference.CommunityGalleryImageID != nil {
				storageImageId = *storageProfile.ImageReference.CommunityGalleryImageID
			}
			if storageProfile.ImageReference != nil && storageProfile.ImageReference.SharedGalleryImageID != nil {
				storageImageId = *storageProfile.ImageReference.SharedGalleryImageID
			}
			d.Set("source_image_id", storageImageId)

			if err := d.Set("source_image_reference", flattenSourceImageReference(storageProfile.ImageReference, storageImageId != "")); err != nil {
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
	}

	if priorityMixPolicy := props.PriorityMixPolicy; priorityMixPolicy != nil {
		if err := d.Set("priority_mix", FlattenOrchestratedVirtualMachineScaleSetPriorityMixPolicy(priorityMixPolicy)); err != nil {
			return fmt.Errorf("setting `priority_mix`: %+v", err)
		}
	}

	d.Set("extension_operations_enabled", extensionOperationsEnabled)

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceOrchestratedVirtualMachineScaleSetDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VMScaleSetClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualMachineScaleSetID(d.Id())
	if err != nil {
		return err
	}

	// Upgrading to the 2021-07-01 exposed a new expand parameter in the GET method
	resp, err := client.Get(ctx, id.ResourceGroup, id.Name, compute.ExpandTypesForGetVMScaleSetsUserData)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return nil
		}

		return fmt.Errorf("retrieving Orchestrated Virtual Machine Scale Set %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	// Sometimes VMSS's aren't fully deleted when the `Delete` call returns - as such we'll try to scale the cluster
	// to 0 nodes first, then delete the cluster - which should ensure there's no Network Interfaces kicking around
	// and work around this Azure API bug:
	// Original Error: Code="InUseSubnetCannotBeDeleted" Message="Subnet internal is in use by
	// /{nicResourceID}/|providers|Microsoft.Compute|virtualMachineScaleSets|acctestvmss-190923101253410278|virtualMachines|0|networkInterfaces|example/ipConfigurations/internal and cannot be deleted.
	// In order to delete the subnet, delete all the resources within the subnet. See aka.ms/deletesubnet.
	if resp.Sku != nil {
		resp.Sku.Capacity = utils.Int64(int64(0))

		log.Printf("[DEBUG] Scaling instances to 0 prior to deletion - this helps avoids networking issues within Azure")
		update := compute.VirtualMachineScaleSetUpdate{
			Sku: resp.Sku,
		}
		future, err := client.Update(ctx, id.ResourceGroup, id.Name, update)
		if err != nil {
			return fmt.Errorf("updating number of instances in Orchestrated Virtual Machine Scale Set %q (Resource Group %q) to scale to 0: %+v", id.Name, id.ResourceGroup, err)
		}

		log.Printf("[DEBUG] Waiting for scaling of instances to 0 prior to deletion - this helps avoids networking issues within Azure")
		err = future.WaitForCompletionRef(ctx, client.Client)
		if err != nil {
			return fmt.Errorf("waiting for number of instances in Orchestrated Virtual Machine Scale Set %q (Resource Group %q) to scale to 0: %+v", id.Name, id.ResourceGroup, err)
		}
		log.Printf("[DEBUG] Scaled instances to 0 prior to deletion - this helps avoids networking issues within Azure")
	} else {
		log.Printf("[DEBUG] Unable to scale instances to `0` since the `sku` block is nil - trying to delete anyway")
	}

	log.Printf("[DEBUG] Deleting Orchestrated Virtual Machine Scale Set %q (Resource Group %q)..", id.Name, id.ResourceGroup)
	// @ArcturusZhang (mimicking from windows_virtual_machine_pluginsdk.go): sending `nil` here omits this value from being sent
	// which matches the previous behaviour - we're only splitting this out so it's clear why
	// TODO: support force deletion once it's out of Preview, if applicable
	var forceDeletion *bool = nil
	future, err := client.Delete(ctx, id.ResourceGroup, id.Name, forceDeletion)
	if err != nil {
		return fmt.Errorf("deleting Orchestrated Virtual Machine Scale Set %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	log.Printf("[DEBUG] Waiting for deletion of Orchestrated Virtual Machine Scale Set %q (Resource Group %q)..", id.Name, id.ResourceGroup)
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of Orchestrated Virtual Machine Scale Set %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	log.Printf("[DEBUG] Deleted Orchestrated Virtual Machine Scale Set %q (Resource Group %q).", id.Name, id.ResourceGroup)

	return nil
}

func expandOrchestratedVirtualMachineScaleSetSku(input string, capacity int) (*compute.Sku, error) {
	skuParts := strings.Split(input, "_")

	if len(skuParts) < 2 || strings.Contains(input, "__") || strings.Contains(input, " ") {
		return nil, fmt.Errorf("'sku_name'(%q) is not formatted properly", input)
	}

	sku := &compute.Sku{
		Name:     utils.String(input),
		Capacity: utils.Int64(int64(capacity)),
		Tier:     utils.String("Standard"),
	}

	return sku, nil
}

func flattenOrchestratedVirtualMachineScaleSetSku(input *compute.Sku) (*string, error) {
	var skuName string
	if input != nil && input.Name != nil {
		if strings.HasPrefix(strings.ToLower(*input.Name), "standard") {
			skuName = *input.Name
		} else {
			skuName = fmt.Sprintf("Standard_%s", *input.Name)
		}

		return &skuName, nil
	}

	return nil, fmt.Errorf("sku struct 'name' is nil")
}

func expandOrchestratedVirtualMachineScaleSetPublicIPSku(input string) *compute.PublicIPAddressSku {
	skuParts := strings.Split(input, "_")

	if len(skuParts) < 2 || strings.Contains(input, "__") || strings.Contains(input, " ") {
		return &compute.PublicIPAddressSku{}
	}

	return &compute.PublicIPAddressSku{
		Name: compute.PublicIPAddressSkuName(skuParts[0]),
		Tier: compute.PublicIPAddressSkuTier(skuParts[1]),
	}
}

func flattenOrchestratedVirtualMachineScaleSetPublicIPSku(input *compute.PublicIPAddressSku) string {
	var skuName string
	if input != nil {
		if string(input.Name) != "" && string(input.Tier) != "" {
			skuName = fmt.Sprintf("%s_%s", string(input.Name), string(input.Tier))
		}
	}

	return skuName
}
