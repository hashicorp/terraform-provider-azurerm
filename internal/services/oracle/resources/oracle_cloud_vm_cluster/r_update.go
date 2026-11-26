// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package oracle_cloud_vm_cluster

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/cloudvmclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.CloudVMClusters
			id, err := cloudvmclusters.ParseCloudVMClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model CloudVmClusterResourceModel
			if err = metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding err: %+v", err)
			}

			if _, err = client.Get(ctx, *id); err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			update := cloudvmclusters.CloudVMClusterUpdate{}
			if metadata.ResourceData.HasChange("tags") {
				update.Tags = pointer.To(model.Tags)
			}

			if metadata.ResourceData.HasChange("file_system_configuration") {
				update.Properties = &cloudvmclusters.CloudVMClusterUpdateProperties{
					FileSystemConfigurationDetails: ExpandFileSystemConfiguration(model.FileSystemConfiguration),
				}
			}

			if err := client.UpdateThenPoll(ctx, *id, update); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
