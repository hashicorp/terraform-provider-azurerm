// Copyright © 2025, Oracle and/or its affiliates. All rights reserved

package oracle

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-03-01/exascaledbstoragevaults"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-03-01/exadbvmclusters"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.Resource = ExadbVmClusterResource{}

type ExadbVmClusterResource struct{}

type ExadbVmClusterResourceModel struct {
	// Azure
	Location          string            `tfschema:"location"`
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Tags              map[string]string `tfschema:"tags"`
	Zones             zones.Schema      `tfschema:"zones"`

	// Required
	DisplayName              string                       `tfschema:"display_name"`
	EnabledEcpuCount         int64                        `tfschema:"enabled_ecpu_count"`
	ExascaleDbStorageVaultId string                       `tfschema:"exascale_db_storage_vault_id"`
	GridImageOcid            string                       `tfschema:"grid_image_ocid"`
	Hostname                 string                       `tfschema:"hostname"`
	NodeCount                int64                        `tfschema:"node_count"`
	Shape                    string                       `tfschema:"shape"`
	SshPublicKeys            []string                     `tfschema:"ssh_public_keys"`
	SubnetId                 string                       `tfschema:"subnet_id"`
	TotalEcpuCount           int64                        `tfschema:"total_ecpu_count"`
	VmFileSystemStorage      []ExadbVmClusterStorageModel `tfschema:"vm_file_system_storage"`
	VnetId                   string                       `tfschema:"virtual_network_id"`

	// Optional
	BackupSubnetCidr       string                            `tfschema:"backup_subnet_cidr"`
	ClusterName            string                            `tfschema:"cluster_name"`
	DataCollectionOptions  []ExadbDataCollectionOptionsModel `tfschema:"data_collection_options"`
	Domain                 string                            `tfschema:"domain"`
	LicenseModel           string                            `tfschema:"license_model"`
	NsgCidrs               []NsgCidrModel                    `tfschema:"nsg_cidrs"`
	Ocid                   string                            `tfschema:"ocid"`
	PrivateZoneOcid        string                            `tfschema:"private_zone_ocid"`
	ScanListenerPortTcp    int64                             `tfschema:"scan_listener_port_tcp"`
	ScanListenerPortTcpSsl int64                             `tfschema:"scan_listener_port_tcp_ssl"`
	SystemVersion          string                            `tfschema:"system_version"`
	TimeZone               string                            `tfschema:"time_zone"`
}

func (ExadbVmClusterResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.Location(),

		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.ExadbVMClusterName,
			ForceNew:     true,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		// Required
		"display_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ExadbVMClusterName,
		},

		"enabled_ecpu_count": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.EcpuCount,
		},

		"exascale_db_storage_vault_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: exascaledbstoragevaults.ValidateExascaleDbStorageVaultID,
		},

		"grid_image_ocid": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},

		"hostname": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"node_count": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(2, 10),
		},

		"shape": {
			Type:     pluginsdk.TypeString,
			Required: true,
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

		"total_ecpu_count": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.EcpuCount,
		},

		"vm_file_system_storage": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"total_size_in_gbs": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},
				},
			},
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
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validate.ClusterName,
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

		"domain": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
			ForceNew: true,
		},

		"license_model": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validate.ExadbLicenseModel,
		},

		"nsg_cidrs": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			ForceNew: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"destination_port_range": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"max": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
									Computed: true,
									ForceNew: true,
								},

								"min": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
									Computed: true,
									ForceNew: true,
								},
							},
						},
					},
					"source": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
						ForceNew: true,
					},
				},
			},
		},

		"private_zone_ocid": {
			Type:     schema.TypeString,
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
			Computed:     true,
			ValidateFunc: validate.ExadbSystemVersion,
		},

		"time_zone": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
			ForceNew: true,
		},

		"tags": commonschema.Tags(),

		"zones": commonschema.ZonesMultipleRequiredForceNew(),
	}
}

