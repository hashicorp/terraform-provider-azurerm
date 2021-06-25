package logic

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/logic/mgmt/2019-05-01/logic"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
)

func resourceLogicAppActionHTTP() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLogicAppActionHTTPCreateUpdate,
		Read:   resourceLogicAppActionHTTPRead,
		Update: resourceLogicAppActionHTTPCreateUpdate,
		Delete: resourceLogicAppActionHTTPDelete,
		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

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

			"method": {
				Type:     pluginsdk.TypeString,
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
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"body": {
				Type:             pluginsdk.TypeString,
				Optional:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
			},

			"headers": {
				Type:     pluginsdk.TypeMap,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},
			"run_after": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				MinItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"action_name": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
						"action_result": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(logic.WorkflowStatusSucceeded),
								string(logic.WorkflowStatusFailed),
								string(logic.WorkflowStatusSkipped),
								string(logic.WorkflowStatusTimedOut),
							}, false),
						},
					},
				},
			},
		},
	}
}

func resourceLogicAppActionHTTPCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	logicAppId := d.Get("logic_app_id").(string)
	name := d.Get("name").(string)

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

	// storing action's body in json object to keep consistent with azure portal
	if bodyRaw, ok := d.GetOk("body"); ok {
		var body map[string]interface{}
		if err := json.Unmarshal([]byte(bodyRaw.(string)), &body); err != nil {
			return fmt.Errorf("error unmarshalling JSON for Action %q: %+v", name, err)
		}
		inputs["body"] = body
	}

	action := map[string]interface{}{
		"inputs": inputs,
		"type":   "http",
	}

	if v, ok := d.GetOk("run_after"); ok {
		action["runAfter"] = expandLogicAppActionRunAfter(v.(*pluginsdk.Set).List())
	}

	err = resourceLogicAppActionUpdate(d, meta, logicAppId, name, action, "azurerm_logic_app_action_http")
	if err != nil {
		return err
	}

	return resourceLogicAppActionHTTPRead(d, meta)
}

func resourceLogicAppActionHTTPRead(d *pluginsdk.ResourceData, meta interface{}) error {
	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	logicAppName := id.Path["workflows"]
	name := id.Path["actions"]

	t, app, err := retrieveLogicAppAction(d, meta, resourceGroup, logicAppName, name)
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
		// TODO: remove in 3.0, this is preserved for backward compatibility
		if v, ok := body.(string); ok {
			d.Set("body", v)
		} else {
			// if user edit workflow in portal, the body becomes json object
			v, err := json.Marshal(body)
			if err != nil {
				return fmt.Errorf("error serializing `body` for Action %q: %+v", name, err)
			}
			d.Set("body", string(v))
		}
	}

	if headers := inputs["headers"]; headers != nil {
		hv := headers.(map[string]interface{})
		if err := d.Set("headers", hv); err != nil {
			return fmt.Errorf("Error setting `headers` for HTTP Action %q: %+v", name, err)
		}
	}

	v = action["runAfter"]
	if v != nil {
		runAfter, ok := v.(map[string]interface{})
		if !ok {
			return fmt.Errorf("Error parsing `runAfter` for HTTP Action %q (Logic App %q / Resource Group %q)", name, logicAppName, resourceGroup)
		}
		if err := d.Set("run_after", flattenLogicAppActionRunAfter(runAfter)); err != nil {
			return fmt.Errorf("Error setting `runAfter` for HTTP Action %q: %+v", name, err)
		}
	}

	return nil
}

func resourceLogicAppActionHTTPDelete(d *pluginsdk.ResourceData, meta interface{}) error {
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
		value, err := tags.TagValueToString(v)
		if err != nil {
			return nil, err
		}

		headers[i] = value
	}

	return &headers, nil
}
