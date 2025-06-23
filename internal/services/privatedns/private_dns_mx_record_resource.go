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
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourcegroups"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/recordsets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/virtualnetworklinks"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/privatedns/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourcePrivateDnsMxRecord() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourcePrivateDnsMxRecordCreateUpdate,
		Read:   resourcePrivateDnsMxRecordRead,
		Update: resourcePrivateDnsMxRecordCreateUpdate,
		Delete: resourcePrivateDnsMxRecordDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			resourceId, err := recordsets.ParseRecordTypeID(id)
			if err != nil {
				return err
			}
			if resourceId.RecordType != recordsets.RecordTypeMX {
				return fmt.Errorf("importing %s wrong type received: expected %s received %s", id, recordsets.RecordTypeMX, resourceId.RecordType)
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

			"private_zone_id": {
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
	if !features.FivePointOh() {
		resource.Schema["private_zone_id"] = &pluginsdk.Schema{
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Computed:      true,
			ForceNew:      true,
			ValidateFunc:  validation.StringIsNotEmpty,
			ConflictsWith: []string{"zone_name", "resource_group_name"},
			AtLeastOneOf:  []string{"zone_name", "resource_group_name", "private_zone_id"},
		}
		// TODO: in 4.0 make `name` case sensitive and replace `resource_group_name` and `zone_name` with `private_zone_id`
		// TODO: make this case sensitive once the API's fixed https://github.com/Azure/azure-rest-api-specs/issues/6641
		resource.Schema["resource_group_name"] = &pluginsdk.Schema{
			Type:             pluginsdk.TypeString,
			Optional:         true,
			Computed:         true,
			ForceNew:         true,
			DiffSuppressFunc: suppress.CaseDifference,
			ValidateFunc:     resourcegroups.ValidateName,
			Deprecated:       "The `resource_group_name` field is deprecated in favor of `private_zone_id`. This will be removed in version 5.0.",
			ConflictsWith:    []string{"private_zone_id"},
			AtLeastOneOf:     []string{"private_zone_id", "zone_name", "resource_group_name"},
		}
		resource.Schema["zone_name"] = &pluginsdk.Schema{
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Computed:      true,
			ForceNew:      true,
			ValidateFunc:  validation.StringIsNotEmpty,
			Deprecated:    "The `zone_name` field is deprecated in favor of `private_zone_id`. This will be removed in version 5.0.",
			ConflictsWith: []string{"private_zone_id"},
			AtLeastOneOf:  []string{"private_zone_id", "zone_name", "resource_group_name"},
		}
	}

	return resource
}

func resourcePrivateDnsMxRecordCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).PrivateDns.RecordSetsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	rawDnsZoneId := d.Get("private_zone_id").(string)
	if !features.FivePointOh() && rawDnsZoneId == "" {
		dnsZoneId := &recordsets.PrivateDnsZoneId{
			ResourceGroupName:  d.Get("resource_group_name").(string),
			PrivateDnsZoneName: d.Get("zone_name").(string),
			SubscriptionId:     subscriptionId,
		}
		rawDnsZoneId = dnsZoneId.ID()
	}
	dnsZoneId, err := virtualnetworklinks.ParsePrivateDnsZoneID(rawDnsZoneId)
	if err != nil {
		return fmt.Errorf("parsing private DNS zone ID %q: %+v", rawDnsZoneId, err)
	}
	id := recordsets.NewRecordTypeID(subscriptionId, dnsZoneId.ResourceGroupName, dnsZoneId.PrivateDnsZoneName, recordsets.RecordTypeMX, d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_private_dns_mx_record", id.ID())
		}
	}

	parameters := recordsets.RecordSet{
		Name: pointer.To(id.RelativeRecordSetName),
		Properties: &recordsets.RecordSetProperties{
			Metadata:  tags.Expand(d.Get("tags").(map[string]interface{})),
			Ttl:       pointer.To(int64(d.Get("ttl").(int))),
			MxRecords: expandAzureRmPrivateDnsMxRecords(d),
		},
	}

	options := recordsets.CreateOrUpdateOperationOptions{
		IfMatch:     pointer.To(""),
		IfNoneMatch: pointer.To(""),
	}
	if _, err := client.CreateOrUpdate(ctx, id, parameters, options); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourcePrivateDnsMxRecordRead(d, meta)
}

func resourcePrivateDnsMxRecordRead(d *pluginsdk.ResourceData, meta interface{}) error {
	dnsClient := meta.(*clients.Client).PrivateDns.RecordSetsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := recordsets.ParseRecordTypeID(d.Id())
	if err != nil {
		return err
	}

	resp, err := dnsClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.RelativeRecordSetName)
	dnsZoneId := &recordsets.PrivateDnsZoneId{
		ResourceGroupName:  id.ResourceGroupName,
		PrivateDnsZoneName: id.PrivateDnsZoneName,
		SubscriptionId:     meta.(*clients.Client).Account.SubscriptionId,
	}
	d.Set("private_zone_id", dnsZoneId.ID())
	if !features.FivePointOh() {
		d.Set("zone_name", id.PrivateDnsZoneName)
		d.Set("resource_group_name", id.ResourceGroupName)
	}

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

	id, err := recordsets.ParseRecordTypeID(d.Id())
	if err != nil {
		return err
	}

	options := recordsets.DeleteOperationOptions{IfMatch: pointer.To("")}

	if _, err = dnsClient.Delete(ctx, *id, options); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func flattenAzureRmPrivateDnsMxRecords(records *[]recordsets.MxRecord) []map[string]interface{} {
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

func expandAzureRmPrivateDnsMxRecords(d *pluginsdk.ResourceData) *[]recordsets.MxRecord {
	recordStrings := d.Get("record").(*pluginsdk.Set).List()
	records := make([]recordsets.MxRecord, len(recordStrings))

	for i, v := range recordStrings {
		if v == nil {
			continue
		}
		record := v.(map[string]interface{})
		mxRecord := recordsets.MxRecord{
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
