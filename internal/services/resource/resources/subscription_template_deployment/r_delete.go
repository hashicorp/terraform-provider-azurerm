// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package subscription_template_deployment

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2023-07-01/deployments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	resourceClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func subscriptionTemplateDeploymentResourceDelete(d *pluginsdk.ResourceData, meta any) error {
	client := resourceClient.From(meta).LegacyDeploymentsClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := deployments.ParseProviderDeploymentID(d.Id())
	if err != nil {
		return err
	}

	// at this time unfortunately the Resources RP doesn't expose a means of deleting top-level objects
	// so we're unable to delete these during deletion - this'll need to be detailed in the docs

	future, err := client.DeleteAtSubscriptionScope(ctx, id.DeploymentName)
	if err != nil {
		return fmt.Errorf("deleting Subscription Template Deployment %q: %+v", id.DeploymentName, err)
	}

	log.Printf("[DEBUG] Waiting for deletion of Subscription Template Deployment %q..", id.DeploymentName)
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of Subscription Template Deployment %q: %+v", id.DeploymentName, err)
	}

	return nil
}
