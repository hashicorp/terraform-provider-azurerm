// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package privatedns

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/recordsets"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/privatedns/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourcePrivateDnsARecord() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourcePrivateDnsARecordCreateUpdate,
		Read:   resourcePrivateDnsARecordRead,
		Update: resourcePrivateDnsARecordCreateUpdate,
		Delete: resourcePrivateDnsARecordDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			resourceId, err := recordsets.ParseRecordTypeID(id)
			if err != nil {
				return err
			}
			if resourceId.RecordType != recordsets.RecordTypeA {
				return fmt.Errorf("importing %s wrong type received: expected %s received %s", id, recordsets.RecordTypeA, resourceId.RecordType)
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
				Required: true,
				ForceNew: true,
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

			"records": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				MaxItems: 20,
				Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
				Set:      pluginsdk.HashString,
			},

			"ttl": {
				Type:     pluginsdk.TypeInt,
				Required: true,
			},

			"fqdn": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourcePrivateDnsARecordCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).PrivateDns.RecordSetsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := recordsets.NewRecordTypeID(subscriptionId, d.Get("resource_group_name").(string), d.Get("zone_name").(string), recordsets.RecordTypeA, d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_private_dns_a_record", id.ID())
		}
	}

	parameters := recordsets.RecordSet{
		Name: utils.String(id.RelativeRecordSetName),
		Properties: &recordsets.RecordSetProperties{
			Metadata: tags.Expand(d.Get("tags").(map[string]interface{})),
			Ttl:      utils.Int64(int64(d.Get("ttl").(int))),
			ARecords: expandAzureRmPrivateDnsARecords(d),
		},
	}

	options := recordsets.CreateOrUpdateOperationOptions{
		IfMatch:     utils.String(""),
		IfNoneMatch: utils.String(""),
	}
	if _, err := client.CreateOrUpdate(ctx, id, parameters, options); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourcePrivateDnsARecordRead(d, meta)
}

func resourcePrivateDnsARecordRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
	d.Set("zone_name", id.PrivateDnsZoneName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("ttl", props.Ttl)
			d.Set("fqdn", props.Fqdn)

			if err := d.Set("records", flattenAzureRmPrivateDnsARecords(props.ARecords)); err != nil {
				return err
			}

			return tags.FlattenAndSet(d, props.Metadata)
		}
	}

	return nil
}

func resourcePrivateDnsARecordDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	dnsClient := meta.(*clients.Client).PrivateDns.RecordSetsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := recordsets.ParseRecordTypeID(d.Id())
	if err != nil {
		return err
	}

	options := recordsets.DeleteOperationOptions{IfMatch: utils.String("")}

	if _, err := dnsClient.Delete(ctx, *id, options); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func flattenAzureRmPrivateDnsARecords(records *[]recordsets.ARecord) []string {
	results := make([]string, 0)
	if records == nil {
		return results
	}

	for _, record := range *records {
		if record.IPv4Address == nil {
			continue
		}

		results = append(results, *record.IPv4Address)
	}

	return results
}

func expandAzureRmPrivateDnsARecords(d *pluginsdk.ResourceData) *[]recordsets.ARecord {
	recordStrings := d.Get("records").(*pluginsdk.Set).List()
	records := make([]recordsets.ARecord, len(recordStrings))

	for i, v := range recordStrings {
		ipv4 := v.(string)
		records[i] = recordsets.ARecord{
			IPv4Address: &ipv4,
		}
	}

	return &records
}
