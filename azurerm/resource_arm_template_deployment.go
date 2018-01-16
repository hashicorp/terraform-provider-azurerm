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
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmTemplateDeployment() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmTemplateDeploymentCreate,
		Read:   resourceArmTemplateDeploymentRead,
		Update: resourceArmTemplateDeploymentCreate,
		Delete: resourceArmTemplateDeploymentDelete,

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
				Type:     schema.TypeMap,
				Optional: true,
			},

			"outputs": {
				Type:     schema.TypeMap,
				Computed: true,
			},

			"deployment_mode": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceArmTemplateDeploymentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient)
	deployClient := client.deploymentsClient
	ctx := client.StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	deploymentMode := d.Get("deployment_mode").(string)

	log.Printf("[INFO] preparing arguments for Azure ARM Template Deployment creation.")
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

	future, err := deployClient.CreateOrUpdate(ctx, resourceGroup, name, deployment)
	if err != nil {
		return fmt.Errorf("Error creating deployment: %+v", err)
	}

	err = future.WaitForCompletion(ctx, deployClient.Client)
	if err != nil {
		return fmt.Errorf("Error creating deployment: %+v", err)
	}

	read, err := deployClient.Get(ctx, resourceGroup, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Template Deployment %s (resource group %s) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	// TODO: is this even needed anymore?
	log.Printf("[DEBUG] Waiting for Template Deployment (%q in Resource Group %q) to become available", name, resourceGroup)
	stateConf := &resource.StateChangeConf{
		Pending: []string{"creating", "updating", "accepted", "running"},
		Target:  []string{"succeeded"},
		Refresh: templateDeploymentStateRefreshFunc(client, ctx, resourceGroup, name),
		Timeout: 40 * time.Minute,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Template Deployment (%q in Resource Group %q) to become available: %+v", name, resourceGroup, err)
	}

	return resourceArmTemplateDeploymentRead(d, meta)
}

func resourceArmTemplateDeploymentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient)
	deployClient := client.deploymentsClient
	ctx := client.StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["deployments"]
	if name == "" {
		name = id.Path["Deployments"]
	}

	resp, err := deployClient.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure RM Template Deployment %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	var outputs map[string]string
	if resp.Properties.Outputs != nil && len(*resp.Properties.Outputs) > 0 {
		outputs = make(map[string]string)
		for key, output := range *resp.Properties.Outputs {
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
	}

	return d.Set("outputs", outputs)
}

func resourceArmTemplateDeploymentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient)
	deployClient := client.deploymentsClient
	ctx := client.StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["deployments"]
	if name == "" {
		name = id.Path["Deployments"]
	}

	future, err := deployClient.Delete(ctx, resourceGroup, name)
	if err != nil {
		return err
	}

	err = future.WaitForCompletion(ctx, deployClient.Client)
	if err != nil {
		return err
	}

	return nil
}

// TODO: move this out into the new `helpers` structure
func expandTemplateBody(template string) (map[string]interface{}, error) {
	var templateBody map[string]interface{}
	err := json.Unmarshal([]byte(template), &templateBody)
	if err != nil {
		return nil, fmt.Errorf("Error Expanding the template_body for Azure RM Template Deployment")
	}
	return templateBody, nil
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

func templateDeploymentStateRefreshFunc(client *ArmClient, ctx context.Context, resourceGroupName string, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.deploymentsClient.Get(ctx, resourceGroupName, name)
		if err != nil {
			return nil, "", fmt.Errorf("Error issuing read request in templateDeploymentStateRefreshFunc to Azure ARM for Template Deployment %q (RG: %q): %+v", name, resourceGroupName, err)
		}

		return res, strings.ToLower(*res.Properties.ProvisioningState), nil
	}
}
