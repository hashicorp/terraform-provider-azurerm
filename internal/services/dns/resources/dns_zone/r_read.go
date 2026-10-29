// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dns_zone

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/recordsets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/zones"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			zonesClient := metadata.Client.Dns.Zones
			recordSetsClient := metadata.Client.Dns.RecordSets

			state := DnsZoneResourceModel{}

			id, err := zones.ParseDnsZoneID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := zonesClient.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			soaRecord := recordsets.NewRecordTypeID(id.SubscriptionId, id.ResourceGroupName, id.DnsZoneName, recordsets.RecordTypeSOA, "@")
			soaRecordResp, err := recordSetsClient.Get(ctx, soaRecord)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state.SoaRecord = flattenDNSZoneSOARecord(soaRecordResp.Model)

			state.Name = id.DnsZoneName
			state.ResourceGroupName = id.ResourceGroupName

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.NumberOfRecordSets = pointer.From(props.NumberOfRecordSets)
					state.MaxNumberOfRecordSets = pointer.From(props.MaxNumberOfRecordSets)
					state.NameServers = pointer.From(props.NameServers)
				}
				state.Tags = pointer.From(model.Tags)
			}

			return metadata.Encode(&state)
		},
	}
}

func flattenDNSZoneSOARecord(input *recordsets.RecordSet) []DnsZoneSoaRecordResourceRecord {
	output := make([]DnsZoneSoaRecordResourceRecord, 0)
	if input != nil {
		if props := input.Properties; props != nil {
			result := DnsZoneSoaRecordResourceRecord{
				Ttl:  pointer.From(props.TTL),
				Tags: pointer.From(props.Metadata),
				Fqdn: pointer.From(props.Fqdn),
			}

			if record := props.SOARecord; record != nil {
				result.Email = pointer.From(record.Email)
				result.HostName = pointer.From(record.Host)
				result.ExpireTime = pointer.From(record.ExpireTime)
				result.MinimumTtl = pointer.From(record.MinimumTTL)
				result.RefreshTime = pointer.From(record.RefreshTime)
				result.RetryTime = pointer.From(record.RetryTime)
				result.SerialNumber = pointer.From(record.SerialNumber)
			}

			output = append(output, result)
		}
	}

	return output
}

type DnsZoneSoaRecordResourceRecord struct {
	Email        string            `tfschema:"email"`
	ExpireTime   int64             `tfschema:"expire_time"`
	MinimumTtl   int64             `tfschema:"minimum_ttl"`
	RefreshTime  int64             `tfschema:"refresh_time"`
	RetryTime    int64             `tfschema:"retry_time"`
	SerialNumber int64             `tfschema:"serial_number"`
	Ttl          int64             `tfschema:"ttl"`
	Fqdn         string            `tfschema:"fqdn"`
	HostName     string            `tfschema:"host_name"`
	Tags         map[string]string `tfschema:"tags"`
}
