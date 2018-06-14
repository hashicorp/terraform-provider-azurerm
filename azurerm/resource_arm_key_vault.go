package azurerm

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/mgmt/2016-10-01/keyvault"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/satori/go.uuid"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

// As can be seen in the API definition, the Sku Family only supports the value
// `A` and is a required field
// https://github.com/Azure/azure-rest-api-specs/blob/master/arm-keyvault/2015-06-01/swagger/keyvault.json#L239
var armKeyVaultSkuFamily = "A"

func resourceArmKeyVault() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmKeyVaultCreate,
		Read:   resourceArmKeyVaultRead,
		Update: resourceArmKeyVaultCreate,
		Delete: resourceArmKeyVaultDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		MigrateState:  resourceAzureRMKeyVaultMigrateState,
		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateKeyVaultName,
			},

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"sku": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(keyvault.Standard),
								string(keyvault.Premium),
							}, false),
						},
					},
				},
			},

			"vault_uri": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tenant_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateUUID,
			},

			"access_policy": {
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 1,
				MaxItems: 16,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tenant_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateUUID,
						},
						"object_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateUUID,
						},
						"application_id": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateUUID,
						},
						"certificate_permissions": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									string(keyvault.Create),
									string(keyvault.Delete),
									string(keyvault.Deleteissuers),
									string(keyvault.Get),
									string(keyvault.Getissuers),
									string(keyvault.Import),
									string(keyvault.List),
									string(keyvault.Listissuers),
									string(keyvault.Managecontacts),
									string(keyvault.Manageissuers),
									string(keyvault.Purge),
									string(keyvault.Recover),
									string(keyvault.Setissuers),
									string(keyvault.Update),
								}, true),
								DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
							},
						},
						"key_permissions": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									string(keyvault.KeyPermissionsBackup),
									string(keyvault.KeyPermissionsCreate),
									string(keyvault.KeyPermissionsDecrypt),
									string(keyvault.KeyPermissionsDelete),
									string(keyvault.KeyPermissionsEncrypt),
									string(keyvault.KeyPermissionsGet),
									string(keyvault.KeyPermissionsImport),
									string(keyvault.KeyPermissionsList),
									string(keyvault.KeyPermissionsPurge),
									string(keyvault.KeyPermissionsRecover),
									string(keyvault.KeyPermissionsRestore),
									string(keyvault.KeyPermissionsSign),
									string(keyvault.KeyPermissionsUnwrapKey),
									string(keyvault.KeyPermissionsUpdate),
									string(keyvault.KeyPermissionsVerify),
									string(keyvault.KeyPermissionsWrapKey),
								}, true),
								DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
							},
						},
						"secret_permissions": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									string(keyvault.SecretPermissionsBackup),
									string(keyvault.SecretPermissionsDelete),
									string(keyvault.SecretPermissionsGet),
									string(keyvault.SecretPermissionsList),
									string(keyvault.SecretPermissionsPurge),
									string(keyvault.SecretPermissionsRecover),
									string(keyvault.SecretPermissionsRestore),
									string(keyvault.SecretPermissionsSet),
								}, true),
								DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
							},
						},
					},
				},
			},

			"enabled_for_deployment": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"enabled_for_disk_encryption": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"enabled_for_template_deployment": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmKeyVaultCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).keyVaultClient
	ctx := meta.(*ArmClient).StopContext
	log.Printf("[INFO] preparing arguments for Azure ARM KeyVault creation.")

	name := d.Get("name").(string)
	location := azureRMNormalizeLocation(d.Get("location").(string))
	resGroup := d.Get("resource_group_name").(string)
	tenantUUID := uuid.FromStringOrNil(d.Get("tenant_id").(string))
	enabledForDeployment := d.Get("enabled_for_deployment").(bool)
	enabledForDiskEncryption := d.Get("enabled_for_disk_encryption").(bool)
	enabledForTemplateDeployment := d.Get("enabled_for_template_deployment").(bool)
	tags := d.Get("tags").(map[string]interface{})

	parameters := keyvault.VaultCreateOrUpdateParameters{
		Location: &location,
		Properties: &keyvault.VaultProperties{
			TenantID:                     &tenantUUID,
			Sku:                          expandKeyVaultSku(d),
			AccessPolicies:               expandKeyVaultAccessPolicies(d),
			EnabledForDeployment:         &enabledForDeployment,
			EnabledForDiskEncryption:     &enabledForDiskEncryption,
			EnabledForTemplateDeployment: &enabledForTemplateDeployment,
		},
		Tags: expandTags(tags),
	}

	_, err := client.CreateOrUpdate(ctx, resGroup, name, parameters)
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read KeyVault %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	if d.IsNewResource() {
		if props := read.Properties; props != nil {
			if vault := props.VaultURI; vault != nil {
				log.Printf("[DEBUG] Waiting for Key Vault %q (Resource Group %q) to become available", name, resGroup)
				stateConf := &resource.StateChangeConf{
					Pending:                   []string{"pending"},
					Target:                    []string{"available"},
					Refresh:                   keyVaultRefreshFunc(*vault),
					Timeout:                   30 * time.Minute,
					Delay:                     30 * time.Second,
					PollInterval:              10 * time.Second,
					ContinuousTargetOccurence: 10,
				}

				if _, err := stateConf.WaitForState(); err != nil {
					return fmt.Errorf("Error waiting for Key Vault %q (Resource Group %q) to become available: %s", name, resGroup, err)
				}
			}
		}
	}

	return resourceArmKeyVaultRead(d, meta)
}

func resourceArmKeyVaultRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).keyVaultClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["vaults"]

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure KeyVault %s: %+v", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	if props := resp.Properties; props != nil {
		d.Set("tenant_id", props.TenantID.String())
		d.Set("enabled_for_deployment", props.EnabledForDeployment)
		d.Set("enabled_for_disk_encryption", props.EnabledForDiskEncryption)
		d.Set("enabled_for_template_deployment", props.EnabledForTemplateDeployment)
		if err := d.Set("sku", flattenKeyVaultSku(props.Sku)); err != nil {
			return fmt.Errorf("Error flattening `sku` for KeyVault %q: %+v", resp.Name, err)
		}
		if err := d.Set("access_policy", flattenKeyVaultAccessPolicies(props.AccessPolicies)); err != nil {
			return fmt.Errorf("Error flattening `access_policy` for KeyVault %q: %+v", resp.Name, err)
		}
		d.Set("vault_uri", props.VaultURI)
	}

	flattenAndSetTags(d, resp.Tags)
	return nil
}

func resourceArmKeyVaultDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).keyVaultClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["vaults"]

	_, err = client.Delete(ctx, resGroup, name)

	return err
}

func expandKeyVaultSku(d *schema.ResourceData) *keyvault.Sku {
	skuSets := d.Get("sku").([]interface{})
	sku := skuSets[0].(map[string]interface{})

	return &keyvault.Sku{
		Family: &armKeyVaultSkuFamily,
		Name:   keyvault.SkuName(sku["name"].(string)),
	}
}

func expandKeyVaultAccessPolicies(d *schema.ResourceData) *[]keyvault.AccessPolicyEntry {
	policies := d.Get("access_policy").([]interface{})
	result := make([]keyvault.AccessPolicyEntry, 0, len(policies))

	for _, policySet := range policies {
		policyRaw := policySet.(map[string]interface{})

		certificatePermissionsRaw := policyRaw["certificate_permissions"].([]interface{})
		certificatePermissions := make([]keyvault.CertificatePermissions, 0)
		for _, v := range certificatePermissionsRaw {
			permission := keyvault.CertificatePermissions(v.(string))
			certificatePermissions = append(certificatePermissions, permission)
		}

		keyPermissionsRaw := policyRaw["key_permissions"].([]interface{})
		keyPermissions := make([]keyvault.KeyPermissions, 0)
		for _, v := range keyPermissionsRaw {
			permission := keyvault.KeyPermissions(v.(string))
			keyPermissions = append(keyPermissions, permission)
		}

		secretPermissionsRaw := policyRaw["secret_permissions"].([]interface{})
		secretPermissions := make([]keyvault.SecretPermissions, 0)
		for _, v := range secretPermissionsRaw {
			permission := keyvault.SecretPermissions(v.(string))
			secretPermissions = append(secretPermissions, permission)
		}

		policy := keyvault.AccessPolicyEntry{
			Permissions: &keyvault.Permissions{
				Certificates: &certificatePermissions,
				Keys:         &keyPermissions,
				Secrets:      &secretPermissions,
			},
		}

		tenantUUID := uuid.FromStringOrNil(policyRaw["tenant_id"].(string))
		policy.TenantID = &tenantUUID
		objectUUID := policyRaw["object_id"].(string)
		policy.ObjectID = &objectUUID

		if v := policyRaw["application_id"]; v != "" {
			applicationUUID := uuid.FromStringOrNil(v.(string))
			policy.ApplicationID = &applicationUUID
		}

		result = append(result, policy)
	}

	return &result
}

func flattenKeyVaultSku(sku *keyvault.Sku) []interface{} {
	result := map[string]interface{}{
		"name": string(sku.Name),
	}

	return []interface{}{result}
}

func flattenKeyVaultAccessPolicies(policies *[]keyvault.AccessPolicyEntry) []interface{} {
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

func validateKeyVaultName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if matched := regexp.MustCompile(`^[a-zA-Z0-9-]{3,24}$`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%q may only contain alphanumeric characters and dashes and must be between 3-24 chars", k))
	}

	return
}

func keyVaultRefreshFunc(vaultUri string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Checking to see if KeyVault %q is available..", vaultUri)

		var PTransport = &http.Transport{Proxy: http.ProxyFromEnvironment}

		client := &http.Client{
			Transport: PTransport,
		}

		conn, err := client.Get(vaultUri)
		if err != nil {
			log.Printf("[DEBUG] Didn't find KeyVault at %q", vaultUri)
			return nil, "pending", fmt.Errorf("Error connecting to %q: %s", vaultUri, err)
		}

		defer conn.Body.Close()

		log.Printf("[DEBUG] Found KeyVault at %q", vaultUri)
		return "available", "available", nil
	}
}
