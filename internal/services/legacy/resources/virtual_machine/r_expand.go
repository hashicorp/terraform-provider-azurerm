// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_machine

import (
	"errors"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/base64"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandAzureRmVirtualMachinePlan(d *pluginsdk.ResourceData) *virtualmachines.Plan {
	planConfigs := d.Get("plan").([]any)
	if len(planConfigs) == 0 {
		return nil
	}

	planConfig := planConfigs[0].(map[string]any)

	return &virtualmachines.Plan{
		Publisher: pointer.To(planConfig["publisher"].(string)),
		Name:      pointer.To(planConfig["name"].(string)),
		Product:   pointer.To(planConfig["product"].(string)),
	}
}

func expandAzureRmVirtualMachineOsProfile(d *pluginsdk.ResourceData) (*virtualmachines.OSProfile, error) {
	osProfiles := d.Get("os_profile").(*pluginsdk.Set).List()

	osProfile := osProfiles[0].(map[string]any)

	adminPassword := osProfile["admin_password"].(string)

	profile := &virtualmachines.OSProfile{
		AdminUsername: pointer.To(osProfile["admin_username"].(string)),
		ComputerName:  pointer.To(osProfile["computer_name"].(string)),
	}

	if adminPassword != "" {
		profile.AdminPassword = &adminPassword
	}

	if _, ok := d.GetOk("os_profile_windows_config"); ok {
		winConfig := expandAzureRmVirtualMachineOsProfileWindowsConfig(d)
		if winConfig != nil {
			profile.WindowsConfiguration = winConfig
		}
	}

	if _, ok := d.GetOk("os_profile_linux_config"); ok {
		linuxConfig := expandAzureRmVirtualMachineOsProfileLinuxConfig(d)
		if linuxConfig != nil {
			profile.LinuxConfiguration = linuxConfig
		}
	}

	if profile.LinuxConfiguration == nil && profile.WindowsConfiguration == nil {
		return nil, errors.New("either a `os_profile_linux_config` or a `os_profile_windows_config` must be specified")
	}

	if _, ok := d.GetOk("os_profile_secrets"); ok {
		secrets := expandAzureRmVirtualMachineOsProfileSecrets(d)
		if secrets != nil {
			profile.Secrets = secrets
		}
	}

	if v := osProfile["custom_data"].(string); v != "" {
		v = base64.EncodeIfNot(v)
		profile.CustomData = &v
	}

	return profile, nil
}

func expandAzureRmVirtualMachineOsProfileSecrets(d *pluginsdk.ResourceData) *[]virtualmachines.VaultSecretGroup {
	secretsConfig := d.Get("os_profile_secrets").([]any)
	secrets := make([]virtualmachines.VaultSecretGroup, 0, len(secretsConfig))

	for _, secretConfig := range secretsConfig {
		if secretConfig == nil {
			continue
		}

		config := secretConfig.(map[string]any)

		vaultSecretGroup := virtualmachines.VaultSecretGroup{
			SourceVault: &virtualmachines.SubResource{
				Id: pointer.To(config["source_vault_id"].(string)),
			},
		}

		if v := config["vault_certificates"]; v != nil {
			certsConfig := v.([]any)
			certs := make([]virtualmachines.VaultCertificate, 0, len(certsConfig))
			for _, certConfig := range certsConfig {
				if certConfig == nil {
					continue
				}
				config := certConfig.(map[string]any)

				cert := virtualmachines.VaultCertificate{
					CertificateURL: pointer.To(config["certificate_url"].(string)),
				}
				if v := config["certificate_store"].(string); v != "" {
					cert.CertificateStore = &v
				}

				certs = append(certs, cert)
			}
			vaultSecretGroup.VaultCertificates = &certs
		}

		secrets = append(secrets, vaultSecretGroup)
	}

	return &secrets
}