func (ExadbVmClusterResource) Attributes() map[string]*pluginsdk.Schema {
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

func (ExadbVmClusterResource) ModelObject() interface{} {
	return &ExadbVmClusterResource{}
}

func (ExadbVmClusterResource) ResourceType() string {
	return "azurerm_oracle_exa_db_vm_cluster"
}

func (r ExadbVmClusterResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 24 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient25.ExadbVMClusters
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model ExadbVmClusterResourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			id := exadbvmclusters.NewExadbVMClusterID(subscriptionId, model.ResourceGroupName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			param := exadbvmclusters.ExadbVMCluster{
				// Azure
				Name:     pointer.To(model.Name),
				Location: model.Location,
				Tags:     pointer.To(model.Tags),
				Zones:    pointer.To(model.Zones),
				Properties: &exadbvmclusters.ExadbVMClusterProperties{
					// Required
					DisplayName:              model.DisplayName,
					EnabledEcpuCount:         model.EnabledEcpuCount,
					ExascaleDbStorageVaultId: model.ExascaleDbStorageVaultId,
					GridImageOcid:            pointer.To(model.GridImageOcid),
					Hostname:                 model.Hostname,
					NodeCount:                model.NodeCount,
					Shape:                    model.Shape,
					SshPublicKeys:            model.SshPublicKeys,
					SubnetId:                 model.SubnetId,
					TotalEcpuCount:           model.TotalEcpuCount,
					VnetId:                   model.VnetId,
				},
			}

			if len(model.VmFileSystemStorage) > 0 {
				param.Properties.VMFileSystemStorage = exadbvmclusters.ExadbVMClusterStorageDetails{
					TotalSizeInGbs: model.VmFileSystemStorage[0].TotalSizeInGbs,
				}
			}

			if model.BackupSubnetCidr != "" {
				param.Properties.BackupSubnetCidr = pointer.To(model.BackupSubnetCidr)
			}
			if model.ClusterName != "" {
				param.Properties.ClusterName = pointer.To(model.ClusterName)
			}
			if len(model.DataCollectionOptions) > 0 {
				param.Properties.DataCollectionOptions = &exadbvmclusters.DataCollectionOptions{
					IsDiagnosticsEventsEnabled: pointer.To(model.DataCollectionOptions[0].IsDiagnosticsEventsEnabled),
					IsHealthMonitoringEnabled:  pointer.To(model.DataCollectionOptions[0].IsHealthMonitoringEnabled),
					IsIncidentLogsEnabled:      pointer.To(model.DataCollectionOptions[0].IsIncidentLogsEnabled),
				}
			}
			if model.Domain != "" {
				param.Properties.Domain = pointer.To(model.Domain)
			}
			if model.LicenseModel != "" {
				param.Properties.LicenseModel = pointer.To(exadbvmclusters.LicenseModel(model.LicenseModel))
			}
			if len(model.NsgCidrs) > 0 {
				param.Properties.NsgCidrs = pointer.To(ExpandNsgCidrs(model.NsgCidrs))
			}
			if model.PrivateZoneOcid != "" {
				param.Properties.PrivateZoneOcid = pointer.To(model.PrivateZoneOcid)
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

			if err := client.CreateOrUpdateThenPoll(ctx, id, param); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ExadbVmClusterResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient25.ExadbVMClusters
			id, err := exadbvmclusters.ParseExadbVMClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ExadbVmClusterResourceModel
			if err = metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding err: %+v", err)
			}

			_, err = client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			update := &exadbvmclusters.ExadbVMClusterUpdate{
				Properties: &exadbvmclusters.ExadbVMClusterUpdateProperties{},
			}

			if metadata.ResourceData.HasChange("node_count") {
				update.Properties.NodeCount = pointer.To(model.NodeCount)
			}

			if metadata.ResourceData.HasChange("tags") {
				update.Tags = pointer.To(model.Tags)
			}

			err = client.UpdateThenPoll(ctx, *id, *update)
			if err != nil {
				return fmt.Errorf("updating %s: %v", id, err)
			}

			return nil
		},
	}
}

