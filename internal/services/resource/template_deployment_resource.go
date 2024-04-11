// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2020-06-01/resources" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceTemplateDeployment() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceTemplateDeploymentCreateUpdate,
		Read:   resourceTemplateDeploymentRead,
		Update: resourceTemplateDeploymentCreateUpdate,
		Delete: resourceTemplateDeploymentDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(180 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(180 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(180 * time.Minute),
		},

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.TemplateDeploymentV0ToV1{},
		}),

		DeprecationMessage: "The resource 'azurerm_template_deployment' has been superseded by the 'azurerm_resource_group_template_deployment' resource and will be removed in v4.0 of the AzureRM Provider.",

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"template_body": {
				Type:      pluginsdk.TypeString,
				Optional:  true,
				Computed:  true,
				StateFunc: utils.NormalizeJson,
			},

			"parameters": {
				Type:          pluginsdk.TypeMap,
				Optional:      true,
				ConflictsWith: []string{"parameters_body"},
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"parameters_body": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				StateFunc:     utils.NormalizeJson,
				ConflictsWith: []string{"parameters"},
			},

			"deployment_mode": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(resources.DeploymentModeComplete),
					string(resources.DeploymentModeIncremental),
				}, false),
			},

			"outputs": {
				Type:     pluginsdk.TypeMap,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},
		},
	}
}

func resourceTemplateDeploymentCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.DeploymentsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewResourceGroupTemplateDeploymentID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.DeploymentName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_template_deployment", id.ID())
		}
	}

	deployment := resources.Deployment{
		Properties: &resources.DeploymentProperties{
			Mode: resources.DeploymentMode(d.Get("deployment_mode").(string)),
		},
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

		deployment.Properties.Parameters = &newParams
	}

	if v, ok := d.GetOk("parameters_body"); ok {
		params, err := expandParametersBody(v.(string))
		if err != nil {
			return err
		}

		deployment.Properties.Parameters = &params
	}

	if v, ok := d.GetOk("template_body"); ok {
		template, err := expandTemplateBody(v.(string))
		if err != nil {
			return err
		}

		deployment.Properties.Template = &template
	}

	if !d.IsNewResource() {
		d.Partial(true)
	}

	deploymentValidateFuture, err := client.Validate(ctx, id.ResourceGroup, id.DeploymentName, deployment)
	if err != nil {
		return fmt.Errorf("requesting Validation of %s: %+v", id, err)
	}

	if err := deploymentValidateFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for Validation of %s: %+v", id, err)
	}
	validationResult, err := deploymentValidateFuture.Result(*client)
	if err != nil {
		return fmt.Errorf("retrieving Validation Result for %s: %+v", id, err)
	}

	if validationResult.Error != nil {
		if validationResult.Error.Message != nil {
			return fmt.Errorf("validating %s for Deployment: %+v", id, *validationResult.Error.Message)
		}
		return fmt.Errorf("validating %s for Deployment: %+v", id, *validationResult.Error)
	}

	if !d.IsNewResource() {
		d.Partial(false)
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.DeploymentName, deployment)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceTemplateDeploymentRead(d, meta)
}

func resourceTemplateDeploymentRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.DeploymentsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ResourceGroupTemplateDeploymentID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.DeploymentName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	outputs := make(map[string]string)
	if props := resp.Properties; props != nil {
		if outs := props.Outputs; outs != nil {
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
		d.Set("outputs", outputs)
	}

	return nil
}

func resourceTemplateDeploymentDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.DeploymentsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ResourceGroupTemplateDeploymentID(d.Id())
	if err != nil {
		return err
	}
	future, err := client.Delete(ctx, id.ResourceGroup, id.DeploymentName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %q: %+v", id, err)
	}

	if err := waitForTemplateDeploymentToBeDeleted(ctx, client, *id); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", id, err)
	}

	return nil
}

// TODO: move this out into the new `helpers` structure
func expandParametersBody(body string) (map[string]interface{}, error) {
	var parametersBody map[string]interface{}
	if err := json.Unmarshal([]byte(body), &parametersBody); err != nil {
		return nil, fmt.Errorf("Expanding the parameters_body for Azure RM Template Deployment")
	}
	return parametersBody, nil
}

func expandTemplateBody(template string) (map[string]interface{}, error) {
	var templateBody map[string]interface{}
	if err := json.Unmarshal([]byte(template), &templateBody); err != nil {
		return nil, fmt.Errorf("Expanding the template_body for Azure RM Template Deployment")
	}
	return templateBody, nil
}

func waitForTemplateDeploymentToBeDeleted(ctx context.Context, client *resources.DeploymentsClient, id parse.ResourceGroupTemplateDeploymentId) error {
	// we can't use the Waiter here since the API returns a 200 once it's deleted which is considered a polling status code..
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("internal-error: context had no deadline")
	}

	log.Printf("[DEBUG] Waiting for %s to be deleted", id)
	stateConf := &pluginsdk.StateChangeConf{
		Pending: []string{"200"},
		Target:  []string{"404"},
		Refresh: templateDeploymentStateStatusCodeRefreshFunc(ctx, client, id),
		Timeout: time.Until(deadline),
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to be deleted: %+v", id, err)
	}

	return nil
}

func templateDeploymentStateStatusCodeRefreshFunc(ctx context.Context, client *resources.DeploymentsClient, id parse.ResourceGroupTemplateDeploymentId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id.ResourceGroup, id.DeploymentName)

		log.Printf("Retrieving %s returned Status %d", id, res.StatusCode)

		if err != nil {
			if utils.ResponseWasNotFound(res.Response) {
				return res, strconv.Itoa(res.StatusCode), nil
			}
			return nil, "", fmt.Errorf("polling for the status of %s: %+v", id, err)
		}

		return res, strconv.Itoa(res.StatusCode), nil
	}
}
