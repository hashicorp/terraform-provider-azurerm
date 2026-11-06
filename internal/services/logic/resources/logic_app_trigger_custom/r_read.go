// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package logic_app_trigger_custom

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/workflowtriggers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceLogicAppTriggerCustomRead(d *pluginsdk.ResourceData, meta any) error {
	id, err := workflowtriggers.ParseTriggerID(d.Id())
	if err != nil {
		return err
	}

	t, app, url, err := helpers.RetrieveLogicAppTrigger(d, meta, *id)
	if err != nil {
		return err
	}

	if t == nil {
		log.Printf("[DEBUG] Logic App %q (Resource Group %q) does not contain Trigger %q - removing from state", id.WorkflowName, id.ResourceGroupName, id.TriggerName)
		d.SetId("")
		return nil
	}

	trigger := *t

	d.Set("name", id.TriggerName)
	d.Set("logic_app_id", app.Id)
	d.Set("callback_url", pointer.From(url))

	// Azure returns an additional field called evaluatedRecurrence in the trigger body which
	// is a copy of the recurrence specified in the body property and breaks the diff suppress logic
	delete(trigger, "evaluatedRecurrence")

	body, err := json.Marshal(trigger)
	if err != nil {
		return fmt.Errorf("serializing `body` for %s: %+v", id.ID(), err)
	}
	log.Printf("[DEBUG] logic_custom_trigger body is: %s", string(body))

	if err := d.Set("body", string(body)); err != nil {
		return fmt.Errorf("setting `body` for %s: %+v", id.ID(), err)
	}

	return nil
}
