// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package privatedns

import (
	"bytes"
	"fmt"
	"math"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/privatedns"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/privatedns/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourcePrivateDnsMxRecord() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourcePrivateDnsMxRecordCreateUpdate,
		Read:   resourcePrivateDnsMxRecordRead,
		Update: resourcePrivateDnsMxRecordCreateUpdate,
		Delete: resourcePrivateDnsMxRecordDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			resourceId, err := privatedns.ParseRecordTypeID(id)
			if err != nil {
				return err
			}
			if resourceId.RecordType != privatedns.RecordTypeMX {
				return fmt.Errorf("importing %s wrong type received: expected %s received %s", id, privatedns.RecordTypeMX, resourceId.RecordType)
			}
			return nil
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				ForceNew: true,
				// name is optional and defaults to root zone (@) if not set
				Optional: true,
				Default:  "@",
				// lower-cased due to the broken API https://github.com/Azure/azure-rest-api-specs/issues/6641
				ValidateFunc: validate.LowerCasedString,
			},

			// TODO: in 4.0 make `name` case sensitive and replace `resource_group_name` and `zone_name` with `private_zone_id`

			// TODO: make this case sensitive once the API's fixed https://github.com/Azure/azure-rest-api-specs/issues/6641
			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"zone_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"record": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"preference": {
							Type:     pluginsdk.TypeInt,
							Required: true,
							// 16 bit uint (rfc 974)
							ValidateFunc: validation.IntBetween(0, 65535),
						},
						"exchange": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"ttl": {
				Type:         pluginsdk.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, math.MaxInt32),
			},

			"fqdn": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourcePrivateDnsMxRecordCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).PrivateDns.RecordSetsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := privatedns.NewRecordTypeID(subscriptionId, d.Get("resource_group_name").(string), d.Get("zone_name").(string), privatedns.RecordTypeMX, d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.RecordSetsGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_private_dns_mx_record", id.ID())
		}
	}

	parameters := privatedns.RecordSet{
		Name: pointer.To(id.RelativeRecordSetName),
		Properties: &privatedns.RecordSetProperties{
			Metadata:  tags.Expand(d.Get("tags").(map[string]interface{})),
			Ttl:       pointer.To(int64(d.Get("ttl").(int))),
			MxRecords: expandAzureRmPrivateDnsMxRecords(d),
		},
	}

	options := privatedns.RecordSetsCreateOrUpdateOperationOptions{
		IfMatch:     pointer.To(""),
		IfNoneMatch: pointer.To(""),
	}
	if _, err := client.RecordSetsCreateOrUpdate(ctx, id, parameters, options); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourcePrivateDnsMxRecordRead(d, meta)
}

func resourcePrivateDnsMxRecordRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
	d.Set("zone_name", id.PrivateDnsZoneName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("ttl", props.Ttl)
			d.Set("fqdn", props.Fqdn)

			if err := d.Set("record", flattenAzureRmPrivateDnsMxRecords(props.MxRecords)); err != nil {
				return err
			}

			return tags.FlattenAndSet(d, props.Metadata)
		}
	}

	return nil
}

func resourcePrivateDnsMxRecordDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	dnsClient := meta.(*clients.Client).PrivateDns.RecordSetsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := privatedns.ParseRecordTypeID(d.Id())
	if err != nil {
		return err
	}

	options := privatedns.RecordSetsDeleteOperationOptions{IfMatch: pointer.To("")}

	if _, err = dnsClient.RecordSetsDelete(ctx, *id, options); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func flattenAzureRmPrivateDnsMxRecords(records *[]privatedns.MxRecord) []map[string]interface{} {
	results := make([]map[string]interface{}, 0)

	if records != nil {
		for _, record := range *records {
			if record.Preference == nil ||
				record.Exchange == nil {
				continue
			}

			results = append(results, map[string]interface{}{
				"preference": *record.Preference,
				"exchange":   *record.Exchange,
			})
		}
	}

	return results
}

func expandAzureRmPrivateDnsMxRecords(d *pluginsdk.ResourceData) *[]privatedns.MxRecord {
	recordStrings := d.Get("record").(*pluginsdk.Set).List()
	records := make([]privatedns.MxRecord, len(recordStrings))

	for i, v := range recordStrings {
		if v == nil {
			continue
		}
		record := v.(map[string]interface{})
		mxRecord := privatedns.MxRecord{
			Preference: pointer.To(int64(record["preference"].(int))),
			Exchange:   pointer.To(record["exchange"].(string)),
		}

		records[i] = mxRecord
	}

	return &records
}

func resourcePrivateDnsMxRecordHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(fmt.Sprintf("%d-", m["preference"].(int)))
		buf.WriteString(fmt.Sprintf("%s-", m["exchange"].(string)))
	}

	return pluginsdk.HashString(buf.String())
}
