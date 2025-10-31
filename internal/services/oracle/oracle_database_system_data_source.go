package oracle

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/dbsystems"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DatabaseSystemDataSource struct{}

type DatabaseSystemDataModel struct {
	Location          string       `tfschema:"location"`
	Name              string       `tfschema:"name"`
	ResourceGroupName string       `tfschema:"resource_group_name"`
	Zones             zones.Schema `tfschema:"zones"`

	// DatabaseSystemProperties
	ClusterName                  string                       `tfschema:"cluster_name"`
	ComputeCount                 int64                        `tfschema:"compute_count"`
	ComputeModel                 string                       `tfschema:"compute_model"`
	DatabaseEdition              string                       `tfschema:"database_edition"`
	DataStorageSizeInGbs         int64                        `tfschema:"data_storage_size_in_gbs"`
	DatabaseSystemOptions        []DatabaseSystemOptionsModel `tfschema:"database_system_options"`
	DatabaseVersion              string                       `tfschema:"database_version"`
	DiskRedundancy               string                       `tfschema:"disk_redundancy"`
	DisplayName                  string                       `tfschema:"display_name"`
	Domain                       string                       `tfschema:"domain"`
	GridImageOcid                string                       `tfschema:"grid_image_ocid"`
	Hostname                     string                       `tfschema:"hostname"`
	LicenseModel                 string                       `tfschema:"license_model"`
	LifecycleDetails             string                       `tfschema:"lifecycle_details"`
	LifecycleState               string                       `tfschema:"lifecycle_state"`
	ListenerPort                 int64                        `tfschema:"listener_port"`
	MemorySizeInGbs              int64                        `tfschema:"memory_size_in_gbs"`
	NetworkAnchorId              string                       `tfschema:"network_anchor_id"`
	NodeCount                    int64                        `tfschema:"node_count"`
	OciUrl                       string                       `tfschema:"oci_url"`
	Ocid                         string                       `tfschema:"ocid"`
	ResourceAnchorId             string                       `tfschema:"resource_anchor_id"`
	ScanDnsName                  string                       `tfschema:"scan_dns_name"`
	ScanIPs                      []string                     `tfschema:"scan_ips"`
	Shape                        string                       `tfschema:"shape"`
	Source                       string                       `tfschema:"source"`
	SshPublicKeys                []string                     `tfschema:"ssh_public_keys"`
	StorageVolumePerformanceMode string                       `tfschema:"storage_volume_performance_mode"`
	TimeZone                     string                       `tfschema:"time_zone"`
	Version                      string                       `tfschema:"version"`
}

type DatabaseSystemOptionsModel struct {
	StorageManagement string `tfschema:"storage_management"`
}

func (d DatabaseSystemDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.DatabaseSystemName,
		},
	}
}

func (d DatabaseSystemDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		// DatabaseSystemProperties
		"compute_count": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"compute_model": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"cluster_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"database_edition": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"data_storage_size_in_gbs": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"database_system_options": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"storage_management": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"database_version": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"disk_redundancy": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"display_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"domain": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"grid_image_ocid": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"hostname": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"license_model": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"lifecycle_details": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"lifecycle_state": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"listener_port": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"memory_size_in_gbs": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"node_count": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"network_anchor_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"oci_url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"ocid": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"resource_anchor_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"scan_dns_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"scan_ips": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"shape": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"ssh_public_keys": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"source": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"storage_volume_performance_mode": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"time_zone": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"version": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"zones": commonschema.ZonesMultipleComputed(),
	}
}

func (d DatabaseSystemDataSource) ModelObject() interface{} {
	return &DatabaseSystemDataModel{}
}

func (d DatabaseSystemDataSource) ResourceType() string {
	return "azurerm_oracle_database_system"
}

func (d DatabaseSystemDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return dbsystems.ValidateDbSystemID
}

func (d DatabaseSystemDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.DbSystems
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state DatabaseSystemDataModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := dbsystems.NewDbSystemID(subscriptionId, state.ResourceGroupName, state.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Zones = pointer.From(model.Zones)
				if props := model.Properties; props != nil {
					state.DatabaseEdition = string(props.DatabaseEdition)
					state.DatabaseVersion = props.DbVersion

					databaseSystemProps := props.DbSystemBaseProperties()

					state.ComputeCount = pointer.From(databaseSystemProps.ComputeCount)
					state.ComputeModel = pointer.FromEnum(databaseSystemProps.ComputeModel)
					state.ClusterName = pointer.From(databaseSystemProps.ClusterName)
					state.DataStorageSizeInGbs = pointer.From(databaseSystemProps.DataStorageSizeInGbs)
					state.DatabaseSystemOptions = FlattenDbSystemOptions(databaseSystemProps.DbSystemOptions)
					state.DiskRedundancy = string(pointer.From(props.DiskRedundancy))
					state.DisplayName = pointer.From(databaseSystemProps.DisplayName)
					state.Domain = pointer.From(databaseSystemProps.Domain)
					state.GridImageOcid = pointer.From(databaseSystemProps.GridImageOcid)
					state.Hostname = databaseSystemProps.Hostname
					state.LicenseModel = string(pointer.From(databaseSystemProps.LicenseModel))
					state.LifecycleDetails = pointer.From(databaseSystemProps.LifecycleDetails)
					state.LifecycleState = string(*databaseSystemProps.LifecycleState)
					state.ListenerPort = pointer.From(databaseSystemProps.ListenerPort)
					state.MemorySizeInGbs = pointer.From(databaseSystemProps.MemorySizeInGbs)
					state.NetworkAnchorId = databaseSystemProps.NetworkAnchorId
					state.NodeCount = pointer.From(databaseSystemProps.NodeCount)
					state.OciUrl = pointer.From(databaseSystemProps.OciURL)
					state.Ocid = pointer.From(databaseSystemProps.Ocid)
					state.ResourceAnchorId = databaseSystemProps.ResourceAnchorId
					state.ScanDnsName = pointer.From(databaseSystemProps.ScanDnsName)
					state.ScanIPs = pointer.From(databaseSystemProps.ScanIPs)
					state.Shape = databaseSystemProps.Shape
					state.Source = string(databaseSystemProps.Source)
					state.SshPublicKeys = databaseSystemProps.SshPublicKeys
					state.StorageVolumePerformanceMode = string(pointer.From(databaseSystemProps.StorageVolumePerformanceMode))
					state.TimeZone = pointer.From(databaseSystemProps.TimeZone)
					state.Version = pointer.From(databaseSystemProps.Version)
				}
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
