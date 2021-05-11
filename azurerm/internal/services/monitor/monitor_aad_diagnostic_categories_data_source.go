package monitor

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"time"
)

func dataSourceMonitorAADDiagnosticCategories() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMonitorAADDiagnosticCategoriesRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"logs": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
		},
	}
}

func dataSourceMonitorAADDiagnosticCategoriesRead(d *schema.ResourceData, meta interface{}) error {
	categoriesClient := meta.(*clients.Client).Monitor.AADDiagnosticSettingsCategoryClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// then retrieve the possible Diagnostics Categories for AAD
	categories, err := categoriesClient.List(ctx)
	if err != nil {
		return fmt.Errorf("Error retrieving Diagnostics Categories for AAD: %+v", err)
	}

	if categories.Value == nil {
		return fmt.Errorf("Error retrieving Diagnostics Categories for AAD: `categories.Value` was nil")
	}

	d.SetId("/providers/microsoft.aadiam/diagnosticSettingsCategories")
	val := *categories.Value

	logs := make([]interface{}, 0)

	for _, v := range val {
		if v.Name == nil {
			continue
		}

		if category := v.DiagnosticSettingsCategory; category != nil {
			logs = append(logs, *v.Name)
		}
	}

	if err := d.Set("logs", logs); err != nil {
		return fmt.Errorf("Error setting `logs`: %+v", err)
	}

	return nil
}
