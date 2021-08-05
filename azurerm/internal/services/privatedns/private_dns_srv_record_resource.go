package privatedns

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/privatedns/mgmt/2018-09-01/privatedns"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/privatedns/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourcePrivateDnsSrvRecord() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourcePrivateDnsSrvRecordCreateUpdate,
		Read:   resourcePrivateDnsSrvRecordRead,
		Update: resourcePrivateDnsSrvRecordCreateUpdate,
		Delete: resourcePrivateDnsSrvRecordDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.SrvRecordID(id)
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
				// lower-cased due to the broken API https://github.com/Azure/azure-rest-api-specs/issues/6641
				ValidateFunc: validate.LowerCasedString,
			},

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
						"priority": {
							Type:         pluginsdk.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 65535),
						},

						"weight": {
							Type:         pluginsdk.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 65535),
						},

						"port": {
							Type:         pluginsdk.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 65535),
						},

						"target": {
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
				ValidateFunc: validation.IntBetween(1, 2147483647),
			},

			"fqdn": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourcePrivateDnsSrvRecordCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).PrivateDns.RecordSetsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceId := parse.NewSrvRecordID(subscriptionId, d.Get("resource_group_name").(string), d.Get("zone_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceId.ResourceGroup, resourceId.PrivateDnsZoneName, privatedns.SRV, resourceId.SRVName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of %s: %+v", resourceId, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_private_dns_srv_record", resourceId.ID())
		}
	}

	ttl := int64(d.Get("ttl").(int))
	t := d.Get("tags").(map[string]interface{})

	parameters := privatedns.RecordSet{
		Name: utils.String(resourceId.SRVName),
		RecordSetProperties: &privatedns.RecordSetProperties{
			Metadata:   tags.Expand(t),
			TTL:        &ttl,
			SrvRecords: expandAzureRmPrivateDnsSrvRecords(d),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, resourceId.ResourceGroup, resourceId.PrivateDnsZoneName, privatedns.SRV, resourceId.SRVName, parameters, "", ""); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", resourceId, err)
	}

	d.SetId(resourceId.ID())
	return resourcePrivateDnsSrvRecordRead(d, meta)
}

func resourcePrivateDnsSrvRecordRead(d *pluginsdk.ResourceData, meta interface{}) error {
	dnsClient := meta.(*clients.Client).PrivateDns.RecordSetsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SrvRecordID(d.Id())
	if err != nil {
		return err
	}

	resp, err := dnsClient.Get(ctx, id.ResourceGroup, id.PrivateDnsZoneName, privatedns.SRV, id.SRVName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.SRVName)
	d.Set("zone_name", id.PrivateDnsZoneName)
	d.Set("resource_group_name", id.ResourceGroup)

	d.Set("ttl", resp.TTL)
	d.Set("fqdn", resp.Fqdn)

	if err := d.Set("record", flattenAzureRmPrivateDnsSrvRecords(resp.SrvRecords)); err != nil {
		return err
	}

	return tags.FlattenAndSet(d, resp.Metadata)
}

func resourcePrivateDnsSrvRecordDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	dnsClient := meta.(*clients.Client).PrivateDns.RecordSetsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SrvRecordID(d.Id())
	if err != nil {
		return err
	}

	if _, err = dnsClient.Delete(ctx, id.ResourceGroup, id.PrivateDnsZoneName, privatedns.SRV, id.SRVName, ""); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func flattenAzureRmPrivateDnsSrvRecords(records *[]privatedns.SrvRecord) []map[string]interface{} {
	results := make([]map[string]interface{}, 0)

	if records != nil {
		for _, record := range *records {
			if record.Priority == nil ||
				record.Weight == nil ||
				record.Port == nil ||
				record.Target == nil {
				continue
			}

			results = append(results, map[string]interface{}{
				"priority": *record.Priority,
				"weight":   *record.Weight,
				"port":     *record.Port,
				"target":   *record.Target,
			})
		}
	}

	return results
}

func expandAzureRmPrivateDnsSrvRecords(d *pluginsdk.ResourceData) *[]privatedns.SrvRecord {
	recordStrings := d.Get("record").(*pluginsdk.Set).List()
	records := make([]privatedns.SrvRecord, len(recordStrings))

	for i, v := range recordStrings {
		record := v.(map[string]interface{})
		priority := int32(record["priority"].(int))
		weight := int32(record["weight"].(int))
		port := int32(record["port"].(int))
		target := record["target"].(string)

		srvRecord := privatedns.SrvRecord{
			Priority: &priority,
			Weight:   &weight,
			Port:     &port,
			Target:   &target,
		}

		records[i] = srvRecord
	}

	return &records
}
