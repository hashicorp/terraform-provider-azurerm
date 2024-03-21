package updatesummaries

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HealthState string

const (
	HealthStateError      HealthState = "Error"
	HealthStateFailure    HealthState = "Failure"
	HealthStateInProgress HealthState = "InProgress"
	HealthStateSuccess    HealthState = "Success"
	HealthStateUnknown    HealthState = "Unknown"
	HealthStateWarning    HealthState = "Warning"
)

func PossibleValuesForHealthState() []string {
	return []string{
		string(HealthStateError),
		string(HealthStateFailure),
		string(HealthStateInProgress),
		string(HealthStateSuccess),
		string(HealthStateUnknown),
		string(HealthStateWarning),
	}
}

func (s *HealthState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseHealthState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseHealthState(input string) (*HealthState, error) {
	vals := map[string]HealthState{
		"error":      HealthStateError,
		"failure":    HealthStateFailure,
		"inprogress": HealthStateInProgress,
		"success":    HealthStateSuccess,
		"unknown":    HealthStateUnknown,
		"warning":    HealthStateWarning,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HealthState(input)
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

type Severity string

const (
	SeverityCritical      Severity = "Critical"
	SeverityHidden        Severity = "Hidden"
	SeverityInformational Severity = "Informational"
	SeverityWarning       Severity = "Warning"
)

func PossibleValuesForSeverity() []string {
	return []string{
		string(SeverityCritical),
		string(SeverityHidden),
		string(SeverityInformational),
		string(SeverityWarning),
	}
}

func (s *Severity) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSeverity(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSeverity(input string) (*Severity, error) {
	vals := map[string]Severity{
		"critical":      SeverityCritical,
		"hidden":        SeverityHidden,
		"informational": SeverityInformational,
		"warning":       SeverityWarning,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Severity(input)
	return &out, nil
}

type Status string

const (
	StatusFailed     Status = "Failed"
	StatusInProgress Status = "InProgress"
	StatusSucceeded  Status = "Succeeded"
)

func PossibleValuesForStatus() []string {
	return []string{
		string(StatusFailed),
		string(StatusInProgress),
		string(StatusSucceeded),
	}
}

func (s *Status) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStatus(input string) (*Status, error) {
	vals := map[string]Status{
		"failed":     StatusFailed,
		"inprogress": StatusInProgress,
		"succeeded":  StatusSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Status(input)
	return &out, nil
}

type UpdateSummariesPropertiesState string

const (
	UpdateSummariesPropertiesStateAppliedSuccessfully   UpdateSummariesPropertiesState = "AppliedSuccessfully"
	UpdateSummariesPropertiesStateNeedsAttention        UpdateSummariesPropertiesState = "NeedsAttention"
	UpdateSummariesPropertiesStatePreparationFailed     UpdateSummariesPropertiesState = "PreparationFailed"
	UpdateSummariesPropertiesStatePreparationInProgress UpdateSummariesPropertiesState = "PreparationInProgress"
	UpdateSummariesPropertiesStateUnknown               UpdateSummariesPropertiesState = "Unknown"
	UpdateSummariesPropertiesStateUpdateAvailable       UpdateSummariesPropertiesState = "UpdateAvailable"
	UpdateSummariesPropertiesStateUpdateFailed          UpdateSummariesPropertiesState = "UpdateFailed"
	UpdateSummariesPropertiesStateUpdateInProgress      UpdateSummariesPropertiesState = "UpdateInProgress"
)

func PossibleValuesForUpdateSummariesPropertiesState() []string {
	return []string{
		string(UpdateSummariesPropertiesStateAppliedSuccessfully),
		string(UpdateSummariesPropertiesStateNeedsAttention),
		string(UpdateSummariesPropertiesStatePreparationFailed),
		string(UpdateSummariesPropertiesStatePreparationInProgress),
		string(UpdateSummariesPropertiesStateUnknown),
		string(UpdateSummariesPropertiesStateUpdateAvailable),
		string(UpdateSummariesPropertiesStateUpdateFailed),
		string(UpdateSummariesPropertiesStateUpdateInProgress),
	}
}

func (s *UpdateSummariesPropertiesState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseUpdateSummariesPropertiesState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseUpdateSummariesPropertiesState(input string) (*UpdateSummariesPropertiesState, error) {
	vals := map[string]UpdateSummariesPropertiesState{
		"appliedsuccessfully":   UpdateSummariesPropertiesStateAppliedSuccessfully,
		"needsattention":        UpdateSummariesPropertiesStateNeedsAttention,
		"preparationfailed":     UpdateSummariesPropertiesStatePreparationFailed,
		"preparationinprogress": UpdateSummariesPropertiesStatePreparationInProgress,
		"unknown":               UpdateSummariesPropertiesStateUnknown,
		"updateavailable":       UpdateSummariesPropertiesStateUpdateAvailable,
		"updatefailed":          UpdateSummariesPropertiesStateUpdateFailed,
		"updateinprogress":      UpdateSummariesPropertiesStateUpdateInProgress,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := UpdateSummariesPropertiesState(input)
	return &out, nil
}
