// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-02-01/profiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-02-01/securitypolicies"
	waf "github.com/hashicorp/go-azure-sdk/resource-manager/frontdoor/2024-02-01/webapplicationfirewallpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.DataSource = CdnFrontDoorSecurityPolicyDataSource{}

type CdnFrontDoorSecurityPolicyDataSource struct{}

type CdnFrontDoorSecurityPolicyDataSourceModel struct {
	Name                  string                                  `tfschema:"name"`
	ProfileName           string                                  `tfschema:"profile_name"`
	ResourceGroupName     string                                  `tfschema:"resource_group_name"`
	CdnFrontDoorProfileId string                                  `tfschema:"cdn_frontdoor_profile_id"`
	SecurityPolicies      []CdnFrontDoorSecurityPolicyPolicyModel `tfschema:"security_policies"`
}

type CdnFrontDoorSecurityPolicyPolicyModel struct {
	Firewall []CdnFrontDoorSecurityPolicyFirewallModel `tfschema:"firewall"`
}

type CdnFrontDoorSecurityPolicyFirewallModel struct {
	Association                  []CdnFrontDoorSecurityPolicyAssociationModel `tfschema:"association"`
	CdnFrontDoorFirewallPolicyId string                                       `tfschema:"cdn_frontdoor_firewall_policy_id"`
}

type CdnFrontDoorSecurityPolicyAssociationModel struct {
	Domain          []CdnFrontDoorSecurityPolicyDomainModel `tfschema:"domain"`
	PatternsToMatch []string                                `tfschema:"patterns_to_match"`
}

type CdnFrontDoorSecurityPolicyDomainModel struct {
	Active               bool   `tfschema:"active"`
	CdnFrontDoorDomainId string `tfschema:"cdn_frontdoor_domain_id"`
}

func (CdnFrontDoorSecurityPolicyDataSource) ResourceType() string {
	return "azurerm_cdn_frontdoor_security_policy"
}

func (CdnFrontDoorSecurityPolicyDataSource) ModelObject() interface{} {
	return &CdnFrontDoorSecurityPolicyDataSourceModel{}
}

func (CdnFrontDoorSecurityPolicyDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.FrontDoorSecurityPolicyName,
		},

		"profile_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.FrontDoorName,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (CdnFrontDoorSecurityPolicyDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"cdn_frontdoor_profile_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"security_policies": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"firewall": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"association": {
									Type:     pluginsdk.TypeList,
									Computed: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"domain": {
												Type:     pluginsdk.TypeList,
												Computed: true,
												Elem: &pluginsdk.Resource{
													Schema: map[string]*pluginsdk.Schema{
														"active": {
															Type:     pluginsdk.TypeBool,
															Computed: true,
														},
														"cdn_frontdoor_domain_id": {
															Type:     pluginsdk.TypeString,
															Computed: true,
														},
													},
												},
											},
											"patterns_to_match": {
												Type:     pluginsdk.TypeList,
												Computed: true,
												Elem: &pluginsdk.Schema{
													Type: pluginsdk.TypeString,
												},
											},
										},
									},
								},
								"cdn_frontdoor_firewall_policy_id": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
							},
						},
					},
				},
			},
		},
	}
}

func (CdnFrontDoorSecurityPolicyDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cdn.FrontDoorSecurityPoliciesClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state CdnFrontDoorSecurityPolicyDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := securitypolicies.NewSecurityPolicyID(subscriptionId, state.ResourceGroupName, state.ProfileName, state.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state.Name = id.SecurityPolicyName
			state.ProfileName = id.ProfileName
			state.ResourceGroupName = id.ResourceGroupName
			state.CdnFrontDoorProfileId = profiles.NewProfileID(id.SubscriptionId, id.ResourceGroupName, id.ProfileName).ID()

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					securityPolicies, err := flattenCdnFrontDoorSecurityPolicyDataSource(props.Parameters)
					if err != nil {
						return fmt.Errorf("flattening `security_policies`: %+v", err)
					}

					state.SecurityPolicies = securityPolicies
				}
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}

func flattenCdnFrontDoorSecurityPolicyDataSource(input securitypolicies.SecurityPolicyPropertiesParameters) ([]CdnFrontDoorSecurityPolicyPolicyModel, error) {
	results := make([]CdnFrontDoorSecurityPolicyPolicyModel, 0)
	if input.SecurityPolicyPropertiesParameters().Type != securitypolicies.SecurityPolicyTypeWebApplicationFirewall {
		return results, fmt.Errorf("unexpected security policy `Type` %q, expected `WebApplicationFirewall`", input.SecurityPolicyPropertiesParameters().Type)
	}

	wafParams := input.(securitypolicies.SecurityPolicyWebApplicationFirewallParameters)
	associations := make([]CdnFrontDoorSecurityPolicyAssociationModel, 0)
	wafPolicyId := ""

	if wafParams.WafPolicy != nil {
		parsedId, err := waf.ParseFrontDoorWebApplicationFirewallPolicyIDInsensitively(pointer.From(wafParams.WafPolicy.Id))
		if err != nil {
			return results, err
		}

		wafPolicyId = parsedId.ID()
	}

	if wafParams.Associations != nil {
		for _, item := range *wafParams.Associations {
			domain, err := flattenCdnFrontDoorSecurityPolicyActivatedResourceReference(item.Domains)
			if err != nil {
				return results, fmt.Errorf("flattening `domain`: %+v", err)
			}

			associations = append(associations, CdnFrontDoorSecurityPolicyAssociationModel{
				Domain:          domain,
				PatternsToMatch: pointer.From(item.PatternsToMatch),
			})
		}
	}

	results = []CdnFrontDoorSecurityPolicyPolicyModel{
		{
			Firewall: []CdnFrontDoorSecurityPolicyFirewallModel{
				{
					Association:                  associations,
					CdnFrontDoorFirewallPolicyId: wafPolicyId,
				},
			},
		},
	}

	return results, nil
}

func flattenCdnFrontDoorSecurityPolicyActivatedResourceReference(input *[]securitypolicies.ActivatedResourceReference) ([]CdnFrontDoorSecurityPolicyDomainModel, error) {
	results := make([]CdnFrontDoorSecurityPolicyDomainModel, 0)
	if input == nil {
		return results, nil
	}

	for _, item := range *input {
		frontDoorDomainId := ""
		if item.Id != nil {
			if parsedFrontDoorCustomDomainId, frontDoorCustomDomainIdErr := parse.FrontDoorCustomDomainIDInsensitively(*item.Id); frontDoorCustomDomainIdErr == nil {
				frontDoorDomainId = parsedFrontDoorCustomDomainId.ID()
			} else if parsedFrontDoorEndpointId, frontDoorEndpointIdErr := parse.FrontDoorEndpointIDInsensitively(*item.Id); frontDoorEndpointIdErr == nil {
				frontDoorDomainId = parsedFrontDoorEndpointId.ID()
			} else {
				return results, fmt.Errorf("flattening `cdn_frontdoor_domain_id`: %+v; %+v", frontDoorCustomDomainIdErr, frontDoorEndpointIdErr)
			}
		}

		results = append(results, CdnFrontDoorSecurityPolicyDomainModel{
			Active:               pointer.From(item.IsActive),
			CdnFrontDoorDomainId: frontDoorDomainId,
		})
	}

	return results, nil
}
