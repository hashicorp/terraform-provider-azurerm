package azurerm

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/structure"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

func resourceArmLogicAppTriggerHttpRequest() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmLogicAppTriggerHttpRequestCreateUpdate,
		Read:   resourceArmLogicAppTriggerHttpRequestRead,
		Update: resourceArmLogicAppTriggerHttpRequestCreateUpdate,
		Delete: resourceArmLogicAppTriggerHttpRequestDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		CustomizeDiff: func(diff *schema.ResourceDiff, v interface{}) error {

			relativePath := diff.Get("relative_path").(string)
			if relativePath != "" {
				method := diff.Get("method").(string)
				if method == "" {
					return fmt.Errorf("`method` must be specified when `relative_path` is set.")
				}
			}

			return nil
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

			"schema": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     validation.ValidateJsonString,
				DiffSuppressFunc: structure.SuppressJsonDiff,
			},

			"method": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					http.MethodDelete,
					http.MethodGet,
					http.MethodPatch,
					http.MethodPost,
					http.MethodPut,
				}, false),
			},

			"relative_path": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateLogicAppTriggerHttpRequestRelativePath,
			},
		},
	}
}

func resourceArmLogicAppTriggerHttpRequestCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	schemaRaw := d.Get("schema").(string)
	var schema map[string]interface{}
	if err := json.Unmarshal([]byte(schemaRaw), &schema); err != nil {
		return fmt.Errorf("Error unmarshalling JSON from Schema: %+v", err)
	}

	inputs := map[string]interface{}{
		"schema": schema,
	}

	if v, ok := d.GetOk("method"); ok {
		inputs["method"] = v.(string)
	}

	if v, ok := d.GetOk("relative_path"); ok {
		inputs["relativePath"] = v.(string)
	}

	trigger := map[string]interface{}{
		"inputs": inputs,
		"kind":   "Http",
		"type":   "Request",
	}

	logicAppId := d.Get("logic_app_id").(string)
	name := d.Get("name").(string)
	if err := resourceLogicAppTriggerUpdate(d, meta, logicAppId, name, trigger, "azurerm_logic_app_trigger_http_request"); err != nil {
		return err
	}

	return resourceArmLogicAppTriggerHttpRequestRead(d, meta)
}

func resourceArmLogicAppTriggerHttpRequestRead(d *schema.ResourceData, meta interface{}) error {
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

	v := trigger["inputs"]
	if v == nil {
		return fmt.Errorf("Error `inputs` was nil for HTTP Trigger %q (Logic App %q / Resource Group %q)", name, logicAppName, resourceGroup)
	}

	inputs, ok := v.(map[string]interface{})
	if !ok {
		return fmt.Errorf("Error parsing `inputs` for HTTP Trigger %q (Logic App %q / Resource Group %q)", name, logicAppName, resourceGroup)
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
			return fmt.Errorf("Error serializing the Schema to JSON: %+v", err)
		}

		d.Set("schema", string(schema))
	}

	return nil
}

func resourceArmLogicAppTriggerHttpRequestDelete(d *schema.ResourceData, meta interface{}) error {
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

func validateLogicAppTriggerHttpRequestRelativePath(v interface{}, _ string) (warnings []string, errors []error) {
	value := v.(string)

	r, _ := regexp.Compile("^[A-Za-z0-9_/}{]+$")
	if !r.MatchString(value) {
		errors = append(errors, fmt.Errorf("Relative Path can only contain alphanumeric characters, underscores, forward slashes and curly braces."))
	}

	return warnings, errors
}
