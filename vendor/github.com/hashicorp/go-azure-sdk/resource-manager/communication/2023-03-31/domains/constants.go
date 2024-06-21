package domains

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DomainManagement string

const (
	DomainManagementAzureManaged                    DomainManagement = "AzureManaged"
	DomainManagementCustomerManaged                 DomainManagement = "CustomerManaged"
	DomainManagementCustomerManagedInExchangeOnline DomainManagement = "CustomerManagedInExchangeOnline"
)

func PossibleValuesForDomainManagement() []string {
	return []string{
		string(DomainManagementAzureManaged),
		string(DomainManagementCustomerManaged),
		string(DomainManagementCustomerManagedInExchangeOnline),
	}
}

func (s *DomainManagement) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDomainManagement(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDomainManagement(input string) (*DomainManagement, error) {
	vals := map[string]DomainManagement{
		"azuremanaged":                    DomainManagementAzureManaged,
		"customermanaged":                 DomainManagementCustomerManaged,
		"customermanagedinexchangeonline": DomainManagementCustomerManagedInExchangeOnline,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DomainManagement(input)
	return &out, nil
}

type DomainsProvisioningState string

const (
	DomainsProvisioningStateCanceled  DomainsProvisioningState = "Canceled"
	DomainsProvisioningStateCreating  DomainsProvisioningState = "Creating"
	DomainsProvisioningStateDeleting  DomainsProvisioningState = "Deleting"
	DomainsProvisioningStateFailed    DomainsProvisioningState = "Failed"
	DomainsProvisioningStateMoving    DomainsProvisioningState = "Moving"
	DomainsProvisioningStateRunning   DomainsProvisioningState = "Running"
	DomainsProvisioningStateSucceeded DomainsProvisioningState = "Succeeded"
	DomainsProvisioningStateUnknown   DomainsProvisioningState = "Unknown"
	DomainsProvisioningStateUpdating  DomainsProvisioningState = "Updating"
)

func PossibleValuesForDomainsProvisioningState() []string {
	return []string{
		string(DomainsProvisioningStateCanceled),
		string(DomainsProvisioningStateCreating),
		string(DomainsProvisioningStateDeleting),
		string(DomainsProvisioningStateFailed),
		string(DomainsProvisioningStateMoving),
		string(DomainsProvisioningStateRunning),
		string(DomainsProvisioningStateSucceeded),
		string(DomainsProvisioningStateUnknown),
		string(DomainsProvisioningStateUpdating),
	}
}

func (s *DomainsProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDomainsProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDomainsProvisioningState(input string) (*DomainsProvisioningState, error) {
	vals := map[string]DomainsProvisioningState{
		"canceled":  DomainsProvisioningStateCanceled,
		"creating":  DomainsProvisioningStateCreating,
		"deleting":  DomainsProvisioningStateDeleting,
		"failed":    DomainsProvisioningStateFailed,
		"moving":    DomainsProvisioningStateMoving,
		"running":   DomainsProvisioningStateRunning,
		"succeeded": DomainsProvisioningStateSucceeded,
		"unknown":   DomainsProvisioningStateUnknown,
		"updating":  DomainsProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DomainsProvisioningState(input)
	return &out, nil
}

type UserEngagementTracking string

const (
	UserEngagementTrackingDisabled UserEngagementTracking = "Disabled"
	UserEngagementTrackingEnabled  UserEngagementTracking = "Enabled"
)

func PossibleValuesForUserEngagementTracking() []string {
	return []string{
		string(UserEngagementTrackingDisabled),
		string(UserEngagementTrackingEnabled),
	}
}

func (s *UserEngagementTracking) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseUserEngagementTracking(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseUserEngagementTracking(input string) (*UserEngagementTracking, error) {
	vals := map[string]UserEngagementTracking{
		"disabled": UserEngagementTrackingDisabled,
		"enabled":  UserEngagementTrackingEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := UserEngagementTracking(input)
	return &out, nil
}

type VerificationStatus string

const (
	VerificationStatusCancellationRequested  VerificationStatus = "CancellationRequested"
	VerificationStatusNotStarted             VerificationStatus = "NotStarted"
	VerificationStatusVerificationFailed     VerificationStatus = "VerificationFailed"
	VerificationStatusVerificationInProgress VerificationStatus = "VerificationInProgress"
	VerificationStatusVerificationRequested  VerificationStatus = "VerificationRequested"
	VerificationStatusVerified               VerificationStatus = "Verified"
)

func PossibleValuesForVerificationStatus() []string {
	return []string{
		string(VerificationStatusCancellationRequested),
		string(VerificationStatusNotStarted),
		string(VerificationStatusVerificationFailed),
		string(VerificationStatusVerificationInProgress),
		string(VerificationStatusVerificationRequested),
		string(VerificationStatusVerified),
	}
}

func (s *VerificationStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVerificationStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVerificationStatus(input string) (*VerificationStatus, error) {
	vals := map[string]VerificationStatus{
		"cancellationrequested":  VerificationStatusCancellationRequested,
		"notstarted":             VerificationStatusNotStarted,
		"verificationfailed":     VerificationStatusVerificationFailed,
		"verificationinprogress": VerificationStatusVerificationInProgress,
		"verificationrequested":  VerificationStatusVerificationRequested,
		"verified":               VerificationStatusVerified,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VerificationStatus(input)
	return &out, nil
}

type VerificationType string

const (
	VerificationTypeDKIM    VerificationType = "DKIM"
	VerificationTypeDKIMTwo VerificationType = "DKIM2"
	VerificationTypeDMARC   VerificationType = "DMARC"
	VerificationTypeDomain  VerificationType = "Domain"
	VerificationTypeSPF     VerificationType = "SPF"
)

func PossibleValuesForVerificationType() []string {
	return []string{
		string(VerificationTypeDKIM),
		string(VerificationTypeDKIMTwo),
		string(VerificationTypeDMARC),
		string(VerificationTypeDomain),
		string(VerificationTypeSPF),
	}
}

func (s *VerificationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVerificationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVerificationType(input string) (*VerificationType, error) {
	vals := map[string]VerificationType{
		"dkim":   VerificationTypeDKIM,
		"dkim2":  VerificationTypeDKIMTwo,
		"dmarc":  VerificationTypeDMARC,
		"domain": VerificationTypeDomain,
		"spf":    VerificationTypeSPF,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VerificationType(input)
	return &out, nil
}
