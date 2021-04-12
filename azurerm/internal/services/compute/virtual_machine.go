package compute

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-12-01/compute"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/validate"
	msiparse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/msi/parse"
	msivalidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/msi/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func virtualMachineAdditionalCapabilitiesSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				// TODO: confirm this command

				// NOTE: requires registration to use:
				// $ az feature show --namespace Microsoft.Compute --name UltraSSDWithVMSS
				// $ az provider register -n Microsoft.Compute
				"ultra_ssd_enabled": {
					Type:     schema.TypeBool,
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

func virtualMachineIdentitySchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Type:     schema.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(compute.ResourceIdentityTypeSystemAssigned),
						string(compute.ResourceIdentityTypeUserAssigned),
						string(compute.ResourceIdentityTypeSystemAssignedUserAssigned),
					}, false),
				},

				"identity_ids": {
					Type:     schema.TypeSet,
					Optional: true,
					Elem: &schema.Schema{
						Type:         schema.TypeString,
						ValidateFunc: msivalidate.UserAssignedIdentityID,
					},
				},

				"principal_id": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"tenant_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func expandVirtualMachineIdentity(input []interface{}) (*compute.VirtualMachineIdentity, error) {
	if len(input) == 0 {
		// TODO: Does this want to be this, or nil?
		return &compute.VirtualMachineIdentity{
			Type: compute.ResourceIdentityTypeNone,
		}, nil
	}

	raw := input[0].(map[string]interface{})

	identity := compute.VirtualMachineIdentity{
		Type: compute.ResourceIdentityType(raw["type"].(string)),
	}

	identityIdsRaw := raw["identity_ids"].(*schema.Set).List()
	identityIds := make(map[string]*compute.VirtualMachineIdentityUserAssignedIdentitiesValue)
	for _, v := range identityIdsRaw {
		identityIds[v.(string)] = &compute.VirtualMachineIdentityUserAssignedIdentitiesValue{}
	}

	if len(identityIds) > 0 {
		if identity.Type != compute.ResourceIdentityTypeUserAssigned && identity.Type != compute.ResourceIdentityTypeSystemAssignedUserAssigned {
			return nil, fmt.Errorf("`identity_ids` can only be specified when `type` includes `UserAssigned`")
		}

		identity.UserAssignedIdentities = identityIds
	}

	return &identity, nil
}

func flattenVirtualMachineIdentity(input *compute.VirtualMachineIdentity) ([]interface{}, error) {
	if input == nil || input.Type == compute.ResourceIdentityTypeNone {
		return []interface{}{}, nil
	}

	identityIds := make([]string, 0)
	if input.UserAssignedIdentities != nil {
		for key := range input.UserAssignedIdentities {
			parsedId, err := msiparse.UserAssignedIdentityIDInsensitively(key)
			if err != nil {
				return nil, err
			}
			identityIds = append(identityIds, parsedId.ID())
		}
	}

	principalId := ""
	if input.PrincipalID != nil {
		principalId = *input.PrincipalID
	}

	tenantId := ""
	if input.TenantID != nil {
		tenantId = *input.TenantID
	}

	return []interface{}{
		map[string]interface{}{
			"type":         string(input.Type),
			"identity_ids": identityIds,
			"principal_id": principalId,
			"tenant_id":    tenantId,
		},
	}, nil
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

func virtualMachineOSDiskSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"caching": {
					Type:     schema.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(compute.CachingTypesNone),
						string(compute.CachingTypesReadOnly),
						string(compute.CachingTypesReadWrite),
					}, false),
				},
				"storage_account_type": {
					Type:     schema.TypeString,
					Required: true,
					// whilst this appears in the Update block the API returns this when changing:
					// Changing property 'osDisk.managedDisk.storageAccountType' is not allowed
					ForceNew: true,
					ValidateFunc: validation.StringInSlice([]string{
						// note: OS Disks don't support Ultra SSDs
						string(compute.StorageAccountTypesPremiumLRS),
						string(compute.StorageAccountTypesStandardLRS),
						string(compute.StorageAccountTypesStandardSSDLRS),
					}, false),
				},

				// Optional
				"diff_disk_settings": {
					Type:     schema.TypeList,
					Optional: true,
					ForceNew: true,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"option": {
								Type:     schema.TypeString,
								Required: true,
								ForceNew: true,
								ValidateFunc: validation.StringInSlice([]string{
									string(compute.Local),
								}, false),
							},
						},
					},
				},

				"disk_encryption_set_id": {
					Type:     schema.TypeString,
					Optional: true,
					// the Compute/VM API is broken and returns the Resource Group name in UPPERCASE
					DiffSuppressFunc: suppress.CaseDifference,
					ValidateFunc:     validate.DiskEncryptionSetID,
				},

				"disk_size_gb": {
					Type:         schema.TypeInt,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.IntBetween(0, 4095),
				},

				"name": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
					Computed: true,
				},

				"write_accelerator_enabled": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
			},
		},
	}
}

