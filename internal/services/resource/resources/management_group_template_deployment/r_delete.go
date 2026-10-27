// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package management_group_template_deployment

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func managementGroupTemplateDeploymentResourceDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Resource.LegacyDeploymentsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ManagementGroupTemplateDeploymentID(d.Id())
	if err != nil {
		return err
	}

	// at this time unfortunately the Resources RP doesn't expose a means of deleting top-level objects
	// so we're unable to delete these during deletion - this'll need to be detailed in the docs

	future, err := client.DeleteAtManagementGroupScope(ctx, id.ManagementGroupName, id.DeploymentName)
	if err != nil {
		return fmt.Errorf("deleting Management Group Template Deployment %q: %+v", id.DeploymentName, err)
	}

	log.Printf("[DEBUG] Waiting for deletion of Management Group Template Deployment %q..", id.DeploymentName)
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of Management Group Template Deployment %q: %+v", id.DeploymentName, err)
	}

	return nil
}
