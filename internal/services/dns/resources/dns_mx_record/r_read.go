// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dns_mx_record

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/recordsets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDnsMxRecordRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Dns.RecordSets
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := recordsets.ParseRecordTypeID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.RelativeRecordSetName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("zone_name", id.DnsZoneName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("ttl", props.TTL)
			d.Set("fqdn", props.Fqdn)

			if err := d.Set("record", flattenAzureRmDnsMxRecords(props.MXRecords)); err != nil {
				return err
			}
			if err := tags.FlattenAndSet(d, props.Metadata); err != nil {
				return err
			}
		}
	}

	return nil
}

// flatten creates an array of map where preference is a string to suit
// the expectations of the ResourceData schema, so that this data can be
// managed by Terradata state.
func flattenAzureRmDnsMxRecords(records *[]recordsets.MxRecord) []map[string]any {
	results := make([]map[string]any, 0)

	if records != nil {
		for _, record := range *records {
			// TODO: convert this to use an int64
			preference := ""
			if record.Preference != nil {
				preference = strconv.Itoa(int(*record.Preference))
			}

			results = append(results, map[string]any{
				"preference": preference,
				"exchange":   pointer.From(record.Exchange),
			})
		}
	}

	return results
}

func resourceDnsMxRecordHash(v any) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]any); ok {
		fmt.Fprintf(&buf, "%s-", m["preference"].(string))
		fmt.Fprintf(&buf, "%s-", m["exchange"].(string))
	}

	return pluginsdk.HashString(buf.String())
}
