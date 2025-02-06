// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2024-05-01/agentpools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2024-05-01/snapshots"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type KubernetesNodePoolSnapshotDataSourceModel struct {
	Name             string            `tfschema:"name"`
	ResourceGroup    string            `tfschema:"resource_group_name"`
	SourceNodePoolId string            `tfschema:"source_node_pool_id"`
	Tags             map[string]string `tfschema:"tags"`
}

type KubernetesNodePoolSnapshotDataSource struct{}

var _ sdk.DataSource = KubernetesNodePoolSnapshotDataSource{}

func (r KubernetesNodePoolSnapshotDataSource) ResourceType() string {
	return "azurerm_kubernetes_node_pool_snapshot"
}

func (r KubernetesNodePoolSnapshotDataSource) ModelObject() interface{} {
	return &KubernetesNodePoolSnapshotDataSourceModel{}
}

func (r KubernetesNodePoolSnapshotDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return snapshots.ValidateSnapshotID
}

func (r KubernetesNodePoolSnapshotDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (r KubernetesNodePoolSnapshotDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"source_node_pool_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": tags.SchemaDataSource(),
	}
}

func (r KubernetesNodePoolSnapshotDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.SnapshotClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state KubernetesNodePoolSnapshotDataSourceModel
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

			state.Name = id.SnapshotName

			if model := resp.Model; model != nil {
				state.Tags = pointer.From(model.Tags)

				if props := model.Properties; props != nil {
					if snapshotType := props.SnapshotType; snapshotType != nil && *snapshotType == snapshots.SnapshotTypeNodePool {
						if props.CreationData != nil && props.CreationData.SourceResourceId != nil {
							nodePoolId, err := agentpools.ParseAgentPoolIDInsensitively(*props.CreationData.SourceResourceId)
							if err != nil {
								return err
							}
							state.SourceNodePoolId = nodePoolId.ID()
						}
					}
				}
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
