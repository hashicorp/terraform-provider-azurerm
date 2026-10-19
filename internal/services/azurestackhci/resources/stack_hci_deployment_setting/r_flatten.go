// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package stack_hci_deployment_setting

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/deploymentsettings"
)

func flattenDeploymentSettingScaleUnits(input []deploymentsettings.ScaleUnits) []ScaleUnitModel {
	if len(input) == 0 {
		return make([]ScaleUnitModel, 0)
	}

	results := make([]ScaleUnitModel, 0, len(input))
	for _, item := range input {
		result := ScaleUnitModel{
			ActiveDirectoryOrganizationalUnitPath: pointer.From(item.DeploymentData.AdouPath),
			Cluster:                               flattenDeploymentSettingCluster(item.DeploymentData.Cluster),
			DomainFqdn:                            pointer.From(item.DeploymentData.DomainFqdn),
			HostNetwork:                           flattenDeploymentSettingHostNetwork(item.DeploymentData.HostNetwork),
			InfrastructureNetwork:                 flattenDeploymentSettingInfrastructureNetwork(item.DeploymentData.InfrastructureNetwork),
			NamePrefix:                            pointer.From(item.DeploymentData.NamingPrefix),
			OptionalService:                       flattenDeploymentSettingOptionalService(item.DeploymentData.OptionalServices),
			PhysicalNode:                          flattenDeploymentSettingPhysicalNode(item.DeploymentData.PhysicalNodes),
			SecretsLocation:                       pointer.From(item.DeploymentData.SecretsLocation),
			Storage:                               flattenDeploymentSettingStorage(item.DeploymentData.Storage),
		}

		if observability := item.DeploymentData.Observability; observability != nil {
			result.EpisodicDataUploadEnabled = pointer.From(observability.EpisodicDataUpload)
			result.EuLocationEnabled = pointer.From(observability.EuLocation)
			result.StreamingDataClientEnabled = pointer.From(observability.StreamingDataClient)
		}

		if securitySettings := item.DeploymentData.SecuritySettings; securitySettings != nil {
			result.BitlockerBootVolumeEnabled = pointer.From(securitySettings.BitlockerBootVolume)
			result.BitlockerDataVolumeEnabled = pointer.From(securitySettings.BitlockerDataVolumes)
			result.CredentialGuardEnabled = pointer.From(securitySettings.CredentialGuardEnforced)
			result.DriftControlEnabled = pointer.From(securitySettings.DriftControlEnforced)
			result.DrtmProtectionEnabled = pointer.From(securitySettings.DrtmProtection)
			result.HvciProtectionEnabled = pointer.From(securitySettings.HvciProtection)
			result.SideChannelMitigationEnabled = pointer.From(securitySettings.SideChannelMitigationEnforced)
			result.SmbClusterEncryptionEnabled = pointer.From(securitySettings.SmbClusterEncryption)
			result.SmbSigningEnabled = pointer.From(securitySettings.SmbSigningEnforced)
			result.WdacEnabled = pointer.From(securitySettings.WdacEnforced)
		}

		results = append(results, result)
	}

	return results
}

func flattenDeploymentSettingCluster(input *deploymentsettings.DeploymentCluster) []ClusterModel {
	if input == nil {
		return make([]ClusterModel, 0)
	}

	return []ClusterModel{{
		AzureServiceEndpoint: pointer.From(input.AzureServiceEndpoint),
		CloudAccountName:     pointer.From(input.CloudAccountName),
		Name:                 pointer.From(input.Name),
		WitnessType:          pointer.From(input.WitnessType),
		WitnessPath:          pointer.From(input.WitnessPath),
	}}
}

func flattenDeploymentSettingHostNetwork(input *deploymentsettings.HostNetwork) []HostNetworkModel {
	if input == nil {
		return make([]HostNetworkModel, 0)
	}

	return []HostNetworkModel{{
		Intent:                               flattenDeploymentSettingHostNetworkIntent(input.Intents),
		StorageAutoIpEnabled:                 pointer.From(input.EnableStorageAutoIP),
		StorageConnectivitySwitchlessEnabled: pointer.From(input.StorageConnectivitySwitchless),
		StorageNetwork:                       flattenDeploymentSettingHostNetworkStorageNetwork(input.StorageNetworks),
	}}
}

func flattenDeploymentSettingHostNetworkIntent(input *[]deploymentsettings.Intents) []HostNetworkIntentModel {
	if input == nil {
		return make([]HostNetworkIntentModel, 0)
	}

	results := make([]HostNetworkIntentModel, 0, len(*input))
	for _, item := range *input {
		results = append(results, HostNetworkIntentModel{
			Adapter:                        pointer.From(item.Adapter),
			AdapterPropertyOverride:        flattenHostNetworkIntentAdapterPropertyOverride(item.AdapterPropertyOverrides),
			Name:                           pointer.From(item.Name),
			AdapterPropertyOverrideEnabled: pointer.From(item.OverrideAdapterProperty),
			QosPolicyOverrideEnabled:       pointer.From(item.OverrideQosPolicy),
			VirtualSwitchConfigurationOverrideEnabled: pointer.From(item.OverrideVirtualSwitchConfiguration),
			QosPolicyOverride:                         flattenHostNetworkIntentQosPolicyOverride(item.QosPolicyOverrides),
			TrafficType:                               pointer.From(item.TrafficType),
			VirtualSwitchConfigurationOverride:        flattenHostNetworkIntentVirtualSwitchConfigurationOverride(item.VirtualSwitchConfigurationOverrides),
		})
	}

	return results
}

