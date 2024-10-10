// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package batch

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/batch/2023-05-01/pool"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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
	commandLine := ""
	if startTask.CommandLine != nil {
		commandLine = *startTask.CommandLine
	}
	result["command_line"] = commandLine

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

	waitForSuccess := false
	if startTask.WaitForSuccess != nil {
		waitForSuccess = *startTask.WaitForSuccess
	}
	result["wait_for_success"] = waitForSuccess

	maxTaskRetryCount := int64(0)
	if startTask.MaxTaskRetryCount != nil {
		maxTaskRetryCount = *startTask.MaxTaskRetryCount
	}

	result["task_retry_maximum"] = maxTaskRetryCount

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

// flattenBatchPoolCertificateReferences flattens a Batch pool certificate reference
func flattenBatchPoolCertificateReferences(armCertificates *[]pool.CertificateReference) []interface{} {
	if armCertificates == nil {
		return []interface{}{}
	}
	output := make([]interface{}, 0)

	for _, armCertificate := range *armCertificates {
		certificate := map[string]interface{}{}

		certificate["id"] = armCertificate.Id
		if armCertificate.StoreLocation != nil {
			certificate["store_location"] = string(*armCertificate.StoreLocation)
		}
		if armCertificate.StoreName != nil {
			certificate["store_name"] = *armCertificate.StoreName
		}
		visibility := &pluginsdk.Set{F: pluginsdk.HashString}
		if armCertificate.Visibility != nil {
			for _, armVisibility := range *armCertificate.Visibility {
				visibility.Add(string(armVisibility))
			}
		}
		certificate["visibility"] = visibility
		output = append(output, certificate)
	}
	return output
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
	numContainerRegistries := 0
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
		for i := 0; i < n; i++ {
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
		return nil
	}

	return mountConfig
}

