// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package orchestrated_virtual_machine_scale_set

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2025-04-01/virtualmachinescalesets"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceOrchestratedVirtualMachineScaleSetCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Compute.VirtualMachineScaleSetsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	isLegacy := true
	id := virtualmachinescalesets.NewVirtualMachineScaleSetID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id, virtualmachinescalesets.DefaultGetOperationOptions())
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %w", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_orchestrated_virtual_machine_scale_set", id.ID())
		}
	}

	t := d.Get("tags").(map[string]any)

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
	rollingUpgradePolicy, err := helpers.ExpandVirtualMachineScaleSetRollingUpgradePolicy(d.Get("rolling_upgrade_policy").([]any), len(zones) > 0, false)
	if err != nil {
		return fmt.Errorf("expanding `rolling_upgrade_policy`: %w", err)
	}

	props.Properties.UpgradePolicy = &virtualmachinescalesets.UpgradePolicy{
		Mode:                 pointer.To(upgradeMode),
		RollingUpgradePolicy: rollingUpgradePolicy,
	}

	virtualMachineProfile := virtualmachinescalesets.VirtualMachineScaleSetVMProfile{
		StorageProfile: &virtualmachinescalesets.VirtualMachineScaleSetStorageProfile{},
	}

	networkProfile := &virtualmachinescalesets.VirtualMachineScaleSetNetworkProfile{
		NetworkApiVersion: pointer.ToEnum[virtualmachinescalesets.NetworkApiVersion](d.Get("network_api_version").(string)),
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
			return fmt.Errorf("expanding `sku_name`: %w", err)
		}
		props.Sku = sku
	}

	if v, ok := d.GetOk("sku_profile"); ok {
		props.Properties.SkuProfile = expandOrchestratedVirtualMachineScaleSetSkuProfile(v.([]any))
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
		virtualMachineProfile.ExtensionProfile, hasHealthExtension, err = helpers.ExpandOrchestratedVirtualMachineScaleSetExtensions(v.(*pluginsdk.Set).List())
		if err != nil {
			return err
		}
	}

	if hasHealthExtension {
		log.Printf("[DEBUG] Orchestrated %s has a Health Extension defined", id)
	}

	// Virtual Machine Scale Set with Flexible Orchestration Mode and 'Rolling' upgradeMode must have Health Extension Present
	if upgradeMode == virtualmachinescalesets.UpgradeModeRolling && !hasHealthExtension {
		return fmt.Errorf("health extension is required when `upgrade_mode` is set to `%s`", string(upgradeMode))
	}

	if v, ok := d.GetOk("extensions_time_budget"); ok {
		if virtualMachineProfile.ExtensionProfile == nil {
			virtualMachineProfile.ExtensionProfile = &virtualmachinescalesets.VirtualMachineScaleSetExtensionProfile{}
		}
		virtualMachineProfile.ExtensionProfile.ExtensionsTimeBudget = pointer.To(v.(string))
	}

	sourceImageReferenceRaw := d.Get("source_image_reference").([]any)
	sourceImageId := d.Get("source_image_id").(string)
	if len(sourceImageReferenceRaw) != 0 || sourceImageId != "" {
		virtualMachineProfile.StorageProfile.ImageReference = helpers.ExpandSourceImageReferenceVMSS(sourceImageReferenceRaw, sourceImageId)
	}

	if userData, ok := d.GetOk("user_data_base64"); ok {
		virtualMachineProfile.UserData = pointer.To(userData.(string))
	}

	var osType virtualmachinescalesets.OperatingSystemTypes
	var winConfigRaw []any
	var linConfigRaw []any
	var vmssOsProfile *virtualmachinescalesets.VirtualMachineScaleSetOSProfile
	extensionOperationsEnabled := d.Get("extension_operations_enabled").(bool)
	osProfileRaw := d.Get("os_profile").([]any)

	if len(osProfileRaw) > 0 && osProfileRaw[0] != nil {
		osProfile := osProfileRaw[0].(map[string]any)
		winConfigRaw = osProfile["windows_configuration"].([]any)
		linConfigRaw = osProfile["linux_configuration"].([]any)
		customData := ""

		// Pass custom data if it is defined in the config file
		if v := osProfile["custom_data"]; v != nil {
			customData = v.(string)
		}

		if len(winConfigRaw) > 0 && winConfigRaw[0] != nil {
			osType = virtualmachinescalesets.OperatingSystemTypesWindows
			winConfig := winConfigRaw[0].(map[string]any)
			provisionVMAgent := winConfig["provision_vm_agent"].(bool)
			patchAssessmentMode := winConfig["patch_assessment_mode"].(string)
			vmssOsProfile = helpers.ExpandOrchestratedVirtualMachineScaleSetOsProfileWithWindowsConfiguration(winConfig, customData)

			// if the Computer Prefix Name was not defined use the computer name
			if vmssOsProfile.ComputerNamePrefix == nil || len(*vmssOsProfile.ComputerNamePrefix) == 0 {
				// validate that the computer name is a valid Computer Prefix Name
				_, errs := validate.WindowsComputerNamePrefix(id.VirtualMachineScaleSetName, "computer_name_prefix")
				if len(errs) > 0 {
					return fmt.Errorf("unable to assume default computer name prefix %s. Please adjust the `name`, or specify an explicit `computer_name_prefix`", errs[0])
				}
				vmssOsProfile.ComputerNamePrefix = pointer.To(id.VirtualMachineScaleSetName)
			}

			if extensionOperationsEnabled && !provisionVMAgent {
				return fmt.Errorf("`extension_operations_enabled` cannot be set to `true` when `provision_vm_agent` is set to `false`")
			}

			if patchAssessmentMode == string(virtualmachinescalesets.WindowsPatchAssessmentModeAutomaticByPlatform) && !provisionVMAgent {
				return fmt.Errorf("when `patch_assessment_mode` is set to `%s`, `provision_vm_agent` must be set to `true`", virtualmachinescalesets.WindowsPatchAssessmentModeAutomaticByPlatform)
			}

			// Validate patch mode and hotpatching configuration
			isHotpatchEnabledImage := helpers.IsValidHotPatchSourceImageReference(sourceImageReferenceRaw, sourceImageId)
			patchMode := winConfig["patch_mode"].(string)
			hotpatchingEnabled := winConfig["hotpatching_enabled"].(bool)

			if isHotpatchEnabledImage {
				// it is a hotpatching enabled image, validate hotpatching enabled settings
				if patchMode != string(virtualmachinescalesets.WindowsVMGuestPatchModeAutomaticByPlatform) {
					return fmt.Errorf("when using a hotpatching enabled image, `patch_mode` must be set to `%s`", virtualmachinescalesets.WindowsVMGuestPatchModeAutomaticByPlatform)
				}

				if !provisionVMAgent {
					return fmt.Errorf("when using a hotpatching enabled image, `provision_vm_agent` must be set to `true`")
				}

				if !hasHealthExtension {
					return fmt.Errorf("when using a hotpatching enabled image, an application health extension must be configured")
				}

				if !hotpatchingEnabled {
					return fmt.Errorf("when using a hotpatching enabled image, `hotpatching_enabled` must be set to `true`")
				}
			} else {
				// not a hotpatching enabled image verify Automatic VM Guest Patching settings
				if patchMode == string(virtualmachinescalesets.WindowsVMGuestPatchModeAutomaticByPlatform) {
					if !provisionVMAgent {
						return fmt.Errorf("when `patch_mode` is set to `%s`, `provision_vm_agent` must be set to `true`", patchMode)
					}

					if !hasHealthExtension {
						return fmt.Errorf("when `patch_mode` is set to `%s`, an application health extension must be configured", patchMode)
					}
				}

				if hotpatchingEnabled {
					return fmt.Errorf("`hotpatching_enabled` can only be used with supported Windows Server images: '2022-datacenter-azure-edition', '2022-datacenter-azure-edition-core-smalldisk', '2022-datacenter-azure-edition-hotpatch', '2022-datacenter-azure-edition-hotpatch-smalldisk', '2025-datacenter-azure-edition', '2025-datacenter-azure-edition-smalldisk', '2025-datacenter-azure-edition-core', or '2025-datacenter-azure-edition-core-smalldisk'")
				}
			}
		}

		if len(linConfigRaw) > 0 && linConfigRaw[0] != nil {
			osType = virtualmachinescalesets.OperatingSystemTypesLinux
			linConfig := linConfigRaw[0].(map[string]any)
			provisionVMAgent := linConfig["provision_vm_agent"].(bool)
			patchAssessmentMode := linConfig["patch_assessment_mode"].(string)
			vmssOsProfile = helpers.ExpandOrchestratedVirtualMachineScaleSetOsProfileWithLinuxConfiguration(linConfig, customData)

			// if the Computer Prefix Name was not defined use the computer name
			if vmssOsProfile.ComputerNamePrefix == nil || len(*vmssOsProfile.ComputerNamePrefix) == 0 {
				// validate that the computer name is a valid Computer Prefix Name
				_, errs := validate.LinuxComputerNamePrefix(id.VirtualMachineScaleSetName, "computer_name_prefix")
				if len(errs) > 0 {
					if errs[0] != nil {
						return fmt.Errorf("unable to assume default computer name prefix `%s`. Please adjust the `name`, or specify an explicit `computer_name_prefix`", errs[0])
					}

					return fmt.Errorf("unable to assume default computer name prefix. Please adjust the `name`, or specify an explicit `computer_name_prefix`")
				}

				vmssOsProfile.ComputerNamePrefix = pointer.To(id.VirtualMachineScaleSetName)
			}

			if extensionOperationsEnabled && !provisionVMAgent {
				return fmt.Errorf("`extension_operations_enabled` cannot be set to `true` when `provision_vm_agent` is set to `false`")
			}

			if patchAssessmentMode == string(virtualmachinescalesets.LinuxPatchAssessmentModeAutomaticByPlatform) && !provisionVMAgent {
				return fmt.Errorf("when `patch_assessment_mode` is set to `%s`, `provision_vm_agent` must be set to `true`", virtualmachinescalesets.LinuxPatchAssessmentModeAutomaticByPlatform)
			}

			// Validate Automatic VM Guest Patching Settings
			patchMode := linConfig["patch_mode"].(string)

			if patchMode == string(virtualmachinescalesets.LinuxVMGuestPatchModeAutomaticByPlatform) {
				if !provisionVMAgent {
					return fmt.Errorf("when `patch_mode` is set to `%s`, `provision_vm_agent` must be set to `true`, got `%s`", patchMode, strconv.FormatBool(provisionVMAgent))
				}

				if !hasHealthExtension {
					return fmt.Errorf("when `patch_mode` is set to `%s`, at least one application health extension must be configured, got 0", patchMode)
				}
			}
		}

		if vmssOsProfile != nil {
			vmssOsProfile.AllowExtensionOperations = pointer.To(extensionOperationsEnabled)
		}

		virtualMachineProfile.OsProfile = vmssOsProfile
	}

	if v, ok := d.GetOk("boot_diagnostics"); ok {
		virtualMachineProfile.DiagnosticsProfile = helpers.ExpandBootDiagnosticsVMSS(v.([]any))
	}

	if v, ok := d.GetOk("priority"); ok {
		virtualMachineProfile.Priority = pointer.ToEnum[virtualmachinescalesets.VirtualMachinePriorityTypes](v.(string))
	}

	if v, ok := d.GetOk("os_disk"); ok {
		virtualMachineProfile.StorageProfile.OsDisk = helpers.ExpandOrchestratedVirtualMachineScaleSetOSDisk(v.([]any), osType)
	}

	additionalCapabilitiesRaw := d.Get("additional_capabilities").([]any)
	props.Properties.AdditionalCapabilities = helpers.ExpandOrchestratedVirtualMachineScaleSetAdditionalCapabilities(additionalCapabilitiesRaw)

	if v, ok := d.GetOk("data_disk"); ok {
		ultraSSDEnabled := d.Get("additional_capabilities.0.ultra_ssd_enabled").(bool)
		dataDisks, err := helpers.ExpandOrchestratedVirtualMachineScaleSetDataDisk(v.([]any), ultraSSDEnabled)
		if err != nil {
			return fmt.Errorf("expanding `data_disk`: %w", err)
		}
		virtualMachineProfile.StorageProfile.DataDisks = dataDisks
	}

	if v, ok := d.GetOk("network_interface"); ok {
		networkInterfaces, err := helpers.ExpandOrchestratedVirtualMachineScaleSetNetworkInterface(v.([]any))
		if err != nil {
			return fmt.Errorf("expanding `network_interface`: %w", err)
		}

		networkProfile.NetworkInterfaceConfigurations = networkInterfaces
		virtualMachineProfile.NetworkProfile = networkProfile
	}

	if v, ok := d.Get("max_bid_price").(float64); ok && v > 0 {
		if *virtualMachineProfile.Priority != virtualmachinescalesets.VirtualMachinePriorityTypesSpot {
			return fmt.Errorf("`max_bid_price` can only be configured when `priority` is set to `%s`", string(virtualmachinescalesets.VirtualMachinePriorityTypesSpot))
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
			return fmt.Errorf("`eviction_policy` can only be specified when `priority` is set to `%s`", string(virtualmachinescalesets.VirtualMachinePriorityTypesSpot))
		}
		virtualMachineProfile.EvictionPolicy = pointer.ToEnum[virtualmachinescalesets.VirtualMachineEvictionPolicyTypes](v.(string))
	} else if *virtualMachineProfile.Priority == virtualmachinescalesets.VirtualMachinePriorityTypesSpot {
		return fmt.Errorf("`eviction_policy` is required when `priority` is set to `%s`", string(virtualmachinescalesets.VirtualMachinePriorityTypesSpot))
	}

	if v, ok := d.GetOk("license_type"); ok {
		virtualMachineProfile.LicenseType = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("termination_notification"); ok {
		virtualMachineProfile.ScheduledEventsProfile = helpers.ExpandOrchestratedVirtualMachineScaleSetScheduledEventsProfile(v.([]any))
	}

	// Only include the virtual machine profile if this is not a legacy configuration
	if !isLegacy {
		if v, ok := d.GetOk("plan"); ok {
			props.Plan = helpers.ExpandPlanVMSS(v.([]any))
		}

		if v, ok := d.GetOk("identity"); ok {
			identityExpanded, err := identity.ExpandSystemAndUserAssignedMap(v.([]any))
			if err != nil {
				return fmt.Errorf("expanding `identity`: %w", err)
			}
			props.Identity = identityExpanded
		}

		if v, ok := d.GetOk("automatic_instance_repair"); ok {
			if !hasHealthExtension {
				return fmt.Errorf("`automatic_instance_repair` can only be enabled when an application health extension is configured")
			}

			props.Properties.AutomaticRepairsPolicy = helpers.ExpandVirtualMachineScaleSetAutomaticRepairsPolicy(v.([]any))
		}

		if v, ok := d.GetOk("zone_balance"); ok && v.(bool) {
			if props.Zones == nil || len(*props.Zones) == 0 {
				return fmt.Errorf("`zone_balance` can only be set to `true` when availability zones are specified")
			}

			props.Properties.ZoneBalance = pointer.To(v.(bool))
		}

		if v, ok := d.GetOk("priority_mix"); ok {
			if *virtualMachineProfile.Priority != virtualmachinescalesets.VirtualMachinePriorityTypesSpot {
				return fmt.Errorf("`priority_mix` can only be specified when `priority` is set to `%s`", string(virtualmachinescalesets.VirtualMachinePriorityTypesSpot))
			}
			props.Properties.PriorityMixPolicy = helpers.ExpandOrchestratedVirtualMachineScaleSetPriorityMixPolicy(v.([]any))
		}

		props.Properties.VirtualMachineProfile = &virtualMachineProfile
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, props, virtualmachinescalesets.DefaultCreateOrUpdateOperationOptions(), sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating Orchestrated %s: %w", id, err)
	}

	d.SetId(id.ID())

	return resourceOrchestratedVirtualMachineScaleSetRead(d, meta)
}

