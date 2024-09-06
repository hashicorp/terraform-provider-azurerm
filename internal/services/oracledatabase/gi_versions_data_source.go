// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

package oracledatabase

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/giversions"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type GiVersionsDataSource struct{}

type GiVersionsModel struct {
	Versions []string `tfschema:"versions"`
}

func (d GiVersionsDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
	}
}

func (d GiVersionsDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{

		// GiVersionProperties
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
	return nil
}

func (d GiVersionsDataSource) ResourceType() string {
	return "azurerm_oracledatabase_gi_versions"
}

func (d GiVersionsDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return giversions.ValidateGiVersionID
}

func (d GiVersionsDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.OracleDatabase.OracleDatabaseClient.GiVersions
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := giversions.NewLocationID(subscriptionId,
				metadata.ResourceData.Get("location_name").(string))

			resp, err := client.ListByLocation(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if model := resp.Model; model != nil {
				output := GiVersionsModel{
					Versions: make([]string, 0),
				}
				for _, element := range *model {
					if element.Properties != nil {
						output.Versions = append(output.Versions, pointer.From(element.Properties.Version))
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
