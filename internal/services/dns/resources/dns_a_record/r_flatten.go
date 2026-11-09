// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dns_a_record

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/recordsets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceDnsARecordFlatten(d *pluginsdk.ResourceData, id *recordsets.RecordTypeId, model *recordsets.RecordSet) error {
	d.Set("name", id.RelativeRecordSetName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("zone_name", id.DnsZoneName)

	if model != nil {
		if props := model.Properties; props != nil {
			d.Set("fqdn", props.Fqdn)
			d.Set("ttl", props.TTL)

			if err := d.Set("records", flattenAzureRmDnsARecords(props.ARecords)); err != nil {
				return fmt.Errorf("setting `records`: %+v", err)
			}

			targetResourceId := ""
			if props.TargetResource != nil && props.TargetResource.Id != nil {
				targetResourceId = *props.TargetResource.Id
			}
			d.Set("target_resource_id", targetResourceId)

			if err := tags.FlattenAndSet(d, props.Metadata); err != nil {
				return err
			}
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

func flattenAzureRmDnsARecords(records *[]recordsets.ARecord) []string {
	if records == nil {
		return []string{}
	}

	results := make([]string, 0)
	for _, record := range *records {
		if record.IPv4Address == nil {
			continue
		}

		results = append(results, *record.IPv4Address)
	}

	return results
}
