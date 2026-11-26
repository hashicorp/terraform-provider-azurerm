// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"context"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func FindLoadBalancerFrontEndIpConfigurationByName(lb *loadbalancers.LoadBalancer, name string) (*loadbalancers.FrontendIPConfiguration, bool) {
	if lb == nil || lb.Properties == nil || lb.Properties.FrontendIPConfigurations == nil {
		return nil, false
	}

	for _, feip := range *lb.Properties.FrontendIPConfigurations {
		if feip.Name != nil && *feip.Name == name {
			return &feip, true
		}
	}

	return nil, false
}

func LoadBalancerSubResourceImporter(parser func(input string) (*loadbalancers.LoadBalancerId, error)) *schema.ResourceImporter {
	return pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
		_, err := parser(id)
		return err
	}, func(ctx context.Context, d *pluginsdk.ResourceData, meta any) ([]*pluginsdk.ResourceData, error) {
		lbId, err := parser(d.Id())
		if err != nil {
			return nil, err
		}

		d.Set("loadbalancer_id", lbId.ID())
		return []*pluginsdk.ResourceData{d}, nil
	})
}
