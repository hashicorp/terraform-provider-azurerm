package dns

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/dns/mgmt/2018-05-01/dns"
	"github.com/hashicorp/go-azure-helpers/response"
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

func resourceArmDnsZone() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDnsZoneCreateUpdate,
		Read:   resourceArmDnsZoneRead,
		Update: resourceArmDnsZoneCreateUpdate,
		Delete: resourceArmDnsZoneDelete,

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

			"tags": tags.Schema(),
		},
	}
}

func resourceArmDnsZoneCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Dns.ZonesClient
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

	if resp.ID == nil {
		return fmt.Errorf("Cannot read DNS Zone %q (Resource Group %q) ID", name, resGroup)
	}

	d.SetId(*resp.ID)

	return resourceArmDnsZoneRead(d, meta)
}

func resourceArmDnsZoneRead(d *schema.ResourceData, meta interface{}) error {
	zonesClient := meta.(*clients.Client).Dns.ZonesClient
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

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmDnsZoneDelete(d *schema.ResourceData, meta interface{}) error {
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
