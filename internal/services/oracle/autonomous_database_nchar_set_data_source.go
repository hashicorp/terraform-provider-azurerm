// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

package oracle

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/autonomousdatabasenationalcharactersets"

	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type AdbsNCharSetsDataSource struct{}

type AdbsNCharSetsModel struct {
	AdbsCharSets []AdbsNCharSetModel `tfschema:"character_sets"`
	Location     string              `tfschema:"location"`
}

type AdbsNCharSetModel struct {
	CharacterSet string `tfschema:"character_set"`
}

func (d AdbsNCharSetsDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.Location(),
	}
}

func (d AdbsNCharSetsDataSource) Attributes() map[string]*pluginsdk.Schema {
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

func (d AdbsNCharSetsDataSource) ModelObject() interface{} {
	return &AdbsNCharSetsModel{}
}

func (d AdbsNCharSetsDataSource) ResourceType() string {
	return "azurerm_oracle_adbs_national_character_sets"
}

func (d AdbsNCharSetsDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return autonomousdatabasenationalcharactersets.ValidateAutonomousDatabaseNationalCharacterSetID
}

func (d AdbsNCharSetsDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.AutonomousDatabaseNationalCharacterSets
			subscriptionId := metadata.Client.Account.SubscriptionId

			state := AdbsNCharSetsModel{
				AdbsCharSets: make([]AdbsNCharSetModel, 0),
			}
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := autonomousdatabasenationalcharactersets.NewLocationID(subscriptionId, state.Location)

			resp, err := client.ListByLocation(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if model := resp.Model; model != nil {
				for _, element := range *model {
					if element.Properties != nil {
						properties := element.Properties
						state.AdbsCharSets = append(state.AdbsCharSets, AdbsNCharSetModel{
							CharacterSet: properties.CharacterSet,
						})
					}
				}
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
