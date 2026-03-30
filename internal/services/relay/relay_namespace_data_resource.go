// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package relay

import (
	"context"
	"fmt"
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

var _ sdk.DataSource = RelayNamespaceDataResource{}

type RelayNamespaceDataResource struct{}

type RelayNamespaceDataResourceModel struct {
	Name                      string            `tfschema:"name"`
	ResourceGroupName         string            `tfschema:"resource_group_name"`
	Location                  string            `tfschema:"location"`
	SkuName                   string            `tfschema:"sku_name"`
	MetricId                  string            `tfschema:"metric_id"`
	PrimaryConnectionString   string            `tfschema:"primary_connection_string"`
	SecondaryConnectionString string            `tfschema:"secondary_connection_string"`
	PrimaryKey                string            `tfschema:"primary_key"`
	SecondaryKey              string            `tfschema:"secondary_key"`
	Tags                      map[string]string `tfschema:"tags"`
}

func (r RelayNamespaceDataResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: r.IDValidationFunc(),
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),
	}
}

func (r RelayNamespaceDataResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"sku_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": commonschema.TagsDataSource(),

		"metric_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"primary_connection_string": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"secondary_connection_string": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"primary_key": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"secondary_key": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},
	}
}

func (r RelayNamespaceDataResource) ModelObject() interface{} {
	return &RelayNamespaceDataResourceModel{}
}

func (r RelayNamespaceDataResource) ResourceType() string {
	return "azurerm_relay_hybrid_connection"
}

func (r RelayNamespaceDataResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: *pluginsdk.DefaultTimeout(5 * time.Minute),

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Relay.NamespacesClient

			id, err := namespaces.ParseNamespaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving: %s: %+v", id, err)
			}

			authRuleId := namespaces.NewAuthorizationRuleID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, "RootManageSharedAccessKey")
			keysResp, err := client.ListKeys(ctx, authRuleId)
			if err != nil {
				return fmt.Errorf("listing keys for %s: %+v", *id, err)
			}

			state := RelayNamespaceDataResourceModel{}

			state.ResourceGroupName = id.ResourceGroupName

			if model := resp.Model; model != nil {
				state.Name = pointer.From(model.Name)
				state.Location = location.NormalizeNilable(&model.Location)

				if sku := model.Sku; sku != nil {
					state.SkuName = string(sku.Name)
				}

				if props := model.Properties; props != nil {
					state.MetricId = pointer.From(props.MetricId)
				}

				state.Tags = pointer.From(model.Tags)
			}

			if model := keysResp.Model; model != nil {
				state.PrimaryConnectionString = pointer.From(model.PrimaryConnectionString)
				state.PrimaryKey = pointer.From(model.PrimaryKey)
				state.SecondaryConnectionString = pointer.From(model.SecondaryConnectionString)
				state.SecondaryKey = pointer.From(model.SecondaryKey)
			}

			return metadata.Encode(&state)
		},
	}
}

func (r RelayNamespaceDataResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validation.StringLenBetween(6, 50)
}
