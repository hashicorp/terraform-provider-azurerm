package azurerm

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/appplatform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmSpringCloudActiveDeployment() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSpringCloudActiveDeploymentCreateUpdate,
		Read:   resourceArmSpringCloudActiveDeploymentRead,
		Update: resourceArmSpringCloudActiveDeploymentCreateUpdate,
		Delete: resourceArmSpringCloudActiveDeploymentDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"spring_cloud_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: appplatform.ValidateSpringCloudName,
			},

			"spring_cloud_app_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: appplatform.ValidateSpringCloudName,
			},

			"spring_cloud_deployment_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: appplatform.ValidateSpringCloudName,
			},
		},
	}
}

func resourceArmSpringCloudActiveDeploymentCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).AppPlatform.ServicesClient
	appsClient := meta.(*ArmClient).AppPlatform.AppsClient
	deploymentsClient := meta.(*ArmClient).AppPlatform.DeploymentsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	springCloudName := d.Get("spring_cloud_name").(string)
	appName := d.Get("spring_cloud_app_name").(string)
	deploymentName := d.Get("spring_cloud_deployment_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if resp, err := client.Get(ctx, resourceGroup, springCloudName); err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("[DEBUG] Spring Cloud Service %q (resource group %q) was not found.", springCloudName, resourceGroup)
		}
		return fmt.Errorf("Error making Read request on AzureRM Spring Cloud Service %q: %+v", springCloudName, err)
	}

	resp, err := appsClient.Get(ctx, resourceGroup, springCloudName, appName, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("[DEBUG] Spring Cloud App %q (Spring Cloud service %q / resource group %q) was not found.", appName, springCloudName, resourceGroup)
		}
		return fmt.Errorf("Error making Read request on AzureRM Spring Cloud App %q (Spring Cloud service %q / resource group %q): %+v", appName, springCloudName, resourceGroup, err)
	}

	if resp, err := deploymentsClient.Get(ctx, resourceGroup, springCloudName, appName, deploymentName); err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("[DEBUG] Spring Cloud deployment %q (Spring Cloud service %q / Spring Cloud App %q / resource group %q) was not found.", deploymentName, springCloudName, appName, resourceGroup)
		}
		return fmt.Errorf("Error making Read request on AzureRM Spring Cloud deployment %q (Spring Cloud service %q / Spring Cloud App %q / resource group %q): %+v", deploymentName, springCloudName, appName, resourceGroup, err)
	}

	resp.Properties.ActiveDeploymentName = &deploymentName

	future, err := appsClient.Update(ctx, resourceGroup, springCloudName, appName, &resp)
	if err != nil {
		return fmt.Errorf("Error swapping active deployment %q (Spring Cloud Service %q / Spring Cloud App %q / Resource Group %q): %+v", deploymentName, springCloudName, appName, resourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error swapping active deployment %q (Spring Cloud Service %q / Spring Cloud App %q / Resource Group %q): %+v", deploymentName, springCloudName, appName, resourceGroup, err)
	}
	d.SetId(*resp.ID)

	return resourceArmSpringCloudActiveDeploymentRead(d, meta)
}

func resourceArmSpringCloudActiveDeploymentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).AppPlatform.AppsClient
	ctx, cancel := timeouts.ForRead(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	springCloudName := id.Path["Spring"]
	appName := id.Path["apps"]

	resp, err := client.Get(ctx, resourceGroup, springCloudName, appName, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Spring Cloud App %q was not found - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on AzureRM Spring Cloud App %q (Spring Cloud Service %q / Resource Group %q): %+v", appName, springCloudName, resourceGroup, err)
	}

	d.Set("spring_cloud_name", springCloudName)
	d.Set("spring_cloud_app_name", resp.Name)
	d.Set("spring_cloud_deployment_name", resp.Properties.ActiveDeploymentName)
	d.Set("resource_group_name", resourceGroup)
	return nil
}

func resourceArmSpringCloudActiveDeploymentDelete(d *schema.ResourceData, meta interface{}) error {
	// There is nothing to delete and the server side can not update app active_deployment_name to empty
	// so return nil
	return nil
}
