package azurerm

import (
	"fmt"
	"log"
	"strings"

	"net/http"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

func resourceArmLogicAppActionHTTP() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmLogicAppActionHTTPCreateUpdate,
		Read:   resourceArmLogicAppActionHTTPRead,
		Update: resourceArmLogicAppActionHTTPCreateUpdate,
		Delete: resourceArmLogicAppActionHTTPDelete,
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

			"method": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					http.MethodDelete,
					http.MethodGet,
					http.MethodPatch,
					http.MethodPost,
					http.MethodPut,
				}, false),
			},

			"uri": {
				Type:     schema.TypeString,
				Required: true,
			},

			"body": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"headers": {
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
	}
}

func resourceArmLogicAppActionHTTPCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	headersRaw := d.Get("headers").(map[string]interface{})
	headers, err := expandLogicAppActionHttpHeaders(headersRaw)
	if err != nil {
		return err
	}

	inputs := map[string]interface{}{
		"method":  d.Get("method").(string),
		"uri":     d.Get("uri").(string),
		"headers": headers,
	}

	if v, ok := d.GetOk("body"); ok {
		inputs["body"] = v.(string)
	}

	action := map[string]interface{}{
		"inputs": inputs,
		"type":   "http",
	}

	logicAppId := d.Get("logic_app_id").(string)
	name := d.Get("name").(string)
	err = resourceLogicAppActionUpdate(d, meta, logicAppId, name, action, "azurerm_logic_app_action_http")
	if err != nil {
		return err
	}

	return resourceArmLogicAppActionHTTPRead(d, meta)
}

func resourceArmLogicAppActionHTTPRead(d *schema.ResourceData, meta interface{}) error {
	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	logicAppName := id.Path["workflows"]
	name := id.Path["actions"]

	t, app, err := retrieveLogicAppAction(meta, resourceGroup, logicAppName, name)
	if err != nil {
		return err
	}

	if t == nil {
		log.Printf("[DEBUG] Logic App %q (Resource Group %q) does not contain Action %q - removing from state", logicAppName, resourceGroup, name)
		d.SetId("")
		return nil
	}

	action := *t

	d.Set("name", name)
	d.Set("logic_app_id", app.ID)

	actionType := action["type"].(string)
	if !strings.EqualFold(actionType, "http") {
		return fmt.Errorf("Expected an HTTP Action for Action %q (Logic App %q / Resource Group %q) - got %q", name, logicAppName, resourceGroup, actionType)
	}

	v := action["inputs"]
	if v == nil {
		return fmt.Errorf("Error`inputs` was nil for HTTP Action %q (Logic App %q / Resource Group %q)", name, logicAppName, resourceGroup)
	}

	inputs, ok := v.(map[string]interface{})
	if !ok {
		return fmt.Errorf("Error parsing `inputs` for HTTP Action %q (Logic App %q / Resource Group %q)", name, logicAppName, resourceGroup)
	}

	if uri := inputs["uri"]; uri != nil {
		d.Set("uri", uri.(string))
	}

	if method := inputs["method"]; method != nil {
		d.Set("method", method.(string))
	}

	if body := inputs["body"]; body != nil {
		d.Set("body", body.(string))
	}

	if headers := inputs["headers"]; headers != nil {
		hv := headers.(map[string]interface{})
		if err := d.Set("headers", hv); err != nil {
			return fmt.Errorf("Error setting `headers` for HTTP Action %q: %+v", name, err)
		}
	}

	return nil
}

func resourceArmLogicAppActionHTTPDelete(d *schema.ResourceData, meta interface{}) error {
	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	logicAppName := id.Path["workflows"]
	name := id.Path["actions"]

	err = resourceLogicAppActionRemove(d, meta, resourceGroup, logicAppName, name)
	if err != nil {
		return fmt.Errorf("Error removing Action %q from Logic App %q (Resource Group %q): %+v", name, logicAppName, resourceGroup, err)
	}

	return nil
}

func expandLogicAppActionHttpHeaders(headersRaw map[string]interface{}) (*map[string]string, error) {
	headers := make(map[string]string)

	for i, v := range headersRaw {
		value, err := tagValueToString(v)
		if err != nil {
			return nil, err
		}

		headers[i] = value
	}

	return &headers, nil
}
