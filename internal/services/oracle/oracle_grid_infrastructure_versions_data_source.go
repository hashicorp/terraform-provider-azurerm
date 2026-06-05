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
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/giversions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type GridInfrastructureVersionsDataSource struct{}

type GridInfrastructureVersionsModel struct {
	GridInfrastructureVersions []GridInfrastructureVersionModel `tfschema:"versions"`
	Location                   string                           `tfschema:"location"`
	Shape                      string                           `tfschema:"shape"`
	Zone                       string                           `tfschema:"zone"`
}

type GridInfrastructureVersionModel struct {
	Id      string `tfschema:"id"`
	Name    string `tfschema:"name"`
	Version string `tfschema:"version"`
}

func (d GridInfrastructureVersionsDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.Location(),
		"shape": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(giversions.PossibleValuesForSystemShapes(), false),
			Description:  "Filter the versions by system shape. Possible values are 'ExaDbXS', 'Exadata.X9M', and 'Exadata.X11M'.",
		},
		"zone": {
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Description: "Filter the versions by zone",
		},
	}
}

func (d GridInfrastructureVersionsDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"versions": {
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Description: "A list of available Oracle Grid Infrastructure versions and their properties.",
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
				},
			},
		},
	}
}

func (d GridInfrastructureVersionsDataSource) ModelObject() interface{} {
	return &GridInfrastructureVersionsModel{}
}

func (d GridInfrastructureVersionsDataSource) ResourceType() string {
	return "azurerm_oracle_grid_infrastructure_versions"
}

func (d GridInfrastructureVersionsDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return giversions.ValidateGiVersionID
}

func (d GridInfrastructureVersionsDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.GiVersions
			subscriptionId := metadata.Client.Account.SubscriptionId

			state := GridInfrastructureVersionsModel{
				GridInfrastructureVersions: make([]GridInfrastructureVersionModel, 0),
			}
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := giversions.NewLocationID(subscriptionId,
				state.Location)

			options := giversions.ListByLocationOperationOptions{}
			if state.Shape != "" {
				options.Shape = pointer.To(giversions.SystemShapes(state.Shape))
			}
			if state.Zone != "" {
				options.Zone = &state.Zone
			}

			if state.Shape == "" || state.Zone == "" {
				log.Printf("[WARN] GI Versions data source: Shape or Zone parameter is empty. This may result in unfiltered results from the API. Consider specifying both Shape and Zone for more precise version filtering.")
			}

			resp, err := client.ListByLocation(ctx, id, options)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if model := resp.Model; model != nil {
				for _, element := range *model {
					if props := element.Properties; props != nil {
						state.GridInfrastructureVersions = append(state.GridInfrastructureVersions, GridInfrastructureVersionModel{
							Id:      pointer.From(element.Id),
							Name:    pointer.From(element.Name),
							Version: props.Version,
						})
					}
				}
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
