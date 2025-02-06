// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datafactory

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type TriggerSchedulesDataSource struct{}

type TriggerSchedulesDataSourceModel struct {
	DataFactoryID string   `tfschema:"data_factory_id"`
	Items         []string `tfschema:"items"`
}

func (d TriggerSchedulesDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"data_factory_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: factories.ValidateFactoryID,
		},
	}
}

func (d TriggerSchedulesDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"items": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (d TriggerSchedulesDataSource) ModelObject() interface{} {
	return &TriggerSchedulesDataSourceModel{}
}

func (d TriggerSchedulesDataSource) ResourceType() string {
	return "azurerm_data_factory_trigger_schedules"
}

func (d TriggerSchedulesDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model TriggerSchedulesDataSourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			client := metadata.Client.DataFactory.TriggersClient

			dataFactoryId, err := factories.ParseFactoryID(model.DataFactoryID)
			if err != nil {
				return err
			}

			iter, err := client.ListByFactoryComplete(ctx, dataFactoryId.ResourceGroupName, dataFactoryId.FactoryName)
			if err != nil {
				if v, ok := err.(autorest.DetailedError); ok {
					if utils.ResponseWasNotFound(autorest.Response{Response: v.Response}) {
						return fmt.Errorf("fetching triggers list for %s", dataFactoryId)
					}
				} else {
					return fmt.Errorf("fetching triggers list for %s: %+v", dataFactoryId, err)
				}
				return fmt.Errorf("fetching triggers list for %s: %+v", dataFactoryId, err)
			}

			triggers := []string{}
			for iter.NotDone() {
				trigger := iter.Value()
				triggers = append(triggers, *trigger.Name)
				if err := iter.NextWithContext(ctx); err != nil {
					return fmt.Errorf("fetching triggers list from Azure Data Factory %q, advancing iterator failed: %+v", dataFactoryId.ID(), err)
				}
			}

			metadata.SetID(dataFactoryId)
			model.Items = triggers

			return metadata.Encode(&model)
		},
	}
}
