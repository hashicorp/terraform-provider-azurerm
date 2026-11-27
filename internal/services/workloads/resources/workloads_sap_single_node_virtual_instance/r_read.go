// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package workloads_sap_single_node_virtual_instance

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2024-09-01/sapvirtualinstances"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Workloads.SAPVirtualInstances

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

			if err := r.flatten(metadata, id, resp.Model); err != nil {
				return err
			}

			return nil
		},
	}
}

type SingleServerConfiguration struct {
	AppResourceGroupName        string                                    `tfschema:"app_resource_group_name"`
	DatabaseType                string                                    `tfschema:"database_type"`
	DiskVolumeConfigurations    []SingleServerDiskVolumeConfiguration     `tfschema:"disk_volume_configuration"`
	IsSecondaryIpEnabled        bool                                      `tfschema:"secondary_ip_enabled"`
	SubnetId                    string                                    `tfschema:"subnet_id"`
	VirtualMachineConfiguration []SingleServerVirtualMachineConfiguration `tfschema:"virtual_machine_configuration"`
	VirtualMachineResourceNames []SingleServerVirtualMachineResourceNames `tfschema:"virtual_machine_resource_names"`
}

type SingleServerDiskVolumeConfiguration struct {
	VolumeName    string `tfschema:"volume_name"`
	NumberOfDisks int64  `tfschema:"number_of_disks"`
	SizeGb        int64  `tfschema:"size_in_gb"`
	SkuName       string `tfschema:"sku_name"`
}

type SingleServerVirtualMachineConfiguration struct {
	ImageReference []SingleServerImageReference `tfschema:"image"`
	OSProfile      []SingleServerOSProfile      `tfschema:"os_profile"`
	VmSize         string                       `tfschema:"virtual_machine_size"`
}

type SingleServerImageReference struct {
	Offer     string `tfschema:"offer"`
	Publisher string `tfschema:"publisher"`
	Sku       string `tfschema:"sku"`
	Version   string `tfschema:"version"`
}

type SingleServerOSProfile struct {
	AdminUsername string `tfschema:"admin_username"`
	SshPrivateKey string `tfschema:"ssh_private_key"`
	SshPublicKey  string `tfschema:"ssh_public_key"`
}

type SingleServerVirtualMachineResourceNames struct {
	DataDisks             []SingleServerDataDisk `tfschema:"data_disk"`
	HostName              string                 `tfschema:"host_name"`
	NetworkInterfaceNames []string               `tfschema:"network_interface_names"`
	OSDiskName            string                 `tfschema:"os_disk_name"`
	VMName                string                 `tfschema:"virtual_machine_name"`
}

type SingleServerDataDisk struct {
	VolumeName string   `tfschema:"volume_name"`
	Names      []string `tfschema:"names"`
}

func hasDuplicateVolumeNameForSAPSingleNodeVirtualInstance(input []any) bool {
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
