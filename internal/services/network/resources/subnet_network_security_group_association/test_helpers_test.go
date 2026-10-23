// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package subnet_network_security_group_association_test

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/subnets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type SubnetResource struct{}

func (SubnetResource) hasNoNetworkSecurityGroup(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) error {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Minute)
	defer cancel()

	id, err := commonids.ParseSubnetID(state.ID)
	if err != nil {
		return err
	}

	resp, err := client.Network.Subnets.Get(ctx, *id, subnets.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s does not exist", id)
		}

		return fmt.Errorf("bad: Get on subnetClient: %+v", err)
	}

	model := resp.Model
	if model == nil {
		return fmt.Errorf("model was nil for %s", id)
	}

	props := resp.Model.Properties
	if props == nil {
		return fmt.Errorf("properties was nil for %s", id)
	}

	if props.NetworkSecurityGroup != nil && ((props.NetworkSecurityGroup.Id == nil) || (props.NetworkSecurityGroup.Id != nil && *props.NetworkSecurityGroup.Id == "")) {
		return fmt.Errorf("no Network Security Group should exist for %s but got %q", id, *props.NetworkSecurityGroup.Id)
	}

	return nil
}
