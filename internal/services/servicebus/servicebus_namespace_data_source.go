// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package servicebus

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2024-01-01/namespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2024-01-01/namespacesauthorizationrule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceServiceBusNamespace() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceServiceBusNamespaceRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"identity": commonschema.SystemAssignedUserAssignedIdentityComputed(),

			"sku": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"capacity": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"premium_messaging_partitions": {
				Type:     pluginsdk.TypeInt,
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

						"identity_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"infrastructure_encryption_enabled": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},
					},
				},
			},

			"local_auth_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"public_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"minimum_tls_version": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"default_primary_connection_string": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"default_secondary_connection_string": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"default_primary_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"default_secondary_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"network_rule_set": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"default_action": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"public_network_access_enabled": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},

						"ip_rules": {
							Type:     pluginsdk.TypeSet,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},

						"trusted_services_allowed": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},

						"network_rules": {
							Type:     pluginsdk.TypeSet,
							Computed: true,
							Set:      networkRuleHash,
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
					},
				},
			},

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceServiceBusNamespaceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceBus.NamespacesClient
	namespaceAuthClient := meta.(*clients.Client).ServiceBus.NamespacesAuthClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := namespaces.NewNamespaceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		identityFlattened, err := identity.FlattenSystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		if err := d.Set("identity", identityFlattened); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if sku := model.Sku; sku != nil {
			d.Set("sku", string(sku.Name))
			d.Set("capacity", sku.Capacity)
		}

		if props := model.Properties; props != nil {
			d.Set("premium_messaging_partitions", props.PremiumMessagingPartitions)

			if customerManagedKey, err := flattenServiceBusNamespaceEncryption(props.Encryption); err == nil {
				d.Set("customer_managed_key", customerManagedKey)
			}

			localAuthEnabled := true
			if props.DisableLocalAuth != nil {
				localAuthEnabled = !*props.DisableLocalAuth
			}
			d.Set("local_auth_enabled", localAuthEnabled)

			publicNetworkAccess := true
			if props.PublicNetworkAccess != nil && *props.PublicNetworkAccess == namespaces.PublicNetworkAccessDisabled {
				publicNetworkAccess = false
			}
			d.Set("public_network_access_enabled", publicNetworkAccess)

			if props.MinimumTlsVersion != nil {
				d.Set("minimum_tls_version", string(pointer.From(props.MinimumTlsVersion)))
			}

			d.Set("endpoint", props.ServiceBusEndpoint)
		}
	}

	authRuleId := namespacesauthorizationrule.NewAuthorizationRuleID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, serviceBusNamespaceDefaultAuthorizationRule)

	keys, err := namespaceAuthClient.NamespacesListKeys(ctx, authRuleId)
	if err != nil {
		log.Printf("[WARN] listing default keys for %s: %+v", id, err)
	} else {
		if keysModel := keys.Model; keysModel != nil {
			d.Set("default_primary_connection_string", keysModel.PrimaryConnectionString)
			d.Set("default_secondary_connection_string", keysModel.SecondaryConnectionString)
			d.Set("default_primary_key", keysModel.PrimaryKey)
			d.Set("default_secondary_key", keysModel.SecondaryKey)
		}
	}

	networkRuleSet, err := client.GetNetworkRuleSet(ctx, id)
	if err != nil {
		return fmt.Errorf("retrieving network rule set %s: %+v", id, err)
	}

	if nrsModel := networkRuleSet.Model; nrsModel != nil {
		if props := nrsModel.Properties; props != nil {
			d.Set("network_rule_set", flattenServiceBusNamespaceNetworkRuleSet(*props))
		}
	}

	return nil
}
