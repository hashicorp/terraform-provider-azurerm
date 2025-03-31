// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

package oracle

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/cloudexadatainfrastructures"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/cloudvmclusters"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.Resource = CloudVmClusterResource{}

type CloudVmClusterResource struct{}

type CloudVmClusterResourceModel struct {
	// Azure
	Location          string            `tfschema:"location"`
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Tags              map[string]string `tfschema:"tags"`

	// Required
	CloudExadataInfrastructureId string   `tfschema:"cloud_exadata_infrastructure_id"`
	CpuCoreCount                 int64    `tfschema:"cpu_core_count"`
	DataStorageSizeInTbs         float64  `tfschema:"data_storage_size_in_tbs"`
	DbNodeStorageSizeInGbs       int64    `tfschema:"db_node_storage_size_in_gbs"`
	DbServers                    []string `tfschema:"db_servers"`
	DisplayName                  string   `tfschema:"display_name"`
	GiVersion                    string   `tfschema:"gi_version"`
	Hostname                     string   `tfschema:"hostname"`
	HostnameActual               string   `tfschema:"hostname_actual"`
	LicenseModel                 string   `tfschema:"license_model"`
	MemorySizeInGbs              int64    `tfschema:"memory_size_in_gbs"`
	SshPublicKeys                []string `tfschema:"ssh_public_keys"`
	SubnetId                     string   `tfschema:"subnet_id"`
	VnetId                       string   `tfschema:"virtual_network_id"`

	// Optional
	BackupSubnetCidr         string                       `tfschema:"backup_subnet_cidr"`
	ClusterName              string                       `tfschema:"cluster_name"`
	DataCollectionOptions    []DataCollectionOptionsModel `tfschema:"data_collection_options"`
	DataStoragePercentage    int64                        `tfschema:"data_storage_percentage"`
	Domain                   string                       `tfschema:"domain"`
	IsLocalBackupEnabled     bool                         `tfschema:"local_backup_enabled"`
	IsSparseDiskgroupEnabled bool                         `tfschema:"sparse_diskgroup_enabled"`
	Ocid                     string                       `tfschema:"ocid"`
	ScanListenerPortTcp      int64                        `tfschema:"scan_listener_port_tcp"`
	ScanListenerPortTcpSsl   int64                        `tfschema:"scan_listener_port_tcp_ssl"`
	SystemVersion            string                       `tfschema:"system_version"`
	TimeZone                 string                       `tfschema:"time_zone"`
	ZoneId                   string                       `tfschema:"zone_id"`
}

func (CloudVmClusterResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.Location(),

		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.CloudVMClusterName,
			ForceNew:     true,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		// Required
		"cloud_exadata_infrastructure_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: cloudexadatainfrastructures.ValidateCloudExadataInfrastructureID,
		},

		"cpu_core_count": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.CpuCoreCount,
		},

		"data_storage_size_in_tbs": {
			Type:         pluginsdk.TypeFloat,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validate.DataStorageSizeInTbs,
		},

		"db_node_storage_size_in_gbs": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Computed: true,
			ForceNew: true,
		},

		"db_servers": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"display_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.CloudVMClusterName,
		},

		"gi_version": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"hostname": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"license_model": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.LicenseModel,
		},

		"memory_size_in_gbs": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Computed: true,
			ForceNew: true,
		},

		"ssh_public_keys": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"subnet_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"virtual_network_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateVirtualNetworkID,
		},

		// Optional
		"backup_subnet_cidr": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsCIDR,
		},

		"cluster_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
			ForceNew: true,
		},

		"data_collection_options": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			ForceNew: true,
			MaxItems: 1,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"diagnostics_events_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Computed: true,
						ForceNew: true,
					},

					"health_monitoring_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Computed: true,
						ForceNew: true,
					},

					"incident_logs_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Computed: true,
						ForceNew: true,
					},
				},
			},
		},

		"data_storage_percentage": {
			Type:         schema.TypeInt,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validate.DataStoragePercentage,
		},

		"domain": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
			ForceNew: true,
		},

		"local_backup_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true,
			ForceNew: true,
		},

		"sparse_diskgroup_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true,
			ForceNew: true,
		},

		"scan_listener_port_tcp": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      1521,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(1024, 8999),
		},

		"scan_listener_port_tcp_ssl": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      2484,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(1024, 8999),
		},

		"system_version": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validate.SystemVersion,
		},

		"time_zone": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
			ForceNew: true,
		},

		"zone_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
			ForceNew: true,
		},

		"tags": commonschema.Tags(),
	}
}

func (CloudVmClusterResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"hostname_actual": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"ocid": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (CloudVmClusterResource) ModelObject() interface{} {
	return &CloudVmClusterResource{}
}

func (CloudVmClusterResource) ResourceType() string {
	return "azurerm_oracle_cloud_vm_cluster"
}

func (r CloudVmClusterResource) Create() sdk.ResourceFunc {
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

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
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
					LicenseModel:                 pointer.To(cloudvmclusters.LicenseModel(model.LicenseModel)),
					SshPublicKeys:                model.SshPublicKeys,
					SubnetId:                     model.SubnetId,
					VnetId:                       model.VnetId,
				},
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

			if err := client.CreateOrUpdateThenPoll(ctx, id, param); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r CloudVmClusterResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.CloudVMClusters
			id, err := cloudvmclusters.ParseCloudVMClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model CloudVmClusterResourceModel
			if err = metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding err: %+v", err)
			}

			_, err = client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if metadata.ResourceData.HasChange("tags") {
				update := cloudvmclusters.CloudVMClusterUpdate{
					Tags: pointer.To(model.Tags),
				}
				if err := client.UpdateThenPoll(ctx, *id, update); err != nil {
					return fmt.Errorf("updating %s: %+v", id, err)
				}
			}

			return nil
		},
	}
}

func (CloudVmClusterResource) Read() sdk.ResourceFunc {
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
					state.LicenseModel = string(pointer.From(props.LicenseModel))
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
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (CloudVmClusterResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.CloudVMClusters

			id, err := cloudvmclusters.ParseCloudVMClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err = client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (CloudVmClusterResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return cloudvmclusters.ValidateCloudVMClusterID
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
