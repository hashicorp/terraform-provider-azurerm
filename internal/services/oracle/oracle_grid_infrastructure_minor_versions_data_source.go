// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package oracle

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/giminorversions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type GridInfrastructureMinorVersionsDataSource struct{}

type GridInfrastructureMinorVersionsModel struct {
	Location                        string                                `tfschema:"location"`
	GridInfrastructureVersion       string                                `tfschema:"grid_infrastructure_version"`
	ShapeFamily                     string                                `tfschema:"shape_family"`
	Zone                            string                                `tfschema:"zone"`
	GridInfrastructureMinorVersions []GridInfrastructureMinorVersionModel `tfschema:"versions"`
}

type GridInfrastructureMinorVersionModel struct {
	Id            string `tfschema:"id"`
	Name          string `tfschema:"name"`
	Version       string `tfschema:"version"`
	GridImageOcid string `tfschema:"grid_image_ocid"`
}

func (d GridInfrastructureMinorVersionsDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.Location(),

		"grid_infrastructure_version": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  "The GI (Grid Infrastructure) version to list minor versions for.",
		},

		"shape_family": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(giminorversions.PossibleValuesForShapeFamily(), false),
			Description:  "Filter the minor versions by shape family. Possible values are 'EXADATA' and 'EXADB_XS'.",
		},

		"zone": {
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Description: "Filter the minor versions by zone.",
		},
	}
}

func (d GridInfrastructureMinorVersionsDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"versions": {
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Description: "A list of available Oracle Grid Infrastructure minor versions and their properties.",
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"version": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"grid_image_ocid": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},
	}
}

func (d GridInfrastructureMinorVersionsDataSource) ModelObject() interface{} {
	return &GridInfrastructureMinorVersionsModel{}
}

func (d GridInfrastructureMinorVersionsDataSource) ResourceType() string {
	return "azurerm_oracle_grid_infrastructure_minor_versions"
}

func (d GridInfrastructureMinorVersionsDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return giminorversions.ValidateGiVersionID
}

func (d GridInfrastructureMinorVersionsDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.GiMinorVersions
			subscriptionId := metadata.Client.Account.SubscriptionId

			state := GridInfrastructureMinorVersionsModel{
				GridInfrastructureMinorVersions: make([]GridInfrastructureMinorVersionModel, 0),
			}
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := giminorversions.NewGiVersionID(subscriptionId, state.Location, state.GridInfrastructureVersion)

			options := giminorversions.DefaultListByParentOperationOptions()
			if state.ShapeFamily != "" {
				options.ShapeFamily = pointer.ToEnum[giminorversions.ShapeFamily](state.ShapeFamily)
			}
			if state.Zone != "" {
				options.Zone = &state.Zone
			}

			if state.ShapeFamily == "" || state.Zone == "" {
				log.Printf("[WARN] GI Minor Versions data source: ShapeFamily or Zone parameter is empty. This may result in unfiltered results from the API. Consider specifying both ShapeFamily and Zone for more precise version filtering.")
			}

			resp, err := client.ListByParent(ctx, id, options)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if model := resp.Model; model != nil {
				for _, element := range *model {
					if props := element.Properties; props != nil {
						state.GridInfrastructureMinorVersions = append(state.GridInfrastructureMinorVersions, GridInfrastructureMinorVersionModel{
							Id:            pointer.From(element.Id),
							Name:          pointer.From(element.Name),
							Version:       props.Version,
							GridImageOcid: pointer.From(props.GridImageOcid),
						})
					}
				}
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
