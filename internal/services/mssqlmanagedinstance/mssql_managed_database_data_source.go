// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssqlmanagedinstance

import (
	"context"
	"fmt"
	"time"

	// nolint: staticcheck

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssqlmanagedinstance/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MsSqlManagedDatabaseDataSourceModel struct {
	Name                    string                    `tfschema:"name"`
	ResourceGroupName       string                    `tfschema:"resource_group_name"`
	ManagedInstanceName     string                    `tfschema:"managed_instance_name"`
	ManagedInstanceId       string                    `tfschema:"managed_instance_id"`
	LongTermRetentionPolicy []LongTermRetentionPolicy `tfschema:"long_term_retention_policy"`
	ShortTermRetentionDays  int64                     `tfschema:"short_term_retention_days"`
	PointInTimeRestore      []PointInTimeRestore      `tfschema:"point_in_time_restore"`
}

var _ sdk.DataSource = MsSqlManagedDatabaseDataSource{}

type MsSqlManagedDatabaseDataSource struct{}

func (d MsSqlManagedDatabaseDataSource) ResourceType() string {
	return "azurerm_mssql_managed_database"
}

func (d MsSqlManagedDatabaseDataSource) ModelObject() interface{} {
	return &MsSqlManagedDatabaseDataSourceModel{}
}

func (d MsSqlManagedDatabaseDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.ManagedDatabaseID
}

func (d MsSqlManagedDatabaseDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ValidateMsSqlManagedInstanceDatabaseName,
		},
		"managed_instance_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ManagedInstanceID,
		},
	}
}

func (d MsSqlManagedDatabaseDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"managed_instance_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"resource_group_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"long_term_retention_policy": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					// WeeklyRetention - The weekly retention policy for an LTR backup in an ISO 8601 format.
					"weekly_retention": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					// MonthlyRetention - The monthly retention policy for an LTR backup in an ISO 8601 format.
					"monthly_retention": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					// YearlyRetention - The yearly retention policy for an LTR backup in an ISO 8601 format.
					"yearly_retention": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					// WeekOfYear - The week of year to take the yearly backup in an ISO 8601 format.
					"week_of_year": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},

					"immutable_backups_enabled": {
						Type:     pluginsdk.TypeBool,
						Computed: true,
					},
				},
			},
		},
		// "long_term_retention_policy": helper.LongTermRetentionPolicySchema(),
		"short_term_retention_days": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},
		"point_in_time_restore": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"restore_point_in_time": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"source_database_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
	}
}

func (d MsSqlManagedDatabaseDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQLManagedInstance.ManagedDatabasesClient
			subscriptionId := metadata.Client.Account.SubscriptionId
			longTermRetentionClient := metadata.Client.MSSQLManagedInstance.ManagedInstancesLongTermRetentionPoliciesClient
			shortTermRetentionClient := metadata.Client.MSSQLManagedInstance.ManagedInstancesShortTermRetentionPoliciesClient

			var state MsSqlManagedDatabaseDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v)", err)
			}

			managedInstanceId, err := parse.ManagedInstanceID(state.ManagedInstanceId)
			if err != nil {
				return err
			}

			id := parse.NewManagedDatabaseID(subscriptionId, managedInstanceId.ResourceGroup, managedInstanceId.Name, state.Name)
			resp, err := client.Get(ctx, id.ResourceGroup, id.ManagedInstanceName, id.DatabaseName)
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %v", id, err)
			}

			model := MsSqlManagedDatabaseDataSourceModel{
				Name:                id.DatabaseName,
				ManagedInstanceName: managedInstanceId.Name,
				ResourceGroupName:   id.ResourceGroup,
				ManagedInstanceId:   managedInstanceId.ID(),
			}

			ltrResp, err := longTermRetentionClient.Get(ctx, id.ResourceGroup, id.ManagedInstanceName, id.DatabaseName)
			if err != nil {
				return fmt.Errorf("retrieving Long Term Retention Policy for  %s: %v", id, err)
			}

			model.LongTermRetentionPolicy = flattenLongTermRetentionPolicy(ltrResp)

			shortTermRetentionResp, err := shortTermRetentionClient.Get(ctx, id.ResourceGroup, id.ManagedInstanceName, id.DatabaseName)
			if err != nil {
				return fmt.Errorf("retrieving Short Term Retention Policy for  %s: %v", id, err)
			}

			model.ShortTermRetentionDays = int64(pointer.From(shortTermRetentionResp.RetentionDays))

			if v, ok := metadata.ResourceData.GetOk("point_in_time_restore"); ok {
				model.PointInTimeRestore = flattenManagedDatabasePointInTimeRestore(v)
			}

			metadata.SetID(id)
			return metadata.Encode(&model)
		},
	}
}
