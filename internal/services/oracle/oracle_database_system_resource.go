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
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/dbversions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/networkanchors"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/resourceanchors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.ResourceWithUpdate = DatabaseSystemResource{}

type DatabaseSystemResource struct{}

type DatabaseSystemResourceModel struct {
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Location          string            `tfschema:"location"`
	Tags              map[string]string `tfschema:"tags"`
	Zones             zones.Schema      `tfschema:"zones"`

	// Required
	AdminPassword    string   `tfschema:"admin_password"`
	ComputeCount     int64    `tfschema:"compute_count"`
	ComputeModel     string   `tfschema:"compute_model"`
	DatabaseEdition  string   `tfschema:"database_edition"`
	DatabaseVersion  string   `tfschema:"database_version"`
	Hostname         string   `tfschema:"hostname"`
	LicenseModel     string   `tfschema:"license_model"`
	NetworkAnchorId  string   `tfschema:"network_anchor_id"`
	ResourceAnchorId string   `tfschema:"resource_anchor_id"`
	Shape            string   `tfschema:"shape"`
	Source           string   `tfschema:"source"`
	SshPublicKeys    []string `tfschema:"ssh_public_keys"`

	// Optional
	ClusterName                  string                       `tfschema:"cluster_name"`
	DatabaseSystemOptions        []DatabaseSystemOptionsModel `tfschema:"database_system_options"`
	DiskRedundancy               string                       `tfschema:"disk_redundancy"`
	DisplayName                  string                       `tfschema:"display_name"`
	Domain                       string                       `tfschema:"domain"`
	InitialDataStorageSizeInGb   int64                        `tfschema:"initial_data_storage_size_in_gb"`
	NodeCount                    int64                        `tfschema:"node_count"`
	PluggableDatabaseName        string                       `tfschema:"pluggable_database_name`
	StorageVolumePerformanceMode string                       `tfschema:"storage_volume_performance_mode"`
	TimeZone                     string                       `tfschema:"time_zone"`
}

func (DatabaseSystemResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DatabaseSystemName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		// Required
		"admin_password": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			Sensitive:    true,
			ForceNew:     true,
			ValidateFunc: validate.DatabaseSystemPassword,
		},

		"compute_count": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.FloatBetween(1.0, 512.0),
		},

		"compute_model": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(dbsystems.PossibleValuesForComputeModel(), false),
		},

		"database_edition": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(dbsystems.PossibleValuesForDbSystemDatabaseEditionType(), false),
		},

		"database_version": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"hostname": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DatabaseSystemName,
		},

		"license_model": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(dbsystems.PossibleValuesForLicenseModel(), false),
		},

		"network_anchor_id": commonschema.ResourceIDReferenceRequiredForceNew(&networkanchors.NetworkAnchorId{}),

		"resource_anchor_id": commonschema.ResourceIDReferenceRequiredForceNew(&resourceanchors.ResourceAnchorId{}),

		"shape": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(dbversions.PossibleValuesForBaseDbSystemShapes(), false),
		},

		"source": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(dbsystems.PossibleValuesForSource(), false),
		},

		"ssh_public_keys": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"zones": commonschema.ZonesMultipleRequiredForceNew(),

		// Optional
		"database_system_options": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"storage_management": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(dbsystems.PossibleValuesForStorageManagementType(), false),
					},
				},
			},
		},

		"initial_data_storage_size_in_gb": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntAtLeast(2),
		},

		"pluggable_database_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validate.PluggableDatabaseName,
		},

		"tags": commonschema.Tags(),

		// Computed
		"cluster_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validate.ClusterName,
		},

		"disk_redundancy": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(dbsystems.PossibleValuesForDiskRedundancyType(), false),
		},

		"display_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validate.DatabaseSystemName,
		},

		"domain": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
			ForceNew: true,
		},

		"node_count": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},

		"storage_volume_performance_mode": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(dbsystems.PossibleValuesForStorageVolumePerformanceMode(), false),
		},

		"time_zone": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
			ForceNew: true,
		},
	}
}

