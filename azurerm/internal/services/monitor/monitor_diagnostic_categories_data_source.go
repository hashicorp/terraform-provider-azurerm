package monitor

import (
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2019-06-01/insights"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
)

func dataSourceMonitorDiagnosticCategories() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMonitorDiagnosticCategoriesRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"logs": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Computed: true,
			},

			"metrics": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Computed: true,
			},
		},
	}
}

func dataSourceMonitorDiagnosticCategoriesRead(d *schema.ResourceData, meta interface{}) error {
	categoriesClient := meta.(*clients.Client).Monitor.DiagnosticSettingsCategoryClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	actualResourceId := d.Get("resource_id").(string)
	// trim off the leading `/` since the CheckExistenceByID / List methods don't expect it
	resourceId := strings.TrimPrefix(actualResourceId, "/")

	// then retrieve the possible Diagnostics Categories for this Resource
	categories, err := categoriesClient.List(ctx, resourceId)
	if err != nil {
		return fmt.Errorf("Error retrieving Diagnostics Categories for Resource %q: %+v", actualResourceId, err)
	}

	if categories.Value == nil {
		return fmt.Errorf("Error retrieving Diagnostics Categories for Resource %q: `categories.Value` was nil", actualResourceId)
	}

	d.SetId(actualResourceId)
	val := *categories.Value

	metrics := make([]string, 0)
	logs := make([]string, 0)

	for _, v := range val {
		if v.Name == nil {
			continue
		}

		if category := v.DiagnosticSettingsCategory; category != nil {
			switch category.CategoryType {
			case insights.Logs:
				logs = append(logs, *v.Name)
			case insights.Metrics:
				metrics = append(metrics, *v.Name)
			default:
				return fmt.Errorf("Unsupported category type %q", string(category.CategoryType))
			}
		}
	}

	if err := d.Set("logs", logs); err != nil {
		return fmt.Errorf("Error setting `logs`: %+v", err)
	}

	if err := d.Set("metrics", metrics); err != nil {
		return fmt.Errorf("Error setting `metrics`: %+v", err)
	}

	return nil
}
