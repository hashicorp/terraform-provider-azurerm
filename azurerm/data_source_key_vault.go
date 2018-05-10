package azurerm

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/mgmt/2016-10-01/keyvault"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmKeyVault() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmKeyVaultRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateKeyVaultName,
			},

			"resource_group_name": resourceGroupNameForDataSourceSchema(),

			"location": locationForDataSourceSchema(),

			"sku": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"vault_uri": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"access_policy": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"object_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"application_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"certificate_permissions": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"key_permissions": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"secret_permissions": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"enabled_for_deployment": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"enabled_for_disk_encryption": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"enabled_for_template_deployment": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"tags": tagsForDataSourceSchema(),
		},
	}
}

func dataSourceArmKeyVaultRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).keyVaultClient
	ctx := meta.(*ArmClient).StopContext

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
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	if props := resp.Properties; props != nil {
		d.Set("tenant_id", props.TenantID.String())
		d.Set("enabled_for_deployment", props.EnabledForDeployment)
		d.Set("enabled_for_disk_encryption", props.EnabledForDiskEncryption)
		d.Set("enabled_for_template_deployment", props.EnabledForTemplateDeployment)
		if err := d.Set("sku", flattenKeyVaultDataSourceSku(props.Sku)); err != nil {
			return fmt.Errorf("Error flattening `sku` for KeyVault %q: %+v", resp.Name, err)
		}
		if err := d.Set("access_policy", flattenKeyVaultDataSourceAccessPolicies(props.AccessPolicies)); err != nil {
			return fmt.Errorf("Error flattening `access_policy` for KeyVault %q: %+v", resp.Name, err)
		}
		d.Set("vault_uri", props.VaultURI)
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func flattenKeyVaultDataSourceSku(sku *keyvault.Sku) []interface{} {
	result := map[string]interface{}{
		"name": string(sku.Name),
	}

	return []interface{}{result}
}

func flattenKeyVaultDataSourceAccessPolicies(policies *[]keyvault.AccessPolicyEntry) []interface{} {
	result := make([]interface{}, 0, len(*policies))

	if policies == nil {
		return result
	}

	for _, policy := range *policies {
		policyRaw := make(map[string]interface{})

		keyPermissionsRaw := make([]interface{}, 0)
		secretPermissionsRaw := make([]interface{}, 0)
		certificatePermissionsRaw := make([]interface{}, 0)

		if permissions := policy.Permissions; permissions != nil {
			if keys := permissions.Keys; keys != nil {
				for _, keyPermission := range *keys {
					keyPermissionsRaw = append(keyPermissionsRaw, string(keyPermission))
				}
			}
			if secrets := permissions.Secrets; secrets != nil {
				for _, secretPermission := range *secrets {
					secretPermissionsRaw = append(secretPermissionsRaw, string(secretPermission))
				}
			}

			if certificates := permissions.Certificates; certificates != nil {
				for _, certificatePermission := range *certificates {
					certificatePermissionsRaw = append(certificatePermissionsRaw, string(certificatePermission))
				}
			}
		}

		policyRaw["tenant_id"] = policy.TenantID.String()
		if policy.ObjectID != nil {
			policyRaw["object_id"] = *policy.ObjectID
		}
		if policy.ApplicationID != nil {
			policyRaw["application_id"] = policy.ApplicationID.String()
		}
		policyRaw["key_permissions"] = keyPermissionsRaw
		policyRaw["secret_permissions"] = secretPermissionsRaw
		policyRaw["certificate_permissions"] = certificatePermissionsRaw

		result = append(result, policyRaw)
	}

	return result
}
