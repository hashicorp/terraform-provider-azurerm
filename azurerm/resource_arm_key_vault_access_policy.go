package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/mgmt/2016-10-01/keyvault"
	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/terraform/helper/schema"
	uuid "github.com/satori/go.uuid"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

var keyVaultResourceName = "azurerm_key_vault"

func resourceArmKeyVaultAccessPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmKeyVaultAccessPolicyCreate,
		Read:   resourceArmKeyVaultAccessPolicyRead,
		Update: resourceArmKeyVaultAccessPolicyUpdate,
		Delete: resourceArmKeyVaultAccessPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"vault_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"tenant_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateUUID,
				ForceNew:     true,
			},
			"object_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateUUID,
				ForceNew:     true,
			},
			"application_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateUUID,
				ForceNew:     true,
			},
			"certificate_permissions": certificatePermissionsSchema(),
			"key_permissions":         keyPermissionsSchema(),
			"secret_permissions":      secretPermissionsSchema(),
		},
	}
}

func resourceArmKeyVaultAccessPolicyCreateOrDelete(d *schema.ResourceData, meta interface{}, action keyvault.AccessPolicyUpdateKind) error {
	client := meta.(*ArmClient).keyVaultClient
	ctx := meta.(*ArmClient).StopContext
	log.Printf("[INFO]Preparing arguments for Key Vault Access Policy: %s.", action)

	name := d.Get("vault_name").(string)
	resGroup := d.Get("resource_group_name").(string)

	accessPolicy := expandKeyVaultAccessPolicy(d)
	accessPolicies := []keyvault.AccessPolicyEntry{accessPolicy}

	log.Printf("SPEW: %s", spew.Sdump(accessPolicies))

	parameters := keyvault.VaultAccessPolicyParameters{
		Name: &name,
		Properties: &keyvault.VaultAccessPolicyProperties{
			AccessPolicies: &accessPolicies,
		},
	}

	_, err := client.UpdateAccessPolicy(ctx, resGroup, name, action, parameters)
	if err != nil {
		return fmt.Errorf("Error updating Key Vault Access Policy (Key Vault %q / Resource Group %q): %+v", name, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Key Vault %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("cannot read KeyVault %s (resource group %s) ID", name, resGroup)
	}

	if d.IsNewResource() {
		d.SetId(*read.ID)
	}

	return nil
}

func resourceArmKeyVaultAccessPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	return resourceArmKeyVaultAccessPolicyCreateOrDelete(d, meta, keyvault.Add)
}

func resourceArmKeyVaultAccessPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	azureRMLockByName(d.Get("vault_name").(string), keyVaultAccessPolicyResourceName)
	defer azureRMUnlockByName(d.Get("vault_name").(string), keyVaultAccessPolicyResourceName)

	return resourceArmKeyVaultAccessPolicyCreateOrDelete(d, meta, keyvault.Remove)
}

func resourceArmKeyVaultAccessPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceArmKeyVaultAccessPolicyCreateOrDelete(d, meta, keyvault.Replace)
}

func resourceArmKeyVaultAccessPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).keyVaultClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["vaults"]

	objectUUID := d.Get("object_id").(string)

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure KeyVault %s: %+v", name, err)
	}

	flattenedPolicy := flattenKeyVaultAccessPolicies(resp.Properties.AccessPolicies)
	policy := findKeyVaultAccessPolicy(objectUUID, flattenedPolicy)

	d.Set("vault_name", resp.Name)
	d.Set("resource_group_name", resGroup)
	d.Set("tenant_id", resp.Properties.TenantID.String())

	d.Set("application_id", policy["application_id"])
	d.Set("key_permissions", policy["key_permissions"])
	d.Set("secret_permissions", policy["secret_permissions"])
	d.Set("certificate_permissions", policy["certificate_permissions"])

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func findKeyVaultAccessPolicy(objectID string, policies []map[string]interface{}) map[string]interface{} {
	for _, policy := range policies {
		if policy["object_id"] != nil && policy["object_id"] == objectID {
			return policy
		}
	}
	return nil
}

func expandKeyVaultAccessPolicy(d *schema.ResourceData) keyvault.AccessPolicyEntry {

	policy := keyvault.AccessPolicyEntry{
		Permissions: &keyvault.Permissions{
			Certificates: expandKeyVaultAccessPolicyCertificatePermissions(d.Get("certificate_permissions").([]interface{})),
			Keys:         expandKeyVaultAccessPolicyKeyPermissions(d.Get("key_permissions").([]interface{})),
			Secrets:      expandKeyVaultAccessPolicySecretPermissions(d.Get("secret_permissions").([]interface{})),
		},
	}

	tenantUUID := uuid.FromStringOrNil(d.Get("tenant_id").(string))
	policy.TenantID = &tenantUUID
	objectUUID := d.Get("object_id").(string)
	policy.ObjectID = &objectUUID

	if v := d.Get("application_id"); v != "" {
		applicationUUID := uuid.FromStringOrNil(v.(string))
		policy.ApplicationID = &applicationUUID
	}

	return policy
}
