// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dns_txt_record

import (
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/recordsets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceDnsTxtRecordFlatten(d *pluginsdk.ResourceData, id *recordsets.RecordTypeId, model *recordsets.RecordSet) error {
	d.Set("name", id.RelativeRecordSetName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("zone_name", id.DnsZoneName)

	if model != nil {
		if props := model.Properties; props != nil {
			d.Set("ttl", props.TTL)
			d.Set("fqdn", props.Fqdn)

			if err := d.Set("record", flattenAzureRmDnsTxtRecords(props.TXTRecords)); err != nil {
				return err
			}
			if err := tags.FlattenAndSet(d, props.Metadata); err != nil {
				return err
			}
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

func flattenAzureRmDnsTxtRecords(records *[]recordsets.TxtRecord) []map[string]any {
	results := make([]map[string]any, 0)

	if records != nil {
		for _, record := range *records {
			value := ""
			if v := record.Value; v != nil {
				value = strings.Join(*v, "")
			}

			results = append(results, map[string]any{
				"value": value,
			})
		}
	}

	return results
}
