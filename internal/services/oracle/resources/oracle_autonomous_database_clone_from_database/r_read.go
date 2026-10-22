// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package oracle_autonomous_database_clone_from_database

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/autonomousdatabases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle/helpers"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.AutonomousDatabases
			id, err := autonomousdatabases.ParseAutonomousDatabaseID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			var state AutonomousDatabaseCloneFromDatabaseResourceModel

			if v, ok := metadata.ResourceData.GetOk("refreshable_model"); ok {
				state.RefreshableModel = v.(string)
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Tags = pointer.From(model.Tags)
				state.Name = id.AutonomousDatabaseName
				state.ResourceGroupName = id.ResourceGroupName

				if cloneProps := model.Properties; cloneProps != nil {
					props, ok := cloneProps.(autonomousdatabases.AutonomousDatabaseCloneProperties)
					if !ok {
						return fmt.Errorf("%s was not of type `Clone`", id)
					}
					state.CloneType = string(props.CloneType)
					state.SourceAutonomousDatabaseId = props.SourceId

					// Base properties
					state.AdminPassword = metadata.ResourceData.Get("admin_password").(string)
					state.BackupRetentionPeriodInDays = pointer.From(props.BackupRetentionPeriodInDays)
					state.CharacterSet = pointer.From(props.CharacterSet)
					state.ComputeCount = pointer.From(props.ComputeCount)
					state.ComputeModel = pointer.FromEnum(props.ComputeModel)
					state.DataStorageSizeInTb = pointer.From(props.DataStorageSizeInTbs)
					state.DatabaseVersion = pointer.From(props.DbVersion)
					state.DatabaseWorkload = pointer.FromEnum(props.DbWorkload)
					state.DisplayName = pointer.From(props.DisplayName)
					state.AutoScalingEnabled = pointer.From(props.IsAutoScalingEnabled)
					state.AutoScalingForStorageEnabled = pointer.From(props.IsAutoScalingForStorageEnabled)
					state.MtlsConnectionRequired = pointer.From(props.IsMtlsConnectionRequired)
					state.LicenseModel = pointer.FromEnum(props.LicenseModel)
					state.NationalCharacterSet = pointer.From(props.NcharacterSet)
					state.AllowedIpAddresses = pointer.From(props.WhitelistedIPs)
					state.CustomerContacts = helpers.FlattenAdbsCustomerContacts(props.CustomerContacts)
					state.SubnetId = pointer.From(props.SubnetId)
					state.VnetId = pointer.From(props.VnetId)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func getSourceWorkloadforClone(ctx context.Context, sourceId string, metadata sdk.ResourceMetaData) (string, error) {
	id, err := autonomousdatabases.ParseAutonomousDatabaseID(sourceId)
	if err != nil {
		return "", err
	}

	client := metadata.Client.Oracle.OracleClient.AutonomousDatabases
	resp, err := client.Get(ctx, *id)
	if err != nil {
		return "", fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if resp.Model == nil {
		return "", fmt.Errorf("retrieving %s: `model` was nil", id)
	}

	if resp.Model.Properties == nil {
		return "", fmt.Errorf("retrieving %s: `properties` was nil", id)
	}

	props := resp.Model.Properties.AutonomousDatabaseBaseProperties()
	if props.DbWorkload == nil {
		return "", fmt.Errorf("unable to determine workload type for %s", id)
	}

	return string(*props.DbWorkload), nil
}
