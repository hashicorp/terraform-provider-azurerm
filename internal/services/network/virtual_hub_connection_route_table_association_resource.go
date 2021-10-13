package network

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2021-02-01/network"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

const virtualHubConnectionResourceName = "azurerm_virtual_hub_connection"

func resourceVirtualHubConnectionRouteTableAssociation() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualHubConnectionRouteTableAssociationCreate,
		Read:   resourceVirtualHubConnectionRouteTableAssociationRead,
		Delete: resourceVirtualHubConnectionRouteTableAssociationDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.SubnetID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"virtual_hub_connection_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"route_table_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},
		},
	}
}

func resourceVirtualHubConnectionRouteTableAssociationCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.HubVirtualNetworkConnectionClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Virtual Hub Connection <-> Route Table Association creation.")

	virtualHubConnectionId := d.Get("virtual_hub_connection_id").(string)
	routeTableId := d.Get("route_table_id").(string)

	parsedVirtualHubConnectionId, err := parse.HubVirtualNetworkConnectionID(virtualHubConnectionId)
	if err != nil {
		return err
	}

	// might need to lock the virtual hub, but not sure yet.

	locks.ByName(parsedVirtualHubConnectionId.Name, virtualHubConnectionResourceName)
	defer locks.UnlockByName(parsedVirtualHubConnectionId.Name, virtualHubConnectionResourceName)

	resourceGroup := parsedVirtualHubConnectionId.ResourceGroup
	virtualHubName := parsedVirtualHubConnectionId.VirtualHubName
	connectionName := parsedVirtualHubConnectionId.Name

	virtualHubConnection, err := client.Get(ctx, resourceGroup, virtualHubName, connectionName)
	if err != nil {
		if utils.ResponseWasNotFound(virtualHubConnection.Response) {
			return fmt.Errorf("virtual hub connection %q (Virtual Hub %q / Resource Group %q) was not found", connectionName, virtualHubName, resourceGroup)
		}

		return fmt.Errorf("retrieving Virtual Hub Connection %q (Virtual Hub %q / Resource Group %q): %+v", connectionName, virtualHubName, resourceGroup, err)
	}

	if props := virtualHubConnection.RoutingConfiguration; props != nil {
		if associatedRouteTable := props.AssociatedRouteTable; associatedRouteTable != nil {
			if associatedRouteTable.ID != nil {
				// do not raise an import error if the connection is associated to a defaultRouteTable and nonRouteTable
				if !strings.HasSuffix(*associatedRouteTable.ID, "/defaultRouteTable") && !strings.HasSuffix(*associatedRouteTable.ID, "/noneRouteTable") {
					return tf.ImportAsExistsError("azurerm_virtual_hub_connection_route_table_association", virtualHubConnectionId)
				}
			}
		}

		virtualHubConnection.RoutingConfiguration.AssociatedRouteTable = &network.SubResource{
			ID: utils.String(routeTableId),
		}
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, virtualHubName, connectionName, virtualHubConnection)
	if err != nil {
		return fmt.Errorf("updating Virtual Hub Association for Route Table %q (Virtual Hub %q / Resource Group %q): %+v", connectionName, virtualHubName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for completion of Virtual Hub Association for Route Table %q (Virtual Hub %q / Resource Group %q): %+v", connectionName, virtualHubName, resourceGroup, err)
	}

	timeout, _ := ctx.Deadline()

	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{string(network.ProvisioningStateUpdating)},
		Target:     []string{string(network.ProvisioningStateSucceeded)},
		Refresh:    VirtualHubConnectionProvisioningStateRefreshFunc(ctx, client, *parsedVirtualHubConnectionId),
		MinTimeout: 1 * time.Minute,
		Timeout:    time.Until(timeout),
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for provisioning state of subnet for Virtual Hub Association for Route Table %q (Virtual Hub %q / Resource Group %q): %+v", connectionName, virtualHubName, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, virtualHubName, connectionName)
	if err != nil {
		return fmt.Errorf("retrieving Virtual Hub Connection %q (Virtual Hub %q / Resource Group %q): %+v", connectionName, virtualHubName, resourceGroup, err)
	}

	d.SetId(*read.ID)

	return resourceVirtualHubConnectionRouteTableAssociationRead(d, meta)
}

