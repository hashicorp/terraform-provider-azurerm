// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package linux_virtual_machine_scale_set

import (
	"fmt"

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

func resourceLinuxVirtualMachineScaleSetCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Compute.VirtualMachineScaleSetsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := virtualmachinescalesets.NewVirtualMachineScaleSetID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		exists, err := client.Get(ctx, id, virtualmachinescalesets.DefaultGetOperationOptions())
		if err != nil {
			if !response.WasNotFound(exists.HttpResponse) {
				return fmt.Errorf("checking for existing Linux %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(exists.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_linux_virtual_machine_scale_set", id.ID())
		}
	}

	location := location.Normalize(d.Get("location").(string))

	additionalCapabilitiesRaw := d.Get("additional_capabilities").([]any)
	additionalCapabilities := helpers.ExpandVirtualMachineScaleSetAdditionalCapabilities(additionalCapabilitiesRaw)

	bootDiagnosticsRaw := d.Get("boot_diagnostics").([]any)
	bootDiagnostics := helpers.ExpandBootDiagnosticsVMSS(bootDiagnosticsRaw)

	dataDisksRaw := d.Get("data_disk").([]any)
	ultraSSDEnabled := d.Get("additional_capabilities.0.ultra_ssd_enabled").(bool)
	dataDisks, err := helpers.ExpandVirtualMachineScaleSetDataDisk(dataDisksRaw, ultraSSDEnabled)
	if err != nil {
		return fmt.Errorf("expanding `data_disk`: %+v", err)
	}

	identityExpanded, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]any))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	networkInterfacesRaw := d.Get("network_interface").([]any)
	networkInterfaces, err := helpers.ExpandVirtualMachineScaleSetNetworkInterface(networkInterfacesRaw)
	if err != nil {
		return fmt.Errorf("expanding `network_interface`: %+v", err)
	}

	osDiskRaw := d.Get("os_disk").([]any)
	osDisk, err := helpers.ExpandVirtualMachineScaleSetOSDisk(osDiskRaw, virtualmachinescalesets.OperatingSystemTypesLinux)
	if err != nil {
		return fmt.Errorf("expanding `os_disk`: %+v", err)
	}
	securityEncryptionType := osDiskRaw[0].(map[string]any)["security_encryption_type"].(string)

	planRaw := d.Get("plan").([]any)
	plan := helpers.ExpandPlanVMSS(planRaw)

	sourceImageReferenceRaw := d.Get("source_image_reference").([]any)
	sourceImageId := d.Get("source_image_id").(string)
	sourceImageReference := helpers.ExpandSourceImageReferenceVMSS(sourceImageReferenceRaw, sourceImageId)

	sshKeysRaw := d.Get("admin_ssh_key").(*pluginsdk.Set).List()
	sshKeys := helpers.ExpandSSHKeysVMSS(sshKeysRaw)

	overProvision := d.Get("overprovision").(bool)
	provisionVMAgent := d.Get("provision_vm_agent").(bool)
	zones := zones.ExpandUntyped(d.Get("zones").(*schema.Set).List())
	healthProbeId := d.Get("health_probe_id").(string)
	upgradeMode := virtualmachinescalesets.UpgradeMode(d.Get("upgrade_mode").(string))
	automaticOSUpgradePolicyRaw := d.Get("automatic_os_upgrade_policy").([]any)
	rollingUpgradePolicyRaw := d.Get("rolling_upgrade_policy").([]any)
	useRollingUpgradePolicy := len(rollingUpgradePolicyRaw) > 0
	automaticOSUpgradePolicy := helpers.ExpandVirtualMachineScaleSetAutomaticUpgradePolicy(automaticOSUpgradePolicyRaw, useRollingUpgradePolicy)
	rollingUpgradePolicy, err := helpers.ExpandVirtualMachineScaleSetRollingUpgradePolicy(rollingUpgradePolicyRaw, len(zones) > 0, overProvision)
	if err != nil {
		return err
	}

	canHaveAutomaticOsUpgradePolicy := upgradeMode == virtualmachinescalesets.UpgradeModeAutomatic || upgradeMode == virtualmachinescalesets.UpgradeModeRolling
	if !canHaveAutomaticOsUpgradePolicy && len(automaticOSUpgradePolicyRaw) > 0 {
		return fmt.Errorf("an `automatic_os_upgrade_policy` block cannot be specified when `upgrade_mode` is not set to `Automatic` or `Rolling`")
	}

	shouldHaveRollingUpgradePolicy := upgradeMode == virtualmachinescalesets.UpgradeModeAutomatic || upgradeMode == virtualmachinescalesets.UpgradeModeRolling
	if !shouldHaveRollingUpgradePolicy && len(rollingUpgradePolicyRaw) > 0 {
		return fmt.Errorf("a `rolling_upgrade_policy` block cannot be specified when `upgrade_mode` is set to %q", string(upgradeMode))
	}
	shouldHaveRollingUpgradePolicy = upgradeMode == virtualmachinescalesets.UpgradeModeRolling
	if shouldHaveRollingUpgradePolicy && len(rollingUpgradePolicyRaw) == 0 {
		return fmt.Errorf("a `rolling_upgrade_policy` block must be specified when `upgrade_mode` is set to %q", string(upgradeMode))
	}

	secretsRaw := d.Get("secret").([]any)
	secrets := helpers.ExpandLinuxSecretsVMSS(secretsRaw)

	var computerNamePrefix string
	if v, ok := d.GetOk("computer_name_prefix"); ok && len(v.(string)) > 0 {
		computerNamePrefix = v.(string)
	} else {
		_, errs := validate.LinuxComputerNamePrefix(d.Get("name"), "computer_name_prefix")
		if len(errs) > 0 {
			return fmt.Errorf("unable to assume default computer name prefix %s. Please adjust the %q, or specify an explicit %q", errs[0], "name", "computer_name_prefix")
		}
		computerNamePrefix = id.VirtualMachineScaleSetName
	}

	disablePasswordAuthentication := d.Get("disable_password_authentication").(bool)
	networkProfile := &virtualmachinescalesets.VirtualMachineScaleSetNetworkProfile{
		NetworkInterfaceConfigurations: networkInterfaces,
	}
	if healthProbeId != "" {
		networkProfile.HealthProbe = &virtualmachinescalesets.ApiEntityReference{
			Id: pointer.To(healthProbeId),
		}
	}

	priority := virtualmachinescalesets.VirtualMachinePriorityTypes(d.Get("priority").(string))
	upgradePolicy := virtualmachinescalesets.UpgradePolicy{
		Mode:                     pointer.To(upgradeMode),
		AutomaticOSUpgradePolicy: automaticOSUpgradePolicy,
		RollingUpgradePolicy:     rollingUpgradePolicy,
	}

	virtualMachineProfile := virtualmachinescalesets.VirtualMachineScaleSetVMProfile{
		Priority: pointer.To(priority),
		OsProfile: &virtualmachinescalesets.VirtualMachineScaleSetOSProfile{
			AdminUsername:      pointer.To(d.Get("admin_username").(string)),
			ComputerNamePrefix: pointer.To(computerNamePrefix),
			LinuxConfiguration: &virtualmachinescalesets.LinuxConfiguration{
				DisablePasswordAuthentication: pointer.To(disablePasswordAuthentication),
				ProvisionVMAgent:              pointer.To(provisionVMAgent),
				Ssh: &virtualmachinescalesets.SshConfiguration{
					PublicKeys: &sshKeys,
				},
			},
			Secrets: secrets,
		},
		DiagnosticsProfile: bootDiagnostics,
		NetworkProfile:     networkProfile,
		StorageProfile: &virtualmachinescalesets.VirtualMachineScaleSetStorageProfile{
			ImageReference: sourceImageReference,
			OsDisk:         osDisk,
			DataDisks:      dataDisks,
		},
	}

	if galleryApplications := helpers.ExpandVirtualMachineScaleSetGalleryApplication(d.Get("gallery_application").([]any)); galleryApplications != nil {
		virtualMachineProfile.ApplicationProfile = &virtualmachinescalesets.ApplicationProfile{
			GalleryApplications: galleryApplications,
		}
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

	hasHealthExtension := false
	if vmExtensionsRaw, ok := d.GetOk("extension"); ok {
		virtualMachineProfile.ExtensionProfile, hasHealthExtension, err = helpers.ExpandVirtualMachineScaleSetExtensions(vmExtensionsRaw.(*pluginsdk.Set).List())
		if err != nil {
			return err
		}
	}

	if v, ok := d.Get("extension_operations_enabled").(bool); ok {
		if v && !provisionVMAgent {
			return fmt.Errorf("`extension_operations_enabled` cannot be set to `true` when `provision_vm_agent` is set to `false`")
		}

		virtualMachineProfile.OsProfile.AllowExtensionOperations = pointer.To(v)
	}

	if v, ok := d.GetOk("extensions_time_budget"); ok {
		if virtualMachineProfile.ExtensionProfile == nil {
			virtualMachineProfile.ExtensionProfile = &virtualmachinescalesets.VirtualMachineScaleSetExtensionProfile{}
		}
		virtualMachineProfile.ExtensionProfile.ExtensionsTimeBudget = pointer.To(v.(string))
	}

	// otherwise the service return the error:
	// Rolling Upgrade mode is not supported for this Virtual Machine Scale Set because a health probe or health extension was not provided.
	if upgradeMode == virtualmachinescalesets.UpgradeModeRolling && (healthProbeId == "" && !hasHealthExtension) {
		return fmt.Errorf("`health_probe_id` must be set or a health extension must be specified when `upgrade_mode` is set to %q", string(upgradeMode))
	}

	if adminPassword, ok := d.GetOk("admin_password"); ok {
		virtualMachineProfile.OsProfile.AdminPassword = pointer.To(adminPassword.(string))
	}

	if v, ok := d.Get("max_bid_price").(float64); ok && v > 0 {
		if priority != virtualmachinescalesets.VirtualMachinePriorityTypesSpot {
			return fmt.Errorf("`max_bid_price` can only be configured when `priority` is set to `Spot`")
		}

		virtualMachineProfile.BillingProfile = &virtualmachinescalesets.BillingProfile{
			MaxPrice: pointer.To(v),
		}
	}

	if v, ok := d.GetOk("custom_data"); ok {
		virtualMachineProfile.OsProfile.CustomData = pointer.To(v.(string))
	}

	if encryptionAtHostEnabled, ok := d.GetOk("encryption_at_host_enabled"); ok {
		virtualMachineProfile.SecurityProfile = &virtualmachinescalesets.SecurityProfile{
			EncryptionAtHost: pointer.To(encryptionAtHostEnabled.(bool)),
		}
	}

	secureBootEnabled := d.Get("secure_boot_enabled").(bool)
	vtpmEnabled := d.Get("vtpm_enabled").(bool)
	if securityEncryptionType != "" {
		if virtualmachinescalesets.SecurityEncryptionTypesDiskWithVMGuestState == virtualmachinescalesets.SecurityEncryptionTypes(securityEncryptionType) && !secureBootEnabled {
			return fmt.Errorf("`secure_boot_enabled` must be set to `true` when `os_disk.0.security_encryption_type` is set to `DiskWithVMGuestState`")
		}
		if !vtpmEnabled {
			return fmt.Errorf("`vtpm_enabled` must be set to `true` when `os_disk.0.security_encryption_type` is set")
		}

		if virtualMachineProfile.SecurityProfile == nil {
			virtualMachineProfile.SecurityProfile = &virtualmachinescalesets.SecurityProfile{}
		}
		virtualMachineProfile.SecurityProfile.SecurityType = pointer.To(virtualmachinescalesets.SecurityTypesConfidentialVM)

		if virtualMachineProfile.SecurityProfile.UefiSettings == nil {
			virtualMachineProfile.SecurityProfile.UefiSettings = &virtualmachinescalesets.UefiSettings{}
		}
		virtualMachineProfile.SecurityProfile.UefiSettings.SecureBootEnabled = pointer.To(secureBootEnabled)
		virtualMachineProfile.SecurityProfile.UefiSettings.VTpmEnabled = pointer.To(vtpmEnabled)
	} else {
		if secureBootEnabled {
			if virtualMachineProfile.SecurityProfile == nil {
				virtualMachineProfile.SecurityProfile = &virtualmachinescalesets.SecurityProfile{}
			}
			if virtualMachineProfile.SecurityProfile.UefiSettings == nil {
				virtualMachineProfile.SecurityProfile.UefiSettings = &virtualmachinescalesets.UefiSettings{}
			}
			virtualMachineProfile.SecurityProfile.SecurityType = pointer.To(virtualmachinescalesets.SecurityTypesTrustedLaunch)
			virtualMachineProfile.SecurityProfile.UefiSettings.SecureBootEnabled = pointer.To(secureBootEnabled)
		}

		if vtpmEnabled {
			if virtualMachineProfile.SecurityProfile == nil {
				virtualMachineProfile.SecurityProfile = &virtualmachinescalesets.SecurityProfile{}
			}
			if virtualMachineProfile.SecurityProfile.UefiSettings == nil {
				virtualMachineProfile.SecurityProfile.UefiSettings = &virtualmachinescalesets.UefiSettings{}
			}
			virtualMachineProfile.SecurityProfile.SecurityType = pointer.To(virtualmachinescalesets.SecurityTypesTrustedLaunch)
			virtualMachineProfile.SecurityProfile.UefiSettings.VTpmEnabled = pointer.To(vtpmEnabled)
		}
	}

	if v, ok := d.GetOk("user_data"); ok {
		virtualMachineProfile.UserData = pointer.To(v.(string))
	}

	// Azure API: "Authentication using either SSH or by user name and password must be enabled in Linux profile."
	if disablePasswordAuthentication && virtualMachineProfile.OsProfile.AdminPassword == nil && len(sshKeys) == 0 {
		return fmt.Errorf("at least one SSH key must be specified if `disable_password_authentication` is enabled")
	}

	if evictionPolicyRaw, ok := d.GetOk("eviction_policy"); ok {
		if *virtualMachineProfile.Priority != virtualmachinescalesets.VirtualMachinePriorityTypesSpot {
			return fmt.Errorf("an `eviction_policy` can only be specified when `priority` is set to `Spot`")
		}
		virtualMachineProfile.EvictionPolicy = pointer.ToEnum[virtualmachinescalesets.VirtualMachineEvictionPolicyTypes](evictionPolicyRaw.(string))
	} else if priority == virtualmachinescalesets.VirtualMachinePriorityTypesSpot {
		return fmt.Errorf("an `eviction_policy` must be specified when `priority` is set to `Spot`")
	}

	if v, ok := d.GetOk("termination_notification"); ok {
		virtualMachineProfile.ScheduledEventsProfile = helpers.ExpandVirtualMachineScaleSetScheduledEventsProfile(v.([]any))
	}

	automaticRepairsPolicyRaw := d.Get("automatic_instance_repair").([]any)
	automaticRepairsPolicy := helpers.ExpandVirtualMachineScaleSetAutomaticRepairsPolicy(automaticRepairsPolicyRaw)

	if automaticRepairsPolicy != nil && healthProbeId == "" && !hasHealthExtension {
		return fmt.Errorf("`automatic_instance_repair` can only be set if there is an application Health extension or a `health_probe_id` defined")
	}

	props := virtualmachinescalesets.VirtualMachineScaleSet{
		ExtendedLocation: helpers.ExpandEdgeZone(d.Get("edge_zone").(string)),
		Location:         location,
		Sku: &virtualmachinescalesets.Sku{
			Name:     pointer.To(d.Get("sku").(string)),
			Capacity: pointer.To(int64(d.Get("instances").(int))),

			// doesn't appear this can be set to anything else, even Promo machines are Standard
			Tier: pointer.To("Standard"),
		},
		Identity: identityExpanded,
		Plan:     plan,
		Tags:     tags.Expand(d.Get("tags").(map[string]any)),
		Properties: &virtualmachinescalesets.VirtualMachineScaleSetProperties{
			AdditionalCapabilities:                 additionalCapabilities,
			AutomaticRepairsPolicy:                 automaticRepairsPolicy,
			DoNotRunExtensionsOnOverprovisionedVMs: pointer.To(d.Get("do_not_run_extensions_on_overprovisioned_machines").(bool)),
			Overprovision:                          pointer.To(d.Get("overprovision").(bool)),
			SinglePlacementGroup:                   pointer.To(d.Get("single_placement_group").(bool)),
			VirtualMachineProfile:                  &virtualMachineProfile,
			UpgradePolicy:                          &upgradePolicy,
			// OrchestrationMode needs to be hardcoded to Uniform, for the
			// standard VMSS resource, since virtualMachineProfile is now supported
			// in both VMSS and Orchestrated VMSS...
			OrchestrationMode: pointer.To(virtualmachinescalesets.OrchestrationModeUniform),
		},
	}

	if v, ok := d.GetOk("scale_in"); ok {
		if v := helpers.ExpandVirtualMachineScaleSetScaleInPolicy(v.([]any)); v != nil {
			props.Properties.ScaleInPolicy = v
		}
	}

	if v, ok := d.GetOk("host_group_id"); ok {
		props.Properties.HostGroup = &virtualmachinescalesets.SubResource{
			Id: pointer.To(v.(string)),
		}
	}

	if spotRestorePolicy := helpers.ExpandVirtualMachineScaleSetSpotRestorePolicy(d.Get("spot_restore").([]any)); spotRestorePolicy != nil {
		props.Properties.SpotRestorePolicy = spotRestorePolicy
	}

	props.Properties.ResiliencyPolicy = helpers.ExpandVirtualMachineScaleSetResiliency(d.Get("resilient_vm_creation_enabled").(bool), d.Get("resilient_vm_deletion_enabled").(bool))

	if len(zones) > 0 {
		props.Zones = &zones
	}

	if v, ok := d.GetOk("platform_fault_domain_count"); ok {
		props.Properties.PlatformFaultDomainCount = pointer.To(int64(v.(int)))
	}

	if v, ok := d.GetOk("proximity_placement_group_id"); ok {
		props.Properties.ProximityPlacementGroup = &virtualmachinescalesets.SubResource{
			Id: pointer.To(v.(string)),
		}
	}

	if v, ok := d.GetOk("zone_balance"); ok && v.(bool) {
		if props.Zones == nil || len(*props.Zones) == 0 {
			return fmt.Errorf("`zone_balance` can only be set to `true` when zones are specified")
		}

		props.Properties.ZoneBalance = pointer.To(v.(bool))
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, props, virtualmachinescalesets.DefaultCreateOrUpdateOperationOptions(), sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating Linux %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceLinuxVirtualMachineScaleSetRead(d, meta)
}