func expandVirtualMachineOSDisk(input []interface{}, osType compute.OperatingSystemTypes) *compute.OSDisk {
	raw := input[0].(map[string]interface{})
	disk := compute.OSDisk{
		Caching: compute.CachingTypes(raw["caching"].(string)),
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

	if osDiskSize := raw["disk_size_gb"].(int); osDiskSize > 0 {
		disk.DiskSizeGB = utils.Int32(int32(osDiskSize))
	}

	if diffDiskSettingsRaw := raw["diff_disk_settings"].([]interface{}); len(diffDiskSettingsRaw) > 0 {
		diffDiskRaw := diffDiskSettingsRaw[0].(map[string]interface{})
		disk.DiffDiskSettings = &compute.DiffDiskSettings{
			Option: compute.DiffDiskOptions(diffDiskRaw["option"].(string)),
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

	return &disk
}

func flattenVirtualMachineOSDisk(ctx context.Context, disksClient *compute.DisksClient, input *compute.OSDisk) ([]interface{}, error) {
	if input == nil {
		return []interface{}{}, nil
	}

	diffDiskSettings := make([]interface{}, 0)
	if input.DiffDiskSettings != nil {
		diffDiskSettings = append(diffDiskSettings, map[string]interface{}{
			"option": string(input.DiffDiskSettings.Option),
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

	if input.ManagedDisk != nil {
		storageAccountType = string(input.ManagedDisk.StorageAccountType)

		if input.ManagedDisk.ID != nil {
			id, err := parse.ManagedDiskID(*input.ManagedDisk.ID)
			if err != nil {
				return nil, err
			}

			disk, err := disksClient.Get(ctx, id.ResourceGroup, id.DiskName)
			if err != nil {
				// turns out ephemeral disks aren't returned/available here
				if !utils.ResponseWasNotFound(disk.Response) {
					return nil, err
				}
			}

			// Ephemeral Disks get an ARM ID but aren't available via the regular API
			// ergo fingers crossed we've got it from the resource because ¯\_(ツ)_/¯
			// where else we'd be able to pull it from
			if !utils.ResponseWasNotFound(disk.Response) {
				// whilst this is available as `input.ManagedDisk.StorageAccountType` it's not returned there
				// however it's only available there for ephemeral os disks
				if disk.Sku != nil && storageAccountType == "" {
					storageAccountType = string(disk.Sku.Name)
				}

				// same goes for Disk Size GB apparently
				if diskSizeGb == 0 && disk.DiskProperties != nil && disk.DiskProperties.DiskSizeGB != nil {
					diskSizeGb = int(*disk.DiskProperties.DiskSizeGB)
				}

				// same goes for Disk Encryption Set Id apparently
				if disk.Encryption != nil && disk.Encryption.DiskEncryptionSetID != nil {
					diskEncryptionSetId = *disk.Encryption.DiskEncryptionSetID
				}
			}
		}
	}

	writeAcceleratorEnabled := false
	if input.WriteAcceleratorEnabled != nil {
		writeAcceleratorEnabled = *input.WriteAcceleratorEnabled
	}
	return []interface{}{
		map[string]interface{}{
			"caching":                   string(input.Caching),
			"disk_size_gb":              diskSizeGb,
			"diff_disk_settings":        diffDiskSettings,
			"disk_encryption_set_id":    diskEncryptionSetId,
			"name":                      name,
			"storage_account_type":      storageAccountType,
			"write_accelerator_enabled": writeAcceleratorEnabled,
		},
	}, nil
}

func virtualMachineDataDiskSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		MaxItems: 1,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"create": virtualMachineCreateDataDiskSchema(),

				"attach": virtualMachineAttachDataDiskSchema(),
			},
		},
	}
}

func virtualMachineCreateDataDiskSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"caching": {
					Type:     schema.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(compute.CachingTypesNone),
						string(compute.CachingTypesReadOnly),
						string(compute.CachingTypesReadWrite),
					}, false),
				},

				"lun": {
					Type:         schema.TypeInt,
					Required:     true,
					ValidateFunc: validation.IntAtLeast(0),
				},

				"name": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"storage_account_type": {
					Type:     schema.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(compute.StorageAccountTypesPremiumLRS),
						string(compute.StorageAccountTypesStandardLRS),
						string(compute.StorageAccountTypesStandardSSDLRS),
						string(compute.StorageAccountTypesUltraSSDLRS),
					}, false),
				},

				"disk_size_gb": {
					Type:     schema.TypeInt,
					Required: true,
					// Max based on P80/S80/E80 Managed Disk type max size
					ValidateFunc: validation.IntBetween(1, 32767),
				},

				"disk_encryption_set_id": {
					Type:     schema.TypeString,
					Optional: true,
					// the Compute/VM API is broken and returns the Resource Group name in UPPERCASE
					DiffSuppressFunc: suppress.CaseDifference,
					ValidateFunc:     validate.DiskEncryptionSetID,
				},

				"write_accelerator_enabled": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
			},
		},
		Set: resourceVirtualMachineCreateDataDiskHash,
	}
}

func virtualMachineAttachDataDiskSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"managed_disk_id": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validate.ManagedDiskID,
				},

				"lun": {
					Type:         schema.TypeInt,
					Required:     true,
					ValidateFunc: validation.IntAtLeast(0),
				},

				"caching": {
					Type:     schema.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(compute.CachingTypesNone),
						string(compute.CachingTypesReadOnly),
						string(compute.CachingTypesReadWrite),
					}, false),
				},

				"storage_account_type": {
					Type:     schema.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(compute.StorageAccountTypesPremiumLRS),
						string(compute.StorageAccountTypesStandardLRS),
						string(compute.StorageAccountTypesStandardSSDLRS),
						string(compute.StorageAccountTypesUltraSSDLRS),
					}, false),
				},

				"write_accelerator_enabled": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
			},
		},
	}
}

func expandVirtualMachineDataDisks(ctx context.Context, d *schema.ResourceData, meta interface{}) (*[]compute.DataDisk, error) {
	input, ok := d.GetOk("data_disks")
	if !ok {
		return &[]compute.DataDisk{}, nil
	}

	var err error
	dataDisksRaw := input.([]interface{})

	result := make([]compute.DataDisk, 0)
	dataDisks := dataDisksRaw[0].(map[string]interface{})
	if newDisksRaw, ok := dataDisks["create"]; ok {
		var newDisks []compute.DataDisk
		if d.IsNewResource() {
			newDisks = expandVirtualMachineCreateDataDisksForCreate(newDisksRaw)
		} else {
			newDisks, err = expandVirtualMachineCreateDataDisksForUpdate(ctx, newDisksRaw, d, meta)
			if err != nil {
				return nil, err
			}
		}
		result = append(result, newDisks...)
	}

	if attachDisksRaw, ok := dataDisks["attach"]; ok {
		attachDisks, err := expandVirtualMachineAttachDataDisksForCreate(ctx, attachDisksRaw, d, meta)
		if err != nil {
			return nil, err
		}

		result = append(result, attachDisks...)
	}

	return &result, nil
}

