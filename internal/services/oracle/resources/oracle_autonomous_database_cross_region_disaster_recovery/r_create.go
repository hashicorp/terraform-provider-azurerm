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
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 120 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.AutonomousDatabases
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model AutonomousDatabaseCrossRegionDisasterRecoveryResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := autonomousdatabases.NewAutonomousDatabaseID(subscriptionId,
				model.ResourceGroupName,
				model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			sourceId, err := autonomousdatabases.ParseAutonomousDatabaseID(model.SourceAutonomousDatabaseId)
			if err != nil {
				return err
			}
			sourceDb, err := client.Get(ctx, *sourceId)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", sourceId, err)
			}
			if sourceDb.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", sourceId)
			}
			if sourceDb.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", sourceId)
			}
			normalizedLocation := location.Normalize(model.Location)
			sourceLocation := location.Normalize(sourceDb.Model.Location)
			if normalizedLocation == sourceLocation {
				return fmt.Errorf("disaster recovery database must reside in a different region from the source database (source is %q, target is %q)", sourceLocation, model.Location)
			}

			subnetId, err := commonids.ParseSubnetID(model.SubnetId)
			if err != nil {
				return err
			}
			vnetId := commonids.NewVirtualNetworkID(subnetId.SubscriptionId, subnetId.ResourceGroupName, subnetId.VirtualNetworkName)

			sourceProps := sourceDb.Model.Properties.AutonomousDatabaseBaseProperties()

			param := autonomousdatabases.AutonomousDatabase{
				Name:     pointer.To(model.Name),
				Location: normalizedLocation,
				Tags:     pointer.To(model.Tags),
				Properties: &autonomousdatabases.AutonomousDatabaseCrossRegionDisasterRecoveryProperties{
					Source:                         autonomousdatabases.SourceCrossRegionDisasterRecovery,
					SourceId:                       model.SourceAutonomousDatabaseId,
					SourceOcid:                     sourceProps.Ocid,
					SourceLocation:                 pointer.To(sourceLocation),
					RemoteDisasterRecoveryType:     autonomousdatabases.DisasterRecoveryTypeAdg,
					IsReplicateAutomaticBackups:    pointer.To(model.ReplicateAutomaticBackupsEnabled),
					AdminPassword:                  sourceProps.AdminPassword,
					BackupRetentionPeriodInDays:    sourceProps.BackupRetentionPeriodInDays,
					CharacterSet:                   sourceProps.CharacterSet,
					ComputeCount:                   sourceProps.ComputeCount,
					ComputeModel:                   sourceProps.ComputeModel,
					CustomerContacts:               sourceProps.CustomerContacts,
					DataBaseType:                   autonomousdatabases.DataBaseTypeCrossRegionDisasterRecovery,
					DataStorageSizeInTbs:           sourceProps.DataStorageSizeInTbs,
					DbWorkload:                     sourceProps.DbWorkload,
					DbVersion:                      sourceProps.DbVersion,
					DisplayName:                    pointer.To(model.DisplayName),
					IsAutoScalingEnabled:           sourceProps.IsAutoScalingEnabled,
					IsAutoScalingForStorageEnabled: sourceProps.IsAutoScalingForStorageEnabled,
					IsMtlsConnectionRequired:       sourceProps.IsMtlsConnectionRequired,
					LicenseModel:                   sourceProps.LicenseModel,
					NcharacterSet:                  sourceProps.NcharacterSet,
					SubnetId:                       pointer.To(model.SubnetId),
					VnetId:                         pointer.To(vnetId.ID()),
				},
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, param, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
