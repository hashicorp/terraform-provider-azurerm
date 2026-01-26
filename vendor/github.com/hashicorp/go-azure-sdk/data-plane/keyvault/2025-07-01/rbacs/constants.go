package rbacs

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataAction string

const (
	DataActionMicrosoftPointKeyVaultManagedHsmBackupStartAction             DataAction = "Microsoft.KeyVault/managedHsm/backup/start/action"
	DataActionMicrosoftPointKeyVaultManagedHsmBackupStatusAction            DataAction = "Microsoft.KeyVault/managedHsm/backup/status/action"
	DataActionMicrosoftPointKeyVaultManagedHsmKeysBackupAction              DataAction = "Microsoft.KeyVault/managedHsm/keys/backup/action"
	DataActionMicrosoftPointKeyVaultManagedHsmKeysCreate                    DataAction = "Microsoft.KeyVault/managedHsm/keys/create"
	DataActionMicrosoftPointKeyVaultManagedHsmKeysDecryptAction             DataAction = "Microsoft.KeyVault/managedHsm/keys/decrypt/action"
	DataActionMicrosoftPointKeyVaultManagedHsmKeysDelete                    DataAction = "Microsoft.KeyVault/managedHsm/keys/delete"
	DataActionMicrosoftPointKeyVaultManagedHsmKeysDeletedKeysDelete         DataAction = "Microsoft.KeyVault/managedHsm/keys/deletedKeys/delete"
	DataActionMicrosoftPointKeyVaultManagedHsmKeysDeletedKeysReadAction     DataAction = "Microsoft.KeyVault/managedHsm/keys/deletedKeys/read/action"
	DataActionMicrosoftPointKeyVaultManagedHsmKeysDeletedKeysRecoverAction  DataAction = "Microsoft.KeyVault/managedHsm/keys/deletedKeys/recover/action"
	DataActionMicrosoftPointKeyVaultManagedHsmKeysEncryptAction             DataAction = "Microsoft.KeyVault/managedHsm/keys/encrypt/action"
	DataActionMicrosoftPointKeyVaultManagedHsmKeysExportAction              DataAction = "Microsoft.KeyVault/managedHsm/keys/export/action"
	DataActionMicrosoftPointKeyVaultManagedHsmKeysImportAction              DataAction = "Microsoft.KeyVault/managedHsm/keys/import/action"
	DataActionMicrosoftPointKeyVaultManagedHsmKeysReadAction                DataAction = "Microsoft.KeyVault/managedHsm/keys/read/action"
	DataActionMicrosoftPointKeyVaultManagedHsmKeysReleaseAction             DataAction = "Microsoft.KeyVault/managedHsm/keys/release/action"
	DataActionMicrosoftPointKeyVaultManagedHsmKeysRestoreAction             DataAction = "Microsoft.KeyVault/managedHsm/keys/restore/action"
	DataActionMicrosoftPointKeyVaultManagedHsmKeysSignAction                DataAction = "Microsoft.KeyVault/managedHsm/keys/sign/action"
	DataActionMicrosoftPointKeyVaultManagedHsmKeysUnwrapAction              DataAction = "Microsoft.KeyVault/managedHsm/keys/unwrap/action"
	DataActionMicrosoftPointKeyVaultManagedHsmKeysVerifyAction              DataAction = "Microsoft.KeyVault/managedHsm/keys/verify/action"
	DataActionMicrosoftPointKeyVaultManagedHsmKeysWrapAction                DataAction = "Microsoft.KeyVault/managedHsm/keys/wrap/action"
	DataActionMicrosoftPointKeyVaultManagedHsmKeysWriteAction               DataAction = "Microsoft.KeyVault/managedHsm/keys/write/action"
	DataActionMicrosoftPointKeyVaultManagedHsmRestoreStartAction            DataAction = "Microsoft.KeyVault/managedHsm/restore/start/action"
	DataActionMicrosoftPointKeyVaultManagedHsmRestoreStatusAction           DataAction = "Microsoft.KeyVault/managedHsm/restore/status/action"
	DataActionMicrosoftPointKeyVaultManagedHsmRngAction                     DataAction = "Microsoft.KeyVault/managedHsm/rng/action"
	DataActionMicrosoftPointKeyVaultManagedHsmRoleAssignmentsDeleteAction   DataAction = "Microsoft.KeyVault/managedHsm/roleAssignments/delete/action"
	DataActionMicrosoftPointKeyVaultManagedHsmRoleAssignmentsReadAction     DataAction = "Microsoft.KeyVault/managedHsm/roleAssignments/read/action"
	DataActionMicrosoftPointKeyVaultManagedHsmRoleAssignmentsWriteAction    DataAction = "Microsoft.KeyVault/managedHsm/roleAssignments/write/action"
	DataActionMicrosoftPointKeyVaultManagedHsmRoleDefinitionsDeleteAction   DataAction = "Microsoft.KeyVault/managedHsm/roleDefinitions/delete/action"
	DataActionMicrosoftPointKeyVaultManagedHsmRoleDefinitionsReadAction     DataAction = "Microsoft.KeyVault/managedHsm/roleDefinitions/read/action"
	DataActionMicrosoftPointKeyVaultManagedHsmRoleDefinitionsWriteAction    DataAction = "Microsoft.KeyVault/managedHsm/roleDefinitions/write/action"
	DataActionMicrosoftPointKeyVaultManagedHsmSecuritydomainDownloadAction  DataAction = "Microsoft.KeyVault/managedHsm/securitydomain/download/action"
	DataActionMicrosoftPointKeyVaultManagedHsmSecuritydomainDownloadRead    DataAction = "Microsoft.KeyVault/managedHsm/securitydomain/download/read"
	DataActionMicrosoftPointKeyVaultManagedHsmSecuritydomainTransferkeyRead DataAction = "Microsoft.KeyVault/managedHsm/securitydomain/transferkey/read"
	DataActionMicrosoftPointKeyVaultManagedHsmSecuritydomainUploadAction    DataAction = "Microsoft.KeyVault/managedHsm/securitydomain/upload/action"
	DataActionMicrosoftPointKeyVaultManagedHsmSecuritydomainUploadRead      DataAction = "Microsoft.KeyVault/managedHsm/securitydomain/upload/read"
)

