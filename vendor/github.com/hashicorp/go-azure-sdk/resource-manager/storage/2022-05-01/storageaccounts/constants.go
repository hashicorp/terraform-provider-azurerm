package storageaccounts

import "strings"

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
