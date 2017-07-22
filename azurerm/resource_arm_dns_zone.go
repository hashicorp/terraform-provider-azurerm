package azurerm

import (
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/arm/dns"
	"github.com/hashicorp/terraform/helper/schema"
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
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": &schema.Schema{
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: resourceAzurermResourceGroupNameDiffSuppress,
			},

			"number_of_record_sets": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"max_number_of_record_sets": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_servers": &schema.Schema{
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmDnsZoneCreate(d *schema.ResourceData, meta interface{}) error {
	zonesClient := meta.(*ArmClient).zonesClient

	zoneName := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	location := "global"

	tags := d.Get("tags").(map[string]interface{})
	metadata := expandTags(tags)

	parameters := dns.Zone{
		Name:     &zoneName,
		Location: &location,
		Tags:     metadata,
	}

	//last parameter is set to empty to allow updates to records after creation
	// (per SDK, set it to '*' to prevent updates, all other values are ignored)
	resp, err := zonesClient.CreateOrUpdate(resGroup, zoneName, parameters, "", "")
	if err != nil {
		return err
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read DNS zone %s (resource group %s) ID", zoneName, resGroup)
	}

	d.SetId(*resp.ID)

	return resourceArmDnsZoneRead(d, meta)
}

func resourceArmDnsZoneRead(d *schema.ResourceData, meta interface{}) error {
	zonesClient := meta.(*ArmClient).zonesClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	zoneName := id.Path["dnszones"]

	resp, err := zonesClient.Get(resGroup, zoneName)
	if err != nil {
		return fmt.Errorf("Error reading DNS zone %s: %v", zoneName, err)
	}
	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}

	d.Set("name", zoneName)
	d.Set("resource_group_name", resGroup)
	d.Set("number_of_record_sets", resp.NumberOfRecordSets)
	d.Set("max_number_of_record_sets", resp.MaxNumberOfRecordSets)

	nameServers := make([]string, 0, len(*resp.NameServers))
	for _, ns := range *resp.NameServers {
		nameServers = append(nameServers, ns)
	}
	if err := d.Set("name_servers", nameServers); err != nil {
		return err
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmDnsZoneDelete(d *schema.ResourceData, meta interface{}) error {
	zonesClient := meta.(*ArmClient).zonesClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	zoneName := id.Path["dnszones"]

	resultChan, errorChan := zonesClient.Delete(resGroup, zoneName, "", make(chan struct{}))
	result := <-resultChan
	error := <-errorChan
	if result.Response.StatusCode != http.StatusOK {
		return fmt.Errorf("Error deleting DNS zone %s: %+v", zoneName, error)
	}

	return nil
}
