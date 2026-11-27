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
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2025-01-01/securitypoliciesinterface"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	servicenetworkingClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/servicenetworking/client"
)

func (f Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := servicenetworkingClient.From(metadata.Client).SecurityPoliciesInterface

			var config SecurityPoliciesModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			trafficControllerId, err := securitypoliciesinterface.ParseTrafficControllerID(config.ApplicationLoadBalancerId)
			if err != nil {
				return err
			}

			id := securitypoliciesinterface.NewSecurityPolicyID(trafficControllerId.SubscriptionId, trafficControllerId.ResourceGroupName, trafficControllerId.TrafficControllerName, config.Name)

			if !features.From(metadata.Client).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				resp, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				if !response.WasNotFound(resp.HttpResponse) {
					return metadata.ResourceRequiresImport(f.ResourceType(), id)
				}
			}

			securityPolicy := securitypoliciesinterface.SecurityPolicy{
				Location: location.Normalize(config.Location),
				Tags:     pointer.To(config.Tags),
				Properties: &securitypoliciesinterface.SecurityPolicyProperties{
					WafPolicy: &securitypoliciesinterface.WafPolicy{
						Id: config.WebApplicationFirewallPolicyId,
					},
				},
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, securityPolicy, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)

			return nil
		},
	}
}
