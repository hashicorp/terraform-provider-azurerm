// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleryapplicationversions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-04-02/disks"
	azValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/compute/2023-03-01/compute"
)

func virtualMachineAdditionalCapabilitiesSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				// TODO: confirm this command

				// NOTE: requires registration to use:
				// $ az feature show --namespace Microsoft.Compute --name UltraSSDWithVMSS
				// $ az provider register -n Microsoft.Compute
				"ultra_ssd_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},
			},
		},
	}
}

func expandVirtualMachineAdditionalCapabilities(input []interface{}) *compute.AdditionalCapabilities {
	capabilities := compute.AdditionalCapabilities{}

	if len(input) > 0 {
		raw := input[0].(map[string]interface{})

		capabilities.UltraSSDEnabled = utils.Bool(raw["ultra_ssd_enabled"].(bool))
	}

	return &capabilities
}

func flattenVirtualMachineAdditionalCapabilities(input *compute.AdditionalCapabilities) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	ultraSsdEnabled := false

	if input.UltraSSDEnabled != nil {
		ultraSsdEnabled = *input.UltraSSDEnabled
	}

	return []interface{}{
		map[string]interface{}{
			"ultra_ssd_enabled": ultraSsdEnabled,
		},
	}
}

func expandVirtualMachineIdentity(input []interface{}) (*compute.VirtualMachineIdentity, error) {
	expanded, err := identity.ExpandSystemAndUserAssignedMap(input)
	if err != nil {
		return nil, err
	}

	out := compute.VirtualMachineIdentity{
		Type: compute.ResourceIdentityType(string(expanded.Type)),
	}
	if expanded.Type == identity.TypeUserAssigned || expanded.Type == identity.TypeSystemAssignedUserAssigned {
		out.UserAssignedIdentities = make(map[string]*compute.UserAssignedIdentitiesValue)
		for k := range expanded.IdentityIds {
			out.UserAssignedIdentities[k] = &compute.UserAssignedIdentitiesValue{
				// intentionally empty
			}
		}
	}
	return &out, nil
}

func flattenVirtualMachineIdentity(input *compute.VirtualMachineIdentity) (*[]interface{}, error) {
	var transform *identity.SystemAndUserAssignedMap

	if input != nil {
		transform = &identity.SystemAndUserAssignedMap{
			Type:        identity.Type(string(input.Type)),
			IdentityIds: make(map[string]identity.UserAssignedIdentityDetails),
		}

		if input.PrincipalID != nil {
			transform.PrincipalId = *input.PrincipalID
		}
		if input.TenantID != nil {
			transform.TenantId = *input.TenantID
		}
		for k, v := range input.UserAssignedIdentities {
			transform.IdentityIds[k] = identity.UserAssignedIdentityDetails{
				ClientId:    v.ClientID,
				PrincipalId: v.PrincipalID,
			}
		}
	}
	return identity.FlattenSystemAndUserAssignedMap(transform)
}

func expandVirtualMachineNetworkInterfaceIDs(input []interface{}) []compute.NetworkInterfaceReference {
	output := make([]compute.NetworkInterfaceReference, 0)

	for i, v := range input {
		output = append(output, compute.NetworkInterfaceReference{
			ID: utils.String(v.(string)),
			NetworkInterfaceReferenceProperties: &compute.NetworkInterfaceReferenceProperties{
				Primary: utils.Bool(i == 0),
			},
		})
	}

	return output
}

func flattenVirtualMachineNetworkInterfaceIDs(input *[]compute.NetworkInterfaceReference) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)

	for _, v := range *input {
		if v.ID == nil {
			continue
		}

		output = append(output, *v.ID)
	}

	return output
}

func virtualMachineOSDiskSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"caching": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(compute.CachingTypesNone),
						string(compute.CachingTypesReadOnly),
						string(compute.CachingTypesReadWrite),
					}, false),
				},
				"storage_account_type": {
					Type:     pluginsdk.TypeString,
					Required: true,
					// whilst this appears in the Update block the API returns this when changing:
					// Changing property 'osDisk.managedDisk.storageAccountType' is not allowed
					ForceNew: true,
					ValidateFunc: validation.StringInSlice([]string{
						// note: OS Disks don't support Ultra SSDs or PremiumV2_LRS
						string(compute.StorageAccountTypesPremiumLRS),
						string(compute.StorageAccountTypesStandardLRS),
						string(compute.StorageAccountTypesStandardSSDLRS),
						string(compute.StorageAccountTypesStandardSSDZRS),
						string(compute.StorageAccountTypesPremiumZRS),
					}, false),
				},

				// Optional
				"diff_disk_settings": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					ForceNew: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"option": {
								Type:     pluginsdk.TypeString,
								Required: true,
								ForceNew: true,
								ValidateFunc: validation.StringInSlice([]string{
									string(compute.DiffDiskOptionsLocal),
								}, false),
							},
							"placement": {
								Type:     pluginsdk.TypeString,
								Optional: true,
								ForceNew: true,
								Default:  string(compute.DiffDiskPlacementCacheDisk),
								ValidateFunc: validation.StringInSlice([]string{
									string(compute.DiffDiskPlacementCacheDisk),
									string(compute.DiffDiskPlacementResourceDisk),
								}, false),
							},
						},
					},
				},

				"disk_encryption_set_id": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					// the Compute/VM API is broken and returns the Resource Group name in UPPERCASE
					DiffSuppressFunc: suppress.CaseDifference,
					ValidateFunc:     validate.DiskEncryptionSetID,
					ConflictsWith:    []string{"os_disk.0.secure_vm_disk_encryption_set_id"},
				},

				"disk_size_gb": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.IntBetween(0, 4095),
				},

				"name": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ForceNew: true,
					Computed: true,
				},

				"secure_vm_disk_encryption_set_id": {
					Type:          pluginsdk.TypeString,
					Optional:      true,
					ForceNew:      true,
					ValidateFunc:  validate.DiskEncryptionSetID,
					ConflictsWith: []string{"os_disk.0.disk_encryption_set_id"},
				},

				"security_encryption_type": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ForceNew: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(compute.SecurityEncryptionTypesVMGuestStateOnly),
						string(compute.SecurityEncryptionTypesDiskWithVMGuestState),
					}, false),
				},

				"write_accelerator_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},
			},
		},
	}
}

func expandVirtualMachineOSDisk(input []interface{}, osType compute.OperatingSystemTypes) (*compute.OSDisk, error) {
	raw := input[0].(map[string]interface{})
	caching := raw["caching"].(string)
	disk := compute.OSDisk{
		Caching: compute.CachingTypes(caching),
		ManagedDisk: &compute.ManagedDiskParameters{
			StorageAccountType: compute.StorageAccountTypes(raw["storage_account_type"].(string)),
		},
		WriteAcceleratorEnabled: utils.Bool(raw["write_accelerator_enabled"].(bool)),

		// these have to be hard-coded so there's no point exposing them
		// for CreateOption, whilst it's possible for this to be "Attach" for an OS Disk
		// from what we can tell this approach has been superseded by provisioning from
		// an image of the machine (e.g. an Image/Shared Image Gallery)
		CreateOption: compute.DiskCreateOptionTypesFromImage,
		OsType:       osType,
	}

	securityEncryptionType := raw["security_encryption_type"].(string)
	if securityEncryptionType != "" {
		disk.ManagedDisk.SecurityProfile = &compute.VMDiskSecurityProfile{
			SecurityEncryptionType: compute.SecurityEncryptionTypes(securityEncryptionType),
		}
	}
	if secureVMDiskEncryptionId := raw["secure_vm_disk_encryption_set_id"].(string); secureVMDiskEncryptionId != "" {
		if compute.SecurityEncryptionTypesDiskWithVMGuestState != compute.SecurityEncryptionTypes(securityEncryptionType) {
			return nil, fmt.Errorf("`secure_vm_disk_encryption_set_id` can only be specified when `security_encryption_type` is set to `DiskWithVMGuestState`")
		}
		disk.ManagedDisk.SecurityProfile.DiskEncryptionSet = &compute.DiskEncryptionSetParameters{
			ID: utils.String(secureVMDiskEncryptionId),
		}
	}

	if osDiskSize := raw["disk_size_gb"].(int); osDiskSize > 0 {
		disk.DiskSizeGB = utils.Int32(int32(osDiskSize))
	}

	if diffDiskSettingsRaw := raw["diff_disk_settings"].([]interface{}); len(diffDiskSettingsRaw) > 0 {
		if caching != string(compute.CachingTypesReadOnly) {
			// Restriction per https://docs.microsoft.com/azure/virtual-machines/ephemeral-os-disks-deploy#vm-template-deployment
			return nil, fmt.Errorf("`diff_disk_settings` can only be set when `caching` is set to `ReadOnly`")
		}

		diffDiskRaw := diffDiskSettingsRaw[0].(map[string]interface{})
		disk.DiffDiskSettings = &compute.DiffDiskSettings{
			Option:    compute.DiffDiskOptions(diffDiskRaw["option"].(string)),
			Placement: compute.DiffDiskPlacement(diffDiskRaw["placement"].(string)),
		}
	}

	if id := raw["disk_encryption_set_id"].(string); id != "" {
		disk.ManagedDisk.DiskEncryptionSet = &compute.DiskEncryptionSetParameters{
			ID: utils.String(id),
		}
	}

	if name := raw["name"].(string); name != "" {
		disk.Name = utils.String(name)
	}

	return &disk, nil
}

