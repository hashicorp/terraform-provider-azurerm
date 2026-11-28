// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_scale_set

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2025-04-01/virtualmachinescalesets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/base64"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func flattenAzureRmVirtualMachineScaleSetOsProfileLinuxConfig(config *virtualmachinescalesets.LinuxConfiguration) []any {
	result := make(map[string]any)

	if v := config.DisablePasswordAuthentication; v != nil {
		result["disable_password_authentication"] = *v
	}

	if ssh := config.Ssh; ssh != nil {
		if keys := ssh.PublicKeys; keys != nil {
			ssh_keys := make([]map[string]any, 0, len(*keys))
			for _, i := range *keys {
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
	}

	return []any{result}
}

func flattenAzureRmVirtualMachineScaleSetOsProfileWindowsConfig(config *virtualmachinescalesets.WindowsConfiguration) []any {
	result := make(map[string]any)

	if config.ProvisionVMAgent != nil {
		result["provision_vm_agent"] = *config.ProvisionVMAgent
	}

	if config.EnableAutomaticUpdates != nil {
		result["enable_automatic_upgrades"] = *config.EnableAutomaticUpdates
	}

	if config.WinRM != nil {
		listeners := make([]map[string]any, 0, len(*config.WinRM.Listeners))
		for _, i := range *config.WinRM.Listeners {
			listener := make(map[string]any)
			listener["protocol"] = i.Protocol

			if i.CertificateURL != nil {
				listener["certificate_url"] = *i.CertificateURL
			}

			listeners = append(listeners, listener)
		}

		result["winrm"] = listeners
	}

	if config.AdditionalUnattendContent != nil {
		content := make([]map[string]any, 0, len(*config.AdditionalUnattendContent))
		for _, i := range *config.AdditionalUnattendContent {
			c := make(map[string]any)
			c["pass"] = i.PassName
			c["component"] = i.ComponentName
			c["setting_name"] = i.SettingName

			if i.Content != nil {
				c["content"] = *i.Content
			}

			content = append(content, c)
		}

		result["additional_unattend_config"] = content
	}

	return []any{result}
}

func flattenAzureRmVirtualMachineScaleSetOsProfileSecrets(secrets *[]virtualmachinescalesets.VaultSecretGroup) []map[string]any {
	result := make([]map[string]any, 0, len(*secrets))
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

func flattenAzureRmVirtualMachineScaleSetBootDiagnostics(bootDiagnostic *virtualmachinescalesets.BootDiagnostics) []any {
	b := make(map[string]any)

	if bootDiagnostic.Enabled != nil {
		b["enabled"] = *bootDiagnostic.Enabled
	}

	if bootDiagnostic.StorageUri != nil {
		b["storage_uri"] = *bootDiagnostic.StorageUri
	}

	return []any{b}
}

func flattenAzureRmVirtualMachineScaleSetRollingUpgradePolicy(rollingUpgradePolicy *virtualmachinescalesets.RollingUpgradePolicy) []any {
	b := make(map[string]any)

	if v := rollingUpgradePolicy.MaxBatchInstancePercent; v != nil {
		b["max_batch_instance_percent"] = *v
	}
	if v := rollingUpgradePolicy.MaxUnhealthyInstancePercent; v != nil {
		b["max_unhealthy_instance_percent"] = *v
	}
	if v := rollingUpgradePolicy.MaxUnhealthyUpgradedInstancePercent; v != nil {
		b["max_unhealthy_upgraded_instance_percent"] = *v
	}
	if v := rollingUpgradePolicy.PauseTimeBetweenBatches; v != nil {
		b["pause_time_between_batches"] = *v
	}

	return []any{b}
}

func flattenAzureRmVirtualMachineScaleSetNetworkProfile(profile *virtualmachinescalesets.VirtualMachineScaleSetNetworkProfile) []map[string]any {
	networkConfigurations := profile.NetworkInterfaceConfigurations
	result := make([]map[string]any, 0, len(*networkConfigurations))
	for _, netConfig := range *networkConfigurations {
		s := map[string]any{
			"name":    netConfig.Name,
			"primary": *netConfig.Properties.Primary,
		}

		if v := netConfig.Properties.EnableAcceleratedNetworking; v != nil {
			s["accelerated_networking"] = *v
		}

		if v := netConfig.Properties.EnableIPForwarding; v != nil {
			s["ip_forwarding"] = *v
		}

		if v := netConfig.Properties.NetworkSecurityGroup; v != nil {
			s["network_security_group_id"] = *v.Id
		}

		if dnsSettings := netConfig.Properties.DnsSettings; dnsSettings != nil {
			dnsServers := make([]string, 0)
			if s := dnsSettings.DnsServers; s != nil {
				dnsServers = *s
			}

			s["dns_settings"] = []any{map[string]any{
				"dns_servers": dnsServers,
			}}
		}

		if netConfig.Properties.IPConfigurations != nil {
			ipConfigs := make([]map[string]any, 0, len(netConfig.Properties.IPConfigurations))
			for _, ipConfig := range netConfig.Properties.IPConfigurations {
				config := make(map[string]any)
				config["name"] = ipConfig.Name

				if properties := ipConfig.Properties; properties != nil {
					if properties.Subnet != nil {
						config["subnet_id"] = *properties.Subnet.Id
					}

					addressPools := make([]any, 0)
					if properties.ApplicationGatewayBackendAddressPools != nil {
						for _, pool := range *properties.ApplicationGatewayBackendAddressPools {
							if v := pool.Id; v != nil {
								addressPools = append(addressPools, *v)
							}
						}
					}
					config["application_gateway_backend_address_pool_ids"] = pluginsdk.NewSet(pluginsdk.HashString, addressPools)

					applicationSecurityGroups := make([]any, 0)
					if properties.ApplicationSecurityGroups != nil {
						for _, asg := range *properties.ApplicationSecurityGroups {
							if v := asg.Id; v != nil {
								applicationSecurityGroups = append(applicationSecurityGroups, *v)
							}
						}
					}
					config["application_security_group_ids"] = pluginsdk.NewSet(pluginsdk.HashString, applicationSecurityGroups)

					if properties.LoadBalancerBackendAddressPools != nil {
						addressPools := make([]any, 0, len(*properties.LoadBalancerBackendAddressPools))
						for _, pool := range *properties.LoadBalancerBackendAddressPools {
							if v := pool.Id; v != nil {
								addressPools = append(addressPools, *v)
							}
						}
						config["load_balancer_backend_address_pool_ids"] = pluginsdk.NewSet(pluginsdk.HashString, addressPools)
					}

					if properties.LoadBalancerInboundNatPools != nil {
						inboundNatPools := make([]any, 0, len(*properties.LoadBalancerInboundNatPools))
						for _, rule := range *properties.LoadBalancerInboundNatPools {
							if v := rule.Id; v != nil {
								inboundNatPools = append(inboundNatPools, *v)
							}
						}
						config["load_balancer_inbound_nat_rules_ids"] = pluginsdk.NewSet(pluginsdk.HashString, inboundNatPools)
					}

					if properties.Primary != nil {
						config["primary"] = *properties.Primary
					}

					if publicIpInfo := properties.PublicIPAddressConfiguration; publicIpInfo != nil {
						publicIpConfigs := make([]map[string]any, 0, 1)
						publicIpConfig := make(map[string]any)
						publicIpConfig["name"] = publicIpInfo.Name
						if publicIpProperties := publicIpInfo.Properties; publicIpProperties != nil {
							if dns := publicIpProperties.DnsSettings; dns != nil {
								publicIpConfig["domain_name_label"] = dns.DomainNameLabel
							}
							if timeout := publicIpProperties.IdleTimeoutInMinutes; timeout != nil {
								publicIpConfig["idle_timeout"] = *timeout
							}
							publicIpConfigs = append(publicIpConfigs, publicIpConfig)
						}
						config["public_ip_address_configuration"] = publicIpConfigs
					}

					ipConfigs = append(ipConfigs, config)
				}
			}

			s["ip_configuration"] = ipConfigs
		}

		result = append(result, s)
	}

	return result
}

func flattenAzureRMVirtualMachineScaleSetOsProfile(d *pluginsdk.ResourceData, profile *virtualmachinescalesets.VirtualMachineScaleSetOSProfile) []any {
	result := make(map[string]any)

	result["computer_name_prefix"] = *profile.ComputerNamePrefix
	result["admin_username"] = *profile.AdminUsername

	// admin password isn't returned, so let's look it up
	if v, ok := d.GetOk("os_profile.0.admin_password"); ok {
		result["admin_password"] = v.(string)
	}

	if profile.CustomData != nil {
		result["custom_data"] = *profile.CustomData
	} else {
		// look up the current custom data
		result["custom_data"] = base64.EncodeIfNot(d.Get("os_profile.0.custom_data").(string))
	}

	return []any{result}
}

func flattenAzureRmVirtualMachineScaleSetStorageProfileOSDisk(profile *virtualmachinescalesets.VirtualMachineScaleSetOSDisk) []any {
	result := make(map[string]any)

	if profile.Name != nil {
		result["name"] = *profile.Name
	}

	if profile.Image != nil {
		result["image"] = *profile.Image.Uri
	}

	containers := make([]any, 0)
	if profile.VhdContainers != nil {
		for _, container := range *profile.VhdContainers {
			containers = append(containers, container)
		}
	}
	result["vhd_containers"] = pluginsdk.NewSet(pluginsdk.HashString, containers)

	if profile.ManagedDisk != nil {
		result["managed_disk_type"] = pointer.FromEnum(profile.ManagedDisk.StorageAccountType)
	}

	result["caching"] = profile.Caching
	result["create_option"] = profile.CreateOption
	result["os_type"] = profile.OsType

	return []any{result}
}

func flattenAzureRmVirtualMachineScaleSetStorageProfileDataDisk(disks *[]virtualmachinescalesets.VirtualMachineScaleSetDataDisk) any {
	result := make([]any, len(*disks))
	for i, disk := range *disks {
		l := make(map[string]any)
		if disk.ManagedDisk != nil {
			l["managed_disk_type"] = pointer.FromEnum(disk.ManagedDisk.StorageAccountType)
		}

		l["create_option"] = disk.CreateOption
		l["caching"] = pointer.FromEnum(disk.Caching)
		if disk.DiskSizeGB != nil {
			l["disk_size_gb"] = *disk.DiskSizeGB
		}

		l["lun"] = disk.Lun

		result[i] = l
	}
	return result
}

func flattenAzureRmVirtualMachineScaleSetStorageProfileImageReference(image *virtualmachinescalesets.ImageReference) []any {
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

func flattenAzureRmVirtualMachineScaleSetSku(sku *virtualmachinescalesets.Sku) []any {
	result := make(map[string]any)
	result["name"] = *sku.Name
	result["capacity"] = *sku.Capacity

	if *sku.Tier != "" {
		result["tier"] = *sku.Tier
	}

	return []any{result}
}

func flattenAzureRmVirtualMachineScaleSetExtensionProfile(profile *virtualmachinescalesets.VirtualMachineScaleSetExtensionProfile) ([]map[string]any, error) {
	if profile.Extensions == nil {
		return []map[string]any{}, nil
	}

	result := make([]map[string]any, 0, len(*profile.Extensions))
	for _, extension := range *profile.Extensions {
		e := make(map[string]any)
		e["name"] = *extension.Name
		properties := extension.Properties
		if properties != nil {
			e["publisher"] = *properties.Publisher
			e["type"] = *properties.Type
			e["type_handler_version"] = *properties.TypeHandlerVersion
			if properties.AutoUpgradeMinorVersion != nil {
				e["auto_upgrade_minor_version"] = *properties.AutoUpgradeMinorVersion
			}

			provisionAfterExtensions := make([]any, 0)
			if properties.ProvisionAfterExtensions != nil {
				for _, provisionAfterExtension := range *properties.ProvisionAfterExtensions {
					provisionAfterExtensions = append(provisionAfterExtensions, provisionAfterExtension)
				}
			}
			e["provision_after_extensions"] = pluginsdk.NewSet(pluginsdk.HashString, provisionAfterExtensions)

			if settings := properties.Settings; settings != nil {
				result, err := json.Marshal(settings)
				if err != nil {
					return nil, fmt.Errorf("unmarshaling `settings`: %+v", err)
				}
				e["settings"] = string(result)
			}
		}

		result = append(result, e)
	}

	return result, nil
}

func flattenAzureRmVirtualMachineScaleSetPlan(plan *virtualmachinescalesets.Plan) []any {
	result := make(map[string]any)

	result["name"] = *plan.Name
	result["publisher"] = *plan.Publisher
	result["product"] = *plan.Product

	return []any{result}
}