func expandAzureRmVirtualMachineOsProfileLinuxConfig(d *pluginsdk.ResourceData) *virtualmachines.LinuxConfiguration {
	osProfilesLinuxConfig := d.Get("os_profile_linux_config").(*pluginsdk.Set).List()

	linuxConfig := osProfilesLinuxConfig[0].(map[string]any)

	config := &virtualmachines.LinuxConfiguration{
		DisablePasswordAuthentication: pointer.To(linuxConfig["disable_password_authentication"].(bool)),
	}

	linuxKeys := linuxConfig["ssh_keys"].([]any)
	sshPublicKeys := make([]virtualmachines.SshPublicKey, 0)
	for _, key := range linuxKeys {
		sshKey, ok := key.(map[string]any)
		if !ok {
			continue
		}

		sshPublicKey := virtualmachines.SshPublicKey{
			Path:    pointer.To(sshKey["path"].(string)),
			KeyData: pointer.To(sshKey["key_data"].(string)),
		}

		sshPublicKeys = append(sshPublicKeys, sshPublicKey)
	}

	if len(sshPublicKeys) > 0 {
		config.Ssh = &virtualmachines.SshConfiguration{
			PublicKeys: &sshPublicKeys,
		}
	}

	return config
}

func expandAzureRmVirtualMachineOsProfileWindowsConfig(d *pluginsdk.ResourceData) *virtualmachines.WindowsConfiguration {
	osProfilesWindowsConfig := d.Get("os_profile_windows_config").(*pluginsdk.Set).List()

	osProfileConfig := osProfilesWindowsConfig[0].(map[string]any)
	config := &virtualmachines.WindowsConfiguration{}

	if v := osProfileConfig["provision_vm_agent"]; v != nil {
		config.ProvisionVMAgent = pointer.To(v.(bool))
	}

	if v := osProfileConfig["enable_automatic_upgrades"]; v != nil {
		config.EnableAutomaticUpdates = pointer.To(v.(bool))
	}

	if v := osProfileConfig["timezone"]; v != nil && v.(string) != "" {
		config.TimeZone = pointer.To(v.(string))
	}

	if v := osProfileConfig["winrm"]; v != nil {
		winRm := v.([]any)
		if len(winRm) > 0 {
			winRmListeners := make([]virtualmachines.WinRMListener, 0, len(winRm))
			for _, winRmConfig := range winRm {
				config := winRmConfig.(map[string]any)

				protocol := config["protocol"].(string)
				winRmListener := virtualmachines.WinRMListener{
					Protocol: pointer.ToEnum[virtualmachines.ProtocolTypes](protocol),
				}
				if v := config["certificate_url"].(string); v != "" {
					winRmListener.CertificateURL = &v
				}

				winRmListeners = append(winRmListeners, winRmListener)
			}
			config.WinRM = &virtualmachines.WinRMConfiguration{
				Listeners: &winRmListeners,
			}
		}
	}
	if v := osProfileConfig["additional_unattend_config"]; v != nil {
		additionalConfig := v.([]any)
		if len(additionalConfig) > 0 {
			additionalConfigContent := make([]virtualmachines.AdditionalUnattendContent, 0, len(additionalConfig))
			for _, addConfig := range additionalConfig {
				config := addConfig.(map[string]any)
				pass := config["pass"].(string)
				component := config["component"].(string)
				settingName := config["setting_name"].(string)
				content := config["content"].(string)

				addContent := virtualmachines.AdditionalUnattendContent{
					PassName:      pointer.ToEnum[virtualmachines.PassNames](pass),
					ComponentName: pointer.ToEnum[virtualmachines.ComponentNames](component),
					SettingName:   pointer.ToEnum[virtualmachines.SettingNames](settingName),
				}

				if content != "" {
					addContent.Content = &content
				}

				additionalConfigContent = append(additionalConfigContent, addContent)
			}
			config.AdditionalUnattendContent = &additionalConfigContent
		}
	}
	return config
}

