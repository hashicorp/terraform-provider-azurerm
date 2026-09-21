// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package pool

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/batch/2024-07-01/pool"
	"github.com/hashicorp/terraform-provider-azurerm/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

// flattenBatchPoolAutoScaleSettings flattens the auto scale settings for a Batch pool
func flattenBatchPoolAutoScaleSettings(settings *pool.AutoScaleSettings) []interface{} {
	results := make([]interface{}, 0)

	if settings == nil {
		log.Printf("[DEBUG] settings is nil")
		return results
	}

	result := make(map[string]interface{})

	if settings.EvaluationInterval != nil {
		result["evaluation_interval"] = *settings.EvaluationInterval
	}

	result["formula"] = settings.Formula

	return append(results, result)
}

// flattenBatchPoolFixedScaleSettings flattens the fixed scale settings for a Batch pool
func flattenBatchPoolFixedScaleSettings(d *pluginsdk.ResourceData, settings *pool.FixedScaleSettings) []interface{} {
	results := make([]interface{}, 0)

	if settings == nil {
		log.Printf("[DEBUG] settings is nil")
		return results
	}

	result := make(map[string]interface{})

	// for now, this is a writeOnly property, so we treat this as secret.
	if v, ok := d.GetOk("fixed_scale.0.node_deallocation_method"); ok {
		result["node_deallocation_method"] = v.(string)
	}

	if settings.TargetDedicatedNodes != nil {
		result["target_dedicated_nodes"] = *settings.TargetDedicatedNodes
	}

	if settings.TargetLowPriorityNodes != nil {
		result["target_low_priority_nodes"] = *settings.TargetLowPriorityNodes
	}

	if settings.ResizeTimeout != nil {
		result["resize_timeout"] = *settings.ResizeTimeout
	}

	return append(results, result)
}

// flattenBatchPoolImageReference flattens the Batch pool image reference
func flattenBatchPoolImageReference(image *pool.ImageReference) []interface{} {
	results := make([]interface{}, 0)
	if image == nil {
		log.Printf("[DEBUG] image is nil")
		return results
	}

	result := make(map[string]interface{})
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

	return append(results, result)
}

