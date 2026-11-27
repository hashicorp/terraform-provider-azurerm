// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package workloads_sap_three_tier_virtual_instance

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2024-09-01/sapvirtualinstances"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model WorkloadsSAPThreeTierVirtualInstanceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Workloads.SAPVirtualInstances
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := sapvirtualinstances.NewSapVirtualInstanceID(subscriptionId, model.ResourceGroupName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			identity, err := identity.ExpandUserAssignedMapFromModel(model.Identity)
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			parameters := &sapvirtualinstances.SAPVirtualInstance{
				Identity: identity,
				Location: location.Normalize(model.Location),
				Properties: &sapvirtualinstances.SAPVirtualInstanceProperties{
					Environment:                       sapvirtualinstances.SAPEnvironmentType(model.Environment),
					ManagedResourcesNetworkAccessType: pointer.ToEnum[sapvirtualinstances.ManagedResourcesNetworkAccessType](model.ManagedResourcesNetworkAccessType),
					SapProduct:                        sapvirtualinstances.SAPProductType(model.SapProduct),
				},
				Tags: &model.Tags,
			}

			deploymentWithOSConfiguration := &sapvirtualinstances.DeploymentWithOSConfiguration{
				AppLocation: pointer.To(location.Normalize(model.AppLocation)),
				OsSapConfiguration: &sapvirtualinstances.OsSapConfiguration{
					SapFqdn: pointer.To(model.SapFqdn),
				},
			}

			threeTierConfiguration, err := expandThreeTierConfiguration(model.ThreeTierConfiguration)
			if err != nil {
				return err
			}
			deploymentWithOSConfiguration.InfrastructureConfiguration = threeTierConfiguration

			parameters.Properties.Configuration = deploymentWithOSConfiguration

			if v := model.ManagedResourceGroupName; v != "" {
				parameters.Properties.ManagedResourceGroupConfiguration = &sapvirtualinstances.ManagedRGConfiguration{
					Name: pointer.To(v),
				}
			}

			if err := client.CreateCallbackThenPoll(ctx, id, *parameters, metadata.SetIDAndIdentityCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
				return err
			}

			return nil
		},
	}
}

func expandVirtualMachineConfiguration(input []VirtualMachineConfiguration) *sapvirtualinstances.VirtualMachineConfiguration {
	if len(input) == 0 {
		return nil
	}

	virtualMachineConfiguration := input[0]

	return &sapvirtualinstances.VirtualMachineConfiguration{
		ImageReference: pointer.From(expandImageReference(virtualMachineConfiguration.ImageReference)),
		OsProfile:      pointer.From(expandOsProfile(virtualMachineConfiguration.OSProfile)),
		VMSize:         virtualMachineConfiguration.VmSize,
	}
}

func expandImageReference(input []ImageReference) *sapvirtualinstances.ImageReference {
	if len(input) == 0 {
		return nil
	}

	imageReference := input[0]

	return &sapvirtualinstances.ImageReference{
		Offer:     pointer.To(imageReference.Offer),
		Publisher: pointer.To(imageReference.Publisher),
		Sku:       pointer.To(imageReference.Sku),
		Version:   pointer.To(imageReference.Version),
	}
}

func expandOsProfile(input []OSProfile) *sapvirtualinstances.OSProfile {
	if len(input) == 0 {
		return nil
	}

	osProfile := input[0]

	return &sapvirtualinstances.OSProfile{
		AdminUsername: pointer.To(osProfile.AdminUsername),
		OsConfiguration: &sapvirtualinstances.LinuxConfiguration{
			DisablePasswordAuthentication: pointer.To(true),
			SshKeyPair: &sapvirtualinstances.SshKeyPair{
				PrivateKey: pointer.To(osProfile.SshPrivateKey),
				PublicKey:  pointer.To(osProfile.SshPublicKey),
			},
		},
	}
}

func expandNetworkInterfaceNames(input []string) *[]sapvirtualinstances.NetworkInterfaceResourceNames {
	result := make([]sapvirtualinstances.NetworkInterfaceResourceNames, 0)
	if len(input) == 0 {
		return &result
	}

	for _, v := range input {
		networkInterfaceName := sapvirtualinstances.NetworkInterfaceResourceNames{
			NetworkInterfaceName: pointer.To(v),
		}

		result = append(result, networkInterfaceName)
	}

	return &result
}

func expandDataDisks(input []DataDisk) *map[string][]string {
	result := make(map[string][]string)
	if len(input) == 0 {
		return &result
	}

	for _, v := range input {
		result[v.VolumeName] = v.Names
	}

	return &result
}

func expandDiskVolumeConfigurations(input []DiskVolumeConfiguration) *sapvirtualinstances.DiskConfiguration {
	if len(input) == 0 {
		return nil
	}

	result := make(map[string]sapvirtualinstances.DiskVolumeConfiguration, 0)

	for _, v := range input {
		result[v.VolumeName] = sapvirtualinstances.DiskVolumeConfiguration{
			Count:  pointer.To(v.NumberOfDisks),
			SizeGB: pointer.To(v.SizeGb),
			Sku: &sapvirtualinstances.DiskSku{
				Name: pointer.ToEnum[sapvirtualinstances.DiskSkuName](v.SkuName),
			},
		}
	}

	return &sapvirtualinstances.DiskConfiguration{
		DiskVolumeConfigurations: &result,
	}
}

