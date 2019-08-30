package azurerm

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	uuid "github.com/satori/go.uuid"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/mgmt/2018-02-14/keyvault"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
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
			"key_vault_id": {
				Type:          schema.TypeString,
				Optional:      true, //todo required in 2.0
				Computed:      true, //todo removed in 2.0
				ForceNew:      true,
				ValidateFunc:  azure.ValidateResourceID,
				ConflictsWith: []string{"vault_name"},
			},

			//todo remove in 2.0
			"vault_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				Deprecated:    "This property has been deprecated in favour of the key_vault_id property. This will prevent a class of bugs as described in https://github.com/terraform-providers/terraform-provider-azurerm/issues/2396 and will be removed in version 2.0 of the provider",
				ValidateFunc:  validate.NoEmptyStrings,
				ConflictsWith: []string{"key_vault_id"},
			},

			//todo remove in 2.0
			"resource_group_name": {
				Type:       schema.TypeString,
				Optional:   true,
				ForceNew:   true,
				Computed:   true,
				Deprecated: "This property has been deprecated as the resource group is now pulled from the vault ID and will be removed in version 2.0 of the provider",
				ValidateFunc: func(v interface{}, k string) (warnings []string, errors []error) {
					value := v.(string)

					if len(value) > 80 {
						errors = append(errors, fmt.Errorf("%q may not exceed 80 characters in length", k))
					}

					if strings.HasSuffix(value, ".") {
						errors = append(errors, fmt.Errorf("%q may not end with a period", k))
					}

					// regex pulled from https://docs.microsoft.com/en-us/rest/api/resources/resourcegroups/createorupdate
					if matched := regexp.MustCompile(`^[-\w\._\(\)]+$`).Match([]byte(value)); !matched {
						errors = append(errors, fmt.Errorf("%q may only contain alphanumeric characters, dash, underscores, parentheses and periods", k))
					}

					return warnings, errors
				},
			},

			"tenant_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.UUID,
			},

			"object_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.UUID,
			},

			"application_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validate.UUID,
			},

			"certificate_permissions": azure.SchemaKeyVaultCertificatePermissions(),

			"key_permissions": azure.SchemaKeyVaultKeyPermissions(),

			"secret_permissions": azure.SchemaKeyVaultSecretPermissions(),

			"storage_permissions": azure.SchemaKeyVaultStoragePermissions(),
		},
	}
}

