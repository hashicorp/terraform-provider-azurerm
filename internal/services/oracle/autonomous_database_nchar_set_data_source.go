// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

package oracle

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/autonomousdatabasecharactersets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/autonomousdatabasenationalcharactersets"

	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type AdbsNCharSetsDataSource struct{}

type AdbsNCharSetsModel struct {
	AdbsNCharSets []AdbsNCharSetModel `tfschema:"character_sets"`
}

type AdbsNCharSetModel struct {
	Name          string `tfschema:"name"`
	NCharacterSet string `tfschema:"character_set"`
}

func (d AdbsNCharSetsDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
	}
}

func (d AdbsNCharSetsDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"character_sets": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
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
	return nil
}

func (d AdbsNCharSetsDataSource) ResourceType() string {
	return "azurerm_oracle_adbs_national_character_sets"
}

func (d AdbsNCharSetsDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return autonomousdatabasecharactersets.ValidateAutonomousDatabaseCharacterSetID
}

func (d AdbsNCharSetsDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.AutonomousDatabaseNationalCharacterSets
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := autonomousdatabasenationalcharactersets.NewLocationID(subscriptionId,
				metadata.ResourceData.Get("location_name").(string))

			resp, err := client.ListByLocation(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if model := resp.Model; model != nil {
				output := AdbsNCharSetsModel{
					AdbsNCharSets: make([]AdbsNCharSetModel, 0),
				}
				for _, element := range *model {
					if element.Properties != nil {
						properties := element.Properties
						output.AdbsNCharSets = append(output.AdbsNCharSets, AdbsNCharSetModel{
							Name:          pointer.From(element.Name),
							NCharacterSet: pointer.From(properties.CharacterSet),
						})
					}
				}
				metadata.SetID(id)
				if err := metadata.Encode(&output); err != nil {
					return fmt.Errorf("encoding %s: %+v", id, err)
				}
			}
			return nil
		},
	}
}
