// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package resource_group_template_deployment

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2020-06-01/resources" //nolint:staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2023-07-01/deployments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceGroupTemplateDeploymentResourceUpdate(d *pluginsdk.ResourceData, meta any) error {
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
		deployment.Properties.DebugSetting = helpers.ExpandTemplateDeploymentDebugSetting(d.Get("debug_level").(string))
	}

	if d.HasChange("deployment_mode") {
		deployment.Properties.Mode = resources.DeploymentMode(d.Get("deployment_mode").(string))
	}

	parameters, err := helpers.ExpandTemplateDeploymentBody(d.Get("parameters_content").(string))
	if err != nil {
		return fmt.Errorf("expanding `parameters_content`: %+v", err)
	}
	deployment.Properties.Parameters = parameters

	if d.HasChange("template_content") {
		templateContents, err := helpers.ExpandTemplateDeploymentBody(d.Get("template_content").(string))
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
		deployment.Tags = tags.Expand(d.Get("tags").(map[string]any))
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
