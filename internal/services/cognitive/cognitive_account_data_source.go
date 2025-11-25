// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cognitive

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2025-06-01/cognitiveservicesaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceCognitiveAccount() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceCognitiveAccountRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"custom_question_answering_search_service_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"custom_subdomain_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"customer_managed_key": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"key_vault_key_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"identity_client_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"dynamic_throttling_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"fqdns": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"identity": commonschema.SystemAssignedUserAssignedIdentityComputed(),

			"kind": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"local_auth_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"metrics_advisor_aad_client_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"metrics_advisor_aad_tenant_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"metrics_advisor_super_user_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"metrics_advisor_website_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"network_acls": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"default_action": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"ip_rules": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},

						"virtual_network_rules": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"subnet_id": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},

									"ignore_missing_vnet_service_endpoint": {
										Type:     pluginsdk.TypeBool,
										Computed: true,
									},
								},
							},
						},

						"bypass": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"network_injection": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"scenario": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"subnet_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"outbound_network_access_restricted": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"primary_access_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"project_management_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"public_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"qna_runtime_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_access_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"sku_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"storage": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"storage_account_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"identity_client_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"tags": commonschema.Tags(),
		},
	}
}

func dataSourceCognitiveAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cognitive.AccountsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := cognitiveservicesaccounts.NewAccountID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	resp, err := client.AccountsGet(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))
		d.Set("kind", model.Kind)

		if sku := model.Sku; sku != nil {
			d.Set("sku_name", sku.Name)
		}

		if props := model.Properties; props != nil {
			if apiProps := props.ApiProperties; apiProps != nil {
				d.Set("custom_question_answering_search_service_id", pointer.From(apiProps.QnaAzureSearchEndpointId))
				d.Set("metrics_advisor_aad_client_id", pointer.From(apiProps.AadClientId))
				d.Set("metrics_advisor_aad_tenant_id", pointer.From(apiProps.AadTenantId))
				d.Set("metrics_advisor_super_user_name", pointer.From(apiProps.SuperUser))
				d.Set("metrics_advisor_website_name", pointer.From(apiProps.WebsiteName))
				d.Set("qna_runtime_endpoint", pointer.From(apiProps.QnaRuntimeEndpoint))
			}

			d.Set("custom_subdomain_name", pointer.From(props.CustomSubDomainName))

			customerManagedKey, err := flattenCognitiveAccountCustomerManagedKey(props.Encryption)
			if err != nil {
				return err
			}

			if err := d.Set("customer_managed_key", customerManagedKey); err != nil {
				return fmt.Errorf("setting `customer_managed_key`: %+v", err)
			}

			d.Set("endpoint", pointer.From(props.Endpoint))
			d.Set("dynamic_throttling_enabled", pointer.From(props.DynamicThrottlingEnabled))
			d.Set("fqdns", pointer.From(props.AllowedFqdnList))

			localAuthEnabled := !pointer.From(props.DisableLocalAuth)
			d.Set("local_auth_enabled", localAuthEnabled)

			if err := d.Set("network_acls", flattenCognitiveAccountDataSourceNetworkAcls(props.NetworkAcls)); err != nil {
				return fmt.Errorf("setting `network_acls` for %s: %+v", id, err)
			}

			networkInjection, err := flattenCognitiveAccountNetworkInjection(props.NetworkInjections)
			if err != nil {
				return err
			}

			if err := d.Set("network_injection", networkInjection); err != nil {
				return fmt.Errorf("setting `network_injection`: %+v", err)
			}

			d.Set("outbound_network_access_restricted", pointer.From(props.RestrictOutboundNetworkAccess))
			d.Set("project_management_enabled", pointer.From(props.AllowProjectManagement))
			d.Set("public_network_access_enabled", pointer.From(props.PublicNetworkAccess) == cognitiveservicesaccounts.PublicNetworkAccessEnabled)

			if err := d.Set("storage", flattenCognitiveAccountStorage(props.UserOwnedStorage)); err != nil {
				return fmt.Errorf("setting `storages` for %s: %+v", id, err)
			}

			if localAuthEnabled {
				keys, err := client.AccountsListKeys(ctx, id)
				if err != nil {
					return fmt.Errorf("listing the Keys for %s: %+v", id, err)
				}

				if model := keys.Model; model != nil {
					d.Set("primary_access_key", pointer.From(model.Key1))
					d.Set("secondary_access_key", pointer.From(model.Key2))
				}
			}
		}

		flattenedIdentity, err := identity.FlattenSystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		if err := d.Set("identity", flattenedIdentity); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		return tags.FlattenAndSet(d, model.Tags)
	}
	return nil
}

func flattenCognitiveAccountDataSourceNetworkAcls(input *cognitiveservicesaccounts.NetworkRuleSet) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	ipRules := make([]interface{}, 0)
	if input.IPRules != nil {
		for _, v := range *input.IPRules {
			ipRules = append(ipRules, v.Value)
		}
	}

	virtualNetworkRules := make([]interface{}, 0)
	if input.VirtualNetworkRules != nil {
		for _, v := range *input.VirtualNetworkRules {
			id := v.Id
			subnetId, err := commonids.ParseSubnetIDInsensitively(v.Id)
			if err == nil {
				id = subnetId.ID()
			}

			virtualNetworkRules = append(virtualNetworkRules, map[string]interface{}{
				"subnet_id":                            id,
				"ignore_missing_vnet_service_endpoint": pointer.From(v.IgnoreMissingVnetServiceEndpoint),
			})
		}
	}

	return []interface{}{map[string]interface{}{
		"bypass":                input.Bypass,
		"default_action":        input.DefaultAction,
		"ip_rules":              ipRules,
		"virtual_network_rules": virtualNetworkRules,
	}}
}
