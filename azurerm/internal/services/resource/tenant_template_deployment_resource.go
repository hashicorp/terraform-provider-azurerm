package resource

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2020-06-01/resources"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/resource/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/resource/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func tenantTemplateDeploymentResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: tenantTemplateDeploymentResourceCreate,
		Read:   tenantTemplateDeploymentResourceRead,
		Update: tenantTemplateDeploymentResourceUpdate,
		Delete: tenantTemplateDeploymentResourceDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.TenantTemplateDeploymentID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(180 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(180 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(180 * time.Minute),
		},

		// (@jackofallops - lintignore needed as we need to make sure the JSON is usable in `output_content`)

		//lintignore:S033
		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.TemplateDeploymentName,
			},

			"location": location.Schema(),

			"template_content": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
				ExactlyOneOf: []string{
					"template_content",
					"template_spec_version_id",
				},
				StateFunc: utils.NormalizeJson,
			},

			"template_spec_version_id": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ExactlyOneOf: []string{
					"template_content",
					"template_spec_version_id",
				},
				ValidateFunc: validate.TemplateSpecVersionID,
			},

			// Optional
			"debug_level": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(templateDeploymentDebugLevels, false),
			},

			"parameters_content": {
				Type:      pluginsdk.TypeString,
				Optional:  true,
				Computed:  true,
				StateFunc: utils.NormalizeJson,
			},

			"tags": tags.Schema(),

			// Computed
			"output_content": {
				Type:     pluginsdk.TypeString,
				Computed: true,
				// NOTE:  outputs can be strings, ints, objects etc - whilst using a nested object was considered
				// parsing the JSON using `jsondecode` allows the users to interact with/map objects as required
			},
		},
	}
}

func tenantTemplateDeploymentResourceCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.DeploymentsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewTenantTemplateDeploymentID(d.Get("name").(string))

	existing, err := client.GetAtTenantScope(ctx, id.DeploymentName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing Tenant Template Deployment %q: %+v", id.DeploymentName, err)
		}
	}
	if existing.Properties != nil {
		return tf.ImportAsExistsError("azurerm_tenant_template_deployment", id.ID())
	}

	deployment := resources.ScopedDeployment{
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		Properties: &resources.DeploymentProperties{
			DebugSetting: expandTemplateDeploymentDebugSetting(d.Get("debug_level").(string)),
			Mode:         resources.DeploymentModeIncremental,
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if templateRaw, ok := d.GetOk("template_content"); ok {
		template, err := expandTemplateDeploymentBody(templateRaw.(string))
		if err != nil {
			return fmt.Errorf("expanding `template_content`: %+v", err)
		}
		deployment.Properties.Template = template
	}

	if templateSpecVersionID, ok := d.GetOk("template_spec_version_id"); ok {
		deployment.Properties.TemplateLink = &resources.TemplateLink{
			ID: utils.String(templateSpecVersionID.(string)),
		}
	}

	if v, ok := d.GetOk("parameters_content"); ok && v != "" {
		parameters, err := expandTemplateDeploymentBody(v.(string))
		if err != nil {
			return fmt.Errorf("expanding `parameters_content`: %+v", err)
		}
		deployment.Properties.Parameters = parameters
	}

	log.Printf("[DEBUG] Running validation of Tenant Template Deployment %q..", id.DeploymentName)
	if err := validateTenantTemplateDeployment(ctx, id, deployment, client); err != nil {
		return fmt.Errorf("validating Tenant Template Deployment %q: %+v", id.DeploymentName, err)
	}
	log.Printf("[DEBUG] Validated Tenant Template Deployment %q..", id.DeploymentName)

	log.Printf("[DEBUG] Provisioning Tenant Template Deployment %q..", id.DeploymentName)
	future, err := client.CreateOrUpdateAtTenantScope(ctx, id.DeploymentName, deployment)
	if err != nil {
		return fmt.Errorf("creating Tenant Template Deployment %q: %+v", id.DeploymentName, err)
	}

	log.Printf("[DEBUG] Waiting for deployment of Tenant Template Deployment %q..", id.DeploymentName)
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of Tenant Template Deployment %q: %+v", id.DeploymentName, err)
	}

	d.SetId(id.ID())
	return tenantTemplateDeploymentResourceRead(d, meta)
}

func tenantTemplateDeploymentResourceUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.DeploymentsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.TenantTemplateDeploymentID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Retrieving Tenant Template Deployment %q..", id.DeploymentName)
	template, err := client.GetAtTenantScope(ctx, id.DeploymentName)
	if err != nil {
		return fmt.Errorf("retrieving Tenant Template Deployment %q: %+v", id.DeploymentName, err)
	}
	if template.Properties == nil {
		return fmt.Errorf("retrieving Tenant Template Deployment %q: `properties` was nil", id.DeploymentName)
	}

	// the API doesn't have a Patch operation, so we'll need to build one
	deployment := resources.ScopedDeployment{
		Location: template.Location,
		Properties: &resources.DeploymentProperties{
			DebugSetting: template.Properties.DebugSetting,
			Mode:         resources.DeploymentModeIncremental,
		},
		Tags: template.Tags,
	}

	if d.HasChange("debug_level") {
		deployment.Properties.DebugSetting = expandTemplateDeploymentDebugSetting(d.Get("debug_level").(string))
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
		exportedTemplate, err := client.ExportTemplateAtTenantScope(ctx, id.DeploymentName)
		if err != nil {
			return fmt.Errorf("retrieving Contents for Tenant Template Deployment %q: %+v", id.DeploymentName, err)
		}

		deployment.Properties.Template = exportedTemplate.Template
	}

	if d.HasChange("template_spec_version_id") {
		deployment.Properties.TemplateLink = &resources.TemplateLink{
			ID: utils.String(d.Get("template_spec_version_id").(string)),
		}
	}

	if d.HasChange("tags") {
		deployment.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	log.Printf("[DEBUG] Running validation of Tenant Template Deployment %q..", id.DeploymentName)
	if err := validateTenantTemplateDeployment(ctx, *id, deployment, client); err != nil {
		return fmt.Errorf("validating Tenant Template Deployment %q: %+v", id.DeploymentName, err)
	}
	log.Printf("[DEBUG] Validated Tenant Template Deployment %q..", id.DeploymentName)

	log.Printf("[DEBUG] Provisioning Tenant Template Deployment %q)..", id.DeploymentName)
	future, err := client.CreateOrUpdateAtTenantScope(ctx, id.DeploymentName, deployment)
	if err != nil {
		return fmt.Errorf("creating Tenant Template Deployment %q: %+v", id.DeploymentName, err)
	}

	log.Printf("[DEBUG] Waiting for deployment of Tenant Template Deployment %q..", id.DeploymentName)
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of Tenant Template Deployment %q: %+v", id.DeploymentName, err)
	}

	return tenantTemplateDeploymentResourceRead(d, meta)
}

func tenantTemplateDeploymentResourceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.DeploymentsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.TenantTemplateDeploymentID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetAtTenantScope(ctx, id.DeploymentName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Tenant Template Deployment %q was not found - removing from state", id.DeploymentName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Tenant Template Deployment %q: %+v", id.DeploymentName, err)
	}

	templateContents, err := client.ExportTemplateAtTenantScope(ctx, id.DeploymentName)
	if err != nil {
		return fmt.Errorf("retrieving Template Content for Tenant Template Deployment %q: %+v", id.DeploymentName, err)
	}

	d.Set("name", id.DeploymentName)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if props := resp.Properties; props != nil {
		d.Set("debug_level", flattenTemplateDeploymentDebugSetting(props.DebugSetting))

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

		templateLinkId := ""
		if props.TemplateLink != nil {
			if props.TemplateLink.ID != nil {
				templateLinkId = *props.TemplateLink.ID
			}
		}
		d.Set("template_spec_version_id", templateLinkId)
	}

	flattenedTemplate, err := flattenTemplateDeploymentBody(templateContents.Template)
	if err != nil {
		return fmt.Errorf("flattening `template_content`: %+v", err)
	}
	d.Set("template_content", flattenedTemplate)

	return tags.FlattenAndSet(d, resp.Tags)
}

func tenantTemplateDeploymentResourceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.DeploymentsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.TenantTemplateDeploymentID(d.Id())
	if err != nil {
		return err
	}

	// at this time unfortunately the Resources RP doesn't expose a means of deleting top-level objects
	// so we're unable to delete these during deletion - this'll need to be detailed in the docs

	log.Printf("[DEBUG] Deleting Tenant Template Deployment %q..", id.DeploymentName)
	future, err := client.DeleteAtTenantScope(ctx, id.DeploymentName)
	if err != nil {
		return fmt.Errorf("deleting Tenant Template Deployment %q: %+v", id.DeploymentName, err)
	}

	log.Printf("[DEBUG] Waiting for deletion of Tenant Template Deployment %q..", id.DeploymentName)
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of Tenant Template Deployment %q: %+v", id.DeploymentName, err)
	}
	log.Printf("[DEBUG] Deleted Tenant Template Deployment %q.", id.DeploymentName)

	return nil
}

func validateTenantTemplateDeployment(ctx context.Context, id parse.TenantTemplateDeploymentId, deployment resources.ScopedDeployment, client *resources.DeploymentsClient) error {
	validationFuture, err := client.ValidateAtTenantScope(ctx, id.DeploymentName, deployment)
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
