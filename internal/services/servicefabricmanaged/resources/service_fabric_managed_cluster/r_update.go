// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package service_fabric_managed_cluster

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicefabricmanagedcluster/2024-04-01/managedcluster"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicefabricmanagedcluster/2024-04-01/nodetype"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (k Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ClusterResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			clusterClient := metadata.Client.ServiceFabricManaged.ManagedClusterClient
			nodeTypeClient := metadata.Client.ServiceFabricManaged.NodeTypeClient

			id, err := managedcluster.ParseManagedClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			cluster := managedcluster.ManagedCluster{
				Location:   model.Location,
				Name:       pointer.To(model.Name),
				Properties: expandClusterProperties(&model),
				Sku: managedcluster.Sku{
					Name: model.Sku,
				},
			}

			tagsMap := make(map[string]string)
			for k, v := range model.Tags {
				tagsMap[k] = v.(string)
			}
			cluster.Tags = &tagsMap

			if err := clusterClient.CreateOrUpdateThenPoll(ctx, *id, cluster); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			toDelete := make([]string, 0)
			if metadata.ResourceData.HasChange("node_type") {
				o, n := metadata.ResourceData.GetChange("node_type")
				ont := o.([]any)
				nnt := n.([]any)

				for _, on := range ont {
					oldNodeType := on.(map[string]any)
					oldId := oldNodeType["name"].(string)
					found := false
					for _, nt := range nnt {
						newNodeType := nt.(map[string]any)
						newId := newNodeType["name"].(string)
						if oldId == newId {
							found = true
							break
						}
					}

					if !found {
						toDelete = append(toDelete, oldId)
					}
				}
			}

			// Delete the old Node Types
			for _, nodeType := range toDelete {
				nodeTypeId := nodetype.NewNodeTypeID(id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, nodeType)

				if err := nodeTypeClient.DeleteThenPoll(ctx, nodeTypeId); err != nil {
					return fmt.Errorf("deleting %s: %+v", nodeTypeId, err)
				}
			}

			// Ensure the remaining Node Types are up-to-date
			for _, nodeType := range model.NodeTypes {
				props, err := expandNodeTypeProperties(&nodeType)
				if err != nil {
					return fmt.Errorf("while expanding node type %q: %+v", nodeType.Name, err)
				}
				nodeTypeId := nodetype.NewNodeTypeID(id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, nodeType.Name)
				payload := nodetype.NodeType{
					Properties: props,
				}
				if err := nodeTypeClient.CreateOrUpdateThenPoll(ctx, nodeTypeId, payload); err != nil {
					return fmt.Errorf("adding/updating %s: %+v", nodeTypeId, err)
				}
			}

			return nil
		},

		Timeout: 90 * time.Minute,
	}
}
