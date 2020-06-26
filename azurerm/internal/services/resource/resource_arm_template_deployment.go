package resource

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/resources"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmTemplateDeployment() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmTemplateDeploymentCreateUpdate,
		Read:   resourceArmTemplateDeploymentRead,
		Update: resourceArmTemplateDeploymentCreateUpdate,
		Delete: resourceArmTemplateDeploymentDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(180 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(180 * time.Minute),
			Delete: schema.DefaultTimeout(180 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			// if resource_group_name is specified a standard ARM deployment will
			// be created, if no resource_group_name is provided, we create a
			// subscription ARM deployment.
			"resource_group_name": azure.SchemaResourceGroupNameOptional(),

			// location is marked as optional, but is required for subscription
			// ARM deployments. The API will return a error message if the location
			// is missing or invalid.
			"location": azure.SchemaLocationOptional(),

			"template_body": {
				Type:      schema.TypeString,
				Optional:  true,
				Computed:  true,
				StateFunc: azure.NormalizeJson,
			},

			"parameters": {
				Type:          schema.TypeMap,
				Optional:      true,
				ConflictsWith: []string{"parameters_body"},
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"parameters_body": {
				Type:          schema.TypeString,
				Optional:      true,
				StateFunc:     azure.NormalizeJson,
				ConflictsWith: []string{"parameters"},
			},

			"deployment_mode": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(resources.Complete),
					string(resources.Incremental),
				}, true),
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"outputs": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceArmTemplateDeploymentCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.DeploymentsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	deploymentMode := d.Get("deployment_mode").(string)

	// determine the deployment target (subscription or resource group)
	// based on the targets we use different methods from the deployments client
	var deploymentTarget string
	if resourceGroup != "" {
		deploymentTarget = "resourceGroup"
	} else {
		deploymentTarget = "subscription"
	}

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		if deploymentTarget == "resourceGroup" {
			// Resource group level deployment
			existing, err := client.Get(ctx, resourceGroup, name)
			if err != nil {
				if !utils.ResponseWasNotFound(existing.Response) {
					return fmt.Errorf("Error checking for presence of existing Template Deployment %s (resource group %s) %v", name, resourceGroup, err)
				}
			}
			if existing.ID != nil && *existing.ID != "" {
				return tf.ImportAsExistsError("azurerm_template_deployment", *existing.ID)
			}
		} else {
			// Subscription level deployment
			existing, err := client.GetAtSubscriptionScope(ctx, name)
			if err != nil {
				if !utils.ResponseWasNotFound(existing.Response) {
					return fmt.Errorf("Error checking for presence of existing Template Deployment %s %v", name, err)
				}
			}
			if existing.ID != nil && *existing.ID != "" {
				return tf.ImportAsExistsError("azurerm_template_deployment", *existing.ID)
			}
		}
	}

	log.Printf("[INFO] preparing arguments for AzureRM Template Deployment creation.")
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

	// append the location if provided. This is required for subscription
	// level deployments.
	if location, ok := d.GetOk("location"); ok {
		deployment.Location = utils.String(location.(string))
	}

	if deploymentTarget == "resourceGroup" {
		// Resource group deployment
		deploymentValidateResponse, err := client.Validate(ctx, resourceGroup, name, deployment)
		if err != nil {
			return fmt.Errorf("Error requesting Validation for Template Deployment %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
		if deploymentValidateResponse.Error != nil {
			if deploymentValidateResponse.Error.Message != nil {
				return fmt.Errorf("Error validating Template for Deployment %q (Resource Group %q): %+v", name, resourceGroup, *deploymentValidateResponse.Error.Message)
			}
			return fmt.Errorf("Error validating Template for Deployment %q (Resource Group %q): %+v", name, resourceGroup, *deploymentValidateResponse.Error)
		}
	} else {
		// Subscription level deployment
		deploymentValidateResponse, err := client.ValidateAtSubscriptionScope(ctx, name, deployment)
		if err != nil {
			return fmt.Errorf("Error requesting Validation for Template Deployment %q: %+v", name, err)
		}
		if deploymentValidateResponse.Error != nil {
			if deploymentValidateResponse.Error.Message != nil {
				return fmt.Errorf("Error validating Template for Deployment %q: %+v", name, *deploymentValidateResponse.Error.Message)
			}
			return fmt.Errorf("Error validating Template for Deployment %q: %+v", name, *deploymentValidateResponse.Error)
		}
	}

	if !d.IsNewResource() {
		d.Partial(false)
	}

	if deploymentTarget == "resourceGroup" {
		// Resource group deployment
		future, err := client.CreateOrUpdate(ctx, resourceGroup, name, deployment)
		if err != nil {
			return fmt.Errorf("Error creating deployment: %+v", err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error waiting for deployment: %+v", err)
		}
		read, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			return err
		}
		if read.ID == nil {
			return fmt.Errorf("Cannot read Template Deployment %s (resource group %s) ID", name, resourceGroup)
		}

		d.SetId(*read.ID)

	} else {
		// Subscription level deployment

		future, err := client.CreateOrUpdateAtSubscriptionScope(ctx, name, deployment)
		if err != nil {
			return fmt.Errorf("Error creating deployment: %+v", err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error waiting for deployment: %+v", err)
		}
		read, err := client.GetAtSubscriptionScope(ctx, name)
		if err != nil {
			return err
		}
		if read.ID == nil {
			return fmt.Errorf("Cannot read Template Deployment %s ID", name)
		}

		d.SetId(*read.ID)
	}

	return resourceArmTemplateDeploymentRead(d, meta)
}

func resourceArmTemplateDeploymentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.DeploymentsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	name := id.Path["deployments"]
	if name == "" {
		name = id.Path["Deployments"]
	}

	var resp resources.DeploymentExtended
	if id.ResourceGroup != "" {
		// Resource group deployment
		response, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				d.SetId("")
				return nil
			}
			return fmt.Errorf("Error making Read request on Azure RM Template Deployment %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
		// surface response
		resp = response
	} else {
		// Subscription deployment
		response, err := client.GetAtSubscriptionScope(ctx, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				d.SetId("")
				return nil
			}
			return fmt.Errorf("Error making Read request on Azure RM Template Deployment %q: %+v", name, err)
		}
		// surface response
		resp = response
	}

	outputs := make(map[string]string)
	if resp.Properties != nil {
		if outs := resp.Properties.Outputs; outs != nil {
			outsVal := outs.(map[string]interface{})
			if len(outsVal) > 0 {
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
			}
		}
	}

	return d.Set("outputs", outputs)
}

func resourceArmTemplateDeploymentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.DeploymentsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["deployments"]
	if name == "" {
		name = id.Path["Deployments"]
	}

	if id.ResourceGroup != "" {
		// Resource group deployment
		_, err = client.Delete(ctx, resourceGroup, name)
		if err != nil {
			return err
		}

	} else {
		// Subscription deployment
		_, err = client.DeleteAtSubscriptionScope(ctx, name)
		if err != nil {
			return err
		}
	}

	return waitForTemplateDeploymentToBeDeleted(ctx, client, resourceGroup, name, d)
}

// TODO: move this out into the new `helpers` structure
func expandParametersBody(body string) (map[string]interface{}, error) {
	var parametersBody map[string]interface{}
	if err := json.Unmarshal([]byte(body), &parametersBody); err != nil {
		return nil, fmt.Errorf("Error Expanding the parameters_body for Azure RM Template Deployment")
	}
	return parametersBody, nil
}

func expandTemplateBody(template string) (map[string]interface{}, error) {
	var templateBody map[string]interface{}
	if err := json.Unmarshal([]byte(template), &templateBody); err != nil {
		return nil, fmt.Errorf("Error Expanding the template_body for Azure RM Template Deployment")
	}
	return templateBody, nil
}

func waitForTemplateDeploymentToBeDeleted(ctx context.Context, client *resources.DeploymentsClient, resourceGroup, name string, d *schema.ResourceData) error {
	// we can't use the Waiter here since the API returns a 200 once it's deleted which is considered a polling status code..
	log.Printf("[DEBUG] Waiting for Template Deployment %q to be deleted", name)
	if resourceGroup != "" {
		// Resource group deployment
		stateConf := &resource.StateChangeConf{
			Pending: []string{"200"},
			Target:  []string{"404"},
			Refresh: templateDeploymentStateStatusCodeRefreshFunc(ctx, client, resourceGroup, name),
			Timeout: d.Timeout(schema.TimeoutDelete),
		}
		if _, err := stateConf.WaitForState(); err != nil {
			return fmt.Errorf("Error waiting for Template Deployment to be deleted: %q (RG: %q): %+v", name, resourceGroup, err)
		}
	} else {
		// Subscription deployment
		stateConf := &resource.StateChangeConf{
			Pending: []string{"200"},
			Target:  []string{"404"},
			Refresh: templateDeploymentStateStatusCodeRefreshFunc(ctx, client, "", name),
			Timeout: d.Timeout(schema.TimeoutDelete),
		}
		if _, err := stateConf.WaitForState(); err != nil {
			return fmt.Errorf("Error waiting for Template Deployment to be deleted: %q: %+v", name, err)
		}
	}

	return nil
}

func templateDeploymentStateStatusCodeRefreshFunc(ctx context.Context, client *resources.DeploymentsClient, resourceGroup, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		if resourceGroup != "" {
			// Resource group deployment
			res, err := client.Get(ctx, resourceGroup, name)
			if err != nil {
				if utils.ResponseWasNotFound(res.Response) {
					return res, strconv.Itoa(res.StatusCode), nil
				}
				return nil, "", fmt.Errorf("Error polling for the status of the Template Deployment %q (RG: %q): %+v", name, resourceGroup, err)
			}

			log.Printf("Retrieving Template Deployment %q (Resource Group %q) returned Status %d", resourceGroup, name, res.StatusCode)
			return res, strconv.Itoa(res.StatusCode), nil

		} else {
			// Subscription deployment
			res, err := client.GetAtSubscriptionScope(ctx, name)
			if err != nil {
				if utils.ResponseWasNotFound(res.Response) {
					return res, strconv.Itoa(res.StatusCode), nil
				}
				return nil, "", fmt.Errorf("Error polling for the status of the Template Deployment %q: %+v", name, err)
			}

			log.Printf("Retrieving Template Deployment %q returned Status %d", name, res.StatusCode)
			return res, strconv.Itoa(res.StatusCode), nil
		}
	}
}
