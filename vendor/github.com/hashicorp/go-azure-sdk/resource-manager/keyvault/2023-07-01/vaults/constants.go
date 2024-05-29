package vaults

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccessPolicyUpdateKind string

const (
	AccessPolicyUpdateKindAdd     AccessPolicyUpdateKind = "add"
	AccessPolicyUpdateKindRemove  AccessPolicyUpdateKind = "remove"
	AccessPolicyUpdateKindReplace AccessPolicyUpdateKind = "replace"
)

func PossibleValuesForAccessPolicyUpdateKind() []string {
	return []string{
		string(AccessPolicyUpdateKindAdd),
		string(AccessPolicyUpdateKindRemove),
		string(AccessPolicyUpdateKindReplace),
	}
}

func (s *AccessPolicyUpdateKind) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAccessPolicyUpdateKind(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAccessPolicyUpdateKind(input string) (*AccessPolicyUpdateKind, error) {
	vals := map[string]AccessPolicyUpdateKind{
		"add":     AccessPolicyUpdateKindAdd,
		"remove":  AccessPolicyUpdateKindRemove,
		"replace": AccessPolicyUpdateKindReplace,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AccessPolicyUpdateKind(input)
	return &out, nil
}

type ActionsRequired string

const (
	ActionsRequiredNone ActionsRequired = "None"
)

func PossibleValuesForActionsRequired() []string {
	return []string{
		string(ActionsRequiredNone),
	}
}

func (s *ActionsRequired) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseActionsRequired(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseActionsRequired(input string) (*ActionsRequired, error) {
	vals := map[string]ActionsRequired{
		"none": ActionsRequiredNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ActionsRequired(input)
	return &out, nil
}

type CertificatePermissions string

const (
	CertificatePermissionsAll            CertificatePermissions = "all"
	CertificatePermissionsBackup         CertificatePermissions = "backup"
	CertificatePermissionsCreate         CertificatePermissions = "create"
	CertificatePermissionsDelete         CertificatePermissions = "delete"
	CertificatePermissionsDeleteissuers  CertificatePermissions = "deleteissuers"
	CertificatePermissionsGet            CertificatePermissions = "get"
	CertificatePermissionsGetissuers     CertificatePermissions = "getissuers"
	CertificatePermissionsImport         CertificatePermissions = "import"
	CertificatePermissionsList           CertificatePermissions = "list"
	CertificatePermissionsListissuers    CertificatePermissions = "listissuers"
	CertificatePermissionsManagecontacts CertificatePermissions = "managecontacts"
	CertificatePermissionsManageissuers  CertificatePermissions = "manageissuers"
	CertificatePermissionsPurge          CertificatePermissions = "purge"
	CertificatePermissionsRecover        CertificatePermissions = "recover"
	CertificatePermissionsRestore        CertificatePermissions = "restore"
	CertificatePermissionsSetissuers     CertificatePermissions = "setissuers"
	CertificatePermissionsUpdate         CertificatePermissions = "update"
)

func PossibleValuesForCertificatePermissions() []string {
	return []string{
		string(CertificatePermissionsAll),
		string(CertificatePermissionsBackup),
		string(CertificatePermissionsCreate),
		string(CertificatePermissionsDelete),
		string(CertificatePermissionsDeleteissuers),
		string(CertificatePermissionsGet),
		string(CertificatePermissionsGetissuers),
		string(CertificatePermissionsImport),
		string(CertificatePermissionsList),
		string(CertificatePermissionsListissuers),
		string(CertificatePermissionsManagecontacts),
		string(CertificatePermissionsManageissuers),
		string(CertificatePermissionsPurge),
		string(CertificatePermissionsRecover),
		string(CertificatePermissionsRestore),
		string(CertificatePermissionsSetissuers),
		string(CertificatePermissionsUpdate),
	}
}

func (s *CertificatePermissions) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCertificatePermissions(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCertificatePermissions(input string) (*CertificatePermissions, error) {
	vals := map[string]CertificatePermissions{
		"all":            CertificatePermissionsAll,
		"backup":         CertificatePermissionsBackup,
		"create":         CertificatePermissionsCreate,
		"delete":         CertificatePermissionsDelete,
		"deleteissuers":  CertificatePermissionsDeleteissuers,
		"get":            CertificatePermissionsGet,
		"getissuers":     CertificatePermissionsGetissuers,
		"import":         CertificatePermissionsImport,
		"list":           CertificatePermissionsList,
		"listissuers":    CertificatePermissionsListissuers,
		"managecontacts": CertificatePermissionsManagecontacts,
		"manageissuers":  CertificatePermissionsManageissuers,
		"purge":          CertificatePermissionsPurge,
		"recover":        CertificatePermissionsRecover,
		"restore":        CertificatePermissionsRestore,
		"setissuers":     CertificatePermissionsSetissuers,
		"update":         CertificatePermissionsUpdate,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CertificatePermissions(input)
	return &out, nil
}

type CreateMode string

const (
	CreateModeDefault CreateMode = "default"
	CreateModeRecover CreateMode = "recover"
)

func PossibleValuesForCreateMode() []string {
	return []string{
		string(CreateModeDefault),
		string(CreateModeRecover),
	}
}

func (s *CreateMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCreateMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCreateMode(input string) (*CreateMode, error) {
	vals := map[string]CreateMode{
		"default": CreateModeDefault,
		"recover": CreateModeRecover,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CreateMode(input)
	return &out, nil
}

type KeyPermissions string

const (
	KeyPermissionsAll               KeyPermissions = "all"
	KeyPermissionsBackup            KeyPermissions = "backup"
	KeyPermissionsCreate            KeyPermissions = "create"
	KeyPermissionsDecrypt           KeyPermissions = "decrypt"
	KeyPermissionsDelete            KeyPermissions = "delete"
	KeyPermissionsEncrypt           KeyPermissions = "encrypt"
	KeyPermissionsGet               KeyPermissions = "get"
	KeyPermissionsGetrotationpolicy KeyPermissions = "getrotationpolicy"
	KeyPermissionsImport            KeyPermissions = "import"
	KeyPermissionsList              KeyPermissions = "list"
	KeyPermissionsPurge             KeyPermissions = "purge"
	KeyPermissionsRecover           KeyPermissions = "recover"
	KeyPermissionsRelease           KeyPermissions = "release"
	KeyPermissionsRestore           KeyPermissions = "restore"
	KeyPermissionsRotate            KeyPermissions = "rotate"
	KeyPermissionsSetrotationpolicy KeyPermissions = "setrotationpolicy"
	KeyPermissionsSign              KeyPermissions = "sign"
	KeyPermissionsUnwrapKey         KeyPermissions = "unwrapKey"
	KeyPermissionsUpdate            KeyPermissions = "update"
	KeyPermissionsVerify            KeyPermissions = "verify"
	KeyPermissionsWrapKey           KeyPermissions = "wrapKey"
)

func PossibleValuesForKeyPermissions() []string {
	return []string{
		string(KeyPermissionsAll),
		string(KeyPermissionsBackup),
		string(KeyPermissionsCreate),
		string(KeyPermissionsDecrypt),
		string(KeyPermissionsDelete),
		string(KeyPermissionsEncrypt),
		string(KeyPermissionsGet),
		string(KeyPermissionsGetrotationpolicy),
		string(KeyPermissionsImport),
		string(KeyPermissionsList),
		string(KeyPermissionsPurge),
		string(KeyPermissionsRecover),
		string(KeyPermissionsRelease),
		string(KeyPermissionsRestore),
		string(KeyPermissionsRotate),
		string(KeyPermissionsSetrotationpolicy),
		string(KeyPermissionsSign),
		string(KeyPermissionsUnwrapKey),
		string(KeyPermissionsUpdate),
		string(KeyPermissionsVerify),
		string(KeyPermissionsWrapKey),
	}
}

func (s *KeyPermissions) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseKeyPermissions(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseKeyPermissions(input string) (*KeyPermissions, error) {
	vals := map[string]KeyPermissions{
		"all":               KeyPermissionsAll,
		"backup":            KeyPermissionsBackup,
		"create":            KeyPermissionsCreate,
		"decrypt":           KeyPermissionsDecrypt,
		"delete":            KeyPermissionsDelete,
		"encrypt":           KeyPermissionsEncrypt,
		"get":               KeyPermissionsGet,
		"getrotationpolicy": KeyPermissionsGetrotationpolicy,
		"import":            KeyPermissionsImport,
		"list":              KeyPermissionsList,
		"purge":             KeyPermissionsPurge,
		"recover":           KeyPermissionsRecover,
		"release":           KeyPermissionsRelease,
		"restore":           KeyPermissionsRestore,
		"rotate":            KeyPermissionsRotate,
		"setrotationpolicy": KeyPermissionsSetrotationpolicy,
		"sign":              KeyPermissionsSign,
		"unwrapkey":         KeyPermissionsUnwrapKey,
		"update":            KeyPermissionsUpdate,
		"verify":            KeyPermissionsVerify,
		"wrapkey":           KeyPermissionsWrapKey,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := KeyPermissions(input)
	return &out, nil
}

type NetworkRuleAction string

const (
	NetworkRuleActionAllow NetworkRuleAction = "Allow"
	NetworkRuleActionDeny  NetworkRuleAction = "Deny"
)

func PossibleValuesForNetworkRuleAction() []string {
	return []string{
		string(NetworkRuleActionAllow),
		string(NetworkRuleActionDeny),
	}
}

func (s *NetworkRuleAction) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNetworkRuleAction(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNetworkRuleAction(input string) (*NetworkRuleAction, error) {
	vals := map[string]NetworkRuleAction{
		"allow": NetworkRuleActionAllow,
		"deny":  NetworkRuleActionDeny,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NetworkRuleAction(input)
	return &out, nil
}

type NetworkRuleBypassOptions string

const (
	NetworkRuleBypassOptionsAzureServices NetworkRuleBypassOptions = "AzureServices"
	NetworkRuleBypassOptionsNone          NetworkRuleBypassOptions = "None"
)

func PossibleValuesForNetworkRuleBypassOptions() []string {
	return []string{
		string(NetworkRuleBypassOptionsAzureServices),
		string(NetworkRuleBypassOptionsNone),
	}
}

func (s *NetworkRuleBypassOptions) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNetworkRuleBypassOptions(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNetworkRuleBypassOptions(input string) (*NetworkRuleBypassOptions, error) {
	vals := map[string]NetworkRuleBypassOptions{
		"azureservices": NetworkRuleBypassOptionsAzureServices,
		"none":          NetworkRuleBypassOptionsNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NetworkRuleBypassOptions(input)
	return &out, nil
}

type PrivateEndpointConnectionProvisioningState string

const (
	PrivateEndpointConnectionProvisioningStateCreating     PrivateEndpointConnectionProvisioningState = "Creating"
	PrivateEndpointConnectionProvisioningStateDeleting     PrivateEndpointConnectionProvisioningState = "Deleting"
	PrivateEndpointConnectionProvisioningStateDisconnected PrivateEndpointConnectionProvisioningState = "Disconnected"
	PrivateEndpointConnectionProvisioningStateFailed       PrivateEndpointConnectionProvisioningState = "Failed"
	PrivateEndpointConnectionProvisioningStateSucceeded    PrivateEndpointConnectionProvisioningState = "Succeeded"
	PrivateEndpointConnectionProvisioningStateUpdating     PrivateEndpointConnectionProvisioningState = "Updating"
)

func PossibleValuesForPrivateEndpointConnectionProvisioningState() []string {
	return []string{
		string(PrivateEndpointConnectionProvisioningStateCreating),
		string(PrivateEndpointConnectionProvisioningStateDeleting),
		string(PrivateEndpointConnectionProvisioningStateDisconnected),
		string(PrivateEndpointConnectionProvisioningStateFailed),
		string(PrivateEndpointConnectionProvisioningStateSucceeded),
		string(PrivateEndpointConnectionProvisioningStateUpdating),
	}
}

func (s *PrivateEndpointConnectionProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePrivateEndpointConnectionProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePrivateEndpointConnectionProvisioningState(input string) (*PrivateEndpointConnectionProvisioningState, error) {
	vals := map[string]PrivateEndpointConnectionProvisioningState{
		"creating":     PrivateEndpointConnectionProvisioningStateCreating,
		"deleting":     PrivateEndpointConnectionProvisioningStateDeleting,
		"disconnected": PrivateEndpointConnectionProvisioningStateDisconnected,
		"failed":       PrivateEndpointConnectionProvisioningStateFailed,
		"succeeded":    PrivateEndpointConnectionProvisioningStateSucceeded,
		"updating":     PrivateEndpointConnectionProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PrivateEndpointConnectionProvisioningState(input)
	return &out, nil
}

type PrivateEndpointServiceConnectionStatus string

const (
	PrivateEndpointServiceConnectionStatusApproved     PrivateEndpointServiceConnectionStatus = "Approved"
	PrivateEndpointServiceConnectionStatusDisconnected PrivateEndpointServiceConnectionStatus = "Disconnected"
	PrivateEndpointServiceConnectionStatusPending      PrivateEndpointServiceConnectionStatus = "Pending"
	PrivateEndpointServiceConnectionStatusRejected     PrivateEndpointServiceConnectionStatus = "Rejected"
)

func PossibleValuesForPrivateEndpointServiceConnectionStatus() []string {
	return []string{
		string(PrivateEndpointServiceConnectionStatusApproved),
		string(PrivateEndpointServiceConnectionStatusDisconnected),
		string(PrivateEndpointServiceConnectionStatusPending),
		string(PrivateEndpointServiceConnectionStatusRejected),
	}
}

func (s *PrivateEndpointServiceConnectionStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePrivateEndpointServiceConnectionStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePrivateEndpointServiceConnectionStatus(input string) (*PrivateEndpointServiceConnectionStatus, error) {
	vals := map[string]PrivateEndpointServiceConnectionStatus{
		"approved":     PrivateEndpointServiceConnectionStatusApproved,
		"disconnected": PrivateEndpointServiceConnectionStatusDisconnected,
		"pending":      PrivateEndpointServiceConnectionStatusPending,
		"rejected":     PrivateEndpointServiceConnectionStatusRejected,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PrivateEndpointServiceConnectionStatus(input)
	return &out, nil
}

type Reason string

const (
	ReasonAccountNameInvalid Reason = "AccountNameInvalid"
	ReasonAlreadyExists      Reason = "AlreadyExists"
)

func PossibleValuesForReason() []string {
	return []string{
		string(ReasonAccountNameInvalid),
		string(ReasonAlreadyExists),
	}
}

func (s *Reason) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseReason(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseReason(input string) (*Reason, error) {
	vals := map[string]Reason{
		"accountnameinvalid": ReasonAccountNameInvalid,
		"alreadyexists":      ReasonAlreadyExists,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Reason(input)
	return &out, nil
}

type SecretPermissions string

const (
	SecretPermissionsAll     SecretPermissions = "all"
	SecretPermissionsBackup  SecretPermissions = "backup"
	SecretPermissionsDelete  SecretPermissions = "delete"
	SecretPermissionsGet     SecretPermissions = "get"
	SecretPermissionsList    SecretPermissions = "list"
	SecretPermissionsPurge   SecretPermissions = "purge"
	SecretPermissionsRecover SecretPermissions = "recover"
	SecretPermissionsRestore SecretPermissions = "restore"
	SecretPermissionsSet     SecretPermissions = "set"
)

func PossibleValuesForSecretPermissions() []string {
	return []string{
		string(SecretPermissionsAll),
		string(SecretPermissionsBackup),
		string(SecretPermissionsDelete),
		string(SecretPermissionsGet),
		string(SecretPermissionsList),
		string(SecretPermissionsPurge),
		string(SecretPermissionsRecover),
		string(SecretPermissionsRestore),
		string(SecretPermissionsSet),
	}
}

func (s *SecretPermissions) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSecretPermissions(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSecretPermissions(input string) (*SecretPermissions, error) {
	vals := map[string]SecretPermissions{
		"all":     SecretPermissionsAll,
		"backup":  SecretPermissionsBackup,
		"delete":  SecretPermissionsDelete,
		"get":     SecretPermissionsGet,
		"list":    SecretPermissionsList,
		"purge":   SecretPermissionsPurge,
		"recover": SecretPermissionsRecover,
		"restore": SecretPermissionsRestore,
		"set":     SecretPermissionsSet,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SecretPermissions(input)
	return &out, nil
}

type SkuFamily string

const (
	SkuFamilyA SkuFamily = "A"
)

func PossibleValuesForSkuFamily() []string {
	return []string{
		string(SkuFamilyA),
	}
}

func (s *SkuFamily) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSkuFamily(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSkuFamily(input string) (*SkuFamily, error) {
	vals := map[string]SkuFamily{
		"a": SkuFamilyA,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SkuFamily(input)
	return &out, nil
}

type SkuName string

const (
	SkuNamePremium  SkuName = "premium"
	SkuNameStandard SkuName = "standard"
)

func PossibleValuesForSkuName() []string {
	return []string{
		string(SkuNamePremium),
		string(SkuNameStandard),
	}
}

func (s *SkuName) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSkuName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSkuName(input string) (*SkuName, error) {
	vals := map[string]SkuName{
		"premium":  SkuNamePremium,
		"standard": SkuNameStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SkuName(input)
	return &out, nil
}

type StoragePermissions string

const (
	StoragePermissionsAll           StoragePermissions = "all"
	StoragePermissionsBackup        StoragePermissions = "backup"
	StoragePermissionsDelete        StoragePermissions = "delete"
	StoragePermissionsDeletesas     StoragePermissions = "deletesas"
	StoragePermissionsGet           StoragePermissions = "get"
	StoragePermissionsGetsas        StoragePermissions = "getsas"
	StoragePermissionsList          StoragePermissions = "list"
	StoragePermissionsListsas       StoragePermissions = "listsas"
	StoragePermissionsPurge         StoragePermissions = "purge"
	StoragePermissionsRecover       StoragePermissions = "recover"
	StoragePermissionsRegeneratekey StoragePermissions = "regeneratekey"
	StoragePermissionsRestore       StoragePermissions = "restore"
	StoragePermissionsSet           StoragePermissions = "set"
	StoragePermissionsSetsas        StoragePermissions = "setsas"
	StoragePermissionsUpdate        StoragePermissions = "update"
)

func PossibleValuesForStoragePermissions() []string {
	return []string{
		string(StoragePermissionsAll),
		string(StoragePermissionsBackup),
		string(StoragePermissionsDelete),
		string(StoragePermissionsDeletesas),
		string(StoragePermissionsGet),
		string(StoragePermissionsGetsas),
		string(StoragePermissionsList),
		string(StoragePermissionsListsas),
		string(StoragePermissionsPurge),
		string(StoragePermissionsRecover),
		string(StoragePermissionsRegeneratekey),
		string(StoragePermissionsRestore),
		string(StoragePermissionsSet),
		string(StoragePermissionsSetsas),
		string(StoragePermissionsUpdate),
	}
}

func (s *StoragePermissions) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStoragePermissions(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStoragePermissions(input string) (*StoragePermissions, error) {
	vals := map[string]StoragePermissions{
		"all":           StoragePermissionsAll,
		"backup":        StoragePermissionsBackup,
		"delete":        StoragePermissionsDelete,
		"deletesas":     StoragePermissionsDeletesas,
		"get":           StoragePermissionsGet,
		"getsas":        StoragePermissionsGetsas,
		"list":          StoragePermissionsList,
		"listsas":       StoragePermissionsListsas,
		"purge":         StoragePermissionsPurge,
		"recover":       StoragePermissionsRecover,
		"regeneratekey": StoragePermissionsRegeneratekey,
		"restore":       StoragePermissionsRestore,
		"set":           StoragePermissionsSet,
		"setsas":        StoragePermissionsSetsas,
		"update":        StoragePermissionsUpdate,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StoragePermissions(input)
	return &out, nil
}

type Type string

const (
	TypeMicrosoftPointKeyVaultVaults Type = "Microsoft.KeyVault/vaults"
)

func PossibleValuesForType() []string {
	return []string{
		string(TypeMicrosoftPointKeyVaultVaults),
	}
}

func (s *Type) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseType(input string) (*Type, error) {
	vals := map[string]Type{
		"microsoft.keyvault/vaults": TypeMicrosoftPointKeyVaultVaults,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Type(input)
	return &out, nil
}

type VaultListFilterTypes string

const (
	VaultListFilterTypesResourceTypeEqMicrosoftPointKeyVaultVaults VaultListFilterTypes = "resourceType eq 'Microsoft.KeyVault/vaults'"
)

func PossibleValuesForVaultListFilterTypes() []string {
	return []string{
		string(VaultListFilterTypesResourceTypeEqMicrosoftPointKeyVaultVaults),
	}
}

func (s *VaultListFilterTypes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVaultListFilterTypes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVaultListFilterTypes(input string) (*VaultListFilterTypes, error) {
	vals := map[string]VaultListFilterTypes{
		"resourcetype eq 'microsoft.keyvault/vaults'": VaultListFilterTypesResourceTypeEqMicrosoftPointKeyVaultVaults,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VaultListFilterTypes(input)
	return &out, nil
}

type VaultProvisioningState string

const (
	VaultProvisioningStateRegisteringDns VaultProvisioningState = "RegisteringDns"
	VaultProvisioningStateSucceeded      VaultProvisioningState = "Succeeded"
)

func PossibleValuesForVaultProvisioningState() []string {
	return []string{
		string(VaultProvisioningStateRegisteringDns),
		string(VaultProvisioningStateSucceeded),
	}
}

func (s *VaultProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVaultProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVaultProvisioningState(input string) (*VaultProvisioningState, error) {
	vals := map[string]VaultProvisioningState{
		"registeringdns": VaultProvisioningStateRegisteringDns,
		"succeeded":      VaultProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VaultProvisioningState(input)
	return &out, nil
}
