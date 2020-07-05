package network

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-03-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
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

			"hub_to_vitual_network_traffic_allowed": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"vitual_network_to_hub_gateways_traffic_allowed": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
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
	client := meta.(*clients.Client).Network.VirtualHubClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := ParseVirtualHubID(d.Get("virtual_hub_id").(string))
	if err != nil {
		return err
	}

	locks.ByName(id.Name, virtualHubResourceName)
	defer locks.UnlockByName(id.Name, virtualHubResourceName)

	virtualHub, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(virtualHub.Response) {
			return fmt.Errorf("Virtual Hub %q was not found in Resource Group %q", id.Name, id.ResourceGroup)
		}

		return fmt.Errorf("Error retrieving Virtual Hub %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if virtualHub.VirtualHubProperties == nil {
		return fmt.Errorf("Error retrieving Virtual Hub %q (Resource Group %q): `properties` was nil", id.Name, id.ResourceGroup)
	}

	name := d.Get("name").(string)

	if features.ShouldResourcesBeImported() {
		if connection, _ := findVirtualHubConnection(name, virtualHub); connection != nil {
			return tf.ImportAsExistsError("azurerm_virtual_hub_connection", *connection.ID)
		}
	}

	props := *virtualHub.VirtualHubProperties

	var connections []network.HubVirtualNetworkConnection
	if props.VirtualNetworkConnections != nil {
		connections = *props.VirtualNetworkConnections
	}

	connection := network.HubVirtualNetworkConnection{
		Name: utils.String(name),
		HubVirtualNetworkConnectionProperties: &network.HubVirtualNetworkConnectionProperties{
			RemoteVirtualNetwork: &network.SubResource{
				ID: utils.String(d.Get("remote_virtual_network_id").(string)),
			},
			AllowHubToRemoteVnetTransit:         utils.Bool(d.Get("hub_to_vitual_network_traffic_allowed").(bool)),
			AllowRemoteVnetToUseHubVnetGateways: utils.Bool(d.Get("vitual_network_to_hub_gateways_traffic_allowed").(bool)),
			EnableInternetSecurity:              utils.Bool(d.Get("internet_security_enabled").(bool)),
		},
	}
	connections = append(connections, connection)
	virtualHub.VirtualHubProperties.VirtualNetworkConnections = &connections

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, virtualHub)
	if err != nil {
		return fmt.Errorf("Error adding Connection %q to Virtual Hub %q (Resource Group %q): %+v", name, id.Name, id.ResourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for addition of Connection %q to Virtual Hub %q (Resource Group %q): %+v", name, id.Name, id.ResourceGroup, err)
	}

	virtualHub, err = client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("Error retrieving Virtual Hub %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	newConnection, err := findVirtualHubConnection(name, virtualHub)
	if err != nil {
		return fmt.Errorf("Error retrieving Virtual Hub %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	if newConnection == nil {
		return fmt.Errorf("Connection %q was not found in Virtual Hub %q / Resource Group %q", name, id.Name, id.ResourceGroup)
	}
	if newConnection.ID == nil && *newConnection.ID == "" {
		return fmt.Errorf("Error retrieving Connection %q (Virtual Hub %q / Resource Group %q): `id` was nil or empty", name, id.Name, id.ResourceGroup)
	}

	d.SetId(*newConnection.ID)
	return resourceArmVirtualHubConnectionRead(d, meta)
}

func resourceArmVirtualHubConnectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualHubClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ParseVirtualHubConnectionID(d.Id())
	if err != nil {
		return err
	}

	virtualHub, err := client.Get(ctx, id.ResourceGroup, id.VirtualHubName)
	if err != nil {
		if utils.ResponseWasNotFound(virtualHub.Response) {
			log.Printf("[DEBUG] Virtual Hub %q was not found in Resource Group %q - so Connection %q can't exist - removing from state", id.Name, id.VirtualHubName, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Virtual Hub %q (Resource Group %q): %+v", id.VirtualHubName, id.ResourceGroup, err)
	}

	connection, err := findVirtualHubConnection(id.Name, virtualHub)
	if err != nil {
		return fmt.Errorf("Error retrieving Connection %q (Virtual Hub %q / Resource Group %q): %+v`", id.Name, id.VirtualHubName, id.ResourceGroup, err)
	}
	if connection == nil {
		log.Printf("[DEBUG] Connection %q was not found within Virtual Hub %q (Resource Group %q) - removing from state", id.Name, id.VirtualHubName, id.ResourceGroup)
		d.SetId("")
		return nil
	}
	if connection.HubVirtualNetworkConnectionProperties == nil {
		return fmt.Errorf("Error retrieving Connection %q (Virtual Hub %q / Resource Group %q): `properties` was nil`", id.Name, id.VirtualHubName, id.ResourceGroup)
	}

	d.Set("name", id.Name)
	d.Set("virtual_hub_id", virtualHub.ID)

	if props := connection.HubVirtualNetworkConnectionProperties; props != nil {
		d.Set("hub_to_vitual_network_traffic_allowed", props.AllowHubToRemoteVnetTransit)
		d.Set("vitual_network_to_hub_gateways_traffic_allowed", props.AllowRemoteVnetToUseHubVnetGateways)
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
	client := meta.(*clients.Client).Network.VirtualHubClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ParseVirtualHubConnectionID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.VirtualHubName, virtualHubResourceName)
	defer locks.UnlockByName(id.VirtualHubName, virtualHubResourceName)

	virtualHub, err := client.Get(ctx, id.ResourceGroup, id.VirtualHubName)
	if err != nil {
		if utils.ResponseWasNotFound(virtualHub.Response) {
			return fmt.Errorf("Virtual Hub %q was not found in Resource Group %q", id.VirtualHubName, id.ResourceGroup)
		}

		return fmt.Errorf("Error retrieving Virtual Hub %q (Resource Group %q): %+v", id.VirtualHubName, id.ResourceGroup, err)
	}

	if virtualHub.VirtualHubProperties == nil {
		return fmt.Errorf("Error retrieving Virtual Hub %q (Resource Group %q): `properties` was nil", id.VirtualHubName, id.ResourceGroup)
	}

	if virtualHub.VirtualHubProperties.VirtualNetworkConnections == nil {
		return fmt.Errorf("Error retrieving Virtual Hub %q (Resource Group %q): `properties.VirtualNetworkConnections` was nil", id.VirtualHubName, id.ResourceGroup)
	}

	var newConnections []network.HubVirtualNetworkConnection
	for _, connection := range *virtualHub.VirtualHubProperties.VirtualNetworkConnections {
		if connection.Name == nil {
			continue
		}

		if *connection.Name == id.Name {
			continue
		}

		newConnections = append(newConnections, connection)
	}
	virtualHub.VirtualHubProperties.VirtualNetworkConnections = &newConnections

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.VirtualHubName, virtualHub)
	if err != nil {
		return fmt.Errorf("Error removing Connection %q to Virtual Hub %q (Resource Group %q): %+v", id.Name, id.VirtualHubName, id.ResourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for removal of Connection %q to Virtual Hub %q (Resource Group %q): %+v", id.Name, id.VirtualHubName, id.ResourceGroup, err)
	}

	return nil
}

func findVirtualHubConnection(name string, virtualHub network.VirtualHub) (*network.HubVirtualNetworkConnection, error) {
	if virtualHub.VirtualHubProperties == nil {
		return nil, fmt.Errorf("`properties` was nil")
	}
	if virtualHub.VirtualHubProperties.VirtualNetworkConnections == nil {
		return nil, fmt.Errorf("`properties.VirtualNetworkConnections` was nil")
	}

	connections := *virtualHub.VirtualHubProperties.VirtualNetworkConnections

	for _, connection := range connections {
		if connection.Name == nil || connection.ID == nil {
			continue
		}

		if *connection.Name == name {
			return &connection, nil
		}
	}

	return nil, nil
}
