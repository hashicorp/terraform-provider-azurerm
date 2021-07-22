package keyvault

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/keyvault/mgmt/2020-04-01-preview/keyvault"
	"github.com/gofrs/uuid"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceKeyVaultAccessPolicy() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceKeyVaultAccessPolicyCreate,
		Read:   resourceKeyVaultAccessPolicyRead,
		Update: resourceKeyVaultAccessPolicyUpdate,
		Delete: resourceKeyVaultAccessPolicyDelete,
		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"key_vault_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"tenant_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"object_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"application_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"certificate_permissions": schemaCertificatePermissions(),

			"key_permissions": schemaKeyPermissions(),

			"secret_permissions": schemaSecretPermissions(),

			"storage_permissions": schemaStoragePermissions(),
		},
	}
}

func resourceKeyVaultAccessPolicyCreateOrDelete(d *pluginsdk.ResourceData, meta interface{}, action keyvault.AccessPolicyUpdateKind) error {
	client := meta.(*clients.Client).KeyVault.VaultsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] Preparing arguments for Key Vault Access Policy: %s.", action)

	vaultId := d.Get("key_vault_id").(string)

	tenantIdRaw := d.Get("tenant_id").(string)
	tenantId, err := uuid.FromString(tenantIdRaw)
	if err != nil {
		return fmt.Errorf("parsing Tenant ID %q as a UUID: %+v", tenantIdRaw, err)
	}

	applicationIdRaw := d.Get("application_id").(string)
	objectId := d.Get("object_id").(string)

	id, err := azure.ParseAzureResourceID(vaultId)
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	vaultName, ok := id.Path["vaults"]
	if !ok {
		return fmt.Errorf("key_value_id does not contain `vaults`: %q", vaultId)
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

		return fmt.Errorf("retrieving Key Vault %q (Resource Group %q): %+v", vaultName, resourceGroup, err)
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

	if d.IsNewResource() {
		props := keyVault.Properties
		if props == nil {
			return fmt.Errorf("parsing Key Vault: `properties` was nil")
		}

		if props.AccessPolicies == nil {
			return fmt.Errorf("parsing Key Vault: `properties.AccessPolicy` was nil")
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

	var accessPolicy keyvault.AccessPolicyEntry
	switch action {
	case keyvault.Remove:
		// To remove a policy correctly, we need to send it with all permissions in the correct case which may have drifted
		// in config over time so we read it back from the vault by objectId
		resp, err := client.Get(ctx, id.ResourceGroup, vaultName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				log.Printf("[ERROR] Key Vault %q (Resource Group %q) was not found - removing from state", vaultName, id.ResourceGroup)
				d.SetId("")
				return nil
			}
			return fmt.Errorf("making Read request on Azure KeyVault %q (Resource Group %q): %+v", vaultName, id.ResourceGroup, err)
		}

		if resp.Properties == nil || resp.Properties.AccessPolicies == nil {
			return fmt.Errorf("failed reading Access Policies for %q (resource group %q)", vaultName, id.ResourceGroup)
		}

		accessPolicyRaw := FindKeyVaultAccessPolicy(resp.Properties.AccessPolicies, objectId, applicationIdRaw)
		if accessPolicyRaw == nil {
			return fmt.Errorf("failed finding this specific Access Policy on Azure KeyVault %q (resource group %q)", vaultName, id.ResourceGroup)
		}
		accessPolicy = *accessPolicyRaw

	default:
		certPermissionsRaw := d.Get("certificate_permissions").([]interface{})
		certPermissions := expandCertificatePermissions(certPermissionsRaw)

		keyPermissionsRaw := d.Get("key_permissions").([]interface{})
		keyPermissions := expandKeyPermissions(keyPermissionsRaw)

		secretPermissionsRaw := d.Get("secret_permissions").([]interface{})
		secretPermissions := expandSecretPermissions(secretPermissionsRaw)

		storagePermissionsRaw := d.Get("storage_permissions").([]interface{})
		storagePermissions := expandStoragePermissions(storagePermissionsRaw)

		accessPolicy = keyvault.AccessPolicyEntry{
			ObjectID: utils.String(objectId),
			TenantID: &tenantId,
			Permissions: &keyvault.Permissions{
				Certificates: certPermissions,
				Keys:         keyPermissions,
				Secrets:      secretPermissions,
				Storage:      storagePermissions,
			},
		}
	}

	if applicationIdRaw != "" {
		applicationId, err2 := uuid.FromString(applicationIdRaw)
		if err2 != nil {
			return fmt.Errorf("parsing Application ID %q as a UUID: %+v", applicationIdRaw, err2)
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
		return fmt.Errorf("updating Access Policy (Object ID %q / Application ID %q) for Key Vault %q (Resource Group %q): %+v", objectId, applicationIdRaw, vaultName, resourceGroup, err)
	}
	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"notfound", "vaultnotfound"},
		Target:                    []string{"found"},
		Refresh:                   accessPolicyRefreshFunc(ctx, client, resourceGroup, vaultName, objectId, applicationIdRaw),
		Delay:                     5 * time.Second,
		ContinuousTargetOccurence: 3,
		Timeout:                   d.Timeout(pluginsdk.TimeoutCreate),
	}

	if action == keyvault.Remove {
		stateConf.Target = []string{"notfound"}
		stateConf.Pending = []string{"found", "vaultnotfound"}
		stateConf.Timeout = d.Timeout(pluginsdk.TimeoutDelete)
	}

	if action == keyvault.Replace {
		stateConf.Timeout = d.Timeout(pluginsdk.TimeoutUpdate)
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("failed waiting for Key Vault Access Policy (Object ID: %q) to apply: %+v", objectId, err)
	}

	read, err := client.Get(ctx, resourceGroup, vaultName)
	if err != nil {
		return fmt.Errorf("retrieving Key Vault %q (Resource Group %q): %+v", vaultName, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read KeyVault %q (Resource Group %q) ID", vaultName, resourceGroup)
	}

	if d.IsNewResource() {
		d.SetId(resourceId)
	}

	return nil
}

func resourceKeyVaultAccessPolicyCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	return resourceKeyVaultAccessPolicyCreateOrDelete(d, meta, keyvault.Add)
}

func resourceKeyVaultAccessPolicyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	return resourceKeyVaultAccessPolicyCreateOrDelete(d, meta, keyvault.Remove)
}

func resourceKeyVaultAccessPolicyUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	return resourceKeyVaultAccessPolicyCreateOrDelete(d, meta, keyvault.Replace)
}

func resourceKeyVaultAccessPolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).KeyVault.VaultsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

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

		return fmt.Errorf("making Read request on Azure KeyVault %q (Resource Group %q): %+v", vaultName, resGroup, err)
	}

	if resp.Properties == nil || resp.Properties.AccessPolicies == nil {
		return fmt.Errorf("failed reading Access Policies for %q (resource group %q)", vaultName, id.ResourceGroup)
	}

	policy := FindKeyVaultAccessPolicy(resp.Properties.AccessPolicies, objectId, applicationId)

	if policy == nil {
		log.Printf("[ERROR] Access Policy (Object ID %q / Application ID %q) was not found in Key Vault %q (Resource Group %q) - removing from state", objectId, applicationId, vaultName, resGroup)
		d.SetId("")
		return nil
	}

	d.Set("key_vault_id", resp.ID)
	d.Set("object_id", objectId)

	if tid := policy.TenantID; tid != nil {
		d.Set("tenant_id", tid.String())
	}

	if aid := policy.ApplicationID; aid != nil {
		d.Set("application_id", aid.String())
	}

	if permissions := policy.Permissions; permissions != nil {
		certificatePermissions := flattenCertificatePermissions(permissions.Certificates)
		if err := d.Set("certificate_permissions", certificatePermissions); err != nil {
			return fmt.Errorf("setting `certificate_permissions`: %+v", err)
		}

		keyPermissions := flattenKeyPermissions(permissions.Keys)
		if err := d.Set("key_permissions", keyPermissions); err != nil {
			return fmt.Errorf("setting `key_permissions`: %+v", err)
		}

		secretPermissions := flattenSecretPermissions(permissions.Secrets)
		if err := d.Set("secret_permissions", secretPermissions); err != nil {
			return fmt.Errorf("setting `secret_permissions`: %+v", err)
		}

		storagePermissions := flattenStoragePermissions(permissions.Storage)
		if err := d.Set("storage_permissions", storagePermissions); err != nil {
			return fmt.Errorf("setting `storage_permissions`: %+v", err)
		}
	}

	return nil
}

func FindKeyVaultAccessPolicy(policies *[]keyvault.AccessPolicyEntry, objectId string, applicationId string) *keyvault.AccessPolicyEntry {
	if policies == nil {
		return nil
	}

	for _, policy := range *policies {
		if id := policy.ObjectID; id != nil {
			if strings.EqualFold(*id, objectId) {
				aid := ""
				if policy.ApplicationID != nil {
					aid = policy.ApplicationID.String()
				}

				if strings.EqualFold(aid, applicationId) {
					return &policy
				}
			}
		}
	}

	return nil
}

func accessPolicyRefreshFunc(ctx context.Context, client *keyvault.VaultsClient, resourceGroup string, vaultName string, objectId string, applicationId string) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Checking for completion of Access Policy create/update")

		read, err := client.Get(ctx, resourceGroup, vaultName)
		if err != nil {
			if utils.ResponseWasNotFound(read.Response) {
				return "vaultnotfound", "vaultnotfound", fmt.Errorf("failed to find vault %q (resource group %q)", vaultName, resourceGroup)
			}
		}

		if read.Properties != nil && read.Properties.AccessPolicies != nil {
			policy := FindKeyVaultAccessPolicy(read.Properties.AccessPolicies, objectId, applicationId)
			if policy != nil {
				return "found", "found", nil
			}
		}

		return "notfound", "notfound", nil
	}
}
