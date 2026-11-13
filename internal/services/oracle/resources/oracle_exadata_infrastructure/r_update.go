// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package oracle_exadata_infrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/cloudexadatainfrastructures"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.CloudExadataInfrastructures
			id, err := cloudexadatainfrastructures.ParseCloudExadataInfrastructureID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ExadataInfraResourceModel
			if err = metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding err: %+v", err)
			}

			if _, err = client.Get(ctx, *id); err != nil {
				return fmt.Errorf("retrieving %s: ", *id)
			}

			if metadata.ResourceData.HasChange("tags") {
				update := &cloudexadatainfrastructures.CloudExadataInfrastructureUpdate{
					Tags: pointer.To(model.Tags),
				}
				if err = client.UpdateThenPoll(ctx, *id, *update); err != nil {
					return fmt.Errorf("updating %s: %v", id, err)
				}
			}
			return nil
		},
	}
}
