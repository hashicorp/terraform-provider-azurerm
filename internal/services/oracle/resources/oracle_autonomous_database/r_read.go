// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package oracle_autonomous_database

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
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
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			client := metadata.Client.Oracle.OracleClient.AutonomousDatabases
			result, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(result.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := AutonomousDatabaseRegularResourceModel{
				Name:              id.AutonomousDatabaseName,
				ResourceGroupName: id.ResourceGroupName,
			}
			if model := result.Model; model != nil {
				props, ok := model.Properties.(autonomousdatabases.AutonomousDatabaseProperties)
				if !ok {
					return fmt.Errorf("%s was not of type `Regular`", id)
				}
				state.AdminPassword = metadata.ResourceData.Get("admin_password").(string)
				state.AutoScalingEnabled = pointer.From(props.IsAutoScalingEnabled)
				state.BackupRetentionPeriodInDays = pointer.From(props.BackupRetentionPeriodInDays)
				state.AutoScalingForStorageEnabled = pointer.From(props.IsAutoScalingForStorageEnabled)
				state.CharacterSet = pointer.From(props.CharacterSet)
				state.ComputeCount = pointer.From(props.ComputeCount)
				state.ComputeModel = pointer.FromEnum(props.ComputeModel)
				state.CustomerContacts = helpers.FlattenAdbsCustomerContacts(props.CustomerContacts)
				state.DataStorageSizeInTbs = pointer.From(props.DataStorageSizeInTbs)
				state.DbWorkload = pointer.FromEnum(props.DbWorkload)
				state.DbVersion = pointer.From(props.DbVersion)
				state.DisplayName = pointer.From(props.DisplayName)
				state.LicenseModel = pointer.FromEnum(props.LicenseModel)
				state.Location = result.Model.Location
				state.MtlsConnectionRequired = pointer.From(props.IsMtlsConnectionRequired)
				state.Name = pointer.From(result.Model.Name)
				state.NationalCharacterSet = pointer.From(props.NcharacterSet)
				state.SubnetId = pointer.From(props.SubnetId)
				state.Tags = pointer.From(result.Model.Tags)
				state.VnetId = pointer.From(props.VnetId)
				state.LongTermBackUpSchedule = helpers.FlattenLongTermBackUpScheduleDetails(props.LongTermBackupSchedule)
				state.AllowedIps = pointer.From(props.WhitelistedIPs)
				state.Ocid = pointer.From(props.Ocid)
			}
			return metadata.Encode(&state)
		},
	}
}

func (r Resource) hasGeneralUpdates(metadata sdk.ResourceMetaData) bool {
	// lintignore:R019 // deliberate subset: only the fields covered by the general update payload; the remaining attributes are updated via separate calls
	return metadata.ResourceData.HasChanges(
		"tags",
		"data_storage_size_in_tbs",
		"compute_count",
		"auto_scaling_enabled",
		"auto_scaling_for_storage_enabled",
		"allowed_ips",
	)
}
