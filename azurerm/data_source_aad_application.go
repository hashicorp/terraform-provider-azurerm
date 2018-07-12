package azurerm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func dataSourceArmAadApplication() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmAadApplicationRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"object_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"application_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"homepage": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceArmAadApplicationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).applicationsClient
	ctx := meta.(*ArmClient).StopContext

	filters, filtersOk := d.GetOk("filter")

	oDataFilter := ""

	if filtersOk {
		oDataFilter = dataSourceAzureFilterBuilder(filters.(*schema.Set))
	}

	log.Printf("[DEBUG] Here we are: %s", oDataFilter)

	listResult, listErr := client.List(ctx, oDataFilter)

	if listErr != nil {
		return fmt.Errorf("Error listing Service Principals: %#v", listErr)
	}

	if listResult.Values() == nil {
		return fmt.Errorf("Unexpected Service Principal query result: %#v", listResult.Values())
	}

	application := listResult.Values()[0]

	d.SetId(fmt.Sprintf("%s", hashcode.String(oDataFilter)))

	d.Set("homepage", application.Homepage)
	d.Set("object_id", application.ObjectID)
	d.Set("application_id", application.AppID)

	return nil
}
