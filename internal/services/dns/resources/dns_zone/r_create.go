// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dns_zone

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/recordsets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/zones"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Dns.Zones
			recordSetsClient := metadata.Client.Dns.RecordSets
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model DnsZoneResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := zones.NewDnsZoneID(subscriptionId, model.ResourceGroupName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
					}
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			parameters := zones.Zone{
				Location: location.Normalize("global"),
				Tags:     pointer.To(model.Tags),
			}

			if _, err := client.CreateOrUpdate(ctx, id, parameters, zones.DefaultCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)

			if len(model.SoaRecord) == 1 {
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
					return fmt.Errorf("creating %s: %+v", soaRecordId, err)
				}
			}

			return nil
		},
	}
}

func expandDNSZoneSOARecord(input DnsZoneSoaRecordResourceRecord) *recordsets.SoaRecord {
	return &recordsets.SoaRecord{
		Email:        pointer.To(input.Email),
		ExpireTime:   pointer.To(input.ExpireTime),
		MinimumTTL:   pointer.To(input.MinimumTtl),
		RefreshTime:  pointer.To(input.RefreshTime),
		RetryTime:    pointer.To(input.RetryTime),
		SerialNumber: pointer.To(input.SerialNumber),
	}
}