func flattenHostNetworkIntentAdapterPropertyOverride(input *deploymentsettings.AdapterPropertyOverrides) []HostNetworkIntentAdapterPropertyOverrideModel {
	if input == nil {
		return make([]HostNetworkIntentAdapterPropertyOverrideModel, 0)
	}

	jumboPacket := pointer.From(input.JumboPacket)
	networkDirect := pointer.From(input.NetworkDirect)
	networkDirectTechnology := pointer.From(input.NetworkDirectTechnology)

	// server will return the block with empty string in all fields by default
	if jumboPacket == "" && networkDirect == "" && networkDirectTechnology == "" {
		return make([]HostNetworkIntentAdapterPropertyOverrideModel, 0)
	}

	return []HostNetworkIntentAdapterPropertyOverrideModel{{
		JumboPacket:             jumboPacket,
		NetworkDirect:           networkDirect,
		NetworkDirectTechnology: networkDirectTechnology,
	}}
}

func flattenHostNetworkIntentQosPolicyOverride(input *deploymentsettings.QosPolicyOverrides) []HostNetworkIntentQosPolicyOverrideModel {
	if input == nil {
		return make([]HostNetworkIntentQosPolicyOverrideModel, 0)
	}

	bandwidthPercentageSMB := pointer.From(input.BandwidthPercentageSMB)
	priorityValue8021ActionCluster := pointer.From(input.PriorityValue8021ActionCluster)
	priorityValue8021ActionSMB := pointer.From(input.PriorityValue8021ActionSMB)

	// server will return the block with empty string in all fields by default
	if bandwidthPercentageSMB == "" && priorityValue8021ActionCluster == "" && priorityValue8021ActionSMB == "" {
		return make([]HostNetworkIntentQosPolicyOverrideModel, 0)
	}

	return []HostNetworkIntentQosPolicyOverrideModel{{
		BandWidthPercentageSMB:         bandwidthPercentageSMB,
		PriorityValue8021ActionCluster: priorityValue8021ActionCluster,
		PriorityValue8021ActionSMB:     priorityValue8021ActionSMB,
	}}
}

func flattenHostNetworkIntentVirtualSwitchConfigurationOverride(input *deploymentsettings.VirtualSwitchConfigurationOverrides) []HostNetworkIntentVirtualSwitchConfigurationOverrideModel {
	if input == nil {
		return make([]HostNetworkIntentVirtualSwitchConfigurationOverrideModel, 0)
	}

	enableIov := pointer.From(input.EnableIov)
	loadBalancingAlgorithm := pointer.From(input.LoadBalancingAlgorithm)

	// server will return the block with empty string in all fields by default
	if enableIov == "" && loadBalancingAlgorithm == "" {
		return make([]HostNetworkIntentVirtualSwitchConfigurationOverrideModel, 0)
	}

	return []HostNetworkIntentVirtualSwitchConfigurationOverrideModel{{
		EnableIov:              enableIov,
		LoadBalancingAlgorithm: loadBalancingAlgorithm,
	}}
}

func flattenDeploymentSettingHostNetworkStorageNetwork(input *[]deploymentsettings.StorageNetworks) []HostNetworkStorageNetworkModel {
	if input == nil {
		return make([]HostNetworkStorageNetworkModel, 0)
	}

	results := make([]HostNetworkStorageNetworkModel, 0, len(*input))
	for _, item := range *input {
		results = append(results, HostNetworkStorageNetworkModel{
			Name:               pointer.From(item.Name),
			NetworkAdapterName: pointer.From(item.NetworkAdapterName),
			VlanId:             pointer.From(item.VlanId),
		})
	}

	return results
}

func flattenDeploymentSettingInfrastructureNetwork(input *[]deploymentsettings.InfrastructureNetwork) []InfrastructureNetworkModel {
	if input == nil {
		return make([]InfrastructureNetworkModel, 0)
	}

	results := make([]InfrastructureNetworkModel, 0, len(*input))
	for _, item := range *input {
		results = append(results, InfrastructureNetworkModel{
			DhcpEnabled: pointer.From(item.UseDhcp),
			DnsServer:   pointer.From(item.DnsServers),
			Gateway:     pointer.From(item.Gateway),
			IpPool:      flattenDeploymentSettingInfrastructureNetworkIpPool(item.IPPools),
			SubnetMask:  pointer.From(item.SubnetMask),
		})
	}

	return results
}

func flattenDeploymentSettingInfrastructureNetworkIpPool(input *[]deploymentsettings.IPPools) []IpPoolModel {
	if input == nil {
		return make([]IpPoolModel, 0)
	}

	results := make([]IpPoolModel, 0, len(*input))
	for _, item := range *input {
		results = append(results, IpPoolModel{
			EndingAddress:   pointer.From(item.EndingAddress),
			StartingAddress: pointer.From(item.StartingAddress),
		})
	}

	return results
}

func flattenDeploymentSettingOptionalService(input *deploymentsettings.OptionalServices) []OptionalServiceModel {
	if input == nil {
		return make([]OptionalServiceModel, 0)
	}

	return []OptionalServiceModel{{
		CustomLocation: pointer.From(input.CustomLocation),
	}}
}

func flattenDeploymentSettingPhysicalNode(input *[]deploymentsettings.PhysicalNodes) []PhysicalNodeModel {
	if input == nil {
		return make([]PhysicalNodeModel, 0)
	}

	results := make([]PhysicalNodeModel, 0, len(*input))
	for _, item := range *input {
		results = append(results, PhysicalNodeModel{
			Ipv4Address: pointer.From(item.IPv4Address),
			Name:        pointer.From(item.Name),
		})
	}

	return results
}

func flattenDeploymentSettingStorage(input *deploymentsettings.Storage) []StorageModel {
	if input == nil {
		return make([]StorageModel, 0)
	}

	return []StorageModel{{
		ConfigurationMode: pointer.From(input.ConfigurationMode),
	}}
}
