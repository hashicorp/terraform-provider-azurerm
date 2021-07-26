package loganalytics

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

func logAnalyticsStorageInsightsImporter(d *pluginsdk.ResourceData, meta interface{}) ([]*pluginsdk.ResourceData, error) {
	if _, err := parse.LogAnalyticsStorageInsightsID(d.Id()); err != nil {
		return []*pluginsdk.ResourceData{d}, err
	}

	if v, ok := d.GetOk("storage_account_key"); ok && v.(string) != "" {
		d.Set("storage_account_key", v)
	}

	return []*pluginsdk.ResourceData{d}, nil
}
