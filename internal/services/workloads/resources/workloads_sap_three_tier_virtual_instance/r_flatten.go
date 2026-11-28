// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package workloads_sap_three_tier_virtual_instance

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2024-09-01/sapvirtualinstances"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (Resource) flatten(metadata sdk.ResourceMetaData, id *sapvirtualinstances.SapVirtualInstanceId, model *sapvirtualinstances.SAPVirtualInstance) error {
	subscriptionId := metadata.Client.Account.SubscriptionId

	state := WorkloadsSAPThreeTierVirtualInstanceModel{
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
					state.AppLocation = location.Normalize(pointer.From(v.AppLocation))

					if osSapConfiguration := v.OsSapConfiguration; osSapConfiguration != nil {
						state.SapFqdn = pointer.From(osSapConfiguration.SapFqdn)
					}

					if configuration := v.InfrastructureConfiguration; configuration != nil {
						if threeTierConfiguration, threeTierConfigurationExists := configuration.(sapvirtualinstances.ThreeTierConfiguration); threeTierConfigurationExists {
							threeTierConfig, err := flattenThreeTierConfiguration(threeTierConfiguration, metadata.ResourceData, subscriptionId)
							if err != nil {
								return err
							}
							state.ThreeTierConfiguration = threeTierConfig
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

func flattenApplicationServer(input sapvirtualinstances.ApplicationServerConfiguration, d *pluginsdk.ResourceData, basePath string) []ApplicationServerConfiguration {
	result := make([]ApplicationServerConfiguration, 0)

	applicationServerConfig := ApplicationServerConfiguration{
		InstanceCount:               input.InstanceCount,
		SubnetId:                    input.SubnetId,
		VirtualMachineConfiguration: flattenVirtualMachineConfiguration(input.VirtualMachineConfiguration, d, basePath+".0.application_server_configuration"),
	}

	return append(result, applicationServerConfig)
}

func flattenCentralServer(input sapvirtualinstances.CentralServerConfiguration, d *pluginsdk.ResourceData, basePath string) []CentralServerConfiguration {
	result := make([]CentralServerConfiguration, 0)

	centralServerConfig := CentralServerConfiguration{
		InstanceCount:               input.InstanceCount,
		SubnetId:                    input.SubnetId,
		VirtualMachineConfiguration: flattenVirtualMachineConfiguration(input.VirtualMachineConfiguration, d, basePath+".0.central_server_configuration"),
	}

	return append(result, centralServerConfig)
}

func flattenDatabaseServer(input sapvirtualinstances.DatabaseConfiguration, d *pluginsdk.ResourceData, basePath string) []DatabaseServerConfiguration {
	result := make([]DatabaseServerConfiguration, 0)

	databaseServerConfig := DatabaseServerConfiguration{
		DiskVolumeConfigurations:    flattenDiskVolumeConfigurations(input.DiskConfiguration),
		InstanceCount:               input.InstanceCount,
		SubnetId:                    input.SubnetId,
		VirtualMachineConfiguration: flattenVirtualMachineConfiguration(input.VirtualMachineConfiguration, d, basePath+".0.database_server_configuration"),
	}

	if v := input.DatabaseType; v != nil {
		databaseServerConfig.DatabaseType = string(*v)
	}

	return append(result, databaseServerConfig)
}

func flattenResourceNames(input sapvirtualinstances.ThreeTierFullResourceNames) []ResourceNames {
	result := make([]ResourceNames, 0)

	resourceNames := ResourceNames{
		ApplicationServer: flattenApplicationServerResourceNames(input.ApplicationServer),
		CentralServer:     flattenCentralServerResourceNames(input.CentralServer),
		DatabaseServer:    flattenDatabaseServerResourceNames(input.DatabaseServer),
		SharedStorage:     flattenSharedStorage(input.SharedStorage),
	}

	return append(result, resourceNames)
}

func flattenApplicationServerResourceNames(input *sapvirtualinstances.ApplicationServerFullResourceNames) []ApplicationServerResourceNames {
	result := make([]ApplicationServerResourceNames, 0)
	if input == nil {
		return result
	}

	return append(result, ApplicationServerResourceNames{
		AvailabilitySetName: pointer.From(input.AvailabilitySetName),
		VirtualMachines:     flattenVirtualMachinesResourceNames(input.VirtualMachines),
	})
}

func flattenVirtualMachinesResourceNames(input *[]sapvirtualinstances.VirtualMachineResourceNames) []VirtualMachineResourceNames {
	result := make([]VirtualMachineResourceNames, 0)
	if input == nil {
		return result
	}

	for _, item := range *input {
		result = append(result, VirtualMachineResourceNames{
			HostName:              pointer.From(item.HostName),
			OSDiskName:            pointer.From(item.OsDiskName),
			VMName:                pointer.From(item.VirtualMachineName),
			DataDisks:             flattenDataDisks(item.DataDiskNames),
			NetworkInterfaceNames: flattenNetworkInterfaceResourceNames(item.NetworkInterfaces),
		})
	}

	return result
}

func flattenCentralServerResourceNames(input *sapvirtualinstances.CentralServerFullResourceNames) []CentralServerResourceNames {
	result := make([]CentralServerResourceNames, 0)
	if input == nil {
		return result
	}

	centralServerResourceNames := CentralServerResourceNames{
		AvailabilitySetName: pointer.From(input.AvailabilitySetName),
		LoadBalancer:        flattenLoadBalancerResourceNames(input.LoadBalancer),
		VirtualMachines:     flattenVirtualMachinesResourceNames(input.VirtualMachines),
	}

	return append(result, centralServerResourceNames)
}

func flattenLoadBalancerResourceNames(input *sapvirtualinstances.LoadBalancerResourceNames) []LoadBalancer {
	result := make([]LoadBalancer, 0)
	if input == nil {
		return result
	}

	return append(result, LoadBalancer{
		Name:                         pointer.From(input.LoadBalancerName),
		BackendPoolNames:             pointer.From(input.BackendPoolNames),
		FrontendIpConfigurationNames: pointer.From(input.FrontendIPConfigurationNames),
		HealthProbeNames:             pointer.From(input.HealthProbeNames),
	})
}

func flattenDatabaseServerResourceNames(input *sapvirtualinstances.DatabaseServerFullResourceNames) []DatabaseServerResourceNames {
	result := make([]DatabaseServerResourceNames, 0)
	if input == nil {
		return result
	}

	return append(result, DatabaseServerResourceNames{
		AvailabilitySetName: pointer.From(input.AvailabilitySetName),
		LoadBalancer:        flattenLoadBalancerResourceNames(input.LoadBalancer),
		VirtualMachines:     flattenVirtualMachinesResourceNames(input.VirtualMachines),
	})
}

func flattenSharedStorage(input *sapvirtualinstances.SharedStorageResourceNames) []SharedStorage {
	result := make([]SharedStorage, 0)
	if input == nil {
		return result
	}

	return append(result, SharedStorage{
		AccountName:         pointer.From(input.SharedStorageAccountName),
		PrivateEndpointName: pointer.From(input.SharedStorageAccountPrivateEndPointName),
	})
}

func flattenDiskVolumeConfigurations(input *sapvirtualinstances.DiskConfiguration) []DiskVolumeConfiguration {
	result := make([]DiskVolumeConfiguration, 0)
	if input == nil || input.DiskVolumeConfigurations == nil {
		return result
	}

	for k, v := range *input.DiskVolumeConfigurations {
		diskVolumeConfiguration := DiskVolumeConfiguration{
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

func flattenNetworkInterfaceResourceNames(input *[]sapvirtualinstances.NetworkInterfaceResourceNames) []string {
	result := make([]string, 0)
	if input == nil {
		return result
	}

	for _, v := range *input {
		result = append(result, pointer.From(v.NetworkInterfaceName))
	}

	return result
}

func flattenDataDisks(input *map[string][]string) []DataDisk {
	results := make([]DataDisk, 0)
	if input == nil {
		return results
	}

	for k, v := range *input {
		dataDisk := DataDisk{
			VolumeName: k,
			Names:      v,
		}

		results = append(results, dataDisk)
	}

	return results
}

func flattenVirtualMachineConfiguration(input sapvirtualinstances.VirtualMachineConfiguration, d *pluginsdk.ResourceData, basePath string) []VirtualMachineConfiguration {
	result := make([]VirtualMachineConfiguration, 0)

	return append(result, VirtualMachineConfiguration{
		ImageReference: flattenImageReference(input.ImageReference),
		OSProfile:      flattenOSProfile(input.OsProfile, d, basePath+".0.virtual_machine_configuration"),
		VmSize:         input.VMSize,
	})
}

func flattenImageReference(input sapvirtualinstances.ImageReference) []ImageReference {
	result := make([]ImageReference, 0)

	return append(result, ImageReference{
		Offer:     pointer.From(input.Offer),
		Publisher: pointer.From(input.Publisher),
		Sku:       pointer.From(input.Sku),
		Version:   pointer.From(input.Version),
	})
}

func flattenOSProfile(input sapvirtualinstances.OSProfile, d *pluginsdk.ResourceData, basePath string) []OSProfile {
	result := make([]OSProfile, 0)

	osProfile := OSProfile{
		AdminUsername: pointer.From(input.AdminUsername),
	}

	if osConfiguration := input.OsConfiguration; osConfiguration != nil {
		if v, ok := osConfiguration.(sapvirtualinstances.LinuxConfiguration); ok {
			if sshKeyPair := v.SshKeyPair; sshKeyPair != nil {
				osProfile.SshPrivateKey = d.Get(basePath + ".0.os_profile.0.ssh_private_key").(string)
				osProfile.SshPublicKey = pointer.From(sshKeyPair.PublicKey)
			}
		}
	}

	return append(result, osProfile)
}

func flattenThreeTierConfiguration(input sapvirtualinstances.ThreeTierConfiguration, d *pluginsdk.ResourceData, subscriptionId string) ([]ThreeTierConfiguration, error) {
	result := make([]ThreeTierConfiguration, 0)

	threeTierConfig := ThreeTierConfiguration{
		AppResourceGroupName:           input.AppResourceGroup,
		ApplicationServerConfiguration: flattenApplicationServer(input.ApplicationServer, d, "three_tier_configuration"),
		CentralServerConfiguration:     flattenCentralServer(input.CentralServer, d, "three_tier_configuration"),
		DatabaseServerConfiguration:    flattenDatabaseServer(input.DatabaseServer, d, "three_tier_configuration"),
	}

	if customResourceNames := input.CustomResourceNames; customResourceNames != nil {
		if v, ok := customResourceNames.(sapvirtualinstances.ThreeTierFullResourceNames); ok {
			threeTierConfig.ResourceNames = flattenResourceNames(v)
		}
	}

	if v := input.HighAvailabilityConfig; v != nil && v.HighAvailabilityType != "" {
		threeTierConfig.HighAvailabilityType = string(v.HighAvailabilityType)
	}

	if v := input.NetworkConfiguration; v != nil && v.IsSecondaryIPEnabled != nil {
		threeTierConfig.IsSecondaryIpEnabled = *v.IsSecondaryIPEnabled
	}

	if storageConfiguration := input.StorageConfiguration; storageConfiguration != nil {
		if transportFileShareConfiguration := storageConfiguration.TransportFileShareConfiguration; transportFileShareConfiguration != nil {
			if createAndMountFileShareConfiguration, ok := transportFileShareConfiguration.(sapvirtualinstances.CreateAndMountFileShareConfiguration); ok {
				transportCreateAndMount := TransportCreateAndMount{
					StorageAccountName: pointer.From(createAndMountFileShareConfiguration.StorageAccountName),
				}

				var resourceGroupId string
				if v := createAndMountFileShareConfiguration.ResourceGroup; v != nil {
					resourceGroupId = commonids.NewResourceGroupID(subscriptionId, *createAndMountFileShareConfiguration.ResourceGroup).ID()
				}
				transportCreateAndMount.ResourceGroupId = resourceGroupId

				threeTierConfig.TransportCreateAndMount = []TransportCreateAndMount{
					transportCreateAndMount,
				}
			}

			if _, ok := transportFileShareConfiguration.(sapvirtualinstances.MountFileShareConfiguration); ok {
				return nil, errors.New("currently, the last segment of the Storage File Share resource manager ID in Swagger is defined as `/shares/` but it's unexpected. The last segment of the Storage File Share resource manager ID should be `/fileshares/` not `/shares/` in Swagger since the backend service is using `/fileshares/`. See more details from https://github.com/Azure/azure-rest-api-specs/issues/25209. So the feature of `TransportMount` isn't supported by TF for now due to this service API bug")
			}
		}
	}

	return append(result, threeTierConfig), nil
}
