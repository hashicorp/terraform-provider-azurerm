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

func dataSourceArmArmSpringCloudDeployment() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmArmSpringCloudDeploymentRead,

		Schema: map[string]*schema.Schema{
			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"name": {
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

			"spring_cloud_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: appplatform.ValidateSpringCloudName,
			},

			"cpu": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"memory_in_gb": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"instance_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"jvm_options": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"env": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"runtime_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceArmArmSpringCloudDeploymentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).AppPlatform.DeploymentsClient
	ctx, cancel := timeouts.ForRead(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	springCloudName := id.Path["Spring"]
	appName := id.Path["apps"]
	name := id.Path["deployments"]

	resp, err := client.Get(ctx, resourceGroup, springCloudName, appName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Spring Cloud Deployments %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Spring Cloud Deployment %q (Spring Cloud Service %q / App Name %q / Resource Group %q): %+v", name, springCloudName, appName, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("spring_cloud_name", springCloudName)
	d.Set("spring_cloud_app_name", appName)
	if deploymentSettings := resp.Properties.DeploymentSettings; deploymentSettings != nil {
		d.Set("cpu", deploymentSettings.CPU)
		d.Set("memory_in_gb", deploymentSettings.MemoryInGB)
		d.Set("jvm_options", deploymentSettings.JvmOptions)
		d.Set("instance_count", deploymentSettings.InstanceCount)
		d.Set("runtime_version", deploymentSettings.RuntimeVersion)
		if err := d.Set("env", flattenSpringCloudDeploymentEnv(deploymentSettings.EnvironmentVariables)); err != nil {
			return fmt.Errorf("Error setting `env`: %+v", err)
		}
	}

	return nil
}
