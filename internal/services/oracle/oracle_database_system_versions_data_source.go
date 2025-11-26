package oracle

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/dbversions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type DatabaseVersionsDataSource struct{}

type DatabaseVersionsModel struct {
	Location                            string                     `tfschema:"location"`
	DatabaseSystemShape                 string                     `tfschema:"database_system_shape"`
	DatabaseSoftwareImageSupportEnabled *bool                      `tfschema:"database_software_image_support_enabled"`
	ShapeFamily                         string                     `tfschema:"shape_family"`
	StorageManagement                   string                     `tfschema:"storage_management"`
	UpgradeSupportEnabled               *bool                      `tfschema:"upgrade_support_enabled"`
	Versions                            []DatabaseVersionItemModel `tfschema:"versions"`
}

type DatabaseVersionItemModel struct {
	Name                         string `tfschema:"name"`
	LatestForMajorVersionEnabled bool   `tfschema:"latest_for_major_version_enabled"`
	SupportsPdb                  bool   `tfschema:"supports_pdb"`
	Version                      string `tfschema:"version"`
}

func (d DatabaseVersionsDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.Location(),

		// Optional filters

		"database_software_image_support_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},
		"database_system_shape": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ValidateFunc:  validation.StringInSlice(dbversions.PossibleValuesForBaseDbSystemShapes(), false),
			ConflictsWith: []string{"shape_family"},
		},
		"shape_family": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ValidateFunc:  validation.StringInSlice(dbversions.PossibleValuesForShapeFamilyType(), false),
			ConflictsWith: []string{"database_system_shape"},
		},
		"storage_management": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(dbversions.PossibleValuesForStorageManagementType(), false),
		},
		"upgrade_support_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},
	}
}

func (d DatabaseVersionsDataSource) Attributes() map[string]*pluginsdk.Schema {
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
					"latest_for_major_version_enabled": {
						Type:     pluginsdk.TypeBool,
						Computed: true,
					},
					"supports_pdb": {
						Type:     pluginsdk.TypeBool,
						Computed: true,
					},
					"version": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},
	}
}

func (d DatabaseVersionsDataSource) ModelObject() interface{} {
	return &DatabaseVersionsModel{}
}

func (d DatabaseVersionsDataSource) ResourceType() string {
	return "azurerm_oracle_database_system_versions"
}

func (d DatabaseVersionsDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return dbversions.ValidateDbSystemDbVersionID
}

func (d DatabaseVersionsDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.DbVersions
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state DatabaseVersionsModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := dbversions.NewLocationID(subscriptionId, location.Normalize(state.Location))

			options := dbversions.DefaultListByLocationOperationOptions()

			options.IsDatabaseSoftwareImageSupported = state.DatabaseSoftwareImageSupportEnabled
			options.IsUpgradeSupported = state.UpgradeSupportEnabled

			if state.DatabaseSystemShape != "" {
				options.DbSystemShape = pointer.ToEnum[dbversions.BaseDbSystemShapes](state.DatabaseSystemShape)
			}
			if state.ShapeFamily != "" {
				options.ShapeFamily = pointer.ToEnum[dbversions.ShapeFamilyType](state.ShapeFamily)
			}
			if state.StorageManagement != "" {
				options.StorageManagement = pointer.ToEnum[dbversions.StorageManagementType](state.StorageManagement)
			}

			resp, err := client.ListByLocation(ctx, id, options)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state.Versions = make([]DatabaseVersionItemModel, 0)

			if model := resp.Model; model != nil {
				for _, element := range *model {
					if props := element.Properties; props != nil {
						item := DatabaseVersionItemModel{
							Name:                         pointer.From(element.Name),
							Version:                      props.Version,
							LatestForMajorVersionEnabled: pointer.From(props.IsLatestForMajorVersion),
							SupportsPdb:                  pointer.From(props.SupportsPdb),
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