func flattenBatchPoolIdentityReferenceToIdentityID(ref *pool.ComputeNodeIdentityReference) string {
	if ref != nil && ref.ResourceId != nil {
		return *ref.ResourceId
	}
	return ""
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
		for i := 0; i < n; i++ {
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

// ExpandBatchPoolImageReference expands Batch pool image reference
func ExpandBatchPoolImageReference(list []interface{}) (*pool.ImageReference, error) {
	if len(list) == 0 || list[0] == nil {
		return nil, fmt.Errorf("storage image reference should be defined")
	}

	storageImageRef := list[0].(map[string]interface{})
	imageRef := &pool.ImageReference{}

	if storageImageRef["id"] != nil && storageImageRef["id"] != "" {
		storageImageRefID := storageImageRef["id"].(string)
		imageRef.Id = &storageImageRefID
	}

	if storageImageRef["offer"] != nil && storageImageRef["offer"] != "" {
		storageImageRefOffer := storageImageRef["offer"].(string)
		imageRef.Offer = &storageImageRefOffer
	}

	if storageImageRef["publisher"] != nil && storageImageRef["publisher"] != "" {
		storageImageRefPublisher := storageImageRef["publisher"].(string)
		imageRef.Publisher = &storageImageRefPublisher
	}

	if storageImageRef["sku"] != nil && storageImageRef["sku"] != "" {
		storageImageRefSku := storageImageRef["sku"].(string)
		imageRef.Sku = &storageImageRefSku
	}

	if storageImageRef["version"] != nil && storageImageRef["version"] != "" {
		storageImageRefVersion := storageImageRef["version"].(string)
		imageRef.Version = &storageImageRefVersion
	}

	return imageRef, nil
}

// ExpandBatchPoolContainerConfiguration expands the Batch pool container configuration
func ExpandBatchPoolContainerConfiguration(list []interface{}) (*pool.ContainerConfiguration, error) {
	if len(list) == 0 || list[0] == nil {
		return nil, nil
	}

	block := list[0].(map[string]interface{})

	containerRegistries, err := expandBatchPoolContainerRegistries(block["container_registries"].([]interface{}))
	if err != nil {
		return nil, err
	}

	obj := &pool.ContainerConfiguration{
		Type:                pool.ContainerType(block["type"].(string)),
		ContainerRegistries: containerRegistries,
		ContainerImageNames: utils.ExpandStringSlice(block["container_image_names"].(*pluginsdk.Set).List()),
	}

	return obj, nil
}

func expandBatchPoolContainerRegistries(list []interface{}) (*[]pool.ContainerRegistry, error) {
	result := []pool.ContainerRegistry{}

	for _, tempItem := range list {
		item := tempItem.(map[string]interface{})
		containerRegistry, err := expandBatchPoolContainerRegistry(item)
		if err != nil {
			return nil, err
		}
		result = append(result, *containerRegistry)
	}
	return &result, nil
}

func expandBatchPoolContainerRegistry(ref map[string]interface{}) (*pool.ContainerRegistry, error) {
	if len(ref) == 0 {
		return nil, fmt.Errorf("container registry reference should be defined")
	}

	containerRegistry := pool.ContainerRegistry{}

	if v := ref["registry_server"]; v != nil && v != "" {
		containerRegistry.RegistryServer = pointer.FromString(v.(string))
	}
	if v := ref["user_name"]; v != nil && v != "" {
		containerRegistry.Username = pointer.FromString(v.(string))
	}
	if v := ref["password"]; v != nil && v != "" {
		containerRegistry.Password = pointer.FromString(v.(string))
	}
	if v := ref["user_assigned_identity_id"]; v != nil && v != "" {
		containerRegistry.IdentityReference = &pool.ComputeNodeIdentityReference{
			ResourceId: pointer.FromString(v.(string)),
		}
	}

	return &containerRegistry, nil
}

// ExpandBatchPoolCertificateReferences expands Batch pool certificate references
func ExpandBatchPoolCertificateReferences(list []interface{}) (*[]pool.CertificateReference, error) {
	var result []pool.CertificateReference

	for _, tempItem := range list {
		item := tempItem.(map[string]interface{})
		certificateReference, err := expandBatchPoolCertificateReference(item)
		if err != nil {
			return nil, err
		}
		result = append(result, *certificateReference)
	}
	return &result, nil
}

func expandBatchPoolCertificateReference(ref map[string]interface{}) (*pool.CertificateReference, error) {
	if len(ref) == 0 {
		return nil, fmt.Errorf("Error: storage image reference should be defined")
	}

	id := ref["id"].(string)
	storeLocation := ref["store_location"].(string)
	storeName := ref["store_name"].(string)
	visibilityRefs := ref["visibility"].(*pluginsdk.Set)
	var visibility []pool.CertificateVisibility
	if visibilityRefs != nil {
		for _, visibilityRef := range visibilityRefs.List() {
			visibility = append(visibility, pool.CertificateVisibility(visibilityRef.(string)))
		}
	}

	certificateReference := &pool.CertificateReference{
		Id:            id,
		StoreLocation: pointer.To(pool.CertificateStoreLocation(storeLocation)),
		StoreName:     &storeName,
		Visibility:    &visibility,
	}
	return certificateReference, nil
}

// ExpandBatchPoolStartTask expands Batch pool start task
func ExpandBatchPoolStartTask(list []interface{}) (*pool.StartTask, error) {
	if len(list) == 0 || list[0] == nil {
		return nil, fmt.Errorf("batch pool start task should be defined")
	}

	startTaskValue := list[0].(map[string]interface{})

	startTaskCmdLine := startTaskValue["command_line"].(string)

	maxTaskRetryCount := int64(1)

	if v := startTaskValue["task_retry_maximum"].(int); v > 0 {
		maxTaskRetryCount = int64(v)
	}

	waitForSuccess := startTaskValue["wait_for_success"].(bool)

	userIdentityList := startTaskValue["user_identity"].([]interface{})
	if len(userIdentityList) == 0 {
		return nil, fmt.Errorf("batch pool start task user identity should be defined")
	}

	userIdentityValue := userIdentityList[0].(map[string]interface{})
	userIdentity := pool.UserIdentity{}

	if autoUserValue, ok := userIdentityValue["auto_user"]; ok {
		autoUser := autoUserValue.([]interface{})
		if len(autoUser) != 0 {
			autoUserMap := autoUser[0].(map[string]interface{})
			userIdentity.AutoUser = &pool.AutoUserSpecification{
				ElevationLevel: pointer.To(pool.ElevationLevel(autoUserMap["elevation_level"].(string))),
				Scope:          pointer.To(pool.AutoUserScope(autoUserMap["scope"].(string))),
			}
		}
	}
	if userNameValue, ok := userIdentityValue["user_name"]; ok {
		userName := userNameValue.(string)
		if len(userName) != 0 {
			userIdentity.UserName = &userName
		}
	}

	resourceFileList := startTaskValue["resource_file"].([]interface{})
	resourceFiles := make([]pool.ResourceFile, 0)
	for _, resourceFileValueTemp := range resourceFileList {
		if resourceFileValueTemp == nil {
			continue
		}
		resourceFileValue := resourceFileValueTemp.(map[string]interface{})
		resourceFile := pool.ResourceFile{}
		if v, ok := resourceFileValue["auto_storage_container_name"]; ok {
			autoStorageContainerName := v.(string)
			if autoStorageContainerName != "" {
				resourceFile.AutoStorageContainerName = &autoStorageContainerName
			}
		}
		if v, ok := resourceFileValue["storage_container_url"]; ok {
			storageContainerURL := v.(string)
			if storageContainerURL != "" {
				resourceFile.StorageContainerURL = &storageContainerURL
			}
		}
		if v, ok := resourceFileValue["http_url"]; ok {
			httpURL := v.(string)
			if httpURL != "" {
				resourceFile.HTTPURL = &httpURL
			}
		}
		if v, ok := resourceFileValue["blob_prefix"]; ok {
			blobPrefix := v.(string)
			if blobPrefix != "" {
				resourceFile.BlobPrefix = &blobPrefix
			}
		}
		if v, ok := resourceFileValue["file_path"]; ok {
			filePath := v.(string)
			if filePath != "" {
				resourceFile.FilePath = &filePath
			}
		}
		if v, ok := resourceFileValue["file_mode"]; ok {
			fileMode := v.(string)
			if fileMode != "" {
				resourceFile.FileMode = &fileMode
			}
		}
		if v, ok := resourceFileValue["user_assigned_identity_id"]; ok {
			resourceId := v.(string)
			if resourceId != "" {
				identityReference := pool.ComputeNodeIdentityReference{
					ResourceId: utils.String(resourceId),
				}
				resourceFile.IdentityReference = &identityReference
			}
		}
		resourceFiles = append(resourceFiles, resourceFile)
	}

	startTask := &pool.StartTask{
		CommandLine:       &startTaskCmdLine,
		MaxTaskRetryCount: &maxTaskRetryCount,
		WaitForSuccess:    &waitForSuccess,
		UserIdentity:      &userIdentity,
		ResourceFiles:     &resourceFiles,
	}

	if v := startTaskValue["common_environment_properties"].(map[string]interface{}); len(v) > 0 {
		startTask.EnvironmentSettings = expandCommonEnvironmentProperties(v)
	}

	if startTaskValue["container"] != nil && len(startTaskValue["container"].([]interface{})) > 0 {
		var containerSettings pool.TaskContainerSettings
		containerSettingsList := startTaskValue["container"].([]interface{})

		if len(containerSettingsList) > 0 && containerSettingsList[0] != nil {
			settingMap := containerSettingsList[0].(map[string]interface{})
			containerSettings.ImageName = settingMap["image_name"].(string)
			if containerRunOptions, ok := settingMap["run_options"]; ok {
				containerSettings.ContainerRunOptions = utils.String(containerRunOptions.(string))
			}
			if registries, ok := settingMap["registry"].([]interface{}); ok && len(registries) > 0 && registries[0] != nil {
				containerRegMap := registries[0].(map[string]interface{})
				if containerRegistryRef, err := expandBatchPoolContainerRegistry(containerRegMap); err == nil {
					containerSettings.Registry = containerRegistryRef
				}
			}
			if workingDir, ok := settingMap["working_directory"]; ok {
				containerSettings.WorkingDirectory = pointer.To(pool.ContainerWorkingDirectory(workingDir.(string)))
			}
		}
		startTask.ContainerSettings = &containerSettings
	}

	return startTask, nil
}

func expandBatchPoolVirtualMachineConfig(d *pluginsdk.ResourceData) (*pool.VirtualMachineConfiguration, error) {
	var result pool.VirtualMachineConfiguration

	result.NodeAgentSkuId = d.Get("node_agent_sku_id").(string)

	storageImageReferenceSet := d.Get("storage_image_reference").([]interface{})
	if imageReference, err := ExpandBatchPoolImageReference(storageImageReferenceSet); err == nil {
		if imageReference != nil {
			// if an image reference ID is specified, the user wants use a custom image. This property is mutually exclusive with other properties.
			if imageReference.Id != nil && (imageReference.Offer != nil || imageReference.Publisher != nil || imageReference.Sku != nil || imageReference.Version != nil) {
				return nil, fmt.Errorf("properties version, offer, publish cannot be defined when using a custom image id")
			} else if imageReference.Id == nil && (imageReference.Offer == nil || imageReference.Publisher == nil || imageReference.Sku == nil || imageReference.Version == nil) {
				return nil, fmt.Errorf("properties version, offer, publish and sku are mandatory when not using a custom image")
			}
			result.ImageReference = *imageReference
		}
	} else {
		return nil, fmt.Errorf("storage_image_reference either is empty or contains parsing errors")
	}

	if v, ok := d.GetOk("container_configuration"); ok {
		if containerConfiguration, err := ExpandBatchPoolContainerConfiguration(v.([]interface{})); err == nil {
			result.ContainerConfiguration = containerConfiguration
		} else {
			return nil, fmt.Errorf("container_configuration either is empty or contains parsing errors")
		}
	}

	if v, ok := d.GetOk("data_disks"); ok {
		result.DataDisks = expandBatchPoolDataDisks(v.([]interface{}))
	}

	if diskEncryptionConfig, diskEncryptionErr := expandBatchPoolDiskEncryptionConfiguration(d.Get("disk_encryption").([]interface{})); diskEncryptionErr == nil {
		result.DiskEncryptionConfiguration = diskEncryptionConfig
	} else {
		return nil, diskEncryptionErr
	}

	if extensions, extErr := expandBatchPoolExtensions(d.Get("extensions").([]interface{})); extErr == nil {
		result.Extensions = extensions
	} else {
		return nil, extErr
	}

	if licenseType, ok := d.GetOk("license_type"); ok {
		result.LicenseType = utils.String(licenseType.(string))
	}

	if v, ok := d.GetOk("node_placement"); ok {
		result.NodePlacementConfiguration = expandBatchPoolNodeReplacementConfig(v.([]interface{}))
	}

	if v, ok := d.GetOk("os_disk_placement"); ok {
		result.OsDisk = expandBatchPoolOSDisk(v)
	}

	if v, ok := d.GetOk("windows"); ok {
		result.WindowsConfiguration = expandBatchPoolWindowsConfiguration(v.([]interface{}))
	}

	return &result, nil
}

func expandBatchPoolOSDisk(ref interface{}) *pool.OSDisk {
	if ref == nil {
		return nil
	}

	return &pool.OSDisk{
		EphemeralOSDiskSettings: &pool.DiffDiskSettings{
			Placement: pointer.To(pool.DiffDiskPlacement(ref.(string))),
		},
	}
}

func expandBatchPoolNodeReplacementConfig(list []interface{}) *pool.NodePlacementConfiguration {
	if len(list) == 0 || list[0] == nil {
		return nil
	}
	item := list[0].(map[string]interface{})["policy"].(string)
	return &pool.NodePlacementConfiguration{
		Policy: pointer.To(pool.NodePlacementPolicyType(item)),
	}
}

func expandBatchPoolWindowsConfiguration(list []interface{}) *pool.WindowsConfiguration {
	if len(list) == 0 {
		return nil
	}

	item := list[0].(map[string]interface{})["enable_automatic_updates"].(bool)
	return &pool.WindowsConfiguration{
		EnableAutomaticUpdates: utils.Bool(item),
	}
}

func expandBatchPoolExtensions(list []interface{}) (*[]pool.VmExtension, error) {
	if len(list) == 0 || list[0] == nil {
		return nil, nil
	}

	var result []pool.VmExtension

	for _, tempItem := range list {
		item := tempItem.(map[string]interface{})
		if batchPoolExtension, err := expandBatchPoolExtension(item); err == nil {
			result = append(result, *batchPoolExtension)
		} else {
			return nil, err
		}
	}

	return &result, nil
}

func expandBatchPoolExtension(ref map[string]interface{}) (*pool.VmExtension, error) {
	if len(ref) == 0 {
		return nil, nil
	}

	result := pool.VmExtension{
		Name:      ref["name"].(string),
		Publisher: ref["publisher"].(string),
		Type:      ref["type"].(string),
	}

	if autoUpgradeMinorVersion, ok := ref["auto_upgrade_minor_version"]; ok {
		result.AutoUpgradeMinorVersion = utils.Bool(autoUpgradeMinorVersion.(bool))
	}

	if autoUpgradeEnabled, ok := ref["automatic_upgrade_enabled"]; ok {
		result.EnableAutomaticUpgrade = utils.Bool(autoUpgradeEnabled.(bool))
	}

	if typeHandlerVersion, ok := ref["type_handler_version"]; ok {
		result.TypeHandlerVersion = utils.String(typeHandlerVersion.(string))
	}

	if settings, ok := ref["settings_json"]; ok {
		err := json.Unmarshal([]byte(settings.(string)), &result.Settings)
		if err != nil {
			return nil, fmt.Errorf("unmarshaling `settings_json`: %+v", err)
		}
	}

	if protectedSettings, ok := ref["protected_settings"]; ok {
		err := json.Unmarshal([]byte(protectedSettings.(string)), &result.ProtectedSettings)
		if err != nil {
			return nil, fmt.Errorf("unmarshaling `protected_settings`: %+v", err)
		}
	}

	if tmpItem, ok := ref["provision_after_extensions"]; ok {
		result.ProvisionAfterExtensions = utils.ExpandStringSlice(tmpItem.(*pluginsdk.Set).List())
	}

	return &result, nil
}

func expandBatchPoolDiskEncryptionConfiguration(list []interface{}) (*pool.DiskEncryptionConfiguration, error) {
	if len(list) == 0 || list[0] == nil {
		return nil, nil
	}
	var result pool.DiskEncryptionConfiguration

	var targetList []pool.DiskEncryptionTarget

	for _, tempItem := range list {
		item := tempItem.(map[string]interface{})
		if dataDiskEncryptionTarget, ok := item["disk_encryption_target"]; ok {
			targetList = append(targetList, pool.DiskEncryptionTarget(dataDiskEncryptionTarget.(string)))
		} else {
			return nil, fmt.Errorf("disk_encryption_target either is empty or contains parsing errors")
		}
	}

	result.Targets = &targetList
	return &result, nil
}

func expandBatchPoolDataDisks(list []interface{}) *[]pool.DataDisk {
	if len(list) == 0 || list[0] == nil {
		return nil
	}
	var result []pool.DataDisk

	for _, tempItem := range list {
		item := tempItem.(map[string]interface{})
		result = append(result, expandBatchPoolDataDisk(item))
	}

	return &result
}

func expandBatchPoolDataDisk(ref map[string]interface{}) pool.DataDisk {
	return pool.DataDisk{
		Lun:                int64(ref["lun"].(int)),
		Caching:            pointer.To(pool.CachingType(ref["caching"].(string))),
		DiskSizeGB:         int64(ref["disk_size_gb"].(int)),
		StorageAccountType: pointer.To(pool.StorageAccountType(ref["storage_account_type"].(string))),
	}
}

func expandCommonEnvironmentProperties(env map[string]interface{}) *[]pool.EnvironmentSetting {
	envSettings := make([]pool.EnvironmentSetting, 0)

	for k, v := range env {
		theValue := v.(string)
		theKey := k
		envSetting := pool.EnvironmentSetting{
			Name:  theKey,
			Value: &theValue,
		}

		envSettings = append(envSettings, envSetting)
	}
	return &envSettings
}

// ExpandBatchMetaData expands Batch pool metadata
func ExpandBatchMetaData(input map[string]interface{}) *[]pool.MetadataItem {
	output := []pool.MetadataItem{}

	for k, v := range input {
		name := k
		value := v.(string)
		output = append(output, pool.MetadataItem{
			Name:  name,
			Value: value,
		})
	}

	return &output
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

func ExpandBatchPoolMountConfigurations(d *pluginsdk.ResourceData) (*[]pool.MountConfiguration, error) {
	var result []pool.MountConfiguration

	if mountConfigs, ok := d.GetOk("mount"); ok {
		mountConfigList := mountConfigs.([]interface{})
		for _, tempItem := range mountConfigList {
			item := tempItem.(map[string]interface{})
			result = append(result, expandBatchPoolMountConfiguration(item))
		}
		return &result, nil
	}

	return nil, fmt.Errorf("mount either is empty or contains parsing errors")
}

func expandBatchPoolMountConfiguration(ref map[string]interface{}) pool.MountConfiguration {
	var result pool.MountConfiguration
	if azureBlobFileSystemConfiguration, err := expandBatchPoolAzureBlobFileSystemConfiguration(ref["azure_blob_file_system"].([]interface{})); err == nil {
		result.AzureBlobFileSystemConfiguration = azureBlobFileSystemConfiguration
	}

	if azureFileShareConfiguration, err := expandBatchPoolAzureFileShareConfiguration(ref["azure_file_share"].([]interface{})); err == nil {
		result.AzureFileShareConfiguration = azureFileShareConfiguration
	}

	if cifsMountConfiguration, err := expandBatchPoolCIFSMountConfiguration(ref["cifs_mount"].([]interface{})); err == nil {
		result.CifsMountConfiguration = cifsMountConfiguration
	}

	if nfsMountConfiguration, err := expandBatchPoolNFSMountConfiguration(ref["nfs_mount"].([]interface{})); err == nil {
		result.NfsMountConfiguration = nfsMountConfiguration
	}

	return result
}

func expandBatchPoolAzureBlobFileSystemConfiguration(list []interface{}) (*pool.AzureBlobFileSystemConfiguration, interface{}) {
	if len(list) == 0 || list[0] == nil {
		return nil, fmt.Errorf("azure_blob_file_system is empty")
	}

	configMap := list[0].(map[string]interface{})
	result := pool.AzureBlobFileSystemConfiguration{
		AccountName:       configMap["account_name"].(string),
		ContainerName:     configMap["container_name"].(string),
		RelativeMountPath: configMap["relative_mount_path"].(string),
	}

	if accountKey, ok := configMap["account_key"]; ok && accountKey != "" {
		result.AccountKey = utils.String(accountKey.(string))
	} else if sasKey, ok := configMap["sas_key"]; ok && sasKey != "" {
		result.SasKey = utils.String(sasKey.(string))
	} else if computedIDRef, err := expandBatchPoolIdentityReference(configMap); err == nil {
		result.IdentityReference = computedIDRef
	}

	if blobfuseOptions, ok := configMap["blobfuse_options"]; ok {
		result.BlobfuseOptions = utils.String(blobfuseOptions.(string))
	}
	return &result, nil
}

func expandBatchPoolAzureFileShareConfiguration(list []interface{}) (*pool.AzureFileShareConfiguration, interface{}) {
	if len(list) == 0 || list[0] == nil {
		return nil, fmt.Errorf("azure_file_share is empty")
	}

	configMap := list[0].(map[string]interface{})
	result := pool.AzureFileShareConfiguration{
		AccountName:       configMap["account_name"].(string),
		AccountKey:        configMap["account_key"].(string),
		AzureFileURL:      configMap["azure_file_url"].(string),
		RelativeMountPath: configMap["relative_mount_path"].(string),
	}

	if mountOptions, ok := configMap["mount_options"]; ok {
		result.MountOptions = utils.String(mountOptions.(string))
	}

	return &result, nil
}

func expandBatchPoolCIFSMountConfiguration(list []interface{}) (*pool.CIFSMountConfiguration, interface{}) {
	if len(list) == 0 || list[0] == nil {
		return nil, fmt.Errorf("cifs_mount is empty")
	}

	configMap := list[0].(map[string]interface{})
	result := pool.CIFSMountConfiguration{
		UserName:          configMap["user_name"].(string),
		Source:            configMap["source"].(string),
		Password:          configMap["password"].(string),
		RelativeMountPath: configMap["relative_mount_path"].(string),
	}

	if mountOptions, ok := configMap["mount_options"]; ok {
		result.MountOptions = utils.String(mountOptions.(string))
	}

	return &result, nil
}

func expandBatchPoolNFSMountConfiguration(list []interface{}) (*pool.NFSMountConfiguration, interface{}) {
	if len(list) == 0 || list[0] == nil {
		return nil, fmt.Errorf("nfs_mount is empty")
	}

	configMap := list[0].(map[string]interface{})
	result := pool.NFSMountConfiguration{
		Source:            configMap["source"].(string),
		RelativeMountPath: configMap["relative_mount_path"].(string),
	}

	if mountOptions, ok := configMap["mount_options"]; ok {
		result.MountOptions = utils.String(mountOptions.(string))
	}
	return &result, nil
}

func expandBatchPoolIdentityReference(ref map[string]interface{}) (*pool.ComputeNodeIdentityReference, error) {
	var result pool.ComputeNodeIdentityReference
	if iid, ok := ref["identity_id"]; ok && iid != "" {
		result.ResourceId = utils.String(iid.(string))
		return &result, nil
	}
	return nil, fmt.Errorf("identity_id is empty")
}

// ExpandBatchPoolNetworkConfiguration expands Batch pool network configuration
func ExpandBatchPoolNetworkConfiguration(list []interface{}) (*pool.NetworkConfiguration, error) {
	if len(list) == 0 || list[0] == nil {
		return nil, nil
	}

	networkConfigValue := list[0].(map[string]interface{})
	networkConfiguration := &pool.NetworkConfiguration{}

	if v, ok := networkConfigValue["dynamic_vnet_assignment_scope"]; ok {
		networkConfiguration.DynamicVnetAssignmentScope = pointer.To(pool.DynamicVNetAssignmentScope(v.(string)))
	}

	if v, ok := networkConfigValue["accelerated_networking_enabled"]; ok {
		networkConfiguration.EnableAcceleratedNetworking = pointer.FromBool(v.(bool))
	}

	if v, ok := networkConfigValue["subnet_id"]; ok {
		if value := v.(string); value != "" {
			networkConfiguration.SubnetId = &value
		}
	}

	if v, ok := networkConfigValue["public_ips"]; ok {
		if networkConfiguration.PublicIPAddressConfiguration == nil {
			networkConfiguration.PublicIPAddressConfiguration = &pool.PublicIPAddressConfiguration{}
		}

		publicIPsRaw := v.(*pluginsdk.Set).List()
		networkConfiguration.PublicIPAddressConfiguration.IPAddressIds = utils.ExpandStringSlice(publicIPsRaw)
	}

	if v, ok := networkConfigValue["endpoint_configuration"]; ok {
		endpoint, err := expandPoolEndpointConfiguration(v.([]interface{}))
		if err != nil {
			return nil, err
		}
		networkConfiguration.EndpointConfiguration = endpoint
	}

	if v, ok := networkConfigValue["public_address_provisioning_type"]; ok {
		if networkConfiguration.PublicIPAddressConfiguration == nil {
			networkConfiguration.PublicIPAddressConfiguration = &pool.PublicIPAddressConfiguration{}
		}

		if value := v.(string); value != "" {
			networkConfiguration.PublicIPAddressConfiguration.Provision = pointer.To(pool.IPAddressProvisioningType(value))
		}
	}

	return networkConfiguration, nil
}

func ExpandBatchPoolTaskSchedulingPolicy(d *pluginsdk.ResourceData) (*pool.TaskSchedulingPolicy, error) {
	var result pool.TaskSchedulingPolicy

	if taskSchedulingPolicyString, ok := d.GetOk("task_scheduling_policy"); ok {
		taskSchedulingPolicy := taskSchedulingPolicyString.([]interface{})
		if len(taskSchedulingPolicy) > 0 {
			item := taskSchedulingPolicy[0].(map[string]interface{})
			result.NodeFillType = pool.ComputeNodeFillType(item["node_fill_type"].(string))
		}
		return &result, nil
	}
	return nil, fmt.Errorf("task_scheduling_policy either is empty or contains parsing errors")
}

func expandPoolEndpointConfiguration(list []interface{}) (*pool.PoolEndpointConfiguration, error) {
	if len(list) == 0 || list[0] == nil {
		return nil, nil
	}

	inboundNatPools := make([]pool.InboundNatPool, len(list))

	for i, inboundNatPoolsValue := range list {
		inboundNatPool := inboundNatPoolsValue.(map[string]interface{})

		name := inboundNatPool["name"].(string)
		protocol := pool.InboundEndpointProtocol(inboundNatPool["protocol"].(string))
		backendPort := int32(inboundNatPool["backend_port"].(int))
		frontendPortRange := inboundNatPool["frontend_port_range"].(string)
		parts := strings.Split(frontendPortRange, "-")
		frontendPortRangeStart, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, err
		}
		frontendPortRangeEnd, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, err
		}

		networkSecurityGroupRules := expandPoolNetworkSecurityGroupRule(inboundNatPool["network_security_group_rules"].([]interface{}))

		inboundNatPools[i] = pool.InboundNatPool{
			Name:                      name,
			Protocol:                  protocol,
			BackendPort:               int64(backendPort),
			FrontendPortRangeStart:    int64(frontendPortRangeStart),
			FrontendPortRangeEnd:      int64(frontendPortRangeEnd),
			NetworkSecurityGroupRules: &networkSecurityGroupRules,
		}
	}

	return &pool.PoolEndpointConfiguration{
		InboundNatPools: inboundNatPools,
	}, nil
}

func expandPoolNetworkSecurityGroupRule(list []interface{}) []pool.NetworkSecurityGroupRule {
	if len(list) == 0 || list[0] == nil {
		return []pool.NetworkSecurityGroupRule{}
	}

	networkSecurityGroupRule := make([]pool.NetworkSecurityGroupRule, 0)
	for _, groupRule := range list {
		groupRuleMap := groupRule.(map[string]interface{})

		priority := int32(groupRuleMap["priority"].(int))
		sourceAddressPrefix := groupRuleMap["source_address_prefix"].(string)
		access := pool.NetworkSecurityGroupRuleAccess(groupRuleMap["access"].(string))

		networkSecurityGroupRuleObject := pool.NetworkSecurityGroupRule{
			Priority:            int64(priority),
			SourceAddressPrefix: sourceAddressPrefix,
			Access:              access,
		}

		portRanges := groupRuleMap["source_port_ranges"].([]interface{})
		if len(portRanges) > 0 {
			portRangesResult := make([]string, 0)
			for _, v := range portRanges {
				portRangesResult = append(portRangesResult, v.(string))
			}
			networkSecurityGroupRuleObject.SourcePortRanges = &portRangesResult
		}

		networkSecurityGroupRule = append(networkSecurityGroupRule, pool.NetworkSecurityGroupRule{
			Priority:            int64(priority),
			SourceAddressPrefix: sourceAddressPrefix,
			Access:              access,
		})
	}

	return networkSecurityGroupRule
}

func flattenBatchPoolNetworkConfiguration(input *pool.NetworkConfiguration) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	subnetId := ""
	if input.SubnetId != nil {
		subnetId = *input.SubnetId
	}

	publicIPAddressIds := make([]interface{}, 0)
	publicAddressProvisioningType := ""
	if config := input.PublicIPAddressConfiguration; config != nil {
		publicIPAddressIds = utils.FlattenStringSlice(config.IPAddressIds)
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
			"subnet_id":                        subnetId,
		},
	}
}

