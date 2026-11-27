// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dns_mx_record

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
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

func resourceDnsMxRecordCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := dnsClient.From(meta).RecordSets
	ctx, cancel := timeouts.ForCreateUpdate(account.StopContext(meta), d)
	subscriptionId := account.From(meta).SubscriptionId
	defer cancel()

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	zoneName := d.Get("zone_name").(string)

	id := recordsets.NewRecordTypeID(subscriptionId, resGroup, zoneName, recordsets.RecordTypeMX, name)
	if d.IsNewResource() {
		if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_dns_mx_record", id.ID())
			}
		}
	}

	ttl := int64(d.Get("ttl").(int))
	t := d.Get("tags").(map[string]any)

	parameters := recordsets.RecordSet{
		Name: &name,
		Properties: &recordsets.RecordSetProperties{
			Metadata:  tags.Expand(t),
			TTL:       &ttl,
			MXRecords: expandAzureRmDnsMxRecords(d),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters, recordsets.DefaultCreateOrUpdateOperationOptions()); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceDnsMxRecordRead(d, meta)
}

// expand creates an array of dns.MxRecord, that is, the array needed
// by azure-sdk-for-go to manipulate azure resources, hence Preference
// is an int32
func expandAzureRmDnsMxRecords(d *pluginsdk.ResourceData) *[]recordsets.MxRecord {
	recordStrings := d.Get("record").(*pluginsdk.Set).List()
	records := make([]recordsets.MxRecord, len(recordStrings))

	for i, v := range recordStrings {
		mxrecord := v.(map[string]any)
		preference := mxrecord["preference"].(string)
		i64, _ := strconv.ParseInt(preference, 10, 32)

		records[i] = recordsets.MxRecord{
			Preference: &i64,
			Exchange:   pointer.To(mxrecord["exchange"].(string)),
		}
	}

	return &records
}
