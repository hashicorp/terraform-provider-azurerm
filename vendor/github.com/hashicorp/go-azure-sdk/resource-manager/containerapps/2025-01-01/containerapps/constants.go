package containerapps

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Action string

const (
	ActionAllow Action = "Allow"
	ActionDeny  Action = "Deny"
)

func PossibleValuesForAction() []string {
	return []string{
		string(ActionAllow),
		string(ActionDeny),
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
		"deny":  ActionDeny,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Action(input)
	return &out, nil
}

type ActiveRevisionsMode string

const (
	ActiveRevisionsModeMultiple ActiveRevisionsMode = "Multiple"
	ActiveRevisionsModeSingle   ActiveRevisionsMode = "Single"
)

func PossibleValuesForActiveRevisionsMode() []string {
	return []string{
		string(ActiveRevisionsModeMultiple),
		string(ActiveRevisionsModeSingle),
	}
}

func (s *ActiveRevisionsMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseActiveRevisionsMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseActiveRevisionsMode(input string) (*ActiveRevisionsMode, error) {
	vals := map[string]ActiveRevisionsMode{
		"multiple": ActiveRevisionsModeMultiple,
		"single":   ActiveRevisionsModeSingle,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ActiveRevisionsMode(input)
	return &out, nil
}

type Affinity string

const (
	AffinityNone   Affinity = "none"
	AffinitySticky Affinity = "sticky"
)

func PossibleValuesForAffinity() []string {
	return []string{
		string(AffinityNone),
		string(AffinitySticky),
	}
}

func (s *Affinity) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAffinity(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAffinity(input string) (*Affinity, error) {
	vals := map[string]Affinity{
		"none":   AffinityNone,
		"sticky": AffinitySticky,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Affinity(input)
	return &out, nil
}

type AppProtocol string

const (
	AppProtocolGrpc AppProtocol = "grpc"
	AppProtocolHTTP AppProtocol = "http"
)

func PossibleValuesForAppProtocol() []string {
	return []string{
		string(AppProtocolGrpc),
		string(AppProtocolHTTP),
	}
}

func (s *AppProtocol) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAppProtocol(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAppProtocol(input string) (*AppProtocol, error) {
	vals := map[string]AppProtocol{
		"grpc": AppProtocolGrpc,
		"http": AppProtocolHTTP,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AppProtocol(input)
	return &out, nil
}

type BindingType string

const (
	BindingTypeDisabled   BindingType = "Disabled"
	BindingTypeSniEnabled BindingType = "SniEnabled"
)

func PossibleValuesForBindingType() []string {
	return []string{
		string(BindingTypeDisabled),
		string(BindingTypeSniEnabled),
	}
}

func (s *BindingType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBindingType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBindingType(input string) (*BindingType, error) {
	vals := map[string]BindingType{
		"disabled":   BindingTypeDisabled,
		"snienabled": BindingTypeSniEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BindingType(input)
	return &out, nil
}

type ContainerAppProvisioningState string

const (
	ContainerAppProvisioningStateCanceled   ContainerAppProvisioningState = "Canceled"
	ContainerAppProvisioningStateDeleting   ContainerAppProvisioningState = "Deleting"
	ContainerAppProvisioningStateFailed     ContainerAppProvisioningState = "Failed"
	ContainerAppProvisioningStateInProgress ContainerAppProvisioningState = "InProgress"
	ContainerAppProvisioningStateSucceeded  ContainerAppProvisioningState = "Succeeded"
)

func PossibleValuesForContainerAppProvisioningState() []string {
	return []string{
		string(ContainerAppProvisioningStateCanceled),
		string(ContainerAppProvisioningStateDeleting),
		string(ContainerAppProvisioningStateFailed),
		string(ContainerAppProvisioningStateInProgress),
		string(ContainerAppProvisioningStateSucceeded),
	}
}

func (s *ContainerAppProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseContainerAppProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseContainerAppProvisioningState(input string) (*ContainerAppProvisioningState, error) {
	vals := map[string]ContainerAppProvisioningState{
		"canceled":   ContainerAppProvisioningStateCanceled,
		"deleting":   ContainerAppProvisioningStateDeleting,
		"failed":     ContainerAppProvisioningStateFailed,
		"inprogress": ContainerAppProvisioningStateInProgress,
		"succeeded":  ContainerAppProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ContainerAppProvisioningState(input)
	return &out, nil
}

type ContainerAppRunningStatus string

const (
	ContainerAppRunningStatusProgressing ContainerAppRunningStatus = "Progressing"
	ContainerAppRunningStatusReady       ContainerAppRunningStatus = "Ready"
	ContainerAppRunningStatusRunning     ContainerAppRunningStatus = "Running"
	ContainerAppRunningStatusStopped     ContainerAppRunningStatus = "Stopped"
	ContainerAppRunningStatusSuspended   ContainerAppRunningStatus = "Suspended"
)

func PossibleValuesForContainerAppRunningStatus() []string {
	return []string{
		string(ContainerAppRunningStatusProgressing),
		string(ContainerAppRunningStatusReady),
		string(ContainerAppRunningStatusRunning),
		string(ContainerAppRunningStatusStopped),
		string(ContainerAppRunningStatusSuspended),
	}
}

func (s *ContainerAppRunningStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseContainerAppRunningStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseContainerAppRunningStatus(input string) (*ContainerAppRunningStatus, error) {
	vals := map[string]ContainerAppRunningStatus{
		"progressing": ContainerAppRunningStatusProgressing,
		"ready":       ContainerAppRunningStatusReady,
		"running":     ContainerAppRunningStatusRunning,
		"stopped":     ContainerAppRunningStatusStopped,
		"suspended":   ContainerAppRunningStatusSuspended,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ContainerAppRunningStatus(input)
	return &out, nil
}

type DnsVerificationTestResult string

const (
	DnsVerificationTestResultFailed  DnsVerificationTestResult = "Failed"
	DnsVerificationTestResultPassed  DnsVerificationTestResult = "Passed"
	DnsVerificationTestResultSkipped DnsVerificationTestResult = "Skipped"
)

func PossibleValuesForDnsVerificationTestResult() []string {
	return []string{
		string(DnsVerificationTestResultFailed),
		string(DnsVerificationTestResultPassed),
		string(DnsVerificationTestResultSkipped),
	}
}

func (s *DnsVerificationTestResult) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDnsVerificationTestResult(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDnsVerificationTestResult(input string) (*DnsVerificationTestResult, error) {
	vals := map[string]DnsVerificationTestResult{
		"failed":  DnsVerificationTestResultFailed,
		"passed":  DnsVerificationTestResultPassed,
		"skipped": DnsVerificationTestResultSkipped,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DnsVerificationTestResult(input)
	return &out, nil
}

type ExtendedLocationTypes string

const (
	ExtendedLocationTypesCustomLocation ExtendedLocationTypes = "CustomLocation"
)

func PossibleValuesForExtendedLocationTypes() []string {
	return []string{
		string(ExtendedLocationTypesCustomLocation),
	}
}

func (s *ExtendedLocationTypes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseExtendedLocationTypes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseExtendedLocationTypes(input string) (*ExtendedLocationTypes, error) {
	vals := map[string]ExtendedLocationTypes{
		"customlocation": ExtendedLocationTypesCustomLocation,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ExtendedLocationTypes(input)
	return &out, nil
}

type IdentitySettingsLifeCycle string

const (
	IdentitySettingsLifeCycleAll  IdentitySettingsLifeCycle = "All"
	IdentitySettingsLifeCycleInit IdentitySettingsLifeCycle = "Init"
	IdentitySettingsLifeCycleMain IdentitySettingsLifeCycle = "Main"
	IdentitySettingsLifeCycleNone IdentitySettingsLifeCycle = "None"
)

func PossibleValuesForIdentitySettingsLifeCycle() []string {
	return []string{
		string(IdentitySettingsLifeCycleAll),
		string(IdentitySettingsLifeCycleInit),
		string(IdentitySettingsLifeCycleMain),
		string(IdentitySettingsLifeCycleNone),
	}
}

func (s *IdentitySettingsLifeCycle) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIdentitySettingsLifeCycle(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIdentitySettingsLifeCycle(input string) (*IdentitySettingsLifeCycle, error) {
	vals := map[string]IdentitySettingsLifeCycle{
		"all":  IdentitySettingsLifeCycleAll,
		"init": IdentitySettingsLifeCycleInit,
		"main": IdentitySettingsLifeCycleMain,
		"none": IdentitySettingsLifeCycleNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IdentitySettingsLifeCycle(input)
	return &out, nil
}

type IngressClientCertificateMode string

const (
	IngressClientCertificateModeAccept  IngressClientCertificateMode = "accept"
	IngressClientCertificateModeIgnore  IngressClientCertificateMode = "ignore"
	IngressClientCertificateModeRequire IngressClientCertificateMode = "require"
)

func PossibleValuesForIngressClientCertificateMode() []string {
	return []string{
		string(IngressClientCertificateModeAccept),
		string(IngressClientCertificateModeIgnore),
		string(IngressClientCertificateModeRequire),
	}
}

func (s *IngressClientCertificateMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIngressClientCertificateMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIngressClientCertificateMode(input string) (*IngressClientCertificateMode, error) {
	vals := map[string]IngressClientCertificateMode{
		"accept":  IngressClientCertificateModeAccept,
		"ignore":  IngressClientCertificateModeIgnore,
		"require": IngressClientCertificateModeRequire,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IngressClientCertificateMode(input)
	return &out, nil
}

type IngressTransportMethod string

const (
	IngressTransportMethodAuto    IngressTransportMethod = "auto"
	IngressTransportMethodHTTP    IngressTransportMethod = "http"
	IngressTransportMethodHTTPTwo IngressTransportMethod = "http2"
	IngressTransportMethodTcp     IngressTransportMethod = "tcp"
)

func PossibleValuesForIngressTransportMethod() []string {
	return []string{
		string(IngressTransportMethodAuto),
		string(IngressTransportMethodHTTP),
		string(IngressTransportMethodHTTPTwo),
		string(IngressTransportMethodTcp),
	}
}

func (s *IngressTransportMethod) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIngressTransportMethod(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIngressTransportMethod(input string) (*IngressTransportMethod, error) {
	vals := map[string]IngressTransportMethod{
		"auto":  IngressTransportMethodAuto,
		"http":  IngressTransportMethodHTTP,
		"http2": IngressTransportMethodHTTPTwo,
		"tcp":   IngressTransportMethodTcp,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IngressTransportMethod(input)
	return &out, nil
}

type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelError LogLevel = "error"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
)

func PossibleValuesForLogLevel() []string {
	return []string{
		string(LogLevelDebug),
		string(LogLevelError),
		string(LogLevelInfo),
		string(LogLevelWarn),
	}
}

func (s *LogLevel) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLogLevel(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLogLevel(input string) (*LogLevel, error) {
	vals := map[string]LogLevel{
		"debug": LogLevelDebug,
		"error": LogLevelError,
		"info":  LogLevelInfo,
		"warn":  LogLevelWarn,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LogLevel(input)
	return &out, nil
}

type RevisionHealthState string

const (
	RevisionHealthStateHealthy   RevisionHealthState = "Healthy"
	RevisionHealthStateNone      RevisionHealthState = "None"
	RevisionHealthStateUnhealthy RevisionHealthState = "Unhealthy"
)

func PossibleValuesForRevisionHealthState() []string {
	return []string{
		string(RevisionHealthStateHealthy),
		string(RevisionHealthStateNone),
		string(RevisionHealthStateUnhealthy),
	}
}

func (s *RevisionHealthState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRevisionHealthState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRevisionHealthState(input string) (*RevisionHealthState, error) {
	vals := map[string]RevisionHealthState{
		"healthy":   RevisionHealthStateHealthy,
		"none":      RevisionHealthStateNone,
		"unhealthy": RevisionHealthStateUnhealthy,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RevisionHealthState(input)
	return &out, nil
}

type RevisionProvisioningState string

const (
	RevisionProvisioningStateDeprovisioned  RevisionProvisioningState = "Deprovisioned"
	RevisionProvisioningStateDeprovisioning RevisionProvisioningState = "Deprovisioning"
	RevisionProvisioningStateFailed         RevisionProvisioningState = "Failed"
	RevisionProvisioningStateProvisioned    RevisionProvisioningState = "Provisioned"
	RevisionProvisioningStateProvisioning   RevisionProvisioningState = "Provisioning"
)

func PossibleValuesForRevisionProvisioningState() []string {
	return []string{
		string(RevisionProvisioningStateDeprovisioned),
		string(RevisionProvisioningStateDeprovisioning),
		string(RevisionProvisioningStateFailed),
		string(RevisionProvisioningStateProvisioned),
		string(RevisionProvisioningStateProvisioning),
	}
}

func (s *RevisionProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRevisionProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRevisionProvisioningState(input string) (*RevisionProvisioningState, error) {
	vals := map[string]RevisionProvisioningState{
		"deprovisioned":  RevisionProvisioningStateDeprovisioned,
		"deprovisioning": RevisionProvisioningStateDeprovisioning,
		"failed":         RevisionProvisioningStateFailed,
		"provisioned":    RevisionProvisioningStateProvisioned,
		"provisioning":   RevisionProvisioningStateProvisioning,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RevisionProvisioningState(input)
	return &out, nil
}

type RevisionRunningState string

const (
	RevisionRunningStateDegraded   RevisionRunningState = "Degraded"
	RevisionRunningStateFailed     RevisionRunningState = "Failed"
	RevisionRunningStateProcessing RevisionRunningState = "Processing"
	RevisionRunningStateRunning    RevisionRunningState = "Running"
	RevisionRunningStateStopped    RevisionRunningState = "Stopped"
	RevisionRunningStateUnknown    RevisionRunningState = "Unknown"
)

func PossibleValuesForRevisionRunningState() []string {
	return []string{
		string(RevisionRunningStateDegraded),
		string(RevisionRunningStateFailed),
		string(RevisionRunningStateProcessing),
		string(RevisionRunningStateRunning),
		string(RevisionRunningStateStopped),
		string(RevisionRunningStateUnknown),
	}
}

func (s *RevisionRunningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRevisionRunningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRevisionRunningState(input string) (*RevisionRunningState, error) {
	vals := map[string]RevisionRunningState{
		"degraded":   RevisionRunningStateDegraded,
		"failed":     RevisionRunningStateFailed,
		"processing": RevisionRunningStateProcessing,
		"running":    RevisionRunningStateRunning,
		"stopped":    RevisionRunningStateStopped,
		"unknown":    RevisionRunningStateUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RevisionRunningState(input)
	return &out, nil
}

type Scheme string

const (
	SchemeHTTP  Scheme = "HTTP"
	SchemeHTTPS Scheme = "HTTPS"
)

func PossibleValuesForScheme() []string {
	return []string{
		string(SchemeHTTP),
		string(SchemeHTTPS),
	}
}

func (s *Scheme) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseScheme(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseScheme(input string) (*Scheme, error) {
	vals := map[string]Scheme{
		"http":  SchemeHTTP,
		"https": SchemeHTTPS,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Scheme(input)
	return &out, nil
}

type StorageType string

const (
	StorageTypeAzureFile    StorageType = "AzureFile"
	StorageTypeEmptyDir     StorageType = "EmptyDir"
	StorageTypeNfsAzureFile StorageType = "NfsAzureFile"
	StorageTypeSecret       StorageType = "Secret"
)

func PossibleValuesForStorageType() []string {
	return []string{
		string(StorageTypeAzureFile),
		string(StorageTypeEmptyDir),
		string(StorageTypeNfsAzureFile),
		string(StorageTypeSecret),
	}
}

func (s *StorageType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStorageType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStorageType(input string) (*StorageType, error) {
	vals := map[string]StorageType{
		"azurefile":    StorageTypeAzureFile,
		"emptydir":     StorageTypeEmptyDir,
		"nfsazurefile": StorageTypeNfsAzureFile,
		"secret":       StorageTypeSecret,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StorageType(input)
	return &out, nil
}

type Type string

const (
	TypeLiveness  Type = "Liveness"
	TypeReadiness Type = "Readiness"
	TypeStartup   Type = "Startup"
)

func PossibleValuesForType() []string {
	return []string{
		string(TypeLiveness),
		string(TypeReadiness),
		string(TypeStartup),
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
		"liveness":  TypeLiveness,
		"readiness": TypeReadiness,
		"startup":   TypeStartup,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Type(input)
	return &out, nil
}