func expandVirtualMachineCreateDataDisksForCreate(input interface{}) []compute.DataDisk {
	var output []compute.DataDisk
	if input == nil || len(input.(*schema.Set).List()) == 0 {
		return output
	}

	for _, v := range input.(*schema.Set).List() {
		disk := v.(map[string]interface{})
		dataDisk := compute.DataDisk{
			Caching:      compute.CachingTypes(disk["caching"].(string)),
			CreateOption: compute.DiskCreateOptionTypesEmpty,
			DiskSizeGB:   utils.Int32(int32(disk["disk_size_gb"].(int))),
			Lun:          utils.Int32(int32(disk["lun"].(int))),
			Name:         utils.String(disk["name"].(string)),
			ManagedDisk: &compute.ManagedDiskParameters{
				StorageAccountType: compute.StorageAccountTypes(disk["storage_account_type"].(string)),
			},
		}
		if diskEncryptionSetID, ok := disk["disk_encryption_set_id"].(string); ok && diskEncryptionSetID != "" {
			dataDisk.ManagedDisk.DiskEncryptionSet = &compute.DiskEncryptionSetParameters{
				ID: utils.String(diskEncryptionSetID),
			}
		}
		output = append(output, dataDisk)
	}

	return output
}

func expandVirtualMachineCreateDataDisksForUpdate(ctx context.Context, input interface{}, d *schema.ResourceData, meta interface{}) ([]compute.DataDisk, error) {
	var output []compute.DataDisk

	if input == nil || len(input.(*schema.Set).List()) == 0 {
		return output, nil
	}

	diskClient := meta.(*clients.Client).Compute.DisksClient

	resourceGroup := d.Get("resource_group_name").(string)

	for _, v := range input.(*schema.Set).List() {
		disk := v.(map[string]interface{})
		name := disk["name"].(string)
		// this filters out the ghost entry we get from the set - less than ideal
		if name == "" {
			continue
		}
		newDiskToAdd := false
		current, err := diskClient.Get(ctx, resourceGroup, name)
		if err != nil {
			// If 404, this is a new disk to add
			if utils.ResponseWasNotFound(current.Response) {
				newDiskToAdd = true
			} else {
				return nil, fmt.Errorf("failure reading Data Disk %q (resource group %q) for update: %+v", name, resourceGroup, err)
			}
		}

		if !newDiskToAdd && current.ID == nil {
			return nil, fmt.Errorf("could not determine ID for Data Disk %q (resource group %q) for update: %+v", name, resourceGroup, err)
		}

		dataDisk := compute.DataDisk{
			Caching:      compute.CachingTypes(disk["caching"].(string)),
			CreateOption: compute.DiskCreateOptionTypesEmpty,
			DiskSizeGB:   utils.Int32(int32(disk["disk_size_gb"].(int))),
			Lun:          utils.Int32(int32(disk["lun"].(int))),
			Name:         utils.String(name),
			ManagedDisk: &compute.ManagedDiskParameters{
				ID:                 current.ID,
				StorageAccountType: compute.StorageAccountTypes(disk["storage_account_type"].(string)),
			},
		}
		if diskEncryptionSetID, ok := disk["disk_encryption_set_id"]; ok && diskEncryptionSetID != "" {
			dataDisk.ManagedDisk.DiskEncryptionSet = &compute.DiskEncryptionSetParameters{
				ID: utils.String(diskEncryptionSetID.(string)),
			}
		}
		if writeAcceleratorEnabled, ok := disk["write_accelerator_enabled"]; ok {
			dataDisk.WriteAcceleratorEnabled = utils.Bool(writeAcceleratorEnabled.(bool))
		}

		output = append(output, dataDisk)
	}

	return output, nil
}

