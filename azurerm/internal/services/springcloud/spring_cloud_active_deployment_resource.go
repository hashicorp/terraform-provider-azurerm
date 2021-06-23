package springcloud

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/appplatform/mgmt/2020-11-01-preview/appplatform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/springcloud/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/springcloud/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceSpringCloudActiveDeployment() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSpringCloudActiveDeploymentCreate,
		Read:   resourceSpringCloudActiveDeploymentRead,
		Update: resourceSpringCloudActiveDeploymentUpdate,
		Delete: resourceSpringCloudActiveDeploymentDelete,

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
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	deploymentName := d.Get("deployment_name").(string)
	appId, err := parse.SpringCloudAppID(d.Get("spring_cloud_app_id").(string))
	if err != nil {
		return err
	}

	resourceId := parse.NewSpringCloudAppID(appId.SubscriptionId, appId.ResourceGroup, appId.SpringName, appId.AppName).ID()
	existing, err := client.Get(ctx, appId.ResourceGroup, appId.SpringName, appId.AppName, "")
	if err != nil {
		return fmt.Errorf("making Read request on AzureRM Spring Cloud App %q (Spring Cloud service %q / resource group %q): %+v", appId.AppName, appId.SpringName, appId.ResourceGroup, err)
	}

	if existing.Properties != nil && existing.Properties.ActiveDeploymentName != nil && *existing.Properties.ActiveDeploymentName != "" {
		return tf.ImportAsExistsError("azurerm_spring_cloud_active_deployment", resourceId)
	}

	if existing.Properties == nil {
		existing.Properties = &appplatform.AppResourceProperties{}
	}
	existing.Properties.ActiveDeploymentName = &deploymentName

	future, err := client.CreateOrUpdate(ctx, appId.ResourceGroup, appId.SpringName, appId.AppName, existing)
	if err != nil {
		return fmt.Errorf("setting active deployment %q (Spring Cloud Service %q / App %q / Resource Group %q): %+v", deploymentName, appId.SpringName, appId.AppName, appId.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for setting active deployment %q (Spring Cloud Service %q / App %q / Resource Group %q): %+v", deploymentName, appId.SpringName, appId.AppName, appId.ResourceGroup, err)
	}

	d.SetId(resourceId)

	return resourceSpringCloudActiveDeploymentRead(d, meta)
}

func resourceSpringCloudActiveDeploymentUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.AppsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SpringCloudAppID(d.Id())
	if err != nil {
		return err
	}

	deploymentName := d.Get("deployment_name").(string)
	app := appplatform.AppResource{
		Properties: &appplatform.AppResourceProperties{
			ActiveDeploymentName: utils.String(deploymentName),
		},
	}

	future, err := client.Update(ctx, id.ResourceGroup, id.SpringName, id.AppName, app)
	if err != nil {
		return fmt.Errorf("updating Active Deployment %q (Spring Cloud Service %q / App %q / Resource Group %q): %+v", deploymentName, id.SpringName, id.AppName, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for Update of Active Deployment %q (Spring Cloud Service %q / App %q / Resource Group %q): %+v", deploymentName, id.SpringName, id.AppName, id.ResourceGroup, err)
	}

	return resourceSpringCloudActiveDeploymentRead(d, meta)
}

func resourceSpringCloudActiveDeploymentRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.AppsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SpringCloudAppID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.SpringName, id.AppName, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Spring Cloud App %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading Active Deployment for Spring Cloud App %q (Spring Cloud Service %q / resource Group %q): %+v", id.AppName, id.SpringName, id.ResourceGroup, err)
	}

	if resp.Properties == nil || resp.Properties.ActiveDeploymentName == nil {
		log.Printf("[DEBUG] Spring Cloud App %q (Spring Cloud Service %q / Resource Group %q) doesn't have Active Deployment - removing from state!", id.AppName, id.SpringName, id.ResourceGroup)
		d.SetId("")
		return nil
	}

	d.Set("spring_cloud_app_id", id.ID())
	d.Set("deployment_name", resp.Properties.ActiveDeploymentName)

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

	app := appplatform.AppResource{
		Properties: &appplatform.AppResourceProperties{
			ActiveDeploymentName: utils.String(""),
		},
	}

	future, err := client.Update(ctx, id.ResourceGroup, id.SpringName, id.AppName, app)
	if err != nil {
		return fmt.Errorf("deleting Active Deployment (Spring Cloud Service %q / App %q / Resource Group %q): %+v", id.SpringName, id.AppName, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deleting active deployment (Spring Cloud Service %q / App %q / Resource Group %q): %+v", id.SpringName, id.AppName, id.ResourceGroup, err)
	}

	return nil
}
