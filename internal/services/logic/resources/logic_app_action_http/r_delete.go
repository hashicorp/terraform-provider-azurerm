// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package logic_app_action_http

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/workflows"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceLogicAppActionHTTPDelete(d *pluginsdk.ResourceData, meta any) error {
	id, err := parse.ActionID(d.Id())
	if err != nil {
		return err
	}

	workflowId := workflows.NewWorkflowID(id.SubscriptionId, id.ResourceGroup, id.WorkflowName)

	if err = helpers.ResourceLogicAppActionRemove(d, meta, workflowId, id.Name); err != nil {
		return fmt.Errorf("removing Action %s: %+v", id, err)
	}

	return nil
}
