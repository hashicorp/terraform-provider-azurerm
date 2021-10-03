package logic

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func resourceLogicAppTriggerCustom() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLogicAppTriggerCustomCreateUpdate,
		Read:   resourceLogicAppTriggerCustomRead,
		Update: resourceLogicAppTriggerCustomCreateUpdate,
		Delete: resourceLogicAppTriggerCustomDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.TriggerID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"logic_app_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"body": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
			},
		},
	}
}

func resourceLogicAppTriggerCustomCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	workflowId, err := parse.WorkflowID(d.Get("logic_app_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewTriggerID(workflowId.SubscriptionId, workflowId.ResourceGroup, workflowId.Name, d.Get("name").(string))

	bodyRaw := d.Get("body").(string)

	var body map[string]interface{}
	if err := json.Unmarshal([]byte(bodyRaw), &body); err != nil {
		return fmt.Errorf("unmarshalling JSON for Custom Trigger %q: %+v", id.Name, err)
	}

	log.Printf("[DEBUG] logic_custom_trigger initial body is: %s", body)

	if err := resourceLogicAppTriggerUpdate(d, meta, *workflowId, id, body, "azurerm_logic_app_trigger_custom"); err != nil {
		return err
	}

	return resourceLogicAppTriggerCustomRead(d, meta)
}

func resourceLogicAppTriggerCustomRead(d *pluginsdk.ResourceData, meta interface{}) error {
	id, err := parse.TriggerID(d.Id())
	if err != nil {
		return err
	}

	t, app, err := retrieveLogicAppTrigger(d, meta, id.ResourceGroup, id.WorkflowName, id.Name)
	if err != nil {
		return err
	}

	if t == nil {
		log.Printf("[DEBUG] Logic App %q (Resource Group %q) does not contain Trigger %q - removing from state", id.WorkflowName, id.ResourceGroup, id.Name)
		d.SetId("")
		return nil
	}

	action := *t

	d.Set("name", id.Name)
	d.Set("logic_app_id", app.ID)

	// Azure returns an additional field called evaluatedRecurrence in the trigger body which
	// is a copy of the recurrence specified in the body property and breaks the diff suppress logic
	delete(action, "evaluatedRecurrence")

	body, err := json.Marshal(action)
	if err != nil {
		return fmt.Errorf("serializing `body` for Trigger %q: %+v", id.Name, err)
	}
	log.Printf("[DEBUG] logic_custom_trigger body is: %s", string(body))

	if err := d.Set("body", string(body)); err != nil {
		return fmt.Errorf("setting `body` for Trigger %q: %+v", id.Name, err)
	}

	return nil
}

func resourceLogicAppTriggerCustomDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	id, err := parse.TriggerID(d.Id())
	if err != nil {
		return err
	}

	err = resourceLogicAppTriggerRemove(d, meta, id.ResourceGroup, id.WorkflowName, id.Name)
	if err != nil {
		return fmt.Errorf("removing Trigger %s: %+v", id, err)
	}

	return nil
}
