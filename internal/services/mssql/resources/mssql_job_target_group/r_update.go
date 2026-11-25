// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_job_target_group

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/jobtargetgroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQL.JobTargetGroupsClient

			id, err := jobtargetgroups.ParseTargetGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config MsSqlJobTargetGroupResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}

			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `model.Properties` was nil", id)
			}

			payload := existing.Model

			if metadata.ResourceData.HasChange("job_target") {
				payload.Properties.Members = expandJobTargets(config.JobTargets)
			}

			if _, err := client.CreateOrUpdate(ctx, *id, *payload); err != nil {
				return fmt.Errorf("updating: %s: %+v", id, err)
			}

			return nil
		},
	}
}
