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
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/managedbackupshorttermretentionpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/manageddatabases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/managedinstancelongtermretentionpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/managedinstances"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssqlmanagedinstance/parse"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQLManagedInstance.ManagedDatabasesClient
			instancesClient := metadata.Client.MSSQLManagedInstance.ManagedInstancesClient
			longTermRetentionClient := metadata.Client.MSSQLManagedInstance.ManagedInstancesLongTermRetentionPoliciesClient
			shortTermRetentionClient := metadata.Client.MSSQLManagedInstance.ManagedInstancesShortTermRetentionPoliciesClient

			var model MsSqlManagedDatabaseModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			managedInstanceId, err := commonids.ParseSqlManagedInstanceID(model.ManagedInstanceId)
			if err != nil {
				return err
			}

			id := commonids.NewSqlManagedInstanceDatabaseID(managedInstanceId.SubscriptionId,
				managedInstanceId.ResourceGroupName, managedInstanceId.ManagedInstanceName, model.Name)

			managedInstance, err := instancesClient.Get(ctx, *managedInstanceId, managedinstances.GetOperationOptions{})
			if err != nil || managedInstance.Model == nil || managedInstance.Model.Location == "" {
				return fmt.Errorf("checking for existence and region of Managed Instance for %s: %+v", id, err)
			}

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			parameters := manageddatabases.ManagedDatabase{
				Location:   managedInstance.Model.Location,
				Properties: pointer.To(manageddatabases.ManagedDatabaseProperties{}),
				Tags:       pointer.To(model.Tags),
			}

			if len(model.PointInTimeRestore) > 0 {
				restorePointInTime := model.PointInTimeRestore[0]
				parameters.Properties.CreateMode = pointer.To(manageddatabases.ManagedDatabaseCreateModePointInTimeRestore)
				parameters.Properties.RestorePointInTime = &restorePointInTime.RestorePointInTime

				if _, err := parse.RestorableDroppedDatabaseID(restorePointInTime.SourceDatabaseId); err == nil {
					parameters.Properties.RestorableDroppedDatabaseId = pointer.To(restorePointInTime.SourceDatabaseId)
				} else {
					parameters.Properties.SourceDatabaseId = pointer.To(restorePointInTime.SourceDatabaseId)
				}
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, parameters, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)

			if len(model.LongTermRetentionPolicy) > 0 {
				longTermRetentionPolicy := managedinstancelongtermretentionpolicies.ManagedInstanceLongTermRetentionPolicy{
					Properties: expandLongTermRetentionPolicy(model.LongTermRetentionPolicy),
				}

				if err := longTermRetentionClient.CreateOrUpdateThenPoll(ctx, id, longTermRetentionPolicy); err != nil {
					return fmt.Errorf("setting Long Term Retention Policies for %s: %+v", id, err)
				}
			}

			if model.ShortTermRetentionDays > 0 {
				shortTermRetentionPolicy := managedbackupshorttermretentionpolicies.ManagedBackupShortTermRetentionPolicy{
					Properties: &managedbackupshorttermretentionpolicies.ManagedBackupShortTermRetentionPolicyProperties{
						RetentionDays: pointer.To(model.ShortTermRetentionDays),
					},
				}
				if err = shortTermRetentionClient.CreateOrUpdateThenPoll(ctx, id, shortTermRetentionPolicy); err != nil {
					return fmt.Errorf("setting Short Term Retention Policy for %s: %+v", id, err)
				}
			}

			return nil
		},
	}
}

func expandLongTermRetentionPolicy(ltrPolicy []LongTermRetentionPolicy) *managedinstancelongtermretentionpolicies.ManagedInstanceLongTermRetentionPolicyProperties {
	if len(ltrPolicy) == 0 {
		return &managedinstancelongtermretentionpolicies.ManagedInstanceLongTermRetentionPolicyProperties{}
	}

	// The API errors if `weekOfYear` is not at least 1 when `yearlyRetention` is provided.
	weekOfYear := int64(1)
	if ltrPolicy[0].WeekOfYear != 0 {
		weekOfYear = ltrPolicy[0].WeekOfYear
	}

	return &managedinstancelongtermretentionpolicies.ManagedInstanceLongTermRetentionPolicyProperties{
		WeeklyRetention:  &ltrPolicy[0].WeeklyRetention,
		MonthlyRetention: &ltrPolicy[0].MonthlyRetention,
		YearlyRetention:  &ltrPolicy[0].YearlyRetention,
		WeekOfYear:       &weekOfYear,
	}
}