// flattenBatchPoolStartTask flattens a Batch pool start task
func flattenBatchPoolStartTask(oldConfig *pluginsdk.ResourceData, startTask *pool.StartTask) []interface{} {
	results := make([]interface{}, 0)

	if startTask == nil {
		log.Printf("[DEBUG] startTask is nil")
		return results
	}

	result := make(map[string]interface{})
	result["command_line"] = pointer.From(startTask.CommandLine)

	if startTask.ContainerSettings != nil {
		containerSettings := make(map[string]interface{})
		containerSettings["image_name"] = startTask.ContainerSettings.ImageName
		if startTask.ContainerSettings.WorkingDirectory != nil {
			containerSettings["working_directory"] = string(*startTask.ContainerSettings.WorkingDirectory)
		}
		if startTask.ContainerSettings.ContainerRunOptions != nil {
			containerSettings["run_options"] = *startTask.ContainerSettings.ContainerRunOptions
		}
		if startTask.ContainerSettings.Registry != nil {
			tmpReg := flattenBatchPoolContainerRegistry(oldConfig, startTask.ContainerSettings.Registry)
			containerSettings["registry"] = []interface{}{
				tmpReg,
			}
		}

		result["container"] = []interface{}{
			containerSettings,
		}
	}

	result["wait_for_success"] = pointer.From(startTask.WaitForSuccess)

	result["task_retry_maximum"] = pointer.From(startTask.MaxTaskRetryCount)

	if startTask.UserIdentity != nil {
		userIdentity := make(map[string]interface{})
		if startTask.UserIdentity.AutoUser != nil {
			autoUser := make(map[string]interface{})

			if startTask.UserIdentity.AutoUser.ElevationLevel != nil {
				autoUser["elevation_level"] = string(*startTask.UserIdentity.AutoUser.ElevationLevel)
			}

			if startTask.UserIdentity.AutoUser.Scope != nil {
				autoUser["scope"] = string(*startTask.UserIdentity.AutoUser.Scope)
			}
			userIdentity["auto_user"] = []interface{}{autoUser}
		} else {
			userIdentity["user_name"] = *startTask.UserIdentity.UserName
		}

		result["user_identity"] = []interface{}{userIdentity}
	}

	resourceFiles := make([]interface{}, 0)
	if startTask.ResourceFiles != nil {
		for _, armResourceFile := range *startTask.ResourceFiles {
			resourceFile := make(map[string]interface{})
			if armResourceFile.AutoStorageContainerName != nil {
				resourceFile["auto_storage_container_name"] = *armResourceFile.AutoStorageContainerName
			}
			if armResourceFile.StorageContainerURL != nil {
				resourceFile["storage_container_url"] = *armResourceFile.StorageContainerURL
			}
			if armResourceFile.HTTPURL != nil {
				resourceFile["http_url"] = *armResourceFile.HTTPURL
			}
			if armResourceFile.BlobPrefix != nil {
				resourceFile["blob_prefix"] = *armResourceFile.BlobPrefix
			}
			if armResourceFile.FilePath != nil {
				resourceFile["file_path"] = *armResourceFile.FilePath
			}
			if armResourceFile.FileMode != nil {
				resourceFile["file_mode"] = *armResourceFile.FileMode
			}
			if armResourceFile.IdentityReference != nil {
				resourceFile["user_assigned_identity_id"] = *armResourceFile.IdentityReference.ResourceId
			}
			resourceFiles = append(resourceFiles, resourceFile)
		}
	}

	environment := make(map[string]interface{})
	if startTask.EnvironmentSettings != nil {
		for _, envSetting := range *startTask.EnvironmentSettings {
			environment[envSetting.Name] = *envSetting.Value
		}
	}

	result["common_environment_properties"] = environment

	result["resource_file"] = resourceFiles

	return append(results, result)
}

// flattenBatchPoolContainerConfiguration flattens a Batch pool container configuration
func flattenBatchPoolContainerConfiguration(d *pluginsdk.ResourceData, armContainerConfiguration *pool.ContainerConfiguration) interface{} {
	result := make(map[string]interface{})

	if armContainerConfiguration == nil {
		return nil
	}

	result["type"] = armContainerConfiguration.Type

	names := &pluginsdk.Set{F: pluginsdk.HashString}
	if armContainerConfiguration.ContainerImageNames != nil {
		for _, armName := range *armContainerConfiguration.ContainerImageNames {
			names.Add(armName)
		}
	}
	result["container_image_names"] = names

	result["container_registries"] = flattenBatchPoolContainerRegistries(d, armContainerConfiguration.ContainerRegistries)

	return []interface{}{result}
}

func flattenBatchPoolContainerRegistries(d *pluginsdk.ResourceData, armContainerRegistries *[]pool.ContainerRegistry) []interface{} {
	results := make([]interface{}, 0)

	if armContainerRegistries == nil {
		return results
	}

	for _, armContainerRegistry := range *armContainerRegistries {
		result := flattenBatchPoolContainerRegistry(d, &armContainerRegistry)
		results = append(results, result)
	}

	return results
}

func flattenBatchPoolContainerRegistry(d *pluginsdk.ResourceData, armContainerRegistry *pool.ContainerRegistry) map[string]interface{} {
	result := make(map[string]interface{})

	if armContainerRegistry == nil {
		return result
	}

	if registryServer := armContainerRegistry.RegistryServer; registryServer != nil {
		result["registry_server"] = *registryServer
	}

	if userName := armContainerRegistry.Username; userName != nil {
		result["user_name"] = *userName
		// Locate the password only if user_name is defined
		result["password"] = findBatchPoolContainerRegistryPassword(d, result["registry_server"].(string), result["user_name"].(string))
	}

	if identity := armContainerRegistry.IdentityReference; identity != nil {
		if identity.ResourceId != nil {
			result["user_assigned_identity_id"] = *identity.ResourceId
		}
	}

	return result
}

