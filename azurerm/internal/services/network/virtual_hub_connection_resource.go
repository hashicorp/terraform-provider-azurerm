package network

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"
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

			"virtual_hub_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"remote_virtual_network_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			// TODO 3.0: remove this property
			"hub_to_vitual_network_traffic_allowed": {
				Type:       schema.TypeBool,
				Optional:   true,
				ForceNew:   true,
				Deprecated: "Due to a breaking behavioural change in the Azure API this property is no longer functional and will be removed in version 3.0 of the provider",
			},

			// TODO 3.0: remove this property
			"vitual_network_to_hub_gateways_traffic_allowed": {
				Type:       schema.TypeBool,
				Optional:   true,
				ForceNew:   true,
				Deprecated: "Due to a breaking behavioural change in the Azure API this property is no longer functional and will be removed in version 3.0 of the provider",
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

	id, err := ParseVirtualHubID(d.Get("virtual_hub_id").(string))
	if err != nil {
		return err
	}

	locks.ByName(id.Name, virtualHubResourceName)
	defer locks.UnlockByName(id.Name, virtualHubResourceName)

	name := d.Get("name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for present of existing Virtual Hub Connection %q (Resource Group %q): %+v", name, id.ResourceGroup, err)
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

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, name, connection)
	if err != nil {
		return fmt.Errorf("creating Virtual Hub Connection %q (Resource Group %q): %+v", name, id.ResourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of Virtual Hub Connection %q (Resource Group %q): %+v", name, id.ResourceGroup, err)
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name, name)
	if err != nil {
		return fmt.Errorf("retrieving Virtual Hub Connection %q (Resource Group %q): %+v", name, id.ResourceGroup, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("Cannot read Virtual Hub Connection %q (Resource Group %q) ID", name, id.ResourceGroup)
	}
	d.SetId(*resp.ID)

	return resourceArmVirtualHubConnectionRead(d, meta)
}

func resourceArmVirtualHubConnectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.HubVirtualNetworkConnectionClient
	virtualHubClient := meta.(*clients.Client).Network.VirtualHubClient
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
		return fmt.Errorf("reading Virtual Hub Connection %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	virtualHubResp, err := virtualHubClient.Get(ctx, id.ResourceGroup, id.VirtualHubName)
	if err != nil {
		return fmt.Errorf("retrieving Virtual Hub %q (Resource Group %q): %+v", id.VirtualHubName, id.ResourceGroup, err)
	}
	if virtualHubResp.ID == nil || *virtualHubResp.ID == "" {
		return fmt.Errorf("Cannot read Virtual Hub %q (Resource Group %q) ID", id.VirtualHubName, id.ResourceGroup)
	}

	d.Set("name", id.Name)
	d.Set("virtual_hub_id", virtualHubResp.ID)

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
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ParseVirtualHubConnectionID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.VirtualHubName, virtualHubResourceName)
	defer locks.UnlockByName(id.VirtualHubName, virtualHubResourceName)

	future, err := client.Delete(ctx, id.ResourceGroup, id.VirtualHubName, id.Name)
	if err != nil {
		return fmt.Errorf("deleting Virtual Hub Connection %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deleting Virtual Hub Connection %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return nil
}