func expandOrchestratedVirtualMachineScaleSetSkuProfile(input []any) *virtualmachinescalesets.SkuProfile {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]any)
	allocationStrategy := v["allocation_strategy"].(string)
	vmSizes := make([]virtualmachinescalesets.SkuProfileVMSize, 0)

	vmSizeRaw := v["virtual_machine_size"].(*pluginsdk.Set).List()
	for _, vmSizeRaw := range vmSizeRaw {
		vmSizeMap := vmSizeRaw.(map[string]any)
		vmSize := virtualmachinescalesets.SkuProfileVMSize{
			Name: pointer.To(vmSizeMap["name"].(string)),
		}

		// `rank` is optional, so an omitted value decodes to Go's zero value.
		// Only send it when the user explicitly configured it.
		if rank := vmSizeMap["rank"].(int); rank > 0 {
			vmSize.Rank = pointer.To(int64(rank - 1))
		}

		vmSizes = append(vmSizes, vmSize)
	}

	return &virtualmachinescalesets.SkuProfile{
		AllocationStrategy: pointer.ToEnum[virtualmachinescalesets.AllocationStrategy](allocationStrategy),
		VMSizes:            pointer.To(vmSizes),
	}
}

func expandOrchestratedVirtualMachineScaleSetSku(input string, capacity int) (*virtualmachinescalesets.Sku, error) {
	skuParts := strings.Split(input, "_")

	if (input != SkuNameMix && len(skuParts) < 2) || strings.Contains(input, "__") || strings.Contains(input, " ") {
		return nil, fmt.Errorf("`sku_name`(`%s`) is not formatted properly", input)
	}

	sku := &virtualmachinescalesets.Sku{
		Name:     pointer.To(input),
		Capacity: pointer.To(int64(capacity)),
	}

	if input != SkuNameMix {
		sku.Tier = pointer.To("Standard")
	}

	return sku, nil
}
