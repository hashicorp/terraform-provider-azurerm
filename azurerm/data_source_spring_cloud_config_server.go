package azurerm

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmSpringCloudConfigServer() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmSpringCloudConfigServerRead,

		Schema: map[string]*schema.Schema{
			"spring_cloud_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"uri": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"host_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"host_key_algorithm": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"label": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"password": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"private_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"repositories": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uri": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_key_algorithm": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"label": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"password": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"pattern": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"private_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"search_paths": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"strict_host_key_checking": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"username": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"search_paths": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"strict_host_key_checking": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"username": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceArmSpringCloudConfigServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).AppPlatform.ServicesClient
	ctx := meta.(*ArmClient).StopContext

	springCloudId := d.Get("spring_cloud_id").(string)
	id, err := azure.ParseAzureResourceID(springCloudId)
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	springCloudName := id.Path["Spring"]

	resp, err := client.Get(ctx, resourceGroup, springCloudName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Spring Cloud %q (Resource Group %q) was not found", springCloudName, resourceGroup)
		}
		return fmt.Errorf("Error reading Spring Cloud %q (Resource Group %q): %+v", springCloudName, resourceGroup, err)
	}

	d.SetId(springCloudId)

	if resp.Properties != nil && resp.Properties.ConfigServerProperties != nil && resp.Properties.ConfigServerProperties.ConfigServer != nil {
		if props := resp.Properties.ConfigServerProperties.ConfigServer.GitProperty; props != nil {
			d.Set("host_key", props.HostKey)
			d.Set("host_key_algorithm", props.HostKeyAlgorithm)
			d.Set("label", props.Label)
			d.Set("password", props.Password)
			d.Set("private_key", props.PrivateKey)
			d.Set("strict_host_key_checking", props.StrictHostKeyChecking)
			d.Set("uri", props.URI)
			d.Set("username", props.Username)
			d.Set("search_paths", props.SearchPaths)
			if err := d.Set("repositories", flattenArmSpringCloudGitPatternRepository(props.Repositories)); err != nil {
				return fmt.Errorf("Error setting `repositories`: %+v", err)
			}
		}
	}

	return nil
}
