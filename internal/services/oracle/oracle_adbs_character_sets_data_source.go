// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

package oracle

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/autonomousdatabasecharactersets"

	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type AdbsCharSetsDataSource struct{}

type AdbsCharSetsModel struct {
	AdbsCharSets []AdbsCharSetModel `tfschema:"character_sets"`
	Location     string             `tfschema:"location"`
}

type AdbsCharSetModel struct {
	CharacterSet string `tfschema:"character_set"`
}

func (d AdbsCharSetsDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.Location(),
	}
}

func (d AdbsCharSetsDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"character_sets": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"character_set": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},
	}
}

func (d AdbsCharSetsDataSource) ModelObject() interface{} {
	return &AdbsCharSetsModel{}
}

func (d AdbsCharSetsDataSource) ResourceType() string {
	return "azurerm_oracle_adbs_character_sets"
}

func (d AdbsCharSetsDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return autonomousdatabasecharactersets.ValidateAutonomousDatabaseCharacterSetID
}

func (d AdbsCharSetsDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.AutonomousDatabaseCharacterSets
			subscriptionId := metadata.Client.Account.SubscriptionId

			state := AdbsCharSetsModel{
				AdbsCharSets: make([]AdbsCharSetModel, 0),
			}
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := autonomousdatabasecharactersets.NewLocationID(subscriptionId, state.Location)

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
						state.AdbsCharSets = append(state.AdbsCharSets, AdbsCharSetModel{
							CharacterSet: props.CharacterSet,
						})
					}
				}
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
