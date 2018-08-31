package azurerm

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2017-05-10/resources"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmTemplateDeployment() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmTemplateDeploymentCreateUpdate,
		Read:   resourceArmTemplateDeploymentRead,
		Update: resourceArmTemplateDeploymentCreateUpdate,
		Delete: resourceArmTemplateDeploymentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(time.Minute * 60),
			Update: schema.DefaultTimeout(time.Minute * 60),
			Delete: schema.DefaultTimeout(time.Minute * 60),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": resourceGroupNameSchema(),

			"template_body": {
				Type:      schema.TypeString,
				Optional:  true,
				Computed:  true,
				StateFunc: normalizeJson,
			},

			"parameters": {
				Type:          schema.TypeMap,
				Optional:      true,
				ConflictsWith: []string{"parameters_body"},
			},

			"parameters_body": {
				Type:          schema.TypeString,
				Optional:      true,
				StateFunc:     normalizeJson,
				ConflictsWith: []string{"parameters"},
			},

			"deployment_mode": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(resources.Complete),
					string(resources.Incremental),
				}, true),
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
			},

			"outputs": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func resourceArmTemplateDeploymentCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).deploymentsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		// first check if there's one in this subscription requiring import
		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Error checking for the existence of Template Deployment %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}

		if resp.ID != nil {
			return tf.ImportAsExistsError("azurerm_template_deployment", *resp.ID)
		}
	}

	log.Printf("[INFO] preparing arguments for AzureRM Template Deployment creation.")
	deploymentMode := d.Get("deployment_mode").(string)
	properties := resources.DeploymentProperties{
		Mode: resources.DeploymentMode(deploymentMode),
	}

	if v, ok := d.GetOk("parameters"); ok {
		params := v.(map[string]interface{})

		newParams := make(map[string]interface{}, len(params))
		for key, val := range params {
			newParams[key] = struct {
				Value interface{}
			}{
				Value: val,
			}
		}

		properties.Parameters = &newParams
	}

	if v, ok := d.GetOk("parameters_body"); ok {
		params, err := expandParametersBody(v.(string))
		if err != nil {
			return err
		}

		properties.Parameters = &params
	}

	if v, ok := d.GetOk("template_body"); ok {
		template, err := expandTemplateBody(v.(string))
		if err != nil {
			return err
		}

		properties.Template = &template
	}

	deployment := resources.Deployment{
		Properties: &properties,
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, deployment)
	if err != nil {
		return fmt.Errorf("Error creating deployment: %+v", err)
	}

	waitCtx, cancel := context.WithTimeout(ctx, d.Timeout(tf.TimeoutForCreateUpdate(d)))
	defer cancel()
	err = future.WaitForCompletionRef(waitCtx, client.Client)
	if err != nil {
		return fmt.Errorf("Error creating deployment: %+v", err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Template Deployment %s (resource group %s) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmTemplateDeploymentRead(d, meta)
}

func resourceArmTemplateDeploymentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).deploymentsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["deployments"]
	if name == "" {
		name = id.Path["Deployments"]
	}

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure RM Template Deployment %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	templ, err := client.ExportTemplate(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error exporting ARM Template for Template Deployment %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)

	if props := resp.Properties; props != nil {
		outputs := flattenTemplateDeploymentOutputs(props.Outputs)
		if err := d.Set("outputs", outputs); err != nil {
			return fmt.Errorf("Error setting `outputs`: %+v", err)
		}

		d.Set("deployment_mode", string(props.Mode))
	}

	template, err := flattenTemplateBody(templ.Template)
	if err != nil {
		return fmt.Errorf("Error flattening `template_body`: %s", err)
	}
	d.Set("template_body", template)

	return nil
}

