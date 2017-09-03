package azurerm

import (
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceArmPublicIPIds() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmPublicIPIdsRead,

		Schema: map[string]*schema.Schema{
			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"minimum_count": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"ids": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceArmPublicIPIdsRead(d *schema.ResourceData, meta interface{}) error {
	publicIPClient := meta.(*ArmClient).publicIPClient

	resGroup := d.Get("resource_group_name").(string)
	minimumCount, minimumCountOk := d.GetOk("minimum_count")
	resp, err := publicIPClient.List(resGroup)
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			d.SetId("")
		}
		return fmt.Errorf("Error making Read request on Azure resource group %s: %s", resGroup, err)
	}
	availableIds := make([]string, 0)
	for _, element := range *resp.Value {
		if element.IPConfiguration == nil {
			availableIds = append(availableIds, *element.ID)
		}
	}

	if minimumCountOk && len(availableIds) < minimumCount.(int) {
		return fmt.Errorf("Not enough unassigned public IP addresses in resource group %s", resGroup)
	}
	d.SetId(time.Now().UTC().String())
	d.Set("ids", availableIds)

	return nil
}