func expandApplicationServer(input []ApplicationServerConfiguration) *sapvirtualinstances.ApplicationServerConfiguration {
	if len(input) == 0 {
		return nil
	}

	applicationServer := input[0]

	return &sapvirtualinstances.ApplicationServerConfiguration{
		InstanceCount:               applicationServer.InstanceCount,
		SubnetId:                    applicationServer.SubnetId,
		VirtualMachineConfiguration: pointer.From(expandVirtualMachineConfiguration(applicationServer.VirtualMachineConfiguration)),
	}
}

func expandCentralServer(input []CentralServerConfiguration) *sapvirtualinstances.CentralServerConfiguration {
	if len(input) == 0 {
		return nil
	}

	centralServer := input[0]

	return &sapvirtualinstances.CentralServerConfiguration{
		InstanceCount:               centralServer.InstanceCount,
		SubnetId:                    centralServer.SubnetId,
		VirtualMachineConfiguration: pointer.From(expandVirtualMachineConfiguration(centralServer.VirtualMachineConfiguration)),
	}
}

func expandDatabaseServer(input []DatabaseServerConfiguration) *sapvirtualinstances.DatabaseConfiguration {
	if len(input) == 0 {
		return nil
	}

	databaseServer := input[0]

	result := &sapvirtualinstances.DatabaseConfiguration{
		DiskConfiguration:           expandDiskVolumeConfigurations(databaseServer.DiskVolumeConfigurations),
		InstanceCount:               databaseServer.InstanceCount,
		SubnetId:                    databaseServer.SubnetId,
		VirtualMachineConfiguration: pointer.From(expandVirtualMachineConfiguration(databaseServer.VirtualMachineConfiguration)),
	}

	if v := databaseServer.DatabaseType; v != "" {
		result.DatabaseType = pointer.ToEnum[sapvirtualinstances.SAPDatabaseType](v)
	}

	return result
}

func expandStorageConfiguration(input ThreeTierConfiguration) (*sapvirtualinstances.StorageConfiguration, error) {
	if len(input.TransportCreateAndMount) == 0 {
		return &sapvirtualinstances.StorageConfiguration{
			TransportFileShareConfiguration: sapvirtualinstances.SkipFileShareConfiguration{},
		}, nil
	}

	result := &sapvirtualinstances.StorageConfiguration{}

	if len(input.TransportCreateAndMount) != 0 {
		transportCreateAndMount, err := expandTransportCreateAndMount(input.TransportCreateAndMount)
		if err != nil {
			return nil, err
		}
		result.TransportFileShareConfiguration = transportCreateAndMount
	}

	return result, nil
}

func expandTransportCreateAndMount(input []TransportCreateAndMount) (*sapvirtualinstances.CreateAndMountFileShareConfiguration, error) {
	if len(input) == 0 {
		return nil, nil
	}

	transportCreateAndMount := input[0]

	result := &sapvirtualinstances.CreateAndMountFileShareConfiguration{}

	if v := transportCreateAndMount.ResourceGroupId; v != "" {
		resourceGroupId, err := commonids.ParseResourceGroupID(v)
		if err != nil {
			return nil, err
		}
		result.ResourceGroup = pointer.To(resourceGroupId.ResourceGroupName)
	}

	if v := transportCreateAndMount.StorageAccountName; v != "" {
		result.StorageAccountName = pointer.To(v)
	}

	return result, nil
}

func expandResourceNames(input []ResourceNames) *sapvirtualinstances.ThreeTierFullResourceNames {
	if len(input) == 0 {
		return nil
	}

	resourceNames := input[0]

	return &sapvirtualinstances.ThreeTierFullResourceNames{
		ApplicationServer: expandApplicationServerResourceNames(resourceNames.ApplicationServer),
		CentralServer:     expandCentralServerResourceNames(resourceNames.CentralServer),
		DatabaseServer:    expandDatabaseServerResourceNames(resourceNames.DatabaseServer),
		SharedStorage:     expandSharedStorage(resourceNames.SharedStorage),
	}
}

func expandApplicationServerResourceNames(input []ApplicationServerResourceNames) *sapvirtualinstances.ApplicationServerFullResourceNames {
	if len(input) == 0 {
		return nil
	}

	applicationServerResourceNames := input[0]

	result := &sapvirtualinstances.ApplicationServerFullResourceNames{
		VirtualMachines: expandVirtualMachinesResourceNames(applicationServerResourceNames.VirtualMachines),
	}

	if v := applicationServerResourceNames.AvailabilitySetName; v != "" {
		result.AvailabilitySetName = pointer.To(v)
	}

	return result
}

