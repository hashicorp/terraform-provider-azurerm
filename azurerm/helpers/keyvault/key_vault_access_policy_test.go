package keyvault

import (
	"github.com/Azure/azure-sdk-for-go/services/keyvault/mgmt/2016-10-01/keyvault"
	"github.com/davecgh/go-spew/spew"
	uuid "github.com/satori/go.uuid"
	"reflect"
	"testing"
)

func TestFlattenKeyVaultAccessPolicies_noPolicies(t *testing.T) {
	policies := make([]keyvault.AccessPolicyEntry, 0, 0)
	flattenedPolicies := FlattenKeyVaultAccessPolicies(&policies)
	if len(flattenedPolicies) != 0 {
		t.Fatalf("Expected empty policy list, received non-empty policy list back")
	}
}

func TestFlattenKeyVaultAccessPolicies_oneEmptyPolicy(t *testing.T) {
	tenantId := uuid.NewV4()
	objectId := uuid.NewV4().String()
	applicationId := uuid.NewV4()

	policies := createPolicies(&tenantId, &objectId, &applicationId, nil)

	expectedPolicy := make(map[string]interface{})

	expectedPolicy["tenant_id"] = tenantId.String()
	expectedPolicy["object_id"] = objectId
	expectedPolicy["application_id"] = applicationId.String()
	expectedPolicy["key_permissions"] = make([]interface{}, 0, 0)
	expectedPolicy["secret_permissions"] = make([]interface{}, 0, 0)
	expectedPolicy["certificate_permissions"] = make([]interface{}, 0, 0)

	flattenedPolicies := FlattenKeyVaultAccessPolicies(policies)

	if !reflect.DeepEqual(flattenedPolicies[0], expectedPolicy) {
		t.Fatalf("Objects are not equal: %+v != %+v", spew.Sdump(flattenedPolicies), spew.Sdump(expectedPolicy))
	}
}

func TestFlattenKeyVaultAccessPolicies_oneFullPolicy(t *testing.T) {
	tenantId := uuid.NewV4()
	objectId := uuid.NewV4().String()
	applicationId := uuid.NewV4()

	certPerms := possibleCertificatePermissionsValues()
	keyPerms := possibleKeyPermissionsValues()
	secretPerms := possibleSecretPermissionsValues()

	permissions := &keyvault.Permissions{
		Certificates: &certPerms,
		Keys:         &keyPerms,
		Secrets:      &secretPerms,
	}

	policies := createPolicies(&tenantId, &objectId, &applicationId, permissions)

	expectedPolicy := make(map[string]interface{})

	expectedPolicy["tenant_id"] = tenantId.String()
	expectedPolicy["object_id"] = objectId
	expectedPolicy["application_id"] = applicationId.String()
	expectedPolicy["key_permissions"] = possibleKeyPermissionsStrings()
	expectedPolicy["secret_permissions"] = possibleSecretPermissionsStrings()
	expectedPolicy["certificate_permissions"] = possibleCertificatePermissionsStrings()
	flattenedPolicies := FlattenKeyVaultAccessPolicies(policies)

	if !reflect.DeepEqual(flattenedPolicies[0], expectedPolicy) {
		t.Fatalf("Objects are not equal: %+v != %+v", spew.Sdump(flattenedPolicies), spew.Sdump(expectedPolicy))
	}
}

func TestExpandKeyVaultAccessPolicyCertificatePermissions_empty(t *testing.T) {
	permissions := make([]interface{}, 0, 0)
	certPerms := ExpandKeyVaultAccessPolicyCertificatePermissions(permissions)

	certPermsKeyVault := make([]keyvault.CertificatePermissions, 0, 0)

	if !reflect.DeepEqual(*certPerms, certPermsKeyVault) {
		t.Fatalf("Objects are not equal: %+v != %+v", spew.Sdump(certPerms), spew.Sdump(certPermsKeyVault))
	}
}

func TestExpandKeyVaultAccessPolicyCertificatePermissions_full(t *testing.T) {
	permissions := possibleCertificatePermissionsStrings()

	certPerms := ExpandKeyVaultAccessPolicyCertificatePermissions(permissions)

	certPermsKeyVault := possibleCertificatePermissionsValues()

	if !reflect.DeepEqual(*certPerms, certPermsKeyVault) {
		t.Fatalf("Objects are not equal: %+v != %+v", spew.Sdump(certPerms), spew.Sdump(certPermsKeyVault))
	}
}

func TestExpandKeyVaultAccessPolicyKeyPermissions_empty(t *testing.T) {
	permissions := make([]interface{}, 0, 0)

	keyPerms := ExpandKeyVaultAccessPolicyKeyPermissions(permissions)

	keyPermsKeyVault := make([]keyvault.KeyPermissions, 0, 0)

	if !reflect.DeepEqual(*keyPerms, keyPermsKeyVault) {
		t.Fatalf("Objects are not equal: %+v != %+v", spew.Sdump(keyPerms), spew.Sdump(keyPermsKeyVault))
	}
}

func TestExpandKeyVaultAccessPolicyKeyPermissions_full(t *testing.T) {
	permissions := possibleKeyPermissionsStrings()

	keyPerms := ExpandKeyVaultAccessPolicyKeyPermissions(permissions)

	keyPermsKeyVault := possibleKeyPermissionsValues()

	if !reflect.DeepEqual(*keyPerms, keyPermsKeyVault) {
		t.Fatalf("Objects are not equal: %+v != %+v", spew.Sdump(keyPerms), spew.Sdump(keyPermsKeyVault))
	}
}

