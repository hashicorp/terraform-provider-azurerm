// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_machine

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-04-02/disks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
)

func flattenAzureRmVirtualMachinePlan(plan *virtualmachines.Plan) []any {
	if plan == nil {
		return []any{}
	}

	result := make(map[string]any)

	if plan.Name != nil {
		result["name"] = *plan.Name
	}
	if plan.Publisher != nil {
		result["publisher"] = *plan.Publisher
	}
	if plan.Product != nil {
		result["product"] = *plan.Product
	}

	return []any{result}
}

func flattenAzureRmVirtualMachineImageReference(image *virtualmachines.ImageReference) []any {
	if image == nil {
		return []any{}
	}

	result := make(map[string]any)
	if image.Publisher != nil {
		result["publisher"] = *image.Publisher
	}
	if image.Offer != nil {
		result["offer"] = *image.Offer
	}
	if image.Sku != nil {
		result["sku"] = *image.Sku
	}
	if image.Version != nil {
		result["version"] = *image.Version
	}
	if image.Id != nil {
		result["id"] = *image.Id
	}

	return []any{result}
}

func flattenAzureRmVirtualMachineDiagnosticsProfile(profile *virtualmachines.BootDiagnostics) []any {
	if profile == nil {
		return []any{}
	}

	result := make(map[string]any)

	if profile.Enabled != nil {
		result["enabled"] = *profile.Enabled
	}

	if profile.StorageUri != nil {
		result["storage_uri"] = *profile.StorageUri
	}

	return []any{result}
}

func flattenAzureRmVirtualMachineAdditionalCapabilities(profile *virtualmachines.AdditionalCapabilities) []any {
	if profile == nil {
		return []any{}
	}

	result := make(map[string]any)
	if v := profile.UltraSSDEnabled; v != nil {
		result["ultra_ssd_enabled"] = *v
	}
	return []any{result}
}

func flattenAzureRmVirtualMachineNetworkInterfaces(profile *virtualmachines.NetworkProfile) []any {
	result := make([]any, 0)
	for _, nic := range *profile.NetworkInterfaces {
		result = append(result, *nic.Id)
	}
	return result
}

func flattenAzureRmVirtualMachineOsProfileSecrets(secrets *[]virtualmachines.VaultSecretGroup) []any {
	if secrets == nil {
		return []any{}
	}

	result := make([]any, 0)
	for _, secret := range *secrets {
		s := map[string]any{
			"source_vault_id": *secret.SourceVault.Id,
		}

		if secret.VaultCertificates != nil {
			certs := make([]map[string]any, 0, len(*secret.VaultCertificates))
			for _, cert := range *secret.VaultCertificates {
				vaultCert := make(map[string]any)
				vaultCert["certificate_url"] = *cert.CertificateURL

				if cert.CertificateStore != nil {
					vaultCert["certificate_store"] = *cert.CertificateStore
				}

				certs = append(certs, vaultCert)
			}

			s["vault_certificates"] = certs
		}

		result = append(result, s)
	}
	return result
}

func flattenAzureRmVirtualMachineDataDisk(disks *[]virtualmachines.DataDisk, disksInfo []*disks.Disk) any {
	result := make([]any, len(*disks))
	for i, disk := range *disks {
		l := make(map[string]any)
		l["name"] = *disk.Name
		if disk.Vhd != nil {
			l["vhd_uri"] = *disk.Vhd.Uri
		}
		if disk.ManagedDisk != nil {
			l["managed_disk_type"] = pointer.FromEnum(disk.ManagedDisk.StorageAccountType)
			if disk.ManagedDisk.Id != nil {
				l["managed_disk_id"] = *disk.ManagedDisk.Id
			}
		}
		l["create_option"] = disk.CreateOption
		l["caching"] = pointer.FromEnum(disk.Caching)
		if disk.DiskSizeGB != nil {
			l["disk_size_gb"] = *disk.DiskSizeGB
		}

		l["lun"] = disk.Lun

		if v := disk.WriteAcceleratorEnabled; v != nil {
			l["write_accelerator_enabled"] = *disk.WriteAcceleratorEnabled
		}

		flattenAzureRmVirtualMachineReviseDiskInfo(l, disksInfo[i])

		result[i] = l
	}
	return result
}

func flattenAzureRmVirtualMachineOsProfile(input *virtualmachines.OSProfile) []any {
	result := make(map[string]any)
	result["computer_name"] = pointer.From(input.ComputerName)
	result["admin_username"] = pointer.From(input.AdminUsername)
	if input.CustomData != nil {
		result["custom_data"] = *input.CustomData
	}

	return []any{result}
}

