package automation

import (
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

func dataSourceAutomationVariableDateTime() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceAutomationVariableDateTimeRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: datasourceAutomationVariableCommonSchema(pluginsdk.TypeString),
	}
}

func dataSourceAutomationVariableDateTimeRead(d *pluginsdk.ResourceData, meta interface{}) error {
	return dataSourceAutomationVariableRead(d, meta, "Datetime")
}