func expandAzureRmVirtualMachineDataDisk(d *pluginsdk.ResourceData) ([]virtualmachines.DataDisk, error) {
	disks := d.Get("storage_data_disk").([]any)
	data_disks := make([]virtualmachines.DataDisk, 0, len(disks))
	for _, disk_config := range disks {
		config := disk_config.(map[string]any)

		createOption := config["create_option"].(string)
		vhdURI := config["vhd_uri"].(string)
		managedDiskType := config["managed_disk_type"].(string)
		managedDiskID := config["managed_disk_id"].(string)
		lun := int64(config["lun"].(int))

		data_disk := virtualmachines.DataDisk{
			Name:         pointer.To(config["name"].(string)),
			Lun:          lun,
			CreateOption: virtualmachines.DiskCreateOptionTypes(createOption),
		}

		if vhdURI != "" {
			data_disk.Vhd = &virtualmachines.VirtualHardDisk{
				Uri: &vhdURI,
			}
		}

		managedDisk := &virtualmachines.ManagedDiskParameters{}

		if managedDiskType != "" {
			managedDisk.StorageAccountType = pointer.ToEnum[virtualmachines.StorageAccountTypes](managedDiskType)
			data_disk.ManagedDisk = managedDisk
		}

		if managedDiskID != "" {
			managedDisk.Id = &managedDiskID
			data_disk.ManagedDisk = managedDisk
		}

		if vhdURI != "" && managedDiskID != "" {
			return nil, errors.New("[ERROR] Conflict between `vhd_uri` and `managed_disk_id` (only one or the other can be used)")
		}
		if vhdURI != "" && managedDiskType != "" {
			return nil, errors.New("[ERROR] Conflict between `vhd_uri` and `managed_disk_type` (only one or the other can be used)")
		}
		if managedDiskID == "" && vhdURI == "" && strings.EqualFold(string(data_disk.CreateOption), string(virtualmachines.DiskCreateOptionTypesAttach)) {
			return nil, errors.New("[ERROR] Must specify `vhd_uri` or `managed_disk_id` to attach")
		}

		if v := config["caching"].(string); v != "" {
			data_disk.Caching = pointer.ToEnum[virtualmachines.CachingTypes](v)
		}

		if v, ok := config["disk_size_gb"].(int); ok {
			data_disk.DiskSizeGB = pointer.To(int64(v))
		}

		if v, ok := config["write_accelerator_enabled"].(bool); ok {
			data_disk.WriteAcceleratorEnabled = pointer.To(v)
		}

		data_disks = append(data_disks, data_disk)
	}

	return data_disks, nil
}

func expandAzureRmVirtualMachineDiagnosticsProfile(d *pluginsdk.ResourceData) *virtualmachines.DiagnosticsProfile {
	bootDiagnostics := d.Get("boot_diagnostics").([]any)

	diagnosticsProfile := &virtualmachines.DiagnosticsProfile{}
	if len(bootDiagnostics) > 0 {
		bootDiagnostic := bootDiagnostics[0].(map[string]any)

		diagnosticsProfile.BootDiagnostics = &virtualmachines.BootDiagnostics{
			Enabled:    pointer.To(bootDiagnostic["enabled"].(bool)),
			StorageUri: pointer.To(bootDiagnostic["storage_uri"].(string)),
		}

		return diagnosticsProfile
	}

	return nil
}

func expandAzureRmVirtualMachineAdditionalCapabilities(d *pluginsdk.ResourceData) *virtualmachines.AdditionalCapabilities {
	additionalCapabilities := d.Get("additional_capabilities").([]any)
	if len(additionalCapabilities) == 0 || additionalCapabilities[0] == nil {
		return nil
	}

	additionalCapability := additionalCapabilities[0].(map[string]any)
	return &virtualmachines.AdditionalCapabilities{
		UltraSSDEnabled: pointer.To(additionalCapability["ultra_ssd_enabled"].(bool)),
	}
}

func expandAzureRmVirtualMachineImageReference(d *pluginsdk.ResourceData) (*virtualmachines.ImageReference, error) {
	storageImageRefs := d.Get("storage_image_reference").(*pluginsdk.Set).List()

	storageImageRef := storageImageRefs[0].(map[string]any)
	imageID := storageImageRef["id"].(string)
	publisher := storageImageRef["publisher"].(string)

	imageReference := virtualmachines.ImageReference{}

	if imageID != "" && publisher != "" {
		return nil, errors.New("conflict between `id` and `publisher` (only one or the other can be used)")
	}

	if imageID != "" {
		imageReference.Id = pointer.To(storageImageRef["id"].(string))
	} else {
		imageReference = virtualmachines.ImageReference{
			Publisher: &publisher,
			Offer:     pointer.To(storageImageRef["offer"].(string)),
			Sku:       pointer.To(storageImageRef["sku"].(string)),
			Version:   pointer.To(storageImageRef["version"].(string)),
		}
	}

	return &imageReference, nil
}

