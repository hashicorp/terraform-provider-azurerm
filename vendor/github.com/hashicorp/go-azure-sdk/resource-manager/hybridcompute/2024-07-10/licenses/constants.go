package licenses

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LicenseCoreType string

const (
	LicenseCoreTypePCore LicenseCoreType = "pCore"
	LicenseCoreTypeVCore LicenseCoreType = "vCore"
)

func PossibleValuesForLicenseCoreType() []string {
	return []string{
		string(LicenseCoreTypePCore),
		string(LicenseCoreTypeVCore),
	}
}

func (s *LicenseCoreType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLicenseCoreType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLicenseCoreType(input string) (*LicenseCoreType, error) {
	vals := map[string]LicenseCoreType{
		"pcore": LicenseCoreTypePCore,
		"vcore": LicenseCoreTypeVCore,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LicenseCoreType(input)
	return &out, nil
}

type LicenseEdition string

const (
	LicenseEditionDatacenter LicenseEdition = "Datacenter"
	LicenseEditionStandard   LicenseEdition = "Standard"
)

func PossibleValuesForLicenseEdition() []string {
	return []string{
		string(LicenseEditionDatacenter),
		string(LicenseEditionStandard),
	}
}

func (s *LicenseEdition) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLicenseEdition(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLicenseEdition(input string) (*LicenseEdition, error) {
	vals := map[string]LicenseEdition{
		"datacenter": LicenseEditionDatacenter,
		"standard":   LicenseEditionStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LicenseEdition(input)
	return &out, nil
}

type LicenseState string

const (
	LicenseStateActivated   LicenseState = "Activated"
	LicenseStateDeactivated LicenseState = "Deactivated"
)

func PossibleValuesForLicenseState() []string {
	return []string{
		string(LicenseStateActivated),
		string(LicenseStateDeactivated),
	}
}

func (s *LicenseState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLicenseState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLicenseState(input string) (*LicenseState, error) {
	vals := map[string]LicenseState{
		"activated":   LicenseStateActivated,
		"deactivated": LicenseStateDeactivated,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LicenseState(input)
	return &out, nil
}

type LicenseTarget string

const (
	LicenseTargetWindowsServerTwoZeroOneTwo     LicenseTarget = "Windows Server 2012"
	LicenseTargetWindowsServerTwoZeroOneTwoRTwo LicenseTarget = "Windows Server 2012 R2"
)

func PossibleValuesForLicenseTarget() []string {
	return []string{
		string(LicenseTargetWindowsServerTwoZeroOneTwo),
		string(LicenseTargetWindowsServerTwoZeroOneTwoRTwo),
	}
}

func (s *LicenseTarget) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLicenseTarget(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLicenseTarget(input string) (*LicenseTarget, error) {
	vals := map[string]LicenseTarget{
		"windows server 2012":    LicenseTargetWindowsServerTwoZeroOneTwo,
		"windows server 2012 r2": LicenseTargetWindowsServerTwoZeroOneTwoRTwo,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LicenseTarget(input)
	return &out, nil
}

type LicenseType string

const (
	LicenseTypeESU LicenseType = "ESU"
)

func PossibleValuesForLicenseType() []string {
	return []string{
		string(LicenseTypeESU),
	}
}

func (s *LicenseType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLicenseType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLicenseType(input string) (*LicenseType, error) {
	vals := map[string]LicenseType{
		"esu": LicenseTypeESU,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LicenseType(input)
	return &out, nil
}

type ProgramYear string

const (
	ProgramYearYearOne   ProgramYear = "Year 1"
	ProgramYearYearThree ProgramYear = "Year 3"
	ProgramYearYearTwo   ProgramYear = "Year 2"
)

func PossibleValuesForProgramYear() []string {
	return []string{
		string(ProgramYearYearOne),
		string(ProgramYearYearThree),
		string(ProgramYearYearTwo),
	}
}

func (s *ProgramYear) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProgramYear(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProgramYear(input string) (*ProgramYear, error) {
	vals := map[string]ProgramYear{
		"year 1": ProgramYearYearOne,
		"year 3": ProgramYearYearThree,
		"year 2": ProgramYearYearTwo,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProgramYear(input)
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