func PossibleValuesForDataAction() []string {
	return []string{
		string(DataActionMicrosoftPointKeyVaultManagedHsmBackupStartAction),
		string(DataActionMicrosoftPointKeyVaultManagedHsmBackupStatusAction),
		string(DataActionMicrosoftPointKeyVaultManagedHsmKeysBackupAction),
		string(DataActionMicrosoftPointKeyVaultManagedHsmKeysCreate),
		string(DataActionMicrosoftPointKeyVaultManagedHsmKeysDecryptAction),
		string(DataActionMicrosoftPointKeyVaultManagedHsmKeysDelete),
		string(DataActionMicrosoftPointKeyVaultManagedHsmKeysDeletedKeysDelete),
		string(DataActionMicrosoftPointKeyVaultManagedHsmKeysDeletedKeysReadAction),
		string(DataActionMicrosoftPointKeyVaultManagedHsmKeysDeletedKeysRecoverAction),
		string(DataActionMicrosoftPointKeyVaultManagedHsmKeysEncryptAction),
		string(DataActionMicrosoftPointKeyVaultManagedHsmKeysExportAction),
		string(DataActionMicrosoftPointKeyVaultManagedHsmKeysImportAction),
		string(DataActionMicrosoftPointKeyVaultManagedHsmKeysReadAction),
		string(DataActionMicrosoftPointKeyVaultManagedHsmKeysReleaseAction),
		string(DataActionMicrosoftPointKeyVaultManagedHsmKeysRestoreAction),
		string(DataActionMicrosoftPointKeyVaultManagedHsmKeysSignAction),
		string(DataActionMicrosoftPointKeyVaultManagedHsmKeysUnwrapAction),
		string(DataActionMicrosoftPointKeyVaultManagedHsmKeysVerifyAction),
		string(DataActionMicrosoftPointKeyVaultManagedHsmKeysWrapAction),
		string(DataActionMicrosoftPointKeyVaultManagedHsmKeysWriteAction),
		string(DataActionMicrosoftPointKeyVaultManagedHsmRestoreStartAction),
		string(DataActionMicrosoftPointKeyVaultManagedHsmRestoreStatusAction),
		string(DataActionMicrosoftPointKeyVaultManagedHsmRngAction),
		string(DataActionMicrosoftPointKeyVaultManagedHsmRoleAssignmentsDeleteAction),
		string(DataActionMicrosoftPointKeyVaultManagedHsmRoleAssignmentsReadAction),
		string(DataActionMicrosoftPointKeyVaultManagedHsmRoleAssignmentsWriteAction),
		string(DataActionMicrosoftPointKeyVaultManagedHsmRoleDefinitionsDeleteAction),
		string(DataActionMicrosoftPointKeyVaultManagedHsmRoleDefinitionsReadAction),
		string(DataActionMicrosoftPointKeyVaultManagedHsmRoleDefinitionsWriteAction),
		string(DataActionMicrosoftPointKeyVaultManagedHsmSecuritydomainDownloadAction),
		string(DataActionMicrosoftPointKeyVaultManagedHsmSecuritydomainDownloadRead),
		string(DataActionMicrosoftPointKeyVaultManagedHsmSecuritydomainTransferkeyRead),
		string(DataActionMicrosoftPointKeyVaultManagedHsmSecuritydomainUploadAction),
		string(DataActionMicrosoftPointKeyVaultManagedHsmSecuritydomainUploadRead),
	}
}