func TestExpandKeyVaultAccessPolicySecretPermissions_empty(t *testing.T) {
	permissions := make([]interface{}, 0, 0)

	secretPerms := ExpandKeyVaultAccessPolicySecretPermissions(permissions)

	secretPermsKeyVault := make([]keyvault.SecretPermissions, 0, 0)

	if !reflect.DeepEqual(*secretPerms, secretPermsKeyVault) {
		t.Fatalf("Objects are not equal: %+v != %+v", spew.Sdump(secretPerms), spew.Sdump(secretPermsKeyVault))
	}
}

func TestExpandKeyVaultAccessPolicySecretPermissions_full(t *testing.T) {
	permissions := possibleSecretPermissionsStrings()

	secretPerms := ExpandKeyVaultAccessPolicySecretPermissions(permissions)

	secretPermsKeyVault := possibleSecretPermissionsValues()

	if !reflect.DeepEqual(*secretPerms, secretPermsKeyVault) {
		t.Fatalf("Objects are not equal: %+v != %+v", spew.Sdump(secretPerms), spew.Sdump(secretPermsKeyVault))
	}
}

func createPolicies(tenantId *uuid.UUID, objectId *string, applicationId *uuid.UUID, permissions *keyvault.Permissions) *[]keyvault.AccessPolicyEntry {

	policies := make([]keyvault.AccessPolicyEntry, 0, 0)

	if permissions == nil {
		certPerms := make([]keyvault.CertificatePermissions, 0, 0)
		keyPerms := make([]keyvault.KeyPermissions, 0, 0)
		secretPerms := make([]keyvault.SecretPermissions, 0, 0)

		permissions = &keyvault.Permissions{
			Certificates: &certPerms,
			Keys:         &keyPerms,
			Secrets:      &secretPerms,
		}
	}

	if applicationId == nil {
		policy := keyvault.AccessPolicyEntry{
			TenantID:    tenantId,
			ObjectID:    objectId,
			Permissions: permissions,
		}
		policies = append(policies, policy)
	} else {
		policy := keyvault.AccessPolicyEntry{
			TenantID:      tenantId,
			ObjectID:      objectId,
			ApplicationID: applicationId,
			Permissions:   permissions,
		}
		policies = append(policies, policy)
	}

	return &policies
}

func possibleCertificatePermissionsStrings() []interface{} {

	var ret []interface{}

	ret = append(
		ret,
		"create",
		"delete",
		"deleteissuers",
		"get",
		"getissuers",
		"import",
		"list",
		"listissuers",
		"managecontacts",
		"manageissuers",
		"purge",
		"recover",
		"setissuers",
		"update",
	)

	return ret
}

func possibleCertificatePermissionsValues() []keyvault.CertificatePermissions {
	return []keyvault.CertificatePermissions{
		keyvault.Create,
		keyvault.Delete,
		keyvault.Deleteissuers,
		keyvault.Get,
		keyvault.Getissuers,
		keyvault.Import,
		keyvault.List,
		keyvault.Listissuers,
		keyvault.Managecontacts,
		keyvault.Manageissuers,
		keyvault.Purge,
		keyvault.Recover,
		keyvault.Setissuers,
		keyvault.Update,
	}
}

func possibleKeyPermissionsStrings() []interface{} {
	var ret []interface{}

	ret = append(
		ret,
		"backup",
		"create",
		"decrypt",
		"delete",
		"encrypt",
		"get",
		"import",
		"list",
		"purge",
		"recover",
		"restore",
		"sign",
		"unwrapKey",
		"update",
		"verify",
		"wrapKey",
	)

	return ret
}

func possibleKeyPermissionsValues() []keyvault.KeyPermissions {
	return []keyvault.KeyPermissions{
		keyvault.KeyPermissionsBackup,
		keyvault.KeyPermissionsCreate,
		keyvault.KeyPermissionsDecrypt,
		keyvault.KeyPermissionsDelete,
		keyvault.KeyPermissionsEncrypt,
		keyvault.KeyPermissionsGet,
		keyvault.KeyPermissionsImport,
		keyvault.KeyPermissionsList,
		keyvault.KeyPermissionsPurge,
		keyvault.KeyPermissionsRecover,
		keyvault.KeyPermissionsRestore,
		keyvault.KeyPermissionsSign,
		keyvault.KeyPermissionsUnwrapKey,
		keyvault.KeyPermissionsUpdate,
		keyvault.KeyPermissionsVerify,
		keyvault.KeyPermissionsWrapKey,
	}
}

func possibleSecretPermissionsStrings() []interface{} {
	var ret []interface{}

	ret = append(
		ret,
		"backup",
		"delete",
		"get",
		"list",
		"purge",
		"recover",
		"restore",
		"set",
	)

	return ret
}

func possibleSecretPermissionsValues() []keyvault.SecretPermissions {
	return []keyvault.SecretPermissions{
		keyvault.SecretPermissionsBackup,
		keyvault.SecretPermissionsDelete,
		keyvault.SecretPermissionsGet,
		keyvault.SecretPermissionsList,
		keyvault.SecretPermissionsPurge,
		keyvault.SecretPermissionsRecover,
		keyvault.SecretPermissionsRestore,
		keyvault.SecretPermissionsSet,
	}
}
