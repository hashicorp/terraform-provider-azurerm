// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package logic_app_workflow

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/workflows"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceLogicAppWorkflowDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Logic.WorkflowClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := workflows.ParseWorkflowID(d.Id())
	if err != nil {
		return err
	}

	// lock to prevent against Actions, Parameters or Triggers conflicting
	locks.ByName(id.WorkflowName, logicAppResourceName)
	defer locks.UnlockByName(id.WorkflowName, logicAppResourceName)

	resp, err := client.Delete(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil
		}

		return fmt.Errorf("issuing delete request for Logic App Workflow %s: %+v", id, err)
	}

	return nil
}