func findBatchPoolContainerRegistryPassword(d *pluginsdk.ResourceData, armServer string, armUsername string) interface{} {
	var numContainerRegistries int
	if n, ok := d.GetOk("container_configuration.0.container_registries.#"); ok {
		numContainerRegistries = n.(int)
	} else {
		return ""
	}

	for i := 0; i < numContainerRegistries; i++ {
		if server, ok := d.GetOk(fmt.Sprintf("container_configuration.0.container_registries.%d.registry_server", i)); !ok || server != armServer {
			continue
		}
		if username, ok := d.GetOk(fmt.Sprintf("container_configuration.0.container_registries.%d.user_name", i)); !ok || username != armUsername {
			continue
		}
		return d.Get(fmt.Sprintf("container_configuration.0.container_registries.%d.password", i))
	}

	return ""
}

func findSensitiveInfoForMountConfig(targetType string, sourceType string, sourceValue string, mountType string, d *pluginsdk.ResourceData) string {
	if num, ok := d.GetOk("mount.#"); ok {
		n := num.(int)
		for i := range n {
			if src, ok := d.GetOk(fmt.Sprintf("mount.%d.%v.0.%v", i, mountType, sourceType)); ok && src == sourceValue {
				return d.Get(fmt.Sprintf("mount.%d.%v.0.%v", i, mountType, targetType)).(string)
			}
		}
	}
	return ""
}

func flattenBatchPoolMountConfig(d *pluginsdk.ResourceData, config *pool.MountConfiguration) map[string]interface{} {
	mountConfig := make(map[string]interface{})

	switch {
	case config.AzureBlobFileSystemConfiguration != nil:
		azureBlobFileSysConfigList := make([]interface{}, 0)
		azureBlobFileSysConfig := make(map[string]interface{})
		azureBlobFileSysConfig["account_name"] = config.AzureBlobFileSystemConfiguration.AccountName
		azureBlobFileSysConfig["container_name"] = config.AzureBlobFileSystemConfiguration.ContainerName
		azureBlobFileSysConfig["relative_mount_path"] = config.AzureBlobFileSystemConfiguration.RelativeMountPath
		azureBlobFileSysConfig["account_key"] = findSensitiveInfoForMountConfig("account_key", "account_name", config.AzureBlobFileSystemConfiguration.AccountName, "azure_blob_file_system", d)
		azureBlobFileSysConfig["sas_key"] = findSensitiveInfoForMountConfig("sas_key", "account_name", config.AzureBlobFileSystemConfiguration.AccountName, "azure_blob_file_system", d)
		if config.AzureBlobFileSystemConfiguration.IdentityReference != nil {
			azureBlobFileSysConfig["identity_id"] = flattenBatchPoolIdentityReferenceToIdentityID(config.AzureBlobFileSystemConfiguration.IdentityReference)
		}
		if config.AzureBlobFileSystemConfiguration.BlobfuseOptions != nil {
			azureBlobFileSysConfig["blobfuse_options"] = *config.AzureBlobFileSystemConfiguration.BlobfuseOptions
		}
		azureBlobFileSysConfigList = append(azureBlobFileSysConfigList, azureBlobFileSysConfig)
		mountConfig["azure_blob_file_system"] = azureBlobFileSysConfigList
	case config.AzureFileShareConfiguration != nil:
		azureFileShareConfigList := make([]interface{}, 0)
		azureFileShareConfig := make(map[string]interface{})
		azureFileShareConfig["account_name"] = config.AzureFileShareConfiguration.AccountName
		azureFileShareConfig["azure_file_url"] = config.AzureFileShareConfiguration.AzureFileURL
		azureFileShareConfig["account_key"] = findSensitiveInfoForMountConfig("account_key", "account_name", config.AzureFileShareConfiguration.AccountName, "azure_file_share", d)
		azureFileShareConfig["relative_mount_path"] = config.AzureFileShareConfiguration.RelativeMountPath

		if config.AzureFileShareConfiguration.MountOptions != nil {
			azureFileShareConfig["mount_options"] = *config.AzureFileShareConfiguration.MountOptions
		}

		azureFileShareConfigList = append(azureFileShareConfigList, azureFileShareConfig)
		mountConfig["azure_file_share"] = azureFileShareConfigList

	case config.CifsMountConfiguration != nil:
		cifsMountConfigList := make([]interface{}, 0)
		cifsMountConfig := make(map[string]interface{})

		cifsMountConfig["user_name"] = config.CifsMountConfiguration.UserName
		cifsMountConfig["password"] = findSensitiveInfoForMountConfig("password", "user_name", config.CifsMountConfiguration.UserName, "cifs_mount", d)
		cifsMountConfig["source"] = config.CifsMountConfiguration.Source
		cifsMountConfig["relative_mount_path"] = config.CifsMountConfiguration.RelativeMountPath

		if config.CifsMountConfiguration.MountOptions != nil {
			cifsMountConfig["mount_options"] = *config.CifsMountConfiguration.MountOptions
		}

		cifsMountConfigList = append(cifsMountConfigList, cifsMountConfig)
		mountConfig["cifs_mount"] = cifsMountConfigList
	case config.NfsMountConfiguration != nil:
		nfsMountConfigList := make([]interface{}, 0)
		nfsMountConfig := make(map[string]interface{})

		nfsMountConfig["source"] = config.NfsMountConfiguration.Source
		nfsMountConfig["relative_mount_path"] = config.NfsMountConfiguration.RelativeMountPath

		if config.NfsMountConfiguration.MountOptions != nil {
			nfsMountConfig["mount_options"] = *config.NfsMountConfiguration.MountOptions
		}

		nfsMountConfigList = append(nfsMountConfigList, nfsMountConfig)
		mountConfig["nfs_mount"] = nfsMountConfigList
	default:
		return map[string]interface{}{}
	}

	return mountConfig
}

