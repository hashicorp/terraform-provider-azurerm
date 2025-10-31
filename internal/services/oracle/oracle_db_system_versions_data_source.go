package oracle

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/dbversions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type DbVersionsDataSource struct{}

type DbVersionsModel struct {
	Location                       string               `tfschema:"location"`
	DatabaseSystemShape            string               `tfschema:"database_system_shape"`
	DatabaseSoftwareImageSupported *bool                `tfschema:"database_software_image_supported"`
	UpgradeSupported               *bool                `tfschema:"upgrade_supported"`
	ShapeFamily                    string               `tfschema:"shape_family"`
	StorageManagement              string               `tfschema:"storage_management"`
	Versions                       []DbVersionItemModel `tfschema:"versions"`
}

type DbVersionItemModel struct {
	Name                  string `tfschema:"name"`
	Version               string `tfschema:"version"`
	LatestForMajorVersion bool   `tfschema:"latest_for_major_version"`
	SupportsPdb           bool   `tfschema:"supports_pdb"`
}

func (d DbVersionsDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.Location(),

		// Optional filters

		"database_system_shape": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(dbversions.PossibleValuesForBaseDbSystemShapes(), false),
		},
		"database_software_image_supported": {
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Description: "Whether to filter versions that support database software images.",
		},
		"upgrade_supported": {
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Description: "Whether to filter versions that support upgrades.",
		},
		"shape_family": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(dbversions.PossibleValuesForShapeFamilyType(), false),
			Description:  "The shape family of the DB System to filter versions by.",
		},
		"storage_management": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(dbversions.PossibleValuesForStorageManagementType(), false),
			Description:  "The storage management type to filter versions by.",
		},
	}
}

func (d DbVersionsDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"versions": {
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Description: "A list of available Oracle Database versions and their properties.",
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"version": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"latest_for_major_version": {
						Type:     pluginsdk.TypeBool,
						Computed: true,
					},
					"supports_pdb": {
						Type:     pluginsdk.TypeBool,
						Computed: true,
					},
				},
			},
		},
	}
}

func (d DbVersionsDataSource) ModelObject() interface{} {
	return &DbVersionsModel{}
}

func (d DbVersionsDataSource) ResourceType() string {
	return "azurerm_oracle_db_system_versions"
}

func (d DbVersionsDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return dbversions.ValidateDbSystemDbVersionID
}

func (d DbVersionsDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.DbVersions
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state DbVersionsModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := dbversions.NewLocationID(subscriptionId,
				state.Location)

			options := dbversions.DefaultListByLocationOperationOptions()

			if state.ShapeFamily != "" || state.DatabaseSystemShape != "" {
				log.Printf("[WARN] Db Versions data source: You can not use both shape_family and database_system_shape together for filtering versions.")
			}

			options.IsDatabaseSoftwareImageSupported = state.DatabaseSoftwareImageSupported
			options.IsUpgradeSupported = state.UpgradeSupported

			if state.DatabaseSystemShape != "" {
				shape := dbversions.BaseDbSystemShapes(state.DatabaseSystemShape)
				options.DbSystemShape = &shape
			}
			if state.ShapeFamily != "" {
				family := dbversions.ShapeFamilyType(state.ShapeFamily)
				options.ShapeFamily = &family
			}
			if state.StorageManagement != "" {
				storage := dbversions.StorageManagementType(state.StorageManagement)
				options.StorageManagement = &storage
			}

			resp, err := client.ListByLocation(ctx, id, options)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state.Versions = make([]DbVersionItemModel, 0)

			if model := resp.Model; model != nil {
				for _, element := range *model {
					if props := element.Properties; props != nil {
						item := DbVersionItemModel{
							Name:                  pointer.From(element.Name),
							Version:               props.Version,
							LatestForMajorVersion: pointer.From(props.IsLatestForMajorVersion),
							SupportsPdb:           pointer.From(props.SupportsPdb),
						}
						state.Versions = append(state.Versions, item)
					}
				}
			}
			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
