package logic

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

func resourceLogicAppTriggerRecurrence() *schema.Resource {
	return &schema.Resource{
		Create: resourceLogicAppTriggerRecurrenceCreateUpdate,
		Read:   resourceLogicAppTriggerRecurrenceRead,
		Update: resourceLogicAppTriggerRecurrenceCreateUpdate,
		Delete: resourceLogicAppTriggerRecurrenceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
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

			"start_time": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsRFC3339Time,
			},

			"time_zone": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateLogicAppTriggerRecurrenceTimeZone(),
			},
		},
	}
}

func resourceLogicAppTriggerRecurrenceCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	trigger := map[string]interface{}{
		"recurrence": map[string]interface{}{
			"frequency": d.Get("frequency").(string),
			"interval":  d.Get("interval").(int),
		},
		"type": "Recurrence",
	}

	if v, ok := d.GetOk("start_time"); ok {
		trigger["recurrence"].(map[string]interface{})["startTime"] = v.(string)

		// time_zone only allowed when start_time is specified
		if v, ok := d.GetOk("time_zone"); ok {
			trigger["recurrence"].(map[string]interface{})["timeZone"] = v.(string)
		}
	}

	logicAppId := d.Get("logic_app_id").(string)
	name := d.Get("name").(string)
	if err := resourceLogicAppTriggerUpdate(d, meta, logicAppId, name, trigger, "azurerm_logic_app_trigger_recurrence"); err != nil {
		return err
	}

	return resourceLogicAppTriggerRecurrenceRead(d, meta)
}

func resourceLogicAppTriggerRecurrenceRead(d *schema.ResourceData, meta interface{}) error {
	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	logicAppName := id.Path["workflows"]
	name := id.Path["triggers"]

	t, app, err := retrieveLogicAppTrigger(d, meta, resourceGroup, logicAppName, name)
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

	if startTime := recurrence["startTime"]; startTime != nil {
		d.Set("start_time", startTime.(string))
	}

	if timeZone := recurrence["timeZone"]; timeZone != nil {
		d.Set("time_zone", timeZone.(string))
	}

	return nil
}

func resourceLogicAppTriggerRecurrenceDelete(d *schema.ResourceData, meta interface{}) error {
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

func validateLogicAppTriggerRecurrenceTimeZone() schema.SchemaValidateFunc {
	// from https://support.microsoft.com/en-us/help/973627/microsoft-time-zone-index-values
	timeZones := []string{
		"Dateline Standard Time",
		"Samoa Standard Time",
		"Hawaiian Standard Time",
		"Alaskan Standard Time",
		"Pacific Standard Time",
		"Mountain Standard Time",
		"Mexico Standard Time",
		"US Mountain Standard Time",
		"Central Standard Time",
		"Canada Central Standard Time",
		"Mexico Standard Time",
		"Central America Standard Time",
		"Eastern Standard Time",
		"US Eastern Standard Time",
		"SA Pacific Standard Time",
		"Atlantic Standard Time",
		"SA Western Standard Time",
		"Pacific SA Standard Time",
		"Newfoundland and Labrador Standard Time",
		"E South America Standard Time",
		"SA Eastern Standard Time",
		"Greenland Standard Time",
		"Mid-Atlantic Standard Time",
		"Azores Standard Time",
		"Cape Verde Standard Time",
		"GMT Standard Time",
		"Greenwich Standard Time",
		"Central Europe Standard Time",
		"Central European Standard Time",
		"Romance Standard Time",
		"W Europe Standard Time",
		"W Central Africa Standard Time",
		"E Europe Standard Time",
		"Egypt Standard Time",
		"FLE Standard Time",
		"GTB Standard Time",
		"Israel Standard Time",
		"South Africa Standard Time",
		"Russian Standard Time",
		"Arab Standard Time",
		"E Africa Standard Time",
		"Arabic Standard Time",
		"Iran Standard Time",
		"Arabian Standard Time",
		"Caucasus Standard Time",
		"Transitional Islamic State of Afghanistan Standard Time",
		"Ekaterinburg Standard Time",
		"West Asia Standard Time",
		"India Standard Time",
		"Nepal Standard Time",
		"Central Asia Standard Time",
		"Sri Lanka Standard Time",
		"N Central Asia Standard Time",
		"Myanmar Standard Time",
		"SE Asia Standard Time",
		"North Asia Standard Time",
		"China Standard Time",
		"Singapore Standard Time",
		"Taipei Standard Time",
		"W Australia Standard Time",
		"North Asia East Standard Time",
		"Korea Standard Time",
		"Tokyo Standard Time",
		"Yakutsk Standard Time",
		"AUS Central Standard Time",
		"Cen Australia Standard Time",
		"AUS Eastern Standard Time",
		"E Australia Standard Time",
		"Tasmania Standard Time",
		"Vladivostok Standard Time",
		"West Pacific Standard Time",
		"Central Pacific Standard Time",
		"Fiji Islands Standard Time",
		"New Zealand Standard Time",
		"Tonga Standard Time",
		"Azerbaijan Standard Time",
		"Middle East Standard Time",
		"Jordan Standard Time",
		"Central Standard Time (Mexico)",
		"Mountain Standard Time (Mexico)",
		"Pacific Standard Time (Mexico)",
		"Namibia Standard Time",
		"Georgian Standard Time",
		"Central Brazilian Standard Time",
		"Montevideo Standard Time",
		"Armenian Standard Time",
		"Venezuela Standard Time",
		"Argentina Standard Time",
		"Morocco Standard Time",
		"Pakistan Standard Time",
		"Mauritius Standard Time",
		"UTC",
		"Paraguay Standard Time",
		"Kamchatka Standard Time",
	}
	return validation.StringInSlice(timeZones, false)
}
