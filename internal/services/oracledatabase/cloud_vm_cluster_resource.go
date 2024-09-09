// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

package oracledatabase

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/cloudvmclusters"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracledatabase/validate"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.Resource = CloudVmClusterResource{}

type CloudVmClusterResource struct{}

type CloudVmClusterResourceModel struct {
	// Azure
	Location          string                 `tfschema:"location"`
	Name              string                 `tfschema:"name"`
	ResourceGroupName string                 `tfschema:"resource_group_name"`
	Tags              map[string]interface{} `tfschema:"tags"`

	// Required
	CloudExadataInfrastructureId string   `tfschema:"cloud_exadata_infrastructure_id"`
	CpuCoreCount                 int64    `tfschema:"cpu_core_count"`
	DataStorageSizeInTbs         float64  `tfschema:"data_storage_size_in_tbs"`
	DbNodeStorageSizeInGbs       int64    `tfschema:"db_node_storage_size_in_gbs"`
	DbServers                    []string `tfschema:"db_servers"`
	DisplayName                  string   `tfschema:"display_name"`
	GiVersion                    string   `tfschema:"gi_version"`
	Hostname                     string   `tfschema:"hostname"`
	LicenseModel                 string   `tfschema:"license_model"`
	MemorySizeInGbs              int64    `tfschema:"memory_size_in_gbs"`
	SshPublicKeys                []string `tfschema:"ssh_public_keys"`
	SubnetId                     string   `tfschema:"subnet_id"`
	VnetId                       string   `tfschema:"vnet_id"`

	// Optional
	BackupSubnetCidr         string                       `tfschema:"backup_subnet_cidr"`
	ClusterName              string                       `tfschema:"cluster_name"`
	DataCollectionOptions    []DataCollectionOptionsModel `tfschema:"data_collection_options"`
	DataStoragePercentage    int64                        `tfschema:"data_storage_percentage"`
	IsLocalBackupEnabled     bool                         `tfschema:"is_local_backup_enabled"`
	IsSparseDiskgroupEnabled bool                         `tfschema:"is_sparse_diskgroup_enabled"`
	TimeZone                 string                       `tfschema:"time_zone"`
}

func (CloudVmClusterResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		// Azure
		"location": commonschema.Location(),

		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.Name,
		},

		"resource_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		// Required
		"cloud_exadata_infrastructure_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"cpu_core_count": {
			Type:     pluginsdk.TypeInt,
			Required: true,
		},

		"data_storage_size_in_tbs": {
			Type:     pluginsdk.TypeFloat,
			Optional: true,
			Computed: true,
			ForceNew: true,
		},

		"db_node_storage_size_in_gbs": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Computed: true,
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
			ValidateFunc: validate.Name,
		},

		"gi_version": {
			Type:             pluginsdk.TypeString,
			Required:         true,
			ForceNew:         true,
			DiffSuppressFunc: GiVersionDiffSuppress,
		},

		"hostname": {
			Type:             pluginsdk.TypeString,
			Required:         true,
			ForceNew:         true,
			DiffSuppressFunc: DbSystemHostnameDiffSuppress,
		},

		"license_model": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"memory_size_in_gbs": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Computed: true,
		},

		"ssh_public_keys": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"subnet_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"vnet_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		// Optional
		"backup_subnet_cidr": {
			Type:     pluginsdk.TypeString,
			Optional: true,
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
			MaxItems: 1,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"is_diagnostics_events_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Computed: true,
					},

					"is_health_monitoring_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Computed: true,
					},

					"is_incident_logs_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Computed: true,
					},
				},
			},
		},

		"data_storage_percentage": {
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
			ForceNew: true,
		},

		"is_local_backup_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true,
			ForceNew: true,
		},

		"is_sparse_diskgroup_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true,
			ForceNew: true,
		},

		"time_zone": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
			ForceNew: true,
		},

		"tags": commonschema.Tags(),
	}
}

