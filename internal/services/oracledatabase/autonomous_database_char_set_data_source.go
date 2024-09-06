// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

package oracledatabase

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/autonomousdatabasecharactersets"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type AdbsCharSetsDataSource struct{}

type AdbsCharSetsModel struct {
	AdbsCharSets []AdbsCharSetModel `tfschema:"character_sets"`
}

type AdbsCharSetModel struct {
	Name         string `tfschema:"name"`
	CharacterSet string `tfschema:"character_set"`
}

func (d AdbsCharSetsDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
	}
}

func (d AdbsCharSetsDataSource) Attributes() map[string]*pluginsdk.Schema {
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

func (d AdbsCharSetsDataSource) ModelObject() interface{} {
	return nil
}

func (d AdbsCharSetsDataSource) ResourceType() string {
	return "azurerm_oracledatabase_adbs_character_sets"
}

func (d AdbsCharSetsDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return autonomousdatabasecharactersets.ValidateAutonomousDatabaseCharacterSetID
}

func (d AdbsCharSetsDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.OracleDatabase.OracleDatabaseClient.AutonomousDatabaseCharacterSets
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := autonomousdatabasecharactersets.NewLocationID(subscriptionId,
				metadata.ResourceData.Get("location_name").(string))

			resp, err := client.ListByLocation(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if model := resp.Model; model != nil {
				output := AdbsCharSetsModel{
					AdbsCharSets: make([]AdbsCharSetModel, 0),
				}
				for _, element := range *model {
					if element.Properties != nil {
						properties := element.Properties
						output.AdbsCharSets = append(output.AdbsCharSets, AdbsCharSetModel{
							Name:         pointer.From(element.Name),
							CharacterSet: pointer.From(properties.CharacterSet),
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
