package cdn

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	track1 "github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceCdnFrontdoorSecurityPolicy() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCdnFrontdoorSecurityPolicyCreate,
		Read:   resourceCdnFrontdoorSecurityPolicyRead,
		Delete: resourceCdnFrontdoorSecurityPolicyDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.FrontdoorSecurityPolicyID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"cdn_frontdoor_profile_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FrontdoorProfileID,
			},

			"web_application_firewall": {
				Type:     pluginsdk.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,

				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"waf_policy_id": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},

						"association": {
							Type:     pluginsdk.TypeList,
							Required: true,
							MaxItems: 100,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{

									"domain": {
										Type:     pluginsdk.TypeList,
										Required: true,
										MaxItems: 25,

										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{

												"id": {
													Type:         pluginsdk.TypeString,
													Required:     true,
													ValidateFunc: validation.StringIsNotEmpty,
												},

												"enabled": {
													Type:     pluginsdk.TypeBool,
													Optional: true,
													Default:  true,
												},
											},
										},
									},

									"patterns_to_match": {
										Type:     pluginsdk.TypeList,
										Required: true,
										MaxItems: 25,

										Elem: &pluginsdk.Schema{
											Type:         pluginsdk.TypeString,
											ValidateFunc: validation.StringIsNotEmpty,
										},
									},
								},
							},
						},
					},
				},
			},

			"cdn_frontdoor_profile_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCdnFrontdoorSecurityPolicyCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontdoorSecurityPoliciesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	profileId, err := parse.FrontdoorProfileID(d.Get("cdn_frontdoor_profile_id").(string))
	if err != nil {
		return err
	}

	securityPolicyName := d.Get("name").(string)
	id := parse.NewFrontdoorSecurityPolicyID(profileId.SubscriptionId, profileId.ResourceGroup, profileId.ProfileName, securityPolicyName)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.SecurityPolicyName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_cdn_frontdoor_security_policy", id.ID())
		}
	}

	params := track1.BasicSecurityPolicyPropertiesParameters(nil)
	if waf, ok := d.GetOk("web_application_firewall"); ok {
		params = expandCdnFrontdoorSecurityPoliciesParameters(waf.([]interface{}), waf)
	} else {
		// Will look for DDoS policy here once it is GA
		return fmt.Errorf("unable to locate %q policy parameters", "web_application_firewall")
	}

	props := track1.SecurityPolicy{
		SecurityPolicyProperties: &track1.SecurityPolicyProperties{
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
	client := meta.(*clients.Client).Cdn.FrontdoorSecurityPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontdoorSecurityPolicyID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.SecurityPolicyName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on Azure Frontdoor Security Policy %q (Resource Group %q): %+v", id.SecurityPolicyName, id.ResourceGroup, err)
	}

	d.Set("name", id.SecurityPolicyName)
	d.Set("cdn_frontdoor_profile_id", parse.NewFrontdoorProfileID(id.SubscriptionId, id.ResourceGroup, id.ProfileName).ID())

	if props := resp.SecurityPolicyProperties; props != nil {
		// If it is not Type WebApplicationFirewall ignore this for now
		switch params := props.Parameters.(type) {
		case track1.SecurityPolicyWebApplicationFirewallParameters:
			if err := d.Set("web_application_firewall", flattenCdnFrontdoorSecurityPoliciesWebApplicationFirewallParameters(&params)); err != nil {
				return fmt.Errorf("setting `web_application_firewall`: %+v", err)
			}
		default:
			// Unknown Security Policy Type
			return fmt.Errorf("unknown security policy type defined in the %q field: %+v", "parameters", err)
		}

		d.Set("profile_name", props.ProfileName)
	}

	return nil
}

func resourceCdnFrontdoorSecurityPolicyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontdoorSecurityPoliciesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontdoorSecurityPolicyID(d.Id())
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

