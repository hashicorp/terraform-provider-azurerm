// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_data_disk_attachment

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Properties struct {
	Output Output `tfschema:"output"`
}

type Disk struct {
	Name                string              `tfschema:"name"`
	ID                  string              `tfschema:"id"`
	CreateOption        string              `tfschema:"create_option"`
	DiskAccessId        string              `tfschema:"disk_access_id"`
	DiskEncryptionSetId string              `tfschema:"disk_encryption_set_id"`
	DiskIOPSReadWrite   int64               `tfschema:"disk_iops_read_write"`
	DiskMBPSReadWrite   int64               `tfschema:"disk_mbps_read_write"`
	DiskSizeGB          int64               `tfschema:"disk_size_in_gb"`
	EncryptionSettings  []EncryptionSetting `tfschema:"encryption_settings"`
	Location            string              `tfschema:"location"`
	ImageReferenceID    string              `tfschema:"image_reference_id"`
	NetworkAccessPolicy string              `tfschema:"network_access_policy"`
	OSType              string              `tfschema:"os_type"`
	SourceResourceID    string              `tfschema:"source_resource_id"`
	SourceURI           string              `tfschema:"source_uri"`
	StorageAccountID    string              `tfschema:"storage_account_id"`
	StorageAccountType  string              `tfschema:"storage_account_type"`
	Tags                map[string]string   `tfschema:"tags"`
	Zones               []string            `tfschema:"zones"`
}

type Output struct {
	AccessSAS string `tfschema:"accessSAS"`
}

type EncryptionSetting struct {
	Enabled            bool                `tfschema:"enabled"`
	DiskEncryptionKeys []DiskEncryptionKey `tfschema:"disk_encryption_keys"`
	KeyEncryptionKeys  []KeyEncryptionKey  `tfschema:"key_encryption_keys"`
}

type DiskEncryptionKey struct {
	SecretURL     string `tfschema:"secret_url"`
	SourceVaultID string `tfschema:"source_vault_id"`
}

type KeyEncryptionKey struct {
	KeyURL        string `tfschema:"key_url"`
	SourceVaultID string `tfschema:"source_vault_id"`
}

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualMachineDataDiskAttachmentCreateUpdate,
		Read:   resourceVirtualMachineDataDiskAttachmentRead,
		Update: resourceVirtualMachineDataDiskAttachmentCreateUpdate,
		Delete: resourceVirtualMachineDataDiskAttachmentDelete,
		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := parse.DataDiskID(id)
			return err
		}, func(ctx context.Context, d *pluginsdk.ResourceData, meta any) ([]*pluginsdk.ResourceData, error) {
			client := meta.(*clients.Client).Compute.VirtualMachinesClient
			id, err := parse.DataDiskID(d.Id())
			if err != nil {
				return nil, err
			}

			virtualMachineId := virtualmachines.NewVirtualMachineID(id.SubscriptionId, id.ResourceGroup, id.VirtualMachineName)
			virtualMachine, err := client.Get(ctx, virtualMachineId, virtualmachines.DefaultGetOperationOptions())
			if err != nil {
				if response.WasNotFound(virtualMachine.HttpResponse) {
					return nil, fmt.Errorf("%s was not found therefore Data Disk Attachment cannot be imported", virtualMachineId)
				}

				return nil, fmt.Errorf("retrieving %s: %+v", id, err)
			}

			var disk *virtualmachines.DataDisk
			if model := virtualMachine.Model; model != nil {
				if props := model.Properties; props != nil {
					if profile := props.StorageProfile; profile != nil {
						if dataDisks := profile.DataDisks; dataDisks != nil {
							for _, dataDisk := range *dataDisks {
								if *dataDisk.Name == id.Name {
									disk = &dataDisk
									break
								}
							}
						}
					}
				}
			}

			if disk == nil {
				return nil, fmt.Errorf("data disk %s was not found", *id)
			}

			if disk.CreateOption != virtualmachines.DiskCreateOptionTypesAttach && disk.CreateOption != virtualmachines.DiskCreateOptionTypesEmpty {
				return nil, fmt.Errorf("the value of `create_option` for the imported `azurerm_virtual_machine_data_disk_attachment` instance must be `Attach` or `Empty`, whereas now is %s", disk.CreateOption)
			}

			return []*pluginsdk.ResourceData{d}, nil
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: virtualMachineDataDiskAttachmentSchema(),
	}
}
