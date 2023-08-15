// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package communication

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/communication/2023-03-31/communicationservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/communication/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.DataSource = CommunicationServiceDataSource{}

type CommunicationServiceDataSource struct{}

type CommunicationServiceDataSourceModel struct {
	Name                      string            `tfschema:"name"`
	ResourceGroupName         string            `tfschema:"resource_group_name"`
	DataLocation              string            `tfschema:"data_location"`
	PrimaryConnectionString   string            `tfschema:"primary_connection_string"`
	PrimaryKey                string            `tfschema:"primary_key"`
	SecondaryConnectionString string            `tfschema:"secondary_connection_string"`
	SecondaryKey              string            `tfschema:"secondary_key"`
	Tags                      map[string]string `tfschema:"tags"`
}

func (CommunicationServiceDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.CommunicationServiceName,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (CommunicationServiceDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"data_location": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"primary_connection_string": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"primary_key": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"secondary_connection_string": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"secondary_key": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": commonschema.TagsDataSource(),
	}
}

func (CommunicationServiceDataSource) ModelObject() interface{} {
	return &CommunicationServiceDataSourceModel{}
}

func (CommunicationServiceDataSource) ResourceType() string {
	return "azurerm_communication_service"
}

func (CommunicationServiceDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Communication.ServiceClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state CommunicationServiceDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := communicationservices.NewCommunicationServiceID(subscriptionId, state.ResourceGroupName, state.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.DataLocation = props.DataLocation
				}

				state.Tags = pointer.From(model.Tags)
			}

			keysResp, err := client.ListKeys(ctx, id)
			if err != nil {
				log.Printf("[WARN] listing keys for %s: %+v", id, err)
			} else if model := keysResp.Model; model != nil {
				state.PrimaryConnectionString = pointer.From(model.PrimaryConnectionString)
				state.PrimaryKey = pointer.From(model.PrimaryKey)
				state.SecondaryConnectionString = pointer.From(model.SecondaryConnectionString)
				state.SecondaryKey = pointer.From(model.SecondaryKey)
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
