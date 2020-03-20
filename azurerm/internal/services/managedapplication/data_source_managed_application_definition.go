package managedapplication

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/managedapplication/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmManagedApplicationDefinition() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmManagedApplicationDefinitionRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.ManagedApplicationDefinitionName,
			},

			"location": azure.SchemaLocationForDataSource(),

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"authorization": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role_definition_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_principal_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"create_ui_definition": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"package_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"lock_level": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"main_template": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"package_file_uri": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceArmManagedApplicationDefinitionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ManagedApplication.ApplicationDefinitionClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Failed: Managed Application Definition (Managed Application Definition Name %q / Resource Group %q) was not found", name, resourceGroup)
		}
		return fmt.Errorf("Failure reading Managed Application Definition (Managed Application Definition Name %q / Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if props := resp.ApplicationDefinitionProperties; props != nil {
		if err := d.Set("authorization", flattenArmManagedApplicationDefinitionAuthorization(props.Authorizations)); err != nil {
			return fmt.Errorf("Failure setting `authorization`: %+v", err)
		}
		d.Set("description", props.Description)
		d.Set("display_name", props.DisplayName)
		d.Set("package_enabled", props.IsEnabled)
		d.Set("lock_level", string(props.LockLevel))
	}
	if v, ok := d.GetOk("create_ui_definition"); ok {
		d.Set("create_ui_definition", v.(string))
	}
	if v, ok := d.GetOk("main_template"); ok {
		d.Set("main_template", v.(string))
	}
	if v, ok := d.GetOk("package_file_uri"); ok {
		d.Set("package_file_uri", v.(string))
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("API returns a nil/empty id on Managed Application Definition %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	d.SetId(*resp.ID)

	return tags.FlattenAndSet(d, resp.Tags)
}
