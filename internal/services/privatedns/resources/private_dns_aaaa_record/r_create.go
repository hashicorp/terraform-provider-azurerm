// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package private_dns_aaaa_record

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/privatedns"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/privatezones"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourcePrivateDnsAaaaRecordCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).PrivateDns.RecordSetsClient

	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	privateDNSZoneID, err := privatezones.ParsePrivateDnsZoneID(d.Get("private_dns_zone_id").(string))
	if err != nil {
		return err
	}

	id := privatedns.NewRecordTypeID(privateDNSZoneID.SubscriptionId, privateDNSZoneID.ResourceGroupName, privateDNSZoneID.PrivateDnsZoneName, privatedns.RecordTypeAAAA, d.Get("name").(string))

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.RecordSetsGet(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_private_dns_aaaa_record", id.ID())
			}
		}
	}

	parameters := privatedns.RecordSet{
		Name: pointer.To(id.RelativeRecordSetName),
		Properties: &privatedns.RecordSetProperties{
			Metadata:    tags.Expand(d.Get("tags").(map[string]any)),
			Ttl:         pointer.To(int64(d.Get("ttl").(int))),
			AaaaRecords: expandAzureRmPrivateDnsAaaaRecords(d),
		},
	}

	options := privatedns.RecordSetsCreateOrUpdateOperationOptions{
		IfMatch:     pointer.To(""),
		IfNoneMatch: pointer.To(""),
	}
	if _, err := client.RecordSetsCreateOrUpdate(ctx, id, parameters, options); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if d.IsNewResource() {
		d.SetId(id.ID())
	}

	return resourcePrivateDnsAaaaRecordRead(d, meta)
}

func expandAzureRmPrivateDnsAaaaRecords(d *pluginsdk.ResourceData) *[]privatedns.AaaaRecord {
	recordStrings := d.Get("records").(*pluginsdk.Set).List()
	records := make([]privatedns.AaaaRecord, len(recordStrings))

	for i, v := range recordStrings {
		records[i] = privatedns.AaaaRecord{
			IPv6Address: pointer.To(v.(string)),
		}
	}

	return &records
}
