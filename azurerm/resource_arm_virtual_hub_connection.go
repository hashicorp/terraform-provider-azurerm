package azurerm

import (
	"fmt"
	"log"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-09-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	networkSvc "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmVirtualHubConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmVirtualHubConnectionCreate,
		Read:   resourceArmVirtualHubConnectionRead,
		Update: resourceArmVirtualHubConnectionUpdate,
		Delete: resourceArmVirtualHubConnectionDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`^[\da-zA-Z][-_.\da-zA-Z]{0,78}[_\da-zA-Z]$`),
					`The name must be between 1 and 80 characters and begin with a letter or number, end with a letter, number or underscore, and may contain only letters, numbers, underscores, periods, or hyphens.`,
				),
			},

			"virtual_hub_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: networkSvc.ValidateVirtualHubID,
			},

			"remote_virtual_network_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				// TODO: more specific validation
				ValidateFunc: azure.ValidateResourceID,
			},

			"allow_hub_to_remote_vnet_transit": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"allow_remote_vnet_to_use_hub_vnet_gateways": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"enable_internet_security": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceArmVirtualHubConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Network.VirtualHubClient
	ctx := meta.(*ArmClient).StopContext

	virtualHubId, err := networkSvc.ParseVirtualHubID(d.Get("virtual_hub_id").(string))
	if err != nil {
		return err
	}

	resourceGroup := virtualHubId.Base.ResourceGroup
	virtualHubName := virtualHubId.Name

	locks.ByName(virtualHubName, virtualHubResourceName)
	defer locks.UnlockByName(virtualHubName, virtualHubResourceName)

	virtualHub, err := client.Get(ctx, resourceGroup, virtualHubName)
	if err != nil {
		if utils.ResponseWasNotFound(virtualHub.Response) {
			return fmt.Errorf("Virtual Hub %q was not found in Resource Group %q", virtualHubName, resourceGroup)
		}

		return fmt.Errorf("Error retrieving Virtual Hub %q (Resource Group %q): %+v", virtualHubName, resourceGroup, err)
	}

	if virtualHub.VirtualHubProperties == nil {
		return fmt.Errorf("Error retrieving Virtual Hub %q (Resource Group %q): `properties` was nil", virtualHubName, resourceGroup)
	}

	props := *virtualHub.VirtualHubProperties

	var connections []network.HubVirtualNetworkConnection
	if props.VirtualNetworkConnections != nil {
		connections = *props.VirtualNetworkConnections
	}

	name := d.Get("name").(string)

	if features.ShouldResourcesBeImported() {
		if connection, _ := findVirtualHubConnection(name, connections); connection != nil {
			return tf.ImportAsExistsError("azurerm_virtual_hub_connection", *connection.ID)
		}
	}

	connection := network.HubVirtualNetworkConnection{
		Name: utils.String(name),
		HubVirtualNetworkConnectionProperties: &network.HubVirtualNetworkConnectionProperties{
			RemoteVirtualNetwork: &network.SubResource{
				ID: utils.String(d.Get("remote_virtual_network_id").(string)),
			},
			AllowHubToRemoteVnetTransit:         utils.Bool(d.Get("allow_hub_to_remote_vnet_transit").(bool)),
			AllowRemoteVnetToUseHubVnetGateways: utils.Bool(d.Get("allow_remote_vnet_to_use_hub_vnet_gateways").(bool)),
			EnableInternetSecurity:              utils.Bool(d.Get("enable_internet_security").(bool)),
		},
	}
	connections = append(connections, connection)
	virtualHub.VirtualHubProperties.VirtualNetworkConnections = &connections

	future, err := client.CreateOrUpdate(ctx, resourceGroup, virtualHubName, virtualHub)
	if err != nil {
		return fmt.Errorf("Error adding Connection %q to Virtual Hub %q (Resource Group %q): %+v", name, virtualHubName, resourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for addition of Connection %q to Virtual Hub %q (Resource Group %q): %+v", name, virtualHubName, resourceGroup, err)
	}

	virtualHub, err = client.Get(ctx, resourceGroup, virtualHubName)
	if err != nil {
		return fmt.Errorf("Error retrieving Virtual Hub %q (Resource Group %q): %+v", virtualHubName, resourceGroup, err)
	}
	if virtualHub.VirtualHubProperties == nil {
		return fmt.Errorf("Error retrieving Virtual Hub %q (Resource Group %q): `properties` was nil", virtualHubName, resourceGroup)
	}
	if virtualHub.VirtualHubProperties.VirtualNetworkConnections == nil {
		return fmt.Errorf("Error retrieving Virtual Hub %q (Resource Group %q): `properties.VirtualNetworkConnections` was nil", virtualHubName, resourceGroup)
	}

	connections = *virtualHub.VirtualHubProperties.VirtualNetworkConnections
	newConnection, _ := findVirtualHubConnection(name, connections)
	if newConnection == nil {
		return fmt.Errorf("Connection %q was not found in Virtual Hub %q / Resource Group %q", name, virtualHubName, resourceGroup)
	}
	if newConnection.ID == nil {
		return fmt.Errorf("Error retrieving Connection %q (Virtual Hub %q / Resource Group %q): `id` was nil", name, virtualHubName, resourceGroup)
	}

	d.SetId(*newConnection.ID)
	return resourceArmVirtualHubConnectionRead(d, meta)
}

func resourceArmVirtualHubConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Network.VirtualHubClient
	ctx := meta.(*ArmClient).StopContext

	virtualHubConnectionId, err := networkSvc.ParseVirtualHubConnectionID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := virtualHubConnectionId.Base.ResourceGroup
	virtualHubName := virtualHubConnectionId.VirtualHubName
	name := virtualHubConnectionId.Name

	locks.ByName(virtualHubName, virtualHubResourceName)
	defer locks.UnlockByName(virtualHubName, virtualHubResourceName)

	virtualHub, err := client.Get(ctx, resourceGroup, virtualHubName)
	if err != nil {
		if utils.ResponseWasNotFound(virtualHub.Response) {
			return fmt.Errorf("Virtual Hub %q was not found in Resource Group %q", virtualHubName, resourceGroup)
		}

		return fmt.Errorf("Error retrieving Virtual Hub %q (Resource Group %q): %+v", virtualHubName, resourceGroup, err)
	}

	if virtualHub.VirtualHubProperties == nil {
		return fmt.Errorf("Error retrieving Virtual Hub %q (Resource Group %q): `properties` was nil", virtualHubName, resourceGroup)
	}

	if virtualHub.VirtualHubProperties.VirtualNetworkConnections == nil {
		return fmt.Errorf("Error retrieving Virtual Hub %q (Resource Group %q): `properties.VirtualNetworkConnections` was nil", virtualHubName, resourceGroup)
	}
	connections := *virtualHub.VirtualHubProperties.VirtualNetworkConnections

	connection, index := findVirtualHubConnection(name, connections)
	if connection == nil {
		return fmt.Errorf("Connection %q was not found within Virtual Hub %q (Resource Group %q): %+v", name, virtualHubName, resourceGroup, err)
	}
	if connection.HubVirtualNetworkConnectionProperties == nil {
		return fmt.Errorf("Error retrieving Connection %q (Virtual Hub %q / Resource Group %q): `properties` was nil`", name, virtualHubName, resourceGroup)
	}

	connectionProps := *connection.HubVirtualNetworkConnectionProperties

	if d.HasChange("allow_hub_to_remote_vnet_transit") {
		connectionProps.AllowHubToRemoteVnetTransit = utils.Bool(d.Get("allow_hub_to_remote_vnet_transit").(bool))
	}

	if d.HasChange("allow_remote_vnet_to_use_hub_vnet_gateways") {
		connectionProps.AllowRemoteVnetToUseHubVnetGateways = utils.Bool(d.Get("allow_remote_vnet_to_use_hub_vnet_gateways").(bool))
	}

	if d.HasChange("enable_internet_security") {
		connectionProps.EnableInternetSecurity = utils.Bool(d.Get("enable_internet_security").(bool))
	}
	connections[index] = *connection
	virtualHub.VirtualHubProperties.VirtualNetworkConnections = &connections

	future, err := client.CreateOrUpdate(ctx, resourceGroup, virtualHubName, virtualHub)
	if err != nil {
		return fmt.Errorf("Error updating Connection %q to Virtual Hub %q (Resource Group %q): %+v", name, virtualHubName, resourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for update of Connection %q to Virtual Hub %q (Resource Group %q): %+v", name, virtualHubName, resourceGroup, err)
	}

	return resourceArmVirtualHubConnectionRead(d, meta)
}

func resourceArmVirtualHubConnectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Network.VirtualHubClient
	ctx := meta.(*ArmClient).StopContext

	virtualHubConnectionId, err := networkSvc.ParseVirtualHubConnectionID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := virtualHubConnectionId.Base.ResourceGroup
	virtualHubName := virtualHubConnectionId.VirtualHubName
	name := virtualHubConnectionId.Name

	locks.ByName(virtualHubName, virtualHubResourceName)
	defer locks.UnlockByName(virtualHubName, virtualHubResourceName)

	virtualHub, err := client.Get(ctx, resourceGroup, virtualHubName)
	if err != nil {
		if utils.ResponseWasNotFound(virtualHub.Response) {
			log.Printf("[DEBUG] Virtual Hub %q was not found in Resource Group %q - so Connection %q can't exist - removing from state", name, virtualHubName, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Virtual Hub %q (Resource Group %q): %+v", virtualHubName, resourceGroup, err)
	}

	if virtualHub.VirtualHubProperties == nil {
		return fmt.Errorf("Error retrieving Virtual Hub %q (Resource Group %q): `properties` was nil", virtualHubName, resourceGroup)
	}

	if virtualHub.VirtualHubProperties.VirtualNetworkConnections == nil {
		return fmt.Errorf("Error retrieving Virtual Hub %q (Resource Group %q): `properties.VirtualNetworkConnections` was nil", virtualHubName, resourceGroup)
	}
	connections := *virtualHub.VirtualHubProperties.VirtualNetworkConnections

	connection, _ := findVirtualHubConnection(name, connections)
	if connection == nil {
		log.Printf("[DEBUG] Connection %q was not found within Virtual Hub %q (Resource Group %q) - removing from state", name, virtualHubName, resourceGroup)
		d.SetId("")
		return nil
	}
	if connection.HubVirtualNetworkConnectionProperties == nil {
		return fmt.Errorf("Error retrieving Connection %q (Virtual Hub %q / Resource Group %q): `properties` was nil`", name, virtualHubName, resourceGroup)
	}
	props := *connection.HubVirtualNetworkConnectionProperties

	d.Set("name", name)
	d.Set("virtual_hub_id", virtualHub.ID)
	d.Set("allow_hub_to_remote_vnet_transit", props.AllowHubToRemoteVnetTransit)
	d.Set("allow_remote_vnet_to_use_hub_vnet_gateways", props.AllowRemoteVnetToUseHubVnetGateways)
	d.Set("enable_internet_security", props.EnableInternetSecurity)

	remoteVirtualNetworkId := ""
	if props.RemoteVirtualNetwork != nil && props.RemoteVirtualNetwork.ID != nil {
		remoteVirtualNetworkId = *props.RemoteVirtualNetwork.ID
	}
	d.Set("remote_virtual_network_id", remoteVirtualNetworkId)

	return nil
}

func resourceArmVirtualHubConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Network.VirtualHubClient
	ctx := meta.(*ArmClient).StopContext

	virtualHubConnectionId, err := networkSvc.ParseVirtualHubConnectionID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := virtualHubConnectionId.Base.ResourceGroup
	virtualHubName := virtualHubConnectionId.VirtualHubName
	name := virtualHubConnectionId.Name

	locks.ByName(virtualHubName, virtualHubResourceName)
	defer locks.UnlockByName(virtualHubName, virtualHubResourceName)

	virtualHub, err := client.Get(ctx, resourceGroup, virtualHubName)
	if err != nil {
		if utils.ResponseWasNotFound(virtualHub.Response) {
			return fmt.Errorf("Virtual Hub %q was not found in Resource Group %q", virtualHubName, resourceGroup)
		}

		return fmt.Errorf("Error retrieving Virtual Hub %q (Resource Group %q): %+v", virtualHubName, resourceGroup, err)
	}

	if virtualHub.VirtualHubProperties == nil {
		return fmt.Errorf("Error retrieving Virtual Hub %q (Resource Group %q): `properties` was nil", virtualHubName, resourceGroup)
	}

	if virtualHub.VirtualHubProperties.VirtualNetworkConnections == nil {
		return fmt.Errorf("Error retrieving Virtual Hub %q (Resource Group %q): `properties.VirtualNetworkConnections` was nil", virtualHubName, resourceGroup)
	}

	var newConnections []network.HubVirtualNetworkConnection
	for _, connection := range *virtualHub.VirtualHubProperties.VirtualNetworkConnections {
		if connection.Name == nil {
			continue
		}

		if *connection.Name == name {
			continue
		}

		newConnections = append(newConnections, connection)
	}
	virtualHub.VirtualHubProperties.VirtualNetworkConnections = &newConnections

	future, err := client.CreateOrUpdate(ctx, resourceGroup, virtualHubName, virtualHub)
	if err != nil {
		return fmt.Errorf("Error removing Connection %q to Virtual Hub %q (Resource Group %q): %+v", name, virtualHubName, resourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for removal of Connection %q to Virtual Hub %q (Resource Group %q): %+v", name, virtualHubName, resourceGroup, err)
	}

	return nil
}

func findVirtualHubConnection(name string, connections []network.HubVirtualNetworkConnection) (conn *network.HubVirtualNetworkConnection, index int) {
	for i, connection := range connections {
		if connection.Name == nil || connection.ID == nil {
			continue
		}

		if *connection.Name == name {
			conn = &connection
			index = i
			return
		}
	}

	return
}
