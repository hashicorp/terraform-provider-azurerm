package virtualmachineinstances

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExtendedLocationTypes string

const (
	ExtendedLocationTypesCustomLocation ExtendedLocationTypes = "CustomLocation"
)

func PossibleValuesForExtendedLocationTypes() []string {
	return []string{
		string(ExtendedLocationTypesCustomLocation),
	}
}

func (s *ExtendedLocationTypes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseExtendedLocationTypes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseExtendedLocationTypes(input string) (*ExtendedLocationTypes, error) {
	vals := map[string]ExtendedLocationTypes{
		"customlocation": ExtendedLocationTypesCustomLocation,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ExtendedLocationTypes(input)
	return &out, nil
}

type OperatingSystemTypes string

const (
	OperatingSystemTypesLinux   OperatingSystemTypes = "Linux"
	OperatingSystemTypesWindows OperatingSystemTypes = "Windows"
)

func PossibleValuesForOperatingSystemTypes() []string {
	return []string{
		string(OperatingSystemTypesLinux),
		string(OperatingSystemTypesWindows),
	}
}

func (s *OperatingSystemTypes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOperatingSystemTypes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOperatingSystemTypes(input string) (*OperatingSystemTypes, error) {
	vals := map[string]OperatingSystemTypes{
		"linux":   OperatingSystemTypesLinux,
		"windows": OperatingSystemTypesWindows,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OperatingSystemTypes(input)
	return &out, nil
}

type PowerStateEnum string

const (
	PowerStateEnumDeallocated  PowerStateEnum = "Deallocated"
	PowerStateEnumDeallocating PowerStateEnum = "Deallocating"
	PowerStateEnumRunning      PowerStateEnum = "Running"
	PowerStateEnumStarting     PowerStateEnum = "Starting"
	PowerStateEnumStopped      PowerStateEnum = "Stopped"
	PowerStateEnumStopping     PowerStateEnum = "Stopping"
	PowerStateEnumUnknown      PowerStateEnum = "Unknown"
)

func PossibleValuesForPowerStateEnum() []string {
	return []string{
		string(PowerStateEnumDeallocated),
		string(PowerStateEnumDeallocating),
		string(PowerStateEnumRunning),
		string(PowerStateEnumStarting),
		string(PowerStateEnumStopped),
		string(PowerStateEnumStopping),
		string(PowerStateEnumUnknown),
	}
}

func (s *PowerStateEnum) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePowerStateEnum(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePowerStateEnum(input string) (*PowerStateEnum, error) {
	vals := map[string]PowerStateEnum{
		"deallocated":  PowerStateEnumDeallocated,
		"deallocating": PowerStateEnumDeallocating,
		"running":      PowerStateEnumRunning,
		"starting":     PowerStateEnumStarting,
		"stopped":      PowerStateEnumStopped,
		"stopping":     PowerStateEnumStopping,
		"unknown":      PowerStateEnumUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PowerStateEnum(input)
	return &out, nil
}

type ProvisioningStateEnum string

const (
	ProvisioningStateEnumAccepted   ProvisioningStateEnum = "Accepted"
	ProvisioningStateEnumCanceled   ProvisioningStateEnum = "Canceled"
	ProvisioningStateEnumDeleting   ProvisioningStateEnum = "Deleting"
	ProvisioningStateEnumFailed     ProvisioningStateEnum = "Failed"
	ProvisioningStateEnumInProgress ProvisioningStateEnum = "InProgress"
	ProvisioningStateEnumSucceeded  ProvisioningStateEnum = "Succeeded"
)

func PossibleValuesForProvisioningStateEnum() []string {
	return []string{
		string(ProvisioningStateEnumAccepted),
		string(ProvisioningStateEnumCanceled),
		string(ProvisioningStateEnumDeleting),
		string(ProvisioningStateEnumFailed),
		string(ProvisioningStateEnumInProgress),
		string(ProvisioningStateEnumSucceeded),
	}
}

func (s *ProvisioningStateEnum) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProvisioningStateEnum(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProvisioningStateEnum(input string) (*ProvisioningStateEnum, error) {
	vals := map[string]ProvisioningStateEnum{
		"accepted":   ProvisioningStateEnumAccepted,
		"canceled":   ProvisioningStateEnumCanceled,
		"deleting":   ProvisioningStateEnumDeleting,
		"failed":     ProvisioningStateEnumFailed,
		"inprogress": ProvisioningStateEnumInProgress,
		"succeeded":  ProvisioningStateEnumSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningStateEnum(input)
	return &out, nil
}

type SecurityTypes string

const (
	SecurityTypesConfidentialVM SecurityTypes = "ConfidentialVM"
	SecurityTypesTrustedLaunch  SecurityTypes = "TrustedLaunch"
)

func PossibleValuesForSecurityTypes() []string {
	return []string{
		string(SecurityTypesConfidentialVM),
		string(SecurityTypesTrustedLaunch),
	}
}

func (s *SecurityTypes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSecurityTypes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSecurityTypes(input string) (*SecurityTypes, error) {
	vals := map[string]SecurityTypes{
		"confidentialvm": SecurityTypesConfidentialVM,
		"trustedlaunch":  SecurityTypesTrustedLaunch,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SecurityTypes(input)
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

type StatusLevelTypes string

const (
	StatusLevelTypesError   StatusLevelTypes = "Error"
	StatusLevelTypesInfo    StatusLevelTypes = "Info"
	StatusLevelTypesWarning StatusLevelTypes = "Warning"
)

func PossibleValuesForStatusLevelTypes() []string {
	return []string{
		string(StatusLevelTypesError),
		string(StatusLevelTypesInfo),
		string(StatusLevelTypesWarning),
	}
}

func (s *StatusLevelTypes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStatusLevelTypes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStatusLevelTypes(input string) (*StatusLevelTypes, error) {
	vals := map[string]StatusLevelTypes{
		"error":   StatusLevelTypesError,
		"info":    StatusLevelTypesInfo,
		"warning": StatusLevelTypesWarning,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StatusLevelTypes(input)
	return &out, nil
}

type StatusTypes string

const (
	StatusTypesFailed     StatusTypes = "Failed"
	StatusTypesInProgress StatusTypes = "InProgress"
	StatusTypesSucceeded  StatusTypes = "Succeeded"
)

func PossibleValuesForStatusTypes() []string {
	return []string{
		string(StatusTypesFailed),
		string(StatusTypesInProgress),
		string(StatusTypesSucceeded),
	}
}

func (s *StatusTypes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStatusTypes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStatusTypes(input string) (*StatusTypes, error) {
	vals := map[string]StatusTypes{
		"failed":     StatusTypesFailed,
		"inprogress": StatusTypesInProgress,
		"succeeded":  StatusTypesSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StatusTypes(input)
	return &out, nil
}

type VMSizeEnum string

const (
	VMSizeEnumCustom                   VMSizeEnum = "Custom"
	VMSizeEnumDefault                  VMSizeEnum = "Default"
	VMSizeEnumStandardAFourVTwo        VMSizeEnum = "Standard_A4_v2"
	VMSizeEnumStandardATwoVTwo         VMSizeEnum = "Standard_A2_v2"
	VMSizeEnumStandardDEightsVThree    VMSizeEnum = "Standard_D8s_v3"
	VMSizeEnumStandardDFoursVThree     VMSizeEnum = "Standard_D4s_v3"
	VMSizeEnumStandardDOneSixsVThree   VMSizeEnum = "Standard_D16s_v3"
	VMSizeEnumStandardDSFiveVTwo       VMSizeEnum = "Standard_DS5_v2"
	VMSizeEnumStandardDSFourVTwo       VMSizeEnum = "Standard_DS4_v2"
	VMSizeEnumStandardDSOneThreeVTwo   VMSizeEnum = "Standard_DS13_v2"
	VMSizeEnumStandardDSThreeVTwo      VMSizeEnum = "Standard_DS3_v2"
	VMSizeEnumStandardDSTwoVTwo        VMSizeEnum = "Standard_DS2_v2"
	VMSizeEnumStandardDThreeTwosVThree VMSizeEnum = "Standard_D32s_v3"
	VMSizeEnumStandardDTwosVThree      VMSizeEnum = "Standard_D2s_v3"
	VMSizeEnumStandardKEightSFiveVOne  VMSizeEnum = "Standard_K8S5_v1"
	VMSizeEnumStandardKEightSFourVOne  VMSizeEnum = "Standard_K8S4_v1"
	VMSizeEnumStandardKEightSThreeVOne VMSizeEnum = "Standard_K8S3_v1"
	VMSizeEnumStandardKEightSTwoVOne   VMSizeEnum = "Standard_K8S2_v1"
	VMSizeEnumStandardKEightSVOne      VMSizeEnum = "Standard_K8S_v1"
	VMSizeEnumStandardNKOneTwo         VMSizeEnum = "Standard_NK12"
	VMSizeEnumStandardNKSix            VMSizeEnum = "Standard_NK6"
	VMSizeEnumStandardNVOneTwo         VMSizeEnum = "Standard_NV12"
	VMSizeEnumStandardNVSix            VMSizeEnum = "Standard_NV6"
)

func PossibleValuesForVMSizeEnum() []string {
	return []string{
		string(VMSizeEnumCustom),
		string(VMSizeEnumDefault),
		string(VMSizeEnumStandardAFourVTwo),
		string(VMSizeEnumStandardATwoVTwo),
		string(VMSizeEnumStandardDEightsVThree),
		string(VMSizeEnumStandardDFoursVThree),
		string(VMSizeEnumStandardDOneSixsVThree),
		string(VMSizeEnumStandardDSFiveVTwo),
		string(VMSizeEnumStandardDSFourVTwo),
		string(VMSizeEnumStandardDSOneThreeVTwo),
		string(VMSizeEnumStandardDSThreeVTwo),
		string(VMSizeEnumStandardDSTwoVTwo),
		string(VMSizeEnumStandardDThreeTwosVThree),
		string(VMSizeEnumStandardDTwosVThree),
		string(VMSizeEnumStandardKEightSFiveVOne),
		string(VMSizeEnumStandardKEightSFourVOne),
		string(VMSizeEnumStandardKEightSThreeVOne),
		string(VMSizeEnumStandardKEightSTwoVOne),
		string(VMSizeEnumStandardKEightSVOne),
		string(VMSizeEnumStandardNKOneTwo),
		string(VMSizeEnumStandardNKSix),
		string(VMSizeEnumStandardNVOneTwo),
		string(VMSizeEnumStandardNVSix),
	}
}

func (s *VMSizeEnum) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVMSizeEnum(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVMSizeEnum(input string) (*VMSizeEnum, error) {
	vals := map[string]VMSizeEnum{
		"custom":           VMSizeEnumCustom,
		"default":          VMSizeEnumDefault,
		"standard_a4_v2":   VMSizeEnumStandardAFourVTwo,
		"standard_a2_v2":   VMSizeEnumStandardATwoVTwo,
		"standard_d8s_v3":  VMSizeEnumStandardDEightsVThree,
		"standard_d4s_v3":  VMSizeEnumStandardDFoursVThree,
		"standard_d16s_v3": VMSizeEnumStandardDOneSixsVThree,
		"standard_ds5_v2":  VMSizeEnumStandardDSFiveVTwo,
		"standard_ds4_v2":  VMSizeEnumStandardDSFourVTwo,
		"standard_ds13_v2": VMSizeEnumStandardDSOneThreeVTwo,
		"standard_ds3_v2":  VMSizeEnumStandardDSThreeVTwo,
		"standard_ds2_v2":  VMSizeEnumStandardDSTwoVTwo,
		"standard_d32s_v3": VMSizeEnumStandardDThreeTwosVThree,
		"standard_d2s_v3":  VMSizeEnumStandardDTwosVThree,
		"standard_k8s5_v1": VMSizeEnumStandardKEightSFiveVOne,
		"standard_k8s4_v1": VMSizeEnumStandardKEightSFourVOne,
		"standard_k8s3_v1": VMSizeEnumStandardKEightSThreeVOne,
		"standard_k8s2_v1": VMSizeEnumStandardKEightSTwoVOne,
		"standard_k8s_v1":  VMSizeEnumStandardKEightSVOne,
		"standard_nk12":    VMSizeEnumStandardNKOneTwo,
		"standard_nk6":     VMSizeEnumStandardNKSix,
		"standard_nv12":    VMSizeEnumStandardNVOneTwo,
		"standard_nv6":     VMSizeEnumStandardNVSix,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VMSizeEnum(input)
	return &out, nil
}
