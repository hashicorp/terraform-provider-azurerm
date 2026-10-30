// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package oracle_cloud_vm_cluster

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/cloudvmclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 24 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.CloudVMClusters
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model CloudVmClusterResourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			id := cloudvmclusters.NewCloudVMClusterID(subscriptionId, model.ResourceGroupName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			param := cloudvmclusters.CloudVMCluster{
				// Azure
				Name:     pointer.To(model.Name),
				Location: model.Location,
				Tags:     pointer.To(model.Tags),
				Properties: &cloudvmclusters.CloudVMClusterProperties{
					// Required
					CloudExadataInfrastructureId: model.CloudExadataInfrastructureId,
					CpuCoreCount:                 model.CpuCoreCount,
					DbServers:                    pointer.To(model.DbServers),
					DisplayName:                  model.DisplayName,
					GiVersion:                    model.GiVersion,
					Hostname:                     model.Hostname,
					LicenseModel:                 pointer.ToEnum[cloudvmclusters.LicenseModel](model.LicenseModel),
					SshPublicKeys:                model.SshPublicKeys,
					SubnetId:                     model.SubnetId,
					VnetId:                       model.VnetId,
				},
			}

			if len(model.FileSystemConfiguration) > 0 {
				param.Properties.FileSystemConfigurationDetails = ExpandFileSystemConfiguration(model.FileSystemConfiguration)
			}

			if model.BackupSubnetCidr != "" {
				param.Properties.BackupSubnetCidr = pointer.To(model.BackupSubnetCidr)
			}
			if model.ClusterName != "" {
				param.Properties.ClusterName = pointer.To(model.ClusterName)
			}
			if len(model.DataCollectionOptions) > 0 {
				param.Properties.DataCollectionOptions = &cloudvmclusters.DataCollectionOptions{
					IsDiagnosticsEventsEnabled: pointer.To(model.DataCollectionOptions[0].IsDiagnosticsEventsEnabled),
					IsHealthMonitoringEnabled:  pointer.To(model.DataCollectionOptions[0].IsHealthMonitoringEnabled),
					IsIncidentLogsEnabled:      pointer.To(model.DataCollectionOptions[0].IsIncidentLogsEnabled),
				}
			}
			if model.Domain != "" {
				param.Properties.Domain = pointer.To(model.Domain)
			}
			if model.ScanListenerPortTcp >= 1024 && model.ScanListenerPortTcp <= 8999 {
				param.Properties.ScanListenerPortTcp = pointer.To(model.ScanListenerPortTcp)
			}
			if model.ScanListenerPortTcpSsl >= 1024 && model.ScanListenerPortTcpSsl <= 8999 {
				param.Properties.ScanListenerPortTcpSsl = pointer.To(model.ScanListenerPortTcpSsl)
			}
			if model.SystemVersion != "" {
				param.Properties.SystemVersion = pointer.To(model.SystemVersion)
			}
			if model.TimeZone != "" {
				param.Properties.TimeZone = pointer.To(model.TimeZone)
			}
			if model.ZoneId != "" {
				param.Properties.ZoneId = pointer.To(model.ZoneId)
			}
			if model.DataStoragePercentage != 0 {
				param.Properties.DataStoragePercentage = pointer.To(model.DataStoragePercentage)
			}
			if model.DataStorageSizeInTbs != 0 {
				param.Properties.DataStorageSizeInTbs = pointer.To(model.DataStorageSizeInTbs)
			}
			if model.DbNodeStorageSizeInGbs != 0 {
				param.Properties.DbNodeStorageSizeInGbs = pointer.To(model.DbNodeStorageSizeInGbs)
			}
			param.Properties.IsLocalBackupEnabled = pointer.To(model.IsLocalBackupEnabled)
			param.Properties.IsSparseDiskgroupEnabled = pointer.To(model.IsSparseDiskgroupEnabled)
			if model.MemorySizeInGbs != 0 {
				param.Properties.MemorySizeInGbs = pointer.To(model.MemorySizeInGbs)
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, param, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)

			return nil
		},
	}
}
