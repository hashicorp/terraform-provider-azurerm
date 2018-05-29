package azurerm

import (
	"fmt"
	"log"
	"net"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/mgmt/2016-10-01/keyvault"
	//"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/go-getter/helper/url"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	uuid "github.com/satori/go.uuid"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

// As can be seen in the API definition, the Sku Family only supports the value
// `A` and is a required field
// https://github.com/Azure/azure-rest-api-specs/blob/master/arm-keyvault/2015-06-01/swagger/keyvault.json#L239
var armKeyVaultSkuFamily = "A"

var keyVaultAccessPolicyResourceName = "azurerm_key_vault_access_policy"

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
				Computed: true,
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
						"certificate_permissions": certificatePermissionsSchema(),
						"key_permissions":         keyPermissionsSchema(),
						"secret_permissions":      secretPermissionsSchema(),
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

	//return fmt.Errorf("%s", spew.Sdump(parameters))

	azureRMLockByName(name, keyVaultResourceName)
	defer azureRMUnlockByName(name, keyVaultResourceName)

	_, err := client.CreateOrUpdate(ctx, resGroup, name, parameters)
	if err != nil {
		return fmt.Errorf("Error updating Key Vault (Key Vault %q / Resource Group %q): %+v", name, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Key Vault %q (Resource Group %q): %+v", name, resGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read KeyVault %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	if d.IsNewResource() {
		if props := read.Properties; props != nil {
			if vault := props.VaultURI; vault != nil {
				err := resource.Retry(120*time.Second, checkKeyVaultDNSIsAvailable(*vault))
				if err != nil {
					return err
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

	azureRMLockByName(d.Get("name").(string), keyVaultAccessPolicyResourceName)
	defer azureRMUnlockByName(d.Get("name").(string), keyVaultAccessPolicyResourceName)

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
	d.Set("tenant_id", resp.Properties.TenantID.String())
	d.Set("enabled_for_deployment", resp.Properties.EnabledForDeployment)
	d.Set("enabled_for_disk_encryption", resp.Properties.EnabledForDiskEncryption)
	d.Set("enabled_for_template_deployment", resp.Properties.EnabledForTemplateDeployment)
	d.Set("sku", flattenKeyVaultSku(resp.Properties.Sku))
	d.Set("access_policy", flattenKeyVaultAccessPolicies(resp.Properties.AccessPolicies))
	d.Set("vault_uri", resp.Properties.VaultURI)

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

	azureRMLockByName(name, keyVaultResourceName)
	defer azureRMUnlockByName(name, keyVaultResourceName)

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

func flattenKeyVaultSku(sku *keyvault.Sku) []interface{} {
	result := map[string]interface{}{
		"name": string(sku.Name),
	}

	return []interface{}{result}
}

func validateKeyVaultName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if matched := regexp.MustCompile(`^[a-zA-Z0-9-]{3,24}$`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%q may only contain alphanumeric characters and dashes and must be between 3-24 chars", k))
	}

	return
}

func checkKeyVaultDNSIsAvailable(vaultUri string) func() *resource.RetryError {
	return func() *resource.RetryError {
		uri, err := url.Parse(vaultUri)
		if err != nil {
			return resource.NonRetryableError(err)
		}

		conn, err := net.Dial("tcp", fmt.Sprintf("%s:443", uri.Host))
		if err != nil {
			return resource.RetryableError(err)
		}

		_ = conn.Close()
		return nil
	}
}

func expandKeyVaultAccessPolicies(d *schema.ResourceData) *[]keyvault.AccessPolicyEntry {
	policies := d.Get("access_policy").([]interface{})
	result := make([]keyvault.AccessPolicyEntry, 0, len(policies))

	for _, policySet := range policies {
		policyRaw := policySet.(map[string]interface{})

		certificatePermissionsRaw := policyRaw["certificate_permissions"].([]interface{})
		keyPermissionsRaw := policyRaw["key_permissions"].([]interface{})
		secretPermissionsRaw := policyRaw["secret_permissions"].([]interface{})

		policy := keyvault.AccessPolicyEntry{
			Permissions: &keyvault.Permissions{
				Certificates: expandKeyVaultAccessPolicyCertificatePermissions(certificatePermissionsRaw),
				Keys:         expandKeyVaultAccessPolicyKeyPermissions(keyPermissionsRaw),
				Secrets:      expandKeyVaultAccessPolicySecretPermissions(secretPermissionsRaw),
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
