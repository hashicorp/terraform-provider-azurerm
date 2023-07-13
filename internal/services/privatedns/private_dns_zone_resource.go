// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package privatedns

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2020-06-01/privatezones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2020-06-01/recordsets"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/privatedns/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourcePrivateDnsZone() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourcePrivateDnsZoneCreateUpdate,
		Read:   resourcePrivateDnsZoneRead,
		Update: resourcePrivateDnsZoneCreateUpdate,
		Delete: resourcePrivateDnsZoneDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := privatezones.ParsePrivateDnsZoneID(id)
			return err
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
			},

			// TODO: this can become case-sensitive with a state migration
			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"number_of_record_sets": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"max_number_of_record_sets": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"max_number_of_virtual_network_links": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"max_number_of_virtual_network_links_with_registration": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"soa_record": {
				Type:     pluginsdk.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"email": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validate.PrivateDnsZoneSOARecordEmail,
						},

						"expire_time": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Default:      2419200,
							ValidateFunc: validation.IntAtLeast(0),
						},

						"minimum_ttl": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Default:      10,
							ValidateFunc: validation.IntAtLeast(0),
						},

						"refresh_time": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Default:      3600,
							ValidateFunc: validation.IntAtLeast(0),
						},

						"retry_time": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Default:      300,
							ValidateFunc: validation.IntAtLeast(0),
						},

						"ttl": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Default:      3600,
							ValidateFunc: validation.IntBetween(0, 2147483647),
						},

						"tags": commonschema.Tags(),

						"fqdn": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						// This field should be able to be updated since DNS Record Sets API allows to update it.
						// So the issue is submitted on https://github.com/Azure/azure-rest-api-specs/issues/11674
						// Once the issue is fixed, the field will be updated to `Required` property.
						"host_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						// This field should be able to be updated since DNS Record Sets API allows to update it.
						// So the issue is submitted on https://github.com/Azure/azure-rest-api-specs/issues/11674
						// Once the issue is fixed, the field will be updated to `Optional` property with `Default` attribute.
						"serial_number": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},
					},
				},
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourcePrivateDnsZoneCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).PrivateDns.PrivateZonesClient
	recordSetsClient := meta.(*clients.Client).PrivateDns.RecordSetsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := privatezones.NewPrivateDnsZoneID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_private_dns_zone", id.ID())
		}
	}

	parameters := privatezones.PrivateZone{
		Location: utils.String("global"),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	options := privatezones.CreateOrUpdateOperationOptions{
		IfMatch:     utils.String(""),
		IfNoneMatch: utils.String(""),
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, parameters, options); err != nil {
		return fmt.Errorf("creating/updating %s: %s", id, err)
	}

	if v, ok := d.GetOk("soa_record"); ok {
		soaRecordRaw := v.([]interface{})[0].(map[string]interface{})
		soaRecord := expandPrivateDNSZoneSOARecord(soaRecordRaw)
		rsParameters := recordsets.RecordSet{
			Properties: &recordsets.RecordSetProperties{
				Ttl:       utils.Int64(int64(soaRecordRaw["ttl"].(int))),
				Metadata:  tags.Expand(soaRecordRaw["tags"].(map[string]interface{})),
				SoaRecord: soaRecord,
			},
		}

		recordId := recordsets.NewRecordTypeID(id.SubscriptionId, id.ResourceGroupName, id.PrivateDnsZoneName, recordsets.RecordTypeSOA, "@")

		val := fmt.Sprintf("%s%s", id.PrivateDnsZoneName, strings.TrimSuffix(*soaRecord.Email, "."))
		if len(val) > 253 {
			return fmt.Errorf("the value %q for `email` which is concatenated with Private DNS Zone `name` cannot exceed 253 characters excluding a trailing period", val)
		}

		createOptions := recordsets.CreateOrUpdateOperationOptions{
			IfMatch:     utils.String(""),
			IfNoneMatch: utils.String(""),
		}

		if _, err := recordSetsClient.CreateOrUpdate(ctx, recordId, rsParameters, createOptions); err != nil {
			return fmt.Errorf("creating/updating %s: %s", recordId, err)
		}
	}

	d.SetId(id.ID())
	return resourcePrivateDnsZoneRead(d, meta)
}

func resourcePrivateDnsZoneRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).PrivateDns.PrivateZonesClient
	recordSetsClient := meta.(*clients.Client).PrivateDns.RecordSetsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := privatezones.ParsePrivateDnsZoneID(d.Id())
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

	recordId := recordsets.NewRecordTypeID(id.SubscriptionId, id.ResourceGroupName, id.PrivateDnsZoneName, recordsets.RecordTypeSOA, "@")
	recordSetResp, err := recordSetsClient.Get(ctx, recordId)
	if err != nil {
		return fmt.Errorf("reading DNS SOA record @: %v", err)
	}

	d.Set("name", id.PrivateDnsZoneName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("number_of_record_sets", props.NumberOfRecordSets)
			d.Set("max_number_of_record_sets", props.MaxNumberOfRecordSets)
			d.Set("max_number_of_virtual_network_links", props.MaxNumberOfVirtualNetworkLinks)
			d.Set("max_number_of_virtual_network_links_with_registration", props.MaxNumberOfVirtualNetworkLinksWithRegistration)
		}

		if err = tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	if err = d.Set("soa_record", flattenPrivateDNSZoneSOARecord(recordSetResp.Model)); err != nil {
		return fmt.Errorf("setting `soa_record`: %+v", err)
	}

	return nil
}

func resourcePrivateDnsZoneDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).PrivateDns.PrivateZonesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := privatezones.ParsePrivateDnsZoneID(d.Id())
	if err != nil {
		return err
	}

	options := privatezones.DeleteOperationOptions{IfMatch: utils.String("")}

	if err = client.DeleteThenPoll(ctx, *id, options); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandPrivateDNSZoneSOARecord(input map[string]interface{}) *recordsets.SoaRecord {
	return &recordsets.SoaRecord{
		Email:       utils.String(input["email"].(string)),
		ExpireTime:  utils.Int64(int64(input["expire_time"].(int))),
		MinimumTtl:  utils.Int64(int64(input["minimum_ttl"].(int))),
		RefreshTime: utils.Int64(int64(input["refresh_time"].(int))),
		RetryTime:   utils.Int64(int64(input["retry_time"].(int))),
	}
}

func flattenPrivateDNSZoneSOARecord(input *recordsets.RecordSet) []interface{} {
	if input == nil || input.Properties == nil {
		return make([]interface{}, 0)
	}

	ttl := 0
	if input.Properties.Ttl != nil {
		ttl = int(*input.Properties.Ttl)
	}

	metaData := make(map[string]interface{})
	if input.Properties.Metadata != nil {
		metaData = tags.Flatten(input.Properties.Metadata)
	}

	fqdn := ""
	if input.Properties.Fqdn != nil {
		fqdn = *input.Properties.Fqdn
	}

	email := ""
	hostName := ""
	expireTime := 0
	minimumTTL := 0
	refreshTime := 0
	retryTime := 0
	serialNumber := 0
	if record := input.Properties.SoaRecord; record != nil {
		if record.Email != nil {
			email = *record.Email
		}

		if record.Host != nil {
			hostName = *record.Host
		}

		if record.ExpireTime != nil {
			expireTime = int(*record.ExpireTime)
		}

		if record.MinimumTtl != nil {
			minimumTTL = int(*record.MinimumTtl)
		}

		if record.RefreshTime != nil {
			refreshTime = int(*record.RefreshTime)
		}

		if record.RetryTime != nil {
			retryTime = int(*record.RetryTime)
		}

		if record.SerialNumber != nil {
			serialNumber = int(*record.SerialNumber)
		}
	}

	return []interface{}{
		map[string]interface{}{
			"email":         email,
			"host_name":     hostName,
			"expire_time":   expireTime,
			"minimum_ttl":   minimumTTL,
			"refresh_time":  refreshTime,
			"retry_time":    retryTime,
			"serial_number": serialNumber,
			"ttl":           ttl,
			"tags":          metaData,
			"fqdn":          fqdn,
		},
	}
}
