package privatedns

import (
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/privatedns/mgmt/2018-09-01/privatedns"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/privatedns/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/privatedns/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourcePrivateDnsZone() *schema.Resource {
	return &schema.Resource{
		Create: resourcePrivateDnsZoneCreateUpdate,
		Read:   resourcePrivateDnsZoneRead,
		Update: resourcePrivateDnsZoneCreateUpdate,
		Delete: resourcePrivateDnsZoneDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.PrivateDnsZoneID(id)
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
				ForceNew: true,
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

func resourcePrivateDnsZoneCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).PrivateDns.PrivateZonesClient
	recordSetsClient := meta.(*clients.Client).PrivateDns.RecordSetsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceId := parse.NewPrivateDnsZoneID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceId.ResourceGroup, resourceId.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Private DNS Zone %q (Resource Group %q): %+v", resourceId.Name, resourceId.ResourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_private_dns_zone", resourceId.ID())
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
	future, err := client.CreateOrUpdate(ctx, resourceId.ResourceGroup, resourceId.Name, parameters, etag, ifNoneMatch)
	if err != nil {
		return fmt.Errorf("creating/updating Private DNS Zone %q (Resource Group %q): %s", resourceId.Name, resourceId.ResourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for create/update of Private DNS Zone %q (Resource Group %q): %+v", resourceId.Name, resourceId.ResourceGroup, err)
	}

	if v, ok := d.GetOk("soa_record"); ok {
		soaRecordRaw := v.([]interface{})[0].(map[string]interface{})
		soaRecord := expandPrivateDNSZoneSOARecord(soaRecordRaw)
		rsParameters := privatedns.RecordSet{
			RecordSetProperties: &privatedns.RecordSetProperties{
				TTL:       utils.Int64(int64(soaRecordRaw["ttl"].(int))),
				Metadata:  tags.Expand(soaRecordRaw["tags"].(map[string]interface{})),
				SoaRecord: soaRecord,
			},
		}

		val := fmt.Sprintf("%s%s", resourceId.Name, strings.TrimSuffix(*soaRecord.Email, "."))
		if len(val) > 253 {
			return fmt.Errorf("the value %q for `email` which is concatenated with Private DNS Zone `name` cannot exceed 253 characters excluding a trailing period", val)
		}

		if _, err := recordSetsClient.CreateOrUpdate(ctx, resourceId.ResourceGroup, resourceId.Name, privatedns.SOA, "@", rsParameters, etag, ifNoneMatch); err != nil {
			return fmt.Errorf("creating/updating Private DNS SOA Record @ (Zone %q / Resource Group %q): %s", resourceId.Name, resourceId.ResourceGroup, err)
		}
	}

	d.SetId(resourceId.ID())
	return resourcePrivateDnsZoneRead(d, meta)
}

func resourcePrivateDnsZoneRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).PrivateDns.PrivateZonesClient
	recordSetsClient := meta.(*clients.Client).PrivateDns.RecordSetsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PrivateDnsZoneID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Private DNS Zone %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	recordSetResp, err := recordSetsClient.Get(ctx, id.ResourceGroup, id.Name, privatedns.SOA, "@")
	if err != nil {
		return fmt.Errorf("reading DNS SOA record @: %v", err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	if props := resp.PrivateZoneProperties; props != nil {
		d.Set("number_of_record_sets", props.NumberOfRecordSets)
		d.Set("max_number_of_record_sets", props.MaxNumberOfRecordSets)
		d.Set("max_number_of_virtual_network_links", props.MaxNumberOfVirtualNetworkLinks)
		d.Set("max_number_of_virtual_network_links_with_registration", props.MaxNumberOfVirtualNetworkLinksWithRegistration)
	}

	if err := d.Set("soa_record", flattenPrivateDNSZoneSOARecord(&recordSetResp)); err != nil {
		return fmt.Errorf("setting `soa_record`: %+v", err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourcePrivateDnsZoneDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).PrivateDns.PrivateZonesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PrivateDnsZoneID(d.Id())
	if err != nil {
		return err
	}

	etag := ""
	future, err := client.Delete(ctx, id.ResourceGroup, id.Name, etag)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("error deleting Private DNS Zone %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("error deleting Private DNS Zone %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return nil
}

func expandPrivateDNSZoneSOARecord(input map[string]interface{}) *privatedns.SoaRecord {
	return &privatedns.SoaRecord{
		Email:       utils.String(input["email"].(string)),
		ExpireTime:  utils.Int64(int64(input["expire_time"].(int))),
		MinimumTTL:  utils.Int64(int64(input["minimum_ttl"].(int))),
		RefreshTime: utils.Int64(int64(input["refresh_time"].(int))),
		RetryTime:   utils.Int64(int64(input["retry_time"].(int))),
	}
}

func flattenPrivateDNSZoneSOARecord(input *privatedns.RecordSet) []interface{} {
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
