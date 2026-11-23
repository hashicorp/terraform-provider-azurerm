// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_scale_set

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2025-04-01/virtualmachinescalesets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/base64"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandVirtualMachineScaleSetSku(d *pluginsdk.ResourceData) *virtualmachinescalesets.Sku {
	skuConfig := d.Get("sku").([]any)
	config := skuConfig[0].(map[string]any)

	sku := &virtualmachinescalesets.Sku{
		Name:     pointer.To(config["name"].(string)),
		Capacity: pointer.To(int64(config["capacity"].(int))),
	}

	if tier, ok := config["tier"].(string); ok && tier != "" {
		sku.Tier = &tier
	}

	return sku
}

func expandAzureRmRollingUpgradePolicy(d *pluginsdk.ResourceData) *virtualmachinescalesets.RollingUpgradePolicy {
	if config, ok := d.GetOk("rolling_upgrade_policy.0"); ok {
		policy := config.(map[string]any)
		return &virtualmachinescalesets.RollingUpgradePolicy{
			MaxBatchInstancePercent:             pointer.To(int64(policy["max_batch_instance_percent"].(int))),
			MaxUnhealthyInstancePercent:         pointer.To(int64(policy["max_unhealthy_instance_percent"].(int))),
			MaxUnhealthyUpgradedInstancePercent: pointer.To(int64(policy["max_unhealthy_upgraded_instance_percent"].(int))),
			PauseTimeBetweenBatches:             pointer.To(policy["pause_time_between_batches"].(string)),
		}
	}
	return nil
}

