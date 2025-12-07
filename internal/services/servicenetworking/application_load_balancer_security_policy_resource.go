// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package servicenetworking

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/webapplicationfirewallpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2025-01-01/securitypoliciesinterface"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.ResourceWithUpdate = SecurityPoliciesResource{}

type SecurityPoliciesResource struct{}

type SecurityPoliciesModel struct {
	ApplicationLoadBalancerId      string            `tfschema:"application_load_balancer_id"`
	Location                       string            `tfschema:"location"`
	Name                           string            `tfschema:"name"`
	WebApplicationFirewallPolicyId string            `tfschema:"web_application_firewall_policy_id"`
	Tags                           map[string]string `tfschema:"tags"`
}

func (f SecurityPoliciesResource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9_.-]{0,62}[a-zA-Z0-9])?$`),
				"`name` must begin with a letter or number, end with a letter and number, and must be between 1 and 64 characters long and can only contain letters, numbers, underscores(_), periods(.), and hyphens(-).",
			),
		},

		"application_load_balancer_id": commonschema.ResourceIDReferenceRequiredForceNew(&securitypoliciesinterface.TrafficControllerId{}),

		"location": commonschema.Location(),

		"web_application_firewall_policy_id": commonschema.ResourceIDReferenceRequiredForceNew(&webapplicationfirewallpolicies.ApplicationGatewayWebApplicationFirewallPolicyId{}),

		"tags": commonschema.Tags(),
	}
}

func (f SecurityPoliciesResource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{}
}

func (f SecurityPoliciesResource) ModelObject() interface{} {
	return &SecurityPoliciesModel{}
}

func (f SecurityPoliciesResource) ResourceType() string {
	return "azurerm_application_load_balancer_security_policy"
}

func (f SecurityPoliciesResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return securitypoliciesinterface.ValidateSecurityPolicyID
}

func (f SecurityPoliciesResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ServiceNetworking.SecurityPoliciesInterface

			var config SecurityPoliciesModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			trafficControllerId, err := securitypoliciesinterface.ParseTrafficControllerID(config.ApplicationLoadBalancerId)
			if err != nil {
				return err
			}

			id := securitypoliciesinterface.NewSecurityPolicyID(trafficControllerId.SubscriptionId, trafficControllerId.ResourceGroupName, trafficControllerId.TrafficControllerName, config.Name)

			resp, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(resp.HttpResponse) {
				return metadata.ResourceRequiresImport(f.ResourceType(), id)
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

			if err := client.CreateOrUpdateThenPoll(ctx, id, securityPolicy); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (f SecurityPoliciesResource) Read() sdk.ResourceFunc {
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

func (t SecurityPoliciesResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ServiceNetworking.SecurityPoliciesInterface

			id, err := securitypoliciesinterface.ParseSecurityPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config SecurityPoliciesModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			payload := securitypoliciesinterface.SecurityPolicyUpdate{}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = pointer.To(config.Tags)
			}

			if _, err := client.Update(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (f SecurityPoliciesResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ServiceNetworking.SecurityPoliciesInterface

			id, err := securitypoliciesinterface.ParseSecurityPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err = client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
