// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package logic_app_trigger_custom

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/workflows"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/workflowtriggers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceLogicAppTriggerCustomDelete(d *pluginsdk.ResourceData, meta any) error {
	id, err := workflowtriggers.ParseTriggerID(d.Id())
	if err != nil {
		return err
	}

	workflowId := workflows.NewWorkflowID(id.SubscriptionId, id.ResourceGroupName, id.WorkflowName)

	if err = helpers.ResourceLogicAppTriggerRemove(d, meta, workflowId, id.TriggerName); err != nil {
		return fmt.Errorf("removing Trigger %s: %+v", id, err)
	}

	return nil
}
