package network

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-07-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceIpGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceIpGroupCreateUpdate,
		Read:   resourceIpGroupRead,
		Update: resourceIpGroupCreateUpdate,
		Delete: resourceIpGroupDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.IpGroupID(id)
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

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"cidrs": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				Set: schema.HashString,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceIpGroupCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.IPGroupsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing IP Group %q (Resource Group %q): %s", name, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_ip_group", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})
	ipAddresses := d.Get("cidrs").(*schema.Set).List()

	sg := network.IPGroup{
		Name:     &name,
		Location: &location,
		IPGroupPropertiesFormat: &network.IPGroupPropertiesFormat{
			IPAddresses: utils.ExpandStringSlice(ipAddresses),
		},
		Tags: tags.Expand(t),
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, name, sg)
	if err != nil {
		return fmt.Errorf("creating/updating IP Group %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the completion of IP Group %q (Resource Group %q): %+v", name, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, name, "")
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("cannot read IP Group %q (resource group %q) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceIpGroupRead(d, meta)
}

func resourceIpGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.IPGroupsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.IpGroupID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on IP Group %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.IPGroupPropertiesFormat; props != nil {
		if props.IPAddresses == nil {
			return fmt.Errorf("list of ipAddresses returned is nil")
		}
		if err := d.Set("cidrs", props.IPAddresses); err != nil {
			return fmt.Errorf("setting `cidrs`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceIpGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.IPGroupsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.IpGroupID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("deleting IP Group %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("deleting IP Group %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return err
}
