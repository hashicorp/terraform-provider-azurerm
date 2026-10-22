// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_managed_database

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/manageddatabases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/managedinstancelongtermretentionpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQLManagedInstance.ManagedDatabasesClient
			longTermRetentionClient := metadata.Client.MSSQLManagedInstance.ManagedInstancesLongTermRetentionPoliciesClient
			shortTermRetentionClient := metadata.Client.MSSQLManagedInstance.ManagedInstancesShortTermRetentionPoliciesClient

			id, err := commonids.ParseManagedInstanceDatabaseID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state MsSqlManagedDatabaseModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			result, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(result.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %v", *id, err)
			}

			managedInstanceId := commonids.NewSqlManagedInstanceID(id.SubscriptionId, id.ResourceGroupName, id.ManagedInstanceName)

			model := MsSqlManagedDatabaseModel{
				Name:              id.DatabaseName,
				ManagedInstanceId: managedInstanceId.ID(),
			}

			// If the database is in a `Stopped` state, attempting to retrieve the long term retention policy results in a 400.
			if result.Model != nil && result.Model.Properties != nil && pointer.From(result.Model.Properties.Status) != manageddatabases.ManagedDatabaseStatusStopped {
				ltrResp, err := longTermRetentionClient.Get(ctx, *id)
				if err != nil {
					return fmt.Errorf("retrieving Long Term Retention Policy for %s: %v", *id, err)
				}

				if ltrResp.Model != nil && ltrResp.Model.Properties != nil {
					model.LongTermRetentionPolicy = flattenLongTermRetentionPolicy(*ltrResp.Model.Properties)
				}
			} else {
				// Preserve state while database is `Stopped`
				model.LongTermRetentionPolicy = state.LongTermRetentionPolicy
			}

			shortTermRetentionResp, err := shortTermRetentionClient.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving Short Term Retention Policy for  %s: %v", *id, err)
			}

			if shortTermRetentionResp.Model != nil && shortTermRetentionResp.Model.Properties != nil {
				model.ShortTermRetentionDays = pointer.From(shortTermRetentionResp.Model.Properties.RetentionDays)
			}

			d := metadata.ResourceData
			if v, ok := d.GetOk("point_in_time_restore"); ok {
				model.PointInTimeRestore = flattenManagedDatabasePointInTimeRestore(v)
			}

			model.Tags = pointer.From(result.Model.Tags)

			return metadata.Encode(&model)
		},
	}
}

func flattenLongTermRetentionPolicy(ltrPolicy managedinstancelongtermretentionpolicies.ManagedInstanceLongTermRetentionPolicyProperties) []LongTermRetentionPolicy {
	ltrModel := LongTermRetentionPolicy{
		WeeklyRetention:  pointer.From(ltrPolicy.WeeklyRetention),
		MonthlyRetention: pointer.From(ltrPolicy.MonthlyRetention),
		YearlyRetention:  pointer.From(ltrPolicy.YearlyRetention),
		WeekOfYear:       pointer.From(ltrPolicy.WeekOfYear),
	}

	return []LongTermRetentionPolicy{ltrModel}
}

func flattenManagedDatabasePointInTimeRestore(input any) []PointInTimeRestore {
	output := make([]PointInTimeRestore, 0)

	if input == nil {
		return output
	}

	attrs := input.([]any)

	for _, attr := range attrs {
		if attr == nil {
			return output
		}

		v := attr.(map[string]any)

		output = append(output, PointInTimeRestore{
			RestorePointInTime: v["restore_point_in_time"].(string),
			SourceDatabaseId:   v["source_database_id"].(string),
		})
	}

	return output
}

type LongTermRetentionPolicy struct {
	WeeklyRetention  string `tfschema:"weekly_retention"`
	MonthlyRetention string `tfschema:"monthly_retention"`
	YearlyRetention  string `tfschema:"yearly_retention"`
	WeekOfYear       int64  `tfschema:"week_of_year"`
}

type PointInTimeRestore struct {
	RestorePointInTime string `tfschema:"restore_point_in_time"`
	SourceDatabaseId   string `tfschema:"source_database_id"`
}
