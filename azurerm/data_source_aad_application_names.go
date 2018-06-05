package azurerm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func dataSourceArmAadApplicationNames() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmAadApplicationReadNames,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),

			"applications": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceArmAadApplicationReadNames(d *schema.ResourceData, meta interface{}) error {
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

	d.SetId(fmt.Sprintf("%d", hashcode.String(oDataFilter)))
	var applicationIds []string

	for _, element := range listResult.Values() {
		applicationIds = append(applicationIds, *element.DisplayName)
	}

	d.Set("applications", applicationIds)

	return nil
}
