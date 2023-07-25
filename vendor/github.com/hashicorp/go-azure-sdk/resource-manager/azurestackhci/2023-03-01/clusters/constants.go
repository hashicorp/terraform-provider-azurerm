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
	ProvisioningStateAccepted           ProvisioningState = "Accepted"
	ProvisioningStateCanceled           ProvisioningState = "Canceled"
	ProvisioningStateConnected          ProvisioningState = "Connected"
	ProvisioningStateCreating           ProvisioningState = "Creating"
	ProvisioningStateDeleted            ProvisioningState = "Deleted"
	ProvisioningStateDeleting           ProvisioningState = "Deleting"
	ProvisioningStateDisableInProgress  ProvisioningState = "DisableInProgress"
	ProvisioningStateDisconnected       ProvisioningState = "Disconnected"
	ProvisioningStateFailed             ProvisioningState = "Failed"
	ProvisioningStateInProgress         ProvisioningState = "InProgress"
	ProvisioningStateMoving             ProvisioningState = "Moving"
	ProvisioningStateNotSpecified       ProvisioningState = "NotSpecified"
	ProvisioningStatePartiallyConnected ProvisioningState = "PartiallyConnected"
	ProvisioningStatePartiallySucceeded ProvisioningState = "PartiallySucceeded"
	ProvisioningStateProvisioning       ProvisioningState = "Provisioning"
	ProvisioningStateSucceeded          ProvisioningState = "Succeeded"
	ProvisioningStateUpdating           ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateAccepted),
		string(ProvisioningStateCanceled),
		string(ProvisioningStateConnected),
		string(ProvisioningStateCreating),
		string(ProvisioningStateDeleted),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateDisableInProgress),
		string(ProvisioningStateDisconnected),
		string(ProvisioningStateFailed),
		string(ProvisioningStateInProgress),
		string(ProvisioningStateMoving),
		string(ProvisioningStateNotSpecified),
		string(ProvisioningStatePartiallyConnected),
		string(ProvisioningStatePartiallySucceeded),
		string(ProvisioningStateProvisioning),
		string(ProvisioningStateSucceeded),
		string(ProvisioningStateUpdating),
	}
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"accepted":           ProvisioningStateAccepted,
		"canceled":           ProvisioningStateCanceled,
		"connected":          ProvisioningStateConnected,
		"creating":           ProvisioningStateCreating,
		"deleted":            ProvisioningStateDeleted,
		"deleting":           ProvisioningStateDeleting,
		"disableinprogress":  ProvisioningStateDisableInProgress,
		"disconnected":       ProvisioningStateDisconnected,
		"failed":             ProvisioningStateFailed,
		"inprogress":         ProvisioningStateInProgress,
		"moving":             ProvisioningStateMoving,
		"notspecified":       ProvisioningStateNotSpecified,
		"partiallyconnected": ProvisioningStatePartiallyConnected,
		"partiallysucceeded": ProvisioningStatePartiallySucceeded,
		"provisioning":       ProvisioningStateProvisioning,
		"succeeded":          ProvisioningStateSucceeded,
		"updating":           ProvisioningStateUpdating,
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
	StatusNotSpecified         Status = "NotSpecified"
	StatusNotYetRegistered     Status = "NotYetRegistered"
)

func PossibleValuesForStatus() []string {
	return []string{
		string(StatusConnectedRecently),
		string(StatusDisconnected),
		string(StatusError),
		string(StatusNotConnectedRecently),
		string(StatusNotSpecified),
		string(StatusNotYetRegistered),
	}
}

func parseStatus(input string) (*Status, error) {
	vals := map[string]Status{
		"connectedrecently":    StatusConnectedRecently,
		"disconnected":         StatusDisconnected,
		"error":                StatusError,
		"notconnectedrecently": StatusNotConnectedRecently,
		"notspecified":         StatusNotSpecified,
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
