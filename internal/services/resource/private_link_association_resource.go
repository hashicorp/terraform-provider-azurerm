package resource

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2020-05-01/privatelinkassociation"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2020-05-01/resourcemanagementprivatelink"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"time"
)

var _ sdk.Resource = PrivateLinkAssociationResource{}

type PrivateLinkAssociationResource struct{}

func (r PrivateLinkAssociationResource) ModelObject() interface{} {
	return &PrivateLinkAssociationResourceSchema{}
}

type PrivateLinkAssociationResourceSchema struct {
	ManagementGroupId          string `tfschema:"management_group_id"`
	Name                       string `tfschema:"name"`
	PrivateLinkId              string `tfschema:"private_link_id"`
	PublicNetworkAccessEnabled bool   `tfschema:"public_network_access_enabled"`
	Scope                      string `tfschema:"scope"`
	TenantID                   string `tfschema:"tenant_id"`
}

func (r PrivateLinkAssociationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return privatelinkassociation.ValidatePrivateLinkAssociationID
}

func (r PrivateLinkAssociationResource) ResourceType() string {
	return "azurerm_private_link_association"
}

func (r PrivateLinkAssociationResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			ForceNew:     true,
			Optional:     true,
			Computed:     true,
			Type:         pluginsdk.TypeString,
			ValidateFunc: validation.IsUUID,
		},
		"management_group_id": {
			ForceNew:     true,
			Required:     true,
			Type:         pluginsdk.TypeString,
			ValidateFunc: commonids.ValidateManagementGroupID,
		},
		"private_link_id": {
			ForceNew:     true,
			Required:     true,
			Type:         pluginsdk.TypeString,
			ValidateFunc: resourcemanagementprivatelink.ValidateResourceManagementPrivateLinkID,
		},
		"public_network_access_enabled": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeBool,
		},
	}
}

func (r PrivateLinkAssociationResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"scope": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},
		"tenant_id": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},
	}
}

func (r PrivateLinkAssociationResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resource.PrivateLinkAssociationClient

			var config PrivateLinkAssociationResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			var name string
			if config.Name != "" {
				name = config.Name
			}

			if name == "" {
				name = uuid.New().String()
			}

			managementGroupId, err := commonids.ParseManagementGroupID(config.ManagementGroupId)
			if err != nil {
				return fmt.Errorf("parse management group id: %+v", err)
			}

			id := privatelinkassociation.NewPrivateLinkAssociationID(managementGroupId.GroupId, name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			payload := privatelinkassociation.PrivateLinkAssociationObject{
				Properties: &privatelinkassociation.PrivateLinkAssociationProperties{
					PrivateLink:         pointer.To(config.PrivateLinkId),
					PublicNetworkAccess: r.expandPublicNetworkAccess(config.PublicNetworkAccessEnabled),
				},
			}

			if _, err := client.Put(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r PrivateLinkAssociationResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resource.PrivateLinkAssociationClient
			schema := PrivateLinkAssociationResourceSchema{}

			id, err := privatelinkassociation.ParsePrivateLinkAssociationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if model := resp.Model; model != nil {
				schema.ManagementGroupId = commonids.NewManagementGroupID(id.GroupId).ID()
				schema.Name = id.PlaId
				if prop := model.Properties; prop != nil {
					schema.PublicNetworkAccessEnabled = r.flattenPublicNetworkAccess(prop.PublicNetworkAccess)
					schema.Scope = pointer.From(prop.Scope)
					schema.TenantID = pointer.From(prop.TenantID)
					schema.PrivateLinkId = pointer.From(prop.PrivateLink)
				}
			}

			return metadata.Encode(&schema)
		},
	}
}

func (r PrivateLinkAssociationResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resource.PrivateLinkAssociationClient

			id, err := privatelinkassociation.ParsePrivateLinkAssociationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r PrivateLinkAssociationResource) expandPublicNetworkAccess(input bool) *privatelinkassociation.PublicNetworkAccessOptions {
	output := privatelinkassociation.PublicNetworkAccessOptionsEnabled
	if !input {
		output = privatelinkassociation.PublicNetworkAccessOptionsDisabled
	}

	return &output
}

func (r PrivateLinkAssociationResource) flattenPublicNetworkAccess(input *privatelinkassociation.PublicNetworkAccessOptions) bool {
	if input == nil || *input == privatelinkassociation.PublicNetworkAccessOptionsEnabled {
		return true
	}

	return false
}
