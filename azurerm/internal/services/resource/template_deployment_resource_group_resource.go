package resource

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2020-06-01/resources"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/resource/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/resource/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceGroupTemplateDeploymentResource() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupTemplateDeploymentResourceCreate,
		Read:   resourceGroupTemplateDeploymentResourceRead,
		Update: resourceGroupTemplateDeploymentResourceUpdate,
		Delete: resourceGroupTemplateDeploymentResourceDelete,
		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.ResourceGroupTemplateDeploymentID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(180 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(180 * time.Minute),
			Delete: schema.DefaultTimeout(180 * time.Minute),
		},

		// (@jackofallops - lintignore needed as we need to make sure the JSON is usable in `output_content`)

		//lintignore:S033
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.TemplateDeploymentName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"deployment_mode": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(resources.Complete),
					string(resources.Incremental),
				}, false),
			},

			"template_content": {
				Type:      schema.TypeString,
				Required:  true,
				StateFunc: utils.NormalizeJson,
			},

			// Optional
			"debug_level": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(templateDeploymentDebugLevels, false),
			},

			"parameters_content": {
				Type:      schema.TypeString,
				Optional:  true,
				Computed:  true,
				StateFunc: utils.NormalizeJson,
			},

			"tags": tags.Schema(),

			// Computed
			"output_content": {
				Type:      schema.TypeString,
				Computed:  true,
				StateFunc: utils.NormalizeJson,
				// NOTE:  outputs can be strings, ints, objects etc - whilst using a nested object was considered
				// parsing the JSON using `jsondecode` allows the users to interact with/map objects as required
			},
		},
	}
}

func resourceGroupTemplateDeploymentResourceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.DeploymentsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewResourceGroupTemplateDeploymentID(d.Get("resource_group_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing Template Deployment %q (Resource Group %q): %+v", id.ResourceGroup, id.Name, err)
		}
	}
	if existing.Properties != nil {
		return tf.ImportAsExistsError("azurerm_resource_group_template_deployment", id.ID(subscriptionId))
	}

	template, err := expandTemplateDeploymentBody(d.Get("template_content").(string))
	if err != nil {
		return fmt.Errorf("expanding `template_content`: %+v", err)
	}
	deployment := resources.Deployment{
		Properties: &resources.DeploymentProperties{
			DebugSetting: expandTemplateDeploymentDebugSetting(d.Get("debug_level").(string)),
			Mode:         resources.DeploymentMode(d.Get("deployment_mode").(string)),
			Template:     template,
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if v, ok := d.GetOk("parameters_content"); ok && v != "" {
		parameters, err := expandTemplateDeploymentBody(v.(string))
		if err != nil {
			return fmt.Errorf("expanding `parameters_content`: %+v", err)
		}
		deployment.Properties.Parameters = parameters
	}

	log.Printf("[DEBUG] Running validation of Template Deployment %q (Resource Group %q)..", id.Name, id.ResourceGroup)
	if err := validateResourceGroupTemplateDeployment(ctx, id, deployment, client); err != nil {
		return fmt.Errorf("validating Template Deployment %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	log.Printf("[DEBUG] Validated Template Deployment %q (Resource Group %q)..", id.Name, id.ResourceGroup)

	log.Printf("[DEBUG] Provisioning Template Deployment %q (Resource Group %q)..", id.Name, id.ResourceGroup)
	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, deployment)
	if err != nil {
		return fmt.Errorf("creating Template Deployment %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	log.Printf("[DEBUG] Waiting for deployment of Template Deployment %q (Resource Group %q)..", id.Name, id.ResourceGroup)
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of Template Deployment %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.SetId(id.ID(subscriptionId))
	return resourceGroupTemplateDeploymentResourceRead(d, meta)
}

func resourceGroupTemplateDeploymentResourceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.DeploymentsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ResourceGroupTemplateDeploymentID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Retrieving Template Deployment %q (Resource Group %q)..", id.Name, id.ResourceGroup)
	template, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("retrieving Template Deployment %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	if template.Properties == nil {
		return fmt.Errorf("retrieving Template Deployment %q (Resource Group %q): `properties` was nil", id.Name, id.ResourceGroup)
	}

	// the API doesn't have a Patch operation, so we'll need to build one
	deployment := resources.Deployment{
		Properties: &resources.DeploymentProperties{
			DebugSetting: template.Properties.DebugSetting,
			Mode:         template.Properties.Mode,
		},
		Tags: template.Tags,
	}

	if d.HasChange("debug_level") {
		deployment.Properties.DebugSetting = expandTemplateDeploymentDebugSetting(d.Get("debug_level").(string))
	}

	if d.HasChange("deployment_mode") {
		deployment.Properties.Mode = resources.DeploymentMode(d.Get("deployment_mode").(string))
	}

	if d.HasChange("parameters_content") {
		parameters, err := expandTemplateDeploymentBody(d.Get("parameters_content").(string))
		if err != nil {
			return fmt.Errorf("expanding `parameters_content`: %+v", err)
		}
		deployment.Properties.Parameters = parameters
	}

	if d.HasChange("template_content") {
		templateContents, err := expandTemplateDeploymentBody(d.Get("template_content").(string))
		if err != nil {
			return fmt.Errorf("expanding `template_content`: %+v", err)
		}

		deployment.Properties.Template = templateContents
	} else {
		// retrieve the existing content and reuse that
		exportedTemplate, err := client.ExportTemplate(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			return fmt.Errorf("retrieving Contents for Template Deployment %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}

		deployment.Properties.Template = exportedTemplate.Template
	}

	if d.HasChange("tags") {
		deployment.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	log.Printf("[DEBUG] Running validation of Template Deployment %q (Resource Group %q)..", id.Name, id.ResourceGroup)
	if err := validateResourceGroupTemplateDeployment(ctx, *id, deployment, client); err != nil {
		return fmt.Errorf("validating Template Deployment %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	log.Printf("[DEBUG] Validated Template Deployment %q (Resource Group %q)..", id.Name, id.ResourceGroup)

	log.Printf("[DEBUG] Provisioning Template Deployment %q (Resource Group %q)..", id.Name, id.ResourceGroup)
	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, deployment)
	if err != nil {
		return fmt.Errorf("creating Template Deployment %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	log.Printf("[DEBUG] Waiting for deployment of Template Deployment %q (Resource Group %q)..", id.Name, id.ResourceGroup)
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of Template Deployment %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return resourceGroupTemplateDeploymentResourceRead(d, meta)
}

func resourceGroupTemplateDeploymentResourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.DeploymentsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ResourceGroupTemplateDeploymentID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Template Deployment %q (Resource Group %q) was not found - removing from state", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Template Deployment %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	templateContents, err := client.ExportTemplate(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("retrieving Template Content for Template Deployment %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	if props := resp.Properties; props != nil {
		d.Set("debug_level", flattenTemplateDeploymentDebugSetting(props.DebugSetting))
		d.Set("deployment_mode", string(props.Mode))

		filteredParams := filterOutTemplateDeploymentParameters(props.Parameters)
		flattenedParams, err := flattenTemplateDeploymentBody(filteredParams)
		if err != nil {
			return fmt.Errorf("flattening `parameters_content`: %+v", err)
		}
		d.Set("parameters_content", flattenedParams)

		flattenedOutputs, err := flattenTemplateDeploymentBody(props.Outputs)
		if err != nil {
			return fmt.Errorf("flattening `output_content`: %+v", err)
		}
		d.Set("output_content", flattenedOutputs)
	}

	flattenedTemplate, err := flattenTemplateDeploymentBody(templateContents.Template)
	if err != nil {
		return fmt.Errorf("flattening `template_content`: %+v", err)
	}
	d.Set("template_content", flattenedTemplate)

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceGroupTemplateDeploymentResourceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.DeploymentsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ResourceGroupTemplateDeploymentID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Retrieving Template Deployment %q (Resource Group %q)..", id.Name, id.ResourceGroup)
	template, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(template.Response) {
			return nil
		}

		return fmt.Errorf("retrieving Template Deployment %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	if template.Properties == nil {
		return fmt.Errorf("`properties` was nil for template`")
	}

	deleteItemsInTemplate := meta.(*clients.Client).Features.TemplateDeployment.DeleteNestedItemsDuringDeletion
	if deleteItemsInTemplate {
		resourceClient := meta.(*clients.Client).Resource
		log.Printf("[DEBUG] Removing items provisioned by the Template Deployment %q (Resource Group %q)..", id.Name, id.ResourceGroup)
		if err := deleteItemsProvisionedByTemplate(ctx, resourceClient, *template.Properties); err != nil {
			return fmt.Errorf("removing items provisioned by this Template Deployment: %+v", err)
		}
		log.Printf("[DEBUG] Removed items provisioned by the Template Deployment %q (Resource Group %q)..", id.Name, id.ResourceGroup)
	} else {
		log.Printf("[DEBUG] Skipping removing items provisioned by the Template Deployment %q (Resource Group %q) as the feature is disabled", id.Name, id.ResourceGroup)
	}

	log.Printf("[DEBUG] Deleting Template Deployment %q (Resource Group %q)..", id.Name, id.ResourceGroup)
	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("deleting Template Deployment %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	log.Printf("[DEBUG] Waiting for deletion of Template Deployment %q (Resource Group %q)..", id.Name, id.ResourceGroup)
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of Template Deployment %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	log.Printf("[DEBUG] Deleted Template Deployment %q (Resource Group %q).", id.Name, id.ResourceGroup)

	return nil
}

func validateResourceGroupTemplateDeployment(ctx context.Context, id parse.ResourceGroupTemplateDeploymentId, deployment resources.Deployment, client *resources.DeploymentsClient) error {
	validationFuture, err := client.Validate(ctx, id.ResourceGroup, id.Name, deployment)
	if err != nil {
		return fmt.Errorf("requesting validating: %+v", err)
	}
	if err := validationFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for validation: %+v", err)
	}
	validationResult, err := validationFuture.Result(*client)
	if err != nil {
		return fmt.Errorf("retrieving validation result: %+v", err)
	}
	if validationResult.Error != nil {
		if validationResult.Error.Message != nil {
			return fmt.Errorf("%s", *validationResult.Error.Message)
		}
		return fmt.Errorf("%+v", *validationResult.Error)
	}

	return nil
}
