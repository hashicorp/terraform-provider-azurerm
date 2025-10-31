// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oracle

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/exadbvmclusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/exascaledbstoragevaults"

	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.Resource = ExascaleDatabaseVirtualMachineClusterResource{}

type ExascaleDatabaseVirtualMachineClusterResource struct{}

type ExascaleDatabaseVirtualMachineClusterResourceModel struct {
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Location          string            `tfschema:"location"`
	Tags              map[string]string `tfschema:"tags"`
	Zones             zones.Schema      `tfschema:"zones"`

	DisplayName                     string                                              `tfschema:"display_name"`
	EnabledEcpuCount                int64                                               `tfschema:"enabled_ecpu_count"`
	ExascaleDbStorageVaultId        string                                              `tfschema:"exascale_database_storage_vault_id"`
	GridImageOcid                   string                                              `tfschema:"grid_image_ocid"`
	Hostname                        string                                              `tfschema:"hostname"`
	NodeCount                       int64                                               `tfschema:"node_count"`
	Shape                           string                                              `tfschema:"shape"`
	SshPublicKeys                   []string                                            `tfschema:"ssh_public_keys"`
	SubnetId                        string                                              `tfschema:"subnet_id"`
	TotalEcpuCount                  int64                                               `tfschema:"total_ecpu_count"`
	VirtualMachineFileSystemStorage []ExascaleDatabaseVirtualMachineClusterStorageModel `tfschema:"virtual_machine_file_system_storage"`
	VnetId                          string                                              `tfschema:"virtual_network_id"`

	BackupSubnetCidr                         string                                      `tfschema:"backup_subnet_cidr"`
	ClusterName                              string                                      `tfschema:"cluster_name"`
	DataCollectionOption                     []ExascaleDatabaseDataCollectionOptionModel `tfschema:"data_collection_option"`
	Domain                                   string                                      `tfschema:"domain"`
	LicenseModel                             string                                      `tfschema:"license_model"`
	NetworkSecurityGroupCidrs                []NetworkSecurityGroupCidrModel             `tfschema:"network_security_group_cidrs"`
	Ocid                                     string                                      `tfschema:"ocid"`
	PrivateZoneOcid                          string                                      `tfschema:"private_zone_ocid"`
	SingleClientAccessNameListenerPortTcp    int64                                       `tfschema:"single_client_access_name_listener_port_tcp"`
	SingleClientAccessNameListenerPortTcpSsl int64                                       `tfschema:"single_client_access_name_listener_port_tcp_ssl"`
	SystemVersion                            string                                      `tfschema:"system_version"`
	TimeZone                                 string                                      `tfschema:"time_zone"`
	ZoneOcid                                 string                                      `tfschema:"zone_ocid"`
}

func (ExascaleDatabaseVirtualMachineClusterResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, 255),
				validate.ExascaleDatabaseVirtualMachineClusterName,
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"display_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, 255),
				validate.ExascaleDatabaseVirtualMachineClusterName,
			),
		},

		"enabled_ecpu_count": {
			Type:     pluginsdk.TypeInt,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.IntBetween(8, 200),
				validation.IntDivisibleBy(4),
			),
		},

		"exascale_database_storage_vault_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: exascaledbstoragevaults.ValidateExascaleDbStorageVaultID,
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

		"total_ecpu_count": {
			Type:     pluginsdk.TypeInt,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.IntBetween(8, 200),
				validation.IntDivisibleBy(4),
			),
		},

		"virtual_machine_file_system_storage": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"total_size_in_gb": {
						Type:     pluginsdk.TypeInt,
						Required: true,
						ForceNew: true,
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

		"backup_subnet_cidr": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsCIDR,
		},

		"cluster_name": {
			Type: pluginsdk.TypeString,
			// O+C if not specified, this gets set to the name of the virtual machine
			Optional: true,
			Computed: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9-]{0,10}$`),
				"The Cluster name must begin with an alphabetic character, be no longer than 11 characters, and may contain alphabets, numbers, and hyphens (-).",
			),
		},

		"data_collection_option": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"diagnostics_events_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						ForceNew: true,
						Default:  false,
					},

					"health_monitoring_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						ForceNew: true,
						Default:  false,
					},

					"incident_logs_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						ForceNew: true,
						Default:  false,
					},
				},
			},
		},

		"domain": {
			Type: pluginsdk.TypeString,
			// O+C if not specified, this gets set to the name of the virtual machine
			Optional: true,
			Computed: true,
			ForceNew: true,
		},

		"grid_image_ocid": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
		},

		"license_model": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(exadbvmclusters.LicenseModelBringYourOwnLicense),
				string(exadbvmclusters.LicenseModelLicenseIncluded),
			}, false),
		},

		"network_security_group_cidrs": {
			Type: pluginsdk.TypeList,
			// O+C if not specified, this gets set to the name of the virtual machine
			Optional: true,
			Computed: true,
			ForceNew: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"destination_port_range": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MinItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"max": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.IntBetween(1, 65535),
								},

								"min": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.IntBetween(1, 65535),
								},
							},
						},
					},
					"source": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringLenBetween(1, 128),
					},
				},
			},
		},

		"private_zone_ocid": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
		},

		"single_client_access_name_listener_port_tcp": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      1521,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(1024, 8999),
		},

		"single_client_access_name_listener_port_tcp_ssl": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      2484,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(1024, 8999),
		},

		"system_version": {
			Type: pluginsdk.TypeString,
			// O+C if not specified, the default value will be provided by API
			Optional: true,
			Computed: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(regexp.MustCompile(`^(19|22|23|24|25)\.[0-9]+(\.[0-9]+)*|[0-9]+(\.[0-9]+)$`),
				"The system_version must match following patterns: (19|22|23|24|25).x(.x(.x)).",
			),
		},

		"time_zone": {
			Type: pluginsdk.TypeString,
			// O+C if not specified, the default value will be provided by API
			Optional: true,
			Computed: true,
			ForceNew: true,
		},

		"tags": commonschema.Tags(),

		"zones": commonschema.ZonesMultipleRequiredForceNew(),
	}
}

func (ExascaleDatabaseVirtualMachineClusterResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"ocid": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"zone_ocid": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (ExascaleDatabaseVirtualMachineClusterResource) ModelObject() interface{} {
	return &ExascaleDatabaseVirtualMachineClusterResource{}
}

func (ExascaleDatabaseVirtualMachineClusterResource) ResourceType() string {
	return "azurerm_oracle_exascale_database_virtual_machine_cluster"
}

func (r ExascaleDatabaseVirtualMachineClusterResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.ExadbVMClusters
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model ExascaleDatabaseVirtualMachineClusterResourceModel
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
				Name:     pointer.To(model.Name),
				Location: model.Location,
				Tags:     pointer.To(model.Tags),
				Zones:    pointer.To(model.Zones),
				Properties: &exadbvmclusters.ExadbVMClusterProperties{
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

			if len(model.VirtualMachineFileSystemStorage) > 0 {
				param.Properties.VMFileSystemStorage = exadbvmclusters.ExadbVMClusterStorageDetails{
					TotalSizeInGbs: model.VirtualMachineFileSystemStorage[0].TotalSizeInGb,
				}
			}

			if model.BackupSubnetCidr != "" {
				param.Properties.BackupSubnetCidr = pointer.To(model.BackupSubnetCidr)
			}
			if model.ClusterName != "" {
				param.Properties.ClusterName = pointer.To(model.ClusterName)
			}
			if len(model.DataCollectionOption) > 0 {
				param.Properties.DataCollectionOptions = &exadbvmclusters.DataCollectionOptions{
					IsDiagnosticsEventsEnabled: pointer.To(model.DataCollectionOption[0].IsDiagnosticsEventsEnabled),
					IsHealthMonitoringEnabled:  pointer.To(model.DataCollectionOption[0].IsHealthMonitoringEnabled),
					IsIncidentLogsEnabled:      pointer.To(model.DataCollectionOption[0].IsIncidentLogsEnabled),
				}
			}
			if model.Domain != "" {
				param.Properties.Domain = pointer.To(model.Domain)
			}
			if model.LicenseModel != "" {
				param.Properties.LicenseModel = pointer.To(exadbvmclusters.LicenseModel(model.LicenseModel))
			}
			if len(model.NetworkSecurityGroupCidrs) > 0 {
				param.Properties.NsgCidrs = pointer.To(expandNsgCidrs(model.NetworkSecurityGroupCidrs))
			}
			if model.PrivateZoneOcid != "" {
				param.Properties.PrivateZoneOcid = pointer.To(model.PrivateZoneOcid)
			}
			if model.SingleClientAccessNameListenerPortTcp >= 1024 && model.SingleClientAccessNameListenerPortTcp <= 8999 {
				param.Properties.ScanListenerPortTcp = pointer.To(model.SingleClientAccessNameListenerPortTcp)
			}
			if model.SingleClientAccessNameListenerPortTcpSsl >= 1024 && model.SingleClientAccessNameListenerPortTcpSsl <= 8999 {
				param.Properties.ScanListenerPortTcpSsl = pointer.To(model.SingleClientAccessNameListenerPortTcpSsl)
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

func (r ExascaleDatabaseVirtualMachineClusterResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.ExadbVMClusters
			id, err := exadbvmclusters.ParseExadbVMClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ExascaleDatabaseVirtualMachineClusterResourceModel
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

func (ExascaleDatabaseVirtualMachineClusterResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := exadbvmclusters.ParseExadbVMClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			client := metadata.Client.Oracle.OracleClient.ExadbVMClusters
			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := ExascaleDatabaseVirtualMachineClusterResourceModel{
				Name:              id.ExadbVmClusterName,
				ResourceGroupName: id.ResourceGroupName,
			}

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
					state.VirtualMachineFileSystemStorage = FlattenVMFileSystemStorage(props.VMFileSystemStorage)
					state.VnetId = props.VnetId
					state.Location = model.Location
					state.Tags = pointer.From(model.Tags)
					state.Zones = pointer.From(model.Zones)

					state.BackupSubnetCidr = pointer.From(props.BackupSubnetCidr)
					state.ClusterName = pointer.From(props.ClusterName)
					state.DataCollectionOption = flattenExadbDataCollectionOptionInterface(metadata.ResourceData.Get("data_collection_option").([]interface{}))
					state.Domain = pointer.From(props.Domain)
					state.LicenseModel = string(pointer.From(props.LicenseModel))
					state.NetworkSecurityGroupCidrs = FlattenNetworkSecurityGroupCidr(props.NsgCidrs)
					state.Ocid = pointer.From(props.Ocid)
					state.PrivateZoneOcid = pointer.From(props.PrivateZoneOcid)
					state.SingleClientAccessNameListenerPortTcp = pointer.From(props.ScanListenerPortTcp)
					state.SingleClientAccessNameListenerPortTcpSsl = pointer.From(props.ScanListenerPortTcpSsl)
					state.SystemVersion = metadata.ResourceData.Get("system_version").(string)
					state.TimeZone = pointer.From(props.TimeZone)
					state.ZoneOcid = pointer.From(props.ZoneOcid)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (ExascaleDatabaseVirtualMachineClusterResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.ExadbVMClusters

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

func (ExascaleDatabaseVirtualMachineClusterResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return exadbvmclusters.ValidateExadbVMClusterID
}

func flattenExadbDataCollectionOptionInterface(input []interface{}) []ExascaleDatabaseDataCollectionOptionModel {
	output := make([]ExascaleDatabaseDataCollectionOptionModel, 0)
	if len(input) == 0 || input[0] == nil {
		return output
	}
	if m, ok := input[0].(map[string]interface{}); ok {
		dataCollection := ExascaleDatabaseDataCollectionOptionModel{
			IsDiagnosticsEventsEnabled: m["diagnostics_events_enabled"].(bool),
			IsHealthMonitoringEnabled:  m["health_monitoring_enabled"].(bool),
			IsIncidentLogsEnabled:      m["incident_logs_enabled"].(bool),
		}
		output = append(output, dataCollection)
	}
	return output
}

func expandNsgCidrs(input []NetworkSecurityGroupCidrModel) []exadbvmclusters.NsgCidr {
	output := make([]exadbvmclusters.NsgCidr, 0, len(input))

	for _, nsgCidr := range input {
		portRangeValue := exadbvmclusters.PortRange{
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
