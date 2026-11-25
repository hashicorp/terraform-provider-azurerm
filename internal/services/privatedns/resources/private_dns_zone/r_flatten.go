// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package private_dns_zone

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/privatedns"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/privatezones"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourcePrivateDnsZoneFlatten(ctx context.Context, d *pluginsdk.ResourceData, meta any, id *privatezones.PrivateDnsZoneId, privateZone *privatezones.PrivateZone, includeResource bool) error {
	recordSetsClient := meta.(*clients.Client).PrivateDns.RecordSetsClient

	d.Set("name", id.PrivateDnsZoneName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if privateZone != nil {
		if props := privateZone.Properties; props != nil {
			d.Set("number_of_record_sets", props.NumberOfRecordSets)
			d.Set("max_number_of_record_sets", props.MaxNumberOfRecordSets)
			d.Set("max_number_of_virtual_network_links", props.MaxNumberOfVirtualNetworkLinks)
			d.Set("max_number_of_virtual_network_links_with_registration", props.MaxNumberOfVirtualNetworkLinksWithRegistration)
		}

		if err := tags.FlattenAndSet(d, privateZone.Tags); err != nil {
			return err
		}

		if includeResource {
			recordSetID := privatedns.NewRecordTypeID(id.SubscriptionId, id.ResourceGroupName, id.PrivateDnsZoneName, privatedns.RecordTypeSOA, "@")
			resp, err := recordSetsClient.RecordSetsGet(ctx, recordSetID)
			if err != nil {
				return fmt.Errorf("retrieving %s: %v", recordSetID, err)
			}

			if err := d.Set("soa_record", flattenPrivateDNSZoneSOARecord(resp.Model)); err != nil {
				return fmt.Errorf("setting `soa_record`: %+v", err)
			}
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

func flattenPrivateDNSZoneSOARecord(input *privatedns.RecordSet) []any {
	if input == nil || input.Properties == nil {
		return make([]any, 0)
	}

	ttl := 0
	if input.Properties.Ttl != nil {
		ttl = int(*input.Properties.Ttl)
	}

	metaData := make(map[string]any)
	if input.Properties.Metadata != nil {
		metaData = tags.Flatten(input.Properties.Metadata)
	}

	fqdn := pointer.From(input.Properties.Fqdn)

	email := ""
	hostName := ""
	expireTime := 0
	minimumTTL := 0
	refreshTime := 0
	retryTime := 0
	serialNumber := 0
	if record := input.Properties.SoaRecord; record != nil {
		if record.Email != nil {
			email = *record.Email
		}

		if record.Host != nil {
			hostName = *record.Host
		}

		if record.ExpireTime != nil {
			expireTime = int(*record.ExpireTime)
		}

		if record.MinimumTtl != nil {
			minimumTTL = int(*record.MinimumTtl)
		}

		if record.RefreshTime != nil {
			refreshTime = int(*record.RefreshTime)
		}

		if record.RetryTime != nil {
			retryTime = int(*record.RetryTime)
		}

		if record.SerialNumber != nil {
			serialNumber = int(*record.SerialNumber)
		}
	}

	return []any{
		map[string]any{
			"email":         email,
			"host_name":     hostName,
			"expire_time":   expireTime,
			"minimum_ttl":   minimumTTL,
			"refresh_time":  refreshTime,
			"retry_time":    retryTime,
			"serial_number": serialNumber,
			"ttl":           ttl,
			"tags":          metaData,
			"fqdn":          fqdn,
		},
	}
}
