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
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2026-08-01-preview/managedbackupshorttermretentionpolicies"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssqlmanagedinstance/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type MsSqlManagedDatabaseDataSourceModel struct {
	Name                                 string                    `tfschema:"name"`
	ResourceGroupName                    string                    `tfschema:"resource_group_name"`
	ManagedInstanceName                  string                    `tfschema:"managed_instance_name"`
	ManagedInstanceId                    string                    `tfschema:"managed_instance_id"`
	LongTermRetentionPolicy              []LongTermRetentionPolicy `tfschema:"long_term_retention_policy"`
	ShortTermRetentionDays               int64                     `tfschema:"short_term_retention_days"`
	ShortTermRetentionImmutabilityLocked bool                      `tfschema:"short_term_retention_immutability_locked"`
	PointInTimeRestore                   []PointInTimeRestore      `tfschema:"point_in_time_restore"`
}

var _ sdk.DataSource = DataSource{}

type DataSource struct{}

func (d DataSource) ResourceType() string {
	return "azurerm_mssql_managed_database"
}

func (d DataSource) Description() string {
	return "Use this data source to access information about an existing Azure SQL Azure Managed Database."
}

func (d DataSource) ModelObject() any {
	return &MsSqlManagedDatabaseDataSourceModel{}
}

func (d DataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return commonids.ValidateSqlManagedInstanceDatabaseID
}

func (d DataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of this Azure SQL Managed Database.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ValidateMsSqlManagedInstanceDatabaseName,
		},
		"managed_instance_id": {
			Description:  "The SQL Managed Instance ID.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSqlManagedInstanceID,
		},
	}
}

func (d DataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"managed_instance_name": {
			Description: "The name of the Managed Instance.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"resource_group_name": {
			Description: "The name of the Resource Group where the Azure SQL Azure Managed Instance exists.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"long_term_retention_policy": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"weekly_retention": {
						Description: "The weekly retention policy for an LTR backup in an ISO 8601 format.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"monthly_retention": {
						Description: "The monthly retention policy for an LTR backup in an ISO 8601 format.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"yearly_retention": {
						Description: "The yearly retention policy for an LTR backup in an ISO 8601 format.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"week_of_year": {
						Description: "The week of year to take the yearly backup.",
						Type:        pluginsdk.TypeInt,
						Computed:    true,
					},
				},
			},
		},

		"short_term_retention_days": {
			Description: "The backup retention period in days. This is how many days Point-in-Time Restore will be supported.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"short_term_retention_immutability_locked": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"point_in_time_restore": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"restore_point_in_time": {
						Description: "The point in time for the restore from `source_database_id`.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
					"source_database_id": {
						Description: "The source database ID that is used to restore from.",
						Type:        schema.TypeString,
						Computed:    true,
					},
				},
			},
		},
	}
}

func (d DataSource) Read() sdk.ResourceFunc {
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

			managedInstanceId, err := commonids.ParseSqlManagedInstanceID(state.ManagedInstanceId)
			if err != nil {
				return err
			}

			id := commonids.NewSqlManagedInstanceDatabaseID(subscriptionId, managedInstanceId.ResourceGroupName, managedInstanceId.ManagedInstanceName, state.Name)
			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %v", id, err)
			}

			model := MsSqlManagedDatabaseDataSourceModel{
				Name:                id.DatabaseName,
				ManagedInstanceName: managedInstanceId.ManagedInstanceName,
				ResourceGroupName:   id.ResourceGroupName,
				ManagedInstanceId:   managedInstanceId.ID(),
			}

			ltrResp, err := longTermRetentionClient.Get(ctx, id)
			if err != nil {
				return fmt.Errorf("retrieving Long Term Retention Policy for  %s: %v", id, err)
			}

			if ltrResp.Model != nil && ltrResp.Model.Properties != nil {
				model.LongTermRetentionPolicy = flattenLongTermRetentionPolicy(*ltrResp.Model.Properties)
			}

			shortTermRetentionResp, err := shortTermRetentionClient.Get(ctx, id)
			if err != nil {
				return fmt.Errorf("retrieving Short Term Retention Policy for  %s: %v", id, err)
			}

			if shortTermRetentionResp.Model != nil && shortTermRetentionResp.Model.Properties != nil {
				model.ShortTermRetentionDays = pointer.From(shortTermRetentionResp.Model.Properties.RetentionDays)
				model.ShortTermRetentionImmutabilityLocked = pointer.From(shortTermRetentionResp.Model.Properties.ImmutabilityStatus) == managedbackupshorttermretentionpolicies.ImmutabilityStatusLocked
			}

			if v, ok := metadata.ResourceData.GetOk("point_in_time_restore"); ok {
				model.PointInTimeRestore = flattenManagedDatabasePointInTimeRestore(v)
			}

			metadata.SetID(id)
			return metadata.Encode(&model)
		},
	}
}
