// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

package oracle

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/giversions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type GiVersionsDataSource struct{}

type GiVersionsModel struct {
	Versions []string `tfschema:"versions"`
	Location string   `tfschema:"location"`
}

func (d GiVersionsDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.Location(),
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

			resp, err := client.ListByLocation(ctx, id)
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