func flattenAzureRmVirtualMachineOsProfileWindowsConfiguration(config *virtualmachines.WindowsConfiguration) []any {
	if config == nil {
		return []any{}
	}

	result := make(map[string]any)

	if config.ProvisionVMAgent != nil {
		result["provision_vm_agent"] = *config.ProvisionVMAgent
	}

	if config.EnableAutomaticUpdates != nil {
		result["enable_automatic_upgrades"] = *config.EnableAutomaticUpdates
	}

	if config.TimeZone != nil {
		result["timezone"] = *config.TimeZone
	}

	listeners := make([]map[string]any, 0)
	if config.WinRM != nil && config.WinRM.Listeners != nil {
		for _, i := range *config.WinRM.Listeners {
			listener := make(map[string]any)
			listener["protocol"] = pointer.FromEnum(i.Protocol)

			if i.CertificateURL != nil {
				listener["certificate_url"] = *i.CertificateURL
			}

			listeners = append(listeners, listener)
		}
	}

	result["winrm"] = listeners

	content := make([]map[string]any, 0)
	if config.AdditionalUnattendContent != nil {
		for _, i := range *config.AdditionalUnattendContent {
			c := make(map[string]any)
			c["pass"] = pointer.FromEnum(i.PassName)
			c["component"] = pointer.FromEnum(i.ComponentName)
			c["setting_name"] = pointer.FromEnum(i.SettingName)

			if i.Content != nil {
				c["content"] = *i.Content
			}

			content = append(content, c)
		}
	}
	result["additional_unattend_config"] = content

	return []any{result}
}

func flattenAzureRmVirtualMachineOsProfileLinuxConfiguration(config *virtualmachines.LinuxConfiguration) []any {
	if config == nil {
		return []any{}
	}

	result := make(map[string]any)

	if config.DisablePasswordAuthentication != nil {
		result["disable_password_authentication"] = *config.DisablePasswordAuthentication
	}

	if config.Ssh != nil && config.Ssh.PublicKeys != nil && len(*config.Ssh.PublicKeys) > 0 {
		ssh_keys := make([]map[string]any, 0)
		for _, i := range *config.Ssh.PublicKeys {
			key := make(map[string]any)

			if i.Path != nil {
				key["path"] = *i.Path
			}

			if i.KeyData != nil {
				key["key_data"] = *i.KeyData
			}

			ssh_keys = append(ssh_keys, key)
		}

		result["ssh_keys"] = ssh_keys
	}

	return []any{result}
}

func flattenAzureRmVirtualMachineOsDisk(disk *virtualmachines.OSDisk, diskInfo *disks.Disk) []any {
	result := make(map[string]any)
	if disk.Name != nil {
		result["name"] = *disk.Name
	}
	if disk.Vhd != nil && disk.Vhd.Uri != nil {
		result["vhd_uri"] = *disk.Vhd.Uri
	}
	if disk.Image != nil && disk.Image.Uri != nil {
		result["image_uri"] = *disk.Image.Uri
	}
	if disk.ManagedDisk != nil {
		result["managed_disk_type"] = pointer.FromEnum(disk.ManagedDisk.StorageAccountType)
		if disk.ManagedDisk.Id != nil {
			result["managed_disk_id"] = *disk.ManagedDisk.Id
		}
	}
	result["create_option"] = disk.CreateOption
	result["caching"] = disk.Caching
	if disk.DiskSizeGB != nil {
		result["disk_size_gb"] = *disk.DiskSizeGB
	}
	result["os_type"] = pointer.FromEnum(disk.OsType)

	if v := disk.WriteAcceleratorEnabled; v != nil {
		result["write_accelerator_enabled"] = *disk.WriteAcceleratorEnabled
	}

	flattenAzureRmVirtualMachineReviseDiskInfo(result, diskInfo)

	return []any{result}
}

func flattenAzureRmVirtualMachineReviseDiskInfo(result map[string]any, diskInfo *disks.Disk) {
	if diskInfo != nil {
		if diskInfo.Sku != nil {
			result["managed_disk_type"] = pointer.FromEnum(diskInfo.Sku.Name)
		}
		if diskInfo.Properties != nil && diskInfo.Properties.DiskSizeGB != nil {
			result["disk_size_gb"] = *diskInfo.Properties.DiskSizeGB
		}
	}
}
