package azurerm

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func dataSourceArmVirtualMachines() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmVMsRead,
		Schema: map[string]*schema.Schema{
			"resource_group_name": resourceGroupNameForDataSourceSchema(),

			"name_prefix": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceArmVMsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).vmClient
	ctx := meta.(*ArmClient).StopContext

	resGroup := d.Get("resource_group_name").(string)
	namePrefix := d.Get("name_prefix").(string)

	resp, err := client.List(ctx, resGroup)
	if err != nil {
		return fmt.Errorf("Error listing Virtual Machines in the Resource Group %q: %v", resGroup, err)
	}

	var ids []string
	var names []string

	for _, vm := range resp.Values() {
		if strings.HasPrefix(*vm.Name, namePrefix) {
			ids = append(ids, *vm.ID)
			names = append(names, *vm.Name)
		}
	}

	d.SetId(time.Now().UTC().String())

	if err := d.Set("ids", ids); err != nil {
		return fmt.Errorf("Error setting `ids`: %+v", err)
	}
	if err := d.Set("names", names); err != nil {
		return fmt.Errorf("Error setting `names`: %+v", err)
	}

	return nil
}
