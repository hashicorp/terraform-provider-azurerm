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
	"github.com/hashicorp/go-azure-sdk/resource-manager/relay/2021-11-01/hybridconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.DataSource = RelayHybridConnectionDataResource{}

type RelayHybridConnectionDataResource struct{}

type RelayHybridConnectionDataResourceModel struct {
	Name                        string `tfschema:"name"`
	ResourceGroupName           string `tfschema:"resource_group_name"`
	RelayNamespaceName          string `tfschema:"relay_namespace_name"`
	RequiresClientAuthorization bool   `tfschema:"requires_client_authorization"`
	UserMetadata                string `tfschema:"user_metadata"`
}

func (RelayHybridConnectionDataResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"relay_namespace_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (RelayHybridConnectionDataResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"requires_client_authorization": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"user_metadata": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (RelayHybridConnectionDataResource) ModelObject() interface{} {
	return &RelayHybridConnectionDataResourceModel{}
}

func (RelayHybridConnectionDataResource) ResourceType() string {
	return "azurerm_relay_hybrid_connection"
}

func (RelayHybridConnectionDataResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: *pluginsdk.DefaultTimeout(5 * time.Minute),

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Relay.HybridConnectionsClient

			id, err := hybridconnections.ParseHybridConnectionID(metadata.ResourceData.Id())
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

			state := RelayHybridConnectionDataResourceModel{}

			state.ResourceGroupName = id.ResourceGroupName
			state.RelayNamespaceName = id.NamespaceName

			if model := resp.Model; model != nil {
				state.Name = pointer.From(model.Name)

				if props := model.Properties; props != nil {
					state.RequiresClientAuthorization = pointer.From(model.Properties.RequiresClientAuthorization)
					state.UserMetadata = pointer.From(model.Properties.UserMetadata)
				}
			}

			return metadata.Encode(&state)
		},
	}
}
