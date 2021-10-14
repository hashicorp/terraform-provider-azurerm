package network

import (
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2021-02-01/network"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

const virtualHubConnectionResourceName = "azurerm_virtual_hub_connection"

func resourceVirtualHubConnectionConfiguration() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualHubConnectionConfigurationCreateOrUpdate,
		Read:   resourceVirtualHubConnectionConfigurationRead,
		Update: resourceVirtualHubConnectionConfigurationCreateOrUpdate,
		Delete: resourceVirtualHubConnectionConfigurationDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.HubVirtualNetworkConnectionID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"virtual_hub_connection_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"associated_route_table_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.HubRouteTableID,
			},

			"propagated_route_table": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"labels": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							AtLeastOneOf: []string{"propagated_route_table.0.labels", "propagated_route_table.0.route_table_ids"},
						},

						"route_table_ids": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validate.HubRouteTableID,
							},
							AtLeastOneOf: []string{"propagated_route_table.0.labels", "propagated_route_table.0.route_table_ids"},
						},
					},
				},
			},

			//lintignore:XS003
			"static_vnet_route": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"address_prefixes": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.IsCIDR,
							},
						},

						"next_hop_ip_address": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsIPv4Address,
						},
					},
				},
			},
		},
	}
}

func resourceVirtualHubConnectionConfigurationCreateOrUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.HubVirtualNetworkConnectionClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	virtualHubConnectionId := d.Get("virtual_hub_connection_id").(string)

	id, err := parse.HubVirtualNetworkConnectionID(virtualHubConnectionId)
	if err != nil {
		return err
	}

	locks.ByName(id.VirtualHubName, virtualHubResourceName)
	defer locks.UnlockByName(id.VirtualHubName, virtualHubResourceName)

	locks.ByName(id.Name, virtualHubConnectionResourceName)
	defer locks.UnlockByName(id.Name, virtualHubConnectionResourceName)

	virtualHubConnection, err := client.Get(ctx, id.ResourceGroup, id.VirtualHubName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(virtualHubConnection.Response) {
			return fmt.Errorf("virtual hub connection %q (Virtual Hub %q / Resource Group %q) was not found", id.Name, id.VirtualHubName, id.ResourceGroup)
		}

		return fmt.Errorf("retrieving Virtual Hub Connection %q (Virtual Hub %q / Resource Group %q): %+v", id.Name, id.VirtualHubName, id.ResourceGroup, err)
	}

	if d.IsNewResource() {
		if err := hasImportError(id, virtualHubConnection); err != nil {
			return err
		}
	}

	virtualHubConnection.HubVirtualNetworkConnectionProperties.RoutingConfiguration = expandVirtualHubConnectionConfigurationRouting(d)

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.VirtualHubName, id.Name, virtualHubConnection)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", virtualHubConnectionId, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of %s: %+v", virtualHubConnectionId, err)
	}

	timeout, _ := ctx.Deadline()

	vnetStateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{string(network.ProvisioningStateUpdating)},
		Target:     []string{string(network.ProvisioningStateSucceeded)},
		Refresh:    virtualHubConnectionProvisioningStateRefreshFunc(ctx, client, *id),
		MinTimeout: 1 * time.Minute,
		Timeout:    time.Until(timeout),
	}
	if _, err = vnetStateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for provisioning state of %s: %+v", virtualHubConnectionId, err)
	}

	d.SetId(virtualHubConnectionId)

	return resourceVirtualHubConnectionConfigurationRead(d, meta)
}

func resourceVirtualHubConnectionConfigurationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.HubVirtualNetworkConnectionClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.HubVirtualNetworkConnectionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.VirtualHubName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading %s: %+v", *id, err)
	}

	if props := resp.HubVirtualNetworkConnectionProperties; props != nil {
		if routing := props.RoutingConfiguration; routing != nil {
			associatedRouteTableId := ""
			if routing.AssociatedRouteTable != nil && routing.AssociatedRouteTable.ID != nil {
				associatedRouteTableId = *routing.AssociatedRouteTable.ID
			}
			d.Set("associated_route_table_id", associatedRouteTableId)

			d.Set("propagated_route_table", flattenVirtualHubConnectionPropagatedRouteTable(routing.PropagatedRouteTables))
			d.Set("static_vnet_route", flattenVirtualHubConnectionVnetStaticRoute(routing.VnetRoutes))
		}
	}

	return nil
}

func resourceVirtualHubConnectionConfigurationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.HubVirtualNetworkConnectionClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.HubVirtualNetworkConnectionID(d.Id())
	if err != nil {
		return err
	}

	// retrieve the connection
	read, err := client.Get(ctx, id.ResourceGroup, id.VirtualHubName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			log.Printf("[DEBUG] Virtual Hub Connection %q (Virtual Hub %q / Resource Group %q) could not be found - removing from state!", id.Name, id.VirtualHubName, id.ResourceGroup)
			return nil
		}

		return fmt.Errorf("retrieving Virtual Hub Connection %q (Virtual Hub %q / Resource Group %q): %+v", id.Name, id.VirtualHubName, id.ResourceGroup, err)
	}

	props := read.RoutingConfiguration
	if props == nil {
		return fmt.Errorf("`Properties` was nil for Virtual Hub Connection %q (Virtual Hub %q / Resource Group %q)", id.Name, id.VirtualHubName, id.ResourceGroup)
	}

	locks.ByName(id.VirtualHubName, virtualHubResourceName)
	defer locks.UnlockByName(id.VirtualHubName, virtualHubResourceName)

	locks.ByName(id.Name, virtualHubConnectionResourceName)
	defer locks.UnlockByName(id.Name, virtualHubConnectionResourceName)

	// then re-retrieve it to ensure we've got the latest state
	read, err = client.Get(ctx, id.ResourceGroup, id.VirtualHubName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			log.Printf("[DEBUG] Virtual Hub Connection %q (Virtual Hub %q / Resource Group %q) could not be found - removing from state!", id.Name, id.VirtualHubName, id.ResourceGroup)
			return nil
		}

		return fmt.Errorf("retrieving Virtual Hub Connection %q (Virtual Hub %q / Resource Group %q): %+v", id.Name, id.VirtualHubName, id.ResourceGroup, err)
	}

	read.RoutingConfiguration.AssociatedRouteTable = nil
	read.RoutingConfiguration.PropagatedRouteTables = nil
	read.RoutingConfiguration.VnetRoutes = nil

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.VirtualHubName, id.Name, read)
	if err != nil {
		return fmt.Errorf("removing Route Configuration from Virtual Hub Connection %q (Virtual Hub %q / Resource Group %q): %+v", id.Name, id.VirtualHubName, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for removal of Route Configuration from Virtual Hub Connection %q (Virtual Hub %q / Resource Group %q): %+v", id.Name, id.VirtualHubName, id.ResourceGroup, err)
	}

	return nil
}

func expandVirtualHubConnectionConfigurationRouting(d *pluginsdk.ResourceData) *network.RoutingConfiguration {
	result := network.RoutingConfiguration{}

	if v, ok := d.GetOk("associated_route_table_id"); ok && v != "" {
		result.AssociatedRouteTable = &network.SubResource{
			ID: utils.String(v.(string)),
		}
	}

	if vnetStaticRoute := d.Get("static_vnet_route").([]interface{}); len(vnetStaticRoute) > 0 {
		result.VnetRoutes = expandVirtualHubConnectionVnetStaticRoute(vnetStaticRoute)
	} else {
		result.VnetRoutes = nil
	}

	if propagatedRouteTable := d.Get("propagated_route_table").([]interface{}); len(propagatedRouteTable) > 0 {
		result.PropagatedRouteTables = expandVirtualHubConnectionPropagatedRouteTable(propagatedRouteTable)
	} else {
		result.PropagatedRouteTables = nil
	}

	return &result
}

func hasImportError(id *parse.HubVirtualNetworkConnectionId, connection network.HubVirtualNetworkConnection) error {
	props := connection.RoutingConfiguration
	if props == nil {
		return nil
	}
	defaultRouteTableID := parse.NewHubRouteTableID(id.SubscriptionId, id.ResourceGroup, id.VirtualHubName, "defaultRouteTable").ID()

	if associatedRouteTable := props.AssociatedRouteTable; associatedRouteTable != nil {
		if associatedRouteTable.ID != nil {
			// do not raise an import error if the connection is associated to a defaultRouteTable
			if *associatedRouteTable.ID != defaultRouteTableID {
				return tf.ImportAsExistsError("azurerm_virtual_hub_connection_configuration", id.ID())
			}
		}
	}

	if propagatedRouteTables := props.PropagatedRouteTables; propagatedRouteTables != nil {
		if propagatedRouteTables.Labels != nil {
			// raise if the labels aren't default
			if !reflect.DeepEqual(propagatedRouteTables.Labels, &[]string{"default"}) {
				return tf.ImportAsExistsError("azurerm_virtual_hub_connection_configuration", id.ID())
			}
		}
		if propagatedRouteTables.Ids != nil {

			r := network.SubResource{
				ID: utils.String(defaultRouteTableID),
			}

			// do not raise an import error if the connection is associated to a defaultRouteTable and nonRouteTable
			if !reflect.DeepEqual(propagatedRouteTables.Ids, &[]network.SubResource{r}) {
				return tf.ImportAsExistsError("azurerm_virtual_hub_connection_configuration", id.ID())
			}
		}
	}

	if vnetRoutes := props.VnetRoutes; vnetRoutes != nil {
		if staticRoutes := vnetRoutes.StaticRoutes; staticRoutes != nil {
			// raise an import error if the connection has existing static routes
			if len(*staticRoutes) > 0 {
				return tf.ImportAsExistsError("azurerm_virtual_hub_connection_configuration", id.ID())
			}
		}
	}

	return nil
}
