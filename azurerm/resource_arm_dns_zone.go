package azurerm

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/dns/mgmt/2018-03-01-preview/dns"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
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

			"zone_type": {
				Type:       schema.TypeString,
				Default:    string(dns.Public),
				Optional:   true,
				Deprecated: "Use the `azurerm_private_dns_zone` resource instead.",
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

			"tags": tags.Schema(),
		},
	}
}

func resourceArmDnsZoneCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).dns.ZonesClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
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
	zoneType := d.Get("zone_type").(string)
	t := d.Get("tags").(map[string]interface{})

	registrationVirtualNetworkIds := expandDnsZoneRegistrationVirtualNetworkIds(d)
	resolutionVirtualNetworkIds := expandDnsZoneResolutionVirtualNetworkIds(d)

	parameters := dns.Zone{
		Location: &location,
		Tags:     tags.Expand(t),
		ZoneProperties: &dns.ZoneProperties{
			ZoneType:                    dns.ZoneType(zoneType),
			RegistrationVirtualNetworks: registrationVirtualNetworkIds,
			ResolutionVirtualNetworks:   resolutionVirtualNetworkIds,
		},
	}

	etag := ""
	ifNoneMatch := "" // set to empty to allow updates to records after creation
	_, err := client.CreateOrUpdate(ctx, resGroup, name, parameters, etag, ifNoneMatch)
	if err != nil {
		return fmt.Errorf("Error creating/updating DNS Zone %q (Resource Group %q): %s", name, resGroup, err)
	}

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving DNS Zone %q (Resource Group %q): %s", name, resGroup, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read DNS Zone %q (Resource Group %q) ID", name, resGroup)
	}

	d.SetId(*resp.ID)

	return resourceArmDnsZoneRead(d, meta)
}

func resourceArmDnsZoneRead(d *schema.ResourceData, meta interface{}) error {
	zonesClient := meta.(*ArmClient).dns.ZonesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
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
		return fmt.Errorf("Error reading DNS Zone %q (Resource Group %q): %+v", name, resGroup, err)
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

	nameServers := make([]string, 0)
	if s := resp.NameServers; s != nil {
		nameServers = *s
	}
	if err := d.Set("name_servers", nameServers); err != nil {
		return err
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmDnsZoneDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).dns.ZonesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
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

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
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
