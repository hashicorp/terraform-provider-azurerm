package licenseprofiles

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EsuEligibility string

const (
	EsuEligibilityEligible   EsuEligibility = "Eligible"
	EsuEligibilityIneligible EsuEligibility = "Ineligible"
	EsuEligibilityUnknown    EsuEligibility = "Unknown"
)

func PossibleValuesForEsuEligibility() []string {
	return []string{
		string(EsuEligibilityEligible),
		string(EsuEligibilityIneligible),
		string(EsuEligibilityUnknown),
	}
}

func (s *EsuEligibility) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEsuEligibility(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEsuEligibility(input string) (*EsuEligibility, error) {
	vals := map[string]EsuEligibility{
		"eligible":   EsuEligibilityEligible,
		"ineligible": EsuEligibilityIneligible,
		"unknown":    EsuEligibilityUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EsuEligibility(input)
	return &out, nil
}

type EsuKeyState string

const (
	EsuKeyStateActive   EsuKeyState = "Active"
	EsuKeyStateInactive EsuKeyState = "Inactive"
)

func PossibleValuesForEsuKeyState() []string {
	return []string{
		string(EsuKeyStateActive),
		string(EsuKeyStateInactive),
	}
}

func (s *EsuKeyState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEsuKeyState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEsuKeyState(input string) (*EsuKeyState, error) {
	vals := map[string]EsuKeyState{
		"active":   EsuKeyStateActive,
		"inactive": EsuKeyStateInactive,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EsuKeyState(input)
	return &out, nil
}

type EsuServerType string

const (
	EsuServerTypeDatacenter EsuServerType = "Datacenter"
	EsuServerTypeStandard   EsuServerType = "Standard"
)

func PossibleValuesForEsuServerType() []string {
	return []string{
		string(EsuServerTypeDatacenter),
		string(EsuServerTypeStandard),
	}
}

func (s *EsuServerType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEsuServerType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEsuServerType(input string) (*EsuServerType, error) {
	vals := map[string]EsuServerType{
		"datacenter": EsuServerTypeDatacenter,
		"standard":   EsuServerTypeStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EsuServerType(input)
	return &out, nil
}

type LicenseProfileProductType string

const (
	LicenseProfileProductTypeWindowsIoTEnterprise LicenseProfileProductType = "WindowsIoTEnterprise"
	LicenseProfileProductTypeWindowsServer        LicenseProfileProductType = "WindowsServer"
)

func PossibleValuesForLicenseProfileProductType() []string {
	return []string{
		string(LicenseProfileProductTypeWindowsIoTEnterprise),
		string(LicenseProfileProductTypeWindowsServer),
	}
}

func (s *LicenseProfileProductType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLicenseProfileProductType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLicenseProfileProductType(input string) (*LicenseProfileProductType, error) {
	vals := map[string]LicenseProfileProductType{
		"windowsiotenterprise": LicenseProfileProductTypeWindowsIoTEnterprise,
		"windowsserver":        LicenseProfileProductTypeWindowsServer,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LicenseProfileProductType(input)
	return &out, nil
}

type LicenseProfileSubscriptionStatus string

const (
	LicenseProfileSubscriptionStatusDisabled  LicenseProfileSubscriptionStatus = "Disabled"
	LicenseProfileSubscriptionStatusDisabling LicenseProfileSubscriptionStatus = "Disabling"
	LicenseProfileSubscriptionStatusEnabled   LicenseProfileSubscriptionStatus = "Enabled"
	LicenseProfileSubscriptionStatusEnabling  LicenseProfileSubscriptionStatus = "Enabling"
	LicenseProfileSubscriptionStatusFailed    LicenseProfileSubscriptionStatus = "Failed"
	LicenseProfileSubscriptionStatusUnknown   LicenseProfileSubscriptionStatus = "Unknown"
)

func PossibleValuesForLicenseProfileSubscriptionStatus() []string {
	return []string{
		string(LicenseProfileSubscriptionStatusDisabled),
		string(LicenseProfileSubscriptionStatusDisabling),
		string(LicenseProfileSubscriptionStatusEnabled),
		string(LicenseProfileSubscriptionStatusEnabling),
		string(LicenseProfileSubscriptionStatusFailed),
		string(LicenseProfileSubscriptionStatusUnknown),
	}
}

func (s *LicenseProfileSubscriptionStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLicenseProfileSubscriptionStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLicenseProfileSubscriptionStatus(input string) (*LicenseProfileSubscriptionStatus, error) {
	vals := map[string]LicenseProfileSubscriptionStatus{
		"disabled":  LicenseProfileSubscriptionStatusDisabled,
		"disabling": LicenseProfileSubscriptionStatusDisabling,
		"enabled":   LicenseProfileSubscriptionStatusEnabled,
		"enabling":  LicenseProfileSubscriptionStatusEnabling,
		"failed":    LicenseProfileSubscriptionStatusFailed,
		"unknown":   LicenseProfileSubscriptionStatusUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LicenseProfileSubscriptionStatus(input)
	return &out, nil
}

type LicenseProfileSubscriptionStatusUpdate string

const (
	LicenseProfileSubscriptionStatusUpdateDisable LicenseProfileSubscriptionStatusUpdate = "Disable"
	LicenseProfileSubscriptionStatusUpdateEnable  LicenseProfileSubscriptionStatusUpdate = "Enable"
)

func PossibleValuesForLicenseProfileSubscriptionStatusUpdate() []string {
	return []string{
		string(LicenseProfileSubscriptionStatusUpdateDisable),
		string(LicenseProfileSubscriptionStatusUpdateEnable),
	}
}

func (s *LicenseProfileSubscriptionStatusUpdate) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLicenseProfileSubscriptionStatusUpdate(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLicenseProfileSubscriptionStatusUpdate(input string) (*LicenseProfileSubscriptionStatusUpdate, error) {
	vals := map[string]LicenseProfileSubscriptionStatusUpdate{
		"disable": LicenseProfileSubscriptionStatusUpdateDisable,
		"enable":  LicenseProfileSubscriptionStatusUpdateEnable,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LicenseProfileSubscriptionStatusUpdate(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateAccepted  ProvisioningState = "Accepted"
	ProvisioningStateCanceled  ProvisioningState = "Canceled"
	ProvisioningStateCreating  ProvisioningState = "Creating"
	ProvisioningStateDeleted   ProvisioningState = "Deleted"
	ProvisioningStateDeleting  ProvisioningState = "Deleting"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
	ProvisioningStateUpdating  ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateAccepted),
		string(ProvisioningStateCanceled),
		string(ProvisioningStateCreating),
		string(ProvisioningStateDeleted),
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
		"accepted":  ProvisioningStateAccepted,
		"canceled":  ProvisioningStateCanceled,
		"creating":  ProvisioningStateCreating,
		"deleted":   ProvisioningStateDeleted,
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
