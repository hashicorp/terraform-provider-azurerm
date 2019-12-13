package azurerm

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	azappplatform "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/appplatform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmSpringCloudApp() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmSpringCloudAppRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azappplatform.ValidateSpringCloudName,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"spring_cloud_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azappplatform.ValidateSpringCloudName,
			},

			"active_deployment_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"created_time": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"persistent_disk": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mount_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size_in_gb": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},

			"public": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"temporary_disk": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mount_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size_in_gb": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},

			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceArmSpringCloudAppRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).AppPlatform.AppsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	springCloudName := d.Get("spring_cloud_name").(string)

	resp, err := client.Get(ctx, resourceGroup, springCloudName, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Spring Cloud App %q (Spring Cloud Service %q / Resource Group %q) was not found", name, springCloudName, resourceGroup)
		}
		return fmt.Errorf("Error reading Spring Cloud App %q (Spring Cloud Service %q / Resource Group %q): %+v", name, springCloudName, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("spring_cloud_name", springCloudName)
	if appResourceProperties := resp.Properties; appResourceProperties != nil {
		d.Set("active_deployment_name", appResourceProperties.ActiveDeploymentName)
		d.Set("created_time", (appResourceProperties.CreatedTime).String())
		if err := d.Set("persistent_disk", flattenArmSpringCloudAppPersistentDisk(appResourceProperties.PersistentDisk)); err != nil {
			return fmt.Errorf("Error setting `persistent_disk`: %+v", err)
		}
		d.Set("public", appResourceProperties.Public)
		if err := d.Set("temporary_disk", flattenArmSpringCloudAppTemporaryDisk(appResourceProperties.TemporaryDisk)); err != nil {
			return fmt.Errorf("Error setting `temporary_disk`: %+v", err)
		}
		d.Set("url", appResourceProperties.URL)
	}

	return nil
}
