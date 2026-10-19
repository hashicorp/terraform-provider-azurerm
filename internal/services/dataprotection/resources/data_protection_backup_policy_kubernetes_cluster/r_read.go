// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_policy_kubernetes_cluster

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/basebackuppolicyresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataProtection.BackupPolicyClient

			id, err := basebackuppolicyresources.ParseBackupPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.BackupPoliciesGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := BackupPolicyKubernatesClusterModel{
				Name:              id.BackupPolicyName,
				ResourceGroupName: id.ResourceGroupName,
				VaultName:         id.BackupVaultName,
			}

			if model := resp.Model; model != nil {
				if properties, ok := model.Properties.(basebackuppolicyresources.BackupPolicy); ok {
					state.DefaultRetentionRule = flattenBackupPolicyKubernetesClusterDefaultRetentionRule(&properties.PolicyRules)
					state.RetentionRule = flattenBackupPolicyKubernetesClusterRetentionRules(&properties.PolicyRules)
					state.BackupRepeatingTimeIntervals = flattenBackupPolicyKubernetesClusterBackupRuleArray(&properties.PolicyRules)
					state.TimeZone = flattenBackupPolicyKubernetesClusterBackupTimeZone(&properties.PolicyRules)
				}
			}

			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
				return err
			}
			return metadata.Encode(&state)
		},
	}
}

type DefaultRetentionRule struct {
	LifeCycle []LifeCycle `tfschema:"life_cycle"`
}

type RetentionRule struct {
	Name      string      `tfschema:"name"`
	Criteria  []Criteria  `tfschema:"criteria"`
	Priority  int64       `tfschema:"priority"`
	LifeCycle []LifeCycle `tfschema:"life_cycle"`
}

type LifeCycle struct {
	DataStoreType string `tfschema:"data_store_type"`
	Duration      string `tfschema:"duration"`
}

type Criteria struct {
	AbsoluteCriteria     string   `tfschema:"absolute_criteria"`
	DaysOfWeek           []string `tfschema:"days_of_week"`
	MonthsOfYear         []string `tfschema:"months_of_year"`
	ScheduledBackupTimes []string `tfschema:"scheduled_backup_times"`
	WeeksOfMonth         []string `tfschema:"weeks_of_month"`
}
