// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2020-06-01/resources" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2022-02-01/templatespecversions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2023-07-01/deployments"
	"github.com/hashicorp/terraform-provider-azurerm/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceGroupTemplateDeploymentResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceGroupTemplateDeploymentResourceCreate,
		Read:   resourceGroupTemplateDeploymentResourceRead,
		Update: resourceGroupTemplateDeploymentResourceUpdate,
		Delete: resourceGroupTemplateDeploymentResourceDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := deployments.ParseResourceGroupProviderDeploymentID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			// v0 -> v1 normalises the casing of IDs imported while this resource parsed them with the
			// case-insensitive legacy parser, so they can be parsed with the case-sensitive SDK parser
			0: migration.ResourceGroupTemplateDeploymentV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(180 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(180 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(180 * time.Minute),
		},

		// (@jackofallops - lintignore needed as we need to make sure the JSON is usable in `output_content`)

		// lintignore:S033
		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.TemplateDeploymentName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"deployment_mode": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(resources.DeploymentModeComplete),
					string(resources.DeploymentModeIncremental),
				}, false),
			},

			"template_content": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
				ExactlyOneOf: []string{
					"template_content",
					"template_spec_version_id",
				},
				StateFunc: helpers.NormalizeJson,
			},

			"template_spec_version_id": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ExactlyOneOf: []string{
					"template_content",
					"template_spec_version_id",
				},
				ValidateFunc: validation.AsGeneratedID(templatespecversions.ParseTemplateSpecVersionIDInsensitively),
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
				StateFunc: helpers.NormalizeJson,
			},

			"tags": commonschema.Tags(),

			// Computed
			"output_content": {
				Type:     pluginsdk.TypeString,
				Computed: true,
				// NOTE:  outputs can be strings, ints, objects etc - whilst using a nested object was considered
				// parsing the JSON using `jsondecode` allows the users to interact with/map objects as required
			},
		},

		// this is needed to fix https://github.com/hashicorp/terraform-provider-azurerm/issues/12828
		// On a change to `template_content` or `parameters_content`, we'll set `output_content` to empty
		// The adverse effect of this is that any change to `template_content` will also cause any resource referencing `output_content` to update
		CustomizeDiff: func(ctx context.Context, d *pluginsdk.ResourceDiff, i interface{}) error {
			if d.HasChange("template_content") {
				o, n := d.GetChange("template_content")

				// the json has to be normalized and then compared against to see if a change has occurred
				if !strings.EqualFold(o.(string), helpers.NormalizeJson(n)) {
					return d.SetNewComputed("output_content")
				}
			}

			if d.HasChange("parameters_content") {
				o, n := d.GetChange("parameters_content")

				// the json has to be normalized and then compared against to see if a change has occurred
				if !strings.EqualFold(o.(string), helpers.NormalizeJson(n)) {
					return d.SetNewComputed("output_content")
				}
			}

			return nil
		},
	}
}

func resourceGroupTemplateDeploymentResourceCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.LegacyDeploymentsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := deployments.NewResourceGroupProviderDeploymentID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id.ResourceGroupName, id.DeploymentName)
		if err != nil {
			if !response.WasNotFound(existing.Response.Response) {
				return fmt.Errorf("checking for presence of existing Template Deployment %q (Resource Group %q): %+v", id.ResourceGroupName, id.DeploymentName, err)
			}
		}

		if !response.WasNotFound(existing.Response.Response) {
			return tf.ImportAsExistsError("azurerm_resource_group_template_deployment", id.ID())
		}
	}

	deployment := resources.Deployment{
		Properties: &resources.DeploymentProperties{
			DebugSetting: expandTemplateDeploymentDebugSetting(d.Get("debug_level").(string)),
			Mode:         resources.DeploymentMode(d.Get("deployment_mode").(string)),
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
			ID: pointer.To(templateSpecVersionID.(string)),
		}
	}

	if v, ok := d.GetOk("parameters_content"); ok && v != "" {
		parameters, err := expandTemplateDeploymentBody(v.(string))
		if err != nil {
			return fmt.Errorf("expanding `parameters_content`: %+v", err)
		}
		deployment.Properties.Parameters = parameters
	}

	log.Printf("[DEBUG] Running validation of Template Deployment %q (Resource Group %q)..", id.DeploymentName, id.ResourceGroupName)
	if err := validateResourceGroupTemplateDeployment(ctx, id, deployment, client); err != nil {
		return fmt.Errorf("validating Template Deployment %q (Resource Group %q): %+v", id.DeploymentName, id.ResourceGroupName, err)
	}
	log.Printf("[DEBUG] Validated Template Deployment %q (Resource Group %q)..", id.DeploymentName, id.ResourceGroupName)

	log.Printf("[DEBUG] Provisioning Template Deployment %q (Resource Group %q)..", id.DeploymentName, id.ResourceGroupName)
	future, err := client.CreateOrUpdate(ctx, id.ResourceGroupName, id.DeploymentName, deployment)
	if err != nil {
		return fmt.Errorf("creating Template Deployment %q (Resource Group %q): %+v", id.DeploymentName, id.ResourceGroupName, err)
	}

	d.SetId(id.ID())

	log.Printf("[DEBUG] Waiting for deployment of Template Deployment %q (Resource Group %q)..", id.DeploymentName, id.ResourceGroupName)
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of Template Deployment %q (Resource Group %q): %+v", id.DeploymentName, id.ResourceGroupName, err)
	}

	return resourceGroupTemplateDeploymentResourceRead(d, meta)
}

func resourceGroupTemplateDeploymentResourceUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.LegacyDeploymentsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := deployments.ParseResourceGroupProviderDeploymentID(d.Id())
	if err != nil {
		return err
	}

	template, err := client.Get(ctx, id.ResourceGroupName, id.DeploymentName)
	if err != nil {
		return fmt.Errorf("retrieving Template Deployment %q (Resource Group %q): %+v", id.DeploymentName, id.ResourceGroupName, err)
	}
	if template.Properties == nil {
		return fmt.Errorf("retrieving Template Deployment %q (Resource Group %q): `properties` was nil", id.DeploymentName, id.ResourceGroupName)
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

	parameters, err := expandTemplateDeploymentBody(d.Get("parameters_content").(string))
	if err != nil {
		return fmt.Errorf("expanding `parameters_content`: %+v", err)
	}
	deployment.Properties.Parameters = parameters

	if d.HasChange("template_content") {
		templateContents, err := expandTemplateDeploymentBody(d.Get("template_content").(string))
		if err != nil {
			return fmt.Errorf("expanding `template_content`: %+v", err)
		}

		deployment.Properties.Template = templateContents
	} else {
		// retrieve the existing content and reuse that
		exportedTemplate, err := client.ExportTemplate(ctx, id.ResourceGroupName, id.DeploymentName)
		if err != nil {
			return fmt.Errorf("retrieving Contents for Template Deployment %q (Resource Group %q): %+v", id.DeploymentName, id.ResourceGroupName, err)
		}

		deployment.Properties.Template = exportedTemplate.Template
	}

	if d.HasChange("template_spec_version_id") {
		deployment.Properties.TemplateLink = &resources.TemplateLink{
			ID: pointer.To(d.Get("template_spec_version_id").(string)),
		}

		if d.Get("template_spec_version_id").(string) != "" {
			deployment.Properties.Template = nil
		}
	}

	if d.HasChange("tags") {
		deployment.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	log.Printf("[DEBUG] Running validation of Template Deployment %q (Resource Group %q)..", id.DeploymentName, id.ResourceGroupName)
	if err := validateResourceGroupTemplateDeployment(ctx, *id, deployment, client); err != nil {
		return fmt.Errorf("validating Template Deployment %q (Resource Group %q): %+v", id.DeploymentName, id.ResourceGroupName, err)
	}
	log.Printf("[DEBUG] Validated Template Deployment %q (Resource Group %q)..", id.DeploymentName, id.ResourceGroupName)

	log.Printf("[DEBUG] Provisioning Template Deployment %q (Resource Group %q)..", id.DeploymentName, id.ResourceGroupName)
	future, err := client.CreateOrUpdate(ctx, id.ResourceGroupName, id.DeploymentName, deployment)
	if err != nil {
		return fmt.Errorf("creating Template Deployment %q (Resource Group %q): %+v", id.DeploymentName, id.ResourceGroupName, err)
	}

	log.Printf("[DEBUG] Waiting for deployment of Template Deployment %q (Resource Group %q)..", id.DeploymentName, id.ResourceGroupName)
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of Template Deployment %q (Resource Group %q): %+v", id.DeploymentName, id.ResourceGroupName, err)
	}

	return resourceGroupTemplateDeploymentResourceRead(d, meta)
}

func resourceGroupTemplateDeploymentResourceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.LegacyDeploymentsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := deployments.ParseResourceGroupProviderDeploymentID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroupName, id.DeploymentName)
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			log.Printf("[DEBUG] Template Deployment %q (Resource Group %q) was not found - removing from state", id.DeploymentName, id.ResourceGroupName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Template Deployment %q (Resource Group %q): %+v", id.DeploymentName, id.ResourceGroupName, err)
	}

	templateContents, err := client.ExportTemplate(ctx, id.ResourceGroupName, id.DeploymentName)
	if err != nil {
		return fmt.Errorf("retrieving Template Content for Template Deployment %q (Resource Group %q): %+v", id.DeploymentName, id.ResourceGroupName, err)
	}

	d.Set("name", id.DeploymentName)
	d.Set("resource_group_name", id.ResourceGroupName)

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

func resourceGroupTemplateDeploymentResourceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.LegacyDeploymentsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := deployments.ParseResourceGroupProviderDeploymentID(d.Id())
	if err != nil {
		return err
	}

	template, err := client.Get(ctx, id.ResourceGroupName, id.DeploymentName)
	if err != nil {
		if response.WasNotFound(template.Response.Response) {
			return nil
		}

		return fmt.Errorf("retrieving Template Deployment %q (Resource Group %q): %+v", id.DeploymentName, id.ResourceGroupName, err)
	}
	if template.Properties == nil {
		return fmt.Errorf("`properties` was nil for template`")
	}

	deleteItemsInTemplate := meta.(*clients.Client).Features.TemplateDeployment.DeleteNestedItemsDuringDeletion
	if deleteItemsInTemplate {
		resourceClient := meta.(*clients.Client).Resource
		log.Printf("[DEBUG] Removing items provisioned by the Template Deployment %q (Resource Group %q)..", id.DeploymentName, id.ResourceGroupName)
		if err := deleteItemsProvisionedByTemplate(ctx, resourceClient, *template.Properties, id.SubscriptionId); err != nil {
			return fmt.Errorf("removing items provisioned by this Template Deployment: %+v", err)
		}
		log.Printf("[DEBUG] Removed items provisioned by the Template Deployment %q (Resource Group %q)..", id.DeploymentName, id.ResourceGroupName)
	} else {
		log.Printf("[DEBUG] Skipping removing items provisioned by the Template Deployment %q (Resource Group %q) as the feature is disabled", id.DeploymentName, id.ResourceGroupName)
	}

	future, err := client.Delete(ctx, id.ResourceGroupName, id.DeploymentName)
	if err != nil {
		return fmt.Errorf("deleting Template Deployment %q (Resource Group %q): %+v", id.DeploymentName, id.ResourceGroupName, err)
	}

	log.Printf("[DEBUG] Waiting for deletion of Template Deployment %q (Resource Group %q)..", id.DeploymentName, id.ResourceGroupName)
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of Template Deployment %q (Resource Group %q): %+v", id.DeploymentName, id.ResourceGroupName, err)
	}

	return nil
}

func validateResourceGroupTemplateDeployment(ctx context.Context, id deployments.ResourceGroupProviderDeploymentId, deployment resources.Deployment, client *resources.DeploymentsClient) error {
	validationFuture, err := client.Validate(ctx, id.ResourceGroupName, id.DeploymentName, deployment)
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
