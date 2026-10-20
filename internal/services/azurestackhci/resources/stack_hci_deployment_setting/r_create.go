// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package stack_hci_deployment_setting

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/deploymentsettings"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 6 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AzureStackHCI.DeploymentSettings

			var config StackHCIDeploymentSettingModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			stackHCIClusterId, err := deploymentsettings.ParseClusterID(config.StackHCIClusterId)
			if err != nil {
				return err
			}
			id := deploymentsettings.NewDeploymentSettingID(stackHCIClusterId.SubscriptionId, stackHCIClusterId.ResourceGroupName, stackHCIClusterId.ClusterName, "default")

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			payload := deploymentsettings.DeploymentSetting{
				Properties: &deploymentsettings.DeploymentSettingsProperties{
					ArcNodeResourceIds: config.ArcResourceIds,
					DeploymentMode:     deploymentsettings.DeploymentModeValidate,
					DeploymentConfiguration: deploymentsettings.DeploymentConfiguration{
						Version:    pointer.To(config.Version),
						ScaleUnits: expandDeploymentSettingScaleUnits(config.ScaleUnit),
					},
				},
			}

			// do validation
			if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
				return fmt.Errorf("validating %s: %+v", id, err)
			}

			// do deployment
			payload.Properties.DeploymentMode = deploymentsettings.DeploymentModeDeploy
			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, payload, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("deploying %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func expandDeploymentSettingScaleUnits(input []ScaleUnitModel) []deploymentsettings.ScaleUnits {
	if len(input) == 0 {
		return nil
	}

	results := make([]deploymentsettings.ScaleUnits, 0, len(input))
	for _, item := range input {
		results = append(results, deploymentsettings.ScaleUnits{
			DeploymentData: deploymentsettings.DeploymentData{
				AdouPath:              pointer.To(item.ActiveDirectoryOrganizationalUnitPath),
				Cluster:               expandDeploymentSettingCluster(item.Cluster),
				DomainFqdn:            pointer.To(item.DomainFqdn),
				HostNetwork:           expandDeploymentSettingHostNetwork(item.HostNetwork),
				InfrastructureNetwork: expandDeploymentSettingInfrastructureNetwork(item.InfrastructureNetwork),
				NamingPrefix:          pointer.To(item.NamePrefix),
				Observability:         expandDeploymentSettingObservability(item),
				OptionalServices:      expandDeploymentSettingOptionalService(item.OptionalService),
				PhysicalNodes:         expandDeploymentSettingPhysicalNode(item.PhysicalNode),
				SecretsLocation:       pointer.To(item.SecretsLocation),
				SecuritySettings:      expandDeploymentSettingSecuritySetting(item),
				Storage:               expandDeploymentSettingStorage(item.Storage),
			},
		})
	}

	return results
}

func expandDeploymentSettingCluster(input []ClusterModel) *deploymentsettings.DeploymentCluster {
	if len(input) == 0 {
		return nil
	}

	v := input[0]

	return &deploymentsettings.DeploymentCluster{
		AzureServiceEndpoint: pointer.To(v.AzureServiceEndpoint),
		CloudAccountName:     pointer.To(v.CloudAccountName),
		Name:                 pointer.To(v.Name),
		WitnessType:          pointer.To(v.WitnessType),
		WitnessPath:          pointer.To(v.WitnessPath),
	}
}

func expandDeploymentSettingHostNetwork(input []HostNetworkModel) *deploymentsettings.HostNetwork {
	if len(input) == 0 {
		return nil
	}

	v := input[0]

	return &deploymentsettings.HostNetwork{
		Intents:                       expandDeploymentSettingHostNetworkIntent(v.Intent),
		EnableStorageAutoIP:           pointer.To(v.StorageAutoIpEnabled),
		StorageConnectivitySwitchless: pointer.To(v.StorageConnectivitySwitchlessEnabled),
		StorageNetworks:               expandDeploymentSettingHostNetworkStorageNetwork(v.StorageNetwork),
	}
}