func expandCdnFrontdoorSecurityPoliciesParameters(input []interface{}, policyType interface{}) track1.SecurityPolicyPropertiesParameters {
	if len(input) == 0 || input[0] == nil {
		return track1.SecurityPolicyPropertiesParameters{}
	}

	results := track1.SecurityPolicyPropertiesParameters{}

	// TODO: Add DDoS when it GA's
	switch policyType.(type) {
	case track1.SecurityPolicyWebApplicationFirewallParameters:
		//results = expandFrontdoorSecurityPoliciesWebApplicationFirewall(input)
		results = track1.SecurityPolicyPropertiesParameters{
			Type: track1.TypeWebApplicationFirewall,
		}
	default:
		// Unknown Security Policy Type
		return results
	}

	return results
}

// func expandFrontdoorSecurityPoliciesWebApplicationFirewall(input []interface{}) track1.SecurityPolicyWebApplicationFirewallParameters {
// 	results := track1.SecurityPolicyWebApplicationFirewallParameters{}
// 	associations := make([]track1.SecurityPolicyWebApplicationFirewallAssociation, 0)
// 	v := input[0].(map[string]interface{})

// 	if id := v["waf_policy_id"].(string); id != "" {
// 		results.WafPolicy = &track1.ResourceReference{
// 			ID: utils.String(id),
// 		}
// 	}

// 	configAssociations := v["association"].([]interface{})

// 	for _, item := range configAssociations {
// 		v := item.(map[string]interface{})

// 		association := track1.SecurityPolicyWebApplicationFirewallAssociation{
// 			Domains:         expandSecurityPoliciesActivatedResourceReference(v["domain"].([]interface{})),
// 			PatternsToMatch: utils.ExpandStringSlice(v["patterns_to_match"].([]interface{})),
// 		}

// 		associations = append(associations, association)
// 	}

// 	results.Associations = &associations

// 	return results
// }

// func expandSecurityPoliciesActivatedResourceReference(input []interface{}) *[]track1.ActivatedResourceReference {
// 	results := make([]track1.ActivatedResourceReference, 0)
// 	if len(input) == 0 {
// 		return &results
// 	}

// 	for _, item := range input {
// 		v := item.(map[string]interface{})
// 		activatedResourceReference := track1.ActivatedResourceReference{}

// 		if id := v["id"].(string); id != "" {
// 			activatedResourceReference.ID = utils.String(id)

// 			enabled := v["enabled"].(bool)

// 			if !enabled {
// 				activatedResourceReference.IsActive = utils.Bool(enabled)
// 			}

// 			results = append(results, activatedResourceReference)
// 		}
// 	}

// 	return &results
// }

func expandCdnFrontdoorSecurityPoliciesUpdateWebApplicationFirewallParameters(input []interface{}) track1.SecurityPolicyUpdateProperties {
	if len(input) == 0 || input[0] == nil {
		return track1.SecurityPolicyUpdateProperties{}
	}

	// TODO: This isn't quite right...
	waf := expandCdnFrontdoorSecurityPoliciesParameters(input, track1.BasicSecurityPolicyPropertiesParameters(track1.SecurityPolicyWebApplicationFirewallParameters{}))

	return track1.SecurityPolicyUpdateProperties{
		Parameters: track1.BasicSecurityPolicyPropertiesParameters(waf),
	}
}

func flattenCdnFrontdoorSecurityPoliciesWebApplicationFirewallParameters(input *track1.SecurityPolicyWebApplicationFirewallParameters) []interface{} {
	params := make([]interface{}, 0)
	associations := make([]interface{}, 0)
	if input == nil {
		return params
	}

	// if we are here we know that the input is a SecurityPolicyWebApplicationFirewallParameters type
	values := make(map[string]interface{})
	values["waf_policy_id"] = *input.WafPolicy.ID

	for _, v := range *input.Associations {
		temp := make(map[string]interface{})
		domains := make([]interface{}, 0)

		for _, x := range *v.Domains {
			domain := make(map[string]interface{})
			if x.ID != nil {
				domain["id"] = *x.ID

				if x.IsActive != nil {
					domain["enabled"] = *x.IsActive
				}
				domains = append(domains, domain)
			}
		}

		association := make([]interface{}, 0)
		temp["domain"] = domains
		temp["patterns_to_match"] = *v.PatternsToMatch

		association = append(association, temp)
		associations = append(associations, association)
	}

	values["association"] = associations
	params = append(params, values)
	return params
}
