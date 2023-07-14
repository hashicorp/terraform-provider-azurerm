// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dns

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/recordsets"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dns/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/set"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceDnsAAAARecord() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDnsAaaaRecordCreateUpdate,
		Read:   resourceDnsAaaaRecordRead,
		Update: resourceDnsAaaaRecordCreateUpdate,
		Delete: resourceDnsAaaaRecordDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			parsed, err := recordsets.ParseRecordTypeID(id)
			if err != nil {
				return err
			}
			if parsed.RecordType != recordsets.RecordTypeAAAA {
				return fmt.Errorf("this resource only supports 'AAAA' records")
			}
			return nil
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.AAAARecordV0ToV1{},
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"zone_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"records": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.IsIPv6Address,
				},
				Set:           set.HashIPv6Address,
				ConflictsWith: []string{"target_resource_id"},
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

			"target_resource_id": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				ValidateFunc:  azure.ValidateResourceID,
				ConflictsWith: []string{"records"},
			},
		},
	}
}

func resourceDnsAaaaRecordCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Dns.RecordSets
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	defer cancel()

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	zoneName := d.Get("zone_name").(string)

	id := recordsets.NewRecordTypeID(subscriptionId, resGroup, zoneName, recordsets.RecordTypeAAAA, name)
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_dns_aaaa_record", id.ID())
		}
	}

	ttl := int64(d.Get("ttl").(int))
	t := d.Get("tags").(map[string]interface{})
	recordsRaw := d.Get("records").(*pluginsdk.Set).List()
	targetResourceId := d.Get("target_resource_id").(string)

	parameters := recordsets.RecordSet{
		Name: &name,
		Properties: &recordsets.RecordSetProperties{
			Metadata:       tags.Expand(t),
			TTL:            &ttl,
			AAAARecords:    expandAzureRmDnsAaaaRecords(recordsRaw),
			TargetResource: &recordsets.SubResource{},
		},
	}

	if targetResourceId != "" {
		parameters.Properties.TargetResource.Id = utils.String(targetResourceId)
	}

	// TODO: this can be removed when the provider SDK is upgraded
	if targetResourceId == "" && len(recordsRaw) == 0 {
		return fmt.Errorf("One of either `records` or `target_resource_id` must be specified")
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters, recordsets.DefaultCreateOrUpdateOperationOptions()); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceDnsAaaaRecordRead(d, meta)
}

func resourceDnsAaaaRecordRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Dns.RecordSets
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := recordsets.ParseRecordTypeID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.RelativeRecordSetName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("zone_name", id.DnsZoneName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("fqdn", props.Fqdn)
			d.Set("ttl", props.TTL)

			if err := d.Set("records", flattenAzureRmDnsAaaaRecords(props.AAAARecords)); err != nil {
				return fmt.Errorf("setting `records`: %+v", err)
			}

			targetResourceId := ""
			if props.TargetResource != nil && props.TargetResource.Id != nil {
				targetResourceId = *props.TargetResource.Id
			}
			d.Set("target_resource_id", targetResourceId)

			if err := tags.FlattenAndSet(d, props.Metadata); err != nil {
				return err
			}
		}
	}

	return nil
}

func resourceDnsAaaaRecordDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Dns.RecordSets
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := recordsets.ParseRecordTypeID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id, recordsets.DefaultDeleteOperationOptions()); err != nil {
		return fmt.Errorf("deleting %s: %+v", id.RelativeRecordSetName, err)
	}

	return nil
}

func expandAzureRmDnsAaaaRecords(input []interface{}) *[]recordsets.AaaaRecord {
	records := make([]recordsets.AaaaRecord, len(input))

	for i, v := range input {
		ipv6 := NormalizeIPv6Address(v)
		records[i] = recordsets.AaaaRecord{
			IPv6Address: &ipv6,
		}
	}

	return &records
}

func flattenAzureRmDnsAaaaRecords(records *[]recordsets.AaaaRecord) []string {
	if records == nil {
		return []string{}
	}

	results := make([]string, 0)
	for _, record := range *records {
		if record.IPv6Address == nil {
			continue
		}

		results = append(results, NormalizeIPv6Address(*record.IPv6Address))
	}
	return results
}