func (DatabaseSystemResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (DatabaseSystemResource) ModelObject() interface{} {
	return &DatabaseSystemResource{}
}

func (DatabaseSystemResource) ResourceType() string {
	return "azurerm_oracle_database_system"
}

func (r DatabaseSystemResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 24 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.DbSystems
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model DatabaseSystemResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := dbsystems.NewDbSystemID(subscriptionId, model.ResourceGroupName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			param := dbsystems.DbSystem{
				// Azure
				Name:     pointer.To(model.Name),
				Location: location.Normalize(model.Location),
				Tags:     pointer.To(model.Tags),
				Zones:    pointer.To(model.Zones),
				Properties: &dbsystems.DbSystemProperties{
					// Required
					AdminPassword:    pointer.To(model.AdminPassword),
					ComputeCount:     pointer.To(model.ComputeCount),
					ComputeModel:     pointer.ToEnum[dbsystems.ComputeModel](model.ComputeModel),
					DatabaseEdition:  dbsystems.DbSystemDatabaseEditionType(model.DatabaseEdition),
					DbVersion:        model.DatabaseVersion,
					Hostname:         model.Hostname,
					NetworkAnchorId:  model.NetworkAnchorId,
					ResourceAnchorId: model.ResourceAnchorId,
					Source:           dbsystems.DbSystemSourceType(model.Source),
					Shape:            model.Shape,
					LicenseModel:     pointer.ToEnum[dbsystems.LicenseModel](model.LicenseModel),
					SshPublicKeys:    model.SshPublicKeys,
				},
			}

			// Optional
			if model.ClusterName != "" {
				param.Properties.ClusterName = pointer.To(model.ClusterName)
			}
			if len(model.DatabaseSystemOptions) > 0 {
				param.Properties.DbSystemOptions = &dbsystems.DbSystemOptions{
					StorageManagement: pointer.ToEnum[dbsystems.StorageManagementType](model.DatabaseSystemOptions[0].StorageManagement),
				}
			}
			if model.DiskRedundancy != "" {
				param.Properties.DiskRedundancy = pointer.ToEnum[dbsystems.DiskRedundancyType](model.DiskRedundancy)
			}
			if model.DisplayName != "" {
				param.Properties.DisplayName = pointer.To(model.DisplayName)
			}
			if model.Domain != "" {
				param.Properties.Domain = pointer.To(model.Domain)
			}
			if model.InitialDataStorageSizeInGb != 0 {
				param.Properties.InitialDataStorageSizeInGb = pointer.To(model.InitialDataStorageSizeInGb)
			}
			if model.NodeCount != 0 {
				param.Properties.NodeCount = pointer.To(model.NodeCount)
			}
			if model.PluggableDatabaseName != "" {
				param.Properties.PdbName = pointer.To(model.PluggableDatabaseName)
			}
			if model.StorageVolumePerformanceMode != "" {
				param.Properties.StorageVolumePerformanceMode = pointer.ToEnum[dbsystems.StorageVolumePerformanceMode](model.StorageVolumePerformanceMode)
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

func (r DatabaseSystemResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.DbSystems
			id, err := dbsystems.ParseDbSystemID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model DatabaseSystemResourceModel
			if err = metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			_, err = client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			update := dbsystems.DbSystemUpdate{}
			if metadata.ResourceData.HasChange("tags") {
				update.Tags = pointer.To(model.Tags)
			}

			if err := client.UpdateThenPoll(ctx, *id, update); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (DatabaseSystemResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := dbsystems.ParseDbSystemID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			client := metadata.Client.Oracle.OracleClient.DbSystems
			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := DatabaseSystemResourceModel{
				Name:              id.DbSystemName,
				ResourceGroupName: id.ResourceGroupName,
			}

			// Azure
			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Tags = pointer.From(model.Tags)
				state.Zones = pointer.From(model.Zones)

				if props := model.Properties; props != nil {
					state.AdminPassword = metadata.ResourceData.Get("admin_password").(string)
					state.ComputeCount = pointer.From(props.ComputeCount)
					state.ComputeModel = string(pointer.FromEnum(props.ComputeModel))
					state.DatabaseEdition = string(props.DatabaseEdition)
					state.DatabaseVersion = props.DbVersion
					state.Hostname = props.Hostname
					state.NetworkAnchorId = props.NetworkAnchorId
					state.ResourceAnchorId = metadata.ResourceData.Get("resource_anchor_id").(string)
					state.Source = string(props.Source)
					state.Shape = props.Shape
					state.Zones = pointer.From(model.Zones)

					state.SshPublicKeys = props.SshPublicKeys
					tmp := make([]string, 0)
					for _, key := range props.SshPublicKeys {
						if key != "" {
							tmp = append(tmp, key)
						}
					}
					state.SshPublicKeys = tmp

					// Optional
					state.ClusterName = pointer.From(props.ClusterName)
					state.DatabaseSystemOptions = FlattenDbSystemOptions(props.DbSystemOptions)
					state.DiskRedundancy = string(pointer.From(props.DiskRedundancy))
					state.DisplayName = pointer.From(props.DisplayName)
					state.Domain = pointer.From(props.Domain)
					state.InitialDataStorageSizeInGb = int64(metadata.ResourceData.Get("initial_data_storage_size_in_gb").(int))
					state.LicenseModel = string(pointer.From(props.LicenseModel))
					state.NodeCount = pointer.From(props.NodeCount)
					state.PluggableDatabaseName = pointer.From(props.PdbName)
					state.StorageVolumePerformanceMode = string(pointer.From(props.StorageVolumePerformanceMode))
					state.TimeZone = pointer.From(props.TimeZone)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (DatabaseSystemResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.DbSystems

			id, err := dbsystems.ParseDbSystemID(metadata.ResourceData.Id())
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

func (DatabaseSystemResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return dbsystems.ValidateDbSystemID
}

func FlattenDbSystemOptions(input *dbsystems.DbSystemOptions) []DatabaseSystemOptionsModel {
	output := make([]DatabaseSystemOptionsModel, 0)
	if input != nil {
		return append(output, DatabaseSystemOptionsModel{
			StorageManagement: string(pointer.From(input.StorageManagement)),
		})
	}
	return output
}
