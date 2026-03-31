// Copyright IBM Corp. 2014, 2025
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
	"github.com/hashicorp/go-azure-sdk/resource-manager/relay/2021-11-01/hybridconnections"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/relay/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.ResourceWithUpdate = RelayHybridConnectionResource{}

type RelayHybridConnectionResource struct{}

type RelayHybridConnectionResourceModel struct {
	Name                        string `tfschema:"name"`
	ResourceGroupName           string `tfschema:"resource_group_name"`
	RelayNamespaceName          string `tfschema:"relay_namespace_name"`
	RequiresClientAuthorization bool   `tfschema:"requires_client_authorization"`
	UserMetadata                string `tfschema:"user_metadata"`
}

func (RelayHybridConnectionResource) Arguments() map[string]*pluginsdk.Schema {
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

		"requires_client_authorization": {
			Type:     pluginsdk.TypeBool,
			Default:  true,
			ForceNew: true,
			Optional: true,
		},

		"user_metadata": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (RelayHybridConnectionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (RelayHybridConnectionResource) ModelObject() interface{} {
	return &RelayHybridConnectionResourceModel{}
}

func (RelayHybridConnectionResource) ResourceType() string {
	return "azurerm_relay_hybrid_connection"
}

func (r RelayHybridConnectionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: *pluginsdk.DefaultTimeout(30 * time.Minute),

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Relay.HybridConnectionsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			log.Printf("[INFO] preparing arguments for Relay Hybrid Connection creation.")

			var config RelayHybridConnectionResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := hybridconnections.NewHybridConnectionID(subscriptionId, config.ResourceGroupName, config.RelayNamespaceName, config.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			parameters := hybridconnections.HybridConnection{
				Properties: &hybridconnections.HybridConnectionProperties{
					RequiresClientAuthorization: &config.RequiresClientAuthorization,
					UserMetadata:                &config.UserMetadata,
				},
			}

			if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r RelayHybridConnectionResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: *pluginsdk.DefaultTimeout(30 * time.Minute),

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Relay.HybridConnectionsClient

			log.Printf("[INFO] preparing arguments for Relay Hybrid Connection update.")

			id, err := hybridconnections.ParseHybridConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config RelayHybridConnectionResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving: %s: %+v", id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}

			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}

			parameters := hybridconnections.HybridConnection{
				Properties: &hybridconnections.HybridConnectionProperties{
					UserMetadata: &config.UserMetadata,
				},
			}

			if _, err := client.CreateOrUpdate(ctx, *id, parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (RelayHybridConnectionResource) Read() sdk.ResourceFunc {
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

			state := RelayHybridConnectionResourceModel{}

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

func (RelayHybridConnectionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: *pluginsdk.DefaultTimeout(30 * time.Minute),

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Relay.HybridConnectionsClient

			id, err := hybridconnections.ParseHybridConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			log.Printf("[INFO] Waiting for %s to be deleted", *id)
			pollerType := custompollers.DeleteRelayHybridConnectionPoller(client, pointer.From(id))
			poller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
			if err := poller.PollUntilDone(ctx); err != nil {
				return err
			}

			return nil
		},
	}
}

func (RelayHybridConnectionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return hybridconnections.ValidateHybridConnectionID
}
