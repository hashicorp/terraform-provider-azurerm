package keyvault

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/keyvault/mgmt/2020-04-01-preview/keyvault"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/set"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceKeyVault() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceKeyVaultRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: func() map[string]*pluginsdk.Schema {
			dsSchema := map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validate.VaultName,
				},

				"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

				"location": azure.SchemaLocationForDataSource(),

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
								Type:     pluginsdk.TypeSet,
								Computed: true,
								Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
								Set:      pluginsdk.HashString,
							},
							"virtual_network_subnet_ids": {
								Type:     pluginsdk.TypeSet,
								Computed: true,
								Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
								Set:      set.HashStringIgnoreCase,
							},
						},
					},
				},

				"purge_protection_enabled": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"tags": tags.SchemaDataSource(),
			}
			if !features.ThreePointOh() {
				dsSchema["soft_delete_enabled"] = &pluginsdk.Schema{
					Type:       pluginsdk.TypeBool,
					Computed:   true,
					Deprecated: `Azure has removed support for disabling Soft Delete as of 2020-12-15, as such this field will always return 'true' and will be removed in version 3.0 of the Azure Provider.`,
				}
			}
			return dsSchema
		}(),
	}
}

func dataSourceKeyVaultRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).KeyVault.VaultsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("KeyVault %q (Resource Group %q) does not exist", name, resourceGroup)
		}
		return fmt.Errorf("Error making Read request on KeyVault %q: %+v", name, err)
	}

	d.SetId(*resp.ID)

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.Properties; props != nil {
		d.Set("tenant_id", props.TenantID.String())
		d.Set("enabled_for_deployment", props.EnabledForDeployment)
		d.Set("enabled_for_disk_encryption", props.EnabledForDiskEncryption)
		d.Set("enabled_for_template_deployment", props.EnabledForTemplateDeployment)
		d.Set("purge_protection_enabled", props.EnablePurgeProtection)
		d.Set("vault_uri", props.VaultURI)

		// TODO: remove in 3.0
		if !features.ThreePointOh() {
			d.Set("soft_delete_enabled", true)
		}

		if sku := props.Sku; sku != nil {
			if err := d.Set("sku_name", string(sku.Name)); err != nil {
				return fmt.Errorf("Error setting `sku_name` for KeyVault %q: %+v", *resp.Name, err)
			}
		} else {
			return fmt.Errorf("Error making Read request on KeyVault %q: Unable to retrieve 'sku' value", *resp.Name)
		}

		flattenedPolicies := flattenAccessPolicies(props.AccessPolicies)
		if err := d.Set("access_policy", flattenedPolicies); err != nil {
			return fmt.Errorf("Error setting `access_policy` for KeyVault %q: %+v", *resp.Name, err)
		}

		if err := d.Set("network_acls", flattenKeyVaultDataSourceNetworkAcls(props.NetworkAcls)); err != nil {
			return fmt.Errorf("Error setting `network_acls` for KeyVault %q: %+v", *resp.Name, err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func flattenKeyVaultDataSourceNetworkAcls(input *keyvault.NetworkRuleSet) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make(map[string]interface{})

	output["bypass"] = string(input.Bypass)
	output["default_action"] = string(input.DefaultAction)

	ipRules := make([]interface{}, 0)
	if input.IPRules != nil {
		for _, v := range *input.IPRules {
			if v.Value == nil {
				continue
			}

			ipRules = append(ipRules, *v.Value)
		}
	}
	output["ip_rules"] = pluginsdk.NewSet(pluginsdk.HashString, ipRules)

	virtualNetworkRules := make([]interface{}, 0)
	if input.VirtualNetworkRules != nil {
		for _, v := range *input.VirtualNetworkRules {
			if v.ID == nil {
				continue
			}

			virtualNetworkRules = append(virtualNetworkRules, *v.ID)
		}
	}
	output["virtual_network_subnet_ids"] = pluginsdk.NewSet(pluginsdk.HashString, virtualNetworkRules)

	return []interface{}{output}
}