func (s *DataAction) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDataAction(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDataAction(input string) (*DataAction, error) {
	vals := map[string]DataAction{
		"microsoft.keyvault/managedhsm/backup/start/action":             DataActionMicrosoftPointKeyVaultManagedHsmBackupStartAction,
		"microsoft.keyvault/managedhsm/backup/status/action":            DataActionMicrosoftPointKeyVaultManagedHsmBackupStatusAction,
		"microsoft.keyvault/managedhsm/keys/backup/action":              DataActionMicrosoftPointKeyVaultManagedHsmKeysBackupAction,
		"microsoft.keyvault/managedhsm/keys/create":                     DataActionMicrosoftPointKeyVaultManagedHsmKeysCreate,
		"microsoft.keyvault/managedhsm/keys/decrypt/action":             DataActionMicrosoftPointKeyVaultManagedHsmKeysDecryptAction,
		"microsoft.keyvault/managedhsm/keys/delete":                     DataActionMicrosoftPointKeyVaultManagedHsmKeysDelete,
		"microsoft.keyvault/managedhsm/keys/deletedkeys/delete":         DataActionMicrosoftPointKeyVaultManagedHsmKeysDeletedKeysDelete,
		"microsoft.keyvault/managedhsm/keys/deletedkeys/read/action":    DataActionMicrosoftPointKeyVaultManagedHsmKeysDeletedKeysReadAction,
		"microsoft.keyvault/managedhsm/keys/deletedkeys/recover/action": DataActionMicrosoftPointKeyVaultManagedHsmKeysDeletedKeysRecoverAction,
		"microsoft.keyvault/managedhsm/keys/encrypt/action":             DataActionMicrosoftPointKeyVaultManagedHsmKeysEncryptAction,
		"microsoft.keyvault/managedhsm/keys/export/action":              DataActionMicrosoftPointKeyVaultManagedHsmKeysExportAction,
		"microsoft.keyvault/managedhsm/keys/import/action":              DataActionMicrosoftPointKeyVaultManagedHsmKeysImportAction,
		"microsoft.keyvault/managedhsm/keys/read/action":                DataActionMicrosoftPointKeyVaultManagedHsmKeysReadAction,
		"microsoft.keyvault/managedhsm/keys/release/action":             DataActionMicrosoftPointKeyVaultManagedHsmKeysReleaseAction,
		"microsoft.keyvault/managedhsm/keys/restore/action":             DataActionMicrosoftPointKeyVaultManagedHsmKeysRestoreAction,
		"microsoft.keyvault/managedhsm/keys/sign/action":                DataActionMicrosoftPointKeyVaultManagedHsmKeysSignAction,
		"microsoft.keyvault/managedhsm/keys/unwrap/action":              DataActionMicrosoftPointKeyVaultManagedHsmKeysUnwrapAction,
		"microsoft.keyvault/managedhsm/keys/verify/action":              DataActionMicrosoftPointKeyVaultManagedHsmKeysVerifyAction,
		"microsoft.keyvault/managedhsm/keys/wrap/action":                DataActionMicrosoftPointKeyVaultManagedHsmKeysWrapAction,
		"microsoft.keyvault/managedhsm/keys/write/action":               DataActionMicrosoftPointKeyVaultManagedHsmKeysWriteAction,
		"microsoft.keyvault/managedhsm/restore/start/action":            DataActionMicrosoftPointKeyVaultManagedHsmRestoreStartAction,
		"microsoft.keyvault/managedhsm/restore/status/action":           DataActionMicrosoftPointKeyVaultManagedHsmRestoreStatusAction,
		"microsoft.keyvault/managedhsm/rng/action":                      DataActionMicrosoftPointKeyVaultManagedHsmRngAction,
		"microsoft.keyvault/managedhsm/roleassignments/delete/action":   DataActionMicrosoftPointKeyVaultManagedHsmRoleAssignmentsDeleteAction,
		"microsoft.keyvault/managedhsm/roleassignments/read/action":     DataActionMicrosoftPointKeyVaultManagedHsmRoleAssignmentsReadAction,
		"microsoft.keyvault/managedhsm/roleassignments/write/action":    DataActionMicrosoftPointKeyVaultManagedHsmRoleAssignmentsWriteAction,
		"microsoft.keyvault/managedhsm/roledefinitions/delete/action":   DataActionMicrosoftPointKeyVaultManagedHsmRoleDefinitionsDeleteAction,
		"microsoft.keyvault/managedhsm/roledefinitions/read/action":     DataActionMicrosoftPointKeyVaultManagedHsmRoleDefinitionsReadAction,
		"microsoft.keyvault/managedhsm/roledefinitions/write/action":    DataActionMicrosoftPointKeyVaultManagedHsmRoleDefinitionsWriteAction,
		"microsoft.keyvault/managedhsm/securitydomain/download/action":  DataActionMicrosoftPointKeyVaultManagedHsmSecuritydomainDownloadAction,
		"microsoft.keyvault/managedhsm/securitydomain/download/read":    DataActionMicrosoftPointKeyVaultManagedHsmSecuritydomainDownloadRead,
		"microsoft.keyvault/managedhsm/securitydomain/transferkey/read": DataActionMicrosoftPointKeyVaultManagedHsmSecuritydomainTransferkeyRead,
		"microsoft.keyvault/managedhsm/securitydomain/upload/action":    DataActionMicrosoftPointKeyVaultManagedHsmSecuritydomainUploadAction,
		"microsoft.keyvault/managedhsm/securitydomain/upload/read":      DataActionMicrosoftPointKeyVaultManagedHsmSecuritydomainUploadRead,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DataAction(input)
	return &out, nil
}

type RoleDefinitionType string

const (
	RoleDefinitionTypeMicrosoftPointAuthorizationRoleDefinitions RoleDefinitionType = "Microsoft.Authorization/roleDefinitions"
)

func PossibleValuesForRoleDefinitionType() []string {
	return []string{
		string(RoleDefinitionTypeMicrosoftPointAuthorizationRoleDefinitions),
	}
}

func (s *RoleDefinitionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRoleDefinitionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRoleDefinitionType(input string) (*RoleDefinitionType, error) {
	vals := map[string]RoleDefinitionType{
		"microsoft.authorization/roledefinitions": RoleDefinitionTypeMicrosoftPointAuthorizationRoleDefinitions,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RoleDefinitionType(input)
	return &out, nil
}

type RoleScope string

const (
	RoleScopeKeys  RoleScope = "/keys"
	RoleScopeSlash RoleScope = "/"
)

func PossibleValuesForRoleScope() []string {
	return []string{
		string(RoleScopeKeys),
		string(RoleScopeSlash),
	}
}

func (s *RoleScope) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRoleScope(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRoleScope(input string) (*RoleScope, error) {
	vals := map[string]RoleScope{
		"/keys": RoleScopeKeys,
		"/":     RoleScopeSlash,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RoleScope(input)
	return &out, nil
}

type RoleType string

const (
	RoleTypeAKVBuiltInRole RoleType = "AKVBuiltInRole"
	RoleTypeCustomRole     RoleType = "CustomRole"
)

func PossibleValuesForRoleType() []string {
	return []string{
		string(RoleTypeAKVBuiltInRole),
		string(RoleTypeCustomRole),
	}
}

func (s *RoleType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRoleType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRoleType(input string) (*RoleType, error) {
	vals := map[string]RoleType{
		"akvbuiltinrole": RoleTypeAKVBuiltInRole,
		"customrole":     RoleTypeCustomRole,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RoleType(input)
	return &out, nil
}
