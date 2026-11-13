// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package private_dns_txt_record

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/privatedns"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/privatezones"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourcePrivateDnsTxtRecordRead(d *pluginsdk.ResourceData, meta any) error {
	dnsClient := meta.(*clients.Client).PrivateDns.RecordSetsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := privatedns.ParseRecordTypeID(d.Id())
	if err != nil {
		return err
	}

	resp, err := dnsClient.RecordSetsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.RelativeRecordSetName)
	d.Set("private_dns_zone_id", privatezones.NewPrivateDnsZoneID(id.SubscriptionId, id.ResourceGroupName, id.PrivateDnsZoneName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("ttl", props.Ttl)
			d.Set("fqdn", props.Fqdn)

			if err := d.Set("record", flattenAzureRmPrivateDnsTxtRecords(props.TxtRecords)); err != nil {
				return fmt.Errorf("setting `record`: %+v", err)
			}

			if err := tags.FlattenAndSet(d, props.Metadata); err != nil {
				return err
			}
		}
	}

	return nil
}

func flattenAzureRmPrivateDnsTxtRecords(records *[]privatedns.TxtRecord) []map[string]any {
	results := make([]map[string]any, 0)

	if records != nil {
		for _, record := range *records {
			txtRecord := make(map[string]any)

			if v := record.Value; v != nil {
				txtRecord["value"] = strings.Join(*v, "")
			}

			results = append(results, txtRecord)
		}
	}

	return results
}