func resourceArmKeyVaultAccessPolicyCreateOrDelete(d *schema.ResourceData, meta interface{}, action keyvault.AccessPolicyUpdateKind) error {
	client := meta.(*ArmClient).keyvault.VaultsClient
	ctx := meta.(*ArmClient).StopContext
	log.Printf("[INFO] Preparing arguments for Key Vault Access Policy: %s.", action)

	vaultId := d.Get("key_vault_id").(string)
	vaultName := d.Get("vault_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	tenantIdRaw := d.Get("tenant_id").(string)
	tenantId, err := uuid.FromString(tenantIdRaw)
	if err != nil {
		return fmt.Errorf("Error parsing Tenant ID %q as a UUID: %+v", tenantIdRaw, err)
	}

	applicationIdRaw := d.Get("application_id").(string)
	objectId := d.Get("object_id").(string)

	if vaultName == "" {
		if vaultId == "" {
			return fmt.Errorf("one of `key_vault_id` or `vault_name` must be set")
		}
		id, err2 := azure.ParseAzureResourceID(vaultId)
		if err2 != nil {
			return err2
		}

		resourceGroup = id.ResourceGroup

		vaultNameTemp, ok := id.Path["vaults"]
		if !ok {
			return fmt.Errorf("key_value_id does not contain `vaults`: %q", vaultId)
		}
		vaultName = vaultNameTemp

	} else if resourceGroup == "" {
		return fmt.Errorf("one of `resource_group_name` must be set when `vault_name` is used")
	}

	keyVault, err := client.Get(ctx, resourceGroup, vaultName)
	if err != nil {
		// If the key vault does not exist but this is not a new resource, the policy
		// which previously existed was deleted with the key vault, so reflect that in
		// state. If this is a new resource and key vault does not exist, it's likely
		// a bad ID was given.
		if utils.ResponseWasNotFound(keyVault.Response) && !d.IsNewResource() {
			log.Printf("[DEBUG] Parent Key Vault %q was not found in Resource Group %q - removing from state!", vaultName, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Key Vault %q (Resource Group %q): %+v", vaultName, resourceGroup, err)
	}

	// This is because azure doesn't have an 'id' for a keyvault access policy
	// In order to compensate for this and allow importing of this resource we are artificially
	// creating an identity for a key vault policy object
	resourceId := fmt.Sprintf("%s/objectId/%s", *keyVault.ID, objectId)
	if applicationIdRaw != "" {
		resourceId = fmt.Sprintf("%s/applicationId/%s", resourceId, applicationIdRaw)
	}

	// Locking to prevent parallel changes causing issues
	locks.ByName(vaultName, keyVaultResourceName)
	defer locks.UnlockByName(vaultName, keyVaultResourceName)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		props := keyVault.Properties
		if props == nil {
			return fmt.Errorf("Error parsing Key Vault: `properties` was nil")
		}

		if props.AccessPolicies == nil {
			return fmt.Errorf("Error parsing Key Vault: `properties.AccessPolicy` was nil")
		}

		for _, policy := range *props.AccessPolicies {
			if policy.TenantID == nil || policy.ObjectID == nil {
				continue
			}

			tenantIdMatches := policy.TenantID.String() == tenantIdRaw
			objectIdMatches := *policy.ObjectID == objectId

			appId := ""
			if a := policy.ApplicationID; a != nil {
				appId = a.String()
			}
			applicationIdMatches := appId == applicationIdRaw
			if tenantIdMatches && objectIdMatches && applicationIdMatches {
				return tf.ImportAsExistsError("azurerm_key_vault_access_policy", resourceId)
			}
		}
	}

	certPermissionsRaw := d.Get("certificate_permissions").([]interface{})
	certPermissions := azure.ExpandCertificatePermissions(certPermissionsRaw)

	keyPermissionsRaw := d.Get("key_permissions").([]interface{})
	keyPermissions := azure.ExpandKeyPermissions(keyPermissionsRaw)

	secretPermissionsRaw := d.Get("secret_permissions").([]interface{})
	secretPermissions := azure.ExpandSecretPermissions(secretPermissionsRaw)

	storagePermissionsRaw := d.Get("storage_permissions").([]interface{})
	storagePermissions := azure.ExpandStoragePermissions(storagePermissionsRaw)

	accessPolicy := keyvault.AccessPolicyEntry{
		ObjectID: utils.String(objectId),
		TenantID: &tenantId,
		Permissions: &keyvault.Permissions{
			Certificates: certPermissions,
			Keys:         keyPermissions,
			Secrets:      secretPermissions,
			Storage:      storagePermissions,
		},
	}

	if applicationIdRaw != "" {
		applicationId, err2 := uuid.FromString(applicationIdRaw)
		if err2 != nil {
			return fmt.Errorf("Error parsing Application ID %q as a UUID: %+v", applicationIdRaw, err2)
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

	if _, err = client.UpdateAccessPolicy(ctx, resourceGroup, vaultName, action, parameters); err != nil {
		return fmt.Errorf("Error updating Access Policy (Object ID %q / Application ID %q) for Key Vault %q (Resource Group %q): %+v", objectId, applicationIdRaw, vaultName, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, vaultName)
	if err != nil {
		return fmt.Errorf("Error retrieving Key Vault %q (Resource Group %q): %+v", vaultName, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read KeyVault %q (Resource Group %q) ID", vaultName, resourceGroup)
	}

	if d.IsNewResource() {
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
	client := meta.(*ArmClient).keyvault.VaultsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())

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

	d.Set("key_vault_id", resp.ID)
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

		storagePermissions := azure.FlattenStoragePermissions(permissions.Storage)
		if err := d.Set("storage_permissions", storagePermissions); err != nil {
			return fmt.Errorf("Error setting `storage_permissions`: %+v", err)
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
