package azurerm

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/privatedns/mgmt/2018-09-01/privatedns"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmPrivateDnsZoneVirtualNetworkLink() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmPrivateDnsZoneVirtualNetworkLinkCreateUpdate,
		Read:   resourceArmPrivateDnsZoneVirtualNetworkLinkRead,
		Update: resourceArmPrivateDnsZoneVirtualNetworkLinkCreateUpdate,
		Delete: resourceArmPrivateDnsZoneVirtualNetworkLinkDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"private_dns_zone_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"virtual_network_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"registration_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"tags": tagsSchema(),
		},
	}
}

func resourceArmPrivateDnsZoneVirtualNetworkLinkCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).privateDns.VirtualNetworkLinksClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	dnsZoneName := d.Get("private_dns_zone_name").(string)
	vNetID := d.Get("virtual_network_id").(string)
	registrationEnabled := d.Get("registration_enabled").(bool)
	resGroup := d.Get("resource_group_name").(string)

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, dnsZoneName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("error checking for presence of existing Private DNS Zone Virtual network link %q (Resource Group %q): %s", name, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_private_dns_zone_virtual_network_link", *existing.ID)
		}
	}

	location := "global"
	tags := d.Get("tags").(map[string]interface{})

	parameters := privatedns.VirtualNetworkLink{
		Location: &location,
		Tags:     expandTags(tags),
		VirtualNetworkLinkProperties: &privatedns.VirtualNetworkLinkProperties{
			VirtualNetwork: &privatedns.SubResource{
				ID: &vNetID,
			},
			RegistrationEnabled: &registrationEnabled,
		},
	}

	etag := ""
	ifNoneMatch := "" // set to empty to allow updates to records after creation

	future, err := client.CreateOrUpdate(ctx, resGroup, dnsZoneName, name, parameters, etag, ifNoneMatch)
	if err != nil {
		return fmt.Errorf("error creating/updating Private DNS Zone Virtual network link %q (Resource Group %q): %s", name, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("error waiting for Private DNS Zone Virtual network link %q to become available: %+v", name, err)
	}

	resp, err := client.Get(ctx, resGroup, dnsZoneName, name)
	if err != nil {
		return fmt.Errorf("error retrieving Private DNS Zone Virtual network link %q (Resource Group %q): %s", name, resGroup, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("cannot read Private DNS Zone Virtual network link %q (Resource Group %q) ID", name, resGroup)
	}

	d.SetId(*resp.ID)

	return resourceArmPrivateDnsZoneVirtualNetworkLinkRead(d, meta)
}

func resourceArmPrivateDnsZoneVirtualNetworkLinkRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).privateDns.VirtualNetworkLinksClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	dnsZoneName := id.Path["privateDnsZones"]
	name := id.Path["virtualNetworkLinks"]

	resp, err := client.Get(ctx, resGroup, dnsZoneName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error reading Private DNS Zone Virtual network link %q (Resource Group %q): %+v", name, resGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)
	d.Set("private_dns_zone_name", dnsZoneName)

	if props := resp.VirtualNetworkLinkProperties; props != nil {
		d.Set("registration_enabled", props.RegistrationEnabled)

		if network := props.VirtualNetwork; network != nil {
			d.Set("virtual_network_id", network.ID)
		}
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmPrivateDnsZoneVirtualNetworkLinkDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).privateDns.VirtualNetworkLinksClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	dnsZoneName := id.Path["privateDnsZones"]
	name := id.Path["virtualNetworkLinks"]

	etag := ""
	if future, err := client.Delete(ctx, resGroup, dnsZoneName, name, etag); err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("error deleting Virtual Network Link %q (Private DNS Zone %q / Resource Group %q): %+v", name, dnsZoneName, resGroup, err)
	}

	// whilst the Delete above returns a Future, the Azure API's broken such that even though it's marked as "gone"
	// it's still kicking around - so we have to poll until this is actually gone
	log.Printf("[DEBUG] Waiting for Virtual Network Link %q (Private DNS Zone %q / Resource Group %q) to be deleted", name, dnsZoneName, resGroup)
	stateConf := &resource.StateChangeConf{
		Pending: []string{"Available"},
		Target:  []string{"NotFound"},
		Refresh: func() (interface{}, string, error) {
			log.Printf("[DEBUG] Checking to see if Virtual Network Link %q (Private DNS Zone %q / Resource Group %q) is available", name, dnsZoneName, resGroup)
			resp, err := client.Get(ctx, resGroup, dnsZoneName, name)
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					log.Printf("[DEBUG] Virtual Network Link %q (Private DNS Zone %q / Resource Group %q) was not found", name, dnsZoneName, resGroup)
					return "NotFound", "NotFound", nil
				}

				return "", "error", err
			}

			log.Printf("[DEBUG] Virtual Network Link %q (Private DNS Zone %q / Resource Group %q) still exists", name, dnsZoneName, resGroup)
			return "Available", "Available", nil
		},
		Timeout:                   30 * time.Minute,
		Delay:                     30 * time.Second,
		PollInterval:              10 * time.Second,
		ContinuousTargetOccurence: 10,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("error waiting for deletion of Virtual Network Link %q (Private DNS Zone %q / Resource Group %q): %+v", name, dnsZoneName, resGroup, err)
	}

	return nil
}
