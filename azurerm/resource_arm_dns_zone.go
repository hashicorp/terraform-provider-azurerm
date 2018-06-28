package azurerm

import (
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/preview/dns/mgmt/2018-03-01-preview/dns"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmDnsZone() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDnsZoneCreate,
		Read:   resourceArmDnsZoneRead,
		Update: resourceArmDnsZoneCreate,
		Delete: resourceArmDnsZoneDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
				}, true),
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
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

func convertZoneType(zoneType string) dns.ZoneType {
	for _, zType := range dns.PossibleZoneTypeValues() {
		if strings.EqualFold(string(zType), zoneType) {
			return zType
		}
	}
	return dns.Public
}

func resourceArmDnsZoneCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).zonesClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	location := "global"
	zoneType := d.Get("zone_type").(string)

	tags := d.Get("tags").(map[string]interface{})
	registrationVirtualNetworks := d.Get("registration_virtual_network_ids").([]interface{})
	registrationVNetSubResources := make([]dns.SubResource, 0, len(registrationVirtualNetworks))
	for _, rvn := range registrationVirtualNetworks {
		id := rvn.(string)
		registrationVNetSubResources = append(registrationVNetSubResources, dns.SubResource{
			ID: &id,
		})
	}
	resolutionVirtualNetworks := d.Get("resolution_virtual_network_ids").([]interface{})
	resolutionVNetSubResources := make([]dns.SubResource, 0, len(resolutionVirtualNetworks))
	for _, rvn := range resolutionVirtualNetworks {
		id := rvn.(string)
		resolutionVNetSubResources = append(resolutionVNetSubResources, dns.SubResource{
			ID: &id,
		})
	}

	parameters := dns.Zone{
		Location: &location,
		Tags:     expandTags(tags),
		ZoneProperties: &dns.ZoneProperties{
			ZoneType:                    convertZoneType(zoneType),
			RegistrationVirtualNetworks: &registrationVNetSubResources,
			ResolutionVirtualNetworks:   &resolutionVNetSubResources,
		},
	}

	etag := ""
	ifNoneMatch := "" // set to empty to allow updates to records after creation
	resp, err := client.CreateOrUpdate(ctx, resGroup, name, parameters, etag, ifNoneMatch)
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

	if resp.RegistrationVirtualNetworks != nil {
		registrationVirtualNetworks := make([]string, 0, len(*resp.RegistrationVirtualNetworks))
		for _, rvn := range *resp.RegistrationVirtualNetworks {
			registrationVirtualNetworks = append(registrationVirtualNetworks, *rvn.ID)
		}
		if err := d.Set("registration_virtual_network_ids", registrationVirtualNetworks); err != nil {
			return err
		}
	}

	if resp.ResolutionVirtualNetworks != nil {
		resolutionVirtualNetworks := make([]string, 0, len(*resp.ResolutionVirtualNetworks))
		for _, rvn := range *resp.ResolutionVirtualNetworks {
			resolutionVirtualNetworks = append(resolutionVirtualNetworks, *rvn.ID)
		}
		if err := d.Set("resolution_virtual_network_ids", resolutionVirtualNetworks); err != nil {
			return err
		}
	}

	if resp.NameServers != nil {
		nameServers := make([]string, 0, len(*resp.NameServers))
		for _, ns := range *resp.NameServers {
			nameServers = append(nameServers, ns)
		}
		if err := d.Set("name_servers", nameServers); err != nil {
			return err
		}
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

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error deleting DNS zone %s (resource group %s): %+v", name, resGroup, err)
	}

	return nil
}
