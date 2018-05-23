package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/mgmt/2016-10-01/keyvault"
	"github.com/hashicorp/terraform/helper/schema"
	uuid "github.com/satori/go.uuid"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmKeyVaultPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmKeyVaultPolicyCreate,
		Read:   resourceArmKeyVaultPolicyRead,
		Update: resourceArmKeyVaultPolicyUpdate,
		Delete: resourceArmKeyVaultPolicyDelete,

		Schema: map[string]*schema.Schema{
			"vault_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"vault_resource_group": {
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

func createKeyVaultAccessPolicy(d *schema.ResourceData) *keyvault.AccessPolicyEntry {

	certificatePermissionsRaw := d.Get("certificate_permissions").([]interface{})
	certificatePermissions := []keyvault.CertificatePermissions{}
	for _, permission := range certificatePermissionsRaw {
		certificatePermissions = append(certificatePermissions, keyvault.CertificatePermissions(permission.(string)))
	}

	keyPermissionsRaw := d.Get("key_permissions").([]interface{})
	keyPermissions := []keyvault.KeyPermissions{}
	for _, permission := range keyPermissionsRaw {
		keyPermissions = append(keyPermissions, keyvault.KeyPermissions(permission.(string)))
	}

	secretPermissionsRaw := d.Get("secret_permissions").([]interface{})
	secretPermissions := []keyvault.SecretPermissions{}
	for _, permission := range secretPermissionsRaw {
		secretPermissions = append(secretPermissions, keyvault.SecretPermissions(permission.(string)))
	}

	policy := keyvault.AccessPolicyEntry{
		Permissions: &keyvault.Permissions{
			Certificates: &certificatePermissions,
			Keys:         &keyPermissions,
			Secrets:      &secretPermissions,
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

	return &policy

}

func resourceArmKeyVaultPolicyCreateOrDelete(d *schema.ResourceData, meta interface{}, action keyvault.AccessPolicyUpdateKind) error {
	client := meta.(*ArmClient).keyVaultClient
	ctx := meta.(*ArmClient).StopContext
	log.Printf("[INFO] preparing arguments for Azure ARM Policy: %s.", action)

	name := d.Get("vault_name").(string)
	resGroup := d.Get("vault_resource_group").(string)

	parameters := keyvault.VaultAccessPolicyParameters{
		Name: &name,
		Properties: &keyvault.VaultAccessPolicyProperties{
			AccessPolicies: &[]keyvault.AccessPolicyEntry{*createKeyVaultAccessPolicy(d)},
		},
	}

	_, err := client.UpdateAccessPolicy(ctx, resGroup, name, action, parameters)
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("cannot read KeyVault %s (resource group %s) ID", name, resGroup)
	}

	if action != keyvault.Remove {
		d.SetId(*read.ID)
		return resourceArmKeyVaultRead(d, meta)
	} else {
		return nil
	}
}

func resourceArmKeyVaultPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	return resourceArmKeyVaultPolicyCreateOrDelete(d, meta, keyvault.Add)
}

func resourceArmKeyVaultPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	return resourceArmKeyVaultPolicyCreateOrDelete(d, meta, keyvault.Remove)
}

func resourceArmKeyVaultPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceArmKeyVaultPolicyCreateOrDelete(d, meta, keyvault.Replace)
}

func resourceArmKeyVaultPolicyRead(d *schema.ResourceData, meta interface{}) error {
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

	policy := findKeyVaultAccessPolicy(objectUUID, resp.Properties.AccessPolicies)

	d.Set("vault_name", resp.Name)
	d.Set("vault_resource_group", resGroup)
	d.Set("tenant_id", resp.Properties.TenantID.String())
	d.Set("application_id", policy["application_id"])
	d.Set("key_permissions", policy["key_permissions"])
	d.Set("secret_permissions", policy["secret_permissions"])
	d.Set("certificate_permissions", policy["certificate_permissions"])

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func findKeyVaultAccessPolicy(objectID string, policies *[]keyvault.AccessPolicyEntry) map[string]interface{} {

	for _, policy := range *policies {
		result := make(map[string]interface{})

		if objectID != *policy.ObjectID {
			continue
		}

		keyPermissionsRaw := make([]interface{}, 0, len(*policy.Permissions.Keys))
		for _, keyPermission := range *policy.Permissions.Keys {
			keyPermissionsRaw = append(keyPermissionsRaw, string(keyPermission))
		}

		secretPermissionsRaw := make([]interface{}, 0, len(*policy.Permissions.Secrets))
		for _, secretPermission := range *policy.Permissions.Secrets {
			secretPermissionsRaw = append(secretPermissionsRaw, string(secretPermission))
		}

		result["tenant_id"] = policy.TenantID.String()
		result["object_id"] = *policy.ObjectID
		if policy.ApplicationID != nil {
			result["application_id"] = policy.ApplicationID.String()
		}
		result["key_permissions"] = keyPermissionsRaw
		result["secret_permissions"] = secretPermissionsRaw

		if policy.Permissions.Certificates != nil {
			certificatePermissionsRaw := make([]interface{}, 0, len(*policy.Permissions.Certificates))
			for _, certificatePermission := range *policy.Permissions.Certificates {
				certificatePermissionsRaw = append(certificatePermissionsRaw, string(certificatePermission))
			}
			result["certificate_permissions"] = certificatePermissionsRaw
		}

		return result
	}

	return nil
}