func expandAzureRmVirtualMachineScaleSetNetworkProfile(d *pluginsdk.ResourceData) *virtualmachinescalesets.VirtualMachineScaleSetNetworkProfile {
	scaleSetNetworkProfileConfigs := d.Get("network_profile").(*pluginsdk.Set).List()
	networkProfileConfig := make([]virtualmachinescalesets.VirtualMachineScaleSetNetworkConfiguration, 0, len(scaleSetNetworkProfileConfigs))

	for _, npProfileConfig := range scaleSetNetworkProfileConfigs {
		config := npProfileConfig.(map[string]any)

		name := config["name"].(string)
		primary := config["primary"].(bool)
		acceleratedNetworking := config["accelerated_networking"].(bool)
		ipForwarding := config["ip_forwarding"].(bool)

		dnsSettingsConfigs := config["dns_settings"].([]any)
		dnsSettings := virtualmachinescalesets.VirtualMachineScaleSetNetworkConfigurationDnsSettings{}
		for _, dnsSettingsConfig := range dnsSettingsConfigs {
			dns_settings := dnsSettingsConfig.(map[string]any)

			if v := dns_settings["dns_servers"]; v != nil {
				dns_servers := dns_settings["dns_servers"].([]any)
				if len(dns_servers) > 0 {
					var dnsServers []string
					for _, v := range dns_servers {
						str := v.(string)
						dnsServers = append(dnsServers, str)
					}
					dnsSettings.DnsServers = &dnsServers
				}
			}
		}
		ipConfigurationConfigs := config["ip_configuration"].([]any)
		ipConfigurations := make([]virtualmachinescalesets.VirtualMachineScaleSetIPConfiguration, 0, len(ipConfigurationConfigs))
		for _, ipConfigConfig := range ipConfigurationConfigs {
			ipconfig := ipConfigConfig.(map[string]any)
			name := ipconfig["name"].(string)

			ipConfiguration := virtualmachinescalesets.VirtualMachineScaleSetIPConfiguration{
				Name: name,
				Properties: &virtualmachinescalesets.VirtualMachineScaleSetIPConfigurationProperties{
					Subnet: &virtualmachinescalesets.ApiEntityReference{
						Id: pointer.To(ipconfig["subnet_id"].(string)),
					},
				},
			}

			ipConfiguration.Properties.Primary = pointer.To(ipconfig["primary"].(bool))

			if v := ipconfig["application_gateway_backend_address_pool_ids"]; v != nil {
				pools := v.(*pluginsdk.Set).List()
				resources := make([]virtualmachinescalesets.SubResource, 0, len(pools))
				for _, p := range pools {
					resources = append(resources, virtualmachinescalesets.SubResource{
						Id: pointer.To(p.(string)),
					})
				}
				ipConfiguration.Properties.ApplicationGatewayBackendAddressPools = &resources
			}

			if v := ipconfig["application_security_group_ids"]; v != nil {
				asgs := v.(*pluginsdk.Set).List()
				resources := make([]virtualmachinescalesets.SubResource, 0, len(asgs))
				for _, p := range asgs {
					resources = append(resources, virtualmachinescalesets.SubResource{
						Id: pointer.To(p.(string)),
					})
				}
				ipConfiguration.Properties.ApplicationSecurityGroups = &resources
			}

			if v := ipconfig["load_balancer_backend_address_pool_ids"]; v != nil {
				pools := v.(*pluginsdk.Set).List()
				resources := make([]virtualmachinescalesets.SubResource, 0, len(pools))
				for _, p := range pools {
					resources = append(resources, virtualmachinescalesets.SubResource{
						Id: pointer.To(p.(string)),
					})
				}
				ipConfiguration.Properties.LoadBalancerBackendAddressPools = &resources
			}

			if v := ipconfig["load_balancer_inbound_nat_rules_ids"]; v != nil {
				rules := v.(*pluginsdk.Set).List()
				rulesResources := make([]virtualmachinescalesets.SubResource, 0, len(rules))
				for _, m := range rules {
					rulesResources = append(rulesResources, virtualmachinescalesets.SubResource{
						Id: pointer.To(m.(string)),
					})
				}
				ipConfiguration.Properties.LoadBalancerInboundNatPools = &rulesResources
			}

			if v := ipconfig["public_ip_address_configuration"]; v != nil {
				publicIpConfigs := v.([]any)
				for _, publicIpConfigConfig := range publicIpConfigs {
					publicIpConfig := publicIpConfigConfig.(map[string]any)

					dnsSettings := virtualmachinescalesets.VirtualMachineScaleSetPublicIPAddressConfigurationDnsSettings{
						DomainNameLabel: publicIpConfig["domain_name_label"].(string),
					}

					prop := virtualmachinescalesets.VirtualMachineScaleSetPublicIPAddressConfigurationProperties{
						DnsSettings:          &dnsSettings,
						IdleTimeoutInMinutes: pointer.To(int64(publicIpConfig["idle_timeout"].(int))),
					}

					config := virtualmachinescalesets.VirtualMachineScaleSetPublicIPAddressConfiguration{
						Name:       publicIpConfig["name"].(string),
						Properties: &prop,
					}
					ipConfiguration.Properties.PublicIPAddressConfiguration = &config
				}
			}

			ipConfigurations = append(ipConfigurations, ipConfiguration)
		}

		nProfile := virtualmachinescalesets.VirtualMachineScaleSetNetworkConfiguration{
			Name: name,
			Properties: &virtualmachinescalesets.VirtualMachineScaleSetNetworkConfigurationProperties{
				Primary:                     &primary,
				IPConfigurations:            ipConfigurations,
				EnableAcceleratedNetworking: &acceleratedNetworking,
				EnableIPForwarding:          &ipForwarding,
				DnsSettings:                 &dnsSettings,
			},
		}

		if v := config["network_security_group_id"].(string); v != "" {
			networkSecurityGroupId := virtualmachinescalesets.SubResource{
				Id: &v,
			}
			nProfile.Properties.NetworkSecurityGroup = &networkSecurityGroupId
		}

		networkProfileConfig = append(networkProfileConfig, nProfile)
	}

	return &virtualmachinescalesets.VirtualMachineScaleSetNetworkProfile{
		NetworkInterfaceConfigurations: &networkProfileConfig,
	}
}

