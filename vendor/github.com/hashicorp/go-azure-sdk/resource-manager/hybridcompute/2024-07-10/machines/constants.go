package machines

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AgentConfigurationMode string

const (
	AgentConfigurationModeFull    AgentConfigurationMode = "full"
	AgentConfigurationModeMonitor AgentConfigurationMode = "monitor"
)

func PossibleValuesForAgentConfigurationMode() []string {
	return []string{
		string(AgentConfigurationModeFull),
		string(AgentConfigurationModeMonitor),
	}
}

func (s *AgentConfigurationMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAgentConfigurationMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAgentConfigurationMode(input string) (*AgentConfigurationMode, error) {
	vals := map[string]AgentConfigurationMode{
		"full":    AgentConfigurationModeFull,
		"monitor": AgentConfigurationModeMonitor,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AgentConfigurationMode(input)
	return &out, nil
}

type ArcKindEnum string

const (
	ArcKindEnumAVS    ArcKindEnum = "AVS"
	ArcKindEnumAWS    ArcKindEnum = "AWS"
	ArcKindEnumEPS    ArcKindEnum = "EPS"
	ArcKindEnumGCP    ArcKindEnum = "GCP"
	ArcKindEnumHCI    ArcKindEnum = "HCI"
	ArcKindEnumSCVMM  ArcKindEnum = "SCVMM"
	ArcKindEnumVMware ArcKindEnum = "VMware"
)

func PossibleValuesForArcKindEnum() []string {
	return []string{
		string(ArcKindEnumAVS),
		string(ArcKindEnumAWS),
		string(ArcKindEnumEPS),
		string(ArcKindEnumGCP),
		string(ArcKindEnumHCI),
		string(ArcKindEnumSCVMM),
		string(ArcKindEnumVMware),
	}
}

func (s *ArcKindEnum) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseArcKindEnum(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseArcKindEnum(input string) (*ArcKindEnum, error) {
	vals := map[string]ArcKindEnum{
		"avs":    ArcKindEnumAVS,
		"aws":    ArcKindEnumAWS,
		"eps":    ArcKindEnumEPS,
		"gcp":    ArcKindEnumGCP,
		"hci":    ArcKindEnumHCI,
		"scvmm":  ArcKindEnumSCVMM,
		"vmware": ArcKindEnumVMware,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ArcKindEnum(input)
	return &out, nil
}

type AssessmentModeTypes string

const (
	AssessmentModeTypesAutomaticByPlatform AssessmentModeTypes = "AutomaticByPlatform"
	AssessmentModeTypesImageDefault        AssessmentModeTypes = "ImageDefault"
)

func PossibleValuesForAssessmentModeTypes() []string {
	return []string{
		string(AssessmentModeTypesAutomaticByPlatform),
		string(AssessmentModeTypesImageDefault),
	}
}

func (s *AssessmentModeTypes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAssessmentModeTypes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAssessmentModeTypes(input string) (*AssessmentModeTypes, error) {
	vals := map[string]AssessmentModeTypes{
		"automaticbyplatform": AssessmentModeTypesAutomaticByPlatform,
		"imagedefault":        AssessmentModeTypesImageDefault,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AssessmentModeTypes(input)
	return &out, nil
}

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

type HotpatchEnablementStatus string

const (
	HotpatchEnablementStatusActionRequired    HotpatchEnablementStatus = "ActionRequired"
	HotpatchEnablementStatusDisabled          HotpatchEnablementStatus = "Disabled"
	HotpatchEnablementStatusEnabled           HotpatchEnablementStatus = "Enabled"
	HotpatchEnablementStatusPendingEvaluation HotpatchEnablementStatus = "PendingEvaluation"
	HotpatchEnablementStatusUnknown           HotpatchEnablementStatus = "Unknown"
)

func PossibleValuesForHotpatchEnablementStatus() []string {
	return []string{
		string(HotpatchEnablementStatusActionRequired),
		string(HotpatchEnablementStatusDisabled),
		string(HotpatchEnablementStatusEnabled),
		string(HotpatchEnablementStatusPendingEvaluation),
		string(HotpatchEnablementStatusUnknown),
	}
}

func (s *HotpatchEnablementStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseHotpatchEnablementStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseHotpatchEnablementStatus(input string) (*HotpatchEnablementStatus, error) {
	vals := map[string]HotpatchEnablementStatus{
		"actionrequired":    HotpatchEnablementStatusActionRequired,
		"disabled":          HotpatchEnablementStatusDisabled,
		"enabled":           HotpatchEnablementStatusEnabled,
		"pendingevaluation": HotpatchEnablementStatusPendingEvaluation,
		"unknown":           HotpatchEnablementStatusUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HotpatchEnablementStatus(input)
	return &out, nil
}

type InstanceViewTypes string

const (
	InstanceViewTypesInstanceView InstanceViewTypes = "instanceView"
)

func PossibleValuesForInstanceViewTypes() []string {
	return []string{
		string(InstanceViewTypesInstanceView),
	}
}

func (s *InstanceViewTypes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseInstanceViewTypes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseInstanceViewTypes(input string) (*InstanceViewTypes, error) {
	vals := map[string]InstanceViewTypes{
		"instanceview": InstanceViewTypesInstanceView,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := InstanceViewTypes(input)
	return &out, nil
}

type LastAttemptStatusEnum string

const (
	LastAttemptStatusEnumFailed  LastAttemptStatusEnum = "Failed"
	LastAttemptStatusEnumSuccess LastAttemptStatusEnum = "Success"
)

func PossibleValuesForLastAttemptStatusEnum() []string {
	return []string{
		string(LastAttemptStatusEnumFailed),
		string(LastAttemptStatusEnumSuccess),
	}
}

func (s *LastAttemptStatusEnum) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLastAttemptStatusEnum(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLastAttemptStatusEnum(input string) (*LastAttemptStatusEnum, error) {
	vals := map[string]LastAttemptStatusEnum{
		"failed":  LastAttemptStatusEnumFailed,
		"success": LastAttemptStatusEnumSuccess,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LastAttemptStatusEnum(input)
	return &out, nil
}

type LicenseAssignmentState string

const (
	LicenseAssignmentStateAssigned    LicenseAssignmentState = "Assigned"
	LicenseAssignmentStateNotAssigned LicenseAssignmentState = "NotAssigned"
)

func PossibleValuesForLicenseAssignmentState() []string {
	return []string{
		string(LicenseAssignmentStateAssigned),
		string(LicenseAssignmentStateNotAssigned),
	}
}

func (s *LicenseAssignmentState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLicenseAssignmentState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLicenseAssignmentState(input string) (*LicenseAssignmentState, error) {
	vals := map[string]LicenseAssignmentState{
		"assigned":    LicenseAssignmentStateAssigned,
		"notassigned": LicenseAssignmentStateNotAssigned,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LicenseAssignmentState(input)
	return &out, nil
}

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

type LicenseStatus string

const (
	LicenseStatusExtendedGrace   LicenseStatus = "ExtendedGrace"
	LicenseStatusLicensed        LicenseStatus = "Licensed"
	LicenseStatusNonGenuineGrace LicenseStatus = "NonGenuineGrace"
	LicenseStatusNotification    LicenseStatus = "Notification"
	LicenseStatusOOBGrace        LicenseStatus = "OOBGrace"
	LicenseStatusOOTGrace        LicenseStatus = "OOTGrace"
	LicenseStatusUnlicensed      LicenseStatus = "Unlicensed"
)

func PossibleValuesForLicenseStatus() []string {
	return []string{
		string(LicenseStatusExtendedGrace),
		string(LicenseStatusLicensed),
		string(LicenseStatusNonGenuineGrace),
		string(LicenseStatusNotification),
		string(LicenseStatusOOBGrace),
		string(LicenseStatusOOTGrace),
		string(LicenseStatusUnlicensed),
	}
}

func (s *LicenseStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLicenseStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLicenseStatus(input string) (*LicenseStatus, error) {
	vals := map[string]LicenseStatus{
		"extendedgrace":   LicenseStatusExtendedGrace,
		"licensed":        LicenseStatusLicensed,
		"nongenuinegrace": LicenseStatusNonGenuineGrace,
		"notification":    LicenseStatusNotification,
		"oobgrace":        LicenseStatusOOBGrace,
		"ootgrace":        LicenseStatusOOTGrace,
		"unlicensed":      LicenseStatusUnlicensed,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LicenseStatus(input)
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

type OsType string

const (
	OsTypeLinux   OsType = "Linux"
	OsTypeWindows OsType = "Windows"
)

func PossibleValuesForOsType() []string {
	return []string{
		string(OsTypeLinux),
		string(OsTypeWindows),
	}
}

func (s *OsType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOsType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOsType(input string) (*OsType, error) {
	vals := map[string]OsType{
		"linux":   OsTypeLinux,
		"windows": OsTypeWindows,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OsType(input)
	return &out, nil
}

type PatchModeTypes string

const (
	PatchModeTypesAutomaticByOS       PatchModeTypes = "AutomaticByOS"
	PatchModeTypesAutomaticByPlatform PatchModeTypes = "AutomaticByPlatform"
	PatchModeTypesImageDefault        PatchModeTypes = "ImageDefault"
	PatchModeTypesManual              PatchModeTypes = "Manual"
)

func PossibleValuesForPatchModeTypes() []string {
	return []string{
		string(PatchModeTypesAutomaticByOS),
		string(PatchModeTypesAutomaticByPlatform),
		string(PatchModeTypesImageDefault),
		string(PatchModeTypesManual),
	}
}

func (s *PatchModeTypes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePatchModeTypes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePatchModeTypes(input string) (*PatchModeTypes, error) {
	vals := map[string]PatchModeTypes{
		"automaticbyos":       PatchModeTypesAutomaticByOS,
		"automaticbyplatform": PatchModeTypesAutomaticByPlatform,
		"imagedefault":        PatchModeTypesImageDefault,
		"manual":              PatchModeTypesManual,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PatchModeTypes(input)
	return &out, nil
}

type PatchOperationStartedBy string

const (
	PatchOperationStartedByPlatform PatchOperationStartedBy = "Platform"
	PatchOperationStartedByUser     PatchOperationStartedBy = "User"
)

func PossibleValuesForPatchOperationStartedBy() []string {
	return []string{
		string(PatchOperationStartedByPlatform),
		string(PatchOperationStartedByUser),
	}
}

func (s *PatchOperationStartedBy) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePatchOperationStartedBy(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePatchOperationStartedBy(input string) (*PatchOperationStartedBy, error) {
	vals := map[string]PatchOperationStartedBy{
		"platform": PatchOperationStartedByPlatform,
		"user":     PatchOperationStartedByUser,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PatchOperationStartedBy(input)
	return &out, nil
}

type PatchOperationStatus string

const (
	PatchOperationStatusCompletedWithWarnings PatchOperationStatus = "CompletedWithWarnings"
	PatchOperationStatusFailed                PatchOperationStatus = "Failed"
	PatchOperationStatusInProgress            PatchOperationStatus = "InProgress"
	PatchOperationStatusSucceeded             PatchOperationStatus = "Succeeded"
	PatchOperationStatusUnknown               PatchOperationStatus = "Unknown"
)

func PossibleValuesForPatchOperationStatus() []string {
	return []string{
		string(PatchOperationStatusCompletedWithWarnings),
		string(PatchOperationStatusFailed),
		string(PatchOperationStatusInProgress),
		string(PatchOperationStatusSucceeded),
		string(PatchOperationStatusUnknown),
	}
}

func (s *PatchOperationStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePatchOperationStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePatchOperationStatus(input string) (*PatchOperationStatus, error) {
	vals := map[string]PatchOperationStatus{
		"completedwithwarnings": PatchOperationStatusCompletedWithWarnings,
		"failed":                PatchOperationStatusFailed,
		"inprogress":            PatchOperationStatusInProgress,
		"succeeded":             PatchOperationStatusSucceeded,
		"unknown":               PatchOperationStatusUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PatchOperationStatus(input)
	return &out, nil
}

type PatchServiceUsed string

const (
	PatchServiceUsedAPT     PatchServiceUsed = "APT"
	PatchServiceUsedUnknown PatchServiceUsed = "Unknown"
	PatchServiceUsedWU      PatchServiceUsed = "WU"
	PatchServiceUsedWUWSUS  PatchServiceUsed = "WU_WSUS"
	PatchServiceUsedYUM     PatchServiceUsed = "YUM"
	PatchServiceUsedZypper  PatchServiceUsed = "Zypper"
)

func PossibleValuesForPatchServiceUsed() []string {
	return []string{
		string(PatchServiceUsedAPT),
		string(PatchServiceUsedUnknown),
		string(PatchServiceUsedWU),
		string(PatchServiceUsedWUWSUS),
		string(PatchServiceUsedYUM),
		string(PatchServiceUsedZypper),
	}
}

func (s *PatchServiceUsed) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePatchServiceUsed(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePatchServiceUsed(input string) (*PatchServiceUsed, error) {
	vals := map[string]PatchServiceUsed{
		"apt":     PatchServiceUsedAPT,
		"unknown": PatchServiceUsedUnknown,
		"wu":      PatchServiceUsedWU,
		"wu_wsus": PatchServiceUsedWUWSUS,
		"yum":     PatchServiceUsedYUM,
		"zypper":  PatchServiceUsedZypper,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PatchServiceUsed(input)
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
	StatusTypesConnected    StatusTypes = "Connected"
	StatusTypesDisconnected StatusTypes = "Disconnected"
	StatusTypesError        StatusTypes = "Error"
)

func PossibleValuesForStatusTypes() []string {
	return []string{
		string(StatusTypesConnected),
		string(StatusTypesDisconnected),
		string(StatusTypesError),
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
		"connected":    StatusTypesConnected,
		"disconnected": StatusTypesDisconnected,
		"error":        StatusTypesError,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StatusTypes(input)
	return &out, nil
}

type VMGuestPatchClassificationLinux string

const (
	VMGuestPatchClassificationLinuxCritical VMGuestPatchClassificationLinux = "Critical"
	VMGuestPatchClassificationLinuxOther    VMGuestPatchClassificationLinux = "Other"
	VMGuestPatchClassificationLinuxSecurity VMGuestPatchClassificationLinux = "Security"
)

func PossibleValuesForVMGuestPatchClassificationLinux() []string {
	return []string{
		string(VMGuestPatchClassificationLinuxCritical),
		string(VMGuestPatchClassificationLinuxOther),
		string(VMGuestPatchClassificationLinuxSecurity),
	}
}

func (s *VMGuestPatchClassificationLinux) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVMGuestPatchClassificationLinux(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVMGuestPatchClassificationLinux(input string) (*VMGuestPatchClassificationLinux, error) {
	vals := map[string]VMGuestPatchClassificationLinux{
		"critical": VMGuestPatchClassificationLinuxCritical,
		"other":    VMGuestPatchClassificationLinuxOther,
		"security": VMGuestPatchClassificationLinuxSecurity,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VMGuestPatchClassificationLinux(input)
	return &out, nil
}

type VMGuestPatchClassificationWindows string

const (
	VMGuestPatchClassificationWindowsCritical     VMGuestPatchClassificationWindows = "Critical"
	VMGuestPatchClassificationWindowsDefinition   VMGuestPatchClassificationWindows = "Definition"
	VMGuestPatchClassificationWindowsFeaturePack  VMGuestPatchClassificationWindows = "FeaturePack"
	VMGuestPatchClassificationWindowsSecurity     VMGuestPatchClassificationWindows = "Security"
	VMGuestPatchClassificationWindowsServicePack  VMGuestPatchClassificationWindows = "ServicePack"
	VMGuestPatchClassificationWindowsTools        VMGuestPatchClassificationWindows = "Tools"
	VMGuestPatchClassificationWindowsUpdateRollUp VMGuestPatchClassificationWindows = "UpdateRollUp"
	VMGuestPatchClassificationWindowsUpdates      VMGuestPatchClassificationWindows = "Updates"
)

func PossibleValuesForVMGuestPatchClassificationWindows() []string {
	return []string{
		string(VMGuestPatchClassificationWindowsCritical),
		string(VMGuestPatchClassificationWindowsDefinition),
		string(VMGuestPatchClassificationWindowsFeaturePack),
		string(VMGuestPatchClassificationWindowsSecurity),
		string(VMGuestPatchClassificationWindowsServicePack),
		string(VMGuestPatchClassificationWindowsTools),
		string(VMGuestPatchClassificationWindowsUpdateRollUp),
		string(VMGuestPatchClassificationWindowsUpdates),
	}
}

func (s *VMGuestPatchClassificationWindows) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVMGuestPatchClassificationWindows(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVMGuestPatchClassificationWindows(input string) (*VMGuestPatchClassificationWindows, error) {
	vals := map[string]VMGuestPatchClassificationWindows{
		"critical":     VMGuestPatchClassificationWindowsCritical,
		"definition":   VMGuestPatchClassificationWindowsDefinition,
		"featurepack":  VMGuestPatchClassificationWindowsFeaturePack,
		"security":     VMGuestPatchClassificationWindowsSecurity,
		"servicepack":  VMGuestPatchClassificationWindowsServicePack,
		"tools":        VMGuestPatchClassificationWindowsTools,
		"updaterollup": VMGuestPatchClassificationWindowsUpdateRollUp,
		"updates":      VMGuestPatchClassificationWindowsUpdates,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VMGuestPatchClassificationWindows(input)
	return &out, nil
}

type VMGuestPatchRebootSetting string

const (
	VMGuestPatchRebootSettingAlways     VMGuestPatchRebootSetting = "Always"
	VMGuestPatchRebootSettingIfRequired VMGuestPatchRebootSetting = "IfRequired"
	VMGuestPatchRebootSettingNever      VMGuestPatchRebootSetting = "Never"
)

func PossibleValuesForVMGuestPatchRebootSetting() []string {
	return []string{
		string(VMGuestPatchRebootSettingAlways),
		string(VMGuestPatchRebootSettingIfRequired),
		string(VMGuestPatchRebootSettingNever),
	}
}

func (s *VMGuestPatchRebootSetting) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVMGuestPatchRebootSetting(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVMGuestPatchRebootSetting(input string) (*VMGuestPatchRebootSetting, error) {
	vals := map[string]VMGuestPatchRebootSetting{
		"always":     VMGuestPatchRebootSettingAlways,
		"ifrequired": VMGuestPatchRebootSettingIfRequired,
		"never":      VMGuestPatchRebootSettingNever,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VMGuestPatchRebootSetting(input)
	return &out, nil
}

type VMGuestPatchRebootStatus string

const (
	VMGuestPatchRebootStatusCompleted VMGuestPatchRebootStatus = "Completed"
	VMGuestPatchRebootStatusFailed    VMGuestPatchRebootStatus = "Failed"
	VMGuestPatchRebootStatusNotNeeded VMGuestPatchRebootStatus = "NotNeeded"
	VMGuestPatchRebootStatusRequired  VMGuestPatchRebootStatus = "Required"
	VMGuestPatchRebootStatusStarted   VMGuestPatchRebootStatus = "Started"
	VMGuestPatchRebootStatusUnknown   VMGuestPatchRebootStatus = "Unknown"
)

func PossibleValuesForVMGuestPatchRebootStatus() []string {
	return []string{
		string(VMGuestPatchRebootStatusCompleted),
		string(VMGuestPatchRebootStatusFailed),
		string(VMGuestPatchRebootStatusNotNeeded),
		string(VMGuestPatchRebootStatusRequired),
		string(VMGuestPatchRebootStatusStarted),
		string(VMGuestPatchRebootStatusUnknown),
	}
}

func (s *VMGuestPatchRebootStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVMGuestPatchRebootStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVMGuestPatchRebootStatus(input string) (*VMGuestPatchRebootStatus, error) {
	vals := map[string]VMGuestPatchRebootStatus{
		"completed": VMGuestPatchRebootStatusCompleted,
		"failed":    VMGuestPatchRebootStatusFailed,
		"notneeded": VMGuestPatchRebootStatusNotNeeded,
		"required":  VMGuestPatchRebootStatusRequired,
		"started":   VMGuestPatchRebootStatusStarted,
		"unknown":   VMGuestPatchRebootStatusUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VMGuestPatchRebootStatus(input)
	return &out, nil
}
