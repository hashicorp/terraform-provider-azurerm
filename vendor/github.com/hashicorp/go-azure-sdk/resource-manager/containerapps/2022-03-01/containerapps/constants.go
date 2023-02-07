package containerapps

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

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

type IngressTransportMethod string

const (
	IngressTransportMethodAuto    IngressTransportMethod = "auto"
	IngressTransportMethodHTTP    IngressTransportMethod = "http"
	IngressTransportMethodHTTPTwo IngressTransportMethod = "http2"
)

func PossibleValuesForIngressTransportMethod() []string {
	return []string{
		string(IngressTransportMethodAuto),
		string(IngressTransportMethodHTTP),
		string(IngressTransportMethodHTTPTwo),
	}
}

func parseIngressTransportMethod(input string) (*IngressTransportMethod, error) {
	vals := map[string]IngressTransportMethod{
		"auto":  IngressTransportMethodAuto,
		"http":  IngressTransportMethodHTTP,
		"http2": IngressTransportMethodHTTPTwo,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IngressTransportMethod(input)
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
	StorageTypeAzureFile StorageType = "AzureFile"
	StorageTypeEmptyDir  StorageType = "EmptyDir"
)

func PossibleValuesForStorageType() []string {
	return []string{
		string(StorageTypeAzureFile),
		string(StorageTypeEmptyDir),
	}
}

func parseStorageType(input string) (*StorageType, error) {
	vals := map[string]StorageType{
		"azurefile": StorageTypeAzureFile,
		"emptydir":  StorageTypeEmptyDir,
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
