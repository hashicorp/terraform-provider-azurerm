// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package workloads_sap_single_node_virtual_instance

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2024-09-01/sapvirtualinstances"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) flatten(metadata sdk.ResourceMetaData, id *sapvirtualinstances.SapVirtualInstanceId, model *sapvirtualinstances.SAPVirtualInstance) error {
	state := WorkloadsSAPSingleNodeVirtualInstanceModel{
		Name:              id.SapVirtualInstanceName,
		ResourceGroupName: id.ResourceGroupName,
	}

	if model != nil {
		state.Location = location.Normalize(model.Location)

		identity, err := identity.FlattenUserAssignedMapToModel(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		state.Identity = pointer.From(identity)

		if props := model.Properties; props != nil {
			state.Environment = string(props.Environment)
			state.ManagedResourcesNetworkAccessType = pointer.FromEnum(props.ManagedResourcesNetworkAccessType)
			state.SapProduct = string(props.SapProduct)
			state.Tags = pointer.From(model.Tags)

			if config := props.Configuration; config != nil {
				if v, ok := config.(sapvirtualinstances.DeploymentWithOSConfiguration); ok {
					appLocation := ""
					if appLocationVal := v.AppLocation; appLocationVal != nil {
						appLocation = *v.AppLocation
					}
					state.AppLocation = location.Normalize(appLocation)

					sapFqdn := ""
					if osSapConfiguration := v.OsSapConfiguration; osSapConfiguration != nil {
						sapFqdn = pointer.From(osSapConfiguration.SapFqdn)
					}
					state.SapFqdn = sapFqdn

					if configuration := v.InfrastructureConfiguration; configuration != nil {
						if singleServerConfiguration, singleServerConfigurationExists := configuration.(sapvirtualinstances.SingleServerConfiguration); singleServerConfigurationExists {
							state.SingleServerConfiguration = flattenSingleServerConfiguration(singleServerConfiguration, metadata.ResourceData)
						}
					}
				}
			}

			if v := props.ManagedResourceGroupConfiguration; v != nil {
				state.ManagedResourceGroupName = pointer.From(v.Name)
			}
		}
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}

func flattenSAPSingleNodeVirtualInstanceDiskVolumeConfigurations(input *sapvirtualinstances.DiskConfiguration) []SingleServerDiskVolumeConfiguration {
	result := make([]SingleServerDiskVolumeConfiguration, 0)
	if input == nil || input.DiskVolumeConfigurations == nil {
		return result
	}

	for k, v := range *input.DiskVolumeConfigurations {
		diskVolumeConfiguration := SingleServerDiskVolumeConfiguration{
			NumberOfDisks: pointer.From(v.Count),
			SizeGb:        pointer.From(v.SizeGB),
			VolumeName:    k,
		}

		if sku := v.Sku; sku != nil {
			diskVolumeConfiguration.SkuName = pointer.FromEnum(sku.Name)
		}

		result = append(result, diskVolumeConfiguration)
	}

	return result
}

func flattenSAPSingleNodeVirtualInstanceVirtualMachineFullResourceNames(input sapvirtualinstances.SingleServerFullResourceNames) []SingleServerVirtualMachineResourceNames {
	result := make([]SingleServerVirtualMachineResourceNames, 0)

	if vm := input.VirtualMachine; vm != nil {
		vmFullResourceNames := SingleServerVirtualMachineResourceNames{
			HostName:              pointer.From(vm.HostName),
			OSDiskName:            pointer.From(vm.OsDiskName),
			VMName:                pointer.From(vm.VirtualMachineName),
			NetworkInterfaceNames: flattenSAPSingleNodeVirtualInstanceNetworkInterfaceResourceNames(vm.NetworkInterfaces),
			DataDisks:             flattenSAPSingleNodeVirtualInstanceDataDisks(vm.DataDiskNames),
		}

		result = append(result, vmFullResourceNames)
	}

	return result
}

func flattenSAPSingleNodeVirtualInstanceNetworkInterfaceResourceNames(input *[]sapvirtualinstances.NetworkInterfaceResourceNames) []string {
	result := make([]string, 0)
	if input == nil {
		return result
	}

	for _, v := range *input {
		result = append(result, pointer.From(v.NetworkInterfaceName))
	}

	return result
}

func flattenSAPSingleNodeVirtualInstanceDataDisks(input *map[string][]string) []SingleServerDataDisk {
	results := make([]SingleServerDataDisk, 0)
	if input == nil {
		return results
	}

	for k, v := range *input {
		dataDisk := SingleServerDataDisk{
			VolumeName: k,
			Names:      v,
		}

		results = append(results, dataDisk)
	}

	return results
}

func flattenSAPSingleNodeVirtualInstanceVirtualMachineConfiguration(input sapvirtualinstances.VirtualMachineConfiguration, d *pluginsdk.ResourceData) []SingleServerVirtualMachineConfiguration {
	result := make([]SingleServerVirtualMachineConfiguration, 0)

	return append(result, SingleServerVirtualMachineConfiguration{
		ImageReference: flattenSAPSingleNodeVirtualInstanceImageReference(input.ImageReference),
		OSProfile:      flattenSAPSingleNodeVirtualInstanceOSProfile(input.OsProfile, d),
		VmSize:         input.VMSize,
	})
}

func flattenSAPSingleNodeVirtualInstanceImageReference(input sapvirtualinstances.ImageReference) []SingleServerImageReference {
	result := make([]SingleServerImageReference, 0)

	return append(result, SingleServerImageReference{
		Offer:     pointer.From(input.Offer),
		Publisher: pointer.From(input.Publisher),
		Sku:       pointer.From(input.Sku),
		Version:   pointer.From(input.Version),
	})
}

func flattenSAPSingleNodeVirtualInstanceOSProfile(input sapvirtualinstances.OSProfile, d *pluginsdk.ResourceData) []SingleServerOSProfile {
	result := make([]SingleServerOSProfile, 0)

	osProfile := SingleServerOSProfile{
		AdminUsername: pointer.From(input.AdminUsername),
	}

	if osConfiguration := input.OsConfiguration; osConfiguration != nil {
		if v, ok := osConfiguration.(sapvirtualinstances.LinuxConfiguration); ok {
			if sshKeyPair := v.SshKeyPair; sshKeyPair != nil {
				osProfile.SshPrivateKey = d.Get("single_server_configuration.0.virtual_machine_configuration.0.os_profile.0.ssh_private_key").(string)
				osProfile.SshPublicKey = pointer.From(sshKeyPair.PublicKey)
			}
		}
	}

	return append(result, osProfile)
}

func flattenSingleServerConfiguration(input sapvirtualinstances.SingleServerConfiguration, d *pluginsdk.ResourceData) []SingleServerConfiguration {
	result := make([]SingleServerConfiguration, 0)

	singleServerConfig := SingleServerConfiguration{
		AppResourceGroupName:        input.AppResourceGroup,
		DatabaseType:                pointer.FromEnum(input.DatabaseType),
		DiskVolumeConfigurations:    flattenSAPSingleNodeVirtualInstanceDiskVolumeConfigurations(input.DbDiskConfiguration),
		SubnetId:                    input.SubnetId,
		VirtualMachineConfiguration: flattenSAPSingleNodeVirtualInstanceVirtualMachineConfiguration(input.VirtualMachineConfiguration, d),
	}

	if networkConfiguration := input.NetworkConfiguration; networkConfiguration != nil {
		singleServerConfig.IsSecondaryIpEnabled = pointer.From(networkConfiguration.IsSecondaryIPEnabled)
	}

	if customResourceNames := input.CustomResourceNames; customResourceNames != nil {
		if v, ok := customResourceNames.(sapvirtualinstances.SingleServerFullResourceNames); ok {
			singleServerConfig.VirtualMachineResourceNames = flattenSAPSingleNodeVirtualInstanceVirtualMachineFullResourceNames(v)
		}
	}

	return append(result, singleServerConfig)
}
