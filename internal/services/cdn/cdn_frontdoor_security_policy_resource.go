package cdn

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	cdnfrontdoorsecurityparams "github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/frontdoorsecurityparams"
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

			"security_policies": {
				Type:     pluginsdk.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,

				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{

						"firewall": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{

									"cdn_frontdoor_firewall_policy_id": {
										Type:     pluginsdk.TypeString,
										Required: true,
										ForceNew: true,
									},

									"association": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										MaxItems: 1,

										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{

												// NOTE: The max number of domains vary depending on sku: 100 Standard, 500 Premium
												"domain": {
													Type:     pluginsdk.TypeList,
													Optional: true,
													MaxItems: 500,

													Elem: &pluginsdk.Resource{
														Schema: map[string]*pluginsdk.Schema{

															"cdn_frontdoor_custom_domain_id": {
																Type:         pluginsdk.TypeString,
																Required:     true,
																ValidateFunc: validate.FrontdoorCustomDomainID,
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
													Optional: true,
													MaxItems: 25,

													Elem: &pluginsdk.Schema{
														Type:    pluginsdk.TypeString,
														Default: "/*",
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

	standardSku := true
	profileClient := meta.(*clients.Client).Cdn.FrontdoorProfileClient
	profile, err := profileClient.Get(ctx, profileId.ResourceGroup, profileId.ProfileName)
	if err != nil {
		return fmt.Errorf("unable to retrieve the %q from the linked %q: %+v", "sku_name", "azurerm_cdn_frontdoor_profile", err)
	}

	skuName := flattenProfileSku(profile.Sku)
	if strings.HasPrefix(skuName, "Premium") {
		standardSku = false
	}

	params := track1.BasicSecurityPolicyPropertiesParameters(nil)

	if secPol, ok := d.GetOk("security_policies"); ok {
		params, err = expandCdnFrontdoorSecurityPoliciesParameters(secPol.([]interface{}), standardSku)
		if err != nil {
			return fmt.Errorf("expanding %q: %+v", "security_policies", err)
		}
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
		return fmt.Errorf("making Read request on Frontdoor Security Policy %q (Resource Group %q): %+v", id.SecurityPolicyName, id.ResourceGroup, err)
	}

	d.Set("name", id.SecurityPolicyName)
	d.Set("cdn_frontdoor_profile_id", parse.NewFrontdoorProfileID(id.SubscriptionId, id.ResourceGroup, id.ProfileName).ID())

	if props := resp.SecurityPolicyProperties; props != nil {
		securityPolicy, err := flattenCdnFrontdoorSecurityPoliciesParameters(props.Parameters)
		if err != nil {
			return fmt.Errorf("flattening %s: %+v", id, err)
		}

		d.Set("security_policies", *securityPolicy)
		d.Set("cdn_frontdoor_profile_name", id.ProfileName)
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

func expandCdnFrontdoorSecurityPoliciesParameters(input []interface{}, isStandardSku bool) (track1.BasicSecurityPolicyPropertiesParameters, error) {
	results, err := cdnfrontdoorsecurityparams.ExpandCdnFrontdoorFirewallPolicyParameters(input, isStandardSku)
	if err != nil {
		return results, err
	}

	if basic, ok := results.AsBasicSecurityPolicyPropertiesParameters(); ok {
		return basic, nil
	}

	return nil, nil
}

func flattenCdnFrontdoorSecurityPoliciesParameters(input track1.BasicSecurityPolicyPropertiesParameters) (*[]interface{}, error) {
	results, err := cdnfrontdoorsecurityparams.FlattenCdnFrontdoorFirewallPolicyParameters(input)
	if err != nil {
		return nil, err
	}

	return &results, nil
}
