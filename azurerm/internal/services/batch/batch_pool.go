package batch

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/batch/mgmt/2019-08-01/batch"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

// flattenBatchPoolAutoScaleSettings flattens the auto scale settings for a Batch pool
func flattenBatchPoolAutoScaleSettings(settings *batch.AutoScaleSettings) []interface{} {
	results := make([]interface{}, 0)

	if settings == nil {
		log.Printf("[DEBUG] settings is nil")
		return results
	}

	result := make(map[string]interface{})

	if settings.EvaluationInterval != nil {
		result["evaluation_interval"] = *settings.EvaluationInterval
	}

	if settings.Formula != nil {
		result["formula"] = *settings.Formula
	}

	return append(results, result)
}

// flattenBatchPoolFixedScaleSettings flattens the fixed scale settings for a Batch pool
func flattenBatchPoolFixedScaleSettings(settings *batch.FixedScaleSettings) []interface{} {
	results := make([]interface{}, 0)

	if settings == nil {
		log.Printf("[DEBUG] settings is nil")
		return results
	}

	result := make(map[string]interface{})

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
func flattenBatchPoolImageReference(image *batch.ImageReference) []interface{} {
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
	if image.ID != nil {
		result["id"] = *image.ID
	}

	return append(results, result)
}

// flattenBatchPoolStartTask flattens a Batch pool start task
func flattenBatchPoolStartTask(startTask *batch.StartTask) []interface{} {
	results := make([]interface{}, 0)

	if startTask == nil {
		log.Printf("[DEBUG] startTask is nil")
		return results
	}

	result := make(map[string]interface{})
	if startTask.CommandLine != nil {
		result["command_line"] = *startTask.CommandLine
	}
	if startTask.WaitForSuccess != nil {
		result["wait_for_success"] = *startTask.WaitForSuccess
	}
	if startTask.MaxTaskRetryCount != nil {
		result["max_task_retry_count"] = *startTask.MaxTaskRetryCount
	}

	if startTask.UserIdentity != nil {
		userIdentity := make(map[string]interface{})
		if startTask.UserIdentity.AutoUser != nil {
			autoUser := make(map[string]interface{})

			elevationLevel := string(startTask.UserIdentity.AutoUser.ElevationLevel)
			scope := string(startTask.UserIdentity.AutoUser.Scope)

			autoUser["elevation_level"] = elevationLevel
			autoUser["scope"] = scope

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
			resourceFiles = append(resourceFiles, resourceFile)
		}
	}

	if startTask.EnvironmentSettings != nil {
		environment := make(map[string]interface{})
		for _, envSetting := range *startTask.EnvironmentSettings {
			environment[*envSetting.Name] = *envSetting.Value
		}

		result["environment"] = environment
	}
	result["resource_file"] = resourceFiles

	return append(results, result)
}

// flattenBatchPoolCertificateReferences flattens a Batch pool certificate reference
func flattenBatchPoolCertificateReferences(armCertificates *[]batch.CertificateReference) []interface{} {
	if armCertificates == nil {
		return []interface{}{}
	}
	output := make([]interface{}, 0)

	for _, armCertificate := range *armCertificates {
		certificate := map[string]interface{}{}
		if armCertificate.ID != nil {
			certificate["id"] = *armCertificate.ID
		}
		certificate["store_location"] = string(armCertificate.StoreLocation)
		if armCertificate.StoreName != nil {
			certificate["store_name"] = *armCertificate.StoreName
		}
		visibility := &schema.Set{F: schema.HashString}
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
func flattenBatchPoolContainerConfiguration(d *schema.ResourceData, armContainerConfiguration *batch.ContainerConfiguration) interface{} {
	result := make(map[string]interface{})

	if armContainerConfiguration == nil {
		return nil
	}

	if armContainerConfiguration.Type != nil {
		result["type"] = *armContainerConfiguration.Type
	}

	names := &schema.Set{F: schema.HashString}
	if armContainerConfiguration.ContainerImageNames != nil {
		for _, armName := range *armContainerConfiguration.ContainerImageNames {
			names.Add(armName)
		}
	}
	result["container_image_names"] = names

	result["container_registries"] = flattenBatchPoolContainerRegistries(d, armContainerConfiguration.ContainerRegistries)

	return []interface{}{result}
}

func flattenBatchPoolContainerRegistries(d *schema.ResourceData, armContainerRegistries *[]batch.ContainerRegistry) []interface{} {
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

func flattenBatchPoolContainerRegistry(d *schema.ResourceData, armContainerRegistry *batch.ContainerRegistry) map[string]interface{} {
	result := make(map[string]interface{})

	if armContainerRegistry == nil {
		return result
	}
	if registryServer := armContainerRegistry.RegistryServer; registryServer != nil {
		result["registry_server"] = *registryServer
	}
	if userName := armContainerRegistry.UserName; userName != nil {
		result["user_name"] = *userName
	}

	// If we didn't specify a registry server and user name, just return what we have now rather than trying to locate the password
	if len(result) != 2 {
		return result
	}

	result["password"] = findBatchPoolContainerRegistryPassword(d, result["registry_server"].(string), result["user_name"].(string))

	return result
}

func findBatchPoolContainerRegistryPassword(d *schema.ResourceData, armServer string, armUsername string) interface{} {
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

// ExpandBatchPoolImageReference expands Batch pool image reference
func ExpandBatchPoolImageReference(list []interface{}) (*batch.ImageReference, error) {
	if len(list) == 0 {
		return nil, fmt.Errorf("Error: storage image reference should be defined")
	}

	storageImageRef := list[0].(map[string]interface{})
	imageRef := &batch.ImageReference{}

	if storageImageRef["id"] != nil && storageImageRef["id"] != "" {
		storageImageRefID := storageImageRef["id"].(string)
		imageRef.ID = &storageImageRefID
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
func ExpandBatchPoolContainerConfiguration(list []interface{}) (*batch.ContainerConfiguration, error) {
	if len(list) == 0 || list[0] == nil {
		return nil, nil
	}

	block := list[0].(map[string]interface{})

	containerRegistries, err := expandBatchPoolContainerRegistries(block["container_registries"].([]interface{}))
	if err != nil {
		return nil, err
	}

	obj := &batch.ContainerConfiguration{
		Type:                utils.String(block["type"].(string)),
		ContainerRegistries: containerRegistries,
		ContainerImageNames: utils.ExpandStringSlice(block["container_image_names"].(*schema.Set).List()),
	}

	return obj, nil
}

func expandBatchPoolContainerRegistries(list []interface{}) (*[]batch.ContainerRegistry, error) {
	result := []batch.ContainerRegistry{}

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

func expandBatchPoolContainerRegistry(ref map[string]interface{}) (*batch.ContainerRegistry, error) {
	if len(ref) == 0 {
		return nil, fmt.Errorf("Error: container registry reference should be defined")
	}

	containerRegistry := batch.ContainerRegistry{
		RegistryServer: utils.String(ref["registry_server"].(string)),
		UserName:       utils.String(ref["user_name"].(string)),
		Password:       utils.String(ref["password"].(string)),
	}
	return &containerRegistry, nil
}

// ExpandBatchPoolCertificateReferences expands Batch pool certificate references
func ExpandBatchPoolCertificateReferences(list []interface{}) (*[]batch.CertificateReference, error) {
	var result []batch.CertificateReference

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

func expandBatchPoolCertificateReference(ref map[string]interface{}) (*batch.CertificateReference, error) {
	if len(ref) == 0 {
		return nil, fmt.Errorf("Error: storage image reference should be defined")
	}

	id := ref["id"].(string)
	storeLocation := ref["store_location"].(string)
	storeName := ref["store_name"].(string)
	visibilityRefs := ref["visibility"].(*schema.Set)
	var visibility []batch.CertificateVisibility
	if visibilityRefs != nil {
		for _, visibilityRef := range visibilityRefs.List() {
			visibility = append(visibility, batch.CertificateVisibility(visibilityRef.(string)))
		}
	}

	certificateReference := &batch.CertificateReference{
		ID:            &id,
		StoreLocation: batch.CertificateStoreLocation(storeLocation),
		StoreName:     &storeName,
		Visibility:    &visibility,
	}
	return certificateReference, nil
}

// ExpandBatchPoolStartTask expands Batch pool start task
func ExpandBatchPoolStartTask(list []interface{}) (*batch.StartTask, error) {
	if len(list) == 0 {
		return nil, fmt.Errorf("Error: batch pool start task should be defined")
	}

	startTaskValue := list[0].(map[string]interface{})

	startTaskCmdLine := startTaskValue["command_line"].(string)
	maxTaskRetryCount := int32(startTaskValue["max_task_retry_count"].(int))
	waitForSuccess := startTaskValue["wait_for_success"].(bool)

	userIdentityList := startTaskValue["user_identity"].([]interface{})
	if len(userIdentityList) == 0 {
		return nil, fmt.Errorf("Error: batch pool start task user identity should be defined")
	}

	userIdentityValue := userIdentityList[0].(map[string]interface{})
	userIdentity := batch.UserIdentity{}

	if autoUserValue, ok := userIdentityValue["auto_user"]; ok {
		autoUser := autoUserValue.([]interface{})
		if len(autoUser) != 0 {
			autoUserMap := autoUser[0].(map[string]interface{})
			userIdentity.AutoUser = &batch.AutoUserSpecification{
				ElevationLevel: batch.ElevationLevel(autoUserMap["elevation_level"].(string)),
				Scope:          batch.AutoUserScope(autoUserMap["scope"].(string)),
			}
		}
	} else if userNameValue, ok := userIdentityValue["username"]; ok {
		userName := userNameValue.(string)
		userIdentity.UserName = &userName
	} else {
		return nil, fmt.Errorf("Error: either auto_user or user_name should be speicfied for Batch pool start task")
	}

	resourceFileList := startTaskValue["resource_file"].([]interface{})
	resourceFiles := make([]batch.ResourceFile, 0)
	for _, resourceFileValueTemp := range resourceFileList {
		resourceFileValue := resourceFileValueTemp.(map[string]interface{})
		resourceFile := batch.ResourceFile{}
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
		resourceFiles = append(resourceFiles, resourceFile)
	}

	startTask := &batch.StartTask{
		CommandLine:       &startTaskCmdLine,
		MaxTaskRetryCount: &maxTaskRetryCount,
		WaitForSuccess:    &waitForSuccess,
		UserIdentity:      &userIdentity,
		ResourceFiles:     &resourceFiles,
	}

	// populate environment settings, if defined
	if environment := startTaskValue["environment"]; environment != nil {
		envMap := environment.(map[string]interface{})
		envSettings := make([]batch.EnvironmentSetting, 0)

		for k, v := range envMap {
			theValue := v.(string)
			theKey := k
			envSetting := batch.EnvironmentSetting{
				Name:  &theKey,
				Value: &theValue,
			}

			envSettings = append(envSettings, envSetting)
		}

		startTask.EnvironmentSettings = &envSettings
	}

	return startTask, nil
}

// ValidateAzureRMBatchPoolName validates the name of a Batch pool
func ValidateAzureRMBatchPoolName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"any combination of alphanumeric characters including hyphens and underscores are allowed in %q: %q", k, value))
	}

	if 1 > len(value) {
		errors = append(errors, fmt.Errorf("%q cannot be less than 1 character: %q", k, value))
	}

	if len(value) > 64 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 64 characters: %q %d", k, value, len(value)))
	}

	return warnings, errors
}

// ExpandBatchMetaData expands Batch pool metadata
func ExpandBatchMetaData(input map[string]interface{}) *[]batch.MetadataItem {
	output := []batch.MetadataItem{}

	for k, v := range input {
		name := k
		value := v.(string)
		output = append(output, batch.MetadataItem{
			Name:  &name,
			Value: &value,
		})
	}

	return &output
}

// FlattenBatchMetaData flattens a Batch pool metadata
func FlattenBatchMetaData(metadatas *[]batch.MetadataItem) map[string]interface{} {
	output := make(map[string]interface{})

	if metadatas == nil {
		return output
	}

	for _, metadata := range *metadatas {
		if metadata.Name == nil || metadata.Value == nil {
			continue
		}

		output[*metadata.Name] = *metadata.Value
	}

	return output
}

// ExpandBatchPoolNetworkConfiguration expands Batch pool network configuration
func ExpandBatchPoolNetworkConfiguration(list []interface{}) (*batch.NetworkConfiguration, error) {
	if len(list) == 0 {
		return nil, nil
	}

	networkConfigValue := list[0].(map[string]interface{})
	networkConfiguration := &batch.NetworkConfiguration{}

	if v, ok := networkConfigValue["subnet_id"]; ok {
		if value := v.(string); value != "" {
			networkConfiguration.SubnetID = &value
		}
	}

	if v, ok := networkConfigValue["public_ips"]; ok {
		publicIPsRaw := v.(*schema.Set).List()
		networkConfiguration.PublicIPs = utils.ExpandStringSlice(publicIPsRaw)
	}

	if v, ok := networkConfigValue["endpoint_configuration"]; ok {
		endpoint, err := ExpandBatchPoolEndpointConfiguration(v.([]interface{}))
		if err != nil {
			return nil, err
		}
		networkConfiguration.EndpointConfiguration = endpoint
	}

	return networkConfiguration, nil
}

// ExpandBatchPoolEndpointConfiguration expands Batch pool endpoint configuration
func ExpandBatchPoolEndpointConfiguration(list []interface{}) (*batch.PoolEndpointConfiguration, error) {
	if len(list) == 0 {
		return nil, nil
	}

	inboundNatPools := make([]batch.InboundNatPool, len(list))

	for i, inboundNatPoolsValue := range list {
		inboundNatPool := inboundNatPoolsValue.(map[string]interface{})

		name := inboundNatPool["name"].(string)
		protocol := batch.InboundEndpointProtocol(inboundNatPool["protocol"].(string))
		backendPort := int32(inboundNatPool["backend_port"].(int))
		frontendPortRange := inboundNatPool["frontend_port_range"].(string)
		parts := strings.Split(frontendPortRange, "-")
		frontendPortRangeStart, _ := strconv.Atoi(parts[0])
		frontendPortRangeEnd, _ := strconv.Atoi(parts[1])

		networkSecurityGroupRules, err := ExpandBatchPoolNetworkSecurityGroupRule(inboundNatPool["network_security_group_rules"].([]interface{}))
		if err != nil {
			return nil, err
		}

		inboundNatPools[i] = batch.InboundNatPool{
			Name:                      &name,
			Protocol:                  protocol,
			BackendPort:               &backendPort,
			FrontendPortRangeStart:    utils.Int32(int32(frontendPortRangeStart)),
			FrontendPortRangeEnd:      utils.Int32(int32(frontendPortRangeEnd)),
			NetworkSecurityGroupRules: &networkSecurityGroupRules,
		}
	}

	return &batch.PoolEndpointConfiguration{
		InboundNatPools: &inboundNatPools,
	}, nil
}

// ExpandBatchPoolNetworkSecurityGroupRule expands Batch pool network security group rule
func ExpandBatchPoolNetworkSecurityGroupRule(list []interface{}) ([]batch.NetworkSecurityGroupRule, error) {
	if len(list) == 0 {
		return nil, nil
	}

	networkSecurityGroupRule := make([]batch.NetworkSecurityGroupRule, len(list))

	for i, groupRule := range list {
		groupRuleMap := groupRule.(map[string]interface{})

		priority := int32(groupRuleMap["priority"].(int))
		sourceAddressPrefix := groupRuleMap["source_address_prefix"].(string)
		access := batch.NetworkSecurityGroupRuleAccess(groupRuleMap["access"].(string))

		networkSecurityGroupRule[i] = batch.NetworkSecurityGroupRule{
			Priority:            &priority,
			SourceAddressPrefix: &sourceAddressPrefix,
			Access:              access,
		}
	}

	return networkSecurityGroupRule, nil
}

// FlattenBatchPoolNetworkConfiguration flattens the network configuration for a Batch pool
func FlattenBatchPoolNetworkConfiguration(networkConfig *batch.NetworkConfiguration) []interface{} {
	results := make([]interface{}, 0)

	if networkConfig == nil {
		log.Printf("[DEBUG] networkConfgiuration is nil")
		return nil
	}

	result := make(map[string]interface{})

	if networkConfig.SubnetID != nil {
		result["subnet_id"] = *networkConfig.SubnetID
	}

	if networkConfig.PublicIPs != nil {
		result["public_ips"] = schema.NewSet(schema.HashString, utils.FlattenStringSlice(networkConfig.PublicIPs))
	}

	if cfg := networkConfig.EndpointConfiguration; cfg != nil && cfg.InboundNatPools != nil && len(*cfg.InboundNatPools) != 0 {
		endpointConfigs := make([]interface{}, len(*cfg.InboundNatPools))

		for i, inboundNatPool := range *cfg.InboundNatPools {
			inboundNatPoolMap := make(map[string]interface{})
			if inboundNatPool.Name != nil {
				inboundNatPoolMap["name"] = *inboundNatPool.Name
			}
			if inboundNatPool.BackendPort != nil {
				inboundNatPoolMap["backend_port"] = *inboundNatPool.BackendPort
			}
			if inboundNatPool.FrontendPortRangeStart != nil && inboundNatPool.FrontendPortRangeEnd != nil {
				inboundNatPoolMap["frontend_port_range"] = fmt.Sprintf("%d-%d", *inboundNatPool.FrontendPortRangeStart, *inboundNatPool.FrontendPortRangeEnd)
			}
			inboundNatPoolMap["protocol"] = inboundNatPool.Protocol

			if sgRules := inboundNatPool.NetworkSecurityGroupRules; sgRules != nil && len(*sgRules) != 0 {
				networkSecurities := make([]interface{}, len(*sgRules))
				for j, networkSecurity := range *sgRules {
					networkSecurityMap := make(map[string]interface{})

					if networkSecurity.Priority != nil {
						networkSecurityMap["priority"] = *networkSecurity.Priority
					}
					if networkSecurity.SourceAddressPrefix != nil {
						networkSecurityMap["source_address_prefix"] = *networkSecurity.SourceAddressPrefix
					}
					networkSecurityMap["access"] = networkSecurity.Access
					networkSecurities[j] = networkSecurityMap
				}
				inboundNatPoolMap["network_security_group_rules"] = networkSecurities
			}
			endpointConfigs[i] = inboundNatPoolMap
		}

		result["endpoint_configuration"] = endpointConfigs
	}

	return append(results, result)
}
