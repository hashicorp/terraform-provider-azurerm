package hybridkubernetes

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridkubernetes/2021-10-01/connectedclusters"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ConnectedClusterResource struct{}

type ConnectedClusterModel struct {
	Name                      string          `tfschema:"name"`
	ResourceGroup             string          `tfschema:"resource_group_name"`
	Location                  string          `tfschema:"location"`
	Identity                  []IdentityModel `tfschema:"identity"`
	AgentPublicKeyCertificate string          `tfschema:"agent_public_key_certificate"`
	Distribution              string          `tfschema:"distribution"`
	Infrastructure            string          `tfschema:"infrastructure"`
}

type IdentityModel struct {
	Type string `tfschema:"type"`
}

func (r ConnectedClusterResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"identity": commonschema.SystemAssignedUserAssignedIdentityRequiredForceNew(),

		"agent_public_key_certificate": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"distribution": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      "generic",
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"infrastructure": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r ConnectedClusterResource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{}
}

func (r ConnectedClusterResource) ModelObject() interface{} {
	return &ConnectedClusterModel{}
}

func (r ConnectedClusterResource) ResourceType() string {
	return "azurerm_connected_cluster"
}

func (r ConnectedClusterResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ConnectedClusterModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			client := metadata.Client.HybridKubernetes.ConnectedClustersClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := connectedclusters.NewConnectedClusterID(subscriptionId, model.ResourceGroup, model.Name)

			existing, err := client.ConnectedClusterGet(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			identity, err := ExpandIdentity(model.Identity)
			if err != nil {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			props := connectedclusters.ConnectedClusterProperties{
				AgentPublicKeyCertificate: model.AgentPublicKeyCertificate,
			}

			if model.Distribution != "" {
				props.Distribution = utils.String(model.Distribution)
			}

			if model.Infrastructure != "" {
				props.Infrastructure = utils.String(model.Infrastructure)
			}

			connectedCluster := connectedclusters.ConnectedCluster{
				Id:         utils.String(id.ID()),
				Identity:   *identity,
				Location:   model.Location,
				Name:       utils.String(model.Name),
				Properties: props,
			}

			if err := client.ConnectedClusterCreateThenPoll(ctx, id, connectedCluster); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ConnectedClusterResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.HybridKubernetes.ConnectedClustersClient
			id, err := connectedclusters.ParseConnectedClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.ConnectedClusterGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			if model := resp.Model; model != nil {
				props := model.Properties

				identity, err := FlattenIdentity(&model.Identity)
				if err != nil {
					return fmt.Errorf("reading %s: %+v", *id, err)
				}
				state := ConnectedClusterModel{
					Name:                      id.ClusterName,
					ResourceGroup:             id.ResourceGroupName,
					Location:                  model.Location,
					Identity:                  identity,
					AgentPublicKeyCertificate: props.AgentPublicKeyCertificate,
				}

				if props.Distribution != nil {
					state.Distribution = *props.Distribution
				}

				if props.Infrastructure != nil {
					state.Infrastructure = *props.Infrastructure
				}

				return metadata.Encode(&state)
			}
			return nil
		},
	}
}

func (r ConnectedClusterResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.HybridKubernetes.ConnectedClustersClient
			id, err := connectedclusters.ParseConnectedClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("deleting %s", *id)

			if err := client.ConnectedClusterDeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ConnectedClusterResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return connectedclusters.ValidateConnectedClusterID
}
