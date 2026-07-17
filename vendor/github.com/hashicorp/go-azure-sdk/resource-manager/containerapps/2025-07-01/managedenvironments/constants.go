package managedenvironments

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CertificateProvisioningState string

const (
	CertificateProvisioningStateCanceled     CertificateProvisioningState = "Canceled"
	CertificateProvisioningStateDeleteFailed CertificateProvisioningState = "DeleteFailed"
	CertificateProvisioningStateDeleting     CertificateProvisioningState = "Deleting"
	CertificateProvisioningStateFailed       CertificateProvisioningState = "Failed"
	CertificateProvisioningStatePending      CertificateProvisioningState = "Pending"
	CertificateProvisioningStateSucceeded    CertificateProvisioningState = "Succeeded"
)

func PossibleValuesForCertificateProvisioningState() []string {
	return []string{
		string(CertificateProvisioningStateCanceled),
		string(CertificateProvisioningStateDeleteFailed),
		string(CertificateProvisioningStateDeleting),
		string(CertificateProvisioningStateFailed),
		string(CertificateProvisioningStatePending),
		string(CertificateProvisioningStateSucceeded),
	}
}

func (s *CertificateProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCertificateProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCertificateProvisioningState(input string) (*CertificateProvisioningState, error) {
	vals := map[string]CertificateProvisioningState{
		"canceled":     CertificateProvisioningStateCanceled,
		"deletefailed": CertificateProvisioningStateDeleteFailed,
		"deleting":     CertificateProvisioningStateDeleting,
		"failed":       CertificateProvisioningStateFailed,
		"pending":      CertificateProvisioningStatePending,
		"succeeded":    CertificateProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CertificateProvisioningState(input)
	return &out, nil
}

type CheckNameAvailabilityReason string

const (
	CheckNameAvailabilityReasonAlreadyExists CheckNameAvailabilityReason = "AlreadyExists"
	CheckNameAvailabilityReasonInvalid       CheckNameAvailabilityReason = "Invalid"
)

func PossibleValuesForCheckNameAvailabilityReason() []string {
	return []string{
		string(CheckNameAvailabilityReasonAlreadyExists),
		string(CheckNameAvailabilityReasonInvalid),
	}
}

func (s *CheckNameAvailabilityReason) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCheckNameAvailabilityReason(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCheckNameAvailabilityReason(input string) (*CheckNameAvailabilityReason, error) {
	vals := map[string]CheckNameAvailabilityReason{
		"alreadyexists": CheckNameAvailabilityReasonAlreadyExists,
		"invalid":       CheckNameAvailabilityReasonInvalid,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CheckNameAvailabilityReason(input)
	return &out, nil
}

type EnvironmentProvisioningState string

const (
	EnvironmentProvisioningStateCanceled                      EnvironmentProvisioningState = "Canceled"
	EnvironmentProvisioningStateFailed                        EnvironmentProvisioningState = "Failed"
	EnvironmentProvisioningStateInfrastructureSetupComplete   EnvironmentProvisioningState = "InfrastructureSetupComplete"
	EnvironmentProvisioningStateInfrastructureSetupInProgress EnvironmentProvisioningState = "InfrastructureSetupInProgress"
	EnvironmentProvisioningStateInitializationInProgress      EnvironmentProvisioningState = "InitializationInProgress"
	EnvironmentProvisioningStateScheduledForDelete            EnvironmentProvisioningState = "ScheduledForDelete"
	EnvironmentProvisioningStateSucceeded                     EnvironmentProvisioningState = "Succeeded"
	EnvironmentProvisioningStateUpgradeFailed                 EnvironmentProvisioningState = "UpgradeFailed"
	EnvironmentProvisioningStateUpgradeRequested              EnvironmentProvisioningState = "UpgradeRequested"
	EnvironmentProvisioningStateWaiting                       EnvironmentProvisioningState = "Waiting"
)

func PossibleValuesForEnvironmentProvisioningState() []string {
	return []string{
		string(EnvironmentProvisioningStateCanceled),
		string(EnvironmentProvisioningStateFailed),
		string(EnvironmentProvisioningStateInfrastructureSetupComplete),
		string(EnvironmentProvisioningStateInfrastructureSetupInProgress),
		string(EnvironmentProvisioningStateInitializationInProgress),
		string(EnvironmentProvisioningStateScheduledForDelete),
		string(EnvironmentProvisioningStateSucceeded),
		string(EnvironmentProvisioningStateUpgradeFailed),
		string(EnvironmentProvisioningStateUpgradeRequested),
		string(EnvironmentProvisioningStateWaiting),
	}
}

func (s *EnvironmentProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEnvironmentProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEnvironmentProvisioningState(input string) (*EnvironmentProvisioningState, error) {
	vals := map[string]EnvironmentProvisioningState{
		"canceled":                      EnvironmentProvisioningStateCanceled,
		"failed":                        EnvironmentProvisioningStateFailed,
		"infrastructuresetupcomplete":   EnvironmentProvisioningStateInfrastructureSetupComplete,
		"infrastructuresetupinprogress": EnvironmentProvisioningStateInfrastructureSetupInProgress,
		"initializationinprogress":      EnvironmentProvisioningStateInitializationInProgress,
		"scheduledfordelete":            EnvironmentProvisioningStateScheduledForDelete,
		"succeeded":                     EnvironmentProvisioningStateSucceeded,
		"upgradefailed":                 EnvironmentProvisioningStateUpgradeFailed,
		"upgraderequested":              EnvironmentProvisioningStateUpgradeRequested,
		"waiting":                       EnvironmentProvisioningStateWaiting,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EnvironmentProvisioningState(input)
	return &out, nil
}

type ManagedCertificateDomainControlValidation string

const (
	ManagedCertificateDomainControlValidationCNAME ManagedCertificateDomainControlValidation = "CNAME"
	ManagedCertificateDomainControlValidationHTTP  ManagedCertificateDomainControlValidation = "HTTP"
	ManagedCertificateDomainControlValidationTXT   ManagedCertificateDomainControlValidation = "TXT"
)

func PossibleValuesForManagedCertificateDomainControlValidation() []string {
	return []string{
		string(ManagedCertificateDomainControlValidationCNAME),
		string(ManagedCertificateDomainControlValidationHTTP),
		string(ManagedCertificateDomainControlValidationTXT),
	}
}

func (s *ManagedCertificateDomainControlValidation) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseManagedCertificateDomainControlValidation(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseManagedCertificateDomainControlValidation(input string) (*ManagedCertificateDomainControlValidation, error) {
	vals := map[string]ManagedCertificateDomainControlValidation{
		"cname": ManagedCertificateDomainControlValidationCNAME,
		"http":  ManagedCertificateDomainControlValidationHTTP,
		"txt":   ManagedCertificateDomainControlValidationTXT,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ManagedCertificateDomainControlValidation(input)
	return &out, nil
}

type PrivateEndpointConnectionProvisioningState string

const (
	PrivateEndpointConnectionProvisioningStateCanceled  PrivateEndpointConnectionProvisioningState = "Canceled"
	PrivateEndpointConnectionProvisioningStateDeleting  PrivateEndpointConnectionProvisioningState = "Deleting"
	PrivateEndpointConnectionProvisioningStateFailed    PrivateEndpointConnectionProvisioningState = "Failed"
	PrivateEndpointConnectionProvisioningStatePending   PrivateEndpointConnectionProvisioningState = "Pending"
	PrivateEndpointConnectionProvisioningStateSucceeded PrivateEndpointConnectionProvisioningState = "Succeeded"
	PrivateEndpointConnectionProvisioningStateUpdating  PrivateEndpointConnectionProvisioningState = "Updating"
	PrivateEndpointConnectionProvisioningStateWaiting   PrivateEndpointConnectionProvisioningState = "Waiting"
)

func PossibleValuesForPrivateEndpointConnectionProvisioningState() []string {
	return []string{
		string(PrivateEndpointConnectionProvisioningStateCanceled),
		string(PrivateEndpointConnectionProvisioningStateDeleting),
		string(PrivateEndpointConnectionProvisioningStateFailed),
		string(PrivateEndpointConnectionProvisioningStatePending),
		string(PrivateEndpointConnectionProvisioningStateSucceeded),
		string(PrivateEndpointConnectionProvisioningStateUpdating),
		string(PrivateEndpointConnectionProvisioningStateWaiting),
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
		"canceled":  PrivateEndpointConnectionProvisioningStateCanceled,
		"deleting":  PrivateEndpointConnectionProvisioningStateDeleting,
		"failed":    PrivateEndpointConnectionProvisioningStateFailed,
		"pending":   PrivateEndpointConnectionProvisioningStatePending,
		"succeeded": PrivateEndpointConnectionProvisioningStateSucceeded,
		"updating":  PrivateEndpointConnectionProvisioningStateUpdating,
		"waiting":   PrivateEndpointConnectionProvisioningStateWaiting,
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
