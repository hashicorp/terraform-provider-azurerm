// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"fmt"
	"log"
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
	azValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/base64"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/compute/2023-03-01/compute"
)

func resourceLinuxVirtualMachineScaleSet() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLinuxVirtualMachineScaleSetCreate,
		Read:   resourceLinuxVirtualMachineScaleSetRead,
		Update: resourceLinuxVirtualMachineScaleSetUpdate,
		Delete: resourceLinuxVirtualMachineScaleSetDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := parse.VirtualMachineScaleSetID(id)
			return err
		}, importVirtualMachineScaleSet(compute.OperatingSystemTypesLinux, "azurerm_linux_virtual_machine_scale_set")),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(time.Minute * 60),
			Read:   pluginsdk.DefaultTimeout(time.Minute * 5),
			Update: pluginsdk.DefaultTimeout(time.Minute * 60),
			Delete: pluginsdk.DefaultTimeout(time.Minute * 60),
		},

		// TODO: exposing requireGuestProvisionSignal once it's available
		// https://github.com/Azure/azure-rest-api-specs/pull/7246

		Schema: resourceLinuxVirtualMachineScaleSetSchema(),
	}
}

func resourceLinuxVirtualMachineScaleSetCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VMScaleSetClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewVirtualMachineScaleSetID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	// Upgrading to the 2021-07-01 exposed a new expand parameter to the GET method
	exists, err := client.Get(ctx, id.ResourceGroup, id.Name, compute.ExpandTypesForGetVMScaleSetsUserData)
	if err != nil {
		if !utils.ResponseWasNotFound(exists.Response) {
			return fmt.Errorf("checking for existing Linux %s: %+v", id, err)
		}
	}

	if !utils.ResponseWasNotFound(exists.Response) {
		return tf.ImportAsExistsError("azurerm_linux_virtual_machine_scale_set", *exists.ID)
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	additionalCapabilitiesRaw := d.Get("additional_capabilities").([]interface{})
	additionalCapabilities := ExpandVirtualMachineScaleSetAdditionalCapabilities(additionalCapabilitiesRaw)

	bootDiagnosticsRaw := d.Get("boot_diagnostics").([]interface{})
	bootDiagnostics := expandBootDiagnostics(bootDiagnosticsRaw)

	dataDisksRaw := d.Get("data_disk").([]interface{})
	ultraSSDEnabled := d.Get("additional_capabilities.0.ultra_ssd_enabled").(bool)
	dataDisks, err := ExpandVirtualMachineScaleSetDataDisk(dataDisksRaw, ultraSSDEnabled)
	if err != nil {
		return fmt.Errorf("expanding `data_disk`: %+v", err)
	}

	identity, err := expandVirtualMachineScaleSetIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	networkInterfacesRaw := d.Get("network_interface").([]interface{})
	networkInterfaces, err := ExpandVirtualMachineScaleSetNetworkInterface(networkInterfacesRaw)
	if err != nil {
		return fmt.Errorf("expanding `network_interface`: %+v", err)
	}

	osDiskRaw := d.Get("os_disk").([]interface{})
	osDisk, err := ExpandVirtualMachineScaleSetOSDisk(osDiskRaw, compute.OperatingSystemTypesLinux)
	if err != nil {
		return fmt.Errorf("expanding `os_disk`: %+v", err)
	}
	securityEncryptionType := osDiskRaw[0].(map[string]interface{})["security_encryption_type"].(string)

	planRaw := d.Get("plan").([]interface{})
	plan := expandPlan(planRaw)

	sourceImageReferenceRaw := d.Get("source_image_reference").([]interface{})
	sourceImageId := d.Get("source_image_id").(string)
	sourceImageReference := expandSourceImageReference(sourceImageReferenceRaw, sourceImageId)

	sshKeysRaw := d.Get("admin_ssh_key").(*pluginsdk.Set).List()
	sshKeys := ExpandSSHKeys(sshKeysRaw)

	provisionVMAgent := d.Get("provision_vm_agent").(bool)
	zones := zones.ExpandUntyped(d.Get("zones").(*schema.Set).List())
	healthProbeId := d.Get("health_probe_id").(string)
	upgradeMode := compute.UpgradeMode(d.Get("upgrade_mode").(string))
	automaticOSUpgradePolicyRaw := d.Get("automatic_os_upgrade_policy").([]interface{})
	automaticOSUpgradePolicy := ExpandVirtualMachineScaleSetAutomaticUpgradePolicy(automaticOSUpgradePolicyRaw)
	rollingUpgradePolicyRaw := d.Get("rolling_upgrade_policy").([]interface{})
	rollingUpgradePolicy, err := ExpandVirtualMachineScaleSetRollingUpgradePolicy(rollingUpgradePolicyRaw, len(zones) > 0)
	if err != nil {
		return err
	}

	canHaveAutomaticOsUpgradePolicy := upgradeMode == compute.UpgradeModeAutomatic || upgradeMode == compute.UpgradeModeRolling
	if !canHaveAutomaticOsUpgradePolicy && len(automaticOSUpgradePolicyRaw) > 0 {
		return fmt.Errorf("an `automatic_os_upgrade_policy` block cannot be specified when `upgrade_mode` is not set to `Automatic` or `Rolling`")
	}

	shouldHaveRollingUpgradePolicy := upgradeMode == compute.UpgradeModeAutomatic || upgradeMode == compute.UpgradeModeRolling
	if !shouldHaveRollingUpgradePolicy && len(rollingUpgradePolicyRaw) > 0 {
		return fmt.Errorf("a `rolling_upgrade_policy` block cannot be specified when `upgrade_mode` is set to %q", string(upgradeMode))
	}
	shouldHaveRollingUpgradePolicy = upgradeMode == compute.UpgradeModeRolling
	if shouldHaveRollingUpgradePolicy && len(rollingUpgradePolicyRaw) == 0 {
		return fmt.Errorf("a `rolling_upgrade_policy` block must be specified when `upgrade_mode` is set to %q", string(upgradeMode))
	}

	secretsRaw := d.Get("secret").([]interface{})
	secrets := expandLinuxSecrets(secretsRaw)

	var computerNamePrefix string
	if v, ok := d.GetOk("computer_name_prefix"); ok && len(v.(string)) > 0 {
		computerNamePrefix = v.(string)
	} else {
		_, errs := validate.LinuxComputerNamePrefix(d.Get("name"), "computer_name_prefix")
		if len(errs) > 0 {
			return fmt.Errorf("unable to assume default computer name prefix %s. Please adjust the %q, or specify an explicit %q", errs[0], "name", "computer_name_prefix")
		}
		computerNamePrefix = id.Name
	}

	disablePasswordAuthentication := d.Get("disable_password_authentication").(bool)
	networkProfile := &compute.VirtualMachineScaleSetNetworkProfile{
		NetworkInterfaceConfigurations: networkInterfaces,
	}
	if healthProbeId != "" {
		networkProfile.HealthProbe = &compute.APIEntityReference{
			ID: utils.String(healthProbeId),
		}
	}

	priority := compute.VirtualMachinePriorityTypes(d.Get("priority").(string))
	upgradePolicy := compute.UpgradePolicy{
		Mode:                     upgradeMode,
		AutomaticOSUpgradePolicy: automaticOSUpgradePolicy,
		RollingUpgradePolicy:     rollingUpgradePolicy,
	}

	virtualMachineProfile := compute.VirtualMachineScaleSetVMProfile{
		Priority: priority,
		OsProfile: &compute.VirtualMachineScaleSetOSProfile{
			AdminUsername:      utils.String(d.Get("admin_username").(string)),
			ComputerNamePrefix: utils.String(computerNamePrefix),
			LinuxConfiguration: &compute.LinuxConfiguration{
				DisablePasswordAuthentication: utils.Bool(disablePasswordAuthentication),
				ProvisionVMAgent:              utils.Bool(provisionVMAgent),
				SSH: &compute.SSHConfiguration{
					PublicKeys: &sshKeys,
				},
			},
			Secrets: secrets,
		},
		DiagnosticsProfile: bootDiagnostics,
		NetworkProfile:     networkProfile,
		StorageProfile: &compute.VirtualMachineScaleSetStorageProfile{
			ImageReference: sourceImageReference,
			OsDisk:         osDisk,
			DataDisks:      dataDisks,
		},
	}

	if !features.FourPointOhBeta() {
		if galleryApplications := expandVirtualMachineScaleSetGalleryApplications(d.Get("gallery_applications").([]interface{})); galleryApplications != nil {
			virtualMachineProfile.ApplicationProfile = &compute.ApplicationProfile{
				GalleryApplications: galleryApplications,
			}
		}
	}

	if galleryApplications := expandVirtualMachineScaleSetGalleryApplication(d.Get("gallery_application").([]interface{})); galleryApplications != nil {
		virtualMachineProfile.ApplicationProfile = &compute.ApplicationProfile{
			GalleryApplications: galleryApplications,
		}
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

	hasHealthExtension := false
	if vmExtensionsRaw, ok := d.GetOk("extension"); ok {
		virtualMachineProfile.ExtensionProfile, hasHealthExtension, err = expandVirtualMachineScaleSetExtensions(vmExtensionsRaw.(*pluginsdk.Set).List())
		if err != nil {
			return err
		}
	}

	if v, ok := d.Get("extension_operations_enabled").(bool); ok {
		if v && !provisionVMAgent {
			return fmt.Errorf("`extension_operations_enabled` cannot be set to `true` when `provision_vm_agent` is set to `false`")
		}

		if !features.FourPointOhBeta() {
			if !pluginsdk.IsExplicitlyNullInConfig(d, "extension_operations_enabled") {
				virtualMachineProfile.OsProfile.AllowExtensionOperations = utils.Bool(v)
			}
		} else {
			virtualMachineProfile.OsProfile.AllowExtensionOperations = utils.Bool(v)
		}
	}

	if v, ok := d.GetOk("extensions_time_budget"); ok {
		if virtualMachineProfile.ExtensionProfile == nil {
			virtualMachineProfile.ExtensionProfile = &compute.VirtualMachineScaleSetExtensionProfile{}
		}
		virtualMachineProfile.ExtensionProfile.ExtensionsTimeBudget = utils.String(v.(string))
	}

	// otherwise the service return the error:
	// Rolling Upgrade mode is not supported for this Virtual Machine Scale Set because a health probe or health extension was not provided.
	if upgradeMode == compute.UpgradeModeRolling && (healthProbeId == "" && !hasHealthExtension) {
		return fmt.Errorf("`health_probe_id` must be set or a health extension must be specified when `upgrade_mode` is set to %q", string(upgradeMode))
	}

	if adminPassword, ok := d.GetOk("admin_password"); ok {
		virtualMachineProfile.OsProfile.AdminPassword = utils.String(adminPassword.(string))
	}

	if v, ok := d.Get("max_bid_price").(float64); ok && v > 0 {
		if priority != compute.VirtualMachinePriorityTypesSpot {
			return fmt.Errorf("`max_bid_price` can only be configured when `priority` is set to `Spot`")
		}

		virtualMachineProfile.BillingProfile = &compute.BillingProfile{
			MaxPrice: utils.Float(v),
		}
	}

	if v, ok := d.GetOk("custom_data"); ok {
		virtualMachineProfile.OsProfile.CustomData = utils.String(v.(string))
	}

	if encryptionAtHostEnabled, ok := d.GetOk("encryption_at_host_enabled"); ok {
		if encryptionAtHostEnabled.(bool) {
			if compute.SecurityEncryptionTypesDiskWithVMGuestState == compute.SecurityEncryptionTypes(securityEncryptionType) {
				return fmt.Errorf("`encryption_at_host_enabled` cannot be set to `true` when `os_disk.0.security_encryption_type` is set to `DiskWithVMGuestState`")
			}
		}

		virtualMachineProfile.SecurityProfile = &compute.SecurityProfile{
			EncryptionAtHost: utils.Bool(encryptionAtHostEnabled.(bool)),
		}
	}

	secureBootEnabled := d.Get("secure_boot_enabled").(bool)
	vtpmEnabled := d.Get("vtpm_enabled").(bool)
	if securityEncryptionType != "" {
		if compute.SecurityEncryptionTypesDiskWithVMGuestState == compute.SecurityEncryptionTypes(securityEncryptionType) && !secureBootEnabled {
			return fmt.Errorf("`secure_boot_enabled` must be set to `true` when `os_disk.0.security_encryption_type` is set to `DiskWithVMGuestState`")
		}
		if !vtpmEnabled {
			return fmt.Errorf("`vtpm_enabled` must be set to `true` when `os_disk.0.security_encryption_type` is set")
		}

		if virtualMachineProfile.SecurityProfile == nil {
			virtualMachineProfile.SecurityProfile = &compute.SecurityProfile{}
		}
		virtualMachineProfile.SecurityProfile.SecurityType = compute.SecurityTypesConfidentialVM

		if virtualMachineProfile.SecurityProfile.UefiSettings == nil {
			virtualMachineProfile.SecurityProfile.UefiSettings = &compute.UefiSettings{}
		}
		virtualMachineProfile.SecurityProfile.UefiSettings.SecureBootEnabled = utils.Bool(secureBootEnabled)
		virtualMachineProfile.SecurityProfile.UefiSettings.VTpmEnabled = utils.Bool(vtpmEnabled)
	} else {
		if secureBootEnabled {
			if virtualMachineProfile.SecurityProfile == nil {
				virtualMachineProfile.SecurityProfile = &compute.SecurityProfile{}
			}
			if virtualMachineProfile.SecurityProfile.UefiSettings == nil {
				virtualMachineProfile.SecurityProfile.UefiSettings = &compute.UefiSettings{}
			}
			virtualMachineProfile.SecurityProfile.SecurityType = compute.SecurityTypesTrustedLaunch
			virtualMachineProfile.SecurityProfile.UefiSettings.SecureBootEnabled = utils.Bool(secureBootEnabled)
		}

		if vtpmEnabled {
			if virtualMachineProfile.SecurityProfile == nil {
				virtualMachineProfile.SecurityProfile = &compute.SecurityProfile{}
			}
			if virtualMachineProfile.SecurityProfile.UefiSettings == nil {
				virtualMachineProfile.SecurityProfile.UefiSettings = &compute.UefiSettings{}
			}
			virtualMachineProfile.SecurityProfile.SecurityType = compute.SecurityTypesTrustedLaunch
			virtualMachineProfile.SecurityProfile.UefiSettings.VTpmEnabled = utils.Bool(vtpmEnabled)
		}
	}

	if v, ok := d.GetOk("user_data"); ok {
		virtualMachineProfile.UserData = utils.String(v.(string))
	}

	// Azure API: "Authentication using either SSH or by user name and password must be enabled in Linux profile."
	if disablePasswordAuthentication && virtualMachineProfile.OsProfile.AdminPassword == nil && len(sshKeys) == 0 {
		return fmt.Errorf("at least one SSH key must be specified if `disable_password_authentication` is enabled")
	}

	if evictionPolicyRaw, ok := d.GetOk("eviction_policy"); ok {
		if virtualMachineProfile.Priority != compute.VirtualMachinePriorityTypesSpot {
			return fmt.Errorf("an `eviction_policy` can only be specified when `priority` is set to `Spot`")
		}
		virtualMachineProfile.EvictionPolicy = compute.VirtualMachineEvictionPolicyTypes(evictionPolicyRaw.(string))
	} else if priority == compute.VirtualMachinePriorityTypesSpot {
		return fmt.Errorf("an `eviction_policy` must be specified when `priority` is set to `Spot`")
	}

	scaleInPolicy := &compute.ScaleInPolicy{
		Rules:         &[]compute.VirtualMachineScaleSetScaleInRules{compute.VirtualMachineScaleSetScaleInRules(string(compute.VirtualMachineScaleSetScaleInRulesDefault))},
		ForceDeletion: utils.Bool(false),
	}

	if !features.FourPointOhBeta() {
		if v, ok := d.GetOk("scale_in_policy"); ok {
			scaleInPolicy.Rules = &[]compute.VirtualMachineScaleSetScaleInRules{compute.VirtualMachineScaleSetScaleInRules(v.(string))}
		}

		if v, ok := d.GetOk("terminate_notification"); ok {
			virtualMachineProfile.ScheduledEventsProfile = ExpandVirtualMachineScaleSetScheduledEventsProfile(v.([]interface{}))
		}
	}

	if v, ok := d.GetOk("scale_in"); ok {
		if v := ExpandVirtualMachineScaleSetScaleInPolicy(v.([]interface{})); v != nil {
			scaleInPolicy = v
		}
	}

	if v, ok := d.GetOk("termination_notification"); ok {
		virtualMachineProfile.ScheduledEventsProfile = ExpandVirtualMachineScaleSetScheduledEventsProfile(v.([]interface{}))
	}

	automaticRepairsPolicyRaw := d.Get("automatic_instance_repair").([]interface{})
	automaticRepairsPolicy := ExpandVirtualMachineScaleSetAutomaticRepairsPolicy(automaticRepairsPolicyRaw)

	props := compute.VirtualMachineScaleSet{
		ExtendedLocation: expandEdgeZone(d.Get("edge_zone").(string)),
		Location:         utils.String(location),
		Sku: &compute.Sku{
			Name:     utils.String(d.Get("sku").(string)),
			Capacity: utils.Int64(int64(d.Get("instances").(int))),

			// doesn't appear this can be set to anything else, even Promo machines are Standard
			Tier: utils.String("Standard"),
		},
		Identity: identity,
		Plan:     plan,
		Tags:     tags.Expand(t),
		VirtualMachineScaleSetProperties: &compute.VirtualMachineScaleSetProperties{
			AdditionalCapabilities:                 additionalCapabilities,
			AutomaticRepairsPolicy:                 automaticRepairsPolicy,
			DoNotRunExtensionsOnOverprovisionedVMs: utils.Bool(d.Get("do_not_run_extensions_on_overprovisioned_machines").(bool)),
			Overprovision:                          utils.Bool(d.Get("overprovision").(bool)),
			SinglePlacementGroup:                   utils.Bool(d.Get("single_placement_group").(bool)),
			VirtualMachineProfile:                  &virtualMachineProfile,
			UpgradePolicy:                          &upgradePolicy,
			// OrchestrationMode needs to be hardcoded to Uniform, for the
			// standard VMSS resource, since virtualMachineProfile is now supported
			// in both VMSS and Orchestrated VMSS...
			OrchestrationMode: compute.OrchestrationModeUniform,
			ScaleInPolicy:     scaleInPolicy,
		},
	}

	if v, ok := d.GetOk("host_group_id"); ok {
		props.VirtualMachineScaleSetProperties.HostGroup = &compute.SubResource{
			ID: utils.String(v.(string)),
		}
	}

	spotRestoreRaw := d.Get("spot_restore").([]interface{})
	if spotRestorePolicy := ExpandVirtualMachineScaleSetSpotRestorePolicy(spotRestoreRaw); spotRestorePolicy != nil {
		props.SpotRestorePolicy = spotRestorePolicy
	}

	if len(zones) > 0 {
		props.Zones = &zones
	}

	if v, ok := d.GetOk("platform_fault_domain_count"); ok {
		props.VirtualMachineScaleSetProperties.PlatformFaultDomainCount = utils.Int32(int32(v.(int)))
	}

	if v, ok := d.GetOk("proximity_placement_group_id"); ok {
		props.VirtualMachineScaleSetProperties.ProximityPlacementGroup = &compute.SubResource{
			ID: utils.String(v.(string)),
		}
	}

	if v, ok := d.GetOk("zone_balance"); ok && v.(bool) {
		if props.Zones == nil || len(*props.Zones) == 0 {
			return fmt.Errorf("`zone_balance` can only be set to `true` when zones are specified")
		}

		props.VirtualMachineScaleSetProperties.ZoneBalance = utils.Bool(v.(bool))
	}

	log.Printf("[DEBUG] Creating Linux %s", id)
	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, props)
	if err != nil {
		return fmt.Errorf("creating Linux %s: %+v", id, err)
	}

	log.Printf("[DEBUG] Waiting for Linux %s to be created..", id)
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of Linux %s: %+v", id, err)
	}
	log.Printf("[DEBUG] %s was created", id)

	d.SetId(id.ID())

	return resourceLinuxVirtualMachineScaleSetRead(d, meta)
}

func resourceLinuxVirtualMachineScaleSetUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VMScaleSetClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualMachineScaleSetID(d.Id())
	if err != nil {
		return err
	}

	updateInstances := false

	// retrieve
	// Upgrading to the 2021-07-01 exposed a new expand parameter to the GET method
	existing, err := client.Get(ctx, id.ResourceGroup, id.Name, compute.ExpandTypesForGetVMScaleSetsUserData)
	if err != nil {
		return fmt.Errorf("retrieving Linux Virtual Machine Scale Set %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	if existing.VirtualMachineScaleSetProperties == nil {
		return fmt.Errorf("retrieving Linux Virtual Machine Scale Set %q (Resource Group %q): `properties` was nil", id.Name, id.ResourceGroup)
	}
	if existing.VirtualMachineScaleSetProperties.VirtualMachineProfile == nil {
		return fmt.Errorf("retrieving Linux Virtual Machine Scale Set %q (Resource Group %q): `properties.virtualMachineProfile` was nil", id.Name, id.ResourceGroup)
	}
	if existing.VirtualMachineScaleSetProperties.VirtualMachineProfile.StorageProfile == nil {
		return fmt.Errorf("retrieving Linux Virtual Machine Scale Set %q (Resource Group %q): `properties.virtualMachineProfile,storageProfile` was nil", id.Name, id.ResourceGroup)
	}

	updateProps := compute.VirtualMachineScaleSetUpdateProperties{
		VirtualMachineProfile: &compute.VirtualMachineScaleSetUpdateVMProfile{
			// if an image reference has been configured previously (it has to be), we would better to include that in this
			// update request to avoid some circumstances that the API will complain ImageReference is null
			// issue tracking: https://github.com/Azure/azure-rest-api-specs/issues/10322
			StorageProfile: &compute.VirtualMachineScaleSetUpdateStorageProfile{
				ImageReference: existing.VirtualMachineScaleSetProperties.VirtualMachineProfile.StorageProfile.ImageReference,
			},
		},
		// if an upgrade policy's been configured previously (which it will have) it must be threaded through
		// this doesn't matter for Manual - but breaks when updating anything on a Automatic and Rolling Mode Scale Set
		UpgradePolicy: existing.VirtualMachineScaleSetProperties.UpgradePolicy,
	}
	update := compute.VirtualMachineScaleSetUpdate{}

	// first try and pull this from existing vm, which covers no changes being made to this block
	automaticOSUpgradeIsEnabled := false
	if policy := existing.VirtualMachineScaleSetProperties.UpgradePolicy; policy != nil {
		if policy.AutomaticOSUpgradePolicy != nil && policy.AutomaticOSUpgradePolicy.EnableAutomaticOSUpgrade != nil {
			automaticOSUpgradeIsEnabled = *policy.AutomaticOSUpgradePolicy.EnableAutomaticOSUpgrade
		}
	}

	if d.HasChange("automatic_os_upgrade_policy") || d.HasChange("rolling_upgrade_policy") {
		upgradePolicy := compute.UpgradePolicy{}
		if existing.VirtualMachineScaleSetProperties.UpgradePolicy == nil {
			upgradePolicy = compute.UpgradePolicy{
				Mode: compute.UpgradeMode(d.Get("upgrade_mode").(string)),
			}
		} else {
			upgradePolicy = *existing.VirtualMachineScaleSetProperties.UpgradePolicy
			upgradePolicy.Mode = compute.UpgradeMode(d.Get("upgrade_mode").(string))
		}

		if d.HasChange("automatic_os_upgrade_policy") {
			automaticRaw := d.Get("automatic_os_upgrade_policy").([]interface{})
			upgradePolicy.AutomaticOSUpgradePolicy = ExpandVirtualMachineScaleSetAutomaticUpgradePolicy(automaticRaw)

			if upgradePolicy.AutomaticOSUpgradePolicy != nil {
				automaticOSUpgradeIsEnabled = *upgradePolicy.AutomaticOSUpgradePolicy.EnableAutomaticOSUpgrade
			}
		}

		if d.HasChange("rolling_upgrade_policy") {
			rollingRaw := d.Get("rolling_upgrade_policy").([]interface{})
			zones := zones.ExpandUntyped(d.Get("zones").(*schema.Set).List())
			rollingUpgradePolicy, err := ExpandVirtualMachineScaleSetRollingUpgradePolicy(rollingRaw, len(zones) > 0)
			if err != nil {
				return err
			}
			upgradePolicy.RollingUpgradePolicy = rollingUpgradePolicy
		}

		updateProps.UpgradePolicy = &upgradePolicy
	}

	priority := compute.VirtualMachinePriorityTypes(d.Get("priority").(string))
	if d.HasChange("max_bid_price") {
		if priority != compute.VirtualMachinePriorityTypesSpot {
			return fmt.Errorf("`max_bid_price` can only be configured when `priority` is set to `Spot`")
		}

		updateProps.VirtualMachineProfile.BillingProfile = &compute.BillingProfile{
			MaxPrice: utils.Float(d.Get("max_bid_price").(float64)),
		}
	}

	if d.HasChange("single_placement_group") {
		singlePlacementGroup := d.Get("single_placement_group").(bool)
		if singlePlacementGroup {
			return fmt.Errorf("%q can not be set to %q once it has been set to %q", "single_placement_group", "true", "false")
		}
		updateProps.SinglePlacementGroup = utils.Bool(singlePlacementGroup)
	}

	if d.HasChange("admin_ssh_key") || d.HasChange("custom_data") || d.HasChange("disable_password_authentication") || d.HasChange("provision_vm_agent") || d.HasChange("secret") {
		osProfile := compute.VirtualMachineScaleSetUpdateOSProfile{}

		if d.HasChange("admin_ssh_key") || d.HasChange("disable_password_authentication") || d.HasChange("provision_vm_agent") {
			linuxConfig := compute.LinuxConfiguration{}

			if d.HasChange("admin_ssh_key") {
				sshKeysRaw := d.Get("admin_ssh_key").(*pluginsdk.Set).List()
				sshKeys := ExpandSSHKeys(sshKeysRaw)
				linuxConfig.SSH = &compute.SSHConfiguration{
					PublicKeys: &sshKeys,
				}
			}

			if d.HasChange("disable_password_authentication") {
				linuxConfig.DisablePasswordAuthentication = utils.Bool(d.Get("disable_password_authentication").(bool))
			}

			if d.HasChange("provision_vm_agent") {
				linuxConfig.ProvisionVMAgent = utils.Bool(d.Get("provision_vm_agent").(bool))
			}

			osProfile.LinuxConfiguration = &linuxConfig
		}

		if d.HasChange("custom_data") {
			updateInstances = true

			// customData can only be sent if it's a base64 encoded string,
			// so it's not possible to remove this without tainting the resource
			if v, ok := d.GetOk("custom_data"); ok {
				osProfile.CustomData = utils.String(v.(string))
			}
		}

		if d.HasChange("secret") {
			secretsRaw := d.Get("secret").([]interface{})
			osProfile.Secrets = expandLinuxSecrets(secretsRaw)
		}

		updateProps.VirtualMachineProfile.OsProfile = &osProfile
	}

	if d.HasChange("data_disk") || d.HasChange("os_disk") || d.HasChange("source_image_id") || d.HasChange("source_image_reference") {
		updateInstances = true

		if updateProps.VirtualMachineProfile.StorageProfile == nil {
			updateProps.VirtualMachineProfile.StorageProfile = &compute.VirtualMachineScaleSetUpdateStorageProfile{}
		}

		if d.HasChange("data_disk") {
			ultraSSDEnabled := d.Get("additional_capabilities.0.ultra_ssd_enabled").(bool)
			dataDisks, err := ExpandVirtualMachineScaleSetDataDisk(d.Get("data_disk").([]interface{}), ultraSSDEnabled)
			if err != nil {
				return fmt.Errorf("expanding `data_disk`: %+v", err)
			}
			updateProps.VirtualMachineProfile.StorageProfile.DataDisks = dataDisks
		}

		if d.HasChange("os_disk") {
			osDiskRaw := d.Get("os_disk").([]interface{})
			updateProps.VirtualMachineProfile.StorageProfile.OsDisk = ExpandVirtualMachineScaleSetOSDiskUpdate(osDiskRaw)
		}

		if d.HasChange("source_image_id") || d.HasChange("source_image_reference") {
			sourceImageReferenceRaw := d.Get("source_image_reference").([]interface{})
			sourceImageId := d.Get("source_image_id").(string)
			sourceImageReference := expandSourceImageReference(sourceImageReferenceRaw, sourceImageId)

			// Must include all storage profile properties when updating disk image.  See: https://github.com/hashicorp/terraform-provider-azurerm/issues/8273
			updateProps.VirtualMachineProfile.StorageProfile.DataDisks = existing.VirtualMachineScaleSetProperties.VirtualMachineProfile.StorageProfile.DataDisks
			updateProps.VirtualMachineProfile.StorageProfile.ImageReference = sourceImageReference
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

	if d.HasChange("network_interface") || d.HasChange("health_probe_id") {
		networkInterfacesRaw := d.Get("network_interface").([]interface{})
		networkInterfaces, err := ExpandVirtualMachineScaleSetNetworkInterfaceUpdate(networkInterfacesRaw)
		if err != nil {
			return fmt.Errorf("expanding `network_interface`: %+v", err)
		}

		updateProps.VirtualMachineProfile.NetworkProfile = &compute.VirtualMachineScaleSetUpdateNetworkProfile{
			NetworkInterfaceConfigurations: networkInterfaces,
		}

		healthProbeId := d.Get("health_probe_id").(string)
		if healthProbeId != "" {
			updateProps.VirtualMachineProfile.NetworkProfile.HealthProbe = &compute.APIEntityReference{
				ID: utils.String(healthProbeId),
			}
		}
	}

	if d.HasChange("boot_diagnostics") {
		updateInstances = true

		bootDiagnosticsRaw := d.Get("boot_diagnostics").([]interface{})
		updateProps.VirtualMachineProfile.DiagnosticsProfile = expandBootDiagnostics(bootDiagnosticsRaw)
	}

	if d.HasChange("do_not_run_extensions_on_overprovisioned_machines") {
		v := d.Get("do_not_run_extensions_on_overprovisioned_machines").(bool)
		updateProps.DoNotRunExtensionsOnOverprovisionedVMs = utils.Bool(v)
	}

	if d.HasChange("overprovision") {
		v := d.Get("overprovision").(bool)
		updateProps.Overprovision = utils.Bool(v)
	}

	if d.HasChange("scale_in") {
		if updateScaleInPolicy := ExpandVirtualMachineScaleSetScaleInPolicy(d.Get("scale_in").([]interface{})); updateScaleInPolicy != nil {
			updateProps.ScaleInPolicy = updateScaleInPolicy
		}
	}

	if !features.FourPointOhBeta() {
		if d.HasChange("scale_in_policy") {
			updateScaleInPolicy := &compute.ScaleInPolicy{}

			if d.HasChange("scale_in_policy") {
				scaleInPolicy := d.Get("scale_in_policy").(string)
				updateScaleInPolicy.Rules = &[]compute.VirtualMachineScaleSetScaleInRules{compute.VirtualMachineScaleSetScaleInRules(scaleInPolicy)}
			}

			updateProps.ScaleInPolicy = updateScaleInPolicy
		}

		if d.HasChange("terminate_notification") {
			notificationRaw := d.Get("terminate_notification").([]interface{})
			updateProps.VirtualMachineProfile.ScheduledEventsProfile = ExpandVirtualMachineScaleSetScheduledEventsProfile(notificationRaw)
		}
	}

	if d.HasChange("termination_notification") {
		notificationRaw := d.Get("termination_notification").([]interface{})
		updateProps.VirtualMachineProfile.ScheduledEventsProfile = ExpandVirtualMachineScaleSetScheduledEventsProfile(notificationRaw)
	}

	if d.HasChange("encryption_at_host_enabled") {
		if d.Get("encryption_at_host_enabled").(bool) {
			osDiskRaw := d.Get("os_disk").([]interface{})
			securityEncryptionType := osDiskRaw[0].(map[string]interface{})["security_encryption_type"].(string)
			if compute.SecurityEncryptionTypesDiskWithVMGuestState == compute.SecurityEncryptionTypes(securityEncryptionType) {
				return fmt.Errorf("`encryption_at_host_enabled` cannot be set to `true` when `os_disk.0.security_encryption_type` is set to `DiskWithVMGuestState`")
			}
		}

		updateProps.VirtualMachineProfile.SecurityProfile = &compute.SecurityProfile{
			EncryptionAtHost: utils.Bool(d.Get("encryption_at_host_enabled").(bool)),
		}
	}

	if d.HasChange("automatic_instance_repair") {
		automaticRepairsPolicyRaw := d.Get("automatic_instance_repair").([]interface{})
		updateProps.AutomaticRepairsPolicy = ExpandVirtualMachineScaleSetAutomaticRepairsPolicy(automaticRepairsPolicyRaw)
	}

	if d.HasChange("identity") {
		identityRaw := d.Get("identity").([]interface{})
		identity, err := expandVirtualMachineScaleSetIdentity(identityRaw)
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}

		update.Identity = identity
	}

	if d.HasChange("plan") {
		planRaw := d.Get("plan").([]interface{})
		update.Plan = expandPlan(planRaw)
	}

	if d.HasChange("sku") || d.HasChange("instances") {
		// in-case ignore_changes is being used, since both fields are required
		// look up the current values and override them as needed
		sku := existing.Sku

		if d.HasChange("sku") {
			updateInstances = true

			sku.Name = utils.String(d.Get("sku").(string))
		}

		if d.HasChange("instances") {
			sku.Capacity = utils.Int64(int64(d.Get("instances").(int)))
		}

		update.Sku = sku
	}

	if d.HasChanges("extension", "extensions_time_budget") {
		updateInstances = true

		extensionProfile, _, err := expandVirtualMachineScaleSetExtensions(d.Get("extension").(*pluginsdk.Set).List())
		if err != nil {
			return err
		}
		updateProps.VirtualMachineProfile.ExtensionProfile = extensionProfile
		updateProps.VirtualMachineProfile.ExtensionProfile.ExtensionsTimeBudget = utils.String(d.Get("extensions_time_budget").(string))
	}

	if d.HasChange("tags") {
		update.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if d.HasChange("user_data") {
		updateInstances = true
		updateProps.VirtualMachineProfile.UserData = utils.String(d.Get("user_data").(string))
	}

	update.VirtualMachineScaleSetUpdateProperties = &updateProps

	metaData := virtualMachineScaleSetUpdateMetaData{
		AutomaticOSUpgradeIsEnabled:  automaticOSUpgradeIsEnabled,
		CanRollInstancesWhenRequired: meta.(*clients.Client).Features.VirtualMachineScaleSet.RollInstancesWhenRequired,
		UpdateInstances:              updateInstances,
		Client:                       meta.(*clients.Client).Compute,
		Existing:                     existing,
		ID:                           id,
		OSType:                       compute.OperatingSystemTypesLinux,
	}

	if err := metaData.performUpdate(ctx, update); err != nil {
		return err
	}

	return resourceLinuxVirtualMachineScaleSetRead(d, meta)
}

func resourceLinuxVirtualMachineScaleSetRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VMScaleSetClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualMachineScaleSetID(d.Id())
	if err != nil {
		return err
	}

	// Upgrading to the 2021-07-01 exposed a new expand parameter to the GET method
	resp, err := client.Get(ctx, id.ResourceGroup, id.Name, compute.ExpandTypesForGetVMScaleSetsUserData)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Linux Virtual Machine Scale Set %q was not found in Resource Group %q - removing from state!", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Linux Virtual Machine Scale Set %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))
	d.Set("edge_zone", flattenEdgeZone(resp.ExtendedLocation))
	d.Set("zones", zones.FlattenUntyped(resp.Zones))

	var skuName *string
	var instances int
	if resp.Sku != nil {
		skuName = resp.Sku.Name
		if resp.Sku.Capacity != nil {
			instances = int(*resp.Sku.Capacity)
		}
	}
	d.Set("instances", instances)
	d.Set("sku", skuName)

	identity, err := flattenVirtualMachineScaleSetIdentity(resp.Identity)
	if err != nil {
		return err
	}
	if err := d.Set("identity", identity); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	if err := d.Set("plan", flattenPlan(resp.Plan)); err != nil {
		return fmt.Errorf("setting `plan`: %+v", err)
	}

	if resp.VirtualMachineScaleSetProperties == nil {
		return fmt.Errorf("retrieving Linux Virtual Machine Scale Set %q (Resource Group %q): `properties` was nil", id.Name, id.ResourceGroup)
	}
	props := *resp.VirtualMachineScaleSetProperties

	if err := d.Set("additional_capabilities", FlattenVirtualMachineScaleSetAdditionalCapabilities(props.AdditionalCapabilities)); err != nil {
		return fmt.Errorf("setting `additional_capabilities`: %+v", props.AdditionalCapabilities)
	}

	if err := d.Set("automatic_instance_repair", FlattenVirtualMachineScaleSetAutomaticRepairsPolicy(props.AutomaticRepairsPolicy)); err != nil {
		return fmt.Errorf("setting `automatic_instance_repair`: %+v", err)
	}

	d.Set("do_not_run_extensions_on_overprovisioned_machines", props.DoNotRunExtensionsOnOverprovisionedVMs)
	if props.HostGroup != nil && props.HostGroup.ID != nil {
		d.Set("host_group_id", props.HostGroup.ID)
	}
	d.Set("overprovision", props.Overprovision)
	proximityPlacementGroupId := ""
	if props.ProximityPlacementGroup != nil && props.ProximityPlacementGroup.ID != nil {
		proximityPlacementGroupId = *props.ProximityPlacementGroup.ID
	}
	d.Set("platform_fault_domain_count", props.PlatformFaultDomainCount)
	d.Set("proximity_placement_group_id", proximityPlacementGroupId)
	d.Set("single_placement_group", props.SinglePlacementGroup)
	d.Set("unique_id", props.UniqueID)
	d.Set("zone_balance", props.ZoneBalance)
	d.Set("scale_in", FlattenVirtualMachineScaleSetScaleInPolicy(props.ScaleInPolicy))

	if !features.FourPointOhBeta() {
		rule := string(compute.VirtualMachineScaleSetScaleInRulesDefault)

		if props.ScaleInPolicy != nil {
			if rules := props.ScaleInPolicy.Rules; rules != nil && len(*rules) > 0 {
				rule = string((*rules)[0])
			}
		}

		d.Set("scale_in_policy", rule)
	}

	if props.SpotRestorePolicy != nil {
		d.Set("spot_restore", FlattenVirtualMachineScaleSetSpotRestorePolicy(props.SpotRestorePolicy))
	}

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

		if profile.ApplicationProfile != nil && profile.ApplicationProfile.GalleryApplications != nil {
			d.Set("gallery_application", flattenVirtualMachineScaleSetGalleryApplication(profile.ApplicationProfile.GalleryApplications))

			if !features.FourPointOhBeta() {
				d.Set("gallery_applications", flattenVirtualMachineScaleSetGalleryApplications(profile.ApplicationProfile.GalleryApplications))
			}
		}

		// the service just return empty when this is not assigned when provisioned
		// See discussion on https://github.com/Azure/azure-rest-api-specs/issues/10971
		priority := compute.VirtualMachinePriorityTypesRegular
		if profile.Priority != "" {
			priority = profile.Priority
		}
		d.Set("priority", priority)

		if storageProfile := profile.StorageProfile; storageProfile != nil {
			if err := d.Set("os_disk", FlattenVirtualMachineScaleSetOSDisk(storageProfile.OsDisk)); err != nil {
				return fmt.Errorf("setting `os_disk`: %+v", err)
			}

			if err := d.Set("data_disk", FlattenVirtualMachineScaleSetDataDisk(storageProfile.DataDisks)); err != nil {
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

		extensionOperationsEnabled := true
		if osProfile := profile.OsProfile; osProfile != nil {
			// admin_password isn't returned, but it's a top level field so we can ignore it without consequence
			d.Set("admin_username", osProfile.AdminUsername)
			d.Set("computer_name_prefix", osProfile.ComputerNamePrefix)

			if osProfile.AllowExtensionOperations != nil {
				extensionOperationsEnabled = *osProfile.AllowExtensionOperations
			}

			if linux := osProfile.LinuxConfiguration; linux != nil {
				d.Set("disable_password_authentication", linux.DisablePasswordAuthentication)
				d.Set("provision_vm_agent", linux.ProvisionVMAgent)

				flattenedSshKeys, err := FlattenSSHKeys(linux.SSH)
				if err != nil {
					return fmt.Errorf("flattening `admin_ssh_key`: %+v", err)
				}
				if err := d.Set("admin_ssh_key", pluginsdk.NewSet(SSHKeySchemaHash, *flattenedSshKeys)); err != nil {
					return fmt.Errorf("setting `admin_ssh_key`: %+v", err)
				}
			}

			if err := d.Set("secret", flattenLinuxSecrets(osProfile.Secrets)); err != nil {
				return fmt.Errorf("setting `secret`: %+v", err)
			}
		}
		d.Set("extension_operations_enabled", extensionOperationsEnabled)

		if nwProfile := profile.NetworkProfile; nwProfile != nil {
			flattenedNics := FlattenVirtualMachineScaleSetNetworkInterface(nwProfile.NetworkInterfaceConfigurations)
			if err := d.Set("network_interface", flattenedNics); err != nil {
				return fmt.Errorf("setting `network_interface`: %+v", err)
			}

			healthProbeId := ""
			if nwProfile.HealthProbe != nil && nwProfile.HealthProbe.ID != nil {
				healthProbeId = *nwProfile.HealthProbe.ID
			}
			d.Set("health_probe_id", healthProbeId)
		}

		if !features.FourPointOhBeta() {
			if scheduleProfile := profile.ScheduledEventsProfile; scheduleProfile != nil {
				if err := d.Set("terminate_notification", FlattenVirtualMachineScaleSetScheduledEventsProfile(scheduleProfile)); err != nil {
					return fmt.Errorf("setting `terminate_notification`: %+v", err)
				}
			}
		}

		if scheduleProfile := profile.ScheduledEventsProfile; scheduleProfile != nil {
			if err := d.Set("termination_notification", FlattenVirtualMachineScaleSetScheduledEventsProfile(scheduleProfile)); err != nil {
				return fmt.Errorf("setting `termination_notification`: %+v", err)
			}
		}

		extensionProfile, err := flattenVirtualMachineScaleSetExtensions(profile.ExtensionProfile, d)
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
		vtpmEnabled := false
		secureBootEnabled := false

		if secprofile := profile.SecurityProfile; secprofile != nil {
			if secprofile.EncryptionAtHost != nil {
				encryptionAtHostEnabled = *secprofile.EncryptionAtHost
			}
			if uefi := profile.SecurityProfile.UefiSettings; uefi != nil {
				if uefi.VTpmEnabled != nil {
					vtpmEnabled = *uefi.VTpmEnabled
				}
				if uefi.SecureBootEnabled != nil {
					secureBootEnabled = *uefi.SecureBootEnabled
				}
			}
		}

		d.Set("encryption_at_host_enabled", encryptionAtHostEnabled)
		d.Set("vtpm_enabled", vtpmEnabled)
		d.Set("secure_boot_enabled", secureBootEnabled)
		d.Set("user_data", profile.UserData)
	}

	if policy := props.UpgradePolicy; policy != nil {
		d.Set("upgrade_mode", string(policy.Mode))

		flattenedAutomatic := FlattenVirtualMachineScaleSetAutomaticOSUpgradePolicy(policy.AutomaticOSUpgradePolicy)
		if err := d.Set("automatic_os_upgrade_policy", flattenedAutomatic); err != nil {
			return fmt.Errorf("setting `automatic_os_upgrade_policy`: %+v", err)
		}

		flattenedRolling := FlattenVirtualMachineScaleSetRollingUpgradePolicy(policy.RollingUpgradePolicy)
		if err := d.Set("rolling_upgrade_policy", flattenedRolling); err != nil {
			return fmt.Errorf("setting `rolling_upgrade_policy`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceLinuxVirtualMachineScaleSetDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VMScaleSetClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualMachineScaleSetID(d.Id())
	if err != nil {
		return err
	}

	// Upgrading to the 2021-07-01 exposed a new expand parameter to the GET method
	resp, err := client.Get(ctx, id.ResourceGroup, id.Name, compute.ExpandTypesForGetVMScaleSetsUserData)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return nil
		}

		return fmt.Errorf("retrieving Linux Virtual Machine Scale Set %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	// If rolling upgrades are configured and running we need to cancel them before trying to delete the VMSS
	if err := meta.(*clients.Client).Compute.CancelRollingUpgradesBeforeDeletion(ctx, *id); err != nil {
		return fmt.Errorf("cancelling rolling upgrades for %s: %+v", *id, err)
	}

	// Sometimes VMSS's aren't fully deleted when the `Delete` call returns - as such we'll try to scale the cluster
	// to 0 nodes first, then delete the cluster - which should ensure there's no Network Interfaces kicking around
	// and work around this Azure API bug:
	// Original Error: Code="InUseSubnetCannotBeDeleted" Message="Subnet internal is in use by
	// /{nicResourceID}/|providers|Microsoft.Compute|virtualMachineScaleSets|acctestvmss-190923101253410278|virtualMachines|0|networkInterfaces|example/ipConfigurations/internal and cannot be deleted.
	// In order to delete the subnet, delete all the resources within the subnet. See aka.ms/deletesubnet.
	scaleToZeroOnDelete := meta.(*clients.Client).Features.VirtualMachineScaleSet.ScaleToZeroOnDelete
	if scaleToZeroOnDelete && resp.Sku != nil {
		resp.Sku.Capacity = utils.Int64(int64(0))

		log.Printf("[DEBUG] Scaling instances to 0 prior to deletion - this helps avoids networking issues within Azure")
		update := compute.VirtualMachineScaleSetUpdate{
			Sku: resp.Sku,
		}
		future, err := client.Update(ctx, id.ResourceGroup, id.Name, update)
		if err != nil {
			return fmt.Errorf("updating number of instances in Linux Virtual Machine Scale Set %q (Resource Group %q) to scale to 0: %+v", id.Name, id.ResourceGroup, err)
		}

		log.Printf("[DEBUG] Waiting for scaling of instances to 0 prior to deletion - this helps avoids networking issues within Azure")
		err = future.WaitForCompletionRef(ctx, client.Client)
		if err != nil {
			return fmt.Errorf("waiting for number of instances in Linux Virtual Machine Scale Set %q (Resource Group %q) to scale to 0: %+v", id.Name, id.ResourceGroup, err)
		}
		log.Printf("[DEBUG] Scaled instances to 0 prior to deletion - this helps avoids networking issues within Azure")
	} else {
		log.Printf("[DEBUG] Unable to scale instances to `0` since the `sku` block is nil - trying to delete anyway")
	}

	log.Printf("[DEBUG] Deleting Linux Virtual Machine Scale Set %q (Resource Group %q)..", id.Name, id.ResourceGroup)
	// @ArcturusZhang (mimicking from linux_virtual_machine_pluginsdk.go): sending `nil` here omits this value from being sent
	// which matches the previous behaviour - we're only splitting this out so it's clear why
	// TODO: support force deletion once it's out of Preview, if applicable
	var forceDeletion *bool = nil
	future, err := client.Delete(ctx, id.ResourceGroup, id.Name, forceDeletion)
	if err != nil {
		return fmt.Errorf("deleting Linux Virtual Machine Scale Set %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	log.Printf("[DEBUG] Waiting for deletion of Linux Virtual Machine Scale Set %q (Resource Group %q)..", id.Name, id.ResourceGroup)
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of Linux Virtual Machine Scale Set %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	log.Printf("[DEBUG] Deleted Linux Virtual Machine Scale Set %q (Resource Group %q).", id.Name, id.ResourceGroup)

	return nil
}

func resourceLinuxVirtualMachineScaleSetSchema() map[string]*pluginsdk.Schema {
	resourceSchema := map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.VirtualMachineName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		// Required
		"admin_username": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"network_interface": VirtualMachineScaleSetNetworkInterfaceSchema(),

		"os_disk": VirtualMachineScaleSetOSDiskSchema(),

		"instances": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      0,
			ValidateFunc: validation.IntAtLeast(0),
		},

		"sku": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		// Optional
		"additional_capabilities": VirtualMachineScaleSetAdditionalCapabilitiesSchema(),

		"admin_password": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ForceNew:         true,
			Sensitive:        true,
			DiffSuppressFunc: adminPasswordDiffSuppressFunc,
		},

		"admin_ssh_key": SSHKeysSchema(false),

		"automatic_os_upgrade_policy": VirtualMachineScaleSetAutomatedOSUpgradePolicySchema(),

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

		"computer_name_prefix": {
			Type:     pluginsdk.TypeString,
			Optional: true,

			// Computed since we reuse the VM name if one's not specified
			Computed: true,
			ForceNew: true,

			ValidateFunc: validate.LinuxComputerNamePrefix,
		},

		"custom_data": base64.OptionalSchema(false),

		"data_disk": VirtualMachineScaleSetDataDiskSchema(),

		"disable_password_authentication": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"do_not_run_extensions_on_overprovisioned_machines": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"edge_zone": commonschema.EdgeZoneOptionalForceNew(),

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

		"extension": VirtualMachineScaleSetExtensionsSchema(),

		"extensions_time_budget": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "PT1H30M",
			ValidateFunc: azValidate.ISO8601DurationBetween("PT15M", "PT2H"),
		},

		"gallery_application": VirtualMachineScaleSetGalleryApplicationSchema(),

		"health_probe_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: azure.ValidateResourceID,
		},

		"host_group_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			// the Compute/VM API is broken and returns the Resource Group name in UPPERCASE
			// tracked by https://github.com/Azure/azure-rest-api-specs/issues/19424
			DiffSuppressFunc: suppress.CaseDifference,
			ValidateFunc:     validate.HostGroupID,
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"max_bid_price": {
			Type:         pluginsdk.TypeFloat,
			Optional:     true,
			Default:      -1,
			ValidateFunc: validate.SpotMaxPrice,
		},

		"overprovision": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"plan": planSchema(),

		"platform_fault_domain_count": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			ForceNew: true,
			Computed: true,
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

		"provision_vm_agent": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
			ForceNew: true,
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

		"secret": linuxSecretSchema(),

		"secure_boot_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
		},

		"single_placement_group": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"source_image_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.Any(
				images.ValidateImageID,
				validate.SharedImageID,
				validate.SharedImageVersionID,
				validate.CommunityGalleryImageID,
				validate.CommunityGalleryImageVersionID,
				validate.SharedGalleryImageID,
				validate.SharedGalleryImageVersionID,
			),
			ExactlyOneOf: []string{
				"source_image_id",
				"source_image_reference",
			},
		},

		"source_image_reference": sourceImageReferenceSchema(false),

		"tags": tags.Schema(),

		"upgrade_mode": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  string(compute.UpgradeModeManual),
			ValidateFunc: validation.StringInSlice([]string{
				string(compute.UpgradeModeAutomatic),
				string(compute.UpgradeModeManual),
				string(compute.UpgradeModeRolling),
			}, false),
		},

		"user_data": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsBase64,
		},

		"vtpm_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
		},

		"zone_balance": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
			Default:  false,
		},

		"scale_in": VirtualMachineScaleSetScaleInPolicySchema(),

		"spot_restore": VirtualMachineScaleSetSpotRestorePolicySchema(),

		"termination_notification": VirtualMachineScaleSetTerminationNotificationSchema(),

		"zones": commonschema.ZonesMultipleOptionalForceNew(),

		// Computed
		"unique_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}

	if !features.FourPointOhBeta() {
		resourceSchema["gallery_applications"] = VirtualMachineScaleSetGalleryApplicationsSchema()
		resourceSchema["terminate_notification"] = VirtualMachineScaleSetTerminateNotificationSchema()

		resourceSchema["scale_in_policy"] = &schema.Schema{
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: !features.FourPointOhBeta(),
			ValidateFunc: validation.StringInSlice([]string{
				string(compute.VirtualMachineScaleSetScaleInRulesDefault),
				string(compute.VirtualMachineScaleSetScaleInRulesNewestVM),
				string(compute.VirtualMachineScaleSetScaleInRulesOldestVM),
			}, false),
			Deprecated:    "`scale_in_policy` will be removed in favour of the `scale_in` code block in version 4.0 of the AzureRM Provider.",
			ConflictsWith: []string{"scale_in"},
		}
	}

	return resourceSchema
}
