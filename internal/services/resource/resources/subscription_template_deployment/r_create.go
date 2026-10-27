// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package subscription_template_deployment

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2020-06-01/resources" //nolint:staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2023-07-01/deployments"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func subscriptionTemplateDeploymentResourceCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Resource.LegacyDeploymentsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := deployments.NewProviderDeploymentID(subscriptionId, d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.GetAtSubscriptionScope(ctx, id.DeploymentName)
		if err != nil {
			if !response.WasNotFound(existing.Response.Response) {
				return fmt.Errorf("checking for presence of existing Subscription Template Deployment %q: %+v", id.DeploymentName, err)
			}
		}

		if !response.WasNotFound(existing.Response.Response) {
			return tf.ImportAsExistsError("azurerm_subscription_template_deployment", id.ID())
		}
	}

	deployment := resources.Deployment{
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Properties: &resources.DeploymentProperties{
			DebugSetting: helpers.ExpandTemplateDeploymentDebugSetting(d.Get("debug_level").(string)),
			Mode:         resources.DeploymentModeIncremental,
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if templateRaw, ok := d.GetOk("template_content"); ok {
		template, err := helpers.ExpandTemplateDeploymentBody(templateRaw.(string))
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
		parameters, err := helpers.ExpandTemplateDeploymentBody(v.(string))
		if err != nil {
			return fmt.Errorf("expanding `parameters_content`: %+v", err)
		}
		deployment.Properties.Parameters = parameters
	}

	log.Printf("[DEBUG] Running validation of Subscription Template Deployment %q..", id.DeploymentName)
	if err := validateSubscriptionTemplateDeployment(ctx, id, deployment, client); err != nil {
		return fmt.Errorf("validating Subscription Template Deployment %q: %+v", id.DeploymentName, err)
	}
	log.Printf("[DEBUG] Validated Subscription Template Deployment %q..", id.DeploymentName)

	log.Printf("[DEBUG] Provisioning Subscription Template Deployment %q..", id.DeploymentName)
	future, err := client.CreateOrUpdateAtSubscriptionScope(ctx, id.DeploymentName, deployment)
	if err != nil {
		return fmt.Errorf("creating Subscription Template Deployment %q: %+v", id.DeploymentName, err)
	}

	d.SetId(id.ID())

	log.Printf("[DEBUG] Waiting for deployment of Subscription Template Deployment %q..", id.DeploymentName)
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of Subscription Template Deployment %q: %+v", id.DeploymentName, err)
	}

	return subscriptionTemplateDeploymentResourceRead(d, meta)
}
