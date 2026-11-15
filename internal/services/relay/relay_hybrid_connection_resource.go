// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package relay

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/relay/2021-11-01/hybridconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type RelayHybridConnectionResource struct{}

var _ sdk.ResourceWithUpdate = RelayHybridConnectionResource{}

type RelayHybridConnectionResourceModel struct {
	Name                        string `tfschema:"name"`
	ResourceGroupName           string `tfschema:"resource_group_name"`
	RelayNamespaceName          string `tfschema:"relay_namespace_name"`
	RequiresClientAuthorization bool   `tfschema:"requires_client_authorization"`
	UserMetadata                string `tfschema:"user_metadata"`
}

func (r RelayHybridConnectionResource) ModelObject() any {
	return &RelayHybridConnectionResourceModel{}
}

func (r RelayHybridConnectionResource) ResourceType() string {
	return "azurerm_relay_hybrid_connection"
}

func (r RelayHybridConnectionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return hybridconnections.ValidateHybridConnectionID
}

func (r RelayHybridConnectionResource) Arguments() map[string]*pluginsdk.Schema {
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

func (r RelayHybridConnectionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r RelayHybridConnectionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Relay.HybridConnectionsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model RelayHybridConnectionResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := hybridconnections.NewHybridConnectionID(subscriptionId, model.ResourceGroupName, model.RelayNamespaceName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
					}
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			parameters := hybridconnections.HybridConnection{
				Properties: &hybridconnections.HybridConnectionProperties{
					RequiresClientAuthorization: pointer.To(model.RequiresClientAuthorization),
					UserMetadata:                pointer.To(model.UserMetadata),
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

func (r RelayHybridConnectionResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
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

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := RelayHybridConnectionResourceModel{
				Name:               id.HybridConnectionName,
				ResourceGroupName:  id.ResourceGroupName,
				RelayNamespaceName: id.NamespaceName,
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.RequiresClientAuthorization = pointer.From(props.RequiresClientAuthorization)
					state.UserMetadata = pointer.From(props.UserMetadata)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r RelayHybridConnectionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
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
			poller := custompollers.NewEventualConsistencyPoller(1, func(pollerCtx context.Context) (*http.Response, error) {
				resp, err := client.Get(pollerCtx, *id)
				return resp.HttpResponse, err
			}, &custompollers.EventualConsistencyPollerOptions{
				Interval:         15 * time.Second,
				TargetStatusCode: pointer.To(http.StatusNotFound),
			})
			if err := poller.PollUntilDone(ctx); err != nil {
				return fmt.Errorf("waiting for %s to be deleted: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r RelayHybridConnectionResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Relay.HybridConnectionsClient

			id, err := hybridconnections.ParseHybridConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			var model RelayHybridConnectionResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existingModel := *resp.Model
			if existingModel.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id)
			}

			if metadata.ResourceData.HasChange("user_metadata") {
				existingModel.Properties.UserMetadata = pointer.To(model.UserMetadata)
			}

			if _, err := client.CreateOrUpdate(ctx, *id, existingModel); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}
			return nil
		},
	}
}