func flattenBatchPoolIdentityReferenceToIdentityID(ref *pool.ComputeNodeIdentityReference) string {
	if ref != nil && ref.ResourceId != nil {
		return *ref.ResourceId
	}
	return ""
}

func flattenBatchPoolSecurityProfile(configProfile *pool.SecurityProfile) []interface{} {
	securityProfile := make([]interface{}, 0)
	securityConfig := make(map[string]interface{})

	securityConfig["host_encryption_enabled"] = pointer.From(configProfile.EncryptionAtHost)
	securityConfig["security_type"] = string(pointer.From(configProfile.SecurityType))

	if configProfile.UefiSettings != nil {
		securityConfig["secure_boot_enabled"] = pointer.From(configProfile.UefiSettings.SecureBootEnabled)
		securityConfig["vtpm_enabled"] = pointer.From(configProfile.UefiSettings.VTpmEnabled)
	}

	securityProfile = append(securityProfile, securityConfig)
	return securityProfile
}

func flattenBatchPoolUserAccount(d *pluginsdk.ResourceData, account *pool.UserAccount) map[string]interface{} {
	userAccount := make(map[string]interface{})
	userAccount["name"] = account.Name
	if account.ElevationLevel != nil {
		userAccount["elevation_level"] = string(*account.ElevationLevel)
	}
	userAccountIndex := -1

	if num, ok := d.GetOk("user_accounts.#"); ok {
		n := num.(int)
		for i := range n {
			if src, nameOk := d.GetOk(fmt.Sprintf("user_accounts.%d.name", i)); nameOk && src == account.Name {
				userAccount["password"] = d.Get(fmt.Sprintf("user_accounts.%d.password", i)).(string)
				userAccountIndex = i
				break
			}
		}
	}

	if account.LinuxUserConfiguration != nil {
		linuxUserConfig := make(map[string]interface{})

		if account.LinuxUserConfiguration.Uid != nil {
			linuxUserConfig["uid"] = *account.LinuxUserConfiguration.Uid
			linuxUserConfig["gid"] = *account.LinuxUserConfiguration.Gid
		}

		if userAccountIndex > -1 {
			if sshPrivateKey, ok := d.GetOk(fmt.Sprintf("user_accounts.%d.linux_user_configuration.0.ssh_private_key", userAccountIndex)); ok {
				linuxUserConfig["ssh_private_key"] = sshPrivateKey
			}
		}

		userAccount["linux_user_configuration"] = []interface{}{
			linuxUserConfig,
		}
	}

	if account.WindowsUserConfiguration != nil {
		loginMode := make(map[string]interface{})
		if account.WindowsUserConfiguration.LoginMode != nil {
			loginMode["login_mode"] = string(*account.WindowsUserConfiguration.LoginMode)
		}
		userAccount["windows_user_configuration"] = []interface{}{
			loginMode,
		}
	}
	return userAccount
}

