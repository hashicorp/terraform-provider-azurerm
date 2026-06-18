package integrationruntime

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataFlowComputeType string

const (
	DataFlowComputeTypeComputeOptimized DataFlowComputeType = "ComputeOptimized"
	DataFlowComputeTypeGeneral          DataFlowComputeType = "General"
	DataFlowComputeTypeMemoryOptimized  DataFlowComputeType = "MemoryOptimized"
)

func PossibleValuesForDataFlowComputeType() []string {
	return []string{
		string(DataFlowComputeTypeComputeOptimized),
		string(DataFlowComputeTypeGeneral),
		string(DataFlowComputeTypeMemoryOptimized),
	}
}

func (s *DataFlowComputeType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDataFlowComputeType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDataFlowComputeType(input string) (*DataFlowComputeType, error) {
	vals := map[string]DataFlowComputeType{
		"computeoptimized": DataFlowComputeTypeComputeOptimized,
		"general":          DataFlowComputeTypeGeneral,
		"memoryoptimized":  DataFlowComputeTypeMemoryOptimized,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DataFlowComputeType(input)
	return &out, nil
}

type IntegrationRuntimeAuthKeyName string

const (
	IntegrationRuntimeAuthKeyNameAuthKeyOne IntegrationRuntimeAuthKeyName = "authKey1"
	IntegrationRuntimeAuthKeyNameAuthKeyTwo IntegrationRuntimeAuthKeyName = "authKey2"
)

func PossibleValuesForIntegrationRuntimeAuthKeyName() []string {
	return []string{
		string(IntegrationRuntimeAuthKeyNameAuthKeyOne),
		string(IntegrationRuntimeAuthKeyNameAuthKeyTwo),
	}
}

func (s *IntegrationRuntimeAuthKeyName) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIntegrationRuntimeAuthKeyName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIntegrationRuntimeAuthKeyName(input string) (*IntegrationRuntimeAuthKeyName, error) {
	vals := map[string]IntegrationRuntimeAuthKeyName{
		"authkey1": IntegrationRuntimeAuthKeyNameAuthKeyOne,
		"authkey2": IntegrationRuntimeAuthKeyNameAuthKeyTwo,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IntegrationRuntimeAuthKeyName(input)
	return &out, nil
}

type IntegrationRuntimeAutoUpdate string

const (
	IntegrationRuntimeAutoUpdateOff IntegrationRuntimeAutoUpdate = "Off"
	IntegrationRuntimeAutoUpdateOn  IntegrationRuntimeAutoUpdate = "On"
)

func PossibleValuesForIntegrationRuntimeAutoUpdate() []string {
	return []string{
		string(IntegrationRuntimeAutoUpdateOff),
		string(IntegrationRuntimeAutoUpdateOn),
	}
}

func (s *IntegrationRuntimeAutoUpdate) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIntegrationRuntimeAutoUpdate(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIntegrationRuntimeAutoUpdate(input string) (*IntegrationRuntimeAutoUpdate, error) {
	vals := map[string]IntegrationRuntimeAutoUpdate{
		"off": IntegrationRuntimeAutoUpdateOff,
		"on":  IntegrationRuntimeAutoUpdateOn,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IntegrationRuntimeAutoUpdate(input)
	return &out, nil
}

type IntegrationRuntimeEdition string

const (
	IntegrationRuntimeEditionEnterprise IntegrationRuntimeEdition = "Enterprise"
	IntegrationRuntimeEditionStandard   IntegrationRuntimeEdition = "Standard"
)

func PossibleValuesForIntegrationRuntimeEdition() []string {
	return []string{
		string(IntegrationRuntimeEditionEnterprise),
		string(IntegrationRuntimeEditionStandard),
	}
}

func (s *IntegrationRuntimeEdition) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIntegrationRuntimeEdition(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIntegrationRuntimeEdition(input string) (*IntegrationRuntimeEdition, error) {
	vals := map[string]IntegrationRuntimeEdition{
		"enterprise": IntegrationRuntimeEditionEnterprise,
		"standard":   IntegrationRuntimeEditionStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IntegrationRuntimeEdition(input)
	return &out, nil
}

type IntegrationRuntimeEntityReferenceType string

const (
	IntegrationRuntimeEntityReferenceTypeIntegrationRuntimeReference IntegrationRuntimeEntityReferenceType = "IntegrationRuntimeReference"
	IntegrationRuntimeEntityReferenceTypeLinkedServiceReference      IntegrationRuntimeEntityReferenceType = "LinkedServiceReference"
)

func PossibleValuesForIntegrationRuntimeEntityReferenceType() []string {
	return []string{
		string(IntegrationRuntimeEntityReferenceTypeIntegrationRuntimeReference),
		string(IntegrationRuntimeEntityReferenceTypeLinkedServiceReference),
	}
}

func (s *IntegrationRuntimeEntityReferenceType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIntegrationRuntimeEntityReferenceType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIntegrationRuntimeEntityReferenceType(input string) (*IntegrationRuntimeEntityReferenceType, error) {
	vals := map[string]IntegrationRuntimeEntityReferenceType{
		"integrationruntimereference": IntegrationRuntimeEntityReferenceTypeIntegrationRuntimeReference,
		"linkedservicereference":      IntegrationRuntimeEntityReferenceTypeLinkedServiceReference,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IntegrationRuntimeEntityReferenceType(input)
	return &out, nil
}

type IntegrationRuntimeInternalChannelEncryptionMode string

const (
	IntegrationRuntimeInternalChannelEncryptionModeNotEncrypted IntegrationRuntimeInternalChannelEncryptionMode = "NotEncrypted"
	IntegrationRuntimeInternalChannelEncryptionModeNotSet       IntegrationRuntimeInternalChannelEncryptionMode = "NotSet"
	IntegrationRuntimeInternalChannelEncryptionModeSslEncrypted IntegrationRuntimeInternalChannelEncryptionMode = "SslEncrypted"
)

func PossibleValuesForIntegrationRuntimeInternalChannelEncryptionMode() []string {
	return []string{
		string(IntegrationRuntimeInternalChannelEncryptionModeNotEncrypted),
		string(IntegrationRuntimeInternalChannelEncryptionModeNotSet),
		string(IntegrationRuntimeInternalChannelEncryptionModeSslEncrypted),
	}
}

func (s *IntegrationRuntimeInternalChannelEncryptionMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIntegrationRuntimeInternalChannelEncryptionMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIntegrationRuntimeInternalChannelEncryptionMode(input string) (*IntegrationRuntimeInternalChannelEncryptionMode, error) {
	vals := map[string]IntegrationRuntimeInternalChannelEncryptionMode{
		"notencrypted": IntegrationRuntimeInternalChannelEncryptionModeNotEncrypted,
		"notset":       IntegrationRuntimeInternalChannelEncryptionModeNotSet,
		"sslencrypted": IntegrationRuntimeInternalChannelEncryptionModeSslEncrypted,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IntegrationRuntimeInternalChannelEncryptionMode(input)
	return &out, nil
}

type IntegrationRuntimeLicenseType string

const (
	IntegrationRuntimeLicenseTypeBasePrice       IntegrationRuntimeLicenseType = "BasePrice"
	IntegrationRuntimeLicenseTypeLicenseIncluded IntegrationRuntimeLicenseType = "LicenseIncluded"
)

func PossibleValuesForIntegrationRuntimeLicenseType() []string {
	return []string{
		string(IntegrationRuntimeLicenseTypeBasePrice),
		string(IntegrationRuntimeLicenseTypeLicenseIncluded),
	}
}

func (s *IntegrationRuntimeLicenseType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIntegrationRuntimeLicenseType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIntegrationRuntimeLicenseType(input string) (*IntegrationRuntimeLicenseType, error) {
	vals := map[string]IntegrationRuntimeLicenseType{
		"baseprice":       IntegrationRuntimeLicenseTypeBasePrice,
		"licenseincluded": IntegrationRuntimeLicenseTypeLicenseIncluded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IntegrationRuntimeLicenseType(input)
	return &out, nil
}

type IntegrationRuntimeSsisCatalogPricingTier string

const (
	IntegrationRuntimeSsisCatalogPricingTierBasic     IntegrationRuntimeSsisCatalogPricingTier = "Basic"
	IntegrationRuntimeSsisCatalogPricingTierPremium   IntegrationRuntimeSsisCatalogPricingTier = "Premium"
	IntegrationRuntimeSsisCatalogPricingTierPremiumRS IntegrationRuntimeSsisCatalogPricingTier = "PremiumRS"
	IntegrationRuntimeSsisCatalogPricingTierStandard  IntegrationRuntimeSsisCatalogPricingTier = "Standard"
)

func PossibleValuesForIntegrationRuntimeSsisCatalogPricingTier() []string {
	return []string{
		string(IntegrationRuntimeSsisCatalogPricingTierBasic),
		string(IntegrationRuntimeSsisCatalogPricingTierPremium),
		string(IntegrationRuntimeSsisCatalogPricingTierPremiumRS),
		string(IntegrationRuntimeSsisCatalogPricingTierStandard),
	}
}

func (s *IntegrationRuntimeSsisCatalogPricingTier) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIntegrationRuntimeSsisCatalogPricingTier(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIntegrationRuntimeSsisCatalogPricingTier(input string) (*IntegrationRuntimeSsisCatalogPricingTier, error) {
	vals := map[string]IntegrationRuntimeSsisCatalogPricingTier{
		"basic":     IntegrationRuntimeSsisCatalogPricingTierBasic,
		"premium":   IntegrationRuntimeSsisCatalogPricingTierPremium,
		"premiumrs": IntegrationRuntimeSsisCatalogPricingTierPremiumRS,
		"standard":  IntegrationRuntimeSsisCatalogPricingTierStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IntegrationRuntimeSsisCatalogPricingTier(input)
	return &out, nil
}

type IntegrationRuntimeState string

const (
	IntegrationRuntimeStateAccessDenied     IntegrationRuntimeState = "AccessDenied"
	IntegrationRuntimeStateInitial          IntegrationRuntimeState = "Initial"
	IntegrationRuntimeStateLimited          IntegrationRuntimeState = "Limited"
	IntegrationRuntimeStateNeedRegistration IntegrationRuntimeState = "NeedRegistration"
	IntegrationRuntimeStateOffline          IntegrationRuntimeState = "Offline"
	IntegrationRuntimeStateOnline           IntegrationRuntimeState = "Online"
	IntegrationRuntimeStateStarted          IntegrationRuntimeState = "Started"
	IntegrationRuntimeStateStarting         IntegrationRuntimeState = "Starting"
	IntegrationRuntimeStateStopped          IntegrationRuntimeState = "Stopped"
	IntegrationRuntimeStateStopping         IntegrationRuntimeState = "Stopping"
)

func PossibleValuesForIntegrationRuntimeState() []string {
	return []string{
		string(IntegrationRuntimeStateAccessDenied),
		string(IntegrationRuntimeStateInitial),
		string(IntegrationRuntimeStateLimited),
		string(IntegrationRuntimeStateNeedRegistration),
		string(IntegrationRuntimeStateOffline),
		string(IntegrationRuntimeStateOnline),
		string(IntegrationRuntimeStateStarted),
		string(IntegrationRuntimeStateStarting),
		string(IntegrationRuntimeStateStopped),
		string(IntegrationRuntimeStateStopping),
	}
}

func (s *IntegrationRuntimeState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIntegrationRuntimeState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIntegrationRuntimeState(input string) (*IntegrationRuntimeState, error) {
	vals := map[string]IntegrationRuntimeState{
		"accessdenied":     IntegrationRuntimeStateAccessDenied,
		"initial":          IntegrationRuntimeStateInitial,
		"limited":          IntegrationRuntimeStateLimited,
		"needregistration": IntegrationRuntimeStateNeedRegistration,
		"offline":          IntegrationRuntimeStateOffline,
		"online":           IntegrationRuntimeStateOnline,
		"started":          IntegrationRuntimeStateStarted,
		"starting":         IntegrationRuntimeStateStarting,
		"stopped":          IntegrationRuntimeStateStopped,
		"stopping":         IntegrationRuntimeStateStopping,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IntegrationRuntimeState(input)
	return &out, nil
}

type IntegrationRuntimeType string

const (
	IntegrationRuntimeTypeManaged    IntegrationRuntimeType = "Managed"
	IntegrationRuntimeTypeSelfHosted IntegrationRuntimeType = "SelfHosted"
)

func PossibleValuesForIntegrationRuntimeType() []string {
	return []string{
		string(IntegrationRuntimeTypeManaged),
		string(IntegrationRuntimeTypeSelfHosted),
	}
}

func (s *IntegrationRuntimeType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIntegrationRuntimeType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIntegrationRuntimeType(input string) (*IntegrationRuntimeType, error) {
	vals := map[string]IntegrationRuntimeType{
		"managed":    IntegrationRuntimeTypeManaged,
		"selfhosted": IntegrationRuntimeTypeSelfHosted,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IntegrationRuntimeType(input)
	return &out, nil
}

type IntegrationRuntimeUpdateResult string

const (
	IntegrationRuntimeUpdateResultFail    IntegrationRuntimeUpdateResult = "Fail"
	IntegrationRuntimeUpdateResultNone    IntegrationRuntimeUpdateResult = "None"
	IntegrationRuntimeUpdateResultSucceed IntegrationRuntimeUpdateResult = "Succeed"
)

func PossibleValuesForIntegrationRuntimeUpdateResult() []string {
	return []string{
		string(IntegrationRuntimeUpdateResultFail),
		string(IntegrationRuntimeUpdateResultNone),
		string(IntegrationRuntimeUpdateResultSucceed),
	}
}

func (s *IntegrationRuntimeUpdateResult) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIntegrationRuntimeUpdateResult(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIntegrationRuntimeUpdateResult(input string) (*IntegrationRuntimeUpdateResult, error) {
	vals := map[string]IntegrationRuntimeUpdateResult{
		"fail":    IntegrationRuntimeUpdateResultFail,
		"none":    IntegrationRuntimeUpdateResultNone,
		"succeed": IntegrationRuntimeUpdateResultSucceed,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IntegrationRuntimeUpdateResult(input)
	return &out, nil
}

type ManagedIntegrationRuntimeNodeStatus string

const (
	ManagedIntegrationRuntimeNodeStatusAvailable   ManagedIntegrationRuntimeNodeStatus = "Available"
	ManagedIntegrationRuntimeNodeStatusRecycling   ManagedIntegrationRuntimeNodeStatus = "Recycling"
	ManagedIntegrationRuntimeNodeStatusStarting    ManagedIntegrationRuntimeNodeStatus = "Starting"
	ManagedIntegrationRuntimeNodeStatusUnavailable ManagedIntegrationRuntimeNodeStatus = "Unavailable"
)

func PossibleValuesForManagedIntegrationRuntimeNodeStatus() []string {
	return []string{
		string(ManagedIntegrationRuntimeNodeStatusAvailable),
		string(ManagedIntegrationRuntimeNodeStatusRecycling),
		string(ManagedIntegrationRuntimeNodeStatusStarting),
		string(ManagedIntegrationRuntimeNodeStatusUnavailable),
	}
}

func (s *ManagedIntegrationRuntimeNodeStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseManagedIntegrationRuntimeNodeStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseManagedIntegrationRuntimeNodeStatus(input string) (*ManagedIntegrationRuntimeNodeStatus, error) {
	vals := map[string]ManagedIntegrationRuntimeNodeStatus{
		"available":   ManagedIntegrationRuntimeNodeStatusAvailable,
		"recycling":   ManagedIntegrationRuntimeNodeStatusRecycling,
		"starting":    ManagedIntegrationRuntimeNodeStatusStarting,
		"unavailable": ManagedIntegrationRuntimeNodeStatusUnavailable,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ManagedIntegrationRuntimeNodeStatus(input)
	return &out, nil
}

type SelfHostedIntegrationRuntimeNodeStatus string

const (
	SelfHostedIntegrationRuntimeNodeStatusInitializeFailed SelfHostedIntegrationRuntimeNodeStatus = "InitializeFailed"
	SelfHostedIntegrationRuntimeNodeStatusInitializing     SelfHostedIntegrationRuntimeNodeStatus = "Initializing"
	SelfHostedIntegrationRuntimeNodeStatusLimited          SelfHostedIntegrationRuntimeNodeStatus = "Limited"
	SelfHostedIntegrationRuntimeNodeStatusNeedRegistration SelfHostedIntegrationRuntimeNodeStatus = "NeedRegistration"
	SelfHostedIntegrationRuntimeNodeStatusOffline          SelfHostedIntegrationRuntimeNodeStatus = "Offline"
	SelfHostedIntegrationRuntimeNodeStatusOnline           SelfHostedIntegrationRuntimeNodeStatus = "Online"
	SelfHostedIntegrationRuntimeNodeStatusUpgrading        SelfHostedIntegrationRuntimeNodeStatus = "Upgrading"
)

func PossibleValuesForSelfHostedIntegrationRuntimeNodeStatus() []string {
	return []string{
		string(SelfHostedIntegrationRuntimeNodeStatusInitializeFailed),
		string(SelfHostedIntegrationRuntimeNodeStatusInitializing),
		string(SelfHostedIntegrationRuntimeNodeStatusLimited),
		string(SelfHostedIntegrationRuntimeNodeStatusNeedRegistration),
		string(SelfHostedIntegrationRuntimeNodeStatusOffline),
		string(SelfHostedIntegrationRuntimeNodeStatusOnline),
		string(SelfHostedIntegrationRuntimeNodeStatusUpgrading),
	}
}

func (s *SelfHostedIntegrationRuntimeNodeStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSelfHostedIntegrationRuntimeNodeStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSelfHostedIntegrationRuntimeNodeStatus(input string) (*SelfHostedIntegrationRuntimeNodeStatus, error) {
	vals := map[string]SelfHostedIntegrationRuntimeNodeStatus{
		"initializefailed": SelfHostedIntegrationRuntimeNodeStatusInitializeFailed,
		"initializing":     SelfHostedIntegrationRuntimeNodeStatusInitializing,
		"limited":          SelfHostedIntegrationRuntimeNodeStatusLimited,
		"needregistration": SelfHostedIntegrationRuntimeNodeStatusNeedRegistration,
		"offline":          SelfHostedIntegrationRuntimeNodeStatusOffline,
		"online":           SelfHostedIntegrationRuntimeNodeStatusOnline,
		"upgrading":        SelfHostedIntegrationRuntimeNodeStatusUpgrading,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SelfHostedIntegrationRuntimeNodeStatus(input)
	return &out, nil
}

type SsisObjectMetadataType string

const (
	SsisObjectMetadataTypeEnvironment SsisObjectMetadataType = "Environment"
	SsisObjectMetadataTypeFolder      SsisObjectMetadataType = "Folder"
	SsisObjectMetadataTypePackage     SsisObjectMetadataType = "Package"
	SsisObjectMetadataTypeProject     SsisObjectMetadataType = "Project"
)

func PossibleValuesForSsisObjectMetadataType() []string {
	return []string{
		string(SsisObjectMetadataTypeEnvironment),
		string(SsisObjectMetadataTypeFolder),
		string(SsisObjectMetadataTypePackage),
		string(SsisObjectMetadataTypeProject),
	}
}

func (s *SsisObjectMetadataType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSsisObjectMetadataType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSsisObjectMetadataType(input string) (*SsisObjectMetadataType, error) {
	vals := map[string]SsisObjectMetadataType{
		"environment": SsisObjectMetadataTypeEnvironment,
		"folder":      SsisObjectMetadataTypeFolder,
		"package":     SsisObjectMetadataTypePackage,
		"project":     SsisObjectMetadataTypeProject,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SsisObjectMetadataType(input)
	return &out, nil
}
