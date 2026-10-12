package exadbvmclusters

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureResourceProvisioningState string

const (
	AzureResourceProvisioningStateCanceled     AzureResourceProvisioningState = "Canceled"
	AzureResourceProvisioningStateFailed       AzureResourceProvisioningState = "Failed"
	AzureResourceProvisioningStateProvisioning AzureResourceProvisioningState = "Provisioning"
	AzureResourceProvisioningStateSucceeded    AzureResourceProvisioningState = "Succeeded"
)

func PossibleValuesForAzureResourceProvisioningState() []string {
	return []string{
		string(AzureResourceProvisioningStateCanceled),
		string(AzureResourceProvisioningStateFailed),
		string(AzureResourceProvisioningStateProvisioning),
		string(AzureResourceProvisioningStateSucceeded),
	}
}

func (s *AzureResourceProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAzureResourceProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAzureResourceProvisioningState(input string) (*AzureResourceProvisioningState, error) {
	vals := map[string]AzureResourceProvisioningState{
		"canceled":     AzureResourceProvisioningStateCanceled,
		"failed":       AzureResourceProvisioningStateFailed,
		"provisioning": AzureResourceProvisioningStateProvisioning,
		"succeeded":    AzureResourceProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AzureResourceProvisioningState(input)
	return &out, nil
}

type ExadbVMClusterLifecycleState string

const (
	ExadbVMClusterLifecycleStateAvailable             ExadbVMClusterLifecycleState = "Available"
	ExadbVMClusterLifecycleStateFailed                ExadbVMClusterLifecycleState = "Failed"
	ExadbVMClusterLifecycleStateMaintenanceInProgress ExadbVMClusterLifecycleState = "MaintenanceInProgress"
	ExadbVMClusterLifecycleStateProvisioning          ExadbVMClusterLifecycleState = "Provisioning"
	ExadbVMClusterLifecycleStateTerminated            ExadbVMClusterLifecycleState = "Terminated"
	ExadbVMClusterLifecycleStateTerminating           ExadbVMClusterLifecycleState = "Terminating"
	ExadbVMClusterLifecycleStateUpdating              ExadbVMClusterLifecycleState = "Updating"
)

func PossibleValuesForExadbVMClusterLifecycleState() []string {
	return []string{
		string(ExadbVMClusterLifecycleStateAvailable),
		string(ExadbVMClusterLifecycleStateFailed),
		string(ExadbVMClusterLifecycleStateMaintenanceInProgress),
		string(ExadbVMClusterLifecycleStateProvisioning),
		string(ExadbVMClusterLifecycleStateTerminated),
		string(ExadbVMClusterLifecycleStateTerminating),
		string(ExadbVMClusterLifecycleStateUpdating),
	}
}

func (s *ExadbVMClusterLifecycleState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseExadbVMClusterLifecycleState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseExadbVMClusterLifecycleState(input string) (*ExadbVMClusterLifecycleState, error) {
	vals := map[string]ExadbVMClusterLifecycleState{
		"available":             ExadbVMClusterLifecycleStateAvailable,
		"failed":                ExadbVMClusterLifecycleStateFailed,
		"maintenanceinprogress": ExadbVMClusterLifecycleStateMaintenanceInProgress,
		"provisioning":          ExadbVMClusterLifecycleStateProvisioning,
		"terminated":            ExadbVMClusterLifecycleStateTerminated,
		"terminating":           ExadbVMClusterLifecycleStateTerminating,
		"updating":              ExadbVMClusterLifecycleStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ExadbVMClusterLifecycleState(input)
	return &out, nil
}

type GridImageType string

const (
	GridImageTypeCustomImage   GridImageType = "CustomImage"
	GridImageTypeReleaseUpdate GridImageType = "ReleaseUpdate"
)

func PossibleValuesForGridImageType() []string {
	return []string{
		string(GridImageTypeCustomImage),
		string(GridImageTypeReleaseUpdate),
	}
}

func (s *GridImageType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseGridImageType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseGridImageType(input string) (*GridImageType, error) {
	vals := map[string]GridImageType{
		"customimage":   GridImageTypeCustomImage,
		"releaseupdate": GridImageTypeReleaseUpdate,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := GridImageType(input)
	return &out, nil
}

type IormLifecycleState string

const (
	IormLifecycleStateBootStrapping IormLifecycleState = "BootStrapping"
	IormLifecycleStateDisabled      IormLifecycleState = "Disabled"
	IormLifecycleStateEnabled       IormLifecycleState = "Enabled"
	IormLifecycleStateFailed        IormLifecycleState = "Failed"
	IormLifecycleStateUpdating      IormLifecycleState = "Updating"
)

func PossibleValuesForIormLifecycleState() []string {
	return []string{
		string(IormLifecycleStateBootStrapping),
		string(IormLifecycleStateDisabled),
		string(IormLifecycleStateEnabled),
		string(IormLifecycleStateFailed),
		string(IormLifecycleStateUpdating),
	}
}

func (s *IormLifecycleState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIormLifecycleState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIormLifecycleState(input string) (*IormLifecycleState, error) {
	vals := map[string]IormLifecycleState{
		"bootstrapping": IormLifecycleStateBootStrapping,
		"disabled":      IormLifecycleStateDisabled,
		"enabled":       IormLifecycleStateEnabled,
		"failed":        IormLifecycleStateFailed,
		"updating":      IormLifecycleStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IormLifecycleState(input)
	return &out, nil
}

type LicenseModel string

const (
	LicenseModelBringYourOwnLicense LicenseModel = "BringYourOwnLicense"
	LicenseModelLicenseIncluded     LicenseModel = "LicenseIncluded"
)

func PossibleValuesForLicenseModel() []string {
	return []string{
		string(LicenseModelBringYourOwnLicense),
		string(LicenseModelLicenseIncluded),
	}
}

func (s *LicenseModel) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLicenseModel(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLicenseModel(input string) (*LicenseModel, error) {
	vals := map[string]LicenseModel{
		"bringyourownlicense": LicenseModelBringYourOwnLicense,
		"licenseincluded":     LicenseModelLicenseIncluded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LicenseModel(input)
	return &out, nil
}

type Objective string

const (
	ObjectiveAuto           Objective = "Auto"
	ObjectiveBalanced       Objective = "Balanced"
	ObjectiveBasic          Objective = "Basic"
	ObjectiveHighThroughput Objective = "HighThroughput"
	ObjectiveLowLatency     Objective = "LowLatency"
)

func PossibleValuesForObjective() []string {
	return []string{
		string(ObjectiveAuto),
		string(ObjectiveBalanced),
		string(ObjectiveBasic),
		string(ObjectiveHighThroughput),
		string(ObjectiveLowLatency),
	}
}

func (s *Objective) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseObjective(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseObjective(input string) (*Objective, error) {
	vals := map[string]Objective{
		"auto":           ObjectiveAuto,
		"balanced":       ObjectiveBalanced,
		"basic":          ObjectiveBasic,
		"highthroughput": ObjectiveHighThroughput,
		"lowlatency":     ObjectiveLowLatency,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Objective(input)
	return &out, nil
}

type ShapeAttribute string

const (
	ShapeAttributeBLOCKSTORAGE ShapeAttribute = "BLOCK_STORAGE"
	ShapeAttributeSMARTSTORAGE ShapeAttribute = "SMART_STORAGE"
)

func PossibleValuesForShapeAttribute() []string {
	return []string{
		string(ShapeAttributeBLOCKSTORAGE),
		string(ShapeAttributeSMARTSTORAGE),
	}
}

func (s *ShapeAttribute) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseShapeAttribute(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseShapeAttribute(input string) (*ShapeAttribute, error) {
	vals := map[string]ShapeAttribute{
		"block_storage": ShapeAttributeBLOCKSTORAGE,
		"smart_storage": ShapeAttributeSMARTSTORAGE,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ShapeAttribute(input)
	return &out, nil
}