func expandAzureRMVirtualMachineScaleSetsOsProfile(d *pluginsdk.ResourceData) *virtualmachinescalesets.VirtualMachineScaleSetOSProfile {
	osProfileConfigs := d.Get("os_profile").([]any)

	osProfileConfig := osProfileConfigs[0].(map[string]any)
	password := osProfileConfig["admin_password"].(string)
	customData := osProfileConfig["custom_data"].(string)

	osProfile := &virtualmachinescalesets.VirtualMachineScaleSetOSProfile{
		ComputerNamePrefix: pointer.To(osProfileConfig["computer_name_prefix"].(string)),
		AdminUsername:      pointer.To(osProfileConfig["admin_username"].(string)),
	}

	if password != "" {
		osProfile.AdminPassword = &password
	}

	if customData != "" {
		customData = base64.EncodeIfNot(customData)
		osProfile.CustomData = &customData
	}

	if _, ok := d.GetOk("os_profile_secrets"); ok {
		secrets := expandAzureRmVirtualMachineScaleSetOsProfileSecrets(d)
		if secrets != nil {
			osProfile.Secrets = secrets
		}
	}

	if _, ok := d.GetOk("os_profile_linux_config"); ok {
		osProfile.LinuxConfiguration = expandAzureRmVirtualMachineScaleSetOsProfileLinuxConfig(d)
	}

	if _, ok := d.GetOk("os_profile_windows_config"); ok {
		winConfig := expandAzureRmVirtualMachineScaleSetOsProfileWindowsConfig(d)
		if winConfig != nil {
			osProfile.WindowsConfiguration = winConfig
		}
	}

	return osProfile
}

func expandAzureRMVirtualMachineScaleSetsDiagnosticProfile(d *pluginsdk.ResourceData) virtualmachinescalesets.DiagnosticsProfile {
	bootDiagnosticConfigs := d.Get("boot_diagnostics").([]any)
	bootDiagnosticConfig := bootDiagnosticConfigs[0].(map[string]any)

	bootDiagnostic := &virtualmachinescalesets.BootDiagnostics{
		Enabled:    pointer.To(bootDiagnosticConfig["enabled"].(bool)),
		StorageUri: pointer.To(bootDiagnosticConfig["storage_uri"].(string)),
	}

	return virtualmachinescalesets.DiagnosticsProfile{
		BootDiagnostics: bootDiagnostic,
	}
}

func expandAzureRMVirtualMachineScaleSetsStorageProfileOsDisk(d *pluginsdk.ResourceData) (*virtualmachinescalesets.VirtualMachineScaleSetOSDisk, error) {
	osDiskConfigs := d.Get("storage_profile_os_disk").(*pluginsdk.Set).List()

	osDiskConfig := osDiskConfigs[0].(map[string]any)
	name := osDiskConfig["name"].(string)
	image := osDiskConfig["image"].(string)
	vhd_containers := osDiskConfig["vhd_containers"].(*pluginsdk.Set).List()
	caching := osDiskConfig["caching"].(string)
	osType := osDiskConfig["os_type"].(string)
	createOption := osDiskConfig["create_option"].(string)
	managedDiskType := osDiskConfig["managed_disk_type"].(string)

	if managedDiskType == "" && name == "" {
		return nil, fmt.Errorf("[ERROR] `name` must be set in `storage_profile_os_disk` for unmanaged disk")
	}

	osDisk := &virtualmachinescalesets.VirtualMachineScaleSetOSDisk{
		Name:         &name,
		Caching:      pointer.ToEnum[virtualmachinescalesets.CachingTypes](caching),
		OsType:       pointer.ToEnum[virtualmachinescalesets.OperatingSystemTypes](osType),
		CreateOption: virtualmachinescalesets.DiskCreateOptionTypes(createOption),
	}

	if image != "" {
		osDisk.Image = &virtualmachinescalesets.VirtualHardDisk{
			Uri: &image,
		}
	}

	if len(vhd_containers) > 0 {
		var vhdContainers []string
		for _, v := range vhd_containers {
			str := v.(string)
			vhdContainers = append(vhdContainers, str)
		}
		osDisk.VhdContainers = &vhdContainers
	}

	managedDisk := &virtualmachinescalesets.VirtualMachineScaleSetManagedDiskParameters{}

	if managedDiskType != "" {
		if name != "" {
			return nil, fmt.Errorf("[ERROR] Conflict between `name` and `managed_disk_type` on `storage_profile_os_disk` (please remove name or set it to blank)")
		}

		osDisk.Name = nil
		managedDisk.StorageAccountType = pointer.ToEnum[virtualmachinescalesets.StorageAccountTypes](managedDiskType)
		osDisk.ManagedDisk = managedDisk
	}

	// BEGIN: code to be removed after GH-13016 is merged
	if image != "" && managedDiskType != "" {
		return nil, fmt.Errorf("[ERROR] Conflict between `image` and `managed_disk_type` on `storage_profile_os_disk` (only one or the other can be used)")
	}

	if len(vhd_containers) > 0 && managedDiskType != "" {
		return nil, fmt.Errorf("[ERROR] Conflict between `vhd_containers` and `managed_disk_type` on `storage_profile_os_disk` (only one or the other can be used)")
	}
	// END: code to be removed after GH-13016 is merged

	return osDisk, nil
}

