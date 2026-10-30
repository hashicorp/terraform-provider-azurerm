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
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/relay/2021-11-01/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type RelayNamespaceResource struct{}

var _ sdk.ResourceWithUpdate = RelayNamespaceResource{}

type RelayNamespaceResourceModel struct {
	Name                      string         `tfschema:"name"`
	Location                  string         `tfschema:"location"`
	ResourceGroupName         string         `tfschema:"resource_group_name"`
	SkuName                   string         `tfschema:"sku_name"`
	Tags                      map[string]any `tfschema:"tags"`
	MetricID                  string         `tfschema:"metric_id"`
	PrimaryConnectionString   string         `tfschema:"primary_connection_string"`
	SecondaryConnectionString string         `tfschema:"secondary_connection_string"`
	PrimaryKey                string         `tfschema:"primary_key"`
	SecondaryKey              string         `tfschema:"secondary_key"`
}

func (r RelayNamespaceResource) ModelObject() any { //address of an empty model so Encode and Decode know the shape
	return &RelayNamespaceResourceModel{}
}

func (r RelayNamespaceResource) ResourceType() string {
	return "azurerm_relay_namespace"
}

func (r RelayNamespaceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return namespaces.ValidateNamespaceID
}

func (r RelayNamespaceResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(6, 50),
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"sku_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(namespaces.PossibleValuesForSkuName(), false),
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

func (r RelayNamespaceResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Relay.NamespacesClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model RelayNamespaceResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := namespaces.NewNamespaceID(subscriptionId, model.ResourceGroupName, model.Name)

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

			//check is sku name has changed, if so, set the sku in the parameters to update it, haschange for every property that does not have forcenew set to true, otherwise the update will not be called
			skuTier := namespaces.SkuTier(model.SkuName)
			parameters := namespaces.RelayNamespace{
				Location: location.Normalize(model.Location),
				Sku: &namespaces.Sku{
					Name: namespaces.SkuName(model.SkuName),
					Tier: &skuTier,
				},
				Properties: &namespaces.RelayNamespaceProperties{},
				Tags:       tags.Expand(model.Tags),
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, parameters, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r RelayNamespaceResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Relay.NamespacesClient

			id, err := namespaces.ParseNamespaceID(metadata.ResourceData.Id())

			if err != nil {
				return fmt.Errorf("parsing namespace ID: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			var state RelayNamespaceResourceModel

			state.Name = id.NamespaceName
			state.ResourceGroupName = id.ResourceGroupName

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
			state.Tags = tags.Flatten(model.Tags)
			if model.Sku != nil {
				state.SkuName = string(model.Sku.Name)
			}

			authRuleId := namespaces.NewAuthorizationRuleID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, "RootManageSharedAccessKey")

			keysResp, err := client.ListKeys(ctx, authRuleId)
			if err != nil {
				return fmt.Errorf("listing keys for %s: %+v", *id, err)
			}
			if keysModel := keysResp.Model; keysModel != nil {
				state.PrimaryConnectionString = pointer.From(keysModel.PrimaryConnectionString)
				state.SecondaryConnectionString = pointer.From(keysModel.SecondaryConnectionString)
				state.PrimaryKey = pointer.From(keysModel.PrimaryKey)
				state.SecondaryKey = pointer.From(keysModel.SecondaryKey)
			}
			return metadata.Encode(&state)
		},
	}
}

func (r RelayNamespaceResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Relay.NamespacesClient
			id, err := namespaces.ParseNamespaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			// we can't make use of the Future here due to a bug where 404 isn't tracked as Successful
			log.Printf("[DEBUG] Waiting for %s to be deleted", *id)
			poller := custompollers.NewEventualConsistencyPoller(1, func(pollerCtx context.Context) (*http.Response, error) {
				resp, err := client.Get(pollerCtx, *id)
				return resp.HttpResponse, err
			}, &custompollers.EventualConsistencyPollerOptions{
				Interval:         15 * time.Second,
				TargetStatusCode: pointer.To(http.StatusNotFound),
			})
			if err := poller.PollUntilDone(ctx); err != nil {
				return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r RelayNamespaceResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Relay.NamespacesClient

			id, err := namespaces.ParseNamespaceID(metadata.ResourceData.Id())

			if err != nil {
				return fmt.Errorf("parsing namespace ID: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil || resp.Model == nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			var model RelayNamespaceResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existingModel := *resp.Model

			//check is sku name has changed, if so, set the sku in the parameters to update it, haschange for every property that does not have forcenew set to true, otherwise the update will not be called
			if metadata.ResourceData.HasChange("sku_name") {
				existingModel.Sku = &namespaces.Sku{
					Name: namespaces.SkuName(model.SkuName),
					Tier: pointer.To(namespaces.SkuTier(model.SkuName)),
				}
			}

			if metadata.ResourceData.HasChange("tags") {
				existingModel.Tags = tags.Expand(model.Tags)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, existingModel); err != nil {
				return fmt.Errorf("creating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
