// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package machine_learning_workspace

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMachineLearningWorkspaceDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).MachineLearning.Workspaces
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := workspaces.ParseWorkspaceID(d.Id())
	if err != nil {
		return fmt.Errorf("parsing Machine Learning Workspace ID `%q`: %+v", d.Id(), err)
	}

	options := workspaces.DefaultDeleteOperationOptions()
	if meta.(*clients.Client).Features.MachineLearning.PurgeSoftDeletedWorkspaceOnDestroy {
		options = workspaces.DeleteOperationOptions{
			ForceToPurge: pointer.To(true),
		}
	}

	future, err := client.Delete(ctx, *id, options)
	if err != nil {
		return fmt.Errorf("deleting Machine Learning Workspace %q (Resource Group %q): %+v", id.WorkspaceName, id.ResourceGroupName, err)
	}

	if err := future.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting for deletion of Machine Learning Workspace %q (Resource Group %q): %+v", id.WorkspaceName, id.ResourceGroupName, err)
	}

	return nil
}