func expandAzureRMVirtualMachineScaleSetsStorageProfileDataDisk(d *pluginsdk.ResourceData) *[]virtualmachinescalesets.VirtualMachineScaleSetDataDisk {
	disks := d.Get("storage_profile_data_disk").([]any)
	dataDisks := make([]virtualmachinescalesets.VirtualMachineScaleSetDataDisk, 0, len(disks))
	for _, diskConfig := range disks {
		config := diskConfig.(map[string]any)
		managedDiskType := config["managed_disk_type"].(string)
		dataDisk := virtualmachinescalesets.VirtualMachineScaleSetDataDisk{
			Lun:          int64(config["lun"].(int)),
			CreateOption: virtualmachinescalesets.DiskCreateOptionTypes(config["create_option"].(string)),
		}

		managedDiskVMSS := &virtualmachinescalesets.VirtualMachineScaleSetManagedDiskParameters{}

		if managedDiskType != "" {
			managedDiskVMSS.StorageAccountType = pointer.ToEnum[virtualmachinescalesets.StorageAccountTypes](managedDiskType)
		} else {
			managedDiskVMSS.StorageAccountType = pointer.To(virtualmachinescalesets.StorageAccountTypesStandardLRS)
		}

		// assume that data disks in VMSS can only be Managed Disks
		dataDisk.ManagedDisk = managedDiskVMSS
		if v := config["caching"].(string); v != "" {
			dataDisk.Caching = pointer.ToEnum[virtualmachinescalesets.CachingTypes](v)
		}

		if v := config["disk_size_gb"]; v != nil {
			dataDisk.DiskSizeGB = pointer.To(int64(config["disk_size_gb"].(int)))
		}

		dataDisks = append(dataDisks, dataDisk)
	}

	return &dataDisks
}

func expandAzureRmVirtualMachineScaleSetStorageProfileImageReference(d *pluginsdk.ResourceData) (*virtualmachinescalesets.ImageReference, error) {
	storageImageRefs := d.Get("storage_profile_image_reference").(*pluginsdk.Set).List()

	storageImageRef := storageImageRefs[0].(map[string]any)

	imageID := storageImageRef["id"].(string)
	publisher := storageImageRef["publisher"].(string)

	imageReference := virtualmachinescalesets.ImageReference{}

	if imageID != "" && publisher != "" {
		return nil, fmt.Errorf("[ERROR] Conflict between `id` and `publisher` (only one or the other can be used)")
	}

	if imageID != "" {
		imageReference.Id = pointer.To(storageImageRef["id"].(string))
	} else {
		offer := storageImageRef["offer"].(string)
		sku := storageImageRef["sku"].(string)
		version := storageImageRef["version"].(string)

		imageReference.Publisher = pointer.To(publisher)
		imageReference.Offer = pointer.To(offer)
		imageReference.Sku = pointer.To(sku)
		imageReference.Version = pointer.To(version)
	}

	return &imageReference, nil
}

