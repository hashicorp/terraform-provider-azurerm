package springcloud

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/appplatform/mgmt/2019-05-01-preview/appplatform"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/springcloud/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/springcloud/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceSpringCloudJavaDeployment() *schema.Resource {
	return &schema.Resource{
		Create: resourceSpringCloudJavaDeploymentCreateUpdate,
		Read:   resourceSpringCloudJavaDeploymentRead,
		Update: resourceSpringCloudJavaDeploymentCreateUpdate,
		Delete: resourceSpringCloudJavaDeploymentDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.SpringCloudDeploymentID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SpringCloudDeploymentName,
			},

			"spring_cloud_app_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SpringCloudAppID,
			},

			"cpu": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 4),
			},

			"environment_variables": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"instance_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 500),
			},

			"jvm_options": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"memory_in_gb": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 8),
			},

			"runtime_version": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(appplatform.Java8),
					string(appplatform.Java11),
				}, false),
				Default: string(appplatform.Java8),
			},
		},
	}
}

func resourceSpringCloudJavaDeploymentCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.DeploymentsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	appId, err := parse.SpringCloudAppID(d.Get("spring_cloud_app_id").(string))
	if err != nil {
		return err
	}

	resourceId := parse.NewSpringCloudDeploymentID(subscriptionId, appId.ResourceGroup, appId.SpringName, appId.AppName, name).ID()
	if d.IsNewResource() {
		existing, err := client.Get(ctx, appId.ResourceGroup, appId.SpringName, appId.AppName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Spring Cloud Deployment %q (Spring Cloud Service %q / App %q / Resource Group %q): %+v", name, appId.SpringName, appId.AppName, appId.ResourceGroup, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_spring_cloud_java_deployment", resourceId)
		}
	}

	deployment := appplatform.DeploymentResource{
		Properties: &appplatform.DeploymentResourceProperties{
			Source: &appplatform.UserSourceInfo{
				Type:         appplatform.Jar,
				RelativePath: utils.String("<default>"),
			},
			DeploymentSettings: &appplatform.DeploymentSettings{
				CPU:                  utils.Int32(int32(d.Get("cpu").(int))),
				MemoryInGB:           utils.Int32(int32(d.Get("memory_in_gb").(int))),
				JvmOptions:           utils.String(d.Get("jvm_options").(string)),
				InstanceCount:        utils.Int32(int32(d.Get("instance_count").(int))),
				EnvironmentVariables: expandSpringCloudDeploymentEnvironmentVariables(d.Get("environment_variables").(map[string]interface{})),
				RuntimeVersion:       appplatform.RuntimeVersion(d.Get("runtime_version").(string)),
			},
		},
	}
	future, err := client.CreateOrUpdate(ctx, appId.ResourceGroup, appId.SpringName, appId.AppName, name, deployment)
	if err != nil {
		return fmt.Errorf("creating/update Spring Cloud Deployment %q (Spring Cloud Service %q / App %q / Resource Group %q): %+v", name, appId.SpringName, appId.AppName, appId.ResourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of Spring Cloud Deployment %q (Spring Cloud Service %q / App %q / Resource Group %q): %+v", name, appId.SpringName, appId.AppName, appId.ResourceGroup, err)
	}

	d.SetId(resourceId)
	return resourceSpringCloudJavaDeploymentRead(d, meta)
}

func resourceSpringCloudJavaDeploymentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.DeploymentsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SpringCloudDeploymentID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.SpringName, id.AppName, id.DeploymentName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Spring Cloud deployment %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading Spring Cloud Deployment %q (Spring Cloud Service %q / App %q / resource Group %q): %+v", id.DeploymentName, id.SpringName, id.AppName, id.ResourceGroup, err)
	}

	d.Set("name", id.DeploymentName)
	d.Set("spring_cloud_app_id", parse.NewSpringCloudAppID(id.SubscriptionId, id.ResourceGroup, id.SpringName, id.AppName).ID())
	if resp.Properties != nil && resp.Properties.DeploymentSettings != nil {
		settings := resp.Properties.DeploymentSettings
		d.Set("cpu", settings.CPU)
		d.Set("memory_in_gb", settings.MemoryInGB)
		d.Set("instance_count", settings.InstanceCount)
		d.Set("jvm_options", settings.JvmOptions)
		d.Set("environment_variables", flattenSpringCloudDeploymentEnvironmentVariables(settings.EnvironmentVariables))
		d.Set("runtime_version", settings.RuntimeVersion)
	}

	return nil
}

func resourceSpringCloudJavaDeploymentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.DeploymentsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SpringCloudDeploymentID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.SpringName, id.AppName, id.DeploymentName); err != nil {
		return fmt.Errorf("deleting Spring Cloud Deployment %q (Spring Cloud Service %q / App %q / resource Group %q): %+v", id.DeploymentName, id.SpringName, id.AppName, id.ResourceGroup, err)
	}

	return nil
}

func expandSpringCloudDeploymentEnvironmentVariables(envMap map[string]interface{}) map[string]*string {
	output := make(map[string]*string, len(envMap))

	for k, v := range envMap {
		output[k] = utils.String(v.(string))
	}

	return output
}

func flattenSpringCloudDeploymentEnvironmentVariables(envMap map[string]*string) map[string]interface{} {
	output := make(map[string]interface{}, len(envMap))
	for i, v := range envMap {
		if v == nil {
			continue
		}
		output[i] = *v
	}
	return output
}
