package cdn

import (
	"fmt"
	"strings"
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

			"parameters": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,

				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"type": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  "WebApplicationFirewall",
							ValidateFunc: validation.StringInSlice([]string{
								"WebApplicationFirewall",
							}, false),
						},

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

	props := securitypolicies.SecurityPolicy{
		Properties: &securitypolicies.SecurityPolicyProperties{
			Parameters: expandFrontdoorSecurityPoliciesParameters(d.Get("parameters").([]interface{})),
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
				if err := d.Set("parameters", flattenFrontdoorSecurityPoliciesParameters(&params)); err != nil {
					return fmt.Errorf("setting `parameters`: %+v", err)
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

	props := securitypolicies.SecurityPolicyUpdateParameters{
		Properties: expandFrontdoorSecurityPoliciesUpdateParameters(d.Get("parameters").([]interface{})),
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

func expandFrontdoorSecurityPoliciesParameters(input []interface{}) securitypolicies.SecurityPolicyPropertiesParameters {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	results := securitypolicies.SecurityPolicyWebApplicationFirewallParameters{}
	associations := make([]securitypolicies.SecurityPolicyWebApplicationFirewallAssociation, 0)
	v := input[0].(map[string]interface{})

	// DDoS is going to come later, currently WebApplicationFirewall is the
	// only supported security policy
	if secPolType := v["type"].(string); secPolType != "" {
		if strings.EqualFold(secPolType, "WebApplicationFirewall") {
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
		}
	}
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

func expandFrontdoorSecurityPoliciesUpdateParameters(input []interface{}) *securitypolicies.SecurityPolicyUpdateProperties {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})

	typeValue := securitypolicies.SecurityPolicyType(v["type"].(string))
	return &securitypolicies.SecurityPolicyUpdateProperties{
		Parameters: typeValue,
	}
}

// func flattenFrontdoorSecurityPoliciesParameters(input *securitypolicies.SecurityPolicyPropertiesParameters) []interface{} {
func flattenFrontdoorSecurityPoliciesParameters(input *securitypolicies.SecurityPolicyWebApplicationFirewallParameters) []interface{} {
	params := make([]interface{}, 0)
	if input == nil {
		return params
	}

	topLevel := make(map[string]interface{}, 0)
	// associations := make([]interface{}, 0)

	// if we are here we know that the input is a SecurityPolicyWebApplicationFirewallParameters type
	topLevel["type"] = *input.Associations
	topLevel["waf_policy_id"] = *input.WafPolicy.Id
	topLevel["association"] = *input.WafPolicy.Id

	// Parameters {
	// 	type = "WebApplicationFirewall"
	// 	waf_policy_id = "foo"

	// 	association {
	// 		domain {
	// 			id      = "foo"
	// 			enabled = true
	// 		}
	// 		domain {
	// 			id      = "foo"
	// 			enabled = true
	// 		}
	// 		domain {
	// 			id      = "foo"
	// 			enabled = true
	// 		}
	// 		patterns_to_match = ["foo"]
	// 	}

	// 	association {
	// 		domain {
	// 			id      = "foo"
	// 			enabled = true
	// 		}
	// 		domain {
	// 			id      = "badf00d"
	// 			enabled = false
	// 		}

	// 		patterns_to_match = ["foo"]
	// 	}
	// }

	// Net New
	// foo["type"] = "WebApplicationFirewall"
	// foo["waf_policy_id"] = waf

	// for _, v := range assocs {
	// 	association := make(map[string]interface{}, 0)
	// 	domains := make([]interface{}, 0)
	// 	domain := make(map[string]interface{}, 0)
	// 	for _, x := range *v.Domains {
	// 		domain["id"] = *x.Id
	// 		domain["enabled"] = *x.IsActive
	// 		domains = append(domains, domain)
	// 	}

	// 	association["domain"] = domains
	// 	association["patterns_to_match"] = *v.PatternsToMatch
	// 	associations = append(associations, association)
	// }

	return params
}
