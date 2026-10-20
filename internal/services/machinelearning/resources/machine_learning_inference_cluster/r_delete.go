// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package machine_learning_inference_cluster

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/machinelearningcomputes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAksInferenceClusterDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).MachineLearning.MachineLearningComputes
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()
	id, err := machinelearningcomputes.ParseComputeID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.ComputeDelete(ctx, *id, machinelearningcomputes.ComputeDeleteOperationOptions{
		UnderlyingResourceAction: pointer.To(machinelearningcomputes.UnderlyingResourceActionDetach),
	})
	if err != nil {
		return fmt.Errorf("deleting Inference Cluster %q in workspace %q (Resource Group %q): %+v",
			id.ComputeName, id.WorkspaceName, id.ResourceGroupName, err)
	}
	if err := future.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting for deletion of Inference Cluster %q in workspace %q (Resource Group %q): %+v",
			id.ComputeName, id.WorkspaceName, id.ResourceGroupName, err)
	}
	return nil
}
