package network

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmVirtualHubConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmVirtualHubConnectionCreate,
		Read:   resourceArmVirtualHubConnectionRead,
		Delete: resourceArmVirtualHubConnectionDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: ValidateVirtualHubConnectionName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"virtual_hub_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: ValidateVirtualHubName,
			},

			"remote_virtual_network_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"hub_to_vitual_network_traffic_allowed": {
				Type:       schema.TypeBool,
				Optional:   true,
				ForceNew:   true,
				Deprecated: "This field has been deprecated since it is maintained internally in the implementation.",
			},

			"vitual_network_to_hub_gateways_traffic_allowed": {
				Type:       schema.TypeBool,
				Optional:   true,
				ForceNew:   true,
				Deprecated: "This field has been deprecated since it is maintained internally in the implementation.",
			},

			"internet_security_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceArmVirtualHubConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.HubVirtualNetworkConnectionClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	virtualHubName := d.Get("virtual_hub_name").(string)

	locks.ByName(virtualHubName, virtualHubResourceName)
	defer locks.UnlockByName(virtualHubName, virtualHubResourceName)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, virtualHubName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for present of existing Virtual Hub Connection %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_virtual_hub_connection", *existing.ID)
		}
	}

	connection := network.HubVirtualNetworkConnection{
		Name: utils.String(name),
		HubVirtualNetworkConnectionProperties: &network.HubVirtualNetworkConnectionProperties{
			RemoteVirtualNetwork: &network.SubResource{
				ID: utils.String(d.Get("remote_virtual_network_id").(string)),
			},
			EnableInternetSecurity: utils.Bool(d.Get("internet_security_enabled").(bool)),
		},
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, virtualHubName, name, connection)
	if err != nil {
		return fmt.Errorf("Error creating Virtual Hub Connection %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of Virtual Hub Connection %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, virtualHubName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Virtual Hub Connection %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read Virtual Hub Connection %q (Resource Group %q) ID", name, resourceGroup)
	}
	d.SetId(*resp.ID)

	return resourceArmVirtualHubConnectionRead(d, meta)
}

func resourceArmVirtualHubConnectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.HubVirtualNetworkConnectionClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ParseVirtualHubConnectionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.VirtualHubName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Virtual Hub Connection %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Virtual Hub Connection %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("virtual_hub_name", id.VirtualHubName)

	if props := resp.HubVirtualNetworkConnectionProperties; props != nil {
		d.Set("internet_security_enabled", props.EnableInternetSecurity)
		remoteVirtualNetworkId := ""
		if props.RemoteVirtualNetwork != nil && props.RemoteVirtualNetwork.ID != nil {
			remoteVirtualNetworkId = *props.RemoteVirtualNetwork.ID
		}
		d.Set("remote_virtual_network_id", remoteVirtualNetworkId)
	}

	return nil
}

func resourceArmVirtualHubConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.HubVirtualNetworkConnectionClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ParseVirtualHubConnectionID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.VirtualHubName, virtualHubResourceName)
	defer locks.UnlockByName(id.VirtualHubName, virtualHubResourceName)

	future, err := client.Delete(ctx, id.ResourceGroup, id.VirtualHubName, id.Name)
	if err != nil {
		return fmt.Errorf("Error deleting Virtual Hub Connection %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error waiting for deleting Virtual Hub Connection %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}
	}

	return nil
}
