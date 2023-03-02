package dns

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/recordsets"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dns/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDnsTxtRecord() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceDnsTxtRecordCreateUpdate,
		Read:   resourceDnsTxtRecordRead,
		Update: resourceDnsTxtRecordCreateUpdate,
		Delete: resourceDnsTxtRecordDelete,

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
			if parsed.RecordType != recordsets.RecordTypeTXT {
				return fmt.Errorf("this resource only supports 'TXT' records")
			}
			return nil
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.TXTRecordV0ToV1{},
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

			"record": {
				Type:     pluginsdk.TypeList,
				Required: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"values": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Computed: !features.FourPointOhBeta(),
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringLenBetween(1, 255),
							},
						},
					},
				},
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

	if !features.FourPointOhBeta() {
		resource.Schema["record"].Elem.(*pluginsdk.Resource).Schema["value"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.StringLenBetween(1, 1024),
			Deprecated:   "`value` will be removed in favour of the property `values` in version 4.0 of the AzureRM Provider.",
		}
	}

	return resource
}

func resourceDnsTxtRecordCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Dns.RecordSets
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	defer cancel()

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	zoneName := d.Get("zone_name").(string)

	id := recordsets.NewRecordTypeID(subscriptionId, resGroup, zoneName, recordsets.RecordTypeTXT, name)
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_dns_txt_record", id.ID())
		}
	}

	ttl := int64(d.Get("ttl").(int))
	t := d.Get("tags").(map[string]interface{})

	parameters := recordsets.RecordSet{
		Name: &name,
		Properties: &recordsets.RecordSetProperties{
			Metadata: tags.Expand(t),
			TTL:      &ttl,
		},
	}

	records, err := expandAzureRmDnsTxtRecords(d)
	if err != nil {
		return err
	}
	parameters.Properties.TXTRecords = records

	if _, err := client.CreateOrUpdate(ctx, id, parameters, recordsets.DefaultCreateOrUpdateOperationOptions()); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceDnsTxtRecordRead(d, meta)
}

func resourceDnsTxtRecordRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
			d.Set("ttl", props.TTL)
			d.Set("fqdn", props.Fqdn)

			if err := d.Set("record", flattenAzureRmDnsTxtRecords(props.TXTRecords)); err != nil {
				return err
			}
			if err := tags.FlattenAndSet(d, props.Metadata); err != nil {
				return err
			}
		}
	}

	return nil
}

func resourceDnsTxtRecordDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Dns.RecordSets
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := recordsets.ParseRecordTypeID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id, recordsets.DefaultDeleteOperationOptions()); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func flattenAzureRmDnsTxtRecords(records *[]recordsets.TxtRecord) []map[string]interface{} {
	results := make([]map[string]interface{}, 0)

	if records != nil {
		for _, recordItem := range *records {
			record := map[string]interface{}{}

			if recordValue := recordItem.Value; recordValue != nil {
				if !features.FourPointOhBeta() {
					value := strings.Join(*recordValue, "")
					record["value"] = value
				}

				var values []string
				record["values"] = append(values, *recordValue...)
			}

			results = append(results, record)
		}
	}

	return results
}

func expandAzureRmDnsTxtRecords(d *pluginsdk.ResourceData) (*[]recordsets.TxtRecord, error) {
	recordStrings := d.Get("record").([]interface{})
	records := make([]recordsets.TxtRecord, len(recordStrings))

	for i, v := range recordStrings {
		record := v.(map[string]interface{})

		// When `Compute: true` is added, d.GetOk() always returns the value of previous apply. So it has to use `d.GetRawConfig()` to check if the property is set in tf config
		isRecordValueSet := !d.GetRawConfig().AsValueMap()["record"].AsValueSlice()[i].AsValueMap()["value"].IsNull()
		isRecordValuesSet := !d.GetRawConfig().AsValueMap()["record"].AsValueSlice()[i].AsValueMap()["values"].IsNull()

		var values []string

		if !features.FourPointOhBeta() {
			if isRecordValueSet && isRecordValuesSet {
				return nil, fmt.Errorf("`record.value` and `record.values` cannot be set together")
			}

			if isRecordValueSet {
				recordValue := record["value"].(string)
				segmentLen := 255

				for len(recordValue) > segmentLen {
					values = append(values, recordValue[:segmentLen])
					recordValue = recordValue[segmentLen:]
				}

				values = append(values, recordValue)
			}
		}

		if isRecordValuesSet {
			recordValues := record["values"].([]interface{})

			for _, recordVal := range recordValues {
				values = append(values, recordVal.(string))
			}
		}

		records[i] = recordsets.TxtRecord{
			Value: &values,
		}
	}

	return &records, nil
}
