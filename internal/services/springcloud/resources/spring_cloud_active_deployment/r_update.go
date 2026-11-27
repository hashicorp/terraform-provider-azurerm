// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_active_deployment

import (
	"fmt"

	appplatform_rm "github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/appplatform/2023-05-01-preview/appplatform"
)

func resourceSpringCloudActiveDeploymentUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).AppPlatform.AppsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	deploymentName := d.Get("deployment_name").(string)
	// the ID is parsed insensitively to preserve the behaviour of the legacy parser this replaced -
	// we are ok with this remaining insensitive as Azure Spring Apps is deprecated and will be removed
	appId, err := appplatform_rm.ParseAppIDInsensitively(d.Get("spring_cloud_app_id").(string))
	if err != nil {
		return err
	}

	parameter := appplatform.ActiveDeploymentCollection{ActiveDeploymentNames: &[]string{deploymentName}}
	future, err := client.SetActiveDeployments(ctx, appId.ResourceGroupName, appId.SpringName, appId.AppName, parameter)
	if err != nil {
		return fmt.Errorf("setting active deployment %q (Spring Cloud Service %q / App %q / Resource Group %q): %+v", deploymentName, appId.SpringName, appId.AppName, appId.ResourceGroupName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for setting active deployment %q (Spring Cloud Service %q / App %q / Resource Group %q): %+v", deploymentName, appId.SpringName, appId.AppName, appId.ResourceGroupName, err)
	}

	d.SetId(appId.ID())

	return resourceSpringCloudActiveDeploymentRead(d, meta)
}
