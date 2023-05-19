package containers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-02-02-preview/managedclusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-02-02-preview/managedclustersnapshots"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-02-02-preview/snapshots"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type KubernetesClusterSnapshotDataSourceModel struct {
	Name            string            `tfschema:"name"`
	ResourceGroup   string            `tfschema:"resource_group_name"`
	SourceClusterId string            `tfschema:"source_cluster_id"`
	Tags            map[string]string `tfschema:"tags"`
}

type KubernetesClusterSnapshotDataSource struct{}

var _ sdk.DataSource = KubernetesClusterSnapshotDataSource{}

func (r KubernetesClusterSnapshotDataSource) ResourceType() string {
	return "azurerm_kubernetes_cluster_snapshot"
}

func (r KubernetesClusterSnapshotDataSource) ModelObject() interface{} {
	return &KubernetesNodePoolSnapshotDataSourceModel{}
}

func (r KubernetesClusterSnapshotDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return snapshots.ValidateSnapshotID
}

func (r KubernetesClusterSnapshotDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (r KubernetesClusterSnapshotDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"source_cluster_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": tags.SchemaDataSource(),
	}
}

func (r KubernetesClusterSnapshotDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.ManagedClusterSnapshotClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state KubernetesClusterSnapshotDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := managedclustersnapshots.NewManagedClusterSnapshotID(subscriptionId, state.ResourceGroup, state.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %v", id, err)
			}

			state.Name = id.ManagedClusterSnapshotName

			if model := resp.Model; model != nil {
				state.Tags = pointer.From(model.Tags)

				if props := model.Properties; props != nil {
					if snapshotType := props.SnapshotType; snapshotType != nil && *snapshotType == managedclustersnapshots.SnapshotTypeManagedCluster {
						if props.CreationData != nil && props.CreationData.SourceResourceId != nil {
							clusterId, err := managedclusters.ParseManagedClusterIDInsensitively(*props.CreationData.SourceResourceId)
							if err != nil {
								return err
							}
							state.SourceClusterId = clusterId.ID()
						}
					}
				}
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