func flattenVirtualMachineOSDisk(ctx context.Context, disksClient *disks.DisksClient, input *compute.OSDisk) ([]interface{}, error) {
	if input == nil {
		return []interface{}{}, nil
	}

	diffDiskSettings := make([]interface{}, 0)
	if input.DiffDiskSettings != nil {
		placement := string(compute.DiffDiskPlacementCacheDisk)
		if input.DiffDiskSettings.Placement != "" {
			placement = string(input.DiffDiskSettings.Placement)
		}

		diffDiskSettings = append(diffDiskSettings, map[string]interface{}{
			"option":    string(input.DiffDiskSettings.Option),
			"placement": placement,
		})
	}

	diskSizeGb := 0
	if input.DiskSizeGB != nil && *input.DiskSizeGB != 0 {
		diskSizeGb = int(*input.DiskSizeGB)
	}

	var name string
	if input.Name != nil {
		name = *input.Name
	}

	diskEncryptionSetId := ""
	storageAccountType := ""
	secureVMDiskEncryptionSetId := ""
	securityEncryptionType := ""

	if input.ManagedDisk != nil {
		storageAccountType = string(input.ManagedDisk.StorageAccountType)

		if input.ManagedDisk.ID != nil {
			id, err := commonids.ParseManagedDiskID(*input.ManagedDisk.ID)
			if err != nil {
				return nil, err
			}

			disk, err := disksClient.Get(ctx, *id)
			if err != nil {
				// turns out ephemeral disks aren't returned/available here
				if !response.WasNotFound(disk.HttpResponse) {
					return nil, err
				}
			}

			// Ephemeral Disks get an ARM ID but aren't available via the regular API
			// ergo fingers crossed we've got it from the resource because ¯\_(ツ)_/¯
			// where else we'd be able to pull it from
			if !response.WasNotFound(disk.HttpResponse) {
				// whilst this is available as `input.ManagedDisk.StorageAccountType` it's not returned there
				// however it's only available there for ephemeral os disks
				if disk.Model.Sku != nil && storageAccountType == "" {
					storageAccountType = string(*disk.Model.Sku.Name)
				}

				// same goes for Disk Size GB apparently
				if diskSizeGb == 0 && disk.Model.Properties != nil && disk.Model.Properties.DiskSizeGB != nil {
					diskSizeGb = int(*disk.Model.Properties.DiskSizeGB)
				}

				// same goes for Disk Encryption Set Id apparently
				if disk.Model.Properties.Encryption != nil && disk.Model.Properties.Encryption.DiskEncryptionSetId != nil {
					diskEncryptionSetId = *disk.Model.Properties.Encryption.DiskEncryptionSetId
				}
			}
		}

		if securityProfile := input.ManagedDisk.SecurityProfile; securityProfile != nil {
			securityEncryptionType = string(securityProfile.SecurityEncryptionType)
			if securityProfile.DiskEncryptionSet != nil && securityProfile.DiskEncryptionSet.ID != nil {
				secureVMDiskEncryptionSetId = *securityProfile.DiskEncryptionSet.ID
			}
		}
	}

	writeAcceleratorEnabled := false
	if input.WriteAcceleratorEnabled != nil {
		writeAcceleratorEnabled = *input.WriteAcceleratorEnabled
	}
	return []interface{}{
		map[string]interface{}{
			"caching":                          string(input.Caching),
			"disk_size_gb":                     diskSizeGb,
			"diff_disk_settings":               diffDiskSettings,
			"disk_encryption_set_id":           diskEncryptionSetId,
			"name":                             name,
			"storage_account_type":             storageAccountType,
			"secure_vm_disk_encryption_set_id": secureVMDiskEncryptionSetId,
			"security_encryption_type":         securityEncryptionType,
			"write_accelerator_enabled":        writeAcceleratorEnabled,
		},
	}, nil
}

func virtualMachineOsImageNotificationSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"timeout": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						"PT15M",
					}, false),
					Default: "PT15M",
				},
			},
		},
	}
}

func virtualMachineTerminationNotificationSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"enabled": {
					Type:     pluginsdk.TypeBool,
					Required: true,
				},
				"timeout": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: azValidate.ISO8601DurationBetween("PT5M", "PT15M"),
					Default:      "PT5M",
				},
			},
		},
	}
}

func expandOsImageNotificationProfile(input []interface{}) *compute.OSImageNotificationProfile {
	if len(input) == 0 {
		return &compute.OSImageNotificationProfile{
			Enable: utils.Bool(false),
		}
	}

	raw := input[0].(map[string]interface{})
	timeout := raw["timeout"].(string)

	return &compute.OSImageNotificationProfile{
		Enable:           utils.Bool(true),
		NotBeforeTimeout: &timeout,
	}
}