func resourceVirtualHubConnectionRouteTableAssociationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.HubVirtualNetworkConnectionClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	parsedVirtualHubConnectionId, err := parse.HubVirtualNetworkConnectionID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := parsedVirtualHubConnectionId.ResourceGroup
	virtualHubName := parsedVirtualHubConnectionId.VirtualHubName
	virtualHubConnectionName := parsedVirtualHubConnectionId.Name

	resp, err := client.Get(ctx, resourceGroup, virtualHubName, virtualHubConnectionName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Virtual Hub Connection %q (Virtual Hub %q / Resource Group %q) could not be found - removing from state!", virtualHubConnectionName, virtualHubName, resourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Virtual Hub Connection %q (Virtual Hub %q / Resource Group %q): %+v", virtualHubConnectionName, virtualHubName, resourceGroup, err)
	}

	props := resp.RoutingConfiguration
	if props == nil {
		return fmt.Errorf("`RoutingConfiguration` was nil for Virtual Hub Connection %q (Virtual Hub %q / Resource Group %q)", virtualHubConnectionName, virtualHubName, resourceGroup)
	}

	associatedRouteTable := props.AssociatedRouteTable
	if associatedRouteTable == nil {
		log.Printf("[DEBUG] Virtual Hub Connection %q (Virtual Hub %q / Resource Group %q) doesn't have an Associated Route Table - removing from state!", virtualHubConnectionName, virtualHubName, resourceGroup)
		d.SetId("")
		return nil
	}

	d.Set("virtual_hub_connection_id", resp.ID)
	d.Set("route_table_id", associatedRouteTable.ID)

	return nil
}

func resourceVirtualHubConnectionRouteTableAssociationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.HubVirtualNetworkConnectionClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	parsedVirtualHubConnectionId, err := parse.HubVirtualNetworkConnectionID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := parsedVirtualHubConnectionId.ResourceGroup
	virtualHubName := parsedVirtualHubConnectionId.VirtualHubName
	virtualHubConnectionName := parsedVirtualHubConnectionId.Name

	// retrieve the connection
	read, err := client.Get(ctx, resourceGroup, virtualHubName, virtualHubConnectionName)
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			log.Printf("[DEBUG] Virtual Hub Connection %q (Virtual Hub %q / Resource Group %q) could not be found - removing from state!", virtualHubConnectionName, virtualHubName, resourceGroup)
			return nil
		}

		return fmt.Errorf("retrieving Virtual Hub Connection %q (Virtual Hub %q / Resource Group %q): %+v", virtualHubConnectionName, virtualHubName, resourceGroup, err)
	}

	props := read.RoutingConfiguration
	if props == nil {
		return fmt.Errorf("`Properties` was nil for Virtual Hub Connection %q (Virtual Hub %q / Resource Group %q)", virtualHubConnectionName, virtualHubName, resourceGroup)
	}

	if props.AssociatedRouteTable == nil || props.AssociatedRouteTable.ID == nil {
		log.Printf("[DEBUG] Virtual Hub Connection %q (Virtual Hub %q / Resource Group %q) has no Route Table - removing from state!", virtualHubConnectionName, virtualHubName, resourceGroup)
		return nil
	}

	locks.ByName(virtualHubName, virtualHubResourceName)
	defer locks.UnlockByName(virtualHubName, virtualHubResourceName)

	locks.ByName(virtualHubConnectionName, virtualHubConnectionResourceName)
	defer locks.UnlockByName(virtualHubConnectionName, virtualHubConnectionResourceName)

	// then re-retrieve it to ensure we've got the latest state
	read, err = client.Get(ctx, resourceGroup, virtualHubName, virtualHubConnectionName)
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			log.Printf("[DEBUG] Virtual Hub Connection %q (Virtual Hub %q / Resource Group %q) could not be found - removing from state!", virtualHubConnectionName, virtualHubName, resourceGroup)
			return nil
		}

		return fmt.Errorf("retrieving Virtual Hub Connection %q (Virtual Hub %q / Resource Group %q): %+v", virtualHubConnectionName, virtualHubName, resourceGroup, err)
	}

	read.RoutingConfiguration.AssociatedRouteTable = nil

	future, err := client.CreateOrUpdate(ctx, resourceGroup, virtualHubName, virtualHubConnectionName, read)
	if err != nil {
		return fmt.Errorf("removing Route Table Association from Virtual Hub Connection %q (Virtual Hub %q / Resource Group %q): %+v", virtualHubConnectionName, virtualHubName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for removal of Route Table Association from Virtual Hub Connection %q (Virtual Hub %q / Resource Group %q): %+v", virtualHubConnectionName, virtualHubName, resourceGroup, err)
	}

	return nil
}
