package caches

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DomainJoinedType string

const (
	DomainJoinedTypeError DomainJoinedType = "Error"
	DomainJoinedTypeNo    DomainJoinedType = "No"
	DomainJoinedTypeYes   DomainJoinedType = "Yes"
)

func PossibleValuesForDomainJoinedType() []string {
	return []string{
		string(DomainJoinedTypeError),
		string(DomainJoinedTypeNo),
		string(DomainJoinedTypeYes),
	}
}

func (s *DomainJoinedType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDomainJoinedType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDomainJoinedType(input string) (*DomainJoinedType, error) {
	vals := map[string]DomainJoinedType{
		"error": DomainJoinedTypeError,
		"no":    DomainJoinedTypeNo,
		"yes":   DomainJoinedTypeYes,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DomainJoinedType(input)
	return &out, nil
}

type FirmwareStatusType string

const (
	FirmwareStatusTypeAvailable   FirmwareStatusType = "available"
	FirmwareStatusTypeUnavailable FirmwareStatusType = "unavailable"
)

func PossibleValuesForFirmwareStatusType() []string {
	return []string{
		string(FirmwareStatusTypeAvailable),
		string(FirmwareStatusTypeUnavailable),
	}
}

func (s *FirmwareStatusType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseFirmwareStatusType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseFirmwareStatusType(input string) (*FirmwareStatusType, error) {
	vals := map[string]FirmwareStatusType{
		"available":   FirmwareStatusTypeAvailable,
		"unavailable": FirmwareStatusTypeUnavailable,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FirmwareStatusType(input)
	return &out, nil
}

type HealthStateType string

const (
	HealthStateTypeDegraded      HealthStateType = "Degraded"
	HealthStateTypeDown          HealthStateType = "Down"
	HealthStateTypeFlushing      HealthStateType = "Flushing"
	HealthStateTypeHealthy       HealthStateType = "Healthy"
	HealthStateTypeStartFailed   HealthStateType = "StartFailed"
	HealthStateTypeStopped       HealthStateType = "Stopped"
	HealthStateTypeStopping      HealthStateType = "Stopping"
	HealthStateTypeTransitioning HealthStateType = "Transitioning"
	HealthStateTypeUnknown       HealthStateType = "Unknown"
	HealthStateTypeUpgradeFailed HealthStateType = "UpgradeFailed"
	HealthStateTypeUpgrading     HealthStateType = "Upgrading"
	HealthStateTypeWaitingForKey HealthStateType = "WaitingForKey"
)

func PossibleValuesForHealthStateType() []string {
	return []string{
		string(HealthStateTypeDegraded),
		string(HealthStateTypeDown),
		string(HealthStateTypeFlushing),
		string(HealthStateTypeHealthy),
		string(HealthStateTypeStartFailed),
		string(HealthStateTypeStopped),
		string(HealthStateTypeStopping),
		string(HealthStateTypeTransitioning),
		string(HealthStateTypeUnknown),
		string(HealthStateTypeUpgradeFailed),
		string(HealthStateTypeUpgrading),
		string(HealthStateTypeWaitingForKey),
	}
}

func (s *HealthStateType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseHealthStateType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseHealthStateType(input string) (*HealthStateType, error) {
	vals := map[string]HealthStateType{
		"degraded":      HealthStateTypeDegraded,
		"down":          HealthStateTypeDown,
		"flushing":      HealthStateTypeFlushing,
		"healthy":       HealthStateTypeHealthy,
		"startfailed":   HealthStateTypeStartFailed,
		"stopped":       HealthStateTypeStopped,
		"stopping":      HealthStateTypeStopping,
		"transitioning": HealthStateTypeTransitioning,
		"unknown":       HealthStateTypeUnknown,
		"upgradefailed": HealthStateTypeUpgradeFailed,
		"upgrading":     HealthStateTypeUpgrading,
		"waitingforkey": HealthStateTypeWaitingForKey,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HealthStateType(input)
	return &out, nil
}

type NfsAccessRuleAccess string

const (
	NfsAccessRuleAccessNo NfsAccessRuleAccess = "no"
	NfsAccessRuleAccessRo NfsAccessRuleAccess = "ro"
	NfsAccessRuleAccessRw NfsAccessRuleAccess = "rw"
)

func PossibleValuesForNfsAccessRuleAccess() []string {
	return []string{
		string(NfsAccessRuleAccessNo),
		string(NfsAccessRuleAccessRo),
		string(NfsAccessRuleAccessRw),
	}
}

func (s *NfsAccessRuleAccess) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNfsAccessRuleAccess(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNfsAccessRuleAccess(input string) (*NfsAccessRuleAccess, error) {
	vals := map[string]NfsAccessRuleAccess{
		"no": NfsAccessRuleAccessNo,
		"ro": NfsAccessRuleAccessRo,
		"rw": NfsAccessRuleAccessRw,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NfsAccessRuleAccess(input)
	return &out, nil
}

type NfsAccessRuleScope string

const (
	NfsAccessRuleScopeDefault NfsAccessRuleScope = "default"
	NfsAccessRuleScopeHost    NfsAccessRuleScope = "host"
	NfsAccessRuleScopeNetwork NfsAccessRuleScope = "network"
)

func PossibleValuesForNfsAccessRuleScope() []string {
	return []string{
		string(NfsAccessRuleScopeDefault),
		string(NfsAccessRuleScopeHost),
		string(NfsAccessRuleScopeNetwork),
	}
}

func (s *NfsAccessRuleScope) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNfsAccessRuleScope(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNfsAccessRuleScope(input string) (*NfsAccessRuleScope, error) {
	vals := map[string]NfsAccessRuleScope{
		"default": NfsAccessRuleScopeDefault,
		"host":    NfsAccessRuleScopeHost,
		"network": NfsAccessRuleScopeNetwork,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NfsAccessRuleScope(input)
	return &out, nil
}

type PrimingJobState string

const (
	PrimingJobStateComplete PrimingJobState = "Complete"
	PrimingJobStatePaused   PrimingJobState = "Paused"
	PrimingJobStateQueued   PrimingJobState = "Queued"
	PrimingJobStateRunning  PrimingJobState = "Running"
)

func PossibleValuesForPrimingJobState() []string {
	return []string{
		string(PrimingJobStateComplete),
		string(PrimingJobStatePaused),
		string(PrimingJobStateQueued),
		string(PrimingJobStateRunning),
	}
}

func (s *PrimingJobState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePrimingJobState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePrimingJobState(input string) (*PrimingJobState, error) {
	vals := map[string]PrimingJobState{
		"complete": PrimingJobStateComplete,
		"paused":   PrimingJobStatePaused,
		"queued":   PrimingJobStateQueued,
		"running":  PrimingJobStateRunning,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PrimingJobState(input)
	return &out, nil
}

type ProvisioningStateType string

const (
	ProvisioningStateTypeCancelled ProvisioningStateType = "Cancelled"
	ProvisioningStateTypeCreating  ProvisioningStateType = "Creating"
	ProvisioningStateTypeDeleting  ProvisioningStateType = "Deleting"
	ProvisioningStateTypeFailed    ProvisioningStateType = "Failed"
	ProvisioningStateTypeSucceeded ProvisioningStateType = "Succeeded"
	ProvisioningStateTypeUpdating  ProvisioningStateType = "Updating"
)

func PossibleValuesForProvisioningStateType() []string {
	return []string{
		string(ProvisioningStateTypeCancelled),
		string(ProvisioningStateTypeCreating),
		string(ProvisioningStateTypeDeleting),
		string(ProvisioningStateTypeFailed),
		string(ProvisioningStateTypeSucceeded),
		string(ProvisioningStateTypeUpdating),
	}
}

func (s *ProvisioningStateType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProvisioningStateType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProvisioningStateType(input string) (*ProvisioningStateType, error) {
	vals := map[string]ProvisioningStateType{
		"cancelled": ProvisioningStateTypeCancelled,
		"creating":  ProvisioningStateTypeCreating,
		"deleting":  ProvisioningStateTypeDeleting,
		"failed":    ProvisioningStateTypeFailed,
		"succeeded": ProvisioningStateTypeSucceeded,
		"updating":  ProvisioningStateTypeUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningStateType(input)
	return &out, nil
}

type UsernameDownloadedType string

const (
	UsernameDownloadedTypeError UsernameDownloadedType = "Error"
	UsernameDownloadedTypeNo    UsernameDownloadedType = "No"
	UsernameDownloadedTypeYes   UsernameDownloadedType = "Yes"
)

func PossibleValuesForUsernameDownloadedType() []string {
	return []string{
		string(UsernameDownloadedTypeError),
		string(UsernameDownloadedTypeNo),
		string(UsernameDownloadedTypeYes),
	}
}

func (s *UsernameDownloadedType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseUsernameDownloadedType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseUsernameDownloadedType(input string) (*UsernameDownloadedType, error) {
	vals := map[string]UsernameDownloadedType{
		"error": UsernameDownloadedTypeError,
		"no":    UsernameDownloadedTypeNo,
		"yes":   UsernameDownloadedTypeYes,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := UsernameDownloadedType(input)
	return &out, nil
}

type UsernameSource string

const (
	UsernameSourceAD   UsernameSource = "AD"
	UsernameSourceFile UsernameSource = "File"
	UsernameSourceLDAP UsernameSource = "LDAP"
	UsernameSourceNone UsernameSource = "None"
)

func PossibleValuesForUsernameSource() []string {
	return []string{
		string(UsernameSourceAD),
		string(UsernameSourceFile),
		string(UsernameSourceLDAP),
		string(UsernameSourceNone),
	}
}

func (s *UsernameSource) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseUsernameSource(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseUsernameSource(input string) (*UsernameSource, error) {
	vals := map[string]UsernameSource{
		"ad":   UsernameSourceAD,
		"file": UsernameSourceFile,
		"ldap": UsernameSourceLDAP,
		"none": UsernameSourceNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := UsernameSource(input)
	return &out, nil
}
