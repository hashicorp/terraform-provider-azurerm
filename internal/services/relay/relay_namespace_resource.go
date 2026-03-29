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
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/relay/2021-11-01/namespaces"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/relay/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.ResourceWithUpdate = RelayNamespaceResource{}

type RelayNamespaceResource struct{}

type RelayNamespaceResourceModel struct {
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

func (r RelayNamespaceResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: r.IDValidationFunc(),
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"sku_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(namespaces.SkuNameStandard),
			}, false),
		},

		"tags": commonschema.Tags(),
	}
}

func (r RelayNamespaceResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
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

func (r RelayNamespaceResource) ModelObject() interface{} {
	return &RelayNamespaceResourceModel{}
}

func (r RelayNamespaceResource) ResourceType() string {
	return "azurerm_relay_hybrid_connection"
}

func (r RelayNamespaceResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: *pluginsdk.DefaultTimeout(30 * time.Minute),

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Relay.NamespacesClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			log.Printf("[INFO] preparing arguments for Relay Namespace create.")

			var config RelayNamespaceResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := namespaces.NewNamespaceID(subscriptionId, config.ResourceGroupName, config.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			parameters := namespaces.RelayNamespace{
				Location: location.Normalize(config.Location),
				Sku: &namespaces.Sku{
					Name: namespaces.SkuName(config.SkuName),
					Tier: pointer.To(namespaces.SkuTier(config.SkuName)),
				},
				Properties: &namespaces.RelayNamespaceProperties{},
				Tags:       pointer.To(config.Tags),
			}

			if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r RelayNamespaceResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: *pluginsdk.DefaultTimeout(30 * time.Minute),

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Relay.NamespacesClient

			log.Printf("[INFO] preparing arguments for Relay Namespace update.")

			id, err := namespaces.ParseNamespaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config RelayNamespaceResourceModel
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

			parameters := namespaces.RelayNamespace{
				Location: location.Normalize(config.Location),
				Sku: &namespaces.Sku{
					Name: namespaces.SkuName(config.SkuName),
					Tier: pointer.To(namespaces.SkuTier(config.SkuName)),
				},
				Properties: &namespaces.RelayNamespaceProperties{},
				Tags:       pointer.To(config.Tags),
			}

			if _, err := client.CreateOrUpdate(ctx, *id, parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r RelayNamespaceResource) Read() sdk.ResourceFunc {
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

			state := RelayNamespaceResourceModel{}

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

func (r RelayNamespaceResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: *pluginsdk.DefaultTimeout(60 * time.Minute),

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Relay.NamespacesClient

			id, err := namespaces.ParseNamespaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			log.Printf("[DEBUG] Waiting for %s to be deleted", *id)
			pollerType := custompollers.DeleteRelayNamespacePoller(client, pointer.From(id))
			poller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
			if err := poller.PollUntilDone(ctx); err != nil {
				return err
			}

			return nil
		},
	}
}

func (r RelayNamespaceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validation.StringLenBetween(6, 50)
}
