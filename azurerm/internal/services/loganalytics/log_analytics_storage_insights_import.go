package loganalytics

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/parse"
)

func logAnalyticsStorageInsightsImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if _, err := parse.LogAnalyticsStorageInsightsID(d.Id()); err != nil {
		return []*schema.ResourceData{d}, err
	}

	if v, ok := d.GetOk("storage_account_key"); ok && v.(string) != "" {
		d.Set("storage_account_key", v)
	}

	return []*schema.ResourceData{d}, nil
}
