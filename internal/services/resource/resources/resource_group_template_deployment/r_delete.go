// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package resource_group_template_deployment

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2023-07-01/deployments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceGroupTemplateDeploymentResourceDelete(d *pluginsdk.ResourceData, meta any) error {
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
