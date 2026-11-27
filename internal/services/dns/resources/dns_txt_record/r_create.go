// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dns_txt_record

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/recordsets"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	dnsClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/dns/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDnsTxtRecordCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := dnsClient.From(meta).RecordSets
	ctx, cancel := timeouts.ForCreateUpdate(account.StopContext(meta), d)
	subscriptionId := account.From(meta).SubscriptionId
	defer cancel()

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	zoneName := d.Get("zone_name").(string)

	id := recordsets.NewRecordTypeID(subscriptionId, resGroup, zoneName, recordsets.RecordTypeTXT, name)
	if d.IsNewResource() {
		if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError(azurermDnsTxtRecordResourceName, id.ID())
			}
		}
	}

	ttl := int64(d.Get("ttl").(int))
	t := d.Get("tags").(map[string]any)

	parameters := recordsets.RecordSet{
		Name: &name,
		Properties: &recordsets.RecordSetProperties{
			Metadata:   tags.Expand(t),
			TTL:        &ttl,
			TXTRecords: expandAzureRmDnsTxtRecords(d),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters, recordsets.DefaultCreateOrUpdateOperationOptions()); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceDnsTxtRecordRead(d, meta)
}

func expandAzureRmDnsTxtRecords(d *pluginsdk.ResourceData) *[]recordsets.TxtRecord {
	recordStrings := d.Get("record").(*pluginsdk.Set).List()
	records := make([]recordsets.TxtRecord, len(recordStrings))

	segmentLen := 254
	for i, v := range recordStrings {
		record := v.(map[string]any)
		v := record["value"].(string)

		var value []string
		for len(v) > segmentLen {
			value = append(value, v[:segmentLen])
			v = v[segmentLen:]
		}
		value = append(value, v)

		records[i] = recordsets.TxtRecord{
			Value: &value,
		}
	}

	return &records
}
