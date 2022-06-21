package keyvault

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/mgmt/2021-10-01/keyvault"
	"github.com/gofrs/uuid"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceKeyVaultAccessPolicy() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceKeyVaultAccessPolicyCreate,
		Read:   resourceKeyVaultAccessPolicyRead,
		Update: resourceKeyVaultAccessPolicyUpdate,
		Delete: resourceKeyVaultAccessPolicyDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.AccessPolicyID(id)
			return err
		}),

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

	vaultId, err := parse.VaultID(d.Get("key_vault_id").(string))
	if err != nil {
		return err
	}

	tenantIdRaw := d.Get("tenant_id").(string)
	tenantId, err := uuid.FromString(tenantIdRaw)
	if err != nil {
		return fmt.Errorf("parsing Tenant ID %q as a UUID: %+v", tenantIdRaw, err)
	}

	objectId := d.Get("object_id").(string)
	applicationIdRaw := d.Get("application_id").(string)

	id := parse.NewAccessPolicyId(*vaultId, objectId, applicationIdRaw)

	keyVault, err := client.Get(ctx, vaultId.ResourceGroup, vaultId.Name)
	if err != nil {
		// If the key vault does not exist but this is not a new resource, the policy
		// which previously existed was deleted with the key vault, so reflect that in
		// state. If this is a new resource and key vault does not exist, it's likely
		// a bad ID was given.
		if utils.ResponseWasNotFound(keyVault.Response) && !d.IsNewResource() {
			log.Printf("[DEBUG] Parent %s was not found - removing from state!", *vaultId)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving parent %s: %+v", *vaultId, err)
	}

	// Locking to prevent parallel changes causing issues
	locks.ByName(vaultId.Name, keyVaultResourceName)
	defer locks.UnlockByName(vaultId.Name, keyVaultResourceName)

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
				return tf.ImportAsExistsError("azurerm_key_vault_access_policy", id.ID())
			}
		}
	}

	var accessPolicy keyvault.AccessPolicyEntry
	switch action {
	case keyvault.AccessPolicyUpdateKindRemove:
		// To remove a policy correctly, we need to send it with all permissions in the correct case which may have drifted
		// in config over time so we read it back from the vault by objectId
		resp, err := client.Get(ctx, vaultId.ResourceGroup, vaultId.Name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				log.Printf("[DEBUG] parent %s was not found - removing from state", *vaultId)
				d.SetId("")
				return nil
			}
			return fmt.Errorf("retrieving parent %s: %+v", vaultId, err)
		}

		if resp.Properties == nil || resp.Properties.AccessPolicies == nil {
			return fmt.Errorf("retrieving parent %s: `accessPolicies` was nil", *vaultId)
		}

		accessPolicyRaw := FindKeyVaultAccessPolicy(resp.Properties.AccessPolicies, objectId, applicationIdRaw)
		if accessPolicyRaw == nil {
			return fmt.Errorf("unable to find Access Policy (Object ID %q / Application ID %q) on %s", id.ObjectID(), id.ApplicationId(), *vaultId)
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
		Name: utils.String(vaultId.Name),
		Properties: &keyvault.VaultAccessPolicyProperties{
			AccessPolicies: &accessPolicies,
		},
	}

	if _, err = client.UpdateAccessPolicy(ctx, vaultId.ResourceGroup, vaultId.Name, action, parameters); err != nil {
		return fmt.Errorf("updating Access Policy (Object ID %q / Application ID %q) for %s: %+v", objectId, applicationIdRaw, *vaultId, err)
	}
	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"notfound", "vaultnotfound"},
		Target:                    []string{"found"},
		Refresh:                   accessPolicyRefreshFunc(ctx, client, vaultId.ResourceGroup, vaultId.Name, objectId, applicationIdRaw),
		Delay:                     5 * time.Second,
		ContinuousTargetOccurence: 3,
		Timeout:                   d.Timeout(pluginsdk.TimeoutCreate),
	}

	if action == keyvault.AccessPolicyUpdateKindRemove {
		stateConf.Target = []string{"notfound"}
		stateConf.Pending = []string{"found", "vaultnotfound"}
		stateConf.Timeout = d.Timeout(pluginsdk.TimeoutDelete)
	}

	if action == keyvault.AccessPolicyUpdateKindReplace {
		stateConf.Timeout = d.Timeout(pluginsdk.TimeoutUpdate)
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("failed waiting for Key Vault Access Policy (Object ID: %q) to apply: %+v", objectId, err)
	}

	if d.IsNewResource() {
		d.SetId(id.ID())
	}

	return nil
}

func resourceKeyVaultAccessPolicyCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	return resourceKeyVaultAccessPolicyCreateOrDelete(d, meta, keyvault.AccessPolicyUpdateKindAdd)
}

func resourceKeyVaultAccessPolicyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	return resourceKeyVaultAccessPolicyCreateOrDelete(d, meta, keyvault.AccessPolicyUpdateKindRemove)
}

func resourceKeyVaultAccessPolicyUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	return resourceKeyVaultAccessPolicyCreateOrDelete(d, meta, keyvault.AccessPolicyUpdateKindReplace)
}

func resourceKeyVaultAccessPolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).KeyVault.VaultsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AccessPolicyID(d.Id())
	if err != nil {
		return err
	}

	vaultId := id.KeyVaultId()

	resp, err := client.Get(ctx, vaultId.ResourceGroup, vaultId.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] parent %q was not found - removing from state", vaultId)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving parent %s: %+v", vaultId, err)
	}

	if resp.Properties == nil || resp.Properties.AccessPolicies == nil {
		return fmt.Errorf("retrieving parent %s: accessPolicies were nil", vaultId)
	}

	policy := FindKeyVaultAccessPolicy(resp.Properties.AccessPolicies, id.ObjectID(), id.ApplicationId())

	if policy == nil {
		log.Printf("[ERROR] Access Policy (Object ID %q / Application ID %q) was not found in %s - removing from state", id.ObjectID(), id.ApplicationId(), vaultId)
		d.SetId("")
		return nil
	}

	d.Set("key_vault_id", id.KeyVaultId().ID())
	d.Set("application_id", id.ApplicationId())
	d.Set("object_id", id.ObjectID())

	tenantId := ""
	if tid := policy.TenantID; tid != nil {
		tenantId = tid.String()
	}
	d.Set("tenant_id", tenantId)

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
