// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package machine_learning_inference_cluster

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/machinelearningcomputes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/machinelearning/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAksInferenceClusterRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).MachineLearning.MachineLearningComputes
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := machinelearningcomputes.ParseComputeID(d.Id())
	if err != nil {
		return err
	}

	d.Set("name", id.ComputeName)

	// Check that Inference Cluster Response can be read
	computeResource, err := client.ComputeGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(computeResource.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on Inference Cluster %q in Workspace %q (Resource Group %q): %+v",
			id.ComputeName, id.WorkspaceName, id.ResourceGroupName, err)
	}

	// Retrieve Machine Learning Workspace ID
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	workspaceId := machinelearningcomputes.NewWorkspaceID(subscriptionId, id.ResourceGroupName, id.WorkspaceName)
	d.Set("machine_learning_workspace_id", workspaceId.ID())

	// use ComputeResource to get to AKS Cluster ID and other properties
	aksComputeProperties := computeResource.Model.Properties.(machinelearningcomputes.AKS)

	// Retrieve AKS Cluster ID
	aksId, err := commonids.ParseKubernetesClusterIDInsensitively(*aksComputeProperties.ResourceId)
	if err != nil {
		return err
	}
	d.Set("kubernetes_cluster_id", aksId.ID())
	clusterPurpose := ""
	if aksComputeProperties.Properties != nil {
		clusterPurpose = pointer.FromEnum(aksComputeProperties.Properties.ClusterPurpose)
	}
	d.Set("cluster_purpose", clusterPurpose)
	d.Set("description", aksComputeProperties.Description)

	// Retrieve location
	if loc := computeResource.Model.Location; loc != nil {
		d.Set("location", location.Normalize(*loc))
	}

	identity, err := helpers.FlattenIdentity(computeResource.Model.Identity)
	if err != nil {
		return fmt.Errorf("flattening `identity`: %+v", err)
	}
	if err := d.Set("identity", identity); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	return tags.FlattenAndSet(d, computeResource.Model.Tags)
}