func expandTerminateNotificationProfile(input []interface{}) *compute.TerminateNotificationProfile {
	if len(input) == 0 {
		return &compute.TerminateNotificationProfile{
			Enable: utils.Bool(false),
		}
	}

	raw := input[0].(map[string]interface{})
	enabled := raw["enabled"].(bool)
	timeout := raw["timeout"].(string)

	return &compute.TerminateNotificationProfile{
		Enable:           &enabled,
		NotBeforeTimeout: &timeout,
	}
}

func flattenOsImageNotificationProfile(input *compute.OSImageNotificationProfile) []interface{} {
	if input == nil || !pointer.From(input.Enable) {
		return nil
	}

	timeout := "PT15M"
	if input.NotBeforeTimeout != nil {
		timeout = *input.NotBeforeTimeout
	}

	return []interface{}{
		map[string]interface{}{
			"timeout": timeout,
		},
	}
}

func flattenTerminateNotificationProfile(input *compute.TerminateNotificationProfile) []interface{} {
	// if enabled is set to false, there will be no ScheduledEventsProfile in response, to avoid plan non empty when
	// a user explicitly set enabled to false, we need to assign a default block to this field

	enabled := false
	if input != nil && input.Enable != nil {
		enabled = *input.Enable
	}

	timeout := "PT5M"
	if input != nil && input.NotBeforeTimeout != nil {
		timeout = *input.NotBeforeTimeout
	}

	return []interface{}{
		map[string]interface{}{
			"enabled": enabled,
			"timeout": timeout,
		},
	}
}

func VirtualMachineGalleryApplicationSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 100,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"version_id": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: galleryapplicationversions.ValidateApplicationVersionID,
				},

				"automatic_upgrade_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				// Example: https://mystorageaccount.blob.core.windows.net/configurations/settings.config
				"configuration_blob_uri": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.IsURLWithHTTPorHTTPS,
				},

				"order": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Default:      0,
					ValidateFunc: validation.IntBetween(0, 2147483647),
				},

				// NOTE: Per the service team, "this is a pass through value that we just add to the model but don't depend on. It can be any string."
				"tag": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"treat_failure_as_deployment_failure_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},
			},
		},
	}
}

func expandVirtualMachineGalleryApplication(input []interface{}) *[]compute.VMGalleryApplication {
	out := make([]compute.VMGalleryApplication, 0)
	if len(input) == 0 {
		return &out
	}

	for _, v := range input {
		packageReferenceId := v.(map[string]interface{})["version_id"].(string)
		automaticUpgradeEnabled := v.(map[string]interface{})["automatic_upgrade_enabled"].(bool)
		configurationReference := v.(map[string]interface{})["configuration_blob_uri"].(string)
		order := v.(map[string]interface{})["order"].(int)
		tag := v.(map[string]interface{})["tag"].(string)
		treatFailureAsDeploymentFailureEnabled := v.(map[string]interface{})["treat_failure_as_deployment_failure_enabled"].(bool)

		app := &compute.VMGalleryApplication{
			PackageReferenceID:              utils.String(packageReferenceId),
			ConfigurationReference:          utils.String(configurationReference),
			Order:                           utils.Int32(int32(order)),
			Tags:                            utils.String(tag),
			EnableAutomaticUpgrade:          utils.Bool(automaticUpgradeEnabled),
			TreatFailureAsDeploymentFailure: utils.Bool(treatFailureAsDeploymentFailureEnabled),
		}

		out = append(out, *app)
	}

	return &out
}

func flattenVirtualMachineGalleryApplication(input *[]compute.VMGalleryApplication) []interface{} {
	if len(*input) == 0 {
		return nil
	}

	out := make([]interface{}, 0)

	for _, v := range *input {
		var packageReferenceId, configurationReference, tag string
		var order int
		var automaticUpgradeEnabled, treatFailureAsDeploymentFailureEnabled bool

		if v.PackageReferenceID != nil {
			packageReferenceId = *v.PackageReferenceID
		}

		if v.ConfigurationReference != nil {
			configurationReference = *v.ConfigurationReference
		}

		if v.EnableAutomaticUpgrade != nil {
			automaticUpgradeEnabled = *v.EnableAutomaticUpgrade
		}

		if v.Order != nil {
			order = int(*v.Order)
		}

		if v.Tags != nil {
			tag = *v.Tags
		}

		if v.TreatFailureAsDeploymentFailure != nil {
			treatFailureAsDeploymentFailureEnabled = *v.TreatFailureAsDeploymentFailure
		}

		app := map[string]interface{}{
			"version_id":                packageReferenceId,
			"automatic_upgrade_enabled": automaticUpgradeEnabled,
			"configuration_blob_uri":    configurationReference,
			"order":                     order,
			"tag":                       tag,
			"treat_failure_as_deployment_failure_enabled": treatFailureAsDeploymentFailureEnabled,
		}

		out = append(out, app)
	}

	return out
}
