// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dns

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/recordsets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/zones"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dns/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dns/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceDnsZone() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDnsZoneCreateUpdate,
		Read:   resourceDnsZoneRead,
		Update: resourceDnsZoneCreateUpdate,
		Delete: resourceDnsZoneDelete,

		SchemaVersion: 2,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.DnsZoneV0ToV1{},
			1: migration.DnsZoneV1ToV2{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := zones.ParseDnsZoneID(id)
			return err
		}),
		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"number_of_record_sets": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"max_number_of_record_sets": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"name_servers": {
				Type:     pluginsdk.TypeSet,
				Computed: true,
				Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
				Set:      pluginsdk.HashString,
			},

			"soa_record": {
				Type:     pluginsdk.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				//ForceNew: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"email": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validate.DnsZoneSOARecordEmail,
						},

						"host_name": {
							Type:     pluginsdk.TypeString,
							Optional: !features.FourPointOhBeta(), // (@jackofallops) - This should not be set or updatable to meet API design, see https://learn.microsoft.com/en-us/azure/dns/dns-zones-records#soa-records
							Computed: true,
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
							Default:      300,
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

						"serial_number": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Default:      1,
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
					},
				},
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceDnsZoneCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Dns.Zones
	recordSetsClient := meta.(*clients.Client).Dns.RecordSets
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := zones.NewDnsZoneID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_dns_zone", id.ID())
		}
	}

	t := d.Get("tags").(map[string]interface{})

	parameters := zones.Zone{
		Location: location.Normalize("global"),
		Tags:     tags.Expand(t),
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters, zones.DefaultCreateOrUpdateOperationOptions()); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if v, ok := d.GetOk("soa_record"); ok {
		soaRecord := v.([]interface{})[0].(map[string]interface{})

		soaRecordID := recordsets.NewRecordTypeID(id.SubscriptionId, id.ResourceGroupName, id.DnsZoneName, recordsets.RecordTypeSOA, "@")
		soaRecordResp, err := recordSetsClient.Get(ctx, soaRecordID)
		if err != nil {
			return fmt.Errorf("retrieving %s to update SOA: %+v", id, err)
		}

		props := soaRecordResp.Model.Properties
		if props == nil || props.SOARecord == nil {
			return fmt.Errorf("could not read SOA properties for %s", id)
		}

		inputSOARecord := expandArmDNSZoneSOARecord(soaRecord)

		inputSOARecord.Host = props.SOARecord.Host

		rsParameters := recordsets.RecordSet{
			Properties: &recordsets.RecordSetProperties{
				TTL:       utils.Int64(int64(soaRecord["ttl"].(int))),
				Metadata:  tags.Expand(soaRecord["tags"].(map[string]interface{})),
				SOARecord: inputSOARecord,
			},
		}

		if len(id.DnsZoneName+strings.TrimSuffix(*rsParameters.Properties.SOARecord.Email, ".")) > 253 {
			return fmt.Errorf("`email` which is concatenated with DNS Zone `name` cannot exceed 253 characters excluding a trailing period")
		}

		soaRecordId := recordsets.NewRecordTypeID(id.SubscriptionId, id.ResourceGroupName, id.DnsZoneName, recordsets.RecordTypeSOA, "@")
		if _, err := recordSetsClient.CreateOrUpdate(ctx, soaRecordId, rsParameters, recordsets.DefaultCreateOrUpdateOperationOptions()); err != nil {
			return fmt.Errorf("creating/updating %s: %+v", soaRecordId, err)
		}
	}

	d.SetId(id.ID())

	return resourceDnsZoneRead(d, meta)
}

func resourceDnsZoneRead(d *pluginsdk.ResourceData, meta interface{}) error {
	zonesClient := meta.(*clients.Client).Dns.Zones
	recordSetsClient := meta.(*clients.Client).Dns.RecordSets
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := zones.ParseDnsZoneID(d.Id())
	if err != nil {
		return err
	}

	resp, err := zonesClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	soaRecord := recordsets.NewRecordTypeID(id.SubscriptionId, id.ResourceGroupName, id.DnsZoneName, recordsets.RecordTypeSOA, "@")
	soaRecordResp, err := recordSetsClient.Get(ctx, soaRecord)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if err := d.Set("soa_record", flattenArmDNSZoneSOARecord(soaRecordResp.Model)); err != nil {
		return fmt.Errorf("setting `soa_record`: %+v", err)
	}

	d.Set("name", id.DnsZoneName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("number_of_record_sets", props.NumberOfRecordSets)
			d.Set("max_number_of_record_sets", props.MaxNumberOfRecordSets)

			nameServers := make([]string, 0)
			if s := props.NameServers; s != nil {
				nameServers = *s
			}
			if err := d.Set("name_servers", nameServers); err != nil {
				return err
			}
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}

func resourceDnsZoneDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Dns.Zones
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := zones.ParseDnsZoneID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id, zones.DefaultDeleteOperationOptions()); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandArmDNSZoneSOARecord(input map[string]interface{}) *recordsets.SoaRecord {
	result := &recordsets.SoaRecord{
		Email:        utils.String(input["email"].(string)),
		ExpireTime:   utils.Int64(int64(input["expire_time"].(int))),
		MinimumTTL:   utils.Int64(int64(input["minimum_ttl"].(int))),
		RefreshTime:  utils.Int64(int64(input["refresh_time"].(int))),
		RetryTime:    utils.Int64(int64(input["retry_time"].(int))),
		SerialNumber: utils.Int64(int64(input["serial_number"].(int))),
	}

	if !features.FourPointOhBeta() && input["host_name"].(string) != "" {
		result.Host = pointer.To(input["host_name"].(string))
	}

	return result
}

func flattenArmDNSZoneSOARecord(input *recordsets.RecordSet) []interface{} {
	output := make([]interface{}, 0)
	if input != nil {
		if props := input.Properties; props != nil {
			ttl := 0
			if props.TTL != nil {
				ttl = int(*props.TTL)
			}

			metaData := make(map[string]interface{})
			if props.Metadata != nil {
				metaData = tags.Flatten(props.Metadata)
			}

			fqdn := ""
			if props.Fqdn != nil {
				fqdn = *props.Fqdn
			}

			email := ""
			hostName := ""
			expireTime := 0
			minimumTTL := 0
			refreshTime := 0
			retryTime := 0
			serialNumber := 0
			if record := props.SOARecord; record != nil {
				if record.Email != nil {
					email = *record.Email
				}

				if record.Host != nil {
					hostName = *record.Host
				}

				if record.ExpireTime != nil {
					expireTime = int(*record.ExpireTime)
				}

				if record.MinimumTTL != nil {
					minimumTTL = int(*record.MinimumTTL)
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

			output = append(output, map[string]interface{}{
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
			})
		}
	}

	return output
}
