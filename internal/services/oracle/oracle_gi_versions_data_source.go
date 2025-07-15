// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

package oracle

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-03-01/giversions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type GiVersionsDataSource struct{}

type GiVersionsModel struct {
	Versions []string `tfschema:"versions"`
	Location string   `tfschema:"location"`
	Shape    string   `tfschema:"shape"`
	Zone     string   `tfschema:"zone"`
}

func (d GiVersionsDataSource) Arguments() map[string]*pluginsdk.Schema {
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

func (d GiVersionsDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"versions": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (d GiVersionsDataSource) ModelObject() interface{} {
	return &GiVersionsModel{}
}

func (d GiVersionsDataSource) ResourceType() string {
	return "azurerm_oracle_gi_versions"
}

func (d GiVersionsDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return giversions.ValidateGiVersionID
}

func (d GiVersionsDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.GiVersions
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state GiVersionsModel
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
				fmt.Printf("[WARN] GI Versions data source: Shape or Zone parameter is empty. This may result in unfiltered results from the API. Consider specifying both Shape and Zone for more precise version filtering.")
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
						state.Versions = append(state.Versions, props.Version)
					}
				}
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
