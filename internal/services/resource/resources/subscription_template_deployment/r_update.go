// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package subscription_template_deployment

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

func subscriptionTemplateDeploymentResourceUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Resource.LegacyDeploymentsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := deployments.ParseProviderDeploymentID(d.Id())
	if err != nil {
		return err
	}

	template, err := client.GetAtSubscriptionScope(ctx, id.DeploymentName)
	if err != nil {
		return fmt.Errorf("retrieving Subscription Template Deployment %q: %+v", id.DeploymentName, err)
	}
	if template.Properties == nil {
		return fmt.Errorf("retrieving Subscription Template Deployment %q: `properties` was nil", id.DeploymentName)
	}

	// the API doesn't have a Patch operation, so we'll need to build one
	deployment := resources.Deployment{
		Location: template.Location,
		Properties: &resources.DeploymentProperties{
			DebugSetting: template.Properties.DebugSetting,
			Mode:         resources.DeploymentModeIncremental,
		},
		Tags: template.Tags,
	}

	if d.HasChange("debug_level") {
		deployment.Properties.DebugSetting = helpers.ExpandTemplateDeploymentDebugSetting(d.Get("debug_level").(string))
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
		exportedTemplate, err := client.ExportTemplateAtSubscriptionScope(ctx, id.DeploymentName)
		if err != nil {
			return fmt.Errorf("retrieving Contents for Subscription Template Deployment %q: %+v", id.DeploymentName, err)
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

	log.Printf("[DEBUG] Running validation of Subscription Template Deployment %q..", id.DeploymentName)
	if err := validateSubscriptionTemplateDeployment(ctx, *id, deployment, client); err != nil {
		return fmt.Errorf("validating Subscription Template Deployment %q: %+v", id.DeploymentName, err)
	}
	log.Printf("[DEBUG] Validated Subscription Template Deployment %q..", id.DeploymentName)

	log.Printf("[DEBUG] Provisioning Subscription Template Deployment %q)..", id.DeploymentName)
	future, err := client.CreateOrUpdateAtSubscriptionScope(ctx, id.DeploymentName, deployment)
	if err != nil {
		return fmt.Errorf("creating Subscription Template Deployment %q: %+v", id.DeploymentName, err)
	}

	log.Printf("[DEBUG] Waiting for deployment of Subscription Template Deployment %q..", id.DeploymentName)
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of Subscription Template Deployment %q: %+v", id.DeploymentName, err)
	}

	return subscriptionTemplateDeploymentResourceRead(d, meta)
}