func expandAzureRmVirtualMachineScaleSetOsProfileLinuxConfig(d *pluginsdk.ResourceData) *virtualmachinescalesets.LinuxConfiguration {
	osProfilesLinuxConfig := d.Get("os_profile_linux_config").(*pluginsdk.Set).List()

	linuxConfig := osProfilesLinuxConfig[0].(map[string]any)

	linuxKeys := linuxConfig["ssh_keys"].([]any)
	sshPublicKeys := make([]virtualmachinescalesets.SshPublicKey, 0, len(linuxKeys))
	for _, key := range linuxKeys {
		if key == nil {
			continue
		}
		sshKey := key.(map[string]any)

		sshPublicKey := virtualmachinescalesets.SshPublicKey{
			Path:    pointer.To(sshKey["path"].(string)),
			KeyData: pointer.To(sshKey["key_data"].(string)),
		}

		sshPublicKeys = append(sshPublicKeys, sshPublicKey)
	}

	return &virtualmachinescalesets.LinuxConfiguration{
		DisablePasswordAuthentication: pointer.To(linuxConfig["disable_password_authentication"].(bool)),
		Ssh: &virtualmachinescalesets.SshConfiguration{
			PublicKeys: &sshPublicKeys,
		},
	}
}

func expandAzureRmVirtualMachineScaleSetOsProfileWindowsConfig(d *pluginsdk.ResourceData) *virtualmachinescalesets.WindowsConfiguration {
	osProfilesWindowsConfig := d.Get("os_profile_windows_config").(*pluginsdk.Set).List()

	osProfileConfig := osProfilesWindowsConfig[0].(map[string]any)
	config := &virtualmachinescalesets.WindowsConfiguration{}

	if v := osProfileConfig["provision_vm_agent"]; v != nil {
		config.ProvisionVMAgent = pointer.To(v.(bool))
	}

	if v := osProfileConfig["enable_automatic_upgrades"]; v != nil {
		config.EnableAutomaticUpdates = pointer.To(v.(bool))
	}

	if v := osProfileConfig["winrm"]; v != nil {
		winRm := v.([]any)
		if len(winRm) > 0 {
			winRmListeners := make([]virtualmachinescalesets.WinRMListener, 0, len(winRm))
			for _, winRmConfig := range winRm {
				config := winRmConfig.(map[string]any)

				protocol := config["protocol"].(string)
				winRmListener := virtualmachinescalesets.WinRMListener{
					Protocol: pointer.ToEnum[virtualmachinescalesets.ProtocolTypes](protocol),
				}
				if v := config["certificate_url"].(string); v != "" {
					winRmListener.CertificateURL = &v
				}

				winRmListeners = append(winRmListeners, winRmListener)
			}
			config.WinRM = &virtualmachinescalesets.WinRMConfiguration{
				Listeners: &winRmListeners,
			}
		}
	}
	if v := osProfileConfig["additional_unattend_config"]; v != nil {
		additionalConfig := v.([]any)
		if len(additionalConfig) > 0 {
			additionalConfigContent := make([]virtualmachinescalesets.AdditionalUnattendContent, 0, len(additionalConfig))
			for _, addConfig := range additionalConfig {
				config := addConfig.(map[string]any)
				pass := config["pass"].(string)
				component := config["component"].(string)
				settingName := config["setting_name"].(string)
				content := config["content"].(string)

				addContent := virtualmachinescalesets.AdditionalUnattendContent{
					PassName:      pointer.ToEnum[virtualmachinescalesets.PassName](pass),
					ComponentName: pointer.ToEnum[virtualmachinescalesets.ComponentName](component),
					SettingName:   pointer.ToEnum[virtualmachinescalesets.SettingNames](settingName),
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

func expandAzureRmVirtualMachineScaleSetOsProfileSecrets(d *pluginsdk.ResourceData) *[]virtualmachinescalesets.VaultSecretGroup {
	secretsConfig := d.Get("os_profile_secrets").(*pluginsdk.Set).List()
	secrets := make([]virtualmachinescalesets.VaultSecretGroup, 0, len(secretsConfig))

	for _, secretConfig := range secretsConfig {
		config := secretConfig.(map[string]any)

		vaultSecretGroup := virtualmachinescalesets.VaultSecretGroup{
			SourceVault: &virtualmachinescalesets.SubResource{
				Id: pointer.To(config["source_vault_id"].(string)),
			},
		}

		if v := config["vault_certificates"]; v != nil {
			certsConfig := v.([]any)
			certs := make([]virtualmachinescalesets.VaultCertificate, 0, len(certsConfig))
			for _, certConfig := range certsConfig {
				config := certConfig.(map[string]any)

				cert := virtualmachinescalesets.VaultCertificate{
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

func expandAzureRMVirtualMachineScaleSetExtensions(d *pluginsdk.ResourceData) (*virtualmachinescalesets.VirtualMachineScaleSetExtensionProfile, error) {
	extensions := d.Get("extension").(*pluginsdk.Set).List()
	resources := make([]virtualmachinescalesets.VirtualMachineScaleSetExtension, 0, len(extensions))
	for _, e := range extensions {
		config := e.(map[string]any)

		extension := virtualmachinescalesets.VirtualMachineScaleSetExtension{
			Name: pointer.To(config["name"].(string)),
			Properties: &virtualmachinescalesets.VirtualMachineScaleSetExtensionProperties{
				Publisher:          pointer.To(config["publisher"].(string)),
				Type:               pointer.To(config["type"].(string)),
				TypeHandlerVersion: pointer.To(config["type_handler_version"].(string)),
			},
		}

		if u := config["auto_upgrade_minor_version"]; u != nil {
			extension.Properties.AutoUpgradeMinorVersion = pointer.To(u.(bool))
		}

		if a := config["provision_after_extensions"]; a != nil {
			provision_after_extensions := config["provision_after_extensions"].(*pluginsdk.Set).List()
			if len(provision_after_extensions) > 0 {
				var provisionAfterExtensions []string
				for _, a := range provision_after_extensions {
					str := a.(string)
					provisionAfterExtensions = append(provisionAfterExtensions, str)
				}
				extension.Properties.ProvisionAfterExtensions = &provisionAfterExtensions
			}
		}

		if s := config["settings"].(string); s != "" {
			var result any
			if err := json.Unmarshal([]byte(s), &result); err != nil {
				return nil, fmt.Errorf("unmarshaling `settings`: %+v", err)
			}
			extension.Properties.Settings = pointer.To(result)
		}

		if s := config["protected_settings"].(string); s != "" {
			var result any
			if err := json.Unmarshal([]byte(s), &result); err != nil {
				return nil, fmt.Errorf("unmarshaling `protected_settings`: %+v", err)
			}
			extension.Properties.ProtectedSettings = pointer.To(result)
		}

		resources = append(resources, extension)
	}

	return &virtualmachinescalesets.VirtualMachineScaleSetExtensionProfile{
		Extensions: &resources,
	}, nil
}

func expandAzureRmVirtualMachineScaleSetPlan(d *pluginsdk.ResourceData) *virtualmachinescalesets.Plan {
	planConfigs := d.Get("plan").(*pluginsdk.Set).List()

	planConfig := planConfigs[0].(map[string]any)

	return &virtualmachinescalesets.Plan{
		Publisher: pointer.To(planConfig["publisher"].(string)),
		Name:      pointer.To(planConfig["name"].(string)),
		Product:   pointer.To(planConfig["product"].(string)),
	}
}

func resourceVirtualMachineScaleSetExtensionHash(v any) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]any); ok {
		fmt.Fprintf(&buf, "%s-", m["name"].(string))
		fmt.Fprintf(&buf, "%s-", m["publisher"].(string))
		fmt.Fprintf(&buf, "%s-", m["type"].(string))
		fmt.Fprintf(&buf, "%s-", m["type_handler_version"].(string))

		if v, ok := m["auto_upgrade_minor_version"]; ok {
			fmt.Fprintf(&buf, "%t-", v.(bool))
		}

		if v, ok := m["provision_after_extensions"]; ok {
			fmt.Fprintf(&buf, "%s-", v.(*pluginsdk.Set).List())
		}

		// we need to ensure the whitespace is consistent
		settings := m["settings"].(string)
		if settings != "" {
			expandedSettings, err := pluginsdk.ExpandJsonFromString(settings)
			if err == nil {
				serializedSettings, err := pluginsdk.FlattenJsonToString(expandedSettings)
				if err == nil {
					fmt.Fprintf(&buf, "%s-", serializedSettings)
				}
			}
		}
	}

	return pluginsdk.HashString(buf.String())
}
