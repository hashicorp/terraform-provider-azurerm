package dns

import (
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/dns/mgmt/2018-05-01/dns"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/dns/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/dns/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceDnsZone() *schema.Resource {
	return &schema.Resource{
		Create: resourceDnsZoneCreateUpdate,
		Read:   resourceDnsZoneRead,
		Update: resourceDnsZoneCreateUpdate,
		Delete: resourceDnsZoneDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.DnsZoneID(id)
			return err
		}),
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"number_of_record_sets": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"max_number_of_record_sets": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"name_servers": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"soa_record": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"email": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.DnsZoneSOARecordEmail,
						},

						"host_name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"expire_time": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      2419200,
							ValidateFunc: validation.IntAtLeast(0),
						},

						"minimum_ttl": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      300,
							ValidateFunc: validation.IntAtLeast(0),
						},

						"refresh_time": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      3600,
							ValidateFunc: validation.IntAtLeast(0),
						},

						"retry_time": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      300,
							ValidateFunc: validation.IntAtLeast(0),
						},

						"serial_number": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      1,
							ValidateFunc: validation.IntAtLeast(0),
						},

						"ttl": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      3600,
							ValidateFunc: validation.IntBetween(0, 2147483647),
						},

						"tags": tags.Schema(),

						"fqdn": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceDnsZoneCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Dns.ZonesClient
	recordSetsClient := meta.(*clients.Client).Dns.RecordSetsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing DNS Zone %q (Resource Group %q): %s", name, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_dns_zone", *existing.ID)
		}
	}

	location := "global"
	t := d.Get("tags").(map[string]interface{})

	parameters := dns.Zone{
		Location: &location,
		Tags:     tags.Expand(t),
	}

	etag := ""
	ifNoneMatch := "" // set to empty to allow updates to records after creation
	if _, err := client.CreateOrUpdate(ctx, resGroup, name, parameters, etag, ifNoneMatch); err != nil {
		return fmt.Errorf("Error creating/updating DNS Zone %q (Resource Group %q): %s", name, resGroup, err)
	}

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving DNS Zone %q (Resource Group %q): %s", name, resGroup, err)
	}

	if v, ok := d.GetOk("soa_record"); ok {
		soaRecord := v.([]interface{})[0].(map[string]interface{})
		rsParameters := dns.RecordSet{
			RecordSetProperties: &dns.RecordSetProperties{
				TTL:       utils.Int64(int64(soaRecord["ttl"].(int))),
				Metadata:  tags.Expand(soaRecord["tags"].(map[string]interface{})),
				SoaRecord: expandArmDNSZoneSOARecord(soaRecord),
			},
		}

		if len(name+strings.TrimSuffix(*rsParameters.RecordSetProperties.SoaRecord.Email, ".")) > 253 {
			return fmt.Errorf("`email` which is concatenated with DNS Zone `name` cannot exceed 253 characters excluding a trailing period")
		}

		if _, err := recordSetsClient.CreateOrUpdate(ctx, resGroup, name, "@", dns.SOA, rsParameters, etag, ifNoneMatch); err != nil {
			return fmt.Errorf("creating/updating DNS SOA Record @ (Zone %q / Resource Group %q): %s", name, resGroup, err)
		}
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read DNS Zone %q (Resource Group %q) ID", name, resGroup)
	}

	d.SetId(*resp.ID)

	return resourceDnsZoneRead(d, meta)
}

func resourceDnsZoneRead(d *schema.ResourceData, meta interface{}) error {
	zonesClient := meta.(*clients.Client).Dns.ZonesClient
	recordSetsClient := meta.(*clients.Client).Dns.RecordSetsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DnsZoneID(d.Id())
	if err != nil {
		return err
	}

	resp, err := zonesClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading DNS Zone %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	d.Set("number_of_record_sets", resp.NumberOfRecordSets)
	d.Set("max_number_of_record_sets", resp.MaxNumberOfRecordSets)

	nameServers := make([]string, 0)
	if s := resp.NameServers; s != nil {
		nameServers = *s
	}
	if err := d.Set("name_servers", nameServers); err != nil {
		return err
	}

	rsResp, err := recordSetsClient.Get(ctx, id.ResourceGroup, id.Name, "@", dns.SOA)
	if err != nil {
		return fmt.Errorf("reading DNS SOA record @: %v", err)
	}

	if err := d.Set("soa_record", flattenArmDNSZoneSOARecord(&rsResp)); err != nil {
		return fmt.Errorf("setting `soa_record`: %+v", err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceDnsZoneDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Dns.ZonesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DnsZoneID(d.Id())
	if err != nil {
		return err
	}

	etag := ""
	future, err := client.Delete(ctx, id.ResourceGroup, id.Name, etag)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error deleting DNS zone %s (resource group %s): %+v", id.Name, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error deleting DNS zone %s (resource group %s): %+v", id.Name, id.ResourceGroup, err)
	}

	return nil
}

func expandArmDNSZoneSOARecord(input map[string]interface{}) *dns.SoaRecord {
	return &dns.SoaRecord{
		Email:        utils.String(input["email"].(string)),
		Host:         utils.String(input["host_name"].(string)),
		ExpireTime:   utils.Int64(int64(input["expire_time"].(int))),
		MinimumTTL:   utils.Int64(int64(input["minimum_ttl"].(int))),
		RefreshTime:  utils.Int64(int64(input["refresh_time"].(int))),
		RetryTime:    utils.Int64(int64(input["retry_time"].(int))),
		SerialNumber: utils.Int64(int64(input["serial_number"].(int))),
	}
}

func flattenArmDNSZoneSOARecord(input *dns.RecordSet) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	ttl := 0
	if input.TTL != nil {
		ttl = int(*input.TTL)
	}

	metaData := make(map[string]interface{})
	if input.Metadata != nil {
		metaData = tags.Flatten(input.Metadata)
	}

	fqdn := ""
	if input.Fqdn != nil {
		fqdn = *input.Fqdn
	}

	email := ""
	hostName := ""
	expireTime := 0
	minimumTTL := 0
	refreshTime := 0
	retryTime := 0
	serialNumber := 0
	if input.SoaRecord != nil {
		if input.SoaRecord.Email != nil {
			email = *input.SoaRecord.Email
		}

		if input.SoaRecord.Host != nil {
			hostName = *input.SoaRecord.Host
		}

		if input.SoaRecord.ExpireTime != nil {
			expireTime = int(*input.SoaRecord.ExpireTime)
		}

		if input.SoaRecord.MinimumTTL != nil {
			minimumTTL = int(*input.SoaRecord.MinimumTTL)
		}

		if input.SoaRecord.RefreshTime != nil {
			refreshTime = int(*input.SoaRecord.RefreshTime)
		}

		if input.SoaRecord.RetryTime != nil {
			retryTime = int(*input.SoaRecord.RetryTime)
		}

		if input.SoaRecord.SerialNumber != nil {
			serialNumber = int(*input.SoaRecord.SerialNumber)
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