func expandVirtualMachineAttachDataDisksForCreate(ctx context.Context, input interface{}, d *schema.ResourceData, meta interface{}) ([]compute.DataDisk, error) {
	var output []compute.DataDisk
	if input == nil || len(input.(*schema.Set).List()) == 0 {
		return output, nil
	}

	disksClient := meta.(*clients.Client).Compute.DisksClient

	for _, v := range input.(*schema.Set).List() {
		disk := v.(map[string]interface{})
		diskIDRaw := disk["managed_disk_id"].(string)

		diskID, err := parse.ManagedDiskID(diskIDRaw)
		if err != nil {
			if d.IsNewResource() {
				return nil, err
			}
			continue
		}
		attach, err := disksClient.Get(ctx, diskID.ResourceGroup, diskID.DiskName)
		if err != nil {
			return nil, fmt.Errorf("failed retrieving details for attached Managed Disk %q (resource group %q: %+v", diskID.DiskName, diskID.ResourceGroup, err)
		}
		dataDisk := compute.DataDisk{}
		dataDisk.Name = &diskID.DiskName
		dataDisk.CreateOption = compute.DiskCreateOptionTypesAttach

		if attach.DiskSizeGB == nil {
			return nil, fmt.Errorf("failed reading `disk_size_gb` from attached Managed Disk %q (resource group %q)", diskID.DiskName, diskID.ResourceGroup)
		}
		dataDisk.DiskSizeGB = attach.DiskSizeGB

		dataDisk.Caching = compute.CachingTypes(disk["caching"].(string))
		dataDisk.Lun = utils.Int32(int32(disk["lun"].(int)))
		dataDisk.ManagedDisk = &compute.ManagedDiskParameters{
			StorageAccountType: compute.StorageAccountTypes(disk["storage_account_type"].(string)),
			ID:                 utils.String(diskIDRaw),
		}

		if writeAcceleratorEnabled, ok := disk["write_accelerator_enabled"]; ok {
			dataDisk.WriteAcceleratorEnabled = utils.Bool(writeAcceleratorEnabled.(bool))
		}

		output = append(output, dataDisk)
	}

	return output, nil
}

func flattenVirtualMachineDataDisks(input *[]compute.DataDisk) ([]interface{}, error) {
	if len(*input) == 0 {
		return []interface{}{}, nil
	}

	var createDisks []interface{}
	var attachDisks []interface{}
	// we need to split into new and "attach", we can use `createOption` as indicator

	for _, v := range *input {
		dataDisk := make(map[string]interface{})

		lun := 0
		if v.Lun != nil {
			lun = int(*v.Lun)
		}
		dataDisk["lun"] = lun
		dataDisk["caching"] = string(v.Caching)

		storageAccountType := ""
		encryptionSetId := ""
		managedDiskID := ""

		writeAcceleratorEnabled := false
		if v.WriteAcceleratorEnabled != nil {
			writeAcceleratorEnabled = *v.WriteAcceleratorEnabled
		}
		dataDisk["write_accelerator_enabled"] = writeAcceleratorEnabled

		createOption := v.CreateOption

		switch createOption {
		case compute.DiskCreateOptionTypesEmpty:
			name := ""
			if v.Name != nil {
				name = *v.Name
			}
			dataDisk["name"] = name

			if managedDisk := v.ManagedDisk; managedDisk != nil {
				storageAccountType = string(managedDisk.StorageAccountType)
				if managedDisk.DiskEncryptionSet != nil && managedDisk.DiskEncryptionSet.ID != nil {
					encryptionSetId = *managedDisk.DiskEncryptionSet.ID
				}
			}
			dataDisk["storage_account_type"] = storageAccountType
			dataDisk["disk_encryption_set_id"] = encryptionSetId

			diskSizeGB := 0
			if v.DiskSizeGB != nil {
				diskSizeGB = int(*v.DiskSizeGB)
			}
			dataDisk["disk_size_gb"] = diskSizeGB

			createDisks = append(createDisks, dataDisk)

		case compute.DiskCreateOptionTypesAttach:
			if managedDisk := v.ManagedDisk; managedDisk != nil {
				storageAccountType = string(managedDisk.StorageAccountType)
				if managedDisk.ID != nil {
					managedDiskID = *managedDisk.ID
				}
			}
			dataDisk["storage_account_type"] = storageAccountType

			dataDisk["managed_disk_id"] = managedDiskID
			attachDisks = append(attachDisks, dataDisk)

		default:
			return nil, fmt.Errorf("unsupported `createOption` type while flattening: %s", string(createOption))
		}
	}
	return []interface{}{
		map[string]interface{}{
			"create": schema.NewSet(resourceVirtualMachineCreateDataDiskHash, createDisks),
			"attach": schema.NewSet(resourceVirtualMachineAttachDataDiskHash, attachDisks),
		},
	}, nil
}

func resourceVirtualMachineCreateDataDiskHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(fmt.Sprintf("%s-", m["name"].(string)))
		buf.WriteString(fmt.Sprintf("%d-", m["lun"].(int)))
		buf.WriteString(fmt.Sprintf("%s-", m["caching"].(string)))
		buf.WriteString(fmt.Sprintf("%s-", m["storage_account_type"].(string)))
		buf.WriteString(fmt.Sprintf("%d-", m["disk_size_gb"].(int)))
		// Due to potential case diff in API response needs to be normalised
		buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(m["disk_encryption_set_id"].(string))))
		buf.WriteString(fmt.Sprintf("%t-", m["write_accelerator_enabled"].(bool)))
	}

	return hashcode.String(buf.String())
}

func resourceVirtualMachineAttachDataDiskHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(fmt.Sprintf("%s-", m["managed_disk_id"].(string)))
		buf.WriteString(fmt.Sprintf("%s-", m["caching"].(string)))
		buf.WriteString(fmt.Sprintf("%t-", m["write_accelerator_enabled"].(bool)))
	}

	return hashcode.String(buf.String())
}

func findInvalidDataDiskChanges(d *schema.ResourceData) (errors []error) {
	vmName := d.Get("name").(string)
	newCreateRaw, oldCreateRaw := d.GetChange("data_disks.0.create")
	for _, n := range newCreateRaw.(*schema.Set).List() {
		newDisk := n.(map[string]interface{})
		for _, o := range oldCreateRaw.(*schema.Set).List() {
			oldDisk := o.(map[string]interface{})
			if newDisk["name"].(string) == oldDisk["name"].(string) {
				if newDisk["storage_account_type"] != oldDisk["storage_account_type"] {
					errors = append(errors, fmt.Errorf("updating 'storage_account_type' for Data Disks is not supported (Disk %q, Virtual Machine %q)", newDisk["name"].(string), vmName))
				}
				if newDisk["lun"] != oldDisk["lun"] {
					errors = append(errors, fmt.Errorf("updating 'lun' for Data Disks is not supported (Disk %q, Virtual Machine %q)", newDisk["name"].(string), vmName))
				}
			}
		}
	}

	newAttachRaw, oldAttachRaw := d.GetChange("data_disks.0.attach")
	for _, n := range newAttachRaw.(*schema.Set).List() {
		newDisk := n.(map[string]interface{})
		for _, o := range oldAttachRaw.(*schema.Set).List() {
			oldDisk := o.(map[string]interface{})
			if newDisk["managed_disk_id"].(string) == oldDisk["managed_disk_id"].(string) {
				if newDisk["storage_account_type"] != oldDisk["storage_account_type"] {
					errors = append(errors, fmt.Errorf("updating 'storage_account_type' for Data Disks is not supported (Disk %q, Virtual Machine %q)", newDisk["managed_disk_id"].(string), vmName))
				}
				if newDisk["lun"] != oldDisk["lun"] {
					errors = append(errors, fmt.Errorf("updating 'lun' for Data Disks is not supported (Disk %q, Virtual Machine %q)", newDisk["managed_disk_id"].(string), vmName))
				}
			}
		}
	}

	return errors
}
