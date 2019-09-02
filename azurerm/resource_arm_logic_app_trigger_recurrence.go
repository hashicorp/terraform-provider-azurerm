package azurerm

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

func resourceArmLogicAppTriggerRecurrence() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmLogicAppTriggerRecurrenceCreateUpdate,
		Read:   resourceArmLogicAppTriggerRecurrenceRead,
		Update: resourceArmLogicAppTriggerRecurrenceCreateUpdate,
		Delete: resourceArmLogicAppTriggerRecurrenceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"logic_app_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"frequency": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Month",
					"Week",
					"Day",
					"Hour",
					"Minute",
					"Hour",
					"Second",
				}, false),
			},

			"interval": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func resourceArmLogicAppTriggerRecurrenceCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	trigger := map[string]interface{}{
		"recurrence": map[string]interface{}{
			"frequency": d.Get("frequency").(string),
			"interval":  d.Get("interval").(int),
		},
		"type": "Recurrence",
	}

	logicAppId := d.Get("logic_app_id").(string)
	name := d.Get("name").(string)
	if err := resourceLogicAppTriggerUpdate(d, meta, logicAppId, name, trigger, "azurerm_logic_app_trigger_recurrence"); err != nil {
		return err
	}

	return resourceArmLogicAppTriggerRecurrenceRead(d, meta)
}

func resourceArmLogicAppTriggerRecurrenceRead(d *schema.ResourceData, meta interface{}) error {
	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	logicAppName := id.Path["workflows"]
	name := id.Path["triggers"]

	t, app, err := retrieveLogicAppTrigger(meta, resourceGroup, logicAppName, name)
	if err != nil {
		return err
	}

	if t == nil {
		log.Printf("[DEBUG] Logic App %q (Resource Group %q) does not contain Trigger %q - removing from state", logicAppName, resourceGroup, name)
		d.SetId("")
		return nil
	}

	trigger := *t

	d.Set("name", name)
	d.Set("logic_app_id", app.ID)

	v := trigger["recurrence"]
	if v == nil {
		return fmt.Errorf("Error `recurrence` was nil for HTTP Trigger %q (Logic App %q / Resource Group %q)", name, logicAppName, resourceGroup)
	}

	recurrence, ok := v.(map[string]interface{})
	if !ok {
		return fmt.Errorf("Error parsing `recurrence` for HTTP Trigger %q (Logic App %q / Resource Group %q)", name, logicAppName, resourceGroup)
	}

	if frequency := recurrence["frequency"]; frequency != nil {
		d.Set("frequency", frequency.(string))
	}

	if interval := recurrence["interval"]; interval != nil {
		d.Set("interval", int(interval.(float64)))
	}

	return nil
}

func resourceArmLogicAppTriggerRecurrenceDelete(d *schema.ResourceData, meta interface{}) error {
	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	logicAppName := id.Path["workflows"]
	name := id.Path["triggers"]

	err = resourceLogicAppTriggerRemove(d, meta, resourceGroup, logicAppName, name)
	if err != nil {
		return fmt.Errorf("Error removing Trigger %q from Logic App %q (Resource Group %q): %+v", name, logicAppName, resourceGroup, err)
	}

	return nil
}
