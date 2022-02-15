package cdn

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01/profiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01/securitypolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceFrontdoorSecurityPolicy() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceFrontdoorSecurityPolicyCreate,
		Read:   resourceFrontdoorSecurityPolicyRead,
		Update: resourceFrontdoorSecurityPolicyUpdate,
		Delete: resourceFrontdoorSecurityPolicyDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := securitypolicies.ParseSecurityPoliciesID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"frontdoor_profile_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: profiles.ValidateProfileID,
			},

			"web_application_firewall": {
				Type:     pluginsdk.TypeList,
				Required: true,
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

			"frontdoor_profile_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceFrontdoorSecurityPolicyCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontdoorSecurityPoliciesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	profileId, err := profiles.ParseProfileID(d.Get("frontdoor_profile_id").(string))
	if err != nil {
		return err
	}

	securityPolicyName := d.Get("name").(string)
	sdkId := securitypolicies.NewSecurityPoliciesID(profileId.SubscriptionId, profileId.ResourceGroupName, profileId.ProfileName, securityPolicyName)
	id := parse.NewFrontdoorSecurityPolicyID(profileId.SubscriptionId, profileId.ResourceGroupName, profileId.ProfileName, securityPolicyName)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, sdkId)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_frontdoor_security_policy", id.ID())
		}
	}

	params := securitypolicies.SecurityPolicyPropertiesParameters(nil)
	if waf, ok := d.GetOk("web_application_firewall"); ok {
		params = expandFrontdoorSecurityPoliciesParameters(waf.([]interface{}), securitypolicies.SecurityPolicyPropertiesParameters(&securitypolicies.SecurityPolicyWebApplicationFirewallParameters{}))
	} else {
		// Will look for DDoS policy here once it is GA
		return fmt.Errorf("unable to locate %q policy parameters", "web_application_firewall")
	}

	props := securitypolicies.SecurityPolicy{
		Properties: &securitypolicies.SecurityPolicyProperties{
			Parameters: params,
		},
	}
	if err := client.CreateThenPoll(ctx, sdkId, props); err != nil {

		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceFrontdoorSecurityPolicyRead(d, meta)
}

func resourceFrontdoorSecurityPolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontdoorSecurityPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	sdkId, err := securitypolicies.ParseSecurityPoliciesID(d.Id())
	if err != nil {
		return err
	}

	id, err := parse.FrontdoorSecurityPolicyID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *sdkId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.SecurityPolicyName)
	d.Set("frontdoor_profile_id", parse.NewFrontdoorProfileID(id.SubscriptionId, id.ResourceGroup, id.ProfileName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			// If it is not Type WebApplicationFirewall ignore this for now
			switch params := props.Parameters.(type) {
			case securitypolicies.SecurityPolicyWebApplicationFirewallParameters:
				if err := d.Set("web_application_firewall", flattenFrontdoorSecurityPoliciesWebApplicationFirewallParameters(&params)); err != nil {
					return fmt.Errorf("setting `web_application_firewall`: %+v", err)
				}
			default:
				// Unknown Security Policy Type
				return fmt.Errorf("unknown security policy type defined in the %q field: %+v", "parameters", err)
			}

			d.Set("profile_name", props.ProfileName)
		}
	}
	return nil
}

func resourceFrontdoorSecurityPolicyUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontdoorSecurityPoliciesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := securitypolicies.ParseSecurityPoliciesID(d.Id())
	if err != nil {
		return err
	}

	params := securitypolicies.SecurityPolicyUpdateProperties{}
	if waf, ok := d.GetOk("web_application_firewall"); ok {
		params = expandFrontdoorSecurityPoliciesUpdateWebApplicationFirewallParameters(waf.([]interface{}))
	} else {
		// Will look for DDoS policy here once it is GA
		return fmt.Errorf("unable to locate %q update parameters", "web_application_firewall")
	}

	props := securitypolicies.SecurityPolicyUpdateParameters{
		Properties: &params,
	}

	if err := client.PatchThenPoll(ctx, *id, props); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceFrontdoorSecurityPolicyRead(d, meta)
}

func resourceFrontdoorSecurityPolicyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontdoorSecurityPoliciesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := securitypolicies.ParseSecurityPoliciesID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {

		return fmt.Errorf("deleting %s: %+v", id, err)
	}
	return nil
}

func expandFrontdoorSecurityPoliciesParameters(input []interface{}, policyType securitypolicies.SecurityPolicyPropertiesParameters) securitypolicies.SecurityPolicyPropertiesParameters {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	results := securitypolicies.SecurityPolicyPropertiesParameters(nil)

	// TODO: Add DDoS when it GA's
	switch policyType.(type) {
	case securitypolicies.SecurityPolicyWebApplicationFirewallParameters:
		results = expandFrontdoorSecurityPoliciesWebApplicationFirewall(input)
	default:
		// Unknown Security Policy Type
		return nil
	}

	return results
}

func expandFrontdoorSecurityPoliciesWebApplicationFirewall(input []interface{}) securitypolicies.SecurityPolicyWebApplicationFirewallParameters {
	results := securitypolicies.SecurityPolicyWebApplicationFirewallParameters{}
	associations := make([]securitypolicies.SecurityPolicyWebApplicationFirewallAssociation, 0)
	v := input[0].(map[string]interface{})

	if id := v["waf_policy_id"].(string); id != "" {
		results.WafPolicy = &securitypolicies.ResourceReference{
			Id: utils.String(id),
		}
	}

	configAssociations := v["association"].([]interface{})

	for _, item := range configAssociations {
		v := item.(map[string]interface{})

		association := securitypolicies.SecurityPolicyWebApplicationFirewallAssociation{
			Domains:         expandSecurityPoliciesActivatedResourceReference(v["domain"].([]interface{})),
			PatternsToMatch: utils.ExpandStringSlice(v["patterns_to_match"].([]interface{})),
		}

		associations = append(associations, association)
	}

	results.Associations = &associations

	return results
}

func expandSecurityPoliciesActivatedResourceReference(input []interface{}) *[]securitypolicies.ActivatedResourceReference {
	results := make([]securitypolicies.ActivatedResourceReference, 0)
	if len(input) == 0 {
		return &results
	}

	for _, item := range input {
		v := item.(map[string]interface{})
		activatedResourceReference := securitypolicies.ActivatedResourceReference{}

		if id := v["id"].(string); id != "" {
			activatedResourceReference.Id = utils.String(id)

			enabled := v["enabled"].(bool)

			if !enabled {
				activatedResourceReference.IsActive = utils.Bool(enabled)
			}

			results = append(results, activatedResourceReference)
		}
	}

	return &results
}

func expandFrontdoorSecurityPoliciesUpdateWebApplicationFirewallParameters(input []interface{}) securitypolicies.SecurityPolicyUpdateProperties {
	if len(input) == 0 || input[0] == nil {
		return securitypolicies.SecurityPolicyUpdateProperties{}
	}

	waf := expandFrontdoorSecurityPoliciesParameters(input, securitypolicies.SecurityPolicyPropertiesParameters(&securitypolicies.SecurityPolicyWebApplicationFirewallParameters{}))

	return securitypolicies.SecurityPolicyUpdateProperties{
		Parameters: waf,
	}
}

func flattenFrontdoorSecurityPoliciesWebApplicationFirewallParameters(input *securitypolicies.SecurityPolicyWebApplicationFirewallParameters) []interface{} {
	params := make([]interface{}, 0)
	associations := make([]interface{}, 0)
	if input == nil {
		return params
	}

	// if we are here we know that the input is a SecurityPolicyWebApplicationFirewallParameters type
	values := make(map[string]interface{})
	values["waf_policy_id"] = *input.WafPolicy.Id

	for _, v := range *input.Associations {
		temp := make(map[string]interface{})
		domains := make([]interface{}, 0)

		for _, x := range *v.Domains {
			domain := make(map[string]interface{})
			if x.Id != nil {
				domain["id"] = *x.Id

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
