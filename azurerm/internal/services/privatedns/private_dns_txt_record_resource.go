package privatedns

import (
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/privatedns/mgmt/2018-09-01/privatedns"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/privatedns/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmPrivateDnsTxtRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmPrivateDnsTxtRecordCreateUpdate,
		Read:   resourceArmPrivateDnsTxtRecordRead,
		Update: resourceArmPrivateDnsTxtRecordCreateUpdate,
		Delete: resourceArmPrivateDnsTxtRecordDelete,
		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.TxtRecordID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				// lower-cased due to the broken API https://github.com/Azure/azure-rest-api-specs/issues/6641
				ValidateFunc: validate.LowerCasedString,
			},

			// TODO: make this case sensitive once the API's fixed https://github.com/Azure/azure-rest-api-specs/issues/6641
			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"zone_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"record": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 1024),
						},
					},
				},
			},

			"ttl": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 2147483647),
			},

			"fqdn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmPrivateDnsTxtRecordCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).PrivateDns.RecordSetsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceId := parse.NewTxtRecordID(subscriptionId, d.Get("resource_group_name").(string), d.Get("zone_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceId.ResourceGroup, resourceId.PrivateDnsZoneName, privatedns.TXT, resourceId.TXTName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of %s: %+v", resourceId, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_private_dns_txt_record", resourceId.ID(""))
		}
	}

	parameters := privatedns.RecordSet{
		Name: utils.String(resourceId.TXTName),
		RecordSetProperties: &privatedns.RecordSetProperties{
			Metadata:   tags.Expand(d.Get("tags").(map[string]interface{})),
			TTL:        utils.Int64(int64(d.Get("ttl").(int))),
			TxtRecords: expandAzureRmPrivateDnsTxtRecords(d),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, resourceId.ResourceGroup, resourceId.PrivateDnsZoneName, privatedns.TXT, resourceId.TXTName, parameters, "", ""); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", resourceId, err)
	}

	d.SetId(resourceId.ID(""))
	return resourceArmPrivateDnsTxtRecordRead(d, meta)
}

func resourceArmPrivateDnsTxtRecordRead(d *schema.ResourceData, meta interface{}) error {
	dnsClient := meta.(*clients.Client).PrivateDns.RecordSetsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.TxtRecordID(d.Id())
	if err != nil {
		return err
	}

	resp, err := dnsClient.Get(ctx, id.ResourceGroup, id.PrivateDnsZoneName, privatedns.TXT, id.TXTName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.TXTName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("zone_name", id.PrivateDnsZoneName)

	d.Set("ttl", resp.TTL)
	d.Set("fqdn", resp.Fqdn)

	if err := d.Set("record", flattenAzureRmPrivateDnsTxtRecords(resp.TxtRecords)); err != nil {
		return fmt.Errorf("setting `record`: %+v", err)
	}

	return tags.FlattenAndSet(d, resp.Metadata)
}

func resourceArmPrivateDnsTxtRecordDelete(d *schema.ResourceData, meta interface{}) error {
	dnsClient := meta.(*clients.Client).PrivateDns.RecordSetsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.TxtRecordID(d.Id())
	if err != nil {
		return err
	}

	if _, err = dnsClient.Delete(ctx, id.ResourceGroup, id.PrivateDnsZoneName, privatedns.TXT, id.PrivateDnsZoneName, ""); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func flattenAzureRmPrivateDnsTxtRecords(records *[]privatedns.TxtRecord) []map[string]interface{} {
	results := make([]map[string]interface{}, 0)

	if records != nil {
		for _, record := range *records {
			txtRecord := make(map[string]interface{})

			if v := record.Value; v != nil {
				value := strings.Join(*v, "")
				txtRecord["value"] = value
			}

			results = append(results, txtRecord)
		}
	}

	return results
}

func expandAzureRmPrivateDnsTxtRecords(d *schema.ResourceData) *[]privatedns.TxtRecord {
	recordStrings := d.Get("record").(*schema.Set).List()
	records := make([]privatedns.TxtRecord, len(recordStrings))

	segmentLen := 254
	for i, v := range recordStrings {
		if v == nil {
			continue
		}

		record := v.(map[string]interface{})
		v := record["value"].(string)

		var value []string
		for len(v) > segmentLen {
			value = append(value, v[:segmentLen])
			v = v[segmentLen:]
		}
		value = append(value, v)

		txtRecord := privatedns.TxtRecord{
			Value: &value,
		}

		records[i] = txtRecord
	}

	return &records
}
