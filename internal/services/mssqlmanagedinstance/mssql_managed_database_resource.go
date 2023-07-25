// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssqlmanagedinstance

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v5.0/sql" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/helper"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssqlmanagedinstance/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MsSqlManagedDatabaseModel struct {
	Name                    string                    `tfschema:"name"`
	ManagedInstanceId       string                    `tfschema:"managed_instance_id"`
	LongTermRetentionPolicy []LongTermRetentionPolicy `tfschema:"long_term_retention_policy"`
	ShortTermRetentionDays  int32                     `tfschema:"short_term_retention_days"`
}

type LongTermRetentionPolicy struct {
	WeeklyRetention  string `tfschema:"weekly_retention"`
	MonthlyRetention string `tfschema:"monthly_retention"`
	YearlyRetention  string `tfschema:"yearly_retention"`
	WeekOfYear       int32  `tfschema:"week_of_year"`
}

var _ sdk.Resource = MsSqlManagedDatabaseResource{}
var _ sdk.ResourceWithUpdate = MsSqlManagedDatabaseResource{}

type MsSqlManagedDatabaseResource struct{}

func (r MsSqlManagedDatabaseResource) ResourceType() string {
	return "azurerm_mssql_managed_database"
}

func (r MsSqlManagedDatabaseResource) ModelObject() interface{} {
	return &MsSqlManagedDatabaseModel{}
}

func (r MsSqlManagedDatabaseResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.ManagedDatabaseID
}

func (r MsSqlManagedDatabaseResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ValidateMsSqlManagedInstanceDatabaseName,
		},

		"managed_instance_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ManagedInstanceID,
		},

		"long_term_retention_policy": helper.LongTermRetentionPolicySchema(),

		"short_term_retention_days": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(1, 35),
			Default:      7,
		},
	}
}

func (r MsSqlManagedDatabaseResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r MsSqlManagedDatabaseResource) Create() sdk.ResourceFunc {
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

			managedInstanceId, err := parse.ManagedInstanceID(model.ManagedInstanceId)
			if err != nil {
				return fmt.Errorf("parsing `managed_instance_id`: %v", err)
			}

			id := parse.NewManagedDatabaseID(managedInstanceId.SubscriptionId,
				managedInstanceId.ResourceGroup, managedInstanceId.Name, model.Name)

			managedInstance, err := instancesClient.Get(ctx, managedInstanceId.ResourceGroup, managedInstanceId.Name, "")
			if err != nil || managedInstance.Location == nil || *managedInstance.Location == "" {
				return fmt.Errorf("checking for existence and region of Managed Instance for %s: %+v", id, err)
			}

			metadata.Logger.Infof("Import check for %s", id)
			existing, err := client.Get(ctx, id.ResourceGroup, id.ManagedInstanceName, id.DatabaseName)
			if err != nil && !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			parameters := sql.ManagedDatabase{
				Location: managedInstance.Location,
			}

			metadata.Logger.Infof("Creating %s", id)

			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ManagedInstanceName, id.DatabaseName, parameters)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for creation of %s: %+v", id, err)
			}

			if len(model.LongTermRetentionPolicy) > 0 {
				longTermRetentionProps := expandLongTermRetentionPolicy(model.LongTermRetentionPolicy)

				longTermRetentionPolicy := sql.ManagedInstanceLongTermRetentionPolicy{
					BaseLongTermRetentionPolicyProperties: &longTermRetentionProps,
				}

				longTermRetentionFuture, err := longTermRetentionClient.CreateOrUpdate(ctx, id.ResourceGroup, id.ManagedInstanceName, id.DatabaseName, longTermRetentionPolicy)
				if err != nil {
					return fmt.Errorf("setting Long Term Retention Policies for %s: %+v", id, err)
				}

				if err = longTermRetentionFuture.WaitForCompletionRef(ctx, longTermRetentionClient.Client); err != nil {
					return fmt.Errorf("waiting for update of Long Term Retention Policies for %s: %+v", id, err)
				}
			}

			if model.ShortTermRetentionDays > 0 {

				shortTermRetentionPolicy := sql.ManagedBackupShortTermRetentionPolicy{
					ManagedBackupShortTermRetentionPolicyProperties: &sql.ManagedBackupShortTermRetentionPolicyProperties{
						RetentionDays: pointer.To(model.ShortTermRetentionDays),
					},
				}
				if _, err := shortTermRetentionClient.CreateOrUpdate(ctx, id.ResourceGroup, id.ManagedInstanceName, id.DatabaseName, shortTermRetentionPolicy); err != nil {
					return fmt.Errorf("setting Short Term Retention Policy for %s: %+v", id, err)
				}
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r MsSqlManagedDatabaseResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			longTermRetentionClient := metadata.Client.MSSQLManagedInstance.ManagedInstancesLongTermRetentionPoliciesClient
			shortTermRetentionClient := metadata.Client.MSSQLManagedInstance.ManagedInstancesShortTermRetentionPoliciesClient

			var model MsSqlManagedDatabaseModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			managedInstanceId, err := parse.ManagedInstanceID(model.ManagedInstanceId)
			if err != nil {
				return fmt.Errorf("parsing `managed_instance_id`: %v", err)
			}

			id := parse.NewManagedDatabaseID(managedInstanceId.SubscriptionId,
				managedInstanceId.ResourceGroup, managedInstanceId.Name, model.Name)

			d := metadata.ResourceData

			if d.HasChange("long_term_retention_policy") {
				longTermRetentionProps := expandLongTermRetentionPolicy(model.LongTermRetentionPolicy)

				longTermRetentionPolicy := sql.ManagedInstanceLongTermRetentionPolicy{
					BaseLongTermRetentionPolicyProperties: &longTermRetentionProps,
				}

				longTermRetentionFuture, err := longTermRetentionClient.CreateOrUpdate(ctx, id.ResourceGroup, id.ManagedInstanceName, id.DatabaseName, longTermRetentionPolicy)
				if err != nil {
					return fmt.Errorf("updating Long Term Retention Policies for %s: %+v", id, err)
				}

				if err = longTermRetentionFuture.WaitForCompletionRef(ctx, longTermRetentionClient.Client); err != nil {
					return fmt.Errorf("waiting for update of Long Term Retention Policies for %s: %+v", id, err)
				}
			}

			if d.HasChange("short_term_retention_days") {

				shortTermRetentionPolicy := sql.ManagedBackupShortTermRetentionPolicy{
					ManagedBackupShortTermRetentionPolicyProperties: &sql.ManagedBackupShortTermRetentionPolicyProperties{
						RetentionDays: pointer.To(model.ShortTermRetentionDays),
					},
				}
				if _, err := shortTermRetentionClient.CreateOrUpdate(ctx, id.ResourceGroup, id.ManagedInstanceName, id.DatabaseName, shortTermRetentionPolicy); err != nil {
					return fmt.Errorf("updating Short Term Retention Policy for %s: %+v", id, err)
				}
			}
			return nil
		},
	}
}

func (r MsSqlManagedDatabaseResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQLManagedInstance.ManagedDatabasesClient
			longTermRetentionClient := metadata.Client.MSSQLManagedInstance.ManagedInstancesLongTermRetentionPoliciesClient
			shortTermRetentionClient := metadata.Client.MSSQLManagedInstance.ManagedInstancesShortTermRetentionPoliciesClient

			id, err := parse.ManagedDatabaseID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("Decoding state for %s", id)
			var state MsSqlManagedDatabaseModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			result, err := client.Get(ctx, id.ResourceGroup, id.ManagedInstanceName, id.DatabaseName)
			if err != nil {
				if utils.ResponseWasNotFound(result.Response) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %v", *id, err)
			}

			managedInstanceId := parse.NewManagedInstanceID(id.SubscriptionId, id.ResourceGroup, id.ManagedInstanceName)

			model := MsSqlManagedDatabaseModel{
				Name:              id.DatabaseName,
				ManagedInstanceId: managedInstanceId.ID(),
			}

			ltrResp, err := longTermRetentionClient.Get(ctx, id.ResourceGroup, id.ManagedInstanceName, id.DatabaseName)
			if err != nil {
				return fmt.Errorf("retrieving Long Term Retention Policy for  %s: %v", *id, err)
			}

			model.LongTermRetentionPolicy = flattenLongTermRetentionPolicy(ltrResp)

			shortTermRetentionResp, err := shortTermRetentionClient.Get(ctx, id.ResourceGroup, id.ManagedInstanceName, id.DatabaseName)
			if err != nil {
				return fmt.Errorf("retrieving Short Term Retention Policy for  %s: %v", *id, err)
			}

			if shortTermRetentionResp.RetentionDays != nil {
				model.ShortTermRetentionDays = *shortTermRetentionResp.RetentionDays
			}

			return metadata.Encode(&model)
		},
	}
}

func (r MsSqlManagedDatabaseResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQLManagedInstance.ManagedDatabasesClient

			id, err := parse.ManagedDatabaseID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			future, err := client.Delete(ctx, id.ResourceGroup, id.ManagedInstanceName, id.DatabaseName)
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for deletion of %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandLongTermRetentionPolicy(ltrPolicy []LongTermRetentionPolicy) sql.BaseLongTermRetentionPolicyProperties {
	return sql.BaseLongTermRetentionPolicyProperties{
		WeeklyRetention:  &ltrPolicy[0].WeeklyRetention,
		MonthlyRetention: &ltrPolicy[0].MonthlyRetention,
		YearlyRetention:  &ltrPolicy[0].YearlyRetention,
		WeekOfYear:       &ltrPolicy[0].WeekOfYear,
	}
}

func flattenLongTermRetentionPolicy(ltrPolicy sql.ManagedInstanceLongTermRetentionPolicy) []LongTermRetentionPolicy {

	ltrModel := LongTermRetentionPolicy{}

	weeklyRetention := ""
	if ltrPolicy.WeeklyRetention != nil {
		weeklyRetention = *ltrPolicy.WeeklyRetention
	}

	monthlyRetention := ""
	if ltrPolicy.MonthlyRetention != nil {
		monthlyRetention = *ltrPolicy.MonthlyRetention
	}

	yearlyRetention := ""
	if ltrPolicy.YearlyRetention != nil {
		yearlyRetention = *ltrPolicy.YearlyRetention
	}

	ltrModel = LongTermRetentionPolicy{
		WeeklyRetention:  weeklyRetention,
		MonthlyRetention: monthlyRetention,
		YearlyRetention:  yearlyRetention,
	}

	if ltrPolicy.WeekOfYear != nil {
		ltrModel.WeekOfYear = *ltrPolicy.WeekOfYear
	}

	return []LongTermRetentionPolicy{ltrModel}
}
