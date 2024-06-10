// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceRouteServer() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceRouteServerCreate,
		Read:   resourceRouteServerRead,
		Update: resourceRouteServerUpdate,
		Delete: resourceRouteServerDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := virtualwans.ParseVirtualHubID(id)
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
				ValidateFunc: commonids.ValidatePublicIPAddressID,
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

func resourceRouteServerCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := virtualwans.NewVirtualHubID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	locks.ByName(id.VirtualHubName, "azurerm_route_server")
	defer locks.UnlockByName(id.VirtualHubName, "azurerm_route_server")

	existing, err := client.VirtualHubsGet(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_route_server", id.ID())
	}

	parameters := virtualwans.VirtualHub{
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Properties: &virtualwans.VirtualHubProperties{
			Sku:                        pointer.To(d.Get("sku").(string)),
			AllowBranchToBranchTraffic: pointer.To(d.Get("branch_to_branch_traffic_enabled").(bool)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if err := client.VirtualHubsCreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	timeout, _ := ctx.Deadline()
	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"Provisioning", "Updating"},
		Target:                    []string{"Succeeded", "Provisioned"},
		Refresh:                   routeServerCreateRefreshFunc(ctx, client, id),
		PollInterval:              15 * time.Second,
		ContinuousTargetOccurence: 5,
		Timeout:                   time.Until(timeout),
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("waiting for creation of %s: %+v", id, err)
	}

	ipConfigName := "ipConfig1"
	ipConfigs := virtualwans.HubIPConfiguration{
		Name: pointer.To(ipConfigName),
		Properties: &virtualwans.HubIPConfigurationPropertiesFormat{
			PublicIPAddress: &virtualwans.PublicIPAddress{
				Id: pointer.To(d.Get("public_ip_address_id").(string)),
			},
			Subnet: &virtualwans.Subnet{
				Id: pointer.To(d.Get("subnet_id").(string)),
			},
		},
	}

	ipConfigId := commonids.NewVirtualHubIPConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.VirtualHubName, ipConfigName)

	if err := client.VirtualHubIPConfigurationCreateOrUpdateThenPoll(ctx, ipConfigId, ipConfigs); err != nil {
		return fmt.Errorf("creating %s: %+v", ipConfigId, err)
	}

	d.SetId(id.ID())

	return resourceRouteServerRead(d, meta)
}

func resourceRouteServerUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualwans.ParseVirtualHubID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.VirtualHubName, "azurerm_route_server")
	defer locks.UnlockByName(id.VirtualHubName, "azurerm_route_server")

	existing, err := client.VirtualHubsGet(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)

	}

	payload := existing.Model

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", *id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", *id)
	}

	if d.HasChange("branch_to_branch_traffic_enabled") {
		payload.Properties.AllowBranchToBranchTraffic = pointer.To(d.Get("branch_to_branch_traffic_enabled").(bool))
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if err := client.VirtualHubsCreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	timeout, _ := ctx.Deadline()
	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"Provisioning", "Updating"},
		Target:                    []string{"Succeeded", "Provisioned"},
		Refresh:                   routeServerCreateRefreshFunc(ctx, client, *id),
		PollInterval:              15 * time.Second,
		ContinuousTargetOccurence: 5,
		Timeout:                   time.Until(timeout),
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("waiting for update of %s: %+v", *id, err)
	}

	d.SetId(id.ID())

	return resourceRouteServerRead(d, meta)
}

func resourceRouteServerRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualwans.ParseVirtualHubID(d.Id())
	if err != nil {
		return err
	}

	routeServer, err := client.VirtualHubsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(routeServer.HttpResponse) {
			log.Printf("[INFO] %s does not exists - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.VirtualHubName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := routeServer.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		if props := model.Properties; props != nil {
			d.Set("sku", props.Sku)
			var virtualRouterIps *[]string
			if props.VirtualRouterIPs != nil {
				virtualRouterIps = props.VirtualRouterIPs
			}
			d.Set("virtual_router_ips", virtualRouterIps)
			if props.AllowBranchToBranchTraffic != nil {
				d.Set("branch_to_branch_traffic_enabled", props.AllowBranchToBranchTraffic)
			}
			if props.VirtualRouterAsn != nil {
				d.Set("virtual_router_asn", props.VirtualRouterAsn)
			}
			d.Set("routing_state", string(pointer.From(props.RoutingState)))
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("flattening `tags`: %+v", err)
		}
	}

	ipConfig, err := client.VirtualHubIPConfigurationList(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving IP Config for %s: %+v", id, err)
	}

	if model := ipConfig.Model; model != nil {
		for _, config := range *model {
			if props := config.Properties; props != nil {
				if props.PublicIPAddress != nil {
					d.Set("public_ip_address_id", pointer.From(props.PublicIPAddress.Id))
				}
				if props.Subnet != nil {
					d.Set("subnet_id", pointer.From(props.Subnet.Id))
				}
			}
		}
	}

	return nil
}

func resourceRouteServerDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualwans.ParseVirtualHubID(d.Id())
	if err != nil {
		return err
	}

	ipConfig, err := client.VirtualHubIPConfigurationList(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving IP Config for %s: %+v", id, err)
	}

	if ipConfig.Model != nil {
		var ipName string
		for _, config := range *ipConfig.Model {
			if config.Name != nil {
				ipName = *config.Name
			}
		}

		ipConfigId := commonids.NewVirtualHubIPConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.VirtualHubName, ipName)

		if err := deleteRouteServerIpConfiguration(ctx, client, ipConfigId); err != nil {
			return err
		}
	}

	if err := client.VirtualHubsDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func deleteRouteServerIpConfiguration(ctx context.Context, client *virtualwans.VirtualWANsClient, id commonids.VirtualHubIPConfigurationId) error {
	if err := client.VirtualHubIPConfigurationDeleteThenPoll(ctx, id); err != nil {
		return fmt.Errorf("deleting IP Config %s: %+v", id, err)
	}

	timeout, _ := ctx.Deadline()
	stateConf := &pluginsdk.StateChangeConf{
		Pending: []string{"200"},
		Target:  []string{"404"},
		Refresh: ipConfigStateRefreshFunc(ctx, client, id),
		Timeout: time.Until(timeout),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s: %+v", id, err)
	}
	return nil
}

func routeServerCreateRefreshFunc(ctx context.Context, client *virtualwans.VirtualWANsClient, id virtualwans.VirtualHubId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.VirtualHubsGet(ctx, id)
		if err != nil {
			if response.WasNotFound(resp.HttpResponse) {
				return nil, "", fmt.Errorf("%s does not exists", id)
			}
			return nil, "", fmt.Errorf("retrieving %s: %+v", id, err)
		}

		if resp.Model != nil && resp.Model.Properties != nil {
			return resp, string(pointer.From(resp.Model.Properties.ProvisioningState)), nil
		}
		return nil, "", fmt.Errorf("unable to read the provisioning state of this %s", id)
	}
}

func ipConfigStateRefreshFunc(ctx context.Context, client *virtualwans.VirtualWANsClient, id commonids.VirtualHubIPConfigurationId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.VirtualHubIPConfigurationGet(ctx, id)
		if err != nil {
			if response.WasNotFound(resp.HttpResponse) {
				return resp, strconv.Itoa(resp.HttpResponse.StatusCode), nil
			}
			return nil, "", fmt.Errorf("polling for the status of %s: %+v", id, err)
		}
		return resp, strconv.Itoa(resp.HttpResponse.StatusCode), nil
	}
}
