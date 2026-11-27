// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package application_load_balancer_security_policy

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/webapplicationfirewallpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2025-01-01/securitypoliciesinterface"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (f Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ServiceNetworking.SecurityPoliciesInterface

			id, err := securitypoliciesinterface.ParseSecurityPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			trafficControllerId := securitypoliciesinterface.NewTrafficControllerID(id.SubscriptionId, id.ResourceGroupName, id.TrafficControllerName)

			state := SecurityPoliciesModel{
				Name:                      id.SecurityPolicyName,
				ApplicationLoadBalancerId: trafficControllerId.ID(),
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Tags = pointer.From(model.Tags)

				if prop := model.Properties; prop != nil {
					if wafPolicy := prop.WafPolicy; wafPolicy != nil {
						webApplicationFirewallPolicyId, err := webapplicationfirewallpolicies.ParseApplicationGatewayWebApplicationFirewallPolicyIDInsensitively(wafPolicy.Id)
						if err != nil {
							return err
						}

						state.WebApplicationFirewallPolicyId = webApplicationFirewallPolicyId.ID()
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}
