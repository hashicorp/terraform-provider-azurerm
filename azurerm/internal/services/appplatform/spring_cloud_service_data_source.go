package appplatform

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/appplatform/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmSpringCloudService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmSpringCloudServiceRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.SpringCloudServiceName,
			},

			"location": azure.SchemaLocationForDataSource(),

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"config_server_git_setting": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uri": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"label": {
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

						"http_basic_auth": DataSourceSchemaConfigServerHttpBasicAuth(),

						"ssh_auth": DataSourceSchemaConfigServerSSHAuth(),

						"repository": {
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
									"label": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"pattern": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"search_paths": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},

									"http_basic_auth": DataSourceSchemaConfigServerHttpBasicAuth(),

									"ssh_auth": DataSourceSchemaConfigServerSSHAuth(),
								},
							},
						},
					},
				},
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceArmSpringCloudServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.ServicesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Spring Cloud %q (Resource Group %q) was not found", name, resourceGroup)
		}
		return fmt.Errorf("Error reading Spring Cloud %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("Error retrieving Spring Cloud Service %q (Resource Group %q): ID was nil or empty", name, resourceGroup)
	}

	d.SetId(*resp.ID)

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if resp.Properties != nil && resp.Properties.ConfigServerProperties != nil && resp.Properties.ConfigServerProperties.ConfigServer != nil {
		if props := resp.Properties.ConfigServerProperties.ConfigServer.GitProperty; props != nil {
			if err := d.Set("config_server_git_setting", flattenArmSpringCloudConfigServerGitProperty(props, d)); err != nil {
				return fmt.Errorf("failure setting AzureRM Spring Cloud Service Config Server error: %+v", err)
			}
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
