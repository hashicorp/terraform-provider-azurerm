package clusters

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterNodeType string

const (
	ClusterNodeTypeFirstParty ClusterNodeType = "FirstParty"
	ClusterNodeTypeThirdParty ClusterNodeType = "ThirdParty"
)

func PossibleValuesForClusterNodeType() []string {
	return []string{
		string(ClusterNodeTypeFirstParty),
		string(ClusterNodeTypeThirdParty),
	}
}

func parseClusterNodeType(input string) (*ClusterNodeType, error) {
	vals := map[string]ClusterNodeType{
		"firstparty": ClusterNodeTypeFirstParty,
		"thirdparty": ClusterNodeTypeThirdParty,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ClusterNodeType(input)
	return &out, nil
}

type DiagnosticLevel string

const (
	DiagnosticLevelBasic    DiagnosticLevel = "Basic"
	DiagnosticLevelEnhanced DiagnosticLevel = "Enhanced"
	DiagnosticLevelOff      DiagnosticLevel = "Off"
)

func PossibleValuesForDiagnosticLevel() []string {
	return []string{
		string(DiagnosticLevelBasic),
		string(DiagnosticLevelEnhanced),
		string(DiagnosticLevelOff),
	}
}

func parseDiagnosticLevel(input string) (*DiagnosticLevel, error) {
	vals := map[string]DiagnosticLevel{
		"basic":    DiagnosticLevelBasic,
		"enhanced": DiagnosticLevelEnhanced,
		"off":      DiagnosticLevelOff,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DiagnosticLevel(input)
	return &out, nil
}

type ImdsAttestation string

const (
	ImdsAttestationDisabled ImdsAttestation = "Disabled"
	ImdsAttestationEnabled  ImdsAttestation = "Enabled"
)

func PossibleValuesForImdsAttestation() []string {
	return []string{
		string(ImdsAttestationDisabled),
		string(ImdsAttestationEnabled),
	}
}

func parseImdsAttestation(input string) (*ImdsAttestation, error) {
	vals := map[string]ImdsAttestation{
		"disabled": ImdsAttestationDisabled,
		"enabled":  ImdsAttestationEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ImdsAttestation(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateAccepted     ProvisioningState = "Accepted"
	ProvisioningStateCanceled     ProvisioningState = "Canceled"
	ProvisioningStateFailed       ProvisioningState = "Failed"
	ProvisioningStateProvisioning ProvisioningState = "Provisioning"
	ProvisioningStateSucceeded    ProvisioningState = "Succeeded"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateAccepted),
		string(ProvisioningStateCanceled),
		string(ProvisioningStateFailed),
		string(ProvisioningStateProvisioning),
		string(ProvisioningStateSucceeded),
	}
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"accepted":     ProvisioningStateAccepted,
		"canceled":     ProvisioningStateCanceled,
		"failed":       ProvisioningStateFailed,
		"provisioning": ProvisioningStateProvisioning,
		"succeeded":    ProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}

type SoftwareAssuranceIntent string

const (
	SoftwareAssuranceIntentDisable SoftwareAssuranceIntent = "Disable"
	SoftwareAssuranceIntentEnable  SoftwareAssuranceIntent = "Enable"
)

func PossibleValuesForSoftwareAssuranceIntent() []string {
	return []string{
		string(SoftwareAssuranceIntentDisable),
		string(SoftwareAssuranceIntentEnable),
	}
}

func parseSoftwareAssuranceIntent(input string) (*SoftwareAssuranceIntent, error) {
	vals := map[string]SoftwareAssuranceIntent{
		"disable": SoftwareAssuranceIntentDisable,
		"enable":  SoftwareAssuranceIntentEnable,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SoftwareAssuranceIntent(input)
	return &out, nil
}

type SoftwareAssuranceStatus string

const (
	SoftwareAssuranceStatusDisabled SoftwareAssuranceStatus = "Disabled"
	SoftwareAssuranceStatusEnabled  SoftwareAssuranceStatus = "Enabled"
)

func PossibleValuesForSoftwareAssuranceStatus() []string {
	return []string{
		string(SoftwareAssuranceStatusDisabled),
		string(SoftwareAssuranceStatusEnabled),
	}
}

func parseSoftwareAssuranceStatus(input string) (*SoftwareAssuranceStatus, error) {
	vals := map[string]SoftwareAssuranceStatus{
		"disabled": SoftwareAssuranceStatusDisabled,
		"enabled":  SoftwareAssuranceStatusEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SoftwareAssuranceStatus(input)
	return &out, nil
}

type Status string

const (
	StatusConnectedRecently    Status = "ConnectedRecently"
	StatusDisconnected         Status = "Disconnected"
	StatusError                Status = "Error"
	StatusNotConnectedRecently Status = "NotConnectedRecently"
	StatusNotYetRegistered     Status = "NotYetRegistered"
)

func PossibleValuesForStatus() []string {
	return []string{
		string(StatusConnectedRecently),
		string(StatusDisconnected),
		string(StatusError),
		string(StatusNotConnectedRecently),
		string(StatusNotYetRegistered),
	}
}

func parseStatus(input string) (*Status, error) {
	vals := map[string]Status{
		"connectedrecently":    StatusConnectedRecently,
		"disconnected":         StatusDisconnected,
		"error":                StatusError,
		"notconnectedrecently": StatusNotConnectedRecently,
		"notyetregistered":     StatusNotYetRegistered,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Status(input)
	return &out, nil
}

type WindowsServerSubscription string

const (
	WindowsServerSubscriptionDisabled WindowsServerSubscription = "Disabled"
	WindowsServerSubscriptionEnabled  WindowsServerSubscription = "Enabled"
)

func PossibleValuesForWindowsServerSubscription() []string {
	return []string{
		string(WindowsServerSubscriptionDisabled),
		string(WindowsServerSubscriptionEnabled),
	}
}

func parseWindowsServerSubscription(input string) (*WindowsServerSubscription, error) {
	vals := map[string]WindowsServerSubscription{
		"disabled": WindowsServerSubscriptionDisabled,
		"enabled":  WindowsServerSubscriptionEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := WindowsServerSubscription(input)
	return &out, nil
}
