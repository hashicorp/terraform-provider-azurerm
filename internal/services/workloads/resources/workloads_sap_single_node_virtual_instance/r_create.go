// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package workloads_sap_single_node_virtual_instance

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
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
			var model WorkloadsSAPSingleNodeVirtualInstanceModel
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
					Configuration: sapvirtualinstances.DeploymentWithOSConfiguration{
						AppLocation:                 pointer.To(location.Normalize(model.AppLocation)),
						InfrastructureConfiguration: expandSingleServerConfiguration(model.SingleServerConfiguration),
						OsSapConfiguration: &sapvirtualinstances.OsSapConfiguration{
							SapFqdn: pointer.To(model.SapFqdn),
						},
					},
					Environment:                       sapvirtualinstances.SAPEnvironmentType(model.Environment),
					ManagedResourcesNetworkAccessType: pointer.ToEnum[sapvirtualinstances.ManagedResourcesNetworkAccessType](model.ManagedResourcesNetworkAccessType),
					SapProduct:                        sapvirtualinstances.SAPProductType(model.SapProduct),
				},
				Tags: &model.Tags,
			}

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

func expandSAPSingleNodeVirtualInstanceVirtualMachineConfiguration(input []SingleServerVirtualMachineConfiguration) *sapvirtualinstances.VirtualMachineConfiguration {
	if len(input) == 0 {
		return nil
	}

	virtualMachineConfiguration := input[0]

	return &sapvirtualinstances.VirtualMachineConfiguration{
		ImageReference: pointer.From(expandSAPSingleNodeVirtualInstanceImageReference(virtualMachineConfiguration.ImageReference)),
		OsProfile:      pointer.From(expandSAPSingleNodeVirtualInstanceOsProfile(virtualMachineConfiguration.OSProfile)),
		VMSize:         virtualMachineConfiguration.VmSize,
	}
}

func expandSAPSingleNodeVirtualInstanceImageReference(input []SingleServerImageReference) *sapvirtualinstances.ImageReference {
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

func expandSAPSingleNodeVirtualInstanceOsProfile(input []SingleServerOSProfile) *sapvirtualinstances.OSProfile {
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

func expandSAPSingleNodeVirtualInstanceVirtualMachineFullResourceNames(input []SingleServerVirtualMachineResourceNames) *sapvirtualinstances.SingleServerFullResourceNames {
	if len(input) == 0 {
		return nil
	}

	virtualMachineFullResourceNames := input[0]

	result := &sapvirtualinstances.SingleServerFullResourceNames{
		VirtualMachine: &sapvirtualinstances.VirtualMachineResourceNames{
			DataDiskNames:     expandSAPSingleNodeVirtualInstanceDataDisks(virtualMachineFullResourceNames.DataDisks),
			NetworkInterfaces: expandSAPSingleNodeVirtualInstanceNetworkInterfaceNames(virtualMachineFullResourceNames.NetworkInterfaceNames),
		},
	}

	if v := virtualMachineFullResourceNames.HostName; v != "" {
		result.VirtualMachine.HostName = pointer.To(v)
	}

	if v := virtualMachineFullResourceNames.OSDiskName; v != "" {
		result.VirtualMachine.OsDiskName = pointer.To(v)
	}

	if v := virtualMachineFullResourceNames.VMName; v != "" {
		result.VirtualMachine.VirtualMachineName = pointer.To(v)
	}

	return result
}

func expandSAPSingleNodeVirtualInstanceNetworkInterfaceNames(input []string) *[]sapvirtualinstances.NetworkInterfaceResourceNames {
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

func expandSAPSingleNodeVirtualInstanceDataDisks(input []SingleServerDataDisk) *map[string][]string {
	result := make(map[string][]string)
	if len(input) == 0 {
		return &result
	}

	for _, v := range input {
		result[v.VolumeName] = v.Names
	}

	return &result
}

func expandSAPSingleNodeVirtualInstanceDiskVolumeConfigurations(input []SingleServerDiskVolumeConfiguration) *sapvirtualinstances.DiskConfiguration {
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

func expandSingleServerConfiguration(input []SingleServerConfiguration) *sapvirtualinstances.SingleServerConfiguration {
	if len(input) == 0 {
		return nil
	}

	singleServerConfiguration := input[0]

	result := &sapvirtualinstances.SingleServerConfiguration{
		AppResourceGroup:    singleServerConfiguration.AppResourceGroupName,
		CustomResourceNames: expandSAPSingleNodeVirtualInstanceVirtualMachineFullResourceNames(singleServerConfiguration.VirtualMachineResourceNames),
		DbDiskConfiguration: expandSAPSingleNodeVirtualInstanceDiskVolumeConfigurations(singleServerConfiguration.DiskVolumeConfigurations),
		NetworkConfiguration: &sapvirtualinstances.NetworkConfiguration{
			IsSecondaryIPEnabled: pointer.To(singleServerConfiguration.IsSecondaryIpEnabled),
		},
		SubnetId:                    singleServerConfiguration.SubnetId,
		VirtualMachineConfiguration: pointer.From(expandSAPSingleNodeVirtualInstanceVirtualMachineConfiguration(singleServerConfiguration.VirtualMachineConfiguration)),
	}

	if v := singleServerConfiguration.DatabaseType; v != "" {
		result.DatabaseType = pointer.ToEnum[sapvirtualinstances.SAPDatabaseType](v)
	}

	return result
}
