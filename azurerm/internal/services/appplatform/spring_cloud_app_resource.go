package appplatform

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/appplatform/mgmt/2019-05-01-preview/appplatform"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/appplatform/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/appplatform/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmSpringCloudApp() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSpringCloudAppCreate,
		Read:   resourceArmSpringCloudAppRead,
		Delete: resourceArmSpringCloudAppDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.SpringCloudAppID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SpringCloudAppName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"service_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SpringCloudServiceName,
			},
		},
	}
}

func resourceArmSpringCloudAppCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.AppsClient
	deploymentClient := meta.(*clients.Client).AppPlatform.DeploymentsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	serviceName := d.Get("service_name").(string)

	existing, err := client.Get(ctx, resourceGroup, serviceName, name, "")
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing Spring Cloud App %q (Spring Cloud Service %q / Resource Group %q): %+v", name, serviceName, resourceGroup, err)
		}
	}
	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_spring_cloud_app", *existing.ID)
	}

	// create app
	appResource := appplatform.AppResource{}
	future, err := client.CreateOrUpdate(ctx, resourceGroup, serviceName, name, appResource)
	if err != nil {
		return fmt.Errorf("creating Spring Cloud App %q (Spring Cloud Service %q / Resource Group %q): %+v", name, serviceName, resourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of Spring Cloud App %q (Spring Cloud Service %q / Resource Group %q): %+v", name, serviceName, resourceGroup, err)
	}

	// create default "hello world" deployment
	deploymentName := "default"
	deploymentResource := appplatform.DeploymentResource{
		Properties: &appplatform.DeploymentResourceProperties{
			AppName: utils.String(name),
			DeploymentSettings: &appplatform.DeploymentSettings{
				CPU:            utils.Int32(1),
				MemoryInGB:     utils.Int32(1),
				RuntimeVersion: appplatform.Java8,
			},
			Source: &appplatform.UserSourceInfo{
				Type:         appplatform.Jar,
				RelativePath: utils.String("<default>"),
			},
		},
	}
	deployFuture, err := deploymentClient.CreateOrUpdate(ctx, resourceGroup, serviceName, name, deploymentName, deploymentResource)
	if err != nil {
		return fmt.Errorf("creating default Deployment %q (Spring Cloud Service %q / App Name %q / Resource Group %q): %+v", deploymentName, serviceName, name, resourceGroup, err)
	}
	if err = deployFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of default Deployment %q (Spring Cloud Service %q / App Name %q /  Resource Group %q): %+v", deploymentName, serviceName, name, resourceGroup, err)
	}

	// binding app with default deployment
	appResource.Properties = &appplatform.AppResourceProperties{
		ActiveDeploymentName: &deploymentName,
	}
	updateFuture, err := client.Update(ctx, resourceGroup, serviceName, name, appResource)
	if err != nil {
		return fmt.Errorf("updating Spring Cloud App %q (Spring Cloud Service %q / Resource Group %q): %+v", name, serviceName, resourceGroup, err)
	}
	if err = updateFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for updating of Spring Cloud App %q (Spring Cloud Service %q / Resource Group %q): %+v", name, serviceName, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, serviceName, name, "")
	if err != nil {
		return fmt.Errorf("retrieving Spring Cloud App %q (Spring Cloud Service %q / Resource Group %q): %+v", name, serviceName, resourceGroup, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("read Spring Cloud App %q (Spring Cloud Service %q / Resource Group %q) ID", name, serviceName, resourceGroup)
	}
	d.SetId(*resp.ID)

	return resourceArmSpringCloudAppRead(d, meta)
}

func resourceArmSpringCloudAppRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.AppsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SpringCloudAppID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ServiceName, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Spring Cloud App %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading Spring Cloud App %q (Spring Cloud Service %q / Resource Group %q): %+v", id.Name, id.ServiceName, id.ResourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("service_name", id.ServiceName)

	return nil
}

func resourceArmSpringCloudAppDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.AppsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SpringCloudAppID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.ServiceName, id.Name); err != nil {
		return fmt.Errorf("deleting Spring Cloud App %q (Spring Cloud Service %q / Resource Group %q): %+v", id.Name, id.ServiceName, id.ResourceGroup, err)
	}

	return nil
}
