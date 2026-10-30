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

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 120 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.AutonomousDatabases
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model AutonomousDatabaseCloneFromDatabaseResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := autonomousdatabases.NewAutonomousDatabaseID(
				subscriptionId,
				model.ResourceGroupName,
				model.Name,
			)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			param := autonomousdatabases.AutonomousDatabase{
				Name:     pointer.To(model.Name),
				Location: location.Normalize(model.Location),
				Tags:     pointer.To(model.Tags),
			}

			param.Properties = &autonomousdatabases.AutonomousDatabaseCloneProperties{
				CloneType:    autonomousdatabases.CloneType(model.CloneType),
				SourceId:     model.SourceAutonomousDatabaseId,
				Source:       pointer.To(autonomousdatabases.SourceTypeDatabase),
				DataBaseType: autonomousdatabases.DataBaseTypeClone,

				AdminPassword:                  pointer.To(model.AdminPassword),
				BackupRetentionPeriodInDays:    pointer.To(model.BackupRetentionPeriodInDays),
				CharacterSet:                   pointer.To(model.CharacterSet),
				ComputeCount:                   pointer.To(model.ComputeCount),
				ComputeModel:                   pointer.ToEnum[autonomousdatabases.ComputeModel](model.ComputeModel),
				CustomerContacts:               pointer.To(helpers.ExpandCloneCustomerContacts(model.CustomerContacts)),
				DataStorageSizeInTbs:           pointer.To(model.DataStorageSizeInTb),
				DbWorkload:                     pointer.ToEnum[autonomousdatabases.WorkloadType](model.DatabaseWorkload),
				DbVersion:                      pointer.To(model.DatabaseVersion),
				DisplayName:                    pointer.To(model.DisplayName),
				IsAutoScalingEnabled:           pointer.To(model.AutoScalingEnabled),
				IsAutoScalingForStorageEnabled: pointer.To(model.AutoScalingForStorageEnabled),
				IsMtlsConnectionRequired:       pointer.To(model.MtlsConnectionRequired),
				LicenseModel:                   pointer.ToEnum[autonomousdatabases.LicenseModel](model.LicenseModel),
				NcharacterSet:                  pointer.To(model.NationalCharacterSet),
				WhitelistedIPs:                 pointer.To(model.AllowedIpAddresses),
			}

			properties := param.Properties.(*autonomousdatabases.AutonomousDatabaseCloneProperties)
			if model.RefreshableModel != "" {
				properties.RefreshableModel = pointer.ToEnum[autonomousdatabases.RefreshableModelType](model.RefreshableModel)
			}

			if model.SubnetId != "" {
				properties.SubnetId = pointer.To(model.SubnetId)
			}

			if model.VnetId != "" {
				properties.VnetId = pointer.To(model.VnetId)
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, param, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)

			return nil
		},
	}
}