func expandDeploymentSettingHostNetworkIntent(input []HostNetworkIntentModel) *[]deploymentsettings.Intents {
	if len(input) == 0 {
		return nil
	}

	results := make([]deploymentsettings.Intents, 0, len(input))
	for _, item := range input {
		results = append(results, deploymentsettings.Intents{
			Adapter:                             pointer.To(item.Adapter),
			AdapterPropertyOverrides:            expandHostNetworkIntentAdapterPropertyOverride(item.AdapterPropertyOverride),
			Name:                                pointer.To(item.Name),
			OverrideAdapterProperty:             pointer.To(item.AdapterPropertyOverrideEnabled),
			OverrideQosPolicy:                   pointer.To(item.QosPolicyOverrideEnabled),
			OverrideVirtualSwitchConfiguration:  pointer.To(item.VirtualSwitchConfigurationOverrideEnabled),
			QosPolicyOverrides:                  expandHostNetworkIntentQosPolicyOverride(item.QosPolicyOverride),
			TrafficType:                         pointer.To(item.TrafficType),
			VirtualSwitchConfigurationOverrides: expandHostNetworkIntentVirtualSwitchConfigurationOverride(item.VirtualSwitchConfigurationOverride),
		})
	}

	return &results
}

func expandHostNetworkIntentAdapterPropertyOverride(input []HostNetworkIntentAdapterPropertyOverrideModel) *deploymentsettings.AdapterPropertyOverrides {
	if len(input) == 0 {
		return &deploymentsettings.AdapterPropertyOverrides{
			JumboPacket:             pointer.To(""),
			NetworkDirect:           pointer.To(""),
			NetworkDirectTechnology: pointer.To(""),
		}
	}

	v := input[0]

	return &deploymentsettings.AdapterPropertyOverrides{
		JumboPacket:             pointer.To(v.JumboPacket),
		NetworkDirect:           pointer.To(v.NetworkDirect),
		NetworkDirectTechnology: pointer.To(v.NetworkDirectTechnology),
	}
}

func expandHostNetworkIntentQosPolicyOverride(input []HostNetworkIntentQosPolicyOverrideModel) *deploymentsettings.QosPolicyOverrides {
	if len(input) == 0 {
		return &deploymentsettings.QosPolicyOverrides{
			BandwidthPercentageSMB:         pointer.To(""),
			PriorityValue8021ActionCluster: pointer.To(""),
			PriorityValue8021ActionSMB:     pointer.To(""),
		}
	}

	v := input[0]

	return &deploymentsettings.QosPolicyOverrides{
		BandwidthPercentageSMB:         pointer.To(v.BandWidthPercentageSMB),
		PriorityValue8021ActionCluster: pointer.To(v.PriorityValue8021ActionCluster),
		PriorityValue8021ActionSMB:     pointer.To(v.PriorityValue8021ActionSMB),
	}
}

func expandHostNetworkIntentVirtualSwitchConfigurationOverride(input []HostNetworkIntentVirtualSwitchConfigurationOverrideModel) *deploymentsettings.VirtualSwitchConfigurationOverrides {
	if len(input) == 0 {
		return &deploymentsettings.VirtualSwitchConfigurationOverrides{
			EnableIov:              pointer.To(""),
			LoadBalancingAlgorithm: pointer.To(""),
		}
	}

	v := input[0]

	return &deploymentsettings.VirtualSwitchConfigurationOverrides{
		EnableIov:              pointer.To(v.EnableIov),
		LoadBalancingAlgorithm: pointer.To(v.LoadBalancingAlgorithm),
	}
}

func expandDeploymentSettingHostNetworkStorageNetwork(input []HostNetworkStorageNetworkModel) *[]deploymentsettings.StorageNetworks {
	if len(input) == 0 {
		return nil
	}

	results := make([]deploymentsettings.StorageNetworks, 0, len(input))
	for _, item := range input {
		results = append(results, deploymentsettings.StorageNetworks{
			Name:               pointer.To(item.Name),
			NetworkAdapterName: pointer.To(item.NetworkAdapterName),
			VlanId:             pointer.To(item.VlanId),
		})
	}

	return &results
}

