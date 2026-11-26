// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package logic_app_action_custom

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/workflows"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceLogicAppActionCustomRead(d *pluginsdk.ResourceData, meta any) error {
	id, err := parse.ActionID(d.Id())
	if err != nil {
		return err
	}

	workflowId := workflows.NewWorkflowID(id.SubscriptionId, id.ResourceGroup, id.WorkflowName)

	t, app, err := helpers.RetrieveLogicAppAction(d, meta, workflowId, id.Name)
	if err != nil {
		return err
	}

	if t == nil {
		log.Printf("[DEBUG] Logic App %q (Resource Group %q) does not contain Action %q - removing from state", id.WorkflowName, id.ResourceGroup, id.Name)
		d.SetId("")
		return nil
	}

	action := *t

	d.Set("name", id.Name)
	d.Set("logic_app_id", app.Id)

	body, err := json.Marshal(action)
	if err != nil {
		return fmt.Errorf("serializing `body` for Action %q: %+v", id.Name, err)
	}

	if err := d.Set("body", string(body)); err != nil {
		return fmt.Errorf("setting `body` for Action %q: %+v", id.Name, err)
	}

	return nil
}
