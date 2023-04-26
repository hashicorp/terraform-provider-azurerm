package containers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-02-02-preview/snapshots"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type KubernetesSnapshotDataSourceModel struct {
	Name             string            `tfschema:"name"`
	ResourceGroup    string            `tfschema:"resource_group_name"`
	SourceResourceId string            `tfschema:"source_resource_id"`
	Tags             map[string]string `tfschema:"tags"`
}

type KubernetesSnapshotDataSource struct{}

var _ sdk.DataSource = KubernetesSnapshotDataSource{}

func (r KubernetesSnapshotDataSource) ResourceType() string {
	return "azurerm_kubernetes_snapshot"
}

func (r KubernetesSnapshotDataSource) ModelObject() interface{} {
	return &KubernetesSnapshotDataSourceModel{}
}

func (r KubernetesSnapshotDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return snapshots.ValidateSnapshotID
}

func (r KubernetesSnapshotDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (r KubernetesSnapshotDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"source_resource_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": tags.SchemaDataSource(),
	}
}

func (r KubernetesSnapshotDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.SnapshotClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state KubernetesSnapshotDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := snapshots.NewSnapshotID(subscriptionId, state.ResourceGroup, state.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %v", id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state.Name = id.SnapshotName

			if model.Tags != nil {
				state.Tags = *model.Tags
			}

			if properties := model.Properties; properties != nil {
				if properties.CreationData != nil {
					if properties.CreationData.SourceResourceId != nil {
						state.SourceResourceId = *properties.CreationData.SourceResourceId
					}
				}
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