func expandAzureRmVirtualMachineNetworkProfile(d *pluginsdk.ResourceData) virtualmachines.NetworkProfile {
	nicIds := d.Get("network_interface_ids").([]any)
	primaryNicId := d.Get("primary_network_interface_id").(string)
	network_interfaces := make([]virtualmachines.NetworkInterfaceReference, 0, len(nicIds))

	network_profile := virtualmachines.NetworkProfile{}

	for _, nic := range nicIds {
		if nic != nil {
			id := nic.(string)

			network_interface := virtualmachines.NetworkInterfaceReference{
				Id: &id,
				Properties: &virtualmachines.NetworkInterfaceReferenceProperties{
					Primary: pointer.To(id == primaryNicId),
				},
			}
			network_interfaces = append(network_interfaces, network_interface)
		}
	}

	network_profile.NetworkInterfaces = &network_interfaces

	return network_profile
}

func expandAzureRmVirtualMachineOsDisk(d *pluginsdk.ResourceData) (*virtualmachines.OSDisk, error) {
	disks := d.Get("storage_os_disk").([]any)

	config := disks[0].(map[string]any)

	createOption := config["create_option"].(string)
	vhdURI := config["vhd_uri"].(string)
	managedDiskType := config["managed_disk_type"].(string)
	managedDiskID := config["managed_disk_id"].(string)

	osDisk := &virtualmachines.OSDisk{
		Name:         pointer.To(config["name"].(string)),
		CreateOption: virtualmachines.DiskCreateOptionTypes(createOption),
	}

	if vhdURI != "" {
		osDisk.Vhd = &virtualmachines.VirtualHardDisk{
			Uri: &vhdURI,
		}
	}

	managedDisk := &virtualmachines.ManagedDiskParameters{}

	if managedDiskType != "" {
		managedDisk.StorageAccountType = pointer.ToEnum[virtualmachines.StorageAccountTypes](managedDiskType)
		osDisk.ManagedDisk = managedDisk
	}

	if managedDiskID != "" {
		managedDisk.Id = &managedDiskID
		osDisk.ManagedDisk = managedDisk
	}

	// BEGIN: code to be removed after GH-13016 is merged
	if vhdURI != "" && managedDiskID != "" {
		return nil, errors.New("conflict between `vhd_uri` and `managed_disk_id` (only one or the other can be used)")
	}
	if vhdURI != "" && managedDiskType != "" {
		return nil, errors.New("conflict between `vhd_uri` and `managed_disk_type` (only one or the other can be used)")
	}
	// END: code to be removed after GH-13016 is merged
	if managedDiskID == "" && vhdURI == "" && strings.EqualFold(string(osDisk.CreateOption), string(virtualmachines.DiskCreateOptionTypesAttach)) {
		return nil, errors.New("must specify `vhd_uri` or `managed_disk_id` to attach")
	}

	if v := config["image_uri"].(string); v != "" {
		osDisk.Image = &virtualmachines.VirtualHardDisk{
			Uri: pointer.To(config["image_uri"].(string)),
		}
	}

	if v := config["os_type"].(string); v != "" {
		osDisk.OsType = pointer.ToEnum[virtualmachines.OperatingSystemTypes](v)
	}

	if v := config["caching"].(string); v != "" {
		osDisk.Caching = pointer.ToEnum[virtualmachines.CachingTypes](v)
	}

	if v := config["disk_size_gb"].(int); v != 0 {
		osDisk.DiskSizeGB = pointer.To(int64(v))
	}

	if v, ok := config["write_accelerator_enabled"].(bool); ok {
		osDisk.WriteAcceleratorEnabled = pointer.To(v)
	}

	return osDisk, nil
}