func expandDeploymentSettingInfrastructureNetwork(input []InfrastructureNetworkModel) *[]deploymentsettings.InfrastructureNetwork {
	if len(input) == 0 {
		return nil
	}

	results := make([]deploymentsettings.InfrastructureNetwork, 0, len(input))
	for _, item := range input {
		results = append(results, deploymentsettings.InfrastructureNetwork{
			DnsServers: pointer.To(item.DnsServer),
			Gateway:    pointer.To(item.Gateway),
			IPPools:    expandDeploymentSettingInfrastructureNetworkIpPool(item.IpPool),
			SubnetMask: pointer.To(item.SubnetMask),
			UseDhcp:    pointer.To(item.DhcpEnabled),
		})
	}

	return &results
}

func expandDeploymentSettingInfrastructureNetworkIpPool(input []IpPoolModel) *[]deploymentsettings.IPPools {
	if len(input) == 0 {
		return nil
	}

	results := make([]deploymentsettings.IPPools, 0, len(input))
	for _, item := range input {
		results = append(results, deploymentsettings.IPPools{
			EndingAddress:   pointer.To(item.EndingAddress),
			StartingAddress: pointer.To(item.StartingAddress),
		})
	}

	return &results
}

func expandDeploymentSettingObservability(input ScaleUnitModel) *deploymentsettings.Observability {
	return &deploymentsettings.Observability{
		EpisodicDataUpload:  pointer.To(input.EpisodicDataUploadEnabled),
		EuLocation:          pointer.To(input.EuLocationEnabled),
		StreamingDataClient: pointer.To(input.StreamingDataClientEnabled),
	}
}

func expandDeploymentSettingOptionalService(input []OptionalServiceModel) *deploymentsettings.OptionalServices {
	if len(input) == 0 {
		return nil
	}

	v := input[0]

	return &deploymentsettings.OptionalServices{
		CustomLocation: pointer.To(v.CustomLocation),
	}
}

func expandDeploymentSettingPhysicalNode(input []PhysicalNodeModel) *[]deploymentsettings.PhysicalNodes {
	if len(input) == 0 {
		return nil
	}

	results := make([]deploymentsettings.PhysicalNodes, 0, len(input))
	for _, item := range input {
		results = append(results, deploymentsettings.PhysicalNodes{
			IPv4Address: pointer.To(item.Ipv4Address),
			Name:        pointer.To(item.Name),
		})
	}

	return &results
}

func expandDeploymentSettingSecuritySetting(input ScaleUnitModel) *deploymentsettings.DeploymentSecuritySettings {
	return &deploymentsettings.DeploymentSecuritySettings{
		BitlockerBootVolume:           pointer.To(input.BitlockerBootVolumeEnabled),
		BitlockerDataVolumes:          pointer.To(input.BitlockerDataVolumeEnabled),
		CredentialGuardEnforced:       pointer.To(input.CredentialGuardEnabled),
		DriftControlEnforced:          pointer.To(input.DriftControlEnabled),
		DrtmProtection:                pointer.To(input.DrtmProtectionEnabled),
		HvciProtection:                pointer.To(input.HvciProtectionEnabled),
		SideChannelMitigationEnforced: pointer.To(input.SideChannelMitigationEnabled),
		SmbClusterEncryption:          pointer.To(input.SmbClusterEncryptionEnabled),
		SmbSigningEnforced:            pointer.To(input.SmbSigningEnabled),
		WdacEnforced:                  pointer.To(input.WdacEnabled),
	}
}

func expandDeploymentSettingStorage(input []StorageModel) *deploymentsettings.Storage {
	if len(input) == 0 {
		return nil
	}

	v := input[0]

	return &deploymentsettings.Storage{
		ConfigurationMode: pointer.To(v.ConfigurationMode),
	}
}
