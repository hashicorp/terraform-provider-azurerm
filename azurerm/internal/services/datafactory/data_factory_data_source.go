package datafactory

import (
	"fmt"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/datafactory/mgmt/2018-06-01/datafactory"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceDataFactory() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDataFactoryRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`^[A-Za-z0-9]+(?:-[A-Za-z0-9]+)*$`),
					`Invalid name for Data Factory, see https://docs.microsoft.com/en-us/azure/data-factory/naming-rules`,
				),
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"location": azure.SchemaLocationForDataSource(),

			"identity": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"principal_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"github_configuration": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"branch_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"git_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"repository_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"root_folder": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"vsts_configuration": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"branch_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"repository_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"root_folder": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceDataFactoryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.FactoriesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error Data Factory %q (Resource Group %q) was not found", name, resourceGroup)
		}

		return fmt.Errorf("Error retrieving Data Factory %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("API returns a nil/empty id on Data Factory %q (resource group %q): %+v", name, resourceGroup, err)
	}
	d.SetId(*resp.ID)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	d.Set("vsts_configuration", []interface{}{})
	d.Set("github_configuration", []interface{}{})
	repoType, repo := flattenDataFactoryRepoConfiguration(&resp)
	if repoType == datafactory.TypeFactoryVSTSConfiguration {
		if err := d.Set("vsts_configuration", repo); err != nil {
			return fmt.Errorf("Error setting `vsts_configuration`: %+v", err)
		}
	}
	if repoType == datafactory.TypeFactoryGitHubConfiguration {
		if err := d.Set("github_configuration", repo); err != nil {
			return fmt.Errorf("Error setting `github_configuration`: %+v", err)
		}
	}
	if repoType == datafactory.TypeFactoryRepoConfiguration {
		d.Set("vsts_configuration", repo)
		d.Set("github_configuration", repo)
	}

	if err := d.Set("identity", flattenDataFactoryIdentity(resp.Identity)); err != nil {
		return fmt.Errorf("Error flattening `identity`: %+v", err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
