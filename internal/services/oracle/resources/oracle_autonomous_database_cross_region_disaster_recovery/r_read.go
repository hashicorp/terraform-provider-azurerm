// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package oracle_autonomous_database_cross_region_disaster_recovery

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/autonomousdatabases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle/helpers"
)

func (Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := autonomousdatabases.ParseAutonomousDatabaseID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			client := metadata.Client.Oracle.OracleClient.AutonomousDatabases
			result, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(result.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := AutonomousDatabaseCrossRegionDisasterRecoveryResourceModel{
				Name:              id.AutonomousDatabaseName,
				ResourceGroupName: id.ResourceGroupName,
			}
			if model := result.Model; model != nil {
				props, ok := model.Properties.(autonomousdatabases.AutonomousDatabaseCrossRegionDisasterRecoveryProperties)
				if !ok {
					return fmt.Errorf("%s was not of type `CrossRegionDisasterRecovery`", id)
				}

				state.ReplicateAutomaticBackupsEnabled = pointer.From(props.IsReplicateAutomaticBackups)
				state.RemoteDisasterRecoveryType = string(props.RemoteDisasterRecoveryType)

				if props.SourceId != "" {
					sourceDatabaseID, err := autonomousdatabases.ParseAutonomousDatabaseIDInsensitively(props.SourceId)
					if err != nil {
						return err
					}
					state.SourceAutonomousDatabaseId = sourceDatabaseID.ID()
				}

				state.AutoScalingEnabled = pointer.From(props.IsAutoScalingEnabled)
				state.BackupRetentionPeriodInDays = pointer.From(props.BackupRetentionPeriodInDays)
				state.AutoScalingForStorageEnabled = pointer.From(props.IsAutoScalingForStorageEnabled)
				state.CharacterSet = pointer.From(props.CharacterSet)
				state.ComputeCount = pointer.From(props.ComputeCount)
				state.ComputeModel = pointer.FromEnum(props.ComputeModel)
				state.CustomerContacts = helpers.FlattenAdbsCustomerContacts(props.CustomerContacts)
				state.DataStorageSizeInTb = pointer.From(props.DataStorageSizeInTbs)
				state.DatabaseWorkload = pointer.FromEnum(props.DbWorkload)
				state.DatabaseVersion = pointer.From(props.DbVersion)
				state.DisplayName = pointer.From(props.DisplayName)
				state.LicenseModel = pointer.FromEnum(props.LicenseModel)
				state.Location = location.Normalize(model.Location)
				state.MtlsConnectionRequired = pointer.From(props.IsMtlsConnectionRequired)
				state.NationalCharacterSet = pointer.From(props.NcharacterSet)

				if props.SubnetId != nil {
					subnetID, err := commonids.ParseSubnetIDInsensitively(*props.SubnetId)
					if err != nil {
						return err
					}
					state.SubnetId = subnetID.ID()
				}

				state.Tags = pointer.From(model.Tags)
			}
			return metadata.Encode(&state)
		},
	}
}