func ExpandBatchPoolUserAccounts(d *pluginsdk.ResourceData) (*[]pool.UserAccount, error) {
	var result []pool.UserAccount

	if userAccountList, ok := d.GetOk("user_accounts"); ok {
		userAccounts := userAccountList.([]interface{})
		if len(userAccounts) > 0 && userAccounts[0] != nil {
			for _, tempItem := range userAccounts {
				item := tempItem.(map[string]interface{})
				result = append(result, expandBatchPoolUserAccount(item))
			}
			return &result, nil
		}
	}

	return nil, fmt.Errorf("user_accounts either is empty or contains parsing errors")
}

func expandBatchPoolUserAccount(ref map[string]interface{}) pool.UserAccount {
	result := pool.UserAccount{
		Name:           ref["name"].(string),
		Password:       ref["password"].(string),
		ElevationLevel: pointer.To(pool.ElevationLevel(ref["elevation_level"].(string))),
	}

	if linuxUserConfig, ok := ref["linux_user_configuration"]; ok {
		if linuxUserConfig != nil && len(linuxUserConfig.([]interface{})) > 0 {
			linuxUserConfigMap := linuxUserConfig.([]interface{})[0].(map[string]interface{})
			var linuxUserConfig pool.LinuxUserConfiguration
			if uid, ok := linuxUserConfigMap["uid"]; ok {
				linuxUserConfig = pool.LinuxUserConfiguration{
					Uid: utils.Int64(int64(uid.(int))),
					Gid: utils.Int64(int64(linuxUserConfigMap["gid"].(int))),
				}
			}
			if sshPrivateKey, ok := linuxUserConfigMap["ssh_private_key"]; ok {
				linuxUserConfig.SshPrivateKey = utils.String(sshPrivateKey.(string))
			}
			result.LinuxUserConfiguration = &linuxUserConfig
		}
	}

	if winUserConfig, ok := ref["windows_user_configuration"]; ok {
		if winUserConfig != nil && len(winUserConfig.([]interface{})) > 0 {
			winUserConfigMap := winUserConfig.([]interface{})[0].(map[string]interface{})
			result.WindowsUserConfiguration = &pool.WindowsUserConfiguration{
				LoginMode: pointer.To(pool.LoginMode(winUserConfigMap["login_mode"].(string))),
			}
		}
	}

	return result
}
