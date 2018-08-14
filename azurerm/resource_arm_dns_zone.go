package azurerm

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/dns/mgmt/2018-03-01-preview/dns"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmDnsZone() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDnsZoneCreateUpdate,
		Read:   resourceArmDnsZoneRead,
		Update: resourceArmDnsZoneCreateUpdate,
		Delete: resourceArmDnsZoneDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(time.Minute * 30),
			Update: schema.DefaultTimeout(time.Minute * 30),
			Delete: schema.DefaultTimeout(time.Minute * 30),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": resourceGroupNameDiffSuppressSchema(),

			"number_of_record_sets": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"max_number_of_record_sets": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"name_servers": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"zone_type": {
				Type:     schema.TypeString,
				Default:  string(dns.Public),
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(dns.Private),
					string(dns.Public),
				}, false),
			},

			"registration_virtual_network_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"resolution_virtual_network_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmDnsZoneCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).zonesClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		// first check if there's one in this subscription requiring import
		resp, err := client.Get(ctx, resGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Error checking for the existence of DNS Zone %q (Resource Group %q): %+v", name, resGroup, err)
			}
		}

		if resp.ID != nil {
			return tf.ImportAsExistsError("azurerm_dns_zone", *resp.ID)
		}
	}

	location := "global"
	zoneType := d.Get("zone_type").(string)
	tags := d.Get("tags").(map[string]interface{})

	registrationVirtualNetworkIds := expandDnsZoneRegistrationVirtualNetworkIds(d)
	resolutionVirtualNetworkIds := expandDnsZoneResolutionVirtualNetworkIds(d)

	parameters := dns.Zone{
		Location: &location,
		Tags:     expandTags(tags),
		ZoneProperties: &dns.ZoneProperties{
			ZoneType:                    dns.ZoneType(zoneType),
			RegistrationVirtualNetworks: registrationVirtualNetworkIds,
			ResolutionVirtualNetworks:   resolutionVirtualNetworkIds,
		},
	}

	etag := ""
	ifNoneMatch := "" // set to empty to allow updates to records after creation
	waitCtx, cancel := context.WithTimeout(ctx, d.Timeout(tf.TimeoutForCreateUpdate(d)))
	defer cancel()
	resp, err := client.CreateOrUpdate(waitCtx, resGroup, name, parameters, etag, ifNoneMatch)
	if err != nil {
		return err
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read DNS zone %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*resp.ID)

	return resourceArmDnsZoneRead(d, meta)
}

func resourceArmDnsZoneRead(d *schema.ResourceData, meta interface{}) error {
	zonesClient := meta.(*ArmClient).zonesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["dnszones"]

	resp, err := zonesClient.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading DNS zone %s (resource group %s): %+v", name, resGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)
	d.Set("number_of_record_sets", resp.NumberOfRecordSets)
	d.Set("max_number_of_record_sets", resp.MaxNumberOfRecordSets)
	d.Set("zone_type", resp.ZoneType)

	registrationVirtualNetworks := flattenDnsZoneRegistrationVirtualNetworkIDs(resp.RegistrationVirtualNetworks)
	if err := d.Set("registration_virtual_network_ids", registrationVirtualNetworks); err != nil {
		return err
	}

	resolutionVirtualNetworks := flattenDnsZoneResolutionVirtualNetworkIDs(resp.ResolutionVirtualNetworks)
	if err := d.Set("resolution_virtual_network_ids", resolutionVirtualNetworks); err != nil {
		return err
	}

	nameServers := flattenDnsZoneNameservers(resp.NameServers)
	if err := d.Set("name_servers", nameServers); err != nil {
		return err
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmDnsZoneDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).zonesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["dnszones"]

	etag := ""
	future, err := client.Delete(ctx, resGroup, name, etag)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error deleting DNS zone %s (resource group %s): %+v", name, resGroup, err)
	}

	waitCtx, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutDelete))
	defer cancel()
	err = future.WaitForCompletionRef(waitCtx, client.Client)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error deleting DNS zone %s (resource group %s): %+v", name, resGroup, err)
	}

	return nil
}

func expandDnsZoneResolutionVirtualNetworkIds(d *schema.ResourceData) *[]dns.SubResource {
	resolutionVirtualNetworks := d.Get("resolution_virtual_network_ids").([]interface{})

	resolutionVNetSubResources := make([]dns.SubResource, 0, len(resolutionVirtualNetworks))
	for _, rvn := range resolutionVirtualNetworks {
		id := rvn.(string)
		resolutionVNetSubResources = append(resolutionVNetSubResources, dns.SubResource{
			ID: &id,
		})
	}

	return &resolutionVNetSubResources
}

func flattenDnsZoneRegistrationVirtualNetworkIDs(input *[]dns.SubResource) []string {
	registrationVirtualNetworks := make([]string, 0)

	if input != nil {
		for _, rvn := range *input {
			registrationVirtualNetworks = append(registrationVirtualNetworks, *rvn.ID)
		}
	}

	return registrationVirtualNetworks
}

func expandDnsZoneRegistrationVirtualNetworkIds(d *schema.ResourceData) *[]dns.SubResource {
	registrationVirtualNetworks := d.Get("registration_virtual_network_ids").([]interface{})

	registrationVNetSubResources := make([]dns.SubResource, 0)
	for _, rvn := range registrationVirtualNetworks {
		id := rvn.(string)
		registrationVNetSubResources = append(registrationVNetSubResources, dns.SubResource{
			ID: &id,
		})
	}

	return &registrationVNetSubResources
}

func flattenDnsZoneResolutionVirtualNetworkIDs(input *[]dns.SubResource) []string {
	resolutionVirtualNetworks := make([]string, 0)

	if input != nil {
		for _, rvn := range *input {
			resolutionVirtualNetworks = append(resolutionVirtualNetworks, *rvn.ID)
		}
	}

	return resolutionVirtualNetworks
}

func flattenDnsZoneNameservers(input *[]string) []string {
	nameServers := make([]string, 0)

	if input != nil {
		for _, ns := range *input {
			nameServers = append(nameServers, ns)
		}
	}

	return nameServers
}
