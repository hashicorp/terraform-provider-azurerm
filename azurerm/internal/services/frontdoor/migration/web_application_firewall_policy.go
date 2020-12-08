package migration

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/frontdoor/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
)

func WebApplicationFirewallPolicyV0Schema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"mode": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"redirect_url": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"custom_block_response_status_code": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"custom_block_response_body": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"custom_rule": {
				Type:     schema.TypeList,
				MaxItems: 100,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},

						"priority": {
							Type:     schema.TypeInt,
							Optional: true,
						},

						"type": {
							Type:     schema.TypeString,
							Required: true,
						},

						"rate_limit_duration_in_minutes": {
							Type:     schema.TypeInt,
							Optional: true,
						},

						"rate_limit_threshold": {
							Type:     schema.TypeInt,
							Optional: true,
						},

						"action": {
							Type:     schema.TypeString,
							Required: true,
						},

						"match_condition": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 100,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"match_variable": {
										Type:     schema.TypeString,
										Required: true,
									},

									"match_values": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},

									"operator": {
										Type:     schema.TypeString,
										Required: true,
									},

									"selector": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},

									"negation_condition": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},

									"transforms": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 5,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},

			"managed_rule": {
				Type:     schema.TypeList,
				MaxItems: 100,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},

						"version": {
							Type:     schema.TypeString,
							Required: true,
						},

						"exclusion": {
							Type:     schema.TypeList,
							MaxItems: 100,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"match_variable": {
										Type:     schema.TypeString,
										Required: true,
									},
									"operator": {
										Type:     schema.TypeString,
										Required: true,
									},
									"selector": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},

						"override": {
							Type:     schema.TypeList,
							MaxItems: 100,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"rule_group_name": {
										Type:     schema.TypeString,
										Required: true,
									},

									"exclusion": {
										Type:     schema.TypeList,
										MaxItems: 100,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"match_variable": {
													Type:     schema.TypeString,
													Required: true,
												},
												"operator": {
													Type:     schema.TypeString,
													Required: true,
												},
												"selector": {
													Type:     schema.TypeString,
													Required: true,
												},
											},
										},
									},

									"rule": {
										Type:     schema.TypeList,
										MaxItems: 1000,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"rule_id": {
													Type:     schema.TypeString,
													Required: true,
												},

												"enabled": {
													Type:     schema.TypeBool,
													Optional: true,
												},

												"exclusion": {
													Type:     schema.TypeList,
													MaxItems: 100,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"match_variable": {
																Type:     schema.TypeString,
																Required: true,
															},
															"operator": {
																Type:     schema.TypeString,
																Required: true,
															},
															"selector": {
																Type:     schema.TypeString,
																Required: true,
															},
														},
													},
												},

												"action": {
													Type:     schema.TypeString,
													Required: true,
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

			"frontend_endpoint_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func WebApplicationFirewallPolicyV0ToV1(rawState map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
	// old
	// 	/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/frontdoorwebapplicationfirewallpolicies/{policyName}
	// new:
	// 	/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/frontDoorWebApplicationFirewallPolicies/{policyName}
	oldId := rawState["id"].(string)
	oldParsedId, err := azure.ParseAzureResourceID(oldId)
	if err != nil {
		return rawState, err
	}

	resourceGroup := oldParsedId.ResourceGroup
	policyName := ""
	for key, value := range oldParsedId.Path {
		if strings.EqualFold(key, "frontDoorWebApplicationFirewallPolicies") {
			policyName = value
			break
		}
	}

	if policyName == "" {
		return rawState, fmt.Errorf("couldn't find the `frontDoorWebApplicationFirewallPolicies` segment in the old resource id %q", oldId)
	}

	newId := parse.NewWebApplicationFirewallPolicyID(resourceGroup, policyName)
	newIdStr := newId.ID(oldParsedId.SubscriptionID)

	log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newIdStr)

	rawState["id"] = newIdStr

	return rawState, nil
}
