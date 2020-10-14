package dns

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/dns/mgmt/2018-05-01/dns"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/dns/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmDnsSoaRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDnsSoaRecordCreateUpdate,
		Read:   resourceArmDnsSoaRecordRead,
		Update: resourceArmDnsSoaRecordCreateUpdate,
		Delete: resourceArmDnsSoaRecordDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.DnsSoaRecordID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"zone_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"record": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"email": {
							Type:     schema.TypeString,
							Required: true,
						},

						"expire_time": {
							Type:     schema.TypeInt,
							Required: true,
						},

						"host_name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"minimum_ttl": {
							Type:     schema.TypeInt,
							Required: true,
						},

						"refresh_time": {
							Type:     schema.TypeInt,
							Required: true,
						},

						"retry_time": {
							Type:     schema.TypeInt,
							Required: true,
						},

						"serial_number": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},

			"ttl": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"fqdn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmDnsSoaRecordCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Dns.RecordSetsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	zoneName := d.Get("zone_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, zoneName, name, dns.SRV)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing DNS SOA Record %q (Zone %q / Resource Group %q): %s", name, zoneName, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_dns_soa_record", *existing.ID)
		}
	}

	ttl := int64(d.Get("ttl").(int))
	t := d.Get("tags").(map[string]interface{})
	record := d.Get("record").([]interface{})

	parameters := dns.RecordSet{
		Name: &name,
		RecordSetProperties: &dns.RecordSetProperties{
			Metadata:  tags.Expand(t),
			TTL:       &ttl,
			SoaRecord: expandAzureRmDnsSoaRecord(record),
		},
	}

	eTag := ""
	ifNoneMatch := "" // set to empty to allow updates to records after creation
	if _, err := client.CreateOrUpdate(ctx, resGroup, zoneName, name, dns.SOA, parameters, eTag, ifNoneMatch); err != nil {
		return fmt.Errorf("Error creating/updating DNS SOA Record %q (Zone %q / Resource Group %q): %s", name, zoneName, resGroup, err)
	}

	resp, err := client.Get(ctx, resGroup, zoneName, name, dns.SOA)
	if err != nil {
		return fmt.Errorf("Error retrieving DNS SOA Record %q (Zone %q / Resource Group %q): %s", name, zoneName, resGroup, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read DNS SOA Record %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*resp.ID)

	return resourceArmDnsSoaRecordRead(d, meta)
}

func resourceArmDnsSoaRecordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Dns.RecordSetsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DnsSoaRecordID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ZoneName, id.Name, dns.SOA)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading DNS SOA record %s: %v", id.Name, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("zone_name", id.ZoneName)
	d.Set("ttl", resp.TTL)
	d.Set("fqdn", resp.Fqdn)

	if err := d.Set("record", flattenAzureRmDnsSoaRecord(resp.SoaRecord)); err != nil {
		return err
	}
	return tags.FlattenAndSet(d, resp.Metadata)
}

func resourceArmDnsSoaRecordDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Dns.RecordSetsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DnsSoaRecordID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, id.ResourceGroup, id.ZoneName, id.Name, dns.SOA, "")
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error deleting DNS SOA Record %s: %+v", id.Name, err)
	}

	return nil
}

func expandAzureRmDnsSoaRecord(input []interface{}) *dns.SoaRecord {
	record := input[0].(map[string]interface{})

	email := record["email"].(string)
	expireTime := int64(record["expire_time"].(int))
	hostName := record["host_name"].(string)
	minimumTTL := int64(record["minimum_ttl"].(int))
	refreshTime := int64(record["refresh_time"].(int))
	retryTime := int64(record["retry_time"].(int))
	serialNumber := int64(record["serial_number"].(int))

	return &dns.SoaRecord{
		Email:        &email,
		ExpireTime:   &expireTime,
		Host:         &hostName,
		MinimumTTL:   &minimumTTL,
		RefreshTime:  &refreshTime,
		RetryTime:    &retryTime,
		SerialNumber: &serialNumber,
	}
}

func flattenAzureRmDnsSoaRecord(input *dns.SoaRecord) []interface{} {

	return []interface{}{
		map[string]interface{}{
			"email":         *input.Email,
			"expire_time":   *input.ExpireTime,
			"host_name":     *input.Host,
			"minimum_ttl":   *input.MinimumTTL,
			"refresh_time":  *input.RefreshTime,
			"retry_time":    *input.RetryTime,
			"serial_number": *input.SerialNumber,
		},
	}
}
