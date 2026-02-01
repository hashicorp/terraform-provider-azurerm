package containerappsrevisions

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

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
