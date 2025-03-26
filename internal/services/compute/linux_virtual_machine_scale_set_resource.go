// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"context"
	"fmt"
	"log"
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
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	azValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/base64"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceLinuxVirtualMachineScaleSet() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLinuxVirtualMachineScaleSetCreate,
		Read:   resourceLinuxVirtualMachineScaleSetRead,
		Update: resourceLinuxVirtualMachineScaleSetUpdate,
		Delete: resourceLinuxVirtualMachineScaleSetDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := commonids.ParseVirtualMachineScaleSetID(id)
			return err
		}, importVirtualMachineScaleSet(virtualmachinescalesets.OperatingSystemTypesLinux, "azurerm_linux_virtual_machine_scale_set")),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(time.Minute * 60),
			Read:   pluginsdk.DefaultTimeout(time.Minute * 5),
			Update: pluginsdk.DefaultTimeout(time.Minute * 60),
			Delete: pluginsdk.DefaultTimeout(time.Minute * 60),
		},

		// TODO: exposing requireGuestProvisionSignal once it's available
		// https://github.com/Azure/azure-rest-api-specs/pull/7246

		Schema: resourceLinuxVirtualMachineScaleSetSchema(),

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
		),
	}
}

func resourceLinuxVirtualMachineScaleSetCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VirtualMachineScaleSetsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := virtualmachinescalesets.NewVirtualMachineScaleSetID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	// Upgrading to the 2021-07-01 exposed a new expand parameter to the GET method
	exists, err := client.Get(ctx, id, virtualmachinescalesets.DefaultGetOperationOptions())
	if err != nil {
		if !response.WasNotFound(exists.HttpResponse) {
			return fmt.Errorf("checking for existing Linux %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(exists.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_linux_virtual_machine_scale_set", id.ID())
	}

	location := azure.NormalizeLocation(d.Get("location").(string))

	additionalCapabilitiesRaw := d.Get("additional_capabilities").([]interface{})
	additionalCapabilities := ExpandVirtualMachineScaleSetAdditionalCapabilities(additionalCapabilitiesRaw)

	bootDiagnosticsRaw := d.Get("boot_diagnostics").([]interface{})
	bootDiagnostics := expandBootDiagnosticsVMSS(bootDiagnosticsRaw)

	dataDisksRaw := d.Get("data_disk").([]interface{})
	ultraSSDEnabled := d.Get("additional_capabilities.0.ultra_ssd_enabled").(bool)
	dataDisks, err := ExpandVirtualMachineScaleSetDataDisk(dataDisksRaw, ultraSSDEnabled)
	if err != nil {
		return fmt.Errorf("expanding `data_disk`: %+v", err)
	}

	identityExpanded, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	networkInterfacesRaw := d.Get("network_interface").([]interface{})
	networkInterfaces, err := ExpandVirtualMachineScaleSetNetworkInterface(networkInterfacesRaw)
	if err != nil {
		return fmt.Errorf("expanding `network_interface`: %+v", err)
	}

	osDiskRaw := d.Get("os_disk").([]interface{})
	osDisk, err := ExpandVirtualMachineScaleSetOSDisk(osDiskRaw, virtualmachinescalesets.OperatingSystemTypesLinux)
	if err != nil {
		return fmt.Errorf("expanding `os_disk`: %+v", err)
	}
	securityEncryptionType := osDiskRaw[0].(map[string]interface{})["security_encryption_type"].(string)

	planRaw := d.Get("plan").([]interface{})
	plan := expandPlanVMSS(planRaw)

	sourceImageReferenceRaw := d.Get("source_image_reference").([]interface{})
	sourceImageId := d.Get("source_image_id").(string)
	sourceImageReference := expandSourceImageReferenceVMSS(sourceImageReferenceRaw, sourceImageId)

	sshKeysRaw := d.Get("admin_ssh_key").(*pluginsdk.Set).List()
	sshKeys := expandSSHKeysVMSS(sshKeysRaw)

	overProvision := d.Get("overprovision").(bool)
	provisionVMAgent := d.Get("provision_vm_agent").(bool)
	zones := zones.ExpandUntyped(d.Get("zones").(*schema.Set).List())
	healthProbeId := d.Get("health_probe_id").(string)
	upgradeMode := virtualmachinescalesets.UpgradeMode(d.Get("upgrade_mode").(string))
	automaticOSUpgradePolicyRaw := d.Get("automatic_os_upgrade_policy").([]interface{})
	automaticOSUpgradePolicy := ExpandVirtualMachineScaleSetAutomaticUpgradePolicy(automaticOSUpgradePolicyRaw)
	rollingUpgradePolicyRaw := d.Get("rolling_upgrade_policy").([]interface{})
	rollingUpgradePolicy, err := ExpandVirtualMachineScaleSetRollingUpgradePolicy(rollingUpgradePolicyRaw, len(zones) > 0, overProvision)
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

	secretsRaw := d.Get("secret").([]interface{})
	secrets := expandLinuxSecretsVMSS(secretsRaw)

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

	if galleryApplications := expandVirtualMachineScaleSetGalleryApplication(d.Get("gallery_application").([]interface{})); galleryApplications != nil {
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
		virtualMachineProfile.ExtensionProfile, hasHealthExtension, err = expandVirtualMachineScaleSetExtensions(vmExtensionsRaw.(*pluginsdk.Set).List())
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
		if encryptionAtHostEnabled.(bool) {
			if virtualmachinescalesets.SecurityEncryptionTypesDiskWithVMGuestState == virtualmachinescalesets.SecurityEncryptionTypes(securityEncryptionType) {
				return fmt.Errorf("`encryption_at_host_enabled` cannot be set to `true` when `os_disk.0.security_encryption_type` is set to `DiskWithVMGuestState`")
			}
		}

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
		virtualMachineProfile.EvictionPolicy = pointer.To(virtualmachinescalesets.VirtualMachineEvictionPolicyTypes(evictionPolicyRaw.(string)))
	} else if priority == virtualmachinescalesets.VirtualMachinePriorityTypesSpot {
		return fmt.Errorf("an `eviction_policy` must be specified when `priority` is set to `Spot`")
	}

	if v, ok := d.GetOk("termination_notification"); ok {
		virtualMachineProfile.ScheduledEventsProfile = ExpandVirtualMachineScaleSetScheduledEventsProfile(v.([]interface{}))
	}

	automaticRepairsPolicyRaw := d.Get("automatic_instance_repair").([]interface{})
	automaticRepairsPolicy := ExpandVirtualMachineScaleSetAutomaticRepairsPolicy(automaticRepairsPolicyRaw)

	if automaticRepairsPolicy != nil && healthProbeId == "" && !hasHealthExtension {
		return fmt.Errorf("`automatic_instance_repair` can only be set if there is an application Health extension or a `health_probe_id` defined")
	}

	props := virtualmachinescalesets.VirtualMachineScaleSet{
		ExtendedLocation: expandEdgeZone(d.Get("edge_zone").(string)),
		Location:         location,
		Sku: &virtualmachinescalesets.Sku{
			Name:     pointer.To(d.Get("sku").(string)),
			Capacity: utils.Int64(int64(d.Get("instances").(int))),

			// doesn't appear this can be set to anything else, even Promo machines are Standard
			Tier: pointer.To("Standard"),
		},
		Identity: identityExpanded,
		Plan:     plan,
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
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
		if v := ExpandVirtualMachineScaleSetScaleInPolicy(v.([]interface{})); v != nil {
			props.Properties.ScaleInPolicy = v
		}
	}

	if v, ok := d.GetOk("host_group_id"); ok {
		props.Properties.HostGroup = &virtualmachinescalesets.SubResource{
			Id: pointer.To(v.(string)),
		}
	}

	spotRestoreRaw := d.Get("spot_restore").([]interface{})
	if spotRestorePolicy := ExpandVirtualMachineScaleSetSpotRestorePolicy(spotRestoreRaw); spotRestorePolicy != nil {
		props.Properties.SpotRestorePolicy = spotRestorePolicy
	}

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

	log.Printf("[DEBUG] Creating Linux %s", id)
	if err := client.CreateOrUpdateThenPoll(ctx, id, props, virtualmachinescalesets.DefaultCreateOrUpdateOperationOptions()); err != nil {
		return fmt.Errorf("creating Linux %s: %+v", id, err)
	}
	log.Printf("[DEBUG] %s was created", id)

	d.SetId(id.ID())

	return resourceLinuxVirtualMachineScaleSetRead(d, meta)
}

func resourceLinuxVirtualMachineScaleSetUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VirtualMachineScaleSetsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualmachinescalesets.ParseVirtualMachineScaleSetID(d.Id())
	if err != nil {
		return err
	}

	updateInstances := false
	options := virtualmachinescalesets.DefaultGetOperationOptions()
	options.Expand = pointer.To(virtualmachinescalesets.ExpandTypesForGetVMScaleSetsUserData)
	existing, err := client.Get(ctx, *id, options)
	if err != nil {
		return fmt.Errorf("retrieving Linux %s: %+v", id, err)
	}
	if existing.Model == nil {
		return fmt.Errorf("retrieving Linux %s: `model` was nil", id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving Linux %s: `properties` was nil", id)
	}
	if existing.Model.Properties.VirtualMachineProfile == nil {
		return fmt.Errorf("retrieving Linux %s: `properties.virtualMachineProfile` was nil", id)
	}
	if existing.Model.Properties.VirtualMachineProfile.StorageProfile == nil {
		return fmt.Errorf("retrieving Linux %s: `properties.virtualMachineProfile,storageProfile` was nil", id)
	}

	updateProps := virtualmachinescalesets.VirtualMachineScaleSetUpdateProperties{
		VirtualMachineProfile: &virtualmachinescalesets.VirtualMachineScaleSetUpdateVMProfile{
			// if an image reference has been configured previously (it has to be), we would better to include that in this
			// update request to avoid some circumstances that the API will complain ImageReference is null
			// issue tracking: https://github.com/Azure/azure-rest-api-specs/issues/10322
			StorageProfile: &virtualmachinescalesets.VirtualMachineScaleSetUpdateStorageProfile{
				ImageReference: existing.Model.Properties.VirtualMachineProfile.StorageProfile.ImageReference,
			},
		},
		// if an upgrade policy's been configured previously (which it will have) it must be threaded through
		// this doesn't matter for Manual - but breaks when updating anything on a Automatic and Rolling Mode Scale Set
		UpgradePolicy: existing.Model.Properties.UpgradePolicy,
	}
	update := virtualmachinescalesets.VirtualMachineScaleSetUpdate{}

	// first try and pull this from existing vm, which covers no changes being made to this block
	automaticOSUpgradeIsEnabled := false
	if policy := existing.Model.Properties.UpgradePolicy; policy != nil {
		if policy.AutomaticOSUpgradePolicy != nil && policy.AutomaticOSUpgradePolicy.EnableAutomaticOSUpgrade != nil {
			automaticOSUpgradeIsEnabled = *policy.AutomaticOSUpgradePolicy.EnableAutomaticOSUpgrade
		}
	}

	if d.HasChange("zones") {
		update.Zones = pointer.To(zones.ExpandUntyped(d.Get("zones").(*schema.Set).List()))
	}

	if d.HasChange("automatic_os_upgrade_policy") || d.HasChange("rolling_upgrade_policy") {
		upgradePolicy := virtualmachinescalesets.UpgradePolicy{}
		if existing.Model.Properties.UpgradePolicy == nil {
			upgradePolicy = virtualmachinescalesets.UpgradePolicy{
				Mode: pointer.To(virtualmachinescalesets.UpgradeMode(d.Get("upgrade_mode").(string))),
			}
		} else {
			upgradePolicy = *existing.Model.Properties.UpgradePolicy
			upgradePolicy.Mode = pointer.To(virtualmachinescalesets.UpgradeMode(d.Get("upgrade_mode").(string)))
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
			rollingUpgradePolicy, err := ExpandVirtualMachineScaleSetRollingUpgradePolicy(rollingRaw, len(zones) > 0, d.Get("overprovision").(bool))
			if err != nil {
				return err
			}
			upgradePolicy.RollingUpgradePolicy = rollingUpgradePolicy
		}

		updateProps.UpgradePolicy = &upgradePolicy
	}

	priority := virtualmachinescalesets.VirtualMachinePriorityTypes(d.Get("priority").(string))
	if d.HasChange("max_bid_price") {
		if priority != virtualmachinescalesets.VirtualMachinePriorityTypesSpot {
			return fmt.Errorf("`max_bid_price` can only be configured when `priority` is set to `Spot`")
		}

		updateProps.VirtualMachineProfile.BillingProfile = &virtualmachinescalesets.BillingProfile{
			MaxPrice: pointer.To(d.Get("max_bid_price").(float64)),
		}
	}

	if d.HasChange("single_placement_group") {
		singlePlacementGroup := d.Get("single_placement_group").(bool)
		if singlePlacementGroup {
			return fmt.Errorf("%q can not be set to %q once it has been set to %q", "single_placement_group", "true", "false")
		}
		updateProps.SinglePlacementGroup = pointer.To(singlePlacementGroup)
	}

	if d.HasChange("admin_ssh_key") || d.HasChange("custom_data") || d.HasChange("disable_password_authentication") || d.HasChange("provision_vm_agent") || d.HasChange("secret") {
		osProfile := virtualmachinescalesets.VirtualMachineScaleSetUpdateOSProfile{}

		if d.HasChange("admin_ssh_key") || d.HasChange("disable_password_authentication") || d.HasChange("provision_vm_agent") {
			linuxConfig := virtualmachinescalesets.LinuxConfiguration{}

			if d.HasChange("admin_ssh_key") {
				sshKeysRaw := d.Get("admin_ssh_key").(*pluginsdk.Set).List()
				sshKeys := expandSSHKeysVMSS(sshKeysRaw)
				linuxConfig.Ssh = &virtualmachinescalesets.SshConfiguration{
					PublicKeys: &sshKeys,
				}
			}

			if d.HasChange("disable_password_authentication") {
				linuxConfig.DisablePasswordAuthentication = pointer.To(d.Get("disable_password_authentication").(bool))
			}

			if d.HasChange("provision_vm_agent") {
				linuxConfig.ProvisionVMAgent = pointer.To(d.Get("provision_vm_agent").(bool))
			}

			osProfile.LinuxConfiguration = &linuxConfig
		}

		if d.HasChange("custom_data") {
			updateInstances = true

			// customData can only be sent if it's a base64 encoded string,
			// so it's not possible to remove this without tainting the resource
			if v, ok := d.GetOk("custom_data"); ok {
				osProfile.CustomData = pointer.To(v.(string))
			}
		}

		if d.HasChange("secret") {
			secretsRaw := d.Get("secret").([]interface{})
			osProfile.Secrets = expandLinuxSecretsVMSS(secretsRaw)
		}

		updateProps.VirtualMachineProfile.OsProfile = &osProfile
	}

	if d.HasChange("data_disk") || d.HasChange("os_disk") || d.HasChange("source_image_id") || d.HasChange("source_image_reference") {
		updateInstances = true

		if updateProps.VirtualMachineProfile.StorageProfile == nil {
			updateProps.VirtualMachineProfile.StorageProfile = &virtualmachinescalesets.VirtualMachineScaleSetUpdateStorageProfile{}
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
			sourceImageReference := expandSourceImageReferenceVMSS(sourceImageReferenceRaw, sourceImageId)

			// Must include all storage profile properties when updating disk image.  See: https://github.com/hashicorp/terraform-provider-azurerm/issues/8273
			updateProps.VirtualMachineProfile.StorageProfile.DataDisks = existing.Model.Properties.VirtualMachineProfile.StorageProfile.DataDisks
			updateProps.VirtualMachineProfile.StorageProfile.ImageReference = sourceImageReference
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

	if d.HasChange("network_interface") || d.HasChange("health_probe_id") {
		networkInterfacesRaw := d.Get("network_interface").([]interface{})
		networkInterfaces, err := ExpandVirtualMachineScaleSetNetworkInterfaceUpdate(networkInterfacesRaw)
		if err != nil {
			return fmt.Errorf("expanding `network_interface`: %+v", err)
		}

		updateProps.VirtualMachineProfile.NetworkProfile = &virtualmachinescalesets.VirtualMachineScaleSetUpdateNetworkProfile{
			NetworkInterfaceConfigurations: networkInterfaces,
		}

		healthProbeId := d.Get("health_probe_id").(string)
		if healthProbeId != "" {
			updateProps.VirtualMachineProfile.NetworkProfile.HealthProbe = &virtualmachinescalesets.ApiEntityReference{
				Id: pointer.To(healthProbeId),
			}
		}
	}

	if d.HasChange("boot_diagnostics") {
		updateInstances = true

		bootDiagnosticsRaw := d.Get("boot_diagnostics").([]interface{})
		updateProps.VirtualMachineProfile.DiagnosticsProfile = expandBootDiagnosticsVMSS(bootDiagnosticsRaw)
	}

	if d.HasChange("do_not_run_extensions_on_overprovisioned_machines") {
		v := d.Get("do_not_run_extensions_on_overprovisioned_machines").(bool)
		updateProps.DoNotRunExtensionsOnOverprovisionedVMs = pointer.To(v)
	}

	if d.HasChange("overprovision") {
		v := d.Get("overprovision").(bool)
		updateProps.Overprovision = pointer.To(v)
	}

	if d.HasChange("scale_in") {
		if updateScaleInPolicy := ExpandVirtualMachineScaleSetScaleInPolicy(d.Get("scale_in").([]interface{})); updateScaleInPolicy != nil {
			updateProps.ScaleInPolicy = updateScaleInPolicy
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
			if virtualmachinescalesets.SecurityEncryptionTypesDiskWithVMGuestState == virtualmachinescalesets.SecurityEncryptionTypes(securityEncryptionType) {
				return fmt.Errorf("`encryption_at_host_enabled` cannot be set to `true` when `os_disk.0.security_encryption_type` is set to `DiskWithVMGuestState`")
			}
		}

		updateProps.VirtualMachineProfile.SecurityProfile = &virtualmachinescalesets.SecurityProfile{
			EncryptionAtHost: pointer.To(d.Get("encryption_at_host_enabled").(bool)),
		}
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

			_, hasHealthProbeId := d.GetOk("health_probe_id")

			if !hasHealthProbeId && !hasHealthExtension {
				return fmt.Errorf("`automatic_instance_repair` can only be set if there is an application Health extension or a `health_probe_id` defined")
			}
		}

		updateProps.AutomaticRepairsPolicy = automaticRepairsPolicy
	}

	if d.HasChange("identity") {
		identityExpanded, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}

		existing.Model.Identity = identityExpanded
		// Removing a user-assigned identity using PATCH requires setting it to `null` in the payload which
		// 1. The go-azure-sdk for resource manager doesn't support at the moment
		// 2. The expand identity function doesn't behave this way
		// For the moment updating the identity with the PUT circumvents this API behaviour
		// See https://github.com/hashicorp/terraform-provider-azurerm/issues/25058 for more details
		if err := client.CreateOrUpdateThenPoll(ctx, *id, *existing.Model, virtualmachinescalesets.DefaultCreateOrUpdateOperationOptions()); err != nil {
			return fmt.Errorf("updating identity for Linux %s: %+v", id, err)
		}
	}

	if d.HasChange("plan") {
		planRaw := d.Get("plan").([]interface{})
		update.Plan = expandPlanVMSS(planRaw)
	}

	if d.HasChange("sku") || d.HasChange("instances") {
		// in-case ignore_changes is being used, since both fields are required
		// look up the current values and override them as needed
		sku := existing.Model.Sku

		if d.HasChange("sku") {
			updateInstances = true

			sku.Name = pointer.To(d.Get("sku").(string))
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
		updateProps.VirtualMachineProfile.ExtensionProfile.ExtensionsTimeBudget = pointer.To(d.Get("extensions_time_budget").(string))
	}

	if d.HasChange("tags") {
		update.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if d.HasChange("user_data") {
		updateInstances = true
		updateProps.VirtualMachineProfile.UserData = pointer.To(d.Get("user_data").(string))
	}

	update.Properties = &updateProps

	metaData := virtualMachineScaleSetUpdateMetaData{
		AutomaticOSUpgradeIsEnabled:  automaticOSUpgradeIsEnabled,
		CanReimageOnManualUpgrade:    meta.(*clients.Client).Features.VirtualMachineScaleSet.ReimageOnManualUpgrade,
		CanRollInstancesWhenRequired: meta.(*clients.Client).Features.VirtualMachineScaleSet.RollInstancesWhenRequired,
		UpdateInstances:              updateInstances,
		Client:                       meta.(*clients.Client).Compute,
		Existing:                     *existing.Model,
		ID:                           id,
		OSType:                       virtualmachinescalesets.OperatingSystemTypesLinux,
	}

	if err := metaData.performUpdate(ctx, update); err != nil {
		return err
	}

	return resourceLinuxVirtualMachineScaleSetRead(d, meta)
}

func resourceLinuxVirtualMachineScaleSetRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
			log.Printf("[DEBUG] Linux %s - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Linux %s: %+v", id, err)
	}

	d.Set("name", id.VirtualMachineScaleSetName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))
		d.Set("edge_zone", flattenEdgeZone(model.ExtendedLocation))
		d.Set("zones", zones.FlattenUntyped(model.Zones))
		var skuName *string
		var instances int
		if model.Sku != nil {
			skuName = model.Sku.Name
			if model.Sku.Capacity != nil {
				instances = int(*model.Sku.Capacity)
			}
		}
		d.Set("instances", instances)
		d.Set("sku", skuName)

		identityFlattened, err := identity.FlattenSystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return err
		}
		if err := d.Set("identity", identityFlattened); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if err := d.Set("plan", flattenPlanVMSS(model.Plan)); err != nil {
			return fmt.Errorf("setting `plan`: %+v", err)
		}

		if props := model.Properties; props != nil {
			if err := d.Set("additional_capabilities", FlattenVirtualMachineScaleSetAdditionalCapabilities(props.AdditionalCapabilities)); err != nil {
				return fmt.Errorf("setting `additional_capabilities`: %+v", props.AdditionalCapabilities)
			}

			if err := d.Set("automatic_instance_repair", FlattenVirtualMachineScaleSetAutomaticRepairsPolicy(props.AutomaticRepairsPolicy)); err != nil {
				return fmt.Errorf("setting `automatic_instance_repair`: %+v", err)
			}

			d.Set("do_not_run_extensions_on_overprovisioned_machines", props.DoNotRunExtensionsOnOverprovisionedVMs)
			if props.HostGroup != nil && props.HostGroup.Id != nil {
				d.Set("host_group_id", props.HostGroup.Id)
			}
			d.Set("overprovision", props.Overprovision)
			proximityPlacementGroupId := ""
			if props.ProximityPlacementGroup != nil && props.ProximityPlacementGroup.Id != nil {
				proximityPlacementGroupId = *props.ProximityPlacementGroup.Id
			}
			d.Set("platform_fault_domain_count", props.PlatformFaultDomainCount)
			d.Set("proximity_placement_group_id", proximityPlacementGroupId)
			d.Set("single_placement_group", props.SinglePlacementGroup)
			d.Set("unique_id", props.UniqueId)
			d.Set("zone_balance", props.ZoneBalance)
			d.Set("scale_in", FlattenVirtualMachineScaleSetScaleInPolicy(props.ScaleInPolicy))

			if props.SpotRestorePolicy != nil {
				d.Set("spot_restore", FlattenVirtualMachineScaleSetSpotRestorePolicy(props.SpotRestorePolicy))
			}

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

				d.Set("eviction_policy", string(pointer.From(profile.EvictionPolicy)))

				if profile.ApplicationProfile != nil && profile.ApplicationProfile.GalleryApplications != nil {
					d.Set("gallery_application", flattenVirtualMachineScaleSetGalleryApplication(profile.ApplicationProfile.GalleryApplications))
				}

				// the service just return empty when this is not assigned when provisioned
				// See discussion on https://github.com/Azure/azure-rest-api-specs/issues/10971
				priority := virtualmachinescalesets.VirtualMachinePriorityTypesRegular
				if pointer.From(profile.Priority) != "" {
					priority = pointer.From(profile.Priority)
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

						flattenedSshKeys, err := flattenSSHKeysVMSS(linux.Ssh)
						if err != nil {
							return fmt.Errorf("flattening `admin_ssh_key`: %+v", err)
						}
						if err := d.Set("admin_ssh_key", pluginsdk.NewSet(SSHKeySchemaHash, *flattenedSshKeys)); err != nil {
							return fmt.Errorf("setting `admin_ssh_key`: %+v", err)
						}
					}

					if err := d.Set("secret", flattenLinuxSecretsVMSS(osProfile.Secrets)); err != nil {
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
					if nwProfile.HealthProbe != nil && nwProfile.HealthProbe.Id != nil {
						healthProbeId = *nwProfile.HealthProbe.Id
					}
					d.Set("health_probe_id", healthProbeId)
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
				d.Set("upgrade_mode", string(pointer.From(policy.Mode)))

				flattenedAutomatic := FlattenVirtualMachineScaleSetAutomaticOSUpgradePolicy(policy.AutomaticOSUpgradePolicy)
				if err := d.Set("automatic_os_upgrade_policy", flattenedAutomatic); err != nil {
					return fmt.Errorf("setting `automatic_os_upgrade_policy`: %+v", err)
				}

				flattenedRolling := FlattenVirtualMachineScaleSetRollingUpgradePolicy(policy.RollingUpgradePolicy)
				if err := d.Set("rolling_upgrade_policy", flattenedRolling); err != nil {
					return fmt.Errorf("setting `rolling_upgrade_policy`: %+v", err)
				}
			}
		}
		return tags.FlattenAndSet(d, model.Tags)
	}
	return nil
}

func resourceLinuxVirtualMachineScaleSetDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VirtualMachineScaleSetsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualmachinescalesets.ParseVirtualMachineScaleSetID(d.Id())
	if err != nil {
		return err
	}

	// Upgrading to the 2021-07-01 exposed a new expand parameter to the GET method
	resp, err := client.Get(ctx, *id, virtualmachinescalesets.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil
		}

		return fmt.Errorf("retrieving Linux %s: %+v", id, err)
	}

	if resp.Model == nil {
		return fmt.Errorf("model was nil for %s", id)
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
	if scaleToZeroOnDelete && resp.Model.Sku != nil {
		resp.Model.Sku.Capacity = utils.Int64(int64(0))

		log.Printf("[DEBUG] Scaling instances to 0 prior to deletion - this helps avoids networking issues within Azure")
		update := virtualmachinescalesets.VirtualMachineScaleSetUpdate{
			Sku: resp.Model.Sku,
		}
		if err := client.UpdateThenPoll(ctx, *id, update, virtualmachinescalesets.DefaultUpdateOperationOptions()); err != nil {
			return fmt.Errorf("updating number of instances in %s to scale to 0: %+v", id, err)
		}
		log.Printf("[DEBUG] Scaled instances to 0 prior to deletion - this helps avoids networking issues within Azure")
	} else {
		log.Printf("[DEBUG] Unable to scale instances to `0` since the `sku` block is nil - trying to delete anyway")
	}

	log.Printf("[DEBUG] Deleting Linux %s", id)
	// @ArcturusZhang (mimicking from linux_virtual_machine_pluginsdk.go): sending `nil` here omits this value from being sent
	// which matches the previous behaviour - we're only splitting this out so it's clear why
	// TODO: support force deletion once it's out of Preview, if applicable
	if err := client.DeleteThenPoll(ctx, *id, virtualmachinescalesets.DefaultDeleteOperationOptions()); err != nil {
		return fmt.Errorf("deleting Linux %s: %+v", id, err)
	}

	log.Printf("[DEBUG] Deleted Linux %s", id)

	return nil
}

func resourceLinuxVirtualMachineScaleSetSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
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
			Default:  string(virtualmachinescalesets.VirtualMachinePriorityTypesRegular),
			ValidateFunc: validation.StringInSlice([]string{
				string(virtualmachinescalesets.VirtualMachinePriorityTypesRegular),
				string(virtualmachinescalesets.VirtualMachinePriorityTypesSpot),
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

		"tags": commonschema.Tags(),

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

		"zones": commonschema.ZonesMultipleOptional(),

		// Computed
		"unique_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