func (ExadbVmClusterResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := exadbvmclusters.ParseExadbVMClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			client := metadata.Client.Oracle.OracleClient25.ExadbVMClusters
			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := ExadbVmClusterResourceModel{
				Name:              id.ExadbVmClusterName,
				ResourceGroupName: id.ResourceGroupName,
			}

			// Azure
			if model := resp.Model; model != nil {
				state.Location = model.Location
				state.Tags = pointer.From(model.Tags)
				state.Zones = pointer.From(model.Zones)

				if props := model.Properties; props != nil {
					state.DisplayName = props.DisplayName
					state.EnabledEcpuCount = props.EnabledEcpuCount
					state.ExascaleDbStorageVaultId = props.ExascaleDbStorageVaultId
					state.GridImageOcid = pointer.From(props.GridImageOcid)
					state.Hostname = props.Hostname
					state.NodeCount = props.NodeCount
					state.Shape = props.Shape
					state.SshPublicKeys = props.SshPublicKeys
					tmp := make([]string, 0)
					for _, key := range props.SshPublicKeys {
						if key != "" {
							tmp = append(tmp, key)
						}
					}
					state.SshPublicKeys = tmp
					state.SubnetId = props.SubnetId
					state.TotalEcpuCount = props.TotalEcpuCount
					state.VmFileSystemStorage = FlattenVMFileSystemStorage(props.VMFileSystemStorage)
					state.VnetId = props.VnetId
					state.Location = model.Location
					state.Tags = pointer.From(model.Tags)
					state.Zones = pointer.From(model.Zones)

					// Optional
					state.BackupSubnetCidr = pointer.From(props.BackupSubnetCidr)
					state.ClusterName = pointer.From(props.ClusterName)
					state.DataCollectionOptions = FlattenExadbDataCollectionOptionsInterface(metadata.ResourceData.Get("data_collection_options").([]interface{}))
					state.Domain = pointer.From(props.Domain)
					state.LicenseModel = string(pointer.From(props.LicenseModel))
					state.NsgCidrs = FlattenNsgCidrs(props.NsgCidrs)
					state.Ocid = pointer.From(props.Ocid)
					state.PrivateZoneOcid = pointer.From(props.PrivateZoneOcid)
					state.ScanListenerPortTcp = pointer.From(props.ScanListenerPortTcp)
					state.ScanListenerPortTcpSsl = pointer.From(props.ScanListenerPortTcpSsl)
					state.SystemVersion = metadata.ResourceData.Get("system_version").(string)
					state.TimeZone = pointer.From(props.TimeZone)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (ExadbVmClusterResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient25.ExadbVMClusters

			id, err := exadbvmclusters.ParseExadbVMClusterID(metadata.ResourceData.Id())
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

func (ExadbVmClusterResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return exadbvmclusters.ValidateExadbVMClusterID
}

func FlattenExadbDataCollectionOptions(dataCollectionOptions *exadbvmclusters.DataCollectionOptions) []ExadbDataCollectionOptionsModel {
	output := make([]ExadbDataCollectionOptionsModel, 0)
	if dataCollectionOptions != nil {
		return append(output, ExadbDataCollectionOptionsModel{
			IsDiagnosticsEventsEnabled: pointer.From(dataCollectionOptions.IsDiagnosticsEventsEnabled),
			IsHealthMonitoringEnabled:  pointer.From(dataCollectionOptions.IsHealthMonitoringEnabled),
			IsIncidentLogsEnabled:      pointer.From(dataCollectionOptions.IsIncidentLogsEnabled),
		})
	}
	return output
}

func FlattenExadbDataCollectionOptionsInterface(input []interface{}) []ExadbDataCollectionOptionsModel {
	output := make([]ExadbDataCollectionOptionsModel, 0)
	if len(input) == 0 || input[0] == nil {
		return output
	}
	if m, ok := input[0].(map[string]interface{}); ok {
		dataCollection := ExadbDataCollectionOptionsModel{
			IsDiagnosticsEventsEnabled: m["diagnostics_events_enabled"].(bool),
			IsHealthMonitoringEnabled:  m["health_monitoring_enabled"].(bool),
			IsIncidentLogsEnabled:      m["incident_logs_enabled"].(bool),
		}
		output = append(output, dataCollection)
	}
	return output
}

func ExpandNsgCidrs(input []NsgCidrModel) []exadbvmclusters.NsgCidr {
	output := make([]exadbvmclusters.NsgCidr, 0, len(input))

	for _, nsgCidr := range input {
		var portRangeValue = exadbvmclusters.PortRange{
			Max: nsgCidr.DestinationPortRange[0].Max,
			Min: nsgCidr.DestinationPortRange[0].Min,
		}
		output = append(output, exadbvmclusters.NsgCidr{
			DestinationPortRange: pointer.To(portRangeValue),
			Source:               nsgCidr.Source,
		})
	}
	return output
}
