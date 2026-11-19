// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_active_deployment

import (
	"context"
	"fmt"
	"log"

	appplatform_rm "github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/appplatform/2023-05-01-preview/appplatform"
)

func resourceSpringCloudActiveDeploymentRead(d *pluginsdk.ResourceData, meta any) error {
	deploymentClient := meta.(*clients.Client).AppPlatform.DeploymentsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// the ID is parsed insensitively to preserve the behaviour of the legacy parser this replaced -
	// we are ok with this remaining insensitive as Azure Spring Apps is deprecated and will be removed
	id, err := appplatform_rm.ParseAppIDInsensitively(d.Id())
	if err != nil {
		return err
	}

	activeDeployments, err := listSpringCloudActiveDeployments(ctx, deploymentClient, id)
	if err != nil {
		return err
	}
	if len(activeDeployments) == 0 {
		log.Printf("[INFO] Spring Cloud App %q does not exist - removing from state", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("spring_cloud_app_id", id.ID())
	d.Set("deployment_name", activeDeployments[0])

	return nil
}

func listSpringCloudActiveDeployments(ctx context.Context, client *appplatform.DeploymentsClient, id *appplatform_rm.AppId) ([]string, error) {
	it, err := client.ListComplete(ctx, id.ResourceGroupName, id.SpringName, id.AppName, nil)
	if err != nil {
		return nil, fmt.Errorf("listing active deployment for %q: %+v", id, err)
	}
	deployments := make([]string, 0)
	for it.NotDone() {
		value := it.Value()
		if value.Properties != nil && value.Properties.Active != nil && *value.Properties.Active {
			deployments = append(deployments, *value.Name)
		}
		if err := it.NextWithContext(ctx); err != nil {
			return nil, fmt.Errorf("listing active deployment for %q: %+v", id, err)
		}
	}
	return deployments, nil
}
