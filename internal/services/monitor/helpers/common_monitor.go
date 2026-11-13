// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2018-04-16/scheduledqueryrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func ExpandMonitorScheduledQueryRulesCommonSource(d *pluginsdk.ResourceData) scheduledqueryrules.Source {
	authorizedResourceIDs := d.Get("authorized_resource_ids").(*pluginsdk.Set).List()
	dataSourceID := d.Get("data_source_id").(string)

	source := scheduledqueryrules.Source{
		AuthorizedResources: pluginsdk.ExpandStringSlice(authorizedResourceIDs),
		DataSourceId:        dataSourceID,
	}

	if query, ok := d.GetOk("query"); ok {
		source.Query = pointer.To(query.(string))
	}
	if queryType, ok := d.GetOk("query_type"); ok {
		source.QueryType = pointer.ToEnum[scheduledqueryrules.QueryType](queryType.(string))
	}

	return source
}

func ExpandStringValues(input []any) []string {
	result := make([]string, 0)
	for _, item := range input {
		if item != nil {
			result = append(result, item.(string))
		} else {
			result = append(result, "")
		}
	}
	return result
}
