// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2021-06-01/cdn" // nolint: staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	cdnfrontdoorsecurityparams "github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/frontdoorsecurityparams"
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
		Delete: resourceCdnFrontdoorSecurityPolicyDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.FrontDoorSecurityPolicyID(id)
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
				ForceNew: true,
				MaxItems: 1,

				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{

						"firewall": {
							Type:     pluginsdk.TypeList,
							Required: true,
							ForceNew: true,
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
										ForceNew: true,
										MaxItems: 1,

										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{

												// NOTE: The max number of domains vary depending on sku: 100 Standard, 500 Premium
												"domain": {
													Type:     pluginsdk.TypeList,
													Required: true,
													ForceNew: true,
													MaxItems: 500,

													Elem: &pluginsdk.Resource{
														Schema: map[string]*pluginsdk.Schema{
															"cdn_frontdoor_domain_id": {
																Type:         pluginsdk.TypeString,
																Required:     true,
																ForceNew:     true,
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
	profile, err := parse.FrontDoorProfileID(d.Get("cdn_frontdoor_profile_id").(string))
	if err != nil {
		return err
	}

	securityPolicyName := d.Get("name").(string)
	id := parse.NewFrontDoorSecurityPolicyID(profile.SubscriptionId, profile.ResourceGroup, profile.ProfileName, securityPolicyName)

	existing, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.SecurityPolicyName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for existing %s: %+v", id, err)
		}
	}

	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_cdn_frontdoor_security_policy", id.ID())
	}

	profileClient := meta.(*clients.Client).Cdn.FrontDoorProfileClient
	resp, err := profileClient.Get(ctx, profile.ResourceGroup, profile.ProfileName)
	if err != nil {
		return fmt.Errorf("unable to retrieve the 'sku_name' from the CDN FrontDoor Profile(Name: %q)': %+v", profile.ProfileName, err)
	}

	if resp.Sku == nil {
		return fmt.Errorf("the CDN FrontDoor Profile(Name: %q) 'sku' was nil", profile.ProfileName)
	}

	isStandardSku := strings.HasPrefix(strings.ToLower(string(resp.Sku.Name)), "standard")

	params, err := cdnfrontdoorsecurityparams.ExpandCdnFrontdoorFirewallPolicyParameters(d.Get("security_policies").([]interface{}), isStandardSku)
	if err != nil {
		return fmt.Errorf("expanding 'security_policies': %+v", err)
	}

	props := cdn.SecurityPolicy{
		SecurityPolicyProperties: &cdn.SecurityPolicyProperties{
			Parameters: params,
		},
	}

	future, err := client.Create(ctx, id.ResourceGroup, id.ProfileName, id.SecurityPolicyName, props)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the creation of %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceCdnFrontdoorSecurityPolicyRead(d, meta)
}

func resourceCdnFrontdoorSecurityPolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorSecurityPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontDoorSecurityPolicyID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.SecurityPolicyName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.SecurityPolicyName)
	d.Set("cdn_frontdoor_profile_id", parse.NewFrontDoorProfileID(id.SubscriptionId, id.ResourceGroup, id.ProfileName).ID())

	if props := resp.SecurityPolicyProperties; props != nil {
		waf, ok := props.Parameters.AsSecurityPolicyWebApplicationFirewallParameters()
		if !ok {
			return fmt.Errorf("flattening %s: %s", id, "expected security policy web application firewall parameters")
		}

		// we know it's a firewall policy at this point,
		// create the objects to hold the policy data
		associations := make([]interface{}, 0)

		wafPolicyId := ""
		if waf.WafPolicy != nil && waf.WafPolicy.ID != nil {
			parsedId, err := parse.FrontDoorFirewallPolicyIDInsensitively(*waf.WafPolicy.ID)
			if err != nil {
				return fmt.Errorf("flattening `cdn_frontdoor_firewall_policy_id`: %+v", err)
			}
			wafPolicyId = parsedId.ID()
		}

		if waf.Associations != nil {
			for _, item := range *waf.Associations {
				domain, err := cdnfrontdoorsecurityparams.FlattenSecurityPoliciesActivatedResourceReference(item.Domains)
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

	return nil
}

func resourceCdnFrontdoorSecurityPolicyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorSecurityPoliciesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontDoorSecurityPolicyID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.ProfileName, id.SecurityPolicyName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the deletion of %s: %+v", *id, err)
	}

	return nil
}
