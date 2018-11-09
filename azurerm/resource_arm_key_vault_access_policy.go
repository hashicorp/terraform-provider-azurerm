package azurerm

import (
	"fmt"
	"log"

	"strings"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/mgmt/2018-02-14/keyvault"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/satori/go.uuid"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

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

			"resource_group_name": resourceGroupNameSchema(),

			"tenant_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateUUID,
			},

			"object_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateUUID,
			},

			"application_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateUUID,
			},

			"certificate_permissions": azure.SchemaKeyVaultCertificatePermissions(),

			"key_permissions": azure.SchemaKeyVaultKeyPermissions(),

			"secret_permissions": azure.SchemaKeyVaultSecretPermissions(),
		},
	}
}

func resourceArmKeyVaultAccessPolicyCreateOrDelete(d *schema.ResourceData, meta interface{}, action keyvault.AccessPolicyUpdateKind) error {
	client := meta.(*ArmClient).keyVaultClient
	ctx := meta.(*ArmClient).StopContext
	log.Printf("[INFO] Preparing arguments for Key Vault Access Policy: %s.", action)

	vaultName := d.Get("vault_name").(string)
	resGroup := d.Get("resource_group_name").(string)

	tenantIdRaw := d.Get("tenant_id").(string)
	tenantId, err := uuid.FromString(tenantIdRaw)
	if err != nil {
		return fmt.Errorf("Error parsing Tenant ID %q as a UUID: %+v", tenantIdRaw, err)
	}

	objectId := d.Get("object_id").(string)

	certPermissionsRaw := d.Get("certificate_permissions").([]interface{})
	certPermissions := azure.ExpandCertificatePermissions(certPermissionsRaw)

	keyPermissionsRaw := d.Get("key_permissions").([]interface{})
	keyPermissions := azure.ExpandKeyPermissions(keyPermissionsRaw)

	secretPermissionsRaw := d.Get("secret_permissions").([]interface{})
	secretPermissions := azure.ExpandSecretPermissions(secretPermissionsRaw)

	accessPolicy := keyvault.AccessPolicyEntry{
		ObjectID: utils.String(objectId),
		TenantID: &tenantId,
		Permissions: &keyvault.Permissions{
			Certificates: certPermissions,
			Keys:         keyPermissions,
			Secrets:      secretPermissions,
		},
	}

	applicationIdRaw := d.Get("application_id").(string)
	if applicationIdRaw != "" {
		applicationId, err2 := uuid.FromString(applicationIdRaw)
		if err2 != nil {
			return fmt.Errorf("Error parsing Appliciation ID %q as a UUID: %+v", applicationIdRaw, err2)
		}

		accessPolicy.ApplicationID = &applicationId
	}

	accessPolicies := []keyvault.AccessPolicyEntry{accessPolicy}

	parameters := keyvault.VaultAccessPolicyParameters{
		Name: &vaultName,
		Properties: &keyvault.VaultAccessPolicyProperties{
			AccessPolicies: &accessPolicies,
		},
	}

	// Locking to prevent parallel changes causing issues
	azureRMLockByName(vaultName, keyVaultResourceName)
	defer azureRMUnlockByName(vaultName, keyVaultResourceName)

	_, err = client.UpdateAccessPolicy(ctx, resGroup, vaultName, action, parameters)
	if err != nil {
		return fmt.Errorf("Error updating Access Policy (Object ID %q / Application ID %q) for Key Vault %q (Resource Group %q): %+v", objectId, applicationIdRaw, vaultName, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, vaultName)
	if err != nil {
		return fmt.Errorf("Error retrieving Key Vault %q (Resource Group %q): %+v", vaultName, resGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read KeyVault %q (Resource Group %q) ID", vaultName, resGroup)
	}

	if d.IsNewResource() {
		// This is because azure doesn't have an 'id' for a keyvault access policy
		// In order to compensate for this and allow importing of this resource we are artificially
		// creating an identity for a key vault policy object
		resourceId := fmt.Sprintf("%s/objectId/%s", *read.ID, objectId)
		if applicationIdRaw != "" {
			resourceId = fmt.Sprintf("%s/applicationId/%s", resourceId, applicationIdRaw)
		}
		d.SetId(resourceId)
	}

	return nil
}

func resourceArmKeyVaultAccessPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	return resourceArmKeyVaultAccessPolicyCreateOrDelete(d, meta, keyvault.Add)
}

func resourceArmKeyVaultAccessPolicyDelete(d *schema.ResourceData, meta interface{}) error {
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
	vaultName := id.Path["vaults"]
	objectId := id.Path["objectId"]
	applicationId := id.Path["applicationId"]

	resp, err := client.Get(ctx, resGroup, vaultName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[ERROR] Key Vault %q (Resource Group %q) was not found - removing from state", vaultName, resGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Azure KeyVault %q (Resource Group %q): %+v", vaultName, resGroup, err)
	}

	policy, err := findKeyVaultAccessPolicy(resp.Properties.AccessPolicies, objectId, applicationId)
	if err != nil {
		return fmt.Errorf("Error locating Access Policy (Object ID %q / Application ID %q) in Key Vault %q (Resource Group %q)", objectId, applicationId, vaultName, resGroup)
	}

	if policy == nil {
		log.Printf("[ERROR] Access Policy (Object ID %q / Application ID %q) was not found in Key Vault %q (Resource Group %q) - removing from state", objectId, applicationId, vaultName, resGroup)
		d.SetId("")
		return nil
	}

	d.Set("vault_name", resp.Name)
	d.Set("resource_group_name", resGroup)
	d.Set("object_id", objectId)

	if tid := policy.TenantID; tid != nil {
		d.Set("tenant_id", tid.String())
	}

	if aid := policy.ApplicationID; aid != nil {
		d.Set("application_id", aid.String())
	}

	if permissions := policy.Permissions; permissions != nil {
		certificatePermissions := azure.FlattenCertificatePermissions(permissions.Certificates)
		if err := d.Set("certificate_permissions", certificatePermissions); err != nil {
			return fmt.Errorf("Error setting `certificate_permissions`: %+v", err)
		}

		keyPermissions := azure.FlattenKeyPermissions(permissions.Keys)
		if err := d.Set("key_permissions", keyPermissions); err != nil {
			return fmt.Errorf("Error setting `key_permissions`: %+v", err)
		}

		secretPermissions := azure.FlattenSecretPermissions(permissions.Secrets)
		if err := d.Set("secret_permissions", secretPermissions); err != nil {
			return fmt.Errorf("Error setting `secret_permissions`: %+v", err)
		}
	}

	return nil
}

func findKeyVaultAccessPolicy(policies *[]keyvault.AccessPolicyEntry, objectId string, applicationId string) (*keyvault.AccessPolicyEntry, error) {
	if policies == nil {
		return nil, nil
	}

	for _, policy := range *policies {
		if id := policy.ObjectID; id != nil {
			if strings.EqualFold(*id, objectId) {
				aid := ""
				if policy.ApplicationID != nil {
					aid = policy.ApplicationID.String()
				}

				if strings.EqualFold(aid, applicationId) {
					return &policy, nil
				}
			}
		}
	}

	return nil, nil
}