// FlattenBatchMetaData flattens a Batch pool metadata
func FlattenBatchMetaData(metadatas *[]pool.MetadataItem) map[string]interface{} {
	output := make(map[string]interface{})

	if metadatas == nil {
		return output
	}

	for _, metadata := range *metadatas {
		output[metadata.Name] = metadata.Value
	}

	return output
}

func flattenBatchPoolNetworkConfiguration(input *pool.NetworkConfiguration) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	publicIPAddressIds := make([]interface{}, 0)
	publicAddressProvisioningType := ""
	if config := input.PublicIPAddressConfiguration; config != nil {
		publicIPAddressIds = helpers.FlattenStringSlice(config.IPAddressIds)
		if config.Provision != nil {
			publicAddressProvisioningType = string(*config.Provision)
		}
	}

	endpointConfigs := make([]interface{}, 0)
	if config := input.EndpointConfiguration; config != nil && config.InboundNatPools != nil {
		for _, inboundNatPool := range config.InboundNatPools {
			name := inboundNatPool.Name

			backendPort := inboundNatPool.BackendPort

			frontendPortRange := fmt.Sprintf("%d-%d", inboundNatPool.FrontendPortRangeStart, inboundNatPool.FrontendPortRangeEnd)

			networkSecurities := make([]interface{}, 0)
			if sgRules := inboundNatPool.NetworkSecurityGroupRules; sgRules != nil {
				for _, networkSecurity := range *sgRules {
					priority := networkSecurity.Priority

					sourceAddressPrefix := networkSecurity.SourceAddressPrefix

					sourcePortRanges := make([]interface{}, 0)
					if networkSecurity.SourcePortRanges != nil {
						for _, sourcePortRange := range *networkSecurity.SourcePortRanges {
							sourcePortRanges = append(sourcePortRanges, sourcePortRange)
						}
					}
					networkSecurities = append(networkSecurities, map[string]interface{}{
						"access":                string(networkSecurity.Access),
						"priority":              priority,
						"source_address_prefix": sourceAddressPrefix,
						"source_port_ranges":    sourcePortRanges,
					})
				}
			}

			endpointConfigs = append(endpointConfigs, map[string]interface{}{
				"backend_port":                 backendPort,
				"frontend_port_range":          frontendPortRange,
				"name":                         name,
				"network_security_group_rules": networkSecurities,
				"protocol":                     string(inboundNatPool.Protocol),
			})
		}
	}

	dynamicVNetAssignmentScope := ""
	if input.DynamicVnetAssignmentScope != nil {
		dynamicVNetAssignmentScope = string(*input.DynamicVnetAssignmentScope)
	}

	return []interface{}{
		map[string]interface{}{
			"dynamic_vnet_assignment_scope":    dynamicVNetAssignmentScope,
			"accelerated_networking_enabled":   pointer.From(input.EnableAcceleratedNetworking),
			"endpoint_configuration":           endpointConfigs,
			"public_address_provisioning_type": publicAddressProvisioningType,
			"public_ips":                       pluginsdk.NewSet(pluginsdk.HashString, publicIPAddressIds),
			"subnet_id":                        pointer.From(input.SubnetId),
		},
	}
}
