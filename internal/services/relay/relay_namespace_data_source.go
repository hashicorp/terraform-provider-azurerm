// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package relay

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/relay/2021-11-01/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type RelayNamespaceDataSource struct{}

type RelayNamespaceDataSourceModel struct {
	Name                      string            `tfschema:"name"`
	ResourceGroupName         string            `tfschema:"resource_group_name"`
	Location                  string            `tfschema:"location"`
	SkuName                   string            `tfschema:"sku_name"`
	Tags                      map[string]string `tfschema:"tags"`
	MetricID                  string            `tfschema:"metric_id"`
	PrimaryConnectionString   string            `tfschema:"primary_connection_string"`
	SecondaryConnectionString string            `tfschema:"secondary_connection_string"`
	PrimaryKey                string            `tfschema:"primary_key"`
	SecondaryKey              string            `tfschema:"secondary_key"`
}

var _ sdk.DataSource = RelayNamespaceDataSource{}

func (d RelayNamespaceDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (d RelayNamespaceDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"metric_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"primary_connection_string": {
			Type:      pluginsdk.TypeString,
			Sensitive: true,
			Computed:  true,
		},

		"primary_key": {
			Type:      pluginsdk.TypeString,
			Sensitive: true,
			Computed:  true,
		},

		"secondary_connection_string": {
			Type:      pluginsdk.TypeString,
			Sensitive: true,
			Computed:  true,
		},

		"secondary_key": {
			Type:      pluginsdk.TypeString,
			Sensitive: true,
			Computed:  true,
		},

		"sku_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": commonschema.TagsDataSource(),
	}
}

func (d RelayNamespaceDataSource) ModelObject() any {
	return &RelayNamespaceDataSourceModel{}
}

func (d RelayNamespaceDataSource) ResourceType() string {
	return "azurerm_relay_namespace"
}

func (d RelayNamespaceDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Relay.NamespacesClient

			var state RelayNamespaceDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := namespaces.NewNamespaceID(metadata.Client.Account.SubscriptionId, state.ResourceGroupName, state.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			props := model.Properties
			if props == nil {
				return fmt.Errorf("retrieving %s: properties were nil", id)
			}

			state.Location = location.Normalize(model.Location)
			state.MetricID = pointer.From(props.MetricId)
			state.Tags = pointer.From(model.Tags)
			if model.Sku != nil {
				state.SkuName = string(model.Sku.Name)
			}

			authRuleId := namespaces.NewAuthorizationRuleID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, "RootManageSharedAccessKey")

			keysResp, err := client.ListKeys(ctx, authRuleId)
			if err != nil {
				log.Printf("[WARN] Unable to List Write keys for %s: %s", id, err)
			} else if keysModel := keysResp.Model; keysModel != nil {
				state.PrimaryConnectionString = pointer.From(keysModel.PrimaryConnectionString)
				state.SecondaryConnectionString = pointer.From(keysModel.SecondaryConnectionString)
				state.PrimaryKey = pointer.From(keysModel.PrimaryKey)
				state.SecondaryKey = pointer.From(keysModel.SecondaryKey)
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
