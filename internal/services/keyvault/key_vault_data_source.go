// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package keyvault

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2023-02-01/vaults"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceKeyVault() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceKeyVaultRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.VaultName,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"sku_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"vault_uri": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tenant_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"access_policy": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"tenant_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"object_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"application_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"certificate_permissions": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
						"key_permissions": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
						"secret_permissions": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
						"storage_permissions": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
					},
				},
			},

			"enabled_for_deployment": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"enabled_for_disk_encryption": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"enabled_for_template_deployment": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			// TODO 4.0: change this from enable_* to *_enabled
			"enable_rbac_authorization": {
				Type:     pluginsdk.TypeBool,
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
						"bypass": {
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
						"virtual_network_subnet_ids": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
					},
				},
			},

			"purge_protection_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"public_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceKeyVaultRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).KeyVault.VaultsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewKeyVaultID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("making read request %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.VaultName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		props := model.Properties
		d.Set("enabled_for_deployment", props.EnabledForDeployment)
		d.Set("enabled_for_disk_encryption", props.EnabledForDiskEncryption)
		d.Set("enabled_for_template_deployment", props.EnabledForTemplateDeployment)
		d.Set("enable_rbac_authorization", props.EnableRbacAuthorization)
		d.Set("purge_protection_enabled", props.EnablePurgeProtection)
		if v := props.PublicNetworkAccess; v != nil {
			d.Set("public_network_access_enabled", *v == "Enabled")
		}
		d.Set("tenant_id", props.TenantId)

		d.Set("vault_uri", props.VaultUri)
		if props.VaultUri != nil {
			meta.(*clients.Client).KeyVault.AddToCache(id, *props.VaultUri)
		}

		skuName := ""
		// the Azure API is inconsistent here, so rewrite this into the casing we expect
		// TODO: this can be removed when the new base layer is enabled?
		for _, v := range vaults.PossibleValuesForSkuName() {
			if strings.EqualFold(v, string(model.Properties.Sku.Name)) {
				skuName = v
			}
		}
		d.Set("sku_name", skuName)

		flattenedPolicies := flattenAccessPolicies(props.AccessPolicies)
		if err := d.Set("access_policy", flattenedPolicies); err != nil {
			return fmt.Errorf("setting `access_policy`: %+v", err)
		}

		if err := d.Set("network_acls", flattenKeyVaultDataSourceNetworkAcls(props.NetworkAcls)); err != nil {
			return fmt.Errorf("setting `network_acls`: %+v", err)
		}
		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("setting `tags`: %+v", err)
		}
	}

	return nil
}

func flattenKeyVaultDataSourceNetworkAcls(input *vaults.NetworkRuleSet) []interface{} {
	output := make([]interface{}, 0)

	if input != nil {
		bypass := ""
		if input.Bypass != nil {
			bypass = string(*input.Bypass)
		}

		defaultAction := ""
		if input.DefaultAction != nil {
			defaultAction = string(*input.DefaultAction)
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
				virtualNetworkRules = append(virtualNetworkRules, v.Id)
			}
		}

		output = append(output, map[string]interface{}{
			"bypass":                     bypass,
			"default_action":             defaultAction,
			"ip_rules":                   ipRules,
			"virtual_network_subnet_ids": virtualNetworkRules,
		})
	}

	return output
}
