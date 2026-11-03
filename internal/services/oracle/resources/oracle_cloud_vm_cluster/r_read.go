// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package oracle_cloud_vm_cluster

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/cloudvmclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := cloudvmclusters.ParseCloudVMClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			client := metadata.Client.Oracle.OracleClient.CloudVMClusters
			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := CloudVmClusterResourceModel{
				Name:              id.CloudVmClusterName,
				ResourceGroupName: id.ResourceGroupName,
			}

			// Azure
			if model := resp.Model; model != nil {
				state.Location = model.Location
				state.Tags = pointer.From(model.Tags)

				if props := model.Properties; props != nil {
					state.CloudExadataInfrastructureId = props.CloudExadataInfrastructureId
					state.CpuCoreCount = props.CpuCoreCount
					state.DataStorageSizeInTbs = pointer.From(props.DataStorageSizeInTbs)
					state.DbNodeStorageSizeInGbs = pointer.From(props.DbNodeStorageSizeInGbs)
					state.DbServers = pointer.From(props.DbServers)
					state.DisplayName = props.DisplayName
					state.GiVersion = props.GiVersion
					state.Hostname = removeHostnameSuffix(props.Hostname)
					state.HostnameActual = props.Hostname
					state.LicenseModel = pointer.FromEnum(props.LicenseModel)
					state.MemorySizeInGbs = pointer.From(props.MemorySizeInGbs)
					state.SshPublicKeys = props.SshPublicKeys
					tmp := make([]string, 0)
					for _, key := range props.SshPublicKeys {
						if key != "" {
							tmp = append(tmp, key)
						}
					}
					state.SshPublicKeys = tmp
					state.SubnetId = props.SubnetId
					state.VnetId = props.VnetId
					// Optional
					state.BackupSubnetCidr = pointer.From(props.BackupSubnetCidr)
					state.ClusterName = pointer.From(props.ClusterName)
					state.DataCollectionOptions = FlattenDataCollectionOptions(props.DataCollectionOptions)
					state.DataStoragePercentage = pointer.From(props.DataStoragePercentage)
					state.Domain = pointer.From(props.Domain)
					state.Ocid = pointer.From(props.Ocid)
					state.IsLocalBackupEnabled = pointer.From(props.IsLocalBackupEnabled)
					state.IsSparseDiskgroupEnabled = pointer.From(props.IsSparseDiskgroupEnabled)
					state.ScanListenerPortTcp = pointer.From(props.ScanListenerPortTcp)
					state.ScanListenerPortTcpSsl = pointer.From(props.ScanListenerPortTcpSsl)
					state.SystemVersion = pointer.From(props.SystemVersion)
					state.TimeZone = pointer.From(props.TimeZone)
					state.ZoneId = pointer.From(props.ZoneId)
					state.FileSystemConfiguration = FlattenFileSystemConfigurationDetails(props.FileSystemConfigurationDetails)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func FlattenDataCollectionOptions(dataCollectionOptions *cloudvmclusters.DataCollectionOptions) []DataCollectionOptionsModel {
	output := make([]DataCollectionOptionsModel, 0)
	if dataCollectionOptions != nil {
		return append(output, DataCollectionOptionsModel{
			IsDiagnosticsEventsEnabled: pointer.From(dataCollectionOptions.IsDiagnosticsEventsEnabled),
			IsHealthMonitoringEnabled:  pointer.From(dataCollectionOptions.IsHealthMonitoringEnabled),
			IsIncidentLogsEnabled:      pointer.From(dataCollectionOptions.IsIncidentLogsEnabled),
		})
	}
	return output
}

func removeHostnameSuffix(hostnameActual string) string {
	suffixIndex := strings.LastIndex(hostnameActual, "-")
	if suffixIndex != -1 {
		return hostnameActual[:suffixIndex]
	} else {
		return hostnameActual
	}
}
