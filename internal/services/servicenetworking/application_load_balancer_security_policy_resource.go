package servicenetworking

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-05-01/webapplicationfirewallpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2025-01-01/securitypoliciesinterface"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.Resource = SecurityPoliciesResource{}

type SecurityPoliciesResource struct{}

type SecurityPoliciesModel struct {
	Name                           string `tfschema:"name"`
	ApplicationLoadBalancerId      string `tfschema:"application_load_balancer_id"`
	WebApplicationFirewallPolicyId string `tfschema:"web_application_firewall_policy_id"`
}

func (f SecurityPoliciesResource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9_.-]{0,62}[a-zA-Z0-9])?$`), "`name` must begin with a letter or number, end with a letter and number, and must be between 1 and 64 characters long and can only contain letters, numbers, underscores(_), periods(.), and hyphens(-)."),
		},

		"application_load_balancer_id": commonschema.ResourceIDReferenceRequiredForceNew(&securitypoliciesinterface.TrafficControllerId{}),

		"web_application_firewall_policy_id": commonschema.ResourceIDReferenceRequiredForceNew(&webapplicationfirewallpolicies.ApplicationGatewayWebApplicationFirewallPolicyId{}),
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
				if prop := model.Properties; prop != nil {
					if wafPolicy := prop.WafPolicy; wafPolicy != nil {
						state.WebApplicationFirewallPolicyId = wafPolicy.Id
					}
				}
			}

			return metadata.Encode(&state)
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
				return fmt.Errorf("deleting %q: %+v", id.ID(), err)
			}

			return nil
		},
	}
}
