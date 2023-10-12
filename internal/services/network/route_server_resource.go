// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-07-01/network"
)

func resourceRouteServer() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceRouteServerCreateUpdate,
		Read:   resourceRouteServerRead,
		Update: resourceRouteServerCreateUpdate,
		Delete: resourceRouteServerDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.VirtualHubID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.RouteServerName(),
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"sku": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Standard"}, false),
			},

			"public_ip_address_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.PublicIpAddressID,
			},

			"subnet_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateSubnetID,
			},

			"branch_to_branch_traffic_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"virtual_router_ips": {
				Type:     pluginsdk.TypeSet,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"routing_state": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"virtual_router_asn": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceRouteServerCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	serverClient := meta.(*clients.Client).Network.VirtualHubClient
	ipClient := meta.(*clients.Client).Network.VirtualHubIPClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewVirtualHubID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	locks.ByName(id.Name, "azurerm_route_server")
	defer locks.UnlockByName(id.Name, "azurerm_route_server")

	if d.IsNewResource() {
		existing, err := serverClient.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for present of existing Route Server %q (Resource Group Name %q): %+v", id.Name, id.ResourceGroup, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_route_server", id.ID())
		}
	}

	location := location.Normalize(d.Get("location").(string))
	t := tags.Expand(d.Get("tags").(map[string]interface{}))

	parameters := network.VirtualHub{
		Location: utils.String(location),
		VirtualHubProperties: &network.VirtualHubProperties{
			Sku:                        utils.String(d.Get("sku").(string)),
			AllowBranchToBranchTraffic: utils.Bool(d.Get("branch_to_branch_traffic_enabled").(bool)),
		},
		Tags: t,
	}

	serverFuture, err := serverClient.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, parameters)
	if err != nil {
		return fmt.Errorf("creating Route Server %q (Resource Group Name %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err := serverFuture.WaitForCompletionRef(ctx, serverClient.Client); err != nil {
		return fmt.Errorf("waiting for creation of %s: %+v", id, err)
	}

	timeout, _ := ctx.Deadline()
	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"Provisioning", "Updating"},
		Target:                    []string{"Succeeded", "Provisioned"},
		Refresh:                   routeServerCreateRefreshFunc(ctx, serverClient, id),
		PollInterval:              15 * time.Second,
		ContinuousTargetOccurence: 5,
		Timeout:                   time.Until(timeout),
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("waiting for creation/update of Route Server %q (Resource Group Name %q): %+v", id.Name, id.ResourceGroup, err)
	}

	ipConfigName := "ipConfig1"
	ipConfigs := network.HubIPConfiguration{
		Name: utils.String(ipConfigName),
		HubIPConfigurationPropertiesFormat: &network.HubIPConfigurationPropertiesFormat{
			PublicIPAddress: &network.PublicIPAddress{
				ID: utils.String(d.Get("public_ip_address_id").(string)),
			},
			Subnet: &network.Subnet{
				ID: utils.String(d.Get("subnet_id").(string)),
			},
		},
	}

	future, err := ipClient.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, ipConfigName, ipConfigs)
	if err != nil {
		return fmt.Errorf("creating/updating IP Configuration %q of Route Server %q (Resource Group Name %q): %+v", ipConfigName, id.Name, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, ipClient.Client); err != nil {
		return fmt.Errorf("waiting on creation/update for IP Configuration %q of Route Server %q (Resource Group Name %q): %+v", ipConfigName, id.Name, id.ResourceGroup, err)
	}
	d.SetId(id.ID())

	return resourceRouteServerRead(d, meta)
}

func resourceRouteServerRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualHubClient
	ipClient := meta.(*clients.Client).Network.VirtualHubIPClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualHubID(d.Id())
	if err != nil {
		return err
	}

	routeServer, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(routeServer.Response) {
			log.Printf("[INFO] Route Server %s does not exists - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading Route Server %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := routeServer.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if props := routeServer.VirtualHubProperties; props != nil {
		d.Set("sku", props.Sku)
		var virtualRouterIps *[]string
		if props.VirtualRouterIps != nil {
			virtualRouterIps = props.VirtualRouterIps
		}
		d.Set("virtual_router_ips", virtualRouterIps)
		if props.AllowBranchToBranchTraffic != nil {
			d.Set("branch_to_branch_traffic_enabled", props.AllowBranchToBranchTraffic)
		}
		if props.VirtualRouterAsn != nil {
			d.Set("virtual_router_asn", props.VirtualRouterAsn)
		}
		d.Set("routing_state", string(props.RoutingState))
	}

	ipConfig, err := ipClient.List(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("retrieving IP Config for Router Server %q (Resource Group Name %q): %+v", id.Name, id.ResourceGroup, err)
	}

	for _, setting := range ipConfig.Values() {
		if ipProps := setting.HubIPConfigurationPropertiesFormat; ipProps != nil {
			if ipProps.PublicIPAddress != nil {
				d.Set("public_ip_address_id", ipProps.PublicIPAddress.ID)
			}
			if ipProps.Subnet != nil {
				d.Set("subnet_id", ipProps.Subnet.ID)
			}
		}
	}
	return tags.FlattenAndSet(d, routeServer.Tags)
}

func resourceRouteServerDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualHubClient
	ipClient := meta.(*clients.Client).Network.VirtualHubIPClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	routeServerId, err := parse.VirtualHubID(d.Id())
	if err != nil {
		return err
	}

	ipConfig, err := ipClient.List(ctx, routeServerId.ResourceGroup, routeServerId.Name)
	if err != nil {
		return fmt.Errorf("retrieving IP Config for Router Server %q (Resource Group Name %q): %+v", routeServerId.Name, routeServerId.ResourceGroup, err)
	}
	var ipName string
	for _, setting := range ipConfig.Values() {
		if setting.Name != nil {
			ipName = *setting.Name
		}
	}
	ipConfigId := parse.NewVirtualHubIpConfigurationID(routeServerId.SubscriptionId, routeServerId.ResourceGroup, routeServerId.Name, ipName)

	if ipConfig.Values() != nil {
		if err := deleteRouteServerIpConfiguration(ctx, ipClient, ipConfigId); err != nil {
			return err
		}
	}

	future, err := client.Delete(ctx, routeServerId.ResourceGroup, routeServerId.Name)
	if err != nil {
		return fmt.Errorf("deleting Route Server %q (Resource Group Name %q): %+v", routeServerId.Name, routeServerId.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("waiting for Route Server %q (Resource Group Name %q): %+v", routeServerId.Name, routeServerId.ResourceGroup, err)
		}
	}

	return nil
}

func deleteRouteServerIpConfiguration(ctx context.Context, client *network.VirtualHubIPConfigurationClient, id parse.VirtualHubIpConfigurationId) error {
	future, err := client.Delete(ctx, id.ResourceGroup, id.VirtualHubName, id.IpConfigurationName)
	if err != nil {
		return fmt.Errorf("deleting Router Server IP Config %s for Route Server %q (Resource Group Name %q): %+v", id.IpConfigurationName, id.VirtualHubName, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("waiting for deletion of Route Server IP Config %s: %+v", id.IpConfigurationName, err)
		}
	}
	timeout, _ := ctx.Deadline()
	stateConf := &pluginsdk.StateChangeConf{
		Pending: []string{"200"},
		Target:  []string{"404"},
		Refresh: ipConfigStateRefreshFunc(ctx, client, id),
		Timeout: time.Until(timeout),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for Router Server IP Config %s for Route Server %q (Resource Group Name %q): %+v", id.IpConfigurationName, id.VirtualHubName, id.ResourceGroup, err)
	}
	return nil
}

func routeServerCreateRefreshFunc(ctx context.Context, client *network.VirtualHubsClient, id parse.VirtualHubId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if utils.ResponseWasNotFound(res.Response) {
				return nil, "", fmt.Errorf("Route Server %q (Resource Group Name %q) does not exists", id.Name, id.ResourceGroup)
			}
			return nil, "", fmt.Errorf("retrieving Route Server %q (Resource Group Name %q) error", id.Name, id.ResourceGroup)
		}

		if res.VirtualHubProperties != nil {
			return res, string(res.VirtualHubProperties.ProvisioningState), nil
		}
		return nil, "", fmt.Errorf("unable to read the provisioning state of this Route Server %q (Resource Group Name %q)", id.Name, id.ResourceGroup)
	}
}

func ipConfigStateRefreshFunc(ctx context.Context, client *network.VirtualHubIPConfigurationClient, id parse.VirtualHubIpConfigurationId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id.ResourceGroup, id.VirtualHubName, id.IpConfigurationName)
		if err != nil {
			if utils.ResponseWasNotFound(res.Response) {
				return res, strconv.Itoa(res.StatusCode), nil
			}
			return nil, "", fmt.Errorf("polling for the status of route server ip config %s: %+v", id, err)
		}
		return res, strconv.Itoa(res.StatusCode), nil
	}
}
