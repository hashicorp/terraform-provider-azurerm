// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_managed_database

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/managedbackupshorttermretentionpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/manageddatabases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/managedinstancelongtermretentionpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/managedinstances"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
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

			d := metadata.ResourceData

			if d.HasChange("tags") {
				managedInstance, err := instancesClient.Get(ctx, *managedInstanceId, managedinstances.GetOperationOptions{})
				if err != nil || managedInstance.Model == nil || managedInstance.Model.Location == "" {
					return fmt.Errorf("checking for existence and region of Managed Instance for %s: %+v", id, err)
				}

				parameters := manageddatabases.ManagedDatabase{
					Location: managedInstance.Model.Location,
					Tags:     pointer.To(model.Tags),
				}

				if err = client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
					return fmt.Errorf("updating %s: %+v", id, err)
				}
			}

			if d.HasChange("long_term_retention_policy") {
				longTermRetentionPolicy := managedinstancelongtermretentionpolicies.ManagedInstanceLongTermRetentionPolicy{
					Properties: expandLongTermRetentionPolicy(model.LongTermRetentionPolicy),
				}

				if err := longTermRetentionClient.CreateOrUpdateThenPoll(ctx, id, longTermRetentionPolicy); err != nil {
					return fmt.Errorf("updating Long Term Retention Policies for %s: %+v", id, err)
				}
			}

			if d.HasChange("short_term_retention_days") {
				shortTermRetentionPolicy := managedbackupshorttermretentionpolicies.ManagedBackupShortTermRetentionPolicy{
					Properties: &managedbackupshorttermretentionpolicies.ManagedBackupShortTermRetentionPolicyProperties{
						RetentionDays: pointer.To(model.ShortTermRetentionDays),
					},
				}
				if err = shortTermRetentionClient.CreateOrUpdateThenPoll(ctx, id, shortTermRetentionPolicy); err != nil {
					return fmt.Errorf("updating Short Term Retention Policy for %s: %+v", id, err)
				}
			}
			return nil
		},
	}
}
