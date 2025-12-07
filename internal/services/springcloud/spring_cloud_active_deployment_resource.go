// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package springcloud

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/appplatform/2023-05-01-preview/appplatform"
)

func resourceSpringCloudActiveDeployment() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		DeprecationMessage: features.DeprecatedInFivePointOh("Azure Spring Apps is now deprecated and will be retired on 2028-05-31 - as such the `azurerm_spring_cloud_active_deployment` resource is deprecated and will be removed in a future major version of the AzureRM Provider. See https://aka.ms/asaretirement for more information."),

		Create: resourceSpringCloudActiveDeploymentCreate,
		Read:   resourceSpringCloudActiveDeploymentRead,
		Update: resourceSpringCloudActiveDeploymentUpdate,
		Delete: resourceSpringCloudActiveDeploymentDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.SpringCloudActiveDeploymentV0ToV1{},
		}),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.SpringCloudAppID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"spring_cloud_app_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SpringCloudAppID,
			},

			"deployment_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.SpringCloudDeploymentName,
			},
		},
	}
}

func resourceSpringCloudActiveDeploymentCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.AppsClient
	deploymentClient := meta.(*clients.Client).AppPlatform.DeploymentsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	deploymentName := d.Get("deployment_name").(string)
	appId, err := parse.SpringCloudAppID(d.Get("spring_cloud_app_id").(string))
	if err != nil {
		return err
	}

	activeDeployments, err := listSpringCloudActiveDeployments(ctx, deploymentClient, appId)
	if err != nil {
		return err
	}
	if len(activeDeployments) != 0 {
		return tf.ImportAsExistsError("azurerm_spring_cloud_active_deployment", appId.ID())
	}

	parameter := appplatform.ActiveDeploymentCollection{ActiveDeploymentNames: &[]string{deploymentName}}
	future, err := client.SetActiveDeployments(ctx, appId.ResourceGroup, appId.SpringName, appId.AppName, parameter)
	if err != nil {
		return fmt.Errorf("setting active deployment %q (Spring Cloud Service %q / App %q / Resource Group %q): %+v", deploymentName, appId.SpringName, appId.AppName, appId.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for setting active deployment %q (Spring Cloud Service %q / App %q / Resource Group %q): %+v", deploymentName, appId.SpringName, appId.AppName, appId.ResourceGroup, err)
	}

	d.SetId(appId.ID())

	return resourceSpringCloudActiveDeploymentRead(d, meta)
}

func resourceSpringCloudActiveDeploymentUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.AppsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	deploymentName := d.Get("deployment_name").(string)
	appId, err := parse.SpringCloudAppID(d.Get("spring_cloud_app_id").(string))
	if err != nil {
		return err
	}

	parameter := appplatform.ActiveDeploymentCollection{ActiveDeploymentNames: &[]string{deploymentName}}
	future, err := client.SetActiveDeployments(ctx, appId.ResourceGroup, appId.SpringName, appId.AppName, parameter)
	if err != nil {
		return fmt.Errorf("setting active deployment %q (Spring Cloud Service %q / App %q / Resource Group %q): %+v", deploymentName, appId.SpringName, appId.AppName, appId.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for setting active deployment %q (Spring Cloud Service %q / App %q / Resource Group %q): %+v", deploymentName, appId.SpringName, appId.AppName, appId.ResourceGroup, err)
	}

	d.SetId(appId.ID())

	return resourceSpringCloudActiveDeploymentRead(d, meta)
}

func resourceSpringCloudActiveDeploymentRead(d *pluginsdk.ResourceData, meta interface{}) error {
	deploymentClient := meta.(*clients.Client).AppPlatform.DeploymentsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SpringCloudAppID(d.Id())
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

func resourceSpringCloudActiveDeploymentDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.AppsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SpringCloudAppID(d.Id())
	if err != nil {
		return err
	}

	parameter := appplatform.ActiveDeploymentCollection{ActiveDeploymentNames: &[]string{}}
	future, err := client.SetActiveDeployments(ctx, id.ResourceGroup, id.SpringName, id.AppName, parameter)
	if err != nil {
		return fmt.Errorf("deleting Active Deployment (Spring Cloud Service %q / App %q / Resource Group %q): %+v", id.SpringName, id.AppName, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deleting active deployment (Spring Cloud Service %q / App %q / Resource Group %q): %+v", id.SpringName, id.AppName, id.ResourceGroup, err)
	}

	return nil
}

func listSpringCloudActiveDeployments(ctx context.Context, client *appplatform.DeploymentsClient, id *parse.SpringCloudAppId) ([]string, error) {
	it, err := client.ListComplete(ctx, id.ResourceGroup, id.SpringName, id.AppName, nil)
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