func (CloudVmClusterResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (CloudVmClusterResource) ModelObject() interface{} {
	return &CloudVmClusterResource{}
}

func (CloudVmClusterResource) ResourceType() string {
	return "azurerm_oracledatabase_cloud_vm_cluster"
}

func (r CloudVmClusterResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.OracleDatabase.OracleDatabaseClient.CloudVMClusters
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model CloudVmClusterResourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			id := cloudvmclusters.NewCloudVMClusterID(subscriptionId,
				model.ResourceGroupName,
				model.Name)

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
				Tags:     tags.Expand(model.Tags),
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
			if model.DataCollectionOptions != nil && len(model.DataCollectionOptions) > 0 {
				param.Properties.DataCollectionOptions = &cloudvmclusters.DataCollectionOptions{
					IsDiagnosticsEventsEnabled: pointer.To(model.DataCollectionOptions[0].IsDiagnosticsEventsEnabled),
					IsHealthMonitoringEnabled:  pointer.To(model.DataCollectionOptions[0].IsHealthMonitoringEnabled),
					IsIncidentLogsEnabled:      pointer.To(model.DataCollectionOptions[0].IsIncidentLogsEnabled),
				}
			}
			if model.TimeZone != "" {
				param.Properties.ClusterName = pointer.To(model.TimeZone)
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

			client := metadata.Client.OracleDatabase.OracleDatabaseClient.CloudVMClusters
			id, err := cloudvmclusters.ParseCloudVMClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model CloudVmClusterResourceModel
			if err = metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding err: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving exists when updating: +%v", *id)
			}
			if existing.Model == nil && existing.Model.Properties == nil {
				return fmt.Errorf("retrieving as nil when updating for %v", *id)
			}

			if metadata.ResourceData.HasChange("tags") {
				update := &cloudvmclusters.CloudVMClusterUpdate{
					Tags: tags.Expand(model.Tags),
				}
				err = client.UpdateThenPoll(ctx, *id, *update)
				if err != nil {
					return fmt.Errorf("updating %s: %v", id, err)
				}
			} else if metadata.ResourceData.HasChangesExcept("tags") {
				return fmt.Errorf("only `tags` currently support updates")
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

			client := metadata.Client.OracleDatabase.OracleDatabaseClient.CloudVMClusters
			result, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(result.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return err
			}

			if result.Model == nil {
				return fmt.Errorf("retrieving %s got nil model", id)
			}
			var output CloudVmClusterResourceModel

			// Azure
			output.Name = pointer.ToString(result.Model.Name)
			output.Location = result.Model.Location
			output.Tags = utils.FlattenPtrMapStringString(result.Model.Tags)
			output.ResourceGroupName = id.ResourceGroupName
			// Required
			output.CloudExadataInfrastructureId = result.Model.Properties.CloudExadataInfrastructureId
			output.CpuCoreCount = result.Model.Properties.CpuCoreCount
			output.DataStorageSizeInTbs = pointer.From(result.Model.Properties.DataStorageSizeInTbs)
			output.DbNodeStorageSizeInGbs = pointer.From(result.Model.Properties.DbNodeStorageSizeInGbs)
			output.DbServers = pointer.From(result.Model.Properties.DbServers)
			output.DisplayName = result.Model.Properties.DisplayName
			output.GiVersion = result.Model.Properties.GiVersion
			output.Hostname = result.Model.Properties.Hostname
			output.LicenseModel = string(pointer.From(result.Model.Properties.LicenseModel))
			output.MemorySizeInGbs = pointer.From(result.Model.Properties.MemorySizeInGbs)
			//output.SshPublicKeys = result.Model.Properties.SshPublicKeys
			tmp := make([]string, 0)
			for _, key := range result.Model.Properties.SshPublicKeys {
				if key != "" {
					tmp = append(tmp, key)
				}
			}
			output.SshPublicKeys = tmp
			output.SubnetId = result.Model.Properties.SubnetId
			output.VnetId = result.Model.Properties.VnetId
			// Optional
			output.BackupSubnetCidr = pointer.From(result.Model.Properties.BackupSubnetCidr)
			output.ClusterName = pointer.From(result.Model.Properties.ClusterName)
			output.DataCollectionOptions = ConvertDataCollectionOptionsToInternal(result.Model.Properties.DataCollectionOptions)
			output.DataStoragePercentage = pointer.From(result.Model.Properties.DataStoragePercentage)
			output.IsLocalBackupEnabled = pointer.From(result.Model.Properties.IsLocalBackupEnabled)
			output.IsSparseDiskgroupEnabled = pointer.From(result.Model.Properties.IsSparseDiskgroupEnabled)
			output.TimeZone = pointer.From(result.Model.Properties.TimeZone)

			return metadata.Encode(&output)
		},
	}
}

func (CloudVmClusterResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.OracleDatabase.OracleDatabaseClient.CloudVMClusters

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