func expandVirtualMachinesResourceNames(input []VirtualMachineResourceNames) *[]sapvirtualinstances.VirtualMachineResourceNames {
	result := make([]sapvirtualinstances.VirtualMachineResourceNames, 0)
	if len(input) == 0 {
		return &result
	}

	for _, item := range input {
		vmResourceNames := sapvirtualinstances.VirtualMachineResourceNames{
			DataDiskNames:     expandDataDisks(item.DataDisks),
			NetworkInterfaces: expandNetworkInterfaceNames(item.NetworkInterfaceNames),
		}

		if v := item.HostName; v != "" {
			vmResourceNames.HostName = pointer.To(v)
		}

		if v := item.OSDiskName; v != "" {
			vmResourceNames.OsDiskName = pointer.To(v)
		}

		if v := item.VMName; v != "" {
			vmResourceNames.VirtualMachineName = pointer.To(v)
		}

		result = append(result, vmResourceNames)
	}

	return &result
}

func expandCentralServerResourceNames(input []CentralServerResourceNames) *sapvirtualinstances.CentralServerFullResourceNames {
	if len(input) == 0 {
		return nil
	}

	centralServerResourceNames := input[0]

	result := &sapvirtualinstances.CentralServerFullResourceNames{
		LoadBalancer:    expandLoadBalancerResourceNames(centralServerResourceNames.LoadBalancer),
		VirtualMachines: expandVirtualMachinesResourceNames(centralServerResourceNames.VirtualMachines),
	}

	if v := centralServerResourceNames.AvailabilitySetName; v != "" {
		result.AvailabilitySetName = pointer.To(v)
	}

	return result
}

func expandLoadBalancerResourceNames(input []LoadBalancer) *sapvirtualinstances.LoadBalancerResourceNames {
	result := &sapvirtualinstances.LoadBalancerResourceNames{}
	if len(input) == 0 {
		return result
	}

	loadBalancerResourceNames := input[0]

	if v := loadBalancerResourceNames.Name; v != "" {
		result.LoadBalancerName = pointer.To(v)
	}

	if v := loadBalancerResourceNames.BackendPoolNames; v != nil {
		result.BackendPoolNames = &v
	}

	if v := loadBalancerResourceNames.FrontendIpConfigurationNames; v != nil {
		result.FrontendIPConfigurationNames = &v
	}

	if v := loadBalancerResourceNames.HealthProbeNames; v != nil {
		result.HealthProbeNames = &v
	}

	return result
}

func expandDatabaseServerResourceNames(input []DatabaseServerResourceNames) *sapvirtualinstances.DatabaseServerFullResourceNames {
	if len(input) == 0 {
		return nil
	}

	databaseServerResourceNames := input[0]

	result := &sapvirtualinstances.DatabaseServerFullResourceNames{
		LoadBalancer:    expandLoadBalancerResourceNames(databaseServerResourceNames.LoadBalancer),
		VirtualMachines: expandVirtualMachinesResourceNames(databaseServerResourceNames.VirtualMachines),
	}

	if v := databaseServerResourceNames.AvailabilitySetName; v != "" {
		result.AvailabilitySetName = pointer.To(v)
	}

	return result
}

func expandSharedStorage(input []SharedStorage) *sapvirtualinstances.SharedStorageResourceNames {
	result := &sapvirtualinstances.SharedStorageResourceNames{}
	if len(input) == 0 {
		return result
	}

	sharedStorage := input[0]

	if v := sharedStorage.AccountName; v != "" {
		result.SharedStorageAccountName = pointer.To(v)
	}

	if v := sharedStorage.PrivateEndpointName; v != "" {
		result.SharedStorageAccountPrivateEndPointName = pointer.To(v)
	}

	return result
}

func expandThreeTierConfiguration(input []ThreeTierConfiguration) (*sapvirtualinstances.ThreeTierConfiguration, error) {
	threeTierConfiguration := input[0]

	result := &sapvirtualinstances.ThreeTierConfiguration{
		ApplicationServer:   pointer.From(expandApplicationServer(threeTierConfiguration.ApplicationServerConfiguration)),
		AppResourceGroup:    threeTierConfiguration.AppResourceGroupName,
		CentralServer:       pointer.From(expandCentralServer(threeTierConfiguration.CentralServerConfiguration)),
		CustomResourceNames: expandResourceNames(threeTierConfiguration.ResourceNames),
		DatabaseServer:      pointer.From(expandDatabaseServer(threeTierConfiguration.DatabaseServerConfiguration)),
		NetworkConfiguration: &sapvirtualinstances.NetworkConfiguration{
			IsSecondaryIPEnabled: pointer.To(threeTierConfiguration.IsSecondaryIpEnabled),
		},
	}

	storageConfiguration, err := expandStorageConfiguration(threeTierConfiguration)
	if err != nil {
		return nil, err
	}
	result.StorageConfiguration = storageConfiguration

	if v := threeTierConfiguration.HighAvailabilityType; v != "" {
		result.HighAvailabilityConfig = &sapvirtualinstances.HighAvailabilityConfiguration{
			HighAvailabilityType: sapvirtualinstances.SAPHighAvailabilityType(v),
		}
	}

	return result, nil
}