func flattenTemplateDeploymentOutputs(input interface{}) map[string]string {
	outputs := make(map[string]string, 0)
	if input == nil {
		return outputs
	}

	outsVal := input.(map[string]interface{})
	for key, output := range outsVal {
		log.Printf("[DEBUG] Processing deployment output %s", key)
		outputMap := output.(map[string]interface{})
		outputValue, ok := outputMap["value"]
		if !ok {
			log.Printf("[DEBUG] No value - skipping")
			continue
		}
		outputType, ok := outputMap["type"]
		if !ok {
			log.Printf("[DEBUG] No type - skipping")
			continue
		}

		var outputValueString string
		switch strings.ToLower(outputType.(string)) {
		case "bool":
			outputValueString = strconv.FormatBool(outputValue.(bool))

		case "string":
			outputValueString = outputValue.(string)

		case "int":
			outputValueString = fmt.Sprint(outputValue)

		default:
			log.Printf("[WARN] Ignoring output %s: Outputs of type %s are not currently supported in azurerm_template_deployment.",
				key, outputType)
			continue
		}
		outputs[key] = outputValueString
	}

	return outputs
}

func resourceArmTemplateDeploymentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).deploymentsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["deployments"]
	if name == "" {
		name = id.Path["Deployments"]
	}

	_, err = client.Delete(ctx, resourceGroup, name)
	if err != nil {
		return err
	}

	timeout := d.Timeout(schema.TimeoutDelete)
	return waitForTemplateDeploymentToBeDeleted(ctx, client, resourceGroup, name, timeout)
}

// TODO: move this out into the new `helpers` structure
func expandParametersBody(body string) (map[string]interface{}, error) {
	var parametersBody map[string]interface{}
	err := json.Unmarshal([]byte(body), &parametersBody)
	if err != nil {
		return nil, fmt.Errorf("Error Expanding the parameters_body for Azure RM Template Deployment")
	}
	return parametersBody, nil
}

func expandTemplateBody(template string) (map[string]interface{}, error) {
	var templateBody map[string]interface{}
	err := json.Unmarshal([]byte(template), &templateBody)
	if err != nil {
		return nil, fmt.Errorf("Error Expanding the template_body for Azure RM Template Deployment")
	}
	return templateBody, nil
}

func flattenTemplateBody(input interface{}) (*string, error) {
	if input == nil {
		return nil, nil
	}

	v := input.(map[string]interface{})
	bytes, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("Error marshalling the ARM Template: %+v", err)
	}

	str := string(bytes)
	return &str, nil
}

func normalizeJson(jsonString interface{}) string {
	if jsonString == nil || jsonString == "" {
		return ""
	}
	var j interface{}
	err := json.Unmarshal([]byte(jsonString.(string)), &j)
	if err != nil {
		return fmt.Sprintf("Error parsing JSON: %+v", err)
	}
	b, _ := json.Marshal(j)
	return string(b[:])
}

func waitForTemplateDeploymentToBeDeleted(ctx context.Context, client resources.DeploymentsClient, resourceGroup, name string, timeout time.Duration) error {
	waitCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	// we can't use the Waiter here since the API returns a 200 once it's deleted which is considered a polling status code..
	log.Printf("[DEBUG] Waiting for Template Deployment (%q in Resource Group %q) to be deleted", name, resourceGroup)
	stateConf := &resource.StateChangeConf{
		Pending: []string{"200"},
		Target:  []string{"404"},
		Refresh: templateDeploymentStateStatusCodeRefreshFunc(waitCtx, client, resourceGroup, name),
		Timeout: timeout,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Template Deployment (%q in Resource Group %q) to be deleted: %+v", name, resourceGroup, err)
	}

	return nil
}

func templateDeploymentStateStatusCodeRefreshFunc(ctx context.Context, client resources.DeploymentsClient, resourceGroup, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, resourceGroup, name)

		log.Printf("Retrieving Template Deployment %q (Resource Group %q) returned Status %d", resourceGroup, name, res.StatusCode)

		if err != nil {
			if utils.ResponseWasNotFound(res.Response) {
				return res, strconv.Itoa(res.StatusCode), nil
			}
			return nil, "", fmt.Errorf("Error polling for the status of the Template Deployment %q (RG: %q): %+v", name, resourceGroup, err)
		}

		return res, strconv.Itoa(res.StatusCode), nil
	}
}
