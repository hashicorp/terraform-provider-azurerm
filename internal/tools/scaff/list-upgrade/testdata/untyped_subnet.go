// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package testdata

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/serviceendpointpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/ipampools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/subnets"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

// resourceSubnet mirrors the real azurerm_subnet: it imports several go-azure-sdk
// resource-manager packages (subnets + serviceendpointpolicies + ipampools),
// builds its commonids.SubnetId directly from config strings (so no parent ID is
// visible in the source), and its Get call is the bare client.Get with a
// subnets-qualified options argument.
func resourceSubnet() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSubnetCreate,
		Read:   resourceSubnetRead,
		Update: resourceSubnetUpdate,
		Delete: resourceSubnetDelete,

		Importer: pluginsdk.ImporterValidatingIdentity(&commonids.SubnetId{}),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&commonids.SubnetId{}),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_group_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},
			"virtual_network_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceSubnetCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.Subnets
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewSubnetID(subscriptionId, d.Get("resource_group_name").(string), d.Get("virtual_network_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id, subnets.DefaultGetOperationOptions())
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for %s: %+v", id, err)
		}
	}

	// Unrelated SDK packages are imported too; the analyzer must still select
	// subnets as the model+client package.
	_ = serviceendpointpolicies.ServiceEndpointPolicyId{}
	_ = ipampools.IPamPoolId{}

	if err := client.CreateOrUpdateThenPoll(ctx, id, subnets.Subnet{}); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSubnetRead(d, meta)
}

func resourceSubnetUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.Subnets
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseSubnetID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Get(ctx, *id, subnets.DefaultGetOperationOptions()); err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return resourceSubnetRead(d, meta)
}

func resourceSubnetRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.Subnets
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseSubnetID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, subnets.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return resourceSubnetFlatten(d, *id, resp.Model)
}

func resourceSubnetDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.Subnets
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseSubnetID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

// resourceSubnetFlatten takes its id by value (commonids.SubnetId), so the list
// generator must dereference the parsed pointer id when calling it.
func resourceSubnetFlatten(d *pluginsdk.ResourceData, id commonids.SubnetId, subnet *subnets.Subnet) error {
	d.Set("name", id.SubnetName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("virtual_network_name", id.VirtualNetworkName)

	if subnet != nil {
		_ = subnet.Name
	}

	return pluginsdk.SetResourceIdentityData(d, &id)
}
