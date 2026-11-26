// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package site_recovery_replication_recovery_plan

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationrecoveryplans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model SiteRecoveryReplicationRecoveryPlanModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}
			client := metadata.Client.RecoveryServices.ReplicationRecoveryPlansClient

			id, err := replicationrecoveryplans.ParseReplicationRecoveryPlanID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parse Site reocvery replication plan id: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			var groupValue []replicationrecoveryplans.RecoveryPlanGroup
			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: model is nil", *id)
			}

			if resp.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: properties is nil", *id)
			}

			if resp.Model.Properties.Groups == nil {
				return fmt.Errorf("retrieving %s: groups is nil", *id)
			}

			groupValue = *resp.Model.Properties.Groups

			if metadata.ResourceData.HasChanges("boot_recovery_group", "failover_recovery_group", "shutdown_recovery_group") {
				groupValue, err = expandRecoveryGroup(model.ShutdownRecoveryGroup, model.FailoverRecoveryGroup, model.BootRecoveryGroup)
			}

			if err != nil {
				return fmt.Errorf("expanding recovery group: %+v", err)
			}

			parameters := replicationrecoveryplans.UpdateRecoveryPlanInput{
				Properties: &replicationrecoveryplans.UpdateRecoveryPlanInputProperties{
					Groups: &groupValue,
				},
			}

			if err = client.UpdateThenPoll(ctx, *id, parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
