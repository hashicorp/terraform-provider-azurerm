package keyvault

import (
	"fmt"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/mgmt/2018-02-14/keyvault"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func resourceAzureRMKeyVaultMigrateState(v int, is *terraform.InstanceState, _ interface{}) (*terraform.InstanceState, error) {
	switch v {
	case 0:
		log.Println("[INFO] Found AzureRM Key Vault State v0; migrating to v1")
		return migrateAzureRMKeyVaultStateV0toV1(is)
	default:
		return is, fmt.Errorf("Unexpected schema version: %d", v)
	}
}

func migrateAzureRMKeyVaultStateV0toV1(is *terraform.InstanceState) (*terraform.InstanceState, error) {
	if is.Empty() {
		log.Println("[DEBUG] Empty InstanceState; nothing to migrate.")
		return is, nil
	}

	log.Printf("[DEBUG] ARM Key Vault Attributes before Migration: %#v", is.Attributes)

	if err := migrateAzureRMKeyVaultStateV0toV1AccessPolicies(is); err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] ARM Key Vault Attributes after State Migration: %#v", is.Attributes)

	return is, nil
}

func migrateAzureRMKeyVaultStateV0toV1AccessPolicies(is *terraform.InstanceState) error {
	keyVaultSchema := resourceArmKeyVault().Schema
	reader := &schema.MapFieldReader{
		Schema: keyVaultSchema,
		Map:    schema.BasicMapReader(is.Attributes),
	}

	// parse and update the existing data
	result, err := reader.ReadField([]string{"access_policy"})
	if err != nil {
		return err
	}

	inputAccessPolicies := result.Value.([]interface{})

	if len(inputAccessPolicies) == 0 {
		return nil
	}

	outputAccessPolicies := make([]interface{}, 0)
	for _, accessPolicy := range inputAccessPolicies {
		policy := accessPolicy.(map[string]interface{})

		if v, ok := policy["certificate_permissions"]; ok {
			inputCertificatePermissions := v.([]interface{})
			outputCertificatePermissions := make([]string, 0)
			for _, p := range inputCertificatePermissions {
				permission := p.(string)
				if strings.ToLower(permission) == "all" {
					outputCertificatePermissions = append(outputCertificatePermissions, string(keyvault.Create))
					outputCertificatePermissions = append(outputCertificatePermissions, string(keyvault.Delete))
					outputCertificatePermissions = append(outputCertificatePermissions, string(keyvault.Deleteissuers))
					outputCertificatePermissions = append(outputCertificatePermissions, string(keyvault.Get))
					outputCertificatePermissions = append(outputCertificatePermissions, string(keyvault.Getissuers))
					outputCertificatePermissions = append(outputCertificatePermissions, string(keyvault.Import))
					outputCertificatePermissions = append(outputCertificatePermissions, string(keyvault.List))
					outputCertificatePermissions = append(outputCertificatePermissions, string(keyvault.Listissuers))
					outputCertificatePermissions = append(outputCertificatePermissions, string(keyvault.Managecontacts))
					outputCertificatePermissions = append(outputCertificatePermissions, string(keyvault.Manageissuers))
					outputCertificatePermissions = append(outputCertificatePermissions, string(keyvault.Setissuers))
					outputCertificatePermissions = append(outputCertificatePermissions, string(keyvault.Update))
					break
				}
			}

			if len(outputCertificatePermissions) > 0 {
				policy["certificate_permissions"] = outputCertificatePermissions
			}
		}

		if v, ok := policy["key_permissions"]; ok {
			inputKeyPermissions := v.([]interface{})
			outputKeyPermissions := make([]string, 0)
			for _, p := range inputKeyPermissions {
				permission := p.(string)
				if strings.ToLower(permission) == "all" {
					outputKeyPermissions = append(outputKeyPermissions, string(keyvault.KeyPermissionsBackup))
					outputKeyPermissions = append(outputKeyPermissions, string(keyvault.KeyPermissionsCreate))
					outputKeyPermissions = append(outputKeyPermissions, string(keyvault.KeyPermissionsDecrypt))
					outputKeyPermissions = append(outputKeyPermissions, string(keyvault.KeyPermissionsDelete))
					outputKeyPermissions = append(outputKeyPermissions, string(keyvault.KeyPermissionsEncrypt))
					outputKeyPermissions = append(outputKeyPermissions, string(keyvault.KeyPermissionsGet))
					outputKeyPermissions = append(outputKeyPermissions, string(keyvault.KeyPermissionsImport))
					outputKeyPermissions = append(outputKeyPermissions, string(keyvault.KeyPermissionsList))
					outputKeyPermissions = append(outputKeyPermissions, string(keyvault.KeyPermissionsPurge))
					outputKeyPermissions = append(outputKeyPermissions, string(keyvault.KeyPermissionsRecover))
					outputKeyPermissions = append(outputKeyPermissions, string(keyvault.KeyPermissionsRestore))
					outputKeyPermissions = append(outputKeyPermissions, string(keyvault.KeyPermissionsSign))
					outputKeyPermissions = append(outputKeyPermissions, string(keyvault.KeyPermissionsUnwrapKey))
					outputKeyPermissions = append(outputKeyPermissions, string(keyvault.KeyPermissionsUpdate))
					outputKeyPermissions = append(outputKeyPermissions, string(keyvault.KeyPermissionsVerify))
					outputKeyPermissions = append(outputKeyPermissions, string(keyvault.KeyPermissionsWrapKey))
					break
				}
			}

			if len(outputKeyPermissions) > 0 {
				policy["key_permissions"] = outputKeyPermissions
			}
		}

		if v, ok := policy["secret_permissions"]; ok {
			inputSecretPermissions := v.([]interface{})
			outputSecretPermissions := make([]string, 0)
			for _, p := range inputSecretPermissions {
				permission := p.(string)
				if strings.ToLower(permission) == "all" {
					outputSecretPermissions = append(outputSecretPermissions, string(keyvault.SecretPermissionsBackup))
					outputSecretPermissions = append(outputSecretPermissions, string(keyvault.SecretPermissionsDelete))
					outputSecretPermissions = append(outputSecretPermissions, string(keyvault.SecretPermissionsGet))
					outputSecretPermissions = append(outputSecretPermissions, string(keyvault.SecretPermissionsList))
					outputSecretPermissions = append(outputSecretPermissions, string(keyvault.SecretPermissionsPurge))
					outputSecretPermissions = append(outputSecretPermissions, string(keyvault.SecretPermissionsRecover))
					outputSecretPermissions = append(outputSecretPermissions, string(keyvault.SecretPermissionsRestore))
					outputSecretPermissions = append(outputSecretPermissions, string(keyvault.SecretPermissionsSet))
					break
				}
			}

			if len(outputSecretPermissions) > 0 {
				policy["secret_permissions"] = outputSecretPermissions
			}
		}

		outputAccessPolicies = append(outputAccessPolicies, policy)
	}

	// remove the existing fields
	for k := range is.Attributes {
		if strings.HasPrefix(k, "access_policy.") {
			delete(is.Attributes, k)
		}
	}

	// write this out
	writer := schema.MapFieldWriter{
		Schema: keyVaultSchema,
	}
	if err := writer.WriteField([]string{"access_policy"}, outputAccessPolicies); err != nil {
		return err
	}
	for k, v := range writer.Map() {
		is.Attributes[k] = v
	}

	return nil
}
