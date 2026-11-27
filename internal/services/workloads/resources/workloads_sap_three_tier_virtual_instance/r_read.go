// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package workloads_sap_three_tier_virtual_instance

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2024-09-01/sapvirtualinstances"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	workloadsClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/workloads/client"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := workloadsClient.From(metadata.Client).SAPVirtualInstances

			id, err := sapvirtualinstances.ParseSapVirtualInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			return r.flatten(metadata, id, resp.Model)
		},
	}
}

type DiskVolumeConfiguration struct {
	VolumeName    string `tfschema:"volume_name"`
	NumberOfDisks int64  `tfschema:"number_of_disks"`
	SizeGb        int64  `tfschema:"size_in_gb"`
	SkuName       string `tfschema:"sku_name"`
}

type VirtualMachineConfiguration struct {
	ImageReference []ImageReference `tfschema:"image"`
	OSProfile      []OSProfile      `tfschema:"os_profile"`
	VmSize         string           `tfschema:"virtual_machine_size"`
}

type ImageReference struct {
	Offer     string `tfschema:"offer"`
	Publisher string `tfschema:"publisher"`
	Sku       string `tfschema:"sku"`
	Version   string `tfschema:"version"`
}

type OSProfile struct {
	AdminUsername string `tfschema:"admin_username"`
	SshPrivateKey string `tfschema:"ssh_private_key"`
	SshPublicKey  string `tfschema:"ssh_public_key"`
}

type VirtualMachineResourceNames struct {
	DataDisks             []DataDisk `tfschema:"data_disk"`
	HostName              string     `tfschema:"host_name"`
	NetworkInterfaceNames []string   `tfschema:"network_interface_names"`
	OSDiskName            string     `tfschema:"os_disk_name"`
	VMName                string     `tfschema:"virtual_machine_name"`
}

type DataDisk struct {
	VolumeName string   `tfschema:"volume_name"`
	Names      []string `tfschema:"names"`
}

type ThreeTierConfiguration struct {
	ApplicationServerConfiguration []ApplicationServerConfiguration `tfschema:"application_server_configuration"`
	AppResourceGroupName           string                           `tfschema:"app_resource_group_name"`
	CentralServerConfiguration     []CentralServerConfiguration     `tfschema:"central_server_configuration"`
	DatabaseServerConfiguration    []DatabaseServerConfiguration    `tfschema:"database_server_configuration"`
	ResourceNames                  []ResourceNames                  `tfschema:"resource_names"`
	HighAvailabilityType           string                           `tfschema:"high_availability_type"`
	IsSecondaryIpEnabled           bool                             `tfschema:"secondary_ip_enabled"`
	TransportCreateAndMount        []TransportCreateAndMount        `tfschema:"transport_create_and_mount"`
}

type TransportCreateAndMount struct {
	ResourceGroupId    string `tfschema:"resource_group_id"`
	StorageAccountName string `tfschema:"storage_account_name"`
}

type ApplicationServerConfiguration struct {
	InstanceCount               int64                         `tfschema:"instance_count"`
	SubnetId                    string                        `tfschema:"subnet_id"`
	VirtualMachineConfiguration []VirtualMachineConfiguration `tfschema:"virtual_machine_configuration"`
}

type CentralServerConfiguration struct {
	InstanceCount               int64                         `tfschema:"instance_count"`
	SubnetId                    string                        `tfschema:"subnet_id"`
	VirtualMachineConfiguration []VirtualMachineConfiguration `tfschema:"virtual_machine_configuration"`
}

type DatabaseServerConfiguration struct {
	DatabaseType                string                        `tfschema:"database_type"`
	DiskVolumeConfigurations    []DiskVolumeConfiguration     `tfschema:"disk_volume_configuration"`
	InstanceCount               int64                         `tfschema:"instance_count"`
	SubnetId                    string                        `tfschema:"subnet_id"`
	VirtualMachineConfiguration []VirtualMachineConfiguration `tfschema:"virtual_machine_configuration"`
}

type ResourceNames struct {
	ApplicationServer []ApplicationServerResourceNames `tfschema:"application_server"`
	CentralServer     []CentralServerResourceNames     `tfschema:"central_server"`
	DatabaseServer    []DatabaseServerResourceNames    `tfschema:"database_server"`
	SharedStorage     []SharedStorage                  `tfschema:"shared_storage"`
}

type ApplicationServerResourceNames struct {
	AvailabilitySetName string                        `tfschema:"availability_set_name"`
	VirtualMachines     []VirtualMachineResourceNames `tfschema:"virtual_machine"`
}

type CentralServerResourceNames struct {
	AvailabilitySetName string                        `tfschema:"availability_set_name"`
	LoadBalancer        []LoadBalancer                `tfschema:"load_balancer"`
	VirtualMachines     []VirtualMachineResourceNames `tfschema:"virtual_machine"`
}

type DatabaseServerResourceNames struct {
	AvailabilitySetName string                        `tfschema:"availability_set_name"`
	LoadBalancer        []LoadBalancer                `tfschema:"load_balancer"`
	VirtualMachines     []VirtualMachineResourceNames `tfschema:"virtual_machine"`
}

type LoadBalancer struct {
	BackendPoolNames             []string `tfschema:"backend_pool_names"`
	FrontendIpConfigurationNames []string `tfschema:"frontend_ip_configuration_names"`
	HealthProbeNames             []string `tfschema:"health_probe_names"`
	Name                         string   `tfschema:"name"`
}

type SharedStorage struct {
	AccountName         string `tfschema:"account_name"`
	PrivateEndpointName string `tfschema:"private_endpoint_name"`
}

func hasDuplicateVolumeName(input []any) bool {
	seen := make(map[string]bool)

	for _, v := range input {
		diskVolume := v.(map[string]any)
		volumeName := diskVolume["volume_name"].(string)

		if seen[volumeName] {
			return true
		}
		seen[volumeName] = true
	}

	return false
}
