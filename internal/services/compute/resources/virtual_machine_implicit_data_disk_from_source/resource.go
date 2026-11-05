// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_implicit_data_disk_from_source

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Properties struct {
	Output Output `tfschema:"output"`
}

type Output struct {
	AccessSAS string `tfschema:"accessSAS"`
}

var (
	_ sdk.ResourceWithUpdate         = Resource{}
	_ sdk.ResourceWithCustomImporter = Resource{}
	_ sdk.ResourceWithCustomizeDiff  = Resource{}
)

type Resource struct{}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.DataDiskID
}

func (r Resource) ResourceType() string {
	return "azurerm_virtual_machine_implicit_data_disk_from_source"
}

func (r Resource) CustomImporter() sdk.ResourceRunFunc {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		client := metadata.Client.Compute.VirtualMachinesClient

		id, err := parse.DataDiskID(metadata.ResourceData.Id())
		if err != nil {
			return err
		}

		virtualMachineId := virtualmachines.NewVirtualMachineID(id.SubscriptionId, id.ResourceGroup, id.VirtualMachineName)
		resp, err := client.Get(ctx, virtualMachineId, virtualmachines.DefaultGetOperationOptions())
		if err != nil {
			return fmt.Errorf("retrieving %s: %+v", virtualMachineId, err)
		}

		if model := resp.Model; model != nil {
			if props := model.Properties; props != nil {
				if profile := props.StorageProfile; profile != nil {
					if dataDisks := profile.DataDisks; dataDisks != nil {
						var disk *virtualmachines.DataDisk
						for _, dataDisk := range *dataDisks {
							if pointer.From(dataDisk.Name) == id.Name {
								disk = &dataDisk
								break
							}
						}

						if disk == nil {
							return fmt.Errorf("unable to retrieve an existing data disk %s", *id)
						}

						if disk.CreateOption != virtualmachines.DiskCreateOptionTypesCopy {
							return fmt.Errorf("the value of `create_option` for the imported `azurerm_virtual_machine_implicit_data_disk_from_source` instance must be `Copy`, whereas now is %s", disk.CreateOption)
						}
					}
				}
			}
		}

		return nil
	}
}
