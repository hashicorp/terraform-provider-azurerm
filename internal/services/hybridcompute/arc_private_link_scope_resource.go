package hybridcompute

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2022-11-10/privatelinkscopes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"time"
)

type PrivateLinkScopeModel struct {
	Name                       string            `tfschema:"name"`
	Location                   string            `tfschema:"location"`
	ResourceGroupName          string            `tfschema:"resource_group_name"`
	Tags                       map[string]string `tfschema:"tags"`
	PublicNetworkAccessEnabled bool              `tfschema:"public_network_access_enabled"`
}

var _ sdk.Resource = ArcPrivateLinkScopeResource{}

// ArcPrivateLinkScopeResource is a Resource implementation for the Azure Arc Private Link Scope resource
type ArcPrivateLinkScopeResource struct{}

func (a ArcPrivateLinkScopeResource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"resource_group_name": commonschema.ResourceGroupName(),
		"location":            commonschema.Location(),
		"public_network_access_enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"tags": commonschema.Tags(),
	}
}

func (a ArcPrivateLinkScopeResource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{}
}

func (a ArcPrivateLinkScopeResource) ModelObject() interface{} {
	return &PrivateLinkScopeModel{}
}

func (a ArcPrivateLinkScopeResource) ResourceType() string {
	return "azurerm_arc_private_link_scope"
}

func (a ArcPrivateLinkScopeResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model PrivateLinkScopeModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.HybridCompute.PrivateLinkScopesClient
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := privatelinkscopes.NewProviderPrivateLinkScopeID(subscriptionId, model.ResourceGroupName, model.Name)
			existing, err := client.PrivateLinkScopesGet(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(a.ResourceType(), id)
			}

			properties := privatelinkscopes.HybridComputePrivateLinkScope{
				Location: location.Normalize(model.Location),
				Name:     &model.Name,
				Tags:     &model.Tags,
			}

			var publicNetwork privatelinkscopes.PublicNetworkAccessType

			if model.PublicNetworkAccessEnabled {
				publicNetwork = privatelinkscopes.PublicNetworkAccessTypeEnabled
			} else {
				publicNetwork = privatelinkscopes.PublicNetworkAccessTypeDisabled
			}

			properties.Properties.PublicNetworkAccess = &publicNetwork

			if _, err := client.PrivateLinkScopesCreateOrUpdate(ctx, id, properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (a ArcPrivateLinkScopeResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.HybridCompute.PrivateLinkScopesClient

			id, err := privatelinkscopes.ParseProviderPrivateLinkScopeID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing id: %+v", err)
			}

			var model PrivateLinkScopeModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.PrivateLinkScopesGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := existing.Model
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			var publicNetwork privatelinkscopes.PublicNetworkAccessType

			if metadata.ResourceData.HasChange("public_network_access_enabled") {
				if model.PublicNetworkAccessEnabled {
					publicNetwork = privatelinkscopes.PublicNetworkAccessTypeEnabled
				} else {
					publicNetwork = privatelinkscopes.PublicNetworkAccessTypeDisabled
				}

				properties.Properties.PublicNetworkAccess = &publicNetwork
			}

			if _, err := client.PrivateLinkScopesCreateOrUpdate(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (a ArcPrivateLinkScopeResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.HybridCompute.PrivateLinkScopesClient

			id, err := privatelinkscopes.ParseProviderPrivateLinkScopeID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing id: %+v", err)
			}

			existing, err := client.PrivateLinkScopesGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := PrivateLinkScopeModel{
				Name:              id.PrivateLinkScopeName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := existing.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				if model.Tags != nil {
					state.Tags = *model.Tags
				}
				if model.Properties != nil && model.Properties.PublicNetworkAccess != nil {
					state.PublicNetworkAccessEnabled = *model.Properties.PublicNetworkAccess == privatelinkscopes.PublicNetworkAccessTypeEnabled
				}
			}

			return metadata.Encode(state)
		},
	}
}

func (a ArcPrivateLinkScopeResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.HybridCompute.PrivateLinkScopesClient

			id, err := privatelinkscopes.ParseProviderPrivateLinkScopeID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing id: %+v", err)
			}

			if err := client.PrivateLinkScopesDeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (a ArcPrivateLinkScopeResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return privatelinkscopes.ValidateProviderPrivateLinkScopeID
}
