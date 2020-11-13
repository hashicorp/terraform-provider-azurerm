package privatedns

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/privatedns/mgmt/2018-09-01/privatedns"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/privatedns/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmPrivateDnsZone() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmPrivateDnsZoneCreateUpdate,
		Read:   resourceArmPrivateDnsZoneRead,
		Update: resourceArmPrivateDnsZoneCreateUpdate,
		Delete: resourceArmPrivateDnsZoneDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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
			},

			"number_of_record_sets": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"max_number_of_record_sets": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"max_number_of_virtual_network_links": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"max_number_of_virtual_network_links_with_registration": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"soa_record": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"email": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.PrivateDnsZoneSOARecordEmail,
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
							Default:      10,
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

						// This field should be able to be updated since DNS Record Sets API allows to update it.
						// So the issue is submitted on https://github.com/Azure/azure-rest-api-specs/issues/11674
						// Once the issue is fixed, the field will be updated to `Required` property.
						"host_name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						// This field should be able to be updated since DNS Record Sets API allows to update it.
						// So the issue is submitted on https://github.com/Azure/azure-rest-api-specs/issues/11674
						// Once the issue is fixed, the field will be updated to `Optional` property with `Default` attribute.
						"serial_number": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmPrivateDnsZoneCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).PrivateDns.PrivateZonesClient
	recordSetsClient := meta.(*clients.Client).PrivateDns.RecordSetsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("error checking for presence of existing Private DNS Zone %q (Resource Group %q): %s", name, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_private_dns_zone", *existing.ID)
		}
	}

	location := "global"
	t := d.Get("tags").(map[string]interface{})

	parameters := privatedns.PrivateZone{
		Location: &location,
		Tags:     tags.Expand(t),
	}

	etag := ""
	ifNoneMatch := "" // set to empty to allow updates to records after creation
	future, err := client.CreateOrUpdate(ctx, resGroup, name, parameters, etag, ifNoneMatch)
	if err != nil {
		return fmt.Errorf("error creating/updating Private DNS Zone %q (Resource Group %q): %s", name, resGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("error waiting for Private DNS Zone %q to become available: %+v", name, err)
	}

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("error retrieving Private DNS Zone %q (Resource Group %q): %s", name, resGroup, err)
	}

	if v, ok := d.GetOk("soa_record"); ok {
		soaRecord := v.([]interface{})[0].(map[string]interface{})
		rsParameters := privatedns.RecordSet{
			RecordSetProperties: &privatedns.RecordSetProperties{
				TTL:       utils.Int64(int64(soaRecord["ttl"].(int))),
				Metadata:  tags.Expand(soaRecord["tags"].(map[string]interface{})),
				SoaRecord: expandArmPrivateDNSZoneSOARecord(soaRecord),
			},
		}

		if _, err := recordSetsClient.CreateOrUpdate(ctx, resGroup, name, privatedns.SOA, "@", rsParameters, etag, ifNoneMatch); err != nil {
			return fmt.Errorf("creating/updating Private DNS SOA Record @ (Zone %q / Resource Group %q): %s", name, resGroup, err)
		}
	}

	if resp.ID == nil {
		return fmt.Errorf("cannot read Private DNS Zone %q (Resource Group %q) ID", name, resGroup)
	}

	d.SetId(*resp.ID)

	return resourceArmPrivateDnsZoneRead(d, meta)
}

func resourceArmPrivateDnsZoneRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).PrivateDns.PrivateZonesClient
	recordSetsClient := meta.(*clients.Client).PrivateDns.RecordSetsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["privateDnsZones"]

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error reading Private DNS Zone %q (Resource Group %q): %+v", name, resGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)

	if props := resp.PrivateZoneProperties; props != nil {
		d.Set("number_of_record_sets", props.NumberOfRecordSets)
		d.Set("max_number_of_record_sets", props.MaxNumberOfRecordSets)
		d.Set("max_number_of_virtual_network_links", props.MaxNumberOfVirtualNetworkLinks)
		d.Set("max_number_of_virtual_network_links_with_registration", props.MaxNumberOfVirtualNetworkLinksWithRegistration)
	}

	rsResp, err := recordSetsClient.Get(ctx, id.ResourceGroup, name, privatedns.SOA, "@")
	if err != nil {
		return fmt.Errorf("reading DNS SOA record @: %v", err)
	}

	if err := d.Set("soa_record", flattenArmPrivateDNSZoneSOARecord(&rsResp)); err != nil {
		return fmt.Errorf("setting `soa_record`: %+v", err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmPrivateDnsZoneDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).PrivateDns.PrivateZonesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["privateDnsZones"]

	etag := ""
	future, err := client.Delete(ctx, resGroup, name, etag)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("error deleting Private DNS Zone %s (resource group %s): %+v", name, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("error deleting Private DNS Zone %s (resource group %s): %+v", name, resGroup, err)
	}

	return nil
}

func expandArmPrivateDNSZoneSOARecord(input map[string]interface{}) *privatedns.SoaRecord {
	return &privatedns.SoaRecord{
		Email:       utils.String(input["email"].(string)),
		ExpireTime:  utils.Int64(int64(input["expire_time"].(int))),
		MinimumTTL:  utils.Int64(int64(input["minimum_ttl"].(int))),
		RefreshTime: utils.Int64(int64(input["refresh_time"].(int))),
		RetryTime:   utils.Int64(int64(input["retry_time"].(int))),
	}
}

func flattenArmPrivateDNSZoneSOARecord(input *privatedns.RecordSet) []interface{} {
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
