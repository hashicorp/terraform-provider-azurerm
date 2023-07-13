package fluxconfiguration

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FluxComplianceState string

const (
	FluxComplianceStateCompliant            FluxComplianceState = "Compliant"
	FluxComplianceStateNonNegativeCompliant FluxComplianceState = "Non-Compliant"
	FluxComplianceStatePending              FluxComplianceState = "Pending"
	FluxComplianceStateSuspended            FluxComplianceState = "Suspended"
	FluxComplianceStateUnknown              FluxComplianceState = "Unknown"
)

func PossibleValuesForFluxComplianceState() []string {
	return []string{
		string(FluxComplianceStateCompliant),
		string(FluxComplianceStateNonNegativeCompliant),
		string(FluxComplianceStatePending),
		string(FluxComplianceStateSuspended),
		string(FluxComplianceStateUnknown),
	}
}

func (s *FluxComplianceState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseFluxComplianceState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseFluxComplianceState(input string) (*FluxComplianceState, error) {
	vals := map[string]FluxComplianceState{
		"compliant":     FluxComplianceStateCompliant,
		"non-compliant": FluxComplianceStateNonNegativeCompliant,
		"pending":       FluxComplianceStatePending,
		"suspended":     FluxComplianceStateSuspended,
		"unknown":       FluxComplianceStateUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FluxComplianceState(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateCanceled  ProvisioningState = "Canceled"
	ProvisioningStateCreating  ProvisioningState = "Creating"
	ProvisioningStateDeleting  ProvisioningState = "Deleting"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
	ProvisioningStateUpdating  ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateCanceled),
		string(ProvisioningStateCreating),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateSucceeded),
		string(ProvisioningStateUpdating),
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
		"canceled":  ProvisioningStateCanceled,
		"creating":  ProvisioningStateCreating,
		"deleting":  ProvisioningStateDeleting,
		"failed":    ProvisioningStateFailed,
		"succeeded": ProvisioningStateSucceeded,
		"updating":  ProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}

type ScopeType string

const (
	ScopeTypeCluster   ScopeType = "cluster"
	ScopeTypeNamespace ScopeType = "namespace"
)

func PossibleValuesForScopeType() []string {
	return []string{
		string(ScopeTypeCluster),
		string(ScopeTypeNamespace),
	}
}

func (s *ScopeType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseScopeType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseScopeType(input string) (*ScopeType, error) {
	vals := map[string]ScopeType{
		"cluster":   ScopeTypeCluster,
		"namespace": ScopeTypeNamespace,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ScopeType(input)
	return &out, nil
}

type SourceKindType string

const (
	SourceKindTypeAzureBlob     SourceKindType = "AzureBlob"
	SourceKindTypeBucket        SourceKindType = "Bucket"
	SourceKindTypeGitRepository SourceKindType = "GitRepository"
)

func PossibleValuesForSourceKindType() []string {
	return []string{
		string(SourceKindTypeAzureBlob),
		string(SourceKindTypeBucket),
		string(SourceKindTypeGitRepository),
	}
}

func (s *SourceKindType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSourceKindType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSourceKindType(input string) (*SourceKindType, error) {
	vals := map[string]SourceKindType{
		"azureblob":     SourceKindTypeAzureBlob,
		"bucket":        SourceKindTypeBucket,
		"gitrepository": SourceKindTypeGitRepository,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SourceKindType(input)
	return &out, nil
}
