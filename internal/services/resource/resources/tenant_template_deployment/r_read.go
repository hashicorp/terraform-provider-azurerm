// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package tenant_template_deployment

import (
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2020-06-01/resources" //nolint:staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func tenantTemplateDeploymentResourceRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Resource.LegacyDeploymentsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.TenantTemplateDeploymentID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetAtTenantScope(ctx, id.DeploymentName)
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
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
		d.Set("debug_level", helpers.FlattenTemplateDeploymentDebugSetting(props.DebugSetting))

		filteredParams := helpers.FilterOutTemplateDeploymentParameters(props.Parameters)
		flattenedParams, err := helpers.FlattenTemplateDeploymentBody(filteredParams)
		if err != nil {
			return fmt.Errorf("flattening `parameters_content`: %+v", err)
		}
		d.Set("parameters_content", flattenedParams)

		flattenedOutputs, err := helpers.FlattenTemplateDeploymentBody(props.Outputs)
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

	flattenedTemplate, err := helpers.FlattenTemplateDeploymentBody(templateContents.Template)
	if err != nil {
		return fmt.Errorf("flattening `template_content`: %+v", err)
	}
	d.Set("template_content", flattenedTemplate)

	return tags.FlattenAndSet(d, resp.Tags)
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
