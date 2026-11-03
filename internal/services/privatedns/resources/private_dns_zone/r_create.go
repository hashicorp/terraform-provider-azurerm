// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package private_dns_zone

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/privatedns"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/privatezones"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourcePrivateDnsZoneCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).PrivateDns.PrivateZonesClient
	recordSetsClient := meta.(*clients.Client).PrivateDns.RecordSetsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := privatezones.NewPrivateDnsZoneID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError(privateDnsZoneResourceName, id.ID())
			}
		}
	}

	parameters := privatezones.PrivateZone{
		Location: pointer.To("global"),
		Tags:     tags.Expand(d.Get("tags").(map[string]any)),
	}

	options := privatezones.CreateOrUpdateOperationOptions{
		IfMatch:     pointer.To(""),
		IfNoneMatch: pointer.To(""),
	}

	if d.IsNewResource() {
		if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, parameters, options, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
			return fmt.Errorf("creating %s: %s", id, err)
		}
		d.SetId(id.ID())
		if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
			return err
		}
	} else {
		if err := client.CreateOrUpdateThenPoll(ctx, id, parameters, options); err != nil {
			return fmt.Errorf("updating %s: %s", id, err)
		}
	}

	if v, ok := d.GetOk("soa_record"); ok {
		soaRecordRaw := v.([]any)[0].(map[string]any)
		soaRecord := expandPrivateDNSZoneSOARecord(soaRecordRaw)
		rsParameters := privatedns.RecordSet{
			Properties: &privatedns.RecordSetProperties{
				Ttl:       pointer.To(int64(soaRecordRaw["ttl"].(int))),
				Metadata:  tags.Expand(soaRecordRaw["tags"].(map[string]any)),
				SoaRecord: soaRecord,
			},
		}

		recordId := privatedns.NewRecordTypeID(id.SubscriptionId, id.ResourceGroupName, id.PrivateDnsZoneName, privatedns.RecordTypeSOA, "@")

		val := fmt.Sprintf("%s%s", id.PrivateDnsZoneName, strings.TrimSuffix(*soaRecord.Email, "."))
		if len(val) > 253 {
			return fmt.Errorf("the value %q for `email` which is concatenated with Private DNS Zone `name` cannot exceed 253 characters excluding a trailing period", val)
		}

		createOptions := privatedns.RecordSetsCreateOrUpdateOperationOptions{
			IfMatch:     pointer.To(""),
			IfNoneMatch: pointer.To(""),
		}

		if _, err := recordSetsClient.RecordSetsCreateOrUpdate(ctx, recordId, rsParameters, createOptions); err != nil {
			return fmt.Errorf("creating/updating %s: %s", recordId, err)
		}
	}

	return resourcePrivateDnsZoneRead(d, meta)
}

func expandPrivateDNSZoneSOARecord(input map[string]any) *privatedns.SoaRecord {
	return &privatedns.SoaRecord{
		Email:       pointer.To(input["email"].(string)),
		ExpireTime:  pointer.To(int64(input["expire_time"].(int))),
		MinimumTtl:  pointer.To(int64(input["minimum_ttl"].(int))),
		RefreshTime: pointer.To(int64(input["refresh_time"].(int))),
		RetryTime:   pointer.To(int64(input["retry_time"].(int))),
	}
}
