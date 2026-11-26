// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dns_zone

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/recordsets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/zones"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Dns.Zones
			recordSetsClient := metadata.Client.Dns.RecordSets

			var model DnsZoneResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id, err := zones.ParseDnsZoneID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}

			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}

			payload := zones.Zone{
				Location: existing.Model.Location,
				Tags:     existing.Model.Tags,
			}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = pointer.To(model.Tags)
			}

			if _, err := client.CreateOrUpdate(ctx, *id, payload, zones.DefaultCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			if metadata.ResourceData.HasChange("soa_record") && len(model.SoaRecord) == 1 {
				soaRecordID := recordsets.NewRecordTypeID(id.SubscriptionId, id.ResourceGroupName, id.DnsZoneName, recordsets.RecordTypeSOA, "@")
				soaRecordResp, err := recordSetsClient.Get(ctx, soaRecordID)
				if err != nil {
					return fmt.Errorf("retrieving %s to update SOA: %+v", id, err)
				}

				props := soaRecordResp.Model.Properties
				if props == nil || props.SOARecord == nil {
					return fmt.Errorf("could not read SOA properties for %s", id)
				}

				inputSOARecord := expandDNSZoneSOARecord(model.SoaRecord[0])

				inputSOARecord.Host = props.SOARecord.Host

				rsParameters := recordsets.RecordSet{
					Properties: &recordsets.RecordSetProperties{
						TTL:       pointer.To(model.SoaRecord[0].Ttl),
						Metadata:  pointer.To(model.SoaRecord[0].Tags),
						SOARecord: inputSOARecord,
					},
				}

				if len(id.DnsZoneName+strings.TrimSuffix(*rsParameters.Properties.SOARecord.Email, ".")) > 253 {
					return fmt.Errorf("`email` which is concatenated with DNS Zone `name` cannot exceed 253 characters excluding a trailing period")
				}

				soaRecordId := recordsets.NewRecordTypeID(id.SubscriptionId, id.ResourceGroupName, id.DnsZoneName, recordsets.RecordTypeSOA, "@")
				if _, err := recordSetsClient.CreateOrUpdate(ctx, soaRecordId, rsParameters, recordsets.DefaultCreateOrUpdateOperationOptions()); err != nil {
					return fmt.Errorf("updating %s: %+v", soaRecordId, err)
				}
			}

			return nil
		},
	}
}
