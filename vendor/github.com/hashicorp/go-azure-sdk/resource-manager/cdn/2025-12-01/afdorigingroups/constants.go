package afdorigingroups

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AfdProvisioningState string

const (
	AfdProvisioningStateCreating  AfdProvisioningState = "Creating"
	AfdProvisioningStateDeleting  AfdProvisioningState = "Deleting"
	AfdProvisioningStateFailed    AfdProvisioningState = "Failed"
	AfdProvisioningStateSucceeded AfdProvisioningState = "Succeeded"
	AfdProvisioningStateUpdating  AfdProvisioningState = "Updating"
)

func PossibleValuesForAfdProvisioningState() []string {
	return []string{
		string(AfdProvisioningStateCreating),
		string(AfdProvisioningStateDeleting),
		string(AfdProvisioningStateFailed),
		string(AfdProvisioningStateSucceeded),
		string(AfdProvisioningStateUpdating),
	}
}

func (s *AfdProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAfdProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAfdProvisioningState(input string) (*AfdProvisioningState, error) {
	vals := map[string]AfdProvisioningState{
		"creating":  AfdProvisioningStateCreating,
		"deleting":  AfdProvisioningStateDeleting,
		"failed":    AfdProvisioningStateFailed,
		"succeeded": AfdProvisioningStateSucceeded,
		"updating":  AfdProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AfdProvisioningState(input)
	return &out, nil
}

type DeploymentStatus string

const (
	DeploymentStatusFailed     DeploymentStatus = "Failed"
	DeploymentStatusInProgress DeploymentStatus = "InProgress"
	DeploymentStatusNotStarted DeploymentStatus = "NotStarted"
	DeploymentStatusSucceeded  DeploymentStatus = "Succeeded"
)

func PossibleValuesForDeploymentStatus() []string {
	return []string{
		string(DeploymentStatusFailed),
		string(DeploymentStatusInProgress),
		string(DeploymentStatusNotStarted),
		string(DeploymentStatusSucceeded),
	}
}

func (s *DeploymentStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDeploymentStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDeploymentStatus(input string) (*DeploymentStatus, error) {
	vals := map[string]DeploymentStatus{
		"failed":     DeploymentStatusFailed,
		"inprogress": DeploymentStatusInProgress,
		"notstarted": DeploymentStatusNotStarted,
		"succeeded":  DeploymentStatusSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DeploymentStatus(input)
	return &out, nil
}

type EnabledState string

const (
	EnabledStateDisabled EnabledState = "Disabled"
	EnabledStateEnabled  EnabledState = "Enabled"
)

func PossibleValuesForEnabledState() []string {
	return []string{
		string(EnabledStateDisabled),
		string(EnabledStateEnabled),
	}
}

func (s *EnabledState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEnabledState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEnabledState(input string) (*EnabledState, error) {
	vals := map[string]EnabledState{
		"disabled": EnabledStateDisabled,
		"enabled":  EnabledStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EnabledState(input)
	return &out, nil
}

type HealthProbeRequestType string

const (
	HealthProbeRequestTypeGET    HealthProbeRequestType = "GET"
	HealthProbeRequestTypeHEAD   HealthProbeRequestType = "HEAD"
	HealthProbeRequestTypeNotSet HealthProbeRequestType = "NotSet"
)

func PossibleValuesForHealthProbeRequestType() []string {
	return []string{
		string(HealthProbeRequestTypeGET),
		string(HealthProbeRequestTypeHEAD),
		string(HealthProbeRequestTypeNotSet),
	}
}

func (s *HealthProbeRequestType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseHealthProbeRequestType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseHealthProbeRequestType(input string) (*HealthProbeRequestType, error) {
	vals := map[string]HealthProbeRequestType{
		"get":    HealthProbeRequestTypeGET,
		"head":   HealthProbeRequestTypeHEAD,
		"notset": HealthProbeRequestTypeNotSet,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HealthProbeRequestType(input)
	return &out, nil
}

type OriginAuthenticationType string

const (
	OriginAuthenticationTypeSystemAssignedIdentity OriginAuthenticationType = "SystemAssignedIdentity"
	OriginAuthenticationTypeUserAssignedIdentity   OriginAuthenticationType = "UserAssignedIdentity"
)

func PossibleValuesForOriginAuthenticationType() []string {
	return []string{
		string(OriginAuthenticationTypeSystemAssignedIdentity),
		string(OriginAuthenticationTypeUserAssignedIdentity),
	}
}

func (s *OriginAuthenticationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOriginAuthenticationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOriginAuthenticationType(input string) (*OriginAuthenticationType, error) {
	vals := map[string]OriginAuthenticationType{
		"systemassignedidentity": OriginAuthenticationTypeSystemAssignedIdentity,
		"userassignedidentity":   OriginAuthenticationTypeUserAssignedIdentity,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OriginAuthenticationType(input)
	return &out, nil
}

type ProbeProtocol string

const (
	ProbeProtocolHTTP   ProbeProtocol = "Http"
	ProbeProtocolHTTPS  ProbeProtocol = "Https"
	ProbeProtocolNotSet ProbeProtocol = "NotSet"
)

func PossibleValuesForProbeProtocol() []string {
	return []string{
		string(ProbeProtocolHTTP),
		string(ProbeProtocolHTTPS),
		string(ProbeProtocolNotSet),
	}
}

func (s *ProbeProtocol) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProbeProtocol(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProbeProtocol(input string) (*ProbeProtocol, error) {
	vals := map[string]ProbeProtocol{
		"http":   ProbeProtocolHTTP,
		"https":  ProbeProtocolHTTPS,
		"notset": ProbeProtocolNotSet,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProbeProtocol(input)
	return &out, nil
}

type UsageUnit string

const (
	UsageUnitCount UsageUnit = "Count"
)

func PossibleValuesForUsageUnit() []string {
	return []string{
		string(UsageUnitCount),
	}
}

func (s *UsageUnit) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseUsageUnit(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseUsageUnit(input string) (*UsageUnit, error) {
	vals := map[string]UsageUnit{
		"count": UsageUnitCount,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := UsageUnit(input)
	return &out, nil
}
