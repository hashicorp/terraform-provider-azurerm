package azurerm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmRedisCache() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmRedisCacheRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": locationForDataSourceSchema(),

			"resource_group_name": resourceGroupNameForDataSourceSchema(),

			"hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"ssl_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"primary_access_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_access_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"tags": tagsForDataSourceSchema(),
		},
	}
}

func dataSourceArmRedisCacheRead(d *schema.ResourceData, meta interface{}) error {
	ctx := meta.(*ArmClient).StopContext
	client := meta.(*ArmClient).redisClient

	resourceGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Redis instance %q (Resource group %q) was not found", name, resourceGroup)
		}
		return fmt.Errorf("Error reading the state of Redis instance %q: %+v", name, err)
	}

	d.SetId(*resp.ID)

	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	d.Set("host_name", *resp.HostName)

	d.Set("sku_name", *resp.Sku)

	if *resp.EnableNonSslPort == false {
		d.Set("port", *resp.Port)
	}

	d.Set("ssl_port", *resp.SslPort)

	keys, err := client.ListKeys(ctx, resourceGroup, name)
	if err != nil {
		return err
	}

	d.Set("primary_access_key", *keys.PrimaryKey)
	d.Set("secondary_access_key", *keys.SecondaryKey)

	flattenAndSetTags(d, resp.Tags)

	return nil

}
