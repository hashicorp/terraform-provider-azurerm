// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package management_group_template_deployment

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2020-06-01/resources" //nolint:staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	mgParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/managementgroup/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func managementGroupTemplateDeploymentResourceCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Resource.LegacyDeploymentsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	managementGroupId, err := mgParse.ManagementGroupID(d.Get("management_group_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewManagementGroupTemplateDeploymentID(managementGroupId.Name, d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.GetAtManagementGroupScope(ctx, id.ManagementGroupName, id.DeploymentName)
		if err != nil {
			if !response.WasNotFound(existing.Response.Response) {
				return fmt.Errorf("checking for presence of existing Management Group Template Deployment %q: %+v", id.DeploymentName, err)
			}
		}
		if existing.Properties != nil {
			return tf.ImportAsExistsError("azurerm_management_group_template_deployment", id.ID())
		}
	}

	deployment := resources.ScopedDeployment{
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

	log.Printf("[DEBUG] Running validation of Management Group Template Deployment %q..", id.DeploymentName)
	if err := validateManagementGroupTemplateDeployment(ctx, id, deployment, client); err != nil {
		return fmt.Errorf("validating Management Group Template Deployment %q: %+v", id.DeploymentName, err)
	}
	log.Printf("[DEBUG] Validated Management Group Template Deployment %q..", id.DeploymentName)

	log.Printf("[DEBUG] Provisioning Management Group Template Deployment %q..", id.DeploymentName)
	future, err := client.CreateOrUpdateAtManagementGroupScope(ctx, id.ManagementGroupName, id.DeploymentName, deployment)
	if err != nil {
		return fmt.Errorf("creating Management Group Template Deployment %q: %+v", id.DeploymentName, err)
	}

	d.SetId(id.ID())

	log.Printf("[DEBUG] Waiting for deployment of Management Group Template Deployment %q..", id.DeploymentName)
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of Management Group Template Deployment %q: %+v", id.DeploymentName, err)
	}

	return managementGroupTemplateDeploymentResourceRead(d, meta)
}
