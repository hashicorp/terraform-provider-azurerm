package devices

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AuthenticationType string

const (
	AuthenticationTypeAzureActiveDirectory AuthenticationType = "AzureActiveDirectory"
	AuthenticationTypeInvalid              AuthenticationType = "Invalid"
)

func PossibleValuesForAuthenticationType() []string {
	return []string{
		string(AuthenticationTypeAzureActiveDirectory),
		string(AuthenticationTypeInvalid),
	}
}

func parseAuthenticationType(input string) (*AuthenticationType, error) {
	vals := map[string]AuthenticationType{
		"azureactivedirectory": AuthenticationTypeAzureActiveDirectory,
		"invalid":              AuthenticationTypeInvalid,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AuthenticationType(input)
	return &out, nil
}

type DataBoxEdgeDeviceKind string

const (
	DataBoxEdgeDeviceKindAzureDataBoxGateway    DataBoxEdgeDeviceKind = "AzureDataBoxGateway"
	DataBoxEdgeDeviceKindAzureModularDataCentre DataBoxEdgeDeviceKind = "AzureModularDataCentre"
	DataBoxEdgeDeviceKindAzureStackEdge         DataBoxEdgeDeviceKind = "AzureStackEdge"
	DataBoxEdgeDeviceKindAzureStackHub          DataBoxEdgeDeviceKind = "AzureStackHub"
)

func PossibleValuesForDataBoxEdgeDeviceKind() []string {
	return []string{
		string(DataBoxEdgeDeviceKindAzureDataBoxGateway),
		string(DataBoxEdgeDeviceKindAzureModularDataCentre),
		string(DataBoxEdgeDeviceKindAzureStackEdge),
		string(DataBoxEdgeDeviceKindAzureStackHub),
	}
}

func parseDataBoxEdgeDeviceKind(input string) (*DataBoxEdgeDeviceKind, error) {
	vals := map[string]DataBoxEdgeDeviceKind{
		"azuredataboxgateway":    DataBoxEdgeDeviceKindAzureDataBoxGateway,
		"azuremodulardatacentre": DataBoxEdgeDeviceKindAzureModularDataCentre,
		"azurestackedge":         DataBoxEdgeDeviceKindAzureStackEdge,
		"azurestackhub":          DataBoxEdgeDeviceKindAzureStackHub,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DataBoxEdgeDeviceKind(input)
	return &out, nil
}

type DataBoxEdgeDeviceStatus string

const (
	DataBoxEdgeDeviceStatusDisconnected          DataBoxEdgeDeviceStatus = "Disconnected"
	DataBoxEdgeDeviceStatusMaintenance           DataBoxEdgeDeviceStatus = "Maintenance"
	DataBoxEdgeDeviceStatusNeedsAttention        DataBoxEdgeDeviceStatus = "NeedsAttention"
	DataBoxEdgeDeviceStatusOffline               DataBoxEdgeDeviceStatus = "Offline"
	DataBoxEdgeDeviceStatusOnline                DataBoxEdgeDeviceStatus = "Online"
	DataBoxEdgeDeviceStatusPartiallyDisconnected DataBoxEdgeDeviceStatus = "PartiallyDisconnected"
	DataBoxEdgeDeviceStatusReadyToSetup          DataBoxEdgeDeviceStatus = "ReadyToSetup"
)

func PossibleValuesForDataBoxEdgeDeviceStatus() []string {
	return []string{
		string(DataBoxEdgeDeviceStatusDisconnected),
		string(DataBoxEdgeDeviceStatusMaintenance),
		string(DataBoxEdgeDeviceStatusNeedsAttention),
		string(DataBoxEdgeDeviceStatusOffline),
		string(DataBoxEdgeDeviceStatusOnline),
		string(DataBoxEdgeDeviceStatusPartiallyDisconnected),
		string(DataBoxEdgeDeviceStatusReadyToSetup),
	}
}

func parseDataBoxEdgeDeviceStatus(input string) (*DataBoxEdgeDeviceStatus, error) {
	vals := map[string]DataBoxEdgeDeviceStatus{
		"disconnected":          DataBoxEdgeDeviceStatusDisconnected,
		"maintenance":           DataBoxEdgeDeviceStatusMaintenance,
		"needsattention":        DataBoxEdgeDeviceStatusNeedsAttention,
		"offline":               DataBoxEdgeDeviceStatusOffline,
		"online":                DataBoxEdgeDeviceStatusOnline,
		"partiallydisconnected": DataBoxEdgeDeviceStatusPartiallyDisconnected,
		"readytosetup":          DataBoxEdgeDeviceStatusReadyToSetup,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DataBoxEdgeDeviceStatus(input)
	return &out, nil
}

type DeviceType string

const (
	DeviceTypeDataBoxEdgeDevice DeviceType = "DataBoxEdgeDevice"
)

func PossibleValuesForDeviceType() []string {
	return []string{
		string(DeviceTypeDataBoxEdgeDevice),
	}
}

func parseDeviceType(input string) (*DeviceType, error) {
	vals := map[string]DeviceType{
		"databoxedgedevice": DeviceTypeDataBoxEdgeDevice,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DeviceType(input)
	return &out, nil
}

type EncryptionAlgorithm string

const (
	EncryptionAlgorithmAESTwoFiveSix        EncryptionAlgorithm = "AES256"
	EncryptionAlgorithmNone                 EncryptionAlgorithm = "None"
	EncryptionAlgorithmRSAESPKCSOneVOneFive EncryptionAlgorithm = "RSAES_PKCS1_v_1_5"
)

func PossibleValuesForEncryptionAlgorithm() []string {
	return []string{
		string(EncryptionAlgorithmAESTwoFiveSix),
		string(EncryptionAlgorithmNone),
		string(EncryptionAlgorithmRSAESPKCSOneVOneFive),
	}
}

func parseEncryptionAlgorithm(input string) (*EncryptionAlgorithm, error) {
	vals := map[string]EncryptionAlgorithm{
		"aes256":            EncryptionAlgorithmAESTwoFiveSix,
		"none":              EncryptionAlgorithmNone,
		"rsaes_pkcs1_v_1_5": EncryptionAlgorithmRSAESPKCSOneVOneFive,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EncryptionAlgorithm(input)
	return &out, nil
}

type InstallRebootBehavior string

const (
	InstallRebootBehaviorNeverReboots   InstallRebootBehavior = "NeverReboots"
	InstallRebootBehaviorRequestReboot  InstallRebootBehavior = "RequestReboot"
	InstallRebootBehaviorRequiresReboot InstallRebootBehavior = "RequiresReboot"
)

func PossibleValuesForInstallRebootBehavior() []string {
	return []string{
		string(InstallRebootBehaviorNeverReboots),
		string(InstallRebootBehaviorRequestReboot),
		string(InstallRebootBehaviorRequiresReboot),
	}
}

func parseInstallRebootBehavior(input string) (*InstallRebootBehavior, error) {
	vals := map[string]InstallRebootBehavior{
		"neverreboots":   InstallRebootBehaviorNeverReboots,
		"requestreboot":  InstallRebootBehaviorRequestReboot,
		"requiresreboot": InstallRebootBehaviorRequiresReboot,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := InstallRebootBehavior(input)
	return &out, nil
}

type KeyVaultSyncStatus string

const (
	KeyVaultSyncStatusKeyVaultNotConfigured KeyVaultSyncStatus = "KeyVaultNotConfigured"
	KeyVaultSyncStatusKeyVaultSyncFailed    KeyVaultSyncStatus = "KeyVaultSyncFailed"
	KeyVaultSyncStatusKeyVaultSyncPending   KeyVaultSyncStatus = "KeyVaultSyncPending"
	KeyVaultSyncStatusKeyVaultSynced        KeyVaultSyncStatus = "KeyVaultSynced"
	KeyVaultSyncStatusKeyVaultSyncing       KeyVaultSyncStatus = "KeyVaultSyncing"
)

func PossibleValuesForKeyVaultSyncStatus() []string {
	return []string{
		string(KeyVaultSyncStatusKeyVaultNotConfigured),
		string(KeyVaultSyncStatusKeyVaultSyncFailed),
		string(KeyVaultSyncStatusKeyVaultSyncPending),
		string(KeyVaultSyncStatusKeyVaultSynced),
		string(KeyVaultSyncStatusKeyVaultSyncing),
	}
}

func parseKeyVaultSyncStatus(input string) (*KeyVaultSyncStatus, error) {
	vals := map[string]KeyVaultSyncStatus{
		"keyvaultnotconfigured": KeyVaultSyncStatusKeyVaultNotConfigured,
		"keyvaultsyncfailed":    KeyVaultSyncStatusKeyVaultSyncFailed,
		"keyvaultsyncpending":   KeyVaultSyncStatusKeyVaultSyncPending,
		"keyvaultsynced":        KeyVaultSyncStatusKeyVaultSynced,
		"keyvaultsyncing":       KeyVaultSyncStatusKeyVaultSyncing,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := KeyVaultSyncStatus(input)
	return &out, nil
}

type MsiIdentityType string

const (
	MsiIdentityTypeNone           MsiIdentityType = "None"
	MsiIdentityTypeSystemAssigned MsiIdentityType = "SystemAssigned"
	MsiIdentityTypeUserAssigned   MsiIdentityType = "UserAssigned"
)

func PossibleValuesForMsiIdentityType() []string {
	return []string{
		string(MsiIdentityTypeNone),
		string(MsiIdentityTypeSystemAssigned),
		string(MsiIdentityTypeUserAssigned),
	}
}

func parseMsiIdentityType(input string) (*MsiIdentityType, error) {
	vals := map[string]MsiIdentityType{
		"none":           MsiIdentityTypeNone,
		"systemassigned": MsiIdentityTypeSystemAssigned,
		"userassigned":   MsiIdentityTypeUserAssigned,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MsiIdentityType(input)
	return &out, nil
}

type NetworkAdapterDHCPStatus string

const (
	NetworkAdapterDHCPStatusDisabled NetworkAdapterDHCPStatus = "Disabled"
	NetworkAdapterDHCPStatusEnabled  NetworkAdapterDHCPStatus = "Enabled"
)

func PossibleValuesForNetworkAdapterDHCPStatus() []string {
	return []string{
		string(NetworkAdapterDHCPStatusDisabled),
		string(NetworkAdapterDHCPStatusEnabled),
	}
}

func parseNetworkAdapterDHCPStatus(input string) (*NetworkAdapterDHCPStatus, error) {
	vals := map[string]NetworkAdapterDHCPStatus{
		"disabled": NetworkAdapterDHCPStatusDisabled,
		"enabled":  NetworkAdapterDHCPStatusEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NetworkAdapterDHCPStatus(input)
	return &out, nil
}

type NetworkAdapterRDMAStatus string

const (
	NetworkAdapterRDMAStatusCapable   NetworkAdapterRDMAStatus = "Capable"
	NetworkAdapterRDMAStatusIncapable NetworkAdapterRDMAStatus = "Incapable"
)

func PossibleValuesForNetworkAdapterRDMAStatus() []string {
	return []string{
		string(NetworkAdapterRDMAStatusCapable),
		string(NetworkAdapterRDMAStatusIncapable),
	}
}

func parseNetworkAdapterRDMAStatus(input string) (*NetworkAdapterRDMAStatus, error) {
	vals := map[string]NetworkAdapterRDMAStatus{
		"capable":   NetworkAdapterRDMAStatusCapable,
		"incapable": NetworkAdapterRDMAStatusIncapable,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NetworkAdapterRDMAStatus(input)
	return &out, nil
}

type NetworkAdapterStatus string

const (
	NetworkAdapterStatusActive   NetworkAdapterStatus = "Active"
	NetworkAdapterStatusInactive NetworkAdapterStatus = "Inactive"
)

func PossibleValuesForNetworkAdapterStatus() []string {
	return []string{
		string(NetworkAdapterStatusActive),
		string(NetworkAdapterStatusInactive),
	}
}

func parseNetworkAdapterStatus(input string) (*NetworkAdapterStatus, error) {
	vals := map[string]NetworkAdapterStatus{
		"active":   NetworkAdapterStatusActive,
		"inactive": NetworkAdapterStatusInactive,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NetworkAdapterStatus(input)
	return &out, nil
}

type NetworkGroup string

const (
	NetworkGroupNonRDMA NetworkGroup = "NonRDMA"
	NetworkGroupNone    NetworkGroup = "None"
	NetworkGroupRDMA    NetworkGroup = "RDMA"
)

func PossibleValuesForNetworkGroup() []string {
	return []string{
		string(NetworkGroupNonRDMA),
		string(NetworkGroupNone),
		string(NetworkGroupRDMA),
	}
}

func parseNetworkGroup(input string) (*NetworkGroup, error) {
	vals := map[string]NetworkGroup{
		"nonrdma": NetworkGroupNonRDMA,
		"none":    NetworkGroupNone,
		"rdma":    NetworkGroupRDMA,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NetworkGroup(input)
	return &out, nil
}

type ResourceMoveStatus string

const (
	ResourceMoveStatusNone                   ResourceMoveStatus = "None"
	ResourceMoveStatusResourceMoveFailed     ResourceMoveStatus = "ResourceMoveFailed"
	ResourceMoveStatusResourceMoveInProgress ResourceMoveStatus = "ResourceMoveInProgress"
)

func PossibleValuesForResourceMoveStatus() []string {
	return []string{
		string(ResourceMoveStatusNone),
		string(ResourceMoveStatusResourceMoveFailed),
		string(ResourceMoveStatusResourceMoveInProgress),
	}
}

func parseResourceMoveStatus(input string) (*ResourceMoveStatus, error) {
	vals := map[string]ResourceMoveStatus{
		"none":                   ResourceMoveStatusNone,
		"resourcemovefailed":     ResourceMoveStatusResourceMoveFailed,
		"resourcemoveinprogress": ResourceMoveStatusResourceMoveInProgress,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ResourceMoveStatus(input)
	return &out, nil
}

type RoleTypes string

const (
	RoleTypesASA                 RoleTypes = "ASA"
	RoleTypesCloudEdgeManagement RoleTypes = "CloudEdgeManagement"
	RoleTypesCognitive           RoleTypes = "Cognitive"
	RoleTypesFunctions           RoleTypes = "Functions"
	RoleTypesIOT                 RoleTypes = "IOT"
	RoleTypesKubernetes          RoleTypes = "Kubernetes"
	RoleTypesMEC                 RoleTypes = "MEC"
)

func PossibleValuesForRoleTypes() []string {
	return []string{
		string(RoleTypesASA),
		string(RoleTypesCloudEdgeManagement),
		string(RoleTypesCognitive),
		string(RoleTypesFunctions),
		string(RoleTypesIOT),
		string(RoleTypesKubernetes),
		string(RoleTypesMEC),
	}
}

func parseRoleTypes(input string) (*RoleTypes, error) {
	vals := map[string]RoleTypes{
		"asa":                 RoleTypesASA,
		"cloudedgemanagement": RoleTypesCloudEdgeManagement,
		"cognitive":           RoleTypesCognitive,
		"functions":           RoleTypesFunctions,
		"iot":                 RoleTypesIOT,
		"kubernetes":          RoleTypesKubernetes,
		"mec":                 RoleTypesMEC,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RoleTypes(input)
	return &out, nil
}

type SkuName string

const (
	SkuNameEdge                 SkuName = "Edge"
	SkuNameEdgeMRMini           SkuName = "EdgeMR_Mini"
	SkuNameEdgePBase            SkuName = "EdgeP_Base"
	SkuNameEdgePHigh            SkuName = "EdgeP_High"
	SkuNameEdgePRBase           SkuName = "EdgePR_Base"
	SkuNameEdgePRBaseUPS        SkuName = "EdgePR_Base_UPS"
	SkuNameGPU                  SkuName = "GPU"
	SkuNameGateway              SkuName = "Gateway"
	SkuNameRCALarge             SkuName = "RCA_Large"
	SkuNameRCASmall             SkuName = "RCA_Small"
	SkuNameRDC                  SkuName = "RDC"
	SkuNameTCALarge             SkuName = "TCA_Large"
	SkuNameTCASmall             SkuName = "TCA_Small"
	SkuNameTDC                  SkuName = "TDC"
	SkuNameTEAFourNodeHeater    SkuName = "TEA_4Node_Heater"
	SkuNameTEAFourNodeUPSHeater SkuName = "TEA_4Node_UPS_Heater"
	SkuNameTEAOneNode           SkuName = "TEA_1Node"
	SkuNameTEAOneNodeHeater     SkuName = "TEA_1Node_Heater"
	SkuNameTEAOneNodeUPS        SkuName = "TEA_1Node_UPS"
	SkuNameTEAOneNodeUPSHeater  SkuName = "TEA_1Node_UPS_Heater"
	SkuNameTMA                  SkuName = "TMA"
)

func PossibleValuesForSkuName() []string {
	return []string{
		string(SkuNameEdge),
		string(SkuNameEdgeMRMini),
		string(SkuNameEdgePBase),
		string(SkuNameEdgePHigh),
		string(SkuNameEdgePRBase),
		string(SkuNameEdgePRBaseUPS),
		string(SkuNameGPU),
		string(SkuNameGateway),
		string(SkuNameRCALarge),
		string(SkuNameRCASmall),
		string(SkuNameRDC),
		string(SkuNameTCALarge),
		string(SkuNameTCASmall),
		string(SkuNameTDC),
		string(SkuNameTEAFourNodeHeater),
		string(SkuNameTEAFourNodeUPSHeater),
		string(SkuNameTEAOneNode),
		string(SkuNameTEAOneNodeHeater),
		string(SkuNameTEAOneNodeUPS),
		string(SkuNameTEAOneNodeUPSHeater),
		string(SkuNameTMA),
	}
}

func parseSkuName(input string) (*SkuName, error) {
	vals := map[string]SkuName{
		"edge":                 SkuNameEdge,
		"edgemr_mini":          SkuNameEdgeMRMini,
		"edgep_base":           SkuNameEdgePBase,
		"edgep_high":           SkuNameEdgePHigh,
		"edgepr_base":          SkuNameEdgePRBase,
		"edgepr_base_ups":      SkuNameEdgePRBaseUPS,
		"gpu":                  SkuNameGPU,
		"gateway":              SkuNameGateway,
		"rca_large":            SkuNameRCALarge,
		"rca_small":            SkuNameRCASmall,
		"rdc":                  SkuNameRDC,
		"tca_large":            SkuNameTCALarge,
		"tca_small":            SkuNameTCASmall,
		"tdc":                  SkuNameTDC,
		"tea_4node_heater":     SkuNameTEAFourNodeHeater,
		"tea_4node_ups_heater": SkuNameTEAFourNodeUPSHeater,
		"tea_1node":            SkuNameTEAOneNode,
		"tea_1node_heater":     SkuNameTEAOneNodeHeater,
		"tea_1node_ups":        SkuNameTEAOneNodeUPS,
		"tea_1node_ups_heater": SkuNameTEAOneNodeUPSHeater,
		"tma":                  SkuNameTMA,
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
	SkuTierStandard SkuTier = "Standard"
)

func PossibleValuesForSkuTier() []string {
	return []string{
		string(SkuTierStandard),
	}
}

func parseSkuTier(input string) (*SkuTier, error) {
	vals := map[string]SkuTier{
		"standard": SkuTierStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SkuTier(input)
	return &out, nil
}

type SubscriptionState string

const (
	SubscriptionStateDeleted      SubscriptionState = "Deleted"
	SubscriptionStateRegistered   SubscriptionState = "Registered"
	SubscriptionStateSuspended    SubscriptionState = "Suspended"
	SubscriptionStateUnregistered SubscriptionState = "Unregistered"
	SubscriptionStateWarned       SubscriptionState = "Warned"
)

func PossibleValuesForSubscriptionState() []string {
	return []string{
		string(SubscriptionStateDeleted),
		string(SubscriptionStateRegistered),
		string(SubscriptionStateSuspended),
		string(SubscriptionStateUnregistered),
		string(SubscriptionStateWarned),
	}
}

func parseSubscriptionState(input string) (*SubscriptionState, error) {
	vals := map[string]SubscriptionState{
		"deleted":      SubscriptionStateDeleted,
		"registered":   SubscriptionStateRegistered,
		"suspended":    SubscriptionStateSuspended,
		"unregistered": SubscriptionStateUnregistered,
		"warned":       SubscriptionStateWarned,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SubscriptionState(input)
	return &out, nil
}

type UpdateOperation string

const (
	UpdateOperationDownload UpdateOperation = "Download"
	UpdateOperationInstall  UpdateOperation = "Install"
	UpdateOperationNone     UpdateOperation = "None"
	UpdateOperationScan     UpdateOperation = "Scan"
)

func PossibleValuesForUpdateOperation() []string {
	return []string{
		string(UpdateOperationDownload),
		string(UpdateOperationInstall),
		string(UpdateOperationNone),
		string(UpdateOperationScan),
	}
}

func parseUpdateOperation(input string) (*UpdateOperation, error) {
	vals := map[string]UpdateOperation{
		"download": UpdateOperationDownload,
		"install":  UpdateOperationInstall,
		"none":     UpdateOperationNone,
		"scan":     UpdateOperationScan,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := UpdateOperation(input)
	return &out, nil
}
