// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package logic_app_trigger_http_request

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/workflowtriggers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceLogicAppTriggerHttpRequestRead(d *pluginsdk.ResourceData, meta any) error {
	id, err := workflowtriggers.ParseTriggerID(d.Id())
	if err != nil {
		return err
	}

	t, app, url, err := helpers.RetrieveLogicAppTrigger(d, meta, *id)
	if err != nil {
		return err
	}

	if t == nil {
		log.Printf("[DEBUG] Logic App %q (Resource Group %q) does not contain %s - removing from state", id.WorkflowName, id.ResourceGroupName, id.ID())
		d.SetId("")
		return nil
	}

	trigger := *t

	d.Set("name", id.TriggerName)
	d.Set("logic_app_id", app.Id)
	d.Set("callback_url", url)

	v := trigger["inputs"]
	if v == nil {
		return fmt.Errorf("`inputs` was nil for HTTP Trigger %s", id)
	}

	inputs, ok := v.(map[string]any)
	if !ok {
		return fmt.Errorf("parsing `inputs` for HTTP Trigger %s", id)
	}

	if method := inputs["method"]; method != nil {
		d.Set("method", method.(string))
	}

	if relativePath := inputs["relativePath"]; relativePath != nil {
		d.Set("relative_path", relativePath.(string))
	}

	if schemaRaw := inputs["schema"]; schemaRaw != nil {
		schema, err := json.Marshal(schemaRaw)
		if err != nil {
			return fmt.Errorf("serializing the Schema to JSON: %+v", err)
		}

		d.Set("schema", string(schema))
	}

	return nil
}
