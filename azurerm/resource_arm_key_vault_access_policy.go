package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/mgmt/2016-10-01/keyvault"
	"github.com/hashicorp/terraform/helper/schema"
	uuid "github.com/satori/go.uuid"
	keyVaultHelper "github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/keyvault"
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
			"certificate_permissions": keyVaultHelper.CertificatePermissionsSchema(),
			"key_permissions":         keyVaultHelper.KeyPermissionsSchema(),
			"secret_permissions":      keyVaultHelper.SecretPermissionsSchema(),
		},
	}
}

func resourceArmKeyVaultAccessPolicyCreateOrDelete(d *schema.ResourceData, meta interface{}, action keyvault.AccessPolicyUpdateKind) error {
	client := meta.(*ArmClient).keyVaultClient
	ctx := meta.(*ArmClient).StopContext
	log.Printf("[INFO] Preparing arguments for Key Vault Access Policy: %s.", action)

	name := d.Get("vault_name").(string)
	resGroup := d.Get("resource_group_name").(string)

	accessPolicies := []keyvault.AccessPolicyEntry{expandKeyVaultAccessPolicy(d)}

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
		// This is because azure doesn't have an 'id' for a keyvault access policy
		// In order to compensate for this and allow importing of this resource we are artificially
		// creating an identity for a key vault policy object
		resourceId := fmt.Sprintf("%s/objectId/%s", *read.ID, d.Get("object_id"))
		if applicationId, ok := d.GetOk("application_id"); ok {
			resourceId = fmt.Sprintf(
				"%s/applicationId/%s",
				resourceId,
				applicationId)
		}
		d.SetId(resourceId)
	}

	return nil
}

func resourceArmKeyVaultAccessPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	return resourceArmKeyVaultAccessPolicyCreateOrDelete(d, meta, keyvault.Add)
}

func resourceArmKeyVaultAccessPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	// Lets lock the keyvault during a deletion to prevent a run updating the policies during a deletion
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

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure KeyVault %s: %+v", name, err)
	}

	if resp.Properties == nil {
		return fmt.Errorf("Properties not returned by api for Azurue KeyVault %s", name)
	}

	flattenedPolicy := keyVaultHelper.FlattenKeyVaultAccessPolicies(resp.Properties.AccessPolicies)

	policyObjectId := id.Path["objectId"]
	policyApplicationId := id.Path["applicationId"]
	policyId := getPolicyIdentity(&policyObjectId, &policyApplicationId)

	policy := findKeyVaultAccessPolicy(policyId, flattenedPolicy)

	if policy == nil {
		d.SetId("")
		return fmt.Errorf("Policy for was not found for vault: %s", name)
	}

	d.Set("vault_name", resp.Name)
	d.Set("resource_group_name", resGroup)
	d.Set("object_id", policyObjectId)
	d.Set("tenant_id", resp.Properties.TenantID.String())
	d.Set("application_id", policyApplicationId)
	d.Set("key_permissions", policy["key_permissions"])
	d.Set("secret_permissions", policy["secret_permissions"])
	d.Set("certificate_permissions", policy["certificate_permissions"])

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func getPolicyIdentity(objectId *string, applicationId *string) string {

	if applicationId != nil && *applicationId != "" {
		return fmt.Sprintf("%s/%s", *objectId, *applicationId)
	} else {
		return fmt.Sprintf("%s", *objectId)
	}
}

func matchAccessPolicy(policyString string, policy map[string]interface{}) bool {

	policyObjectId := policy["object_id"].(string)

	if policyApplicationId, ok := policy["application_id"]; ok {
		return policyString == fmt.Sprintf("%s/%s", policyObjectId, policyApplicationId)
	} else {
		return policyString == policyObjectId
	}
}

func findKeyVaultAccessPolicy(policyString string, policies []map[string]interface{}) map[string]interface{} {
	for _, policy := range policies {
		if matchAccessPolicy(policyString, policy) {
			return policy
		}
	}
	return nil
}

func expandKeyVaultAccessPolicy(d *schema.ResourceData) keyvault.AccessPolicyEntry {

	policy := keyvault.AccessPolicyEntry{
		Permissions: &keyvault.Permissions{
			Certificates: keyVaultHelper.ExpandKeyVaultAccessPolicyCertificatePermissions(d.Get("certificate_permissions").([]interface{})),
			Keys:         keyVaultHelper.ExpandKeyVaultAccessPolicyKeyPermissions(d.Get("key_permissions").([]interface{})),
			Secrets:      keyVaultHelper.ExpandKeyVaultAccessPolicySecretPermissions(d.Get("secret_permissions").([]interface{})),
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
