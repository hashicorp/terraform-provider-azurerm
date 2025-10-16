// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-02-01/profiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-02-01/securitypolicies"
	waf "github.com/hashicorp/go-azure-sdk/resource-manager/frontdoor/2024-02-01/webapplicationfirewallpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceCdnFrontDoorSecurityPolicy() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCdnFrontdoorSecurityPolicyCreate,
		Read:   resourceCdnFrontdoorSecurityPolicyRead,
		Update: resourceCdnFrontdoorSecurityPolicyUpdate,
		Delete: resourceCdnFrontdoorSecurityPolicyDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := securitypolicies.ParseSecurityPolicyID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"cdn_frontdoor_profile_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FrontDoorProfileID,
			},

			"security_policies": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,

				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"firewall": {
							Type:     pluginsdk.TypeList,
							Required: true,
							MaxItems: 1,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"cdn_frontdoor_firewall_policy_id": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ForceNew:     true,
										ValidateFunc: validate.FrontDoorFirewallPolicyID,
									},

									"association": {
										Type:     pluginsdk.TypeList,
										Required: true,
										MaxItems: 1,

										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												// NOTE: The max number of domains vary depending on sku: 100 Standard, 500 Premium
												"domain": {
													Type:     pluginsdk.TypeList,
													Required: true,
													MaxItems: 500,

													Elem: &pluginsdk.Resource{
														Schema: map[string]*pluginsdk.Schema{
															"cdn_frontdoor_domain_id": {
																Type:         pluginsdk.TypeString,
																Required:     true,
																ValidateFunc: validate.FrontDoorSecurityPolicyDomainID,
															},

															"active": {
																Type:     pluginsdk.TypeBool,
																Computed: true,
															},
														},
													},
												},

												// NOTE: Per the service team the only acceptable value as of GA is "/*"
												"patterns_to_match": {
													Type:     pluginsdk.TypeList,
													Required: true,
													ForceNew: true,
													MaxItems: 1,

													Elem: &pluginsdk.Schema{
														Type: pluginsdk.TypeString,
														ValidateFunc: validation.StringInSlice([]string{
															"/*",
														}, false),
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceCdnFrontdoorSecurityPolicyCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorSecurityPoliciesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// NOTE: The profile id is used to retrieve properties from the related profile that must match in this security policy
	profileId, err := profiles.ParseProfileID(d.Get("cdn_frontdoor_profile_id").(string))
	if err != nil {
		return err
	}

	securityPolicyName := d.Get("name").(string)
	id := securitypolicies.NewSecurityPolicyID(profileId.SubscriptionId, profileId.ResourceGroupName, profileId.ProfileName, securityPolicyName)

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_cdn_frontdoor_security_policy", id.ID())
	}

	profileClient := meta.(*clients.Client).Cdn.FrontDoorProfilesClient
	resp, err := profileClient.Get(ctx, pointer.From(profileId))
	if err != nil {
		return fmt.Errorf("unable to retrieve the 'sku_name' from the CDN FrontDoor Profile(Name: %q)': %+v", profileId.ProfileName, err)
	}

	profileModel := resp.Model

	if profileModel == nil {
		return errors.New("profileModel is 'nil'")
	}

	isStandardSku := true
	if profileModel.Sku.Name != nil {
		isStandardSku = strings.HasPrefix(strings.ToLower(pointer.FromEnum(profileModel.Sku.Name)), "standard")
	}

	params, err := expandCdnFrontdoorFirewallPolicyParameters(d.Get("security_policies").([]interface{}), isStandardSku)
	if err != nil {
		return fmt.Errorf("expanding 'security_policies': %+v", err)
	}

	props := securitypolicies.SecurityPolicy{
		Properties: &securitypolicies.SecurityPolicyProperties{
			Parameters: params,
		},
	}

	err = client.CreateThenPoll(ctx, id, props)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceCdnFrontdoorSecurityPolicyRead(d, meta)
}

func resourceCdnFrontdoorSecurityPolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorSecurityPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := securitypolicies.ParseSecurityPolicyID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, pointer.From(id))
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.SecurityPolicyName)
	d.Set("cdn_frontdoor_profile_id", profiles.NewProfileID(id.SubscriptionId, id.ResourceGroupName, id.ProfileName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			if props.Parameters.SecurityPolicyPropertiesParameters().Type != securitypolicies.SecurityPolicyTypeWebApplicationFirewall {
				return fmt.Errorf("'model.Properties.Parameters.Type' of %q is unexpected, want security policy 'Type' of 'WebApplicationFirewall': %s", props.Parameters.SecurityPolicyPropertiesParameters().Type, id)
			}

			// we know it's a firewall policy at this point,
			// create the objects to hold the policy data
			wafParams := props.Parameters.(securitypolicies.SecurityPolicyWebApplicationFirewallParameters)
			associations := make([]interface{}, 0)
			wafPolicyId := ""

			if wafParams.WafPolicy != nil && wafParams.WafPolicy.Id != nil {
				parsedId, err := waf.ParseFrontDoorWebApplicationFirewallPolicyIDInsensitively(*wafParams.WafPolicy.Id)
				if err != nil {
					return fmt.Errorf("flattening `cdn_frontdoor_firewall_policy_id`: %+v", err)
				}
				wafPolicyId = parsedId.ID()
			}

			if wafParams.Associations != nil {
				for _, item := range *wafParams.Associations {
					domain, err := flattenSecurityPoliciesActivatedResourceReference(item.Domains)
					if err != nil {
						return fmt.Errorf("flattening `ActivatedResourceReference`: %+v", err)
					}

					associations = append(associations, map[string]interface{}{
						"domain":            domain,
						"patterns_to_match": utils.FlattenStringSlice(item.PatternsToMatch),
					})
				}
			}

			securityPolicy := []interface{}{
				map[string]interface{}{
					"firewall": []interface{}{
						map[string]interface{}{
							"association":                      associations,
							"cdn_frontdoor_firewall_policy_id": wafPolicyId,
						},
					},
				},
			}

			d.Set("security_policies", securityPolicy)
		}
	}

	return nil
}

func resourceCdnFrontdoorSecurityPolicyUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorSecurityPoliciesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// NOTE: The profile id is used to retrieve properties from the related profile that must match in this security policy
	profileId, err := profiles.ParseProfileID(d.Get("cdn_frontdoor_profile_id").(string))
	if err != nil {
		return err
	}

	securityPolicyName := d.Get("name").(string)
	id := securitypolicies.NewSecurityPolicyID(profileId.SubscriptionId, profileId.ResourceGroupName, profileId.ProfileName, securityPolicyName)

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for existing %s: %+v", id, err)
		}
	}

	profileClient := meta.(*clients.Client).Cdn.FrontDoorProfilesClient
	resp, err := profileClient.Get(ctx, pointer.From(profileId))
	if err != nil {
		return fmt.Errorf("unable to retrieve the `sku_name` from %s: %+v", *profileId, err)
	}

	profileModel := resp.Model

	if profileModel == nil {
		return errors.New("profileModel is 'nil'")
	}

	isStandardSku := true
	if profileModel.Sku.Name != nil {
		isStandardSku = strings.HasPrefix(strings.ToLower(pointer.FromEnum(profileModel.Sku.Name)), "standard")
	}

	params, err := expandCdnFrontdoorFirewallPolicyParameters(d.Get("security_policies").([]interface{}), isStandardSku)
	if err != nil {
		return fmt.Errorf("expanding 'security_policies': %+v", err)
	}

	props := securitypolicies.SecurityPolicy{
		Properties: &securitypolicies.SecurityPolicyProperties{
			Parameters: params,
		},
	}

	// Using 'Create' for update because it is a PUT operation
	err = client.CreateThenPoll(ctx, id, props)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceCdnFrontdoorSecurityPolicyRead(d, meta)
}

func resourceCdnFrontdoorSecurityPolicyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorSecurityPoliciesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := securitypolicies.ParseSecurityPolicyID(d.Id())
	if err != nil {
		return err
	}

	err = client.DeleteThenPoll(ctx, pointer.From(id))
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandCdnFrontdoorFirewallPolicyParameters(input []interface{}, isStandardSku bool) (*securitypolicies.SecurityPolicyWebApplicationFirewallParameters, error) {
	results := securitypolicies.SecurityPolicyWebApplicationFirewallParameters{}
	if len(input) == 0 {
		return &results, nil
	}

	associations := make([]securitypolicies.SecurityPolicyWebApplicationFirewallAssociation, 0)

	// pull off only the firewall policy from the security_policies list
	policyType := input[0].(map[string]interface{})
	firewallPolicy := policyType["firewall"].([]interface{})
	v := firewallPolicy[0].(map[string]interface{})

	if id := v["cdn_frontdoor_firewall_policy_id"].(string); id != "" {
		results.WafPolicy = &securitypolicies.ResourceReference{
			Id: pointer.To(id),
		}
	}

	configAssociations := v["association"].([]interface{})

	for _, item := range configAssociations {
		v := item.(map[string]interface{})
		domains := expandSecurityPoliciesActivatedResourceReference(v["domain"].([]interface{}))

		if isStandardSku {
			if len(*domains) > 100 {
				return &results, fmt.Errorf("the 'Standard_AzureFrontDoor' sku is only allowed to have 100 or less domains associated with the firewall policy, got %d", len(*domains))
			}
		} else {
			if len(*domains) > 500 {
				return &results, fmt.Errorf("the 'Premium_AzureFrontDoor' sku is only allowed to have 500 or less domains associated with the firewall policy, got %d", len(*domains))
			}
		}

		association := securitypolicies.SecurityPolicyWebApplicationFirewallAssociation{
			Domains:         domains,
			PatternsToMatch: utils.ExpandStringSlice(v["patterns_to_match"].([]interface{})),
		}

		associations = append(associations, association)
	}

	results.Associations = &associations

	return &results, nil
}

func expandSecurityPoliciesActivatedResourceReference(input []interface{}) *[]securitypolicies.ActivatedResourceReference {
	results := make([]securitypolicies.ActivatedResourceReference, 0)
	if len(input) == 0 {
		return &results
	}

	for _, item := range input {
		v := item.(map[string]interface{})

		if id := v["cdn_frontdoor_domain_id"].(string); id != "" {
			results = append(results, securitypolicies.ActivatedResourceReference{
				Id: pointer.To(id),
			})
		}
	}

	return &results
}

func flattenSecurityPoliciesActivatedResourceReference(input *[]securitypolicies.ActivatedResourceReference) ([]interface{}, error) {
	results := make([]interface{}, 0)
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
				return nil, fmt.Errorf("flattening `cdn_frontdoor_domain_id`: %+v; %+v", frontDoorCustomDomainIdErr, frontDoorEndpointIdErr)
			}
		}

		results = append(results, map[string]interface{}{
			"active":                  pointer.From(item.IsActive),
			"cdn_frontdoor_domain_id": frontDoorDomainId,
		})
	}

	return results, nil
}
