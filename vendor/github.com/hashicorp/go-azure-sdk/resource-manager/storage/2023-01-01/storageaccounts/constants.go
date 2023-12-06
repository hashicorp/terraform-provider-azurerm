package storageaccounts

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccessTier string

const (
	AccessTierCool    AccessTier = "Cool"
	AccessTierHot     AccessTier = "Hot"
	AccessTierPremium AccessTier = "Premium"
)

func PossibleValuesForAccessTier() []string {
	return []string{
		string(AccessTierCool),
		string(AccessTierHot),
		string(AccessTierPremium),
	}
}

func (s *AccessTier) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAccessTier(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAccessTier(input string) (*AccessTier, error) {
	vals := map[string]AccessTier{
		"cool":    AccessTierCool,
		"hot":     AccessTierHot,
		"premium": AccessTierPremium,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AccessTier(input)
	return &out, nil
}

type AccountImmutabilityPolicyState string

const (
	AccountImmutabilityPolicyStateDisabled AccountImmutabilityPolicyState = "Disabled"
	AccountImmutabilityPolicyStateLocked   AccountImmutabilityPolicyState = "Locked"
	AccountImmutabilityPolicyStateUnlocked AccountImmutabilityPolicyState = "Unlocked"
)

func PossibleValuesForAccountImmutabilityPolicyState() []string {
	return []string{
		string(AccountImmutabilityPolicyStateDisabled),
		string(AccountImmutabilityPolicyStateLocked),
		string(AccountImmutabilityPolicyStateUnlocked),
	}
}

func (s *AccountImmutabilityPolicyState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAccountImmutabilityPolicyState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAccountImmutabilityPolicyState(input string) (*AccountImmutabilityPolicyState, error) {
	vals := map[string]AccountImmutabilityPolicyState{
		"disabled": AccountImmutabilityPolicyStateDisabled,
		"locked":   AccountImmutabilityPolicyStateLocked,
		"unlocked": AccountImmutabilityPolicyStateUnlocked,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AccountImmutabilityPolicyState(input)
	return &out, nil
}

type AccountStatus string

const (
	AccountStatusAvailable   AccountStatus = "available"
	AccountStatusUnavailable AccountStatus = "unavailable"
)

func PossibleValuesForAccountStatus() []string {
	return []string{
		string(AccountStatusAvailable),
		string(AccountStatusUnavailable),
	}
}

func (s *AccountStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAccountStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAccountStatus(input string) (*AccountStatus, error) {
	vals := map[string]AccountStatus{
		"available":   AccountStatusAvailable,
		"unavailable": AccountStatusUnavailable,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AccountStatus(input)
	return &out, nil
}

type AccountType string

const (
	AccountTypeComputer AccountType = "Computer"
	AccountTypeUser     AccountType = "User"
)

func PossibleValuesForAccountType() []string {
	return []string{
		string(AccountTypeComputer),
		string(AccountTypeUser),
	}
}

func (s *AccountType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAccountType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAccountType(input string) (*AccountType, error) {
	vals := map[string]AccountType{
		"computer": AccountTypeComputer,
		"user":     AccountTypeUser,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AccountType(input)
	return &out, nil
}

type Action string

const (
	ActionAllow Action = "Allow"
)

func PossibleValuesForAction() []string {
	return []string{
		string(ActionAllow),
	}
}

func (s *Action) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAction(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAction(input string) (*Action, error) {
	vals := map[string]Action{
		"allow": ActionAllow,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Action(input)
	return &out, nil
}

type AllowedCopyScope string

const (
	AllowedCopyScopeAAD         AllowedCopyScope = "AAD"
	AllowedCopyScopePrivateLink AllowedCopyScope = "PrivateLink"
)

func PossibleValuesForAllowedCopyScope() []string {
	return []string{
		string(AllowedCopyScopeAAD),
		string(AllowedCopyScopePrivateLink),
	}
}

func (s *AllowedCopyScope) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAllowedCopyScope(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAllowedCopyScope(input string) (*AllowedCopyScope, error) {
	vals := map[string]AllowedCopyScope{
		"aad":         AllowedCopyScopeAAD,
		"privatelink": AllowedCopyScopePrivateLink,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AllowedCopyScope(input)
	return &out, nil
}

type BlobRestoreProgressStatus string

const (
	BlobRestoreProgressStatusComplete   BlobRestoreProgressStatus = "Complete"
	BlobRestoreProgressStatusFailed     BlobRestoreProgressStatus = "Failed"
	BlobRestoreProgressStatusInProgress BlobRestoreProgressStatus = "InProgress"
)

func PossibleValuesForBlobRestoreProgressStatus() []string {
	return []string{
		string(BlobRestoreProgressStatusComplete),
		string(BlobRestoreProgressStatusFailed),
		string(BlobRestoreProgressStatusInProgress),
	}
}

func (s *BlobRestoreProgressStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBlobRestoreProgressStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBlobRestoreProgressStatus(input string) (*BlobRestoreProgressStatus, error) {
	vals := map[string]BlobRestoreProgressStatus{
		"complete":   BlobRestoreProgressStatusComplete,
		"failed":     BlobRestoreProgressStatusFailed,
		"inprogress": BlobRestoreProgressStatusInProgress,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BlobRestoreProgressStatus(input)
	return &out, nil
}

type Bypass string

const (
	BypassAzureServices Bypass = "AzureServices"
	BypassLogging       Bypass = "Logging"
	BypassMetrics       Bypass = "Metrics"
	BypassNone          Bypass = "None"
)

func PossibleValuesForBypass() []string {
	return []string{
		string(BypassAzureServices),
		string(BypassLogging),
		string(BypassMetrics),
		string(BypassNone),
	}
}

func (s *Bypass) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBypass(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBypass(input string) (*Bypass, error) {
	vals := map[string]Bypass{
		"azureservices": BypassAzureServices,
		"logging":       BypassLogging,
		"metrics":       BypassMetrics,
		"none":          BypassNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Bypass(input)
	return &out, nil
}

type DefaultAction string

const (
	DefaultActionAllow DefaultAction = "Allow"
	DefaultActionDeny  DefaultAction = "Deny"
)

func PossibleValuesForDefaultAction() []string {
	return []string{
		string(DefaultActionAllow),
		string(DefaultActionDeny),
	}
}

func (s *DefaultAction) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDefaultAction(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDefaultAction(input string) (*DefaultAction, error) {
	vals := map[string]DefaultAction{
		"allow": DefaultActionAllow,
		"deny":  DefaultActionDeny,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DefaultAction(input)
	return &out, nil
}

type DefaultSharePermission string

const (
	DefaultSharePermissionNone                                       DefaultSharePermission = "None"
	DefaultSharePermissionStorageFileDataSmbShareContributor         DefaultSharePermission = "StorageFileDataSmbShareContributor"
	DefaultSharePermissionStorageFileDataSmbShareElevatedContributor DefaultSharePermission = "StorageFileDataSmbShareElevatedContributor"
	DefaultSharePermissionStorageFileDataSmbShareReader              DefaultSharePermission = "StorageFileDataSmbShareReader"
)

func PossibleValuesForDefaultSharePermission() []string {
	return []string{
		string(DefaultSharePermissionNone),
		string(DefaultSharePermissionStorageFileDataSmbShareContributor),
		string(DefaultSharePermissionStorageFileDataSmbShareElevatedContributor),
		string(DefaultSharePermissionStorageFileDataSmbShareReader),
	}
}

func (s *DefaultSharePermission) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDefaultSharePermission(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDefaultSharePermission(input string) (*DefaultSharePermission, error) {
	vals := map[string]DefaultSharePermission{
		"none":                               DefaultSharePermissionNone,
		"storagefiledatasmbsharecontributor": DefaultSharePermissionStorageFileDataSmbShareContributor,
		"storagefiledatasmbshareelevatedcontributor": DefaultSharePermissionStorageFileDataSmbShareElevatedContributor,
		"storagefiledatasmbsharereader":              DefaultSharePermissionStorageFileDataSmbShareReader,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DefaultSharePermission(input)
	return &out, nil
}

type DirectoryServiceOptions string

const (
	DirectoryServiceOptionsAADDS   DirectoryServiceOptions = "AADDS"
	DirectoryServiceOptionsAADKERB DirectoryServiceOptions = "AADKERB"
	DirectoryServiceOptionsAD      DirectoryServiceOptions = "AD"
	DirectoryServiceOptionsNone    DirectoryServiceOptions = "None"
)

func PossibleValuesForDirectoryServiceOptions() []string {
	return []string{
		string(DirectoryServiceOptionsAADDS),
		string(DirectoryServiceOptionsAADKERB),
		string(DirectoryServiceOptionsAD),
		string(DirectoryServiceOptionsNone),
	}
}

func (s *DirectoryServiceOptions) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDirectoryServiceOptions(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDirectoryServiceOptions(input string) (*DirectoryServiceOptions, error) {
	vals := map[string]DirectoryServiceOptions{
		"aadds":   DirectoryServiceOptionsAADDS,
		"aadkerb": DirectoryServiceOptionsAADKERB,
		"ad":      DirectoryServiceOptionsAD,
		"none":    DirectoryServiceOptionsNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DirectoryServiceOptions(input)
	return &out, nil
}

type DnsEndpointType string

const (
	DnsEndpointTypeAzureDnsZone DnsEndpointType = "AzureDnsZone"
	DnsEndpointTypeStandard     DnsEndpointType = "Standard"
)

func PossibleValuesForDnsEndpointType() []string {
	return []string{
		string(DnsEndpointTypeAzureDnsZone),
		string(DnsEndpointTypeStandard),
	}
}

func (s *DnsEndpointType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDnsEndpointType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDnsEndpointType(input string) (*DnsEndpointType, error) {
	vals := map[string]DnsEndpointType{
		"azurednszone": DnsEndpointTypeAzureDnsZone,
		"standard":     DnsEndpointTypeStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DnsEndpointType(input)
	return &out, nil
}

type ExpirationAction string

const (
	ExpirationActionLog ExpirationAction = "Log"
)

func PossibleValuesForExpirationAction() []string {
	return []string{
		string(ExpirationActionLog),
	}
}

func (s *ExpirationAction) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseExpirationAction(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseExpirationAction(input string) (*ExpirationAction, error) {
	vals := map[string]ExpirationAction{
		"log": ExpirationActionLog,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ExpirationAction(input)
	return &out, nil
}

type FailoverType string

const (
	FailoverTypePlanned FailoverType = "Planned"
)

func PossibleValuesForFailoverType() []string {
	return []string{
		string(FailoverTypePlanned),
	}
}

func (s *FailoverType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseFailoverType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseFailoverType(input string) (*FailoverType, error) {
	vals := map[string]FailoverType{
		"planned": FailoverTypePlanned,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FailoverType(input)
	return &out, nil
}

type GeoReplicationStatus string

const (
	GeoReplicationStatusBootstrap   GeoReplicationStatus = "Bootstrap"
	GeoReplicationStatusLive        GeoReplicationStatus = "Live"
	GeoReplicationStatusUnavailable GeoReplicationStatus = "Unavailable"
)

func PossibleValuesForGeoReplicationStatus() []string {
	return []string{
		string(GeoReplicationStatusBootstrap),
		string(GeoReplicationStatusLive),
		string(GeoReplicationStatusUnavailable),
	}
}

func (s *GeoReplicationStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseGeoReplicationStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseGeoReplicationStatus(input string) (*GeoReplicationStatus, error) {
	vals := map[string]GeoReplicationStatus{
		"bootstrap":   GeoReplicationStatusBootstrap,
		"live":        GeoReplicationStatusLive,
		"unavailable": GeoReplicationStatusUnavailable,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := GeoReplicationStatus(input)
	return &out, nil
}

type HTTPProtocol string

const (
	HTTPProtocolHTTPS     HTTPProtocol = "https"
	HTTPProtocolHTTPSHttp HTTPProtocol = "https,http"
)

func PossibleValuesForHTTPProtocol() []string {
	return []string{
		string(HTTPProtocolHTTPS),
		string(HTTPProtocolHTTPSHttp),
	}
}

func (s *HTTPProtocol) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseHTTPProtocol(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseHTTPProtocol(input string) (*HTTPProtocol, error) {
	vals := map[string]HTTPProtocol{
		"https":      HTTPProtocolHTTPS,
		"https,http": HTTPProtocolHTTPSHttp,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HTTPProtocol(input)
	return &out, nil
}

type KeyPermission string

const (
	KeyPermissionFull KeyPermission = "Full"
	KeyPermissionRead KeyPermission = "Read"
)

func PossibleValuesForKeyPermission() []string {
	return []string{
		string(KeyPermissionFull),
		string(KeyPermissionRead),
	}
}

func (s *KeyPermission) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseKeyPermission(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseKeyPermission(input string) (*KeyPermission, error) {
	vals := map[string]KeyPermission{
		"full": KeyPermissionFull,
		"read": KeyPermissionRead,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := KeyPermission(input)
	return &out, nil
}

type KeySource string

const (
	KeySourceMicrosoftPointKeyvault KeySource = "Microsoft.Keyvault"
	KeySourceMicrosoftPointStorage  KeySource = "Microsoft.Storage"
)

func PossibleValuesForKeySource() []string {
	return []string{
		string(KeySourceMicrosoftPointKeyvault),
		string(KeySourceMicrosoftPointStorage),
	}
}

func (s *KeySource) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseKeySource(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseKeySource(input string) (*KeySource, error) {
	vals := map[string]KeySource{
		"microsoft.keyvault": KeySourceMicrosoftPointKeyvault,
		"microsoft.storage":  KeySourceMicrosoftPointStorage,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := KeySource(input)
	return &out, nil
}

type KeyType string

const (
	KeyTypeAccount KeyType = "Account"
	KeyTypeService KeyType = "Service"
)

func PossibleValuesForKeyType() []string {
	return []string{
		string(KeyTypeAccount),
		string(KeyTypeService),
	}
}

func (s *KeyType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseKeyType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseKeyType(input string) (*KeyType, error) {
	vals := map[string]KeyType{
		"account": KeyTypeAccount,
		"service": KeyTypeService,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := KeyType(input)
	return &out, nil
}

type Kind string

const (
	KindBlobStorage      Kind = "BlobStorage"
	KindBlockBlobStorage Kind = "BlockBlobStorage"
	KindFileStorage      Kind = "FileStorage"
	KindStorage          Kind = "Storage"
	KindStorageVTwo      Kind = "StorageV2"
)

func PossibleValuesForKind() []string {
	return []string{
		string(KindBlobStorage),
		string(KindBlockBlobStorage),
		string(KindFileStorage),
		string(KindStorage),
		string(KindStorageVTwo),
	}
}

func (s *Kind) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseKind(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseKind(input string) (*Kind, error) {
	vals := map[string]Kind{
		"blobstorage":      KindBlobStorage,
		"blockblobstorage": KindBlockBlobStorage,
		"filestorage":      KindFileStorage,
		"storage":          KindStorage,
		"storagev2":        KindStorageVTwo,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Kind(input)
	return &out, nil
}

type LargeFileSharesState string

const (
	LargeFileSharesStateDisabled LargeFileSharesState = "Disabled"
	LargeFileSharesStateEnabled  LargeFileSharesState = "Enabled"
)

func PossibleValuesForLargeFileSharesState() []string {
	return []string{
		string(LargeFileSharesStateDisabled),
		string(LargeFileSharesStateEnabled),
	}
}

func (s *LargeFileSharesState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLargeFileSharesState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLargeFileSharesState(input string) (*LargeFileSharesState, error) {
	vals := map[string]LargeFileSharesState{
		"disabled": LargeFileSharesStateDisabled,
		"enabled":  LargeFileSharesStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LargeFileSharesState(input)
	return &out, nil
}

type ListKeyExpand string

const (
	ListKeyExpandKerb ListKeyExpand = "kerb"
)

func PossibleValuesForListKeyExpand() []string {
	return []string{
		string(ListKeyExpandKerb),
	}
}

func (s *ListKeyExpand) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseListKeyExpand(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseListKeyExpand(input string) (*ListKeyExpand, error) {
	vals := map[string]ListKeyExpand{
		"kerb": ListKeyExpandKerb,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ListKeyExpand(input)
	return &out, nil
}

type MinimumTlsVersion string

const (
	MinimumTlsVersionTLSOneOne  MinimumTlsVersion = "TLS1_1"
	MinimumTlsVersionTLSOneTwo  MinimumTlsVersion = "TLS1_2"
	MinimumTlsVersionTLSOneZero MinimumTlsVersion = "TLS1_0"
)

func PossibleValuesForMinimumTlsVersion() []string {
	return []string{
		string(MinimumTlsVersionTLSOneOne),
		string(MinimumTlsVersionTLSOneTwo),
		string(MinimumTlsVersionTLSOneZero),
	}
}

func (s *MinimumTlsVersion) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMinimumTlsVersion(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMinimumTlsVersion(input string) (*MinimumTlsVersion, error) {
	vals := map[string]MinimumTlsVersion{
		"tls1_1": MinimumTlsVersionTLSOneOne,
		"tls1_2": MinimumTlsVersionTLSOneTwo,
		"tls1_0": MinimumTlsVersionTLSOneZero,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MinimumTlsVersion(input)
	return &out, nil
}

type Permissions string

const (
	PermissionsA Permissions = "a"
	PermissionsC Permissions = "c"
	PermissionsD Permissions = "d"
	PermissionsL Permissions = "l"
	PermissionsP Permissions = "p"
	PermissionsR Permissions = "r"
	PermissionsU Permissions = "u"
	PermissionsW Permissions = "w"
)

func PossibleValuesForPermissions() []string {
	return []string{
		string(PermissionsA),
		string(PermissionsC),
		string(PermissionsD),
		string(PermissionsL),
		string(PermissionsP),
		string(PermissionsR),
		string(PermissionsU),
		string(PermissionsW),
	}
}

func (s *Permissions) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePermissions(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePermissions(input string) (*Permissions, error) {
	vals := map[string]Permissions{
		"a": PermissionsA,
		"c": PermissionsC,
		"d": PermissionsD,
		"l": PermissionsL,
		"p": PermissionsP,
		"r": PermissionsR,
		"u": PermissionsU,
		"w": PermissionsW,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Permissions(input)
	return &out, nil
}

type PostFailoverRedundancy string

const (
	PostFailoverRedundancyStandardLRS PostFailoverRedundancy = "Standard_LRS"
	PostFailoverRedundancyStandardZRS PostFailoverRedundancy = "Standard_ZRS"
)

func PossibleValuesForPostFailoverRedundancy() []string {
	return []string{
		string(PostFailoverRedundancyStandardLRS),
		string(PostFailoverRedundancyStandardZRS),
	}
}

func (s *PostFailoverRedundancy) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePostFailoverRedundancy(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePostFailoverRedundancy(input string) (*PostFailoverRedundancy, error) {
	vals := map[string]PostFailoverRedundancy{
		"standard_lrs": PostFailoverRedundancyStandardLRS,
		"standard_zrs": PostFailoverRedundancyStandardZRS,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PostFailoverRedundancy(input)
	return &out, nil
}

type PostPlannedFailoverRedundancy string

const (
	PostPlannedFailoverRedundancyStandardGRS    PostPlannedFailoverRedundancy = "Standard_GRS"
	PostPlannedFailoverRedundancyStandardGZRS   PostPlannedFailoverRedundancy = "Standard_GZRS"
	PostPlannedFailoverRedundancyStandardRAGRS  PostPlannedFailoverRedundancy = "Standard_RAGRS"
	PostPlannedFailoverRedundancyStandardRAGZRS PostPlannedFailoverRedundancy = "Standard_RAGZRS"
)

func PossibleValuesForPostPlannedFailoverRedundancy() []string {
	return []string{
		string(PostPlannedFailoverRedundancyStandardGRS),
		string(PostPlannedFailoverRedundancyStandardGZRS),
		string(PostPlannedFailoverRedundancyStandardRAGRS),
		string(PostPlannedFailoverRedundancyStandardRAGZRS),
	}
}

func (s *PostPlannedFailoverRedundancy) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePostPlannedFailoverRedundancy(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePostPlannedFailoverRedundancy(input string) (*PostPlannedFailoverRedundancy, error) {
	vals := map[string]PostPlannedFailoverRedundancy{
		"standard_grs":    PostPlannedFailoverRedundancyStandardGRS,
		"standard_gzrs":   PostPlannedFailoverRedundancyStandardGZRS,
		"standard_ragrs":  PostPlannedFailoverRedundancyStandardRAGRS,
		"standard_ragzrs": PostPlannedFailoverRedundancyStandardRAGZRS,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PostPlannedFailoverRedundancy(input)
	return &out, nil
}

type PrivateEndpointConnectionProvisioningState string

const (
	PrivateEndpointConnectionProvisioningStateCreating  PrivateEndpointConnectionProvisioningState = "Creating"
	PrivateEndpointConnectionProvisioningStateDeleting  PrivateEndpointConnectionProvisioningState = "Deleting"
	PrivateEndpointConnectionProvisioningStateFailed    PrivateEndpointConnectionProvisioningState = "Failed"
	PrivateEndpointConnectionProvisioningStateSucceeded PrivateEndpointConnectionProvisioningState = "Succeeded"
)

func PossibleValuesForPrivateEndpointConnectionProvisioningState() []string {
	return []string{
		string(PrivateEndpointConnectionProvisioningStateCreating),
		string(PrivateEndpointConnectionProvisioningStateDeleting),
		string(PrivateEndpointConnectionProvisioningStateFailed),
		string(PrivateEndpointConnectionProvisioningStateSucceeded),
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
		"creating":  PrivateEndpointConnectionProvisioningStateCreating,
		"deleting":  PrivateEndpointConnectionProvisioningStateDeleting,
		"failed":    PrivateEndpointConnectionProvisioningStateFailed,
		"succeeded": PrivateEndpointConnectionProvisioningStateSucceeded,
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
	PrivateEndpointServiceConnectionStatusApproved PrivateEndpointServiceConnectionStatus = "Approved"
	PrivateEndpointServiceConnectionStatusPending  PrivateEndpointServiceConnectionStatus = "Pending"
	PrivateEndpointServiceConnectionStatusRejected PrivateEndpointServiceConnectionStatus = "Rejected"
)

func PossibleValuesForPrivateEndpointServiceConnectionStatus() []string {
	return []string{
		string(PrivateEndpointServiceConnectionStatusApproved),
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
		"approved": PrivateEndpointServiceConnectionStatusApproved,
		"pending":  PrivateEndpointServiceConnectionStatusPending,
		"rejected": PrivateEndpointServiceConnectionStatusRejected,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PrivateEndpointServiceConnectionStatus(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateCreating     ProvisioningState = "Creating"
	ProvisioningStateResolvingDNS ProvisioningState = "ResolvingDNS"
	ProvisioningStateSucceeded    ProvisioningState = "Succeeded"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateCreating),
		string(ProvisioningStateResolvingDNS),
		string(ProvisioningStateSucceeded),
	}
}

func (s *ProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"creating":     ProvisioningStateCreating,
		"resolvingdns": ProvisioningStateResolvingDNS,
		"succeeded":    ProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}

type PublicNetworkAccess string

const (
	PublicNetworkAccessDisabled PublicNetworkAccess = "Disabled"
	PublicNetworkAccessEnabled  PublicNetworkAccess = "Enabled"
)

func PossibleValuesForPublicNetworkAccess() []string {
	return []string{
		string(PublicNetworkAccessDisabled),
		string(PublicNetworkAccessEnabled),
	}
}

func (s *PublicNetworkAccess) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePublicNetworkAccess(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePublicNetworkAccess(input string) (*PublicNetworkAccess, error) {
	vals := map[string]PublicNetworkAccess{
		"disabled": PublicNetworkAccessDisabled,
		"enabled":  PublicNetworkAccessEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PublicNetworkAccess(input)
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

type RoutingChoice string

const (
	RoutingChoiceInternetRouting  RoutingChoice = "InternetRouting"
	RoutingChoiceMicrosoftRouting RoutingChoice = "MicrosoftRouting"
)

func PossibleValuesForRoutingChoice() []string {
	return []string{
		string(RoutingChoiceInternetRouting),
		string(RoutingChoiceMicrosoftRouting),
	}
}

func (s *RoutingChoice) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRoutingChoice(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRoutingChoice(input string) (*RoutingChoice, error) {
	vals := map[string]RoutingChoice{
		"internetrouting":  RoutingChoiceInternetRouting,
		"microsoftrouting": RoutingChoiceMicrosoftRouting,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RoutingChoice(input)
	return &out, nil
}

type Services string

const (
	ServicesB Services = "b"
	ServicesF Services = "f"
	ServicesQ Services = "q"
	ServicesT Services = "t"
)

func PossibleValuesForServices() []string {
	return []string{
		string(ServicesB),
		string(ServicesF),
		string(ServicesQ),
		string(ServicesT),
	}
}

func (s *Services) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseServices(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseServices(input string) (*Services, error) {
	vals := map[string]Services{
		"b": ServicesB,
		"f": ServicesF,
		"q": ServicesQ,
		"t": ServicesT,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Services(input)
	return &out, nil
}

type SignedResource string

const (
	SignedResourceB SignedResource = "b"
	SignedResourceC SignedResource = "c"
	SignedResourceF SignedResource = "f"
	SignedResourceS SignedResource = "s"
)

func PossibleValuesForSignedResource() []string {
	return []string{
		string(SignedResourceB),
		string(SignedResourceC),
		string(SignedResourceF),
		string(SignedResourceS),
	}
}

func (s *SignedResource) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSignedResource(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSignedResource(input string) (*SignedResource, error) {
	vals := map[string]SignedResource{
		"b": SignedResourceB,
		"c": SignedResourceC,
		"f": SignedResourceF,
		"s": SignedResourceS,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SignedResource(input)
	return &out, nil
}

type SignedResourceTypes string

const (
	SignedResourceTypesC SignedResourceTypes = "c"
	SignedResourceTypesO SignedResourceTypes = "o"
	SignedResourceTypesS SignedResourceTypes = "s"
)

func PossibleValuesForSignedResourceTypes() []string {
	return []string{
		string(SignedResourceTypesC),
		string(SignedResourceTypesO),
		string(SignedResourceTypesS),
	}
}

func (s *SignedResourceTypes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSignedResourceTypes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSignedResourceTypes(input string) (*SignedResourceTypes, error) {
	vals := map[string]SignedResourceTypes{
		"c": SignedResourceTypesC,
		"o": SignedResourceTypesO,
		"s": SignedResourceTypesS,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SignedResourceTypes(input)
	return &out, nil
}

type SkuConversionStatus string

const (
	SkuConversionStatusFailed     SkuConversionStatus = "Failed"
	SkuConversionStatusInProgress SkuConversionStatus = "InProgress"
	SkuConversionStatusSucceeded  SkuConversionStatus = "Succeeded"
)

func PossibleValuesForSkuConversionStatus() []string {
	return []string{
		string(SkuConversionStatusFailed),
		string(SkuConversionStatusInProgress),
		string(SkuConversionStatusSucceeded),
	}
}

func (s *SkuConversionStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSkuConversionStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSkuConversionStatus(input string) (*SkuConversionStatus, error) {
	vals := map[string]SkuConversionStatus{
		"failed":     SkuConversionStatusFailed,
		"inprogress": SkuConversionStatusInProgress,
		"succeeded":  SkuConversionStatusSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SkuConversionStatus(input)
	return &out, nil
}

type SkuName string

const (
	SkuNamePremiumLRS     SkuName = "Premium_LRS"
	SkuNamePremiumZRS     SkuName = "Premium_ZRS"
	SkuNameStandardGRS    SkuName = "Standard_GRS"
	SkuNameStandardGZRS   SkuName = "Standard_GZRS"
	SkuNameStandardLRS    SkuName = "Standard_LRS"
	SkuNameStandardRAGRS  SkuName = "Standard_RAGRS"
	SkuNameStandardRAGZRS SkuName = "Standard_RAGZRS"
	SkuNameStandardZRS    SkuName = "Standard_ZRS"
)

func PossibleValuesForSkuName() []string {
	return []string{
		string(SkuNamePremiumLRS),
		string(SkuNamePremiumZRS),
		string(SkuNameStandardGRS),
		string(SkuNameStandardGZRS),
		string(SkuNameStandardLRS),
		string(SkuNameStandardRAGRS),
		string(SkuNameStandardRAGZRS),
		string(SkuNameStandardZRS),
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
		"premium_lrs":     SkuNamePremiumLRS,
		"premium_zrs":     SkuNamePremiumZRS,
		"standard_grs":    SkuNameStandardGRS,
		"standard_gzrs":   SkuNameStandardGZRS,
		"standard_lrs":    SkuNameStandardLRS,
		"standard_ragrs":  SkuNameStandardRAGRS,
		"standard_ragzrs": SkuNameStandardRAGZRS,
		"standard_zrs":    SkuNameStandardZRS,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SkuName(input)
	return &out, nil
}

type SkuTier string

const (
	SkuTierPremium  SkuTier = "Premium"
	SkuTierStandard SkuTier = "Standard"
)

func PossibleValuesForSkuTier() []string {
	return []string{
		string(SkuTierPremium),
		string(SkuTierStandard),
	}
}

func (s *SkuTier) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSkuTier(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSkuTier(input string) (*SkuTier, error) {
	vals := map[string]SkuTier{
		"premium":  SkuTierPremium,
		"standard": SkuTierStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SkuTier(input)
	return &out, nil
}

type State string

const (
	StateDeprovisioning       State = "Deprovisioning"
	StateFailed               State = "Failed"
	StateNetworkSourceDeleted State = "NetworkSourceDeleted"
	StateProvisioning         State = "Provisioning"
	StateSucceeded            State = "Succeeded"
)

func PossibleValuesForState() []string {
	return []string{
		string(StateDeprovisioning),
		string(StateFailed),
		string(StateNetworkSourceDeleted),
		string(StateProvisioning),
		string(StateSucceeded),
	}
}

func (s *State) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseState(input string) (*State, error) {
	vals := map[string]State{
		"deprovisioning":       StateDeprovisioning,
		"failed":               StateFailed,
		"networksourcedeleted": StateNetworkSourceDeleted,
		"provisioning":         StateProvisioning,
		"succeeded":            StateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := State(input)
	return &out, nil
}

type StorageAccountExpand string

const (
	StorageAccountExpandBlobRestoreStatus   StorageAccountExpand = "blobRestoreStatus"
	StorageAccountExpandGeoReplicationStats StorageAccountExpand = "geoReplicationStats"
)

func PossibleValuesForStorageAccountExpand() []string {
	return []string{
		string(StorageAccountExpandBlobRestoreStatus),
		string(StorageAccountExpandGeoReplicationStats),
	}
}

func (s *StorageAccountExpand) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStorageAccountExpand(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStorageAccountExpand(input string) (*StorageAccountExpand, error) {
	vals := map[string]StorageAccountExpand{
		"blobrestorestatus":   StorageAccountExpandBlobRestoreStatus,
		"georeplicationstats": StorageAccountExpandGeoReplicationStats,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StorageAccountExpand(input)
	return &out, nil
}

type Type string

const (
	TypeMicrosoftPointStorageStorageAccounts Type = "Microsoft.Storage/storageAccounts"
)

func PossibleValuesForType() []string {
	return []string{
		string(TypeMicrosoftPointStorageStorageAccounts),
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
		"microsoft.storage/storageaccounts": TypeMicrosoftPointStorageStorageAccounts,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Type(input)
	return &out, nil
}
