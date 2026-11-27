// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_active_deployment

import (
	"fmt"

	appplatform_rm "github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	springcloudClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/appplatform/2023-05-01-preview/appplatform"
)

func resourceSpringCloudActiveDeploymentDelete(d *pluginsdk.ResourceData, meta any) error {
	client := springcloudClient.From(meta).AppsClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	// the ID is parsed insensitively to preserve the behaviour of the legacy parser this replaced -
	// we are ok with this remaining insensitive as Azure Spring Apps is deprecated and will be removed
	id, err := appplatform_rm.ParseAppIDInsensitively(d.Id())
	if err != nil {
		return err
	}

	parameter := appplatform.ActiveDeploymentCollection{ActiveDeploymentNames: &[]string{}}
	future, err := client.SetActiveDeployments(ctx, id.ResourceGroupName, id.SpringName, id.AppName, parameter)
	if err != nil {
		return fmt.Errorf("deleting Active Deployment (Spring Cloud Service %q / App %q / Resource Group %q): %+v", id.SpringName, id.AppName, id.ResourceGroupName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deleting active deployment (Spring Cloud Service %q / App %q / Resource Group %q): %+v", id.SpringName, id.AppName, id.ResourceGroupName, err)
	}

	return nil
}
