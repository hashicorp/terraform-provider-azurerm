package network

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	networkValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmVirtualHubSecurityPartnerProvider() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmVirtualHubSecurityPartnerProviderCreate,
		Read:   resourceArmVirtualHubSecurityPartnerProviderRead,
		Update: resourceArmVirtualHubSecurityPartnerProviderUpdate,
		Delete: resourceArmVirtualHubSecurityPartnerProviderDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.VirtualHubSecurityPartnerProviderID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"security_provider_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(network.ZScaler),
					string(network.IBoss),
					string(network.Checkpoint),
				}, false),
			},

			"virtual_hub_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: networkValidate.ValidateVirtualHubID,
			},

			"tags": tags.Schema(),
		},
	}
}
func resourceArmVirtualHubSecurityPartnerProviderCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.SecurityPartnerProviderClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	existing, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for present of existing Security Partner Provider %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_virtual_hub_security_partner_provider", *existing.ID)
	}

	parameters := network.SecurityPartnerProvider{
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		SecurityPartnerProviderPropertiesFormat: &network.SecurityPartnerProviderPropertiesFormat{
			SecurityProviderName: network.SecurityProviderName(d.Get("security_provider_name").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if v, ok := d.GetOk("virtual_hub_id"); ok {
		parameters.SecurityPartnerProviderPropertiesFormat.VirtualHub = &network.SubResource{
			ID: utils.String(v.(string)),
		}
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, parameters)
	if err != nil {
		return fmt.Errorf("creating Security Partner Provider %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on creating future for Security Partner Provider %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("retrieving Security Partner Provider %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for Security Partner Provider %q (Resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*resp.ID)

	return resourceArmVirtualHubSecurityPartnerProviderRead(d, meta)
}

func resourceArmVirtualHubSecurityPartnerProviderRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.SecurityPartnerProviderClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualHubSecurityPartnerProviderID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] security partner provider %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Security Partner Provider %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if props := resp.SecurityPartnerProviderPropertiesFormat; props != nil {
		d.Set("security_provider_name", props.SecurityProviderName)

		if props.VirtualHub != nil && props.VirtualHub.ID != nil {
			d.Set("virtual_hub_id", props.VirtualHub.ID)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmVirtualHubSecurityPartnerProviderUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.SecurityPartnerProviderClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualHubSecurityPartnerProviderID(d.Id())
	if err != nil {
		return err
	}

	parameters := network.TagsObject{}

	if d.HasChange("tags") {
		parameters.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if _, err := client.UpdateTags(ctx, id.ResourceGroup, id.Name, parameters); err != nil {
		return fmt.Errorf("updating Security Partner Provider %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return resourceArmVirtualHubSecurityPartnerProviderRead(d, meta)
}

func resourceArmVirtualHubSecurityPartnerProviderDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.SecurityPartnerProviderClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualHubSecurityPartnerProviderID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("deleting Security Partner Provider %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on deleting future for Security Partner Provider %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return nil
}
