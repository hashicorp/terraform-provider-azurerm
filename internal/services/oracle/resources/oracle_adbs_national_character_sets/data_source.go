// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package oracle_adbs_national_character_sets

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/autonomousdatabasenationalcharactersets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DataSource struct{}

type AdbsNCharSetsModel struct {
	AdbsCharSets []AdbsNCharSetModel `tfschema:"character_sets"`
	Location     string              `tfschema:"location"`
}

type AdbsNCharSetModel struct {
	CharacterSet string `tfschema:"character_set"`
}

func (d DataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.Location(),
	}
}

func (d DataSource) Attributes() map[string]*pluginsdk.Schema {
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

func (d DataSource) ModelObject() any {
	return &AdbsNCharSetsModel{}
}

func (d DataSource) ResourceType() string {
	return "azurerm_oracle_adbs_national_character_sets"
}

func (d DataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return autonomousdatabasenationalcharactersets.ValidateAutonomousDatabaseNationalCharacterSetID
}

func (d DataSource) Read() sdk.ResourceFunc {
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
