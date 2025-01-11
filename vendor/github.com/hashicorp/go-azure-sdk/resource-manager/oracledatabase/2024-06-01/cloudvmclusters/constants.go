package cloudvmclusters

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

type CloudVMClusterLifecycleState string

const (
	CloudVMClusterLifecycleStateAvailable             CloudVMClusterLifecycleState = "Available"
	CloudVMClusterLifecycleStateFailed                CloudVMClusterLifecycleState = "Failed"
	CloudVMClusterLifecycleStateMaintenanceInProgress CloudVMClusterLifecycleState = "MaintenanceInProgress"
	CloudVMClusterLifecycleStateProvisioning          CloudVMClusterLifecycleState = "Provisioning"
	CloudVMClusterLifecycleStateTerminated            CloudVMClusterLifecycleState = "Terminated"
	CloudVMClusterLifecycleStateTerminating           CloudVMClusterLifecycleState = "Terminating"
	CloudVMClusterLifecycleStateUpdating              CloudVMClusterLifecycleState = "Updating"
)

func PossibleValuesForCloudVMClusterLifecycleState() []string {
	return []string{
		string(CloudVMClusterLifecycleStateAvailable),
		string(CloudVMClusterLifecycleStateFailed),
		string(CloudVMClusterLifecycleStateMaintenanceInProgress),
		string(CloudVMClusterLifecycleStateProvisioning),
		string(CloudVMClusterLifecycleStateTerminated),
		string(CloudVMClusterLifecycleStateTerminating),
		string(CloudVMClusterLifecycleStateUpdating),
	}
}

func (s *CloudVMClusterLifecycleState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCloudVMClusterLifecycleState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCloudVMClusterLifecycleState(input string) (*CloudVMClusterLifecycleState, error) {
	vals := map[string]CloudVMClusterLifecycleState{
		"available":             CloudVMClusterLifecycleStateAvailable,
		"failed":                CloudVMClusterLifecycleStateFailed,
		"maintenanceinprogress": CloudVMClusterLifecycleStateMaintenanceInProgress,
		"provisioning":          CloudVMClusterLifecycleStateProvisioning,
		"terminated":            CloudVMClusterLifecycleStateTerminated,
		"terminating":           CloudVMClusterLifecycleStateTerminating,
		"updating":              CloudVMClusterLifecycleStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CloudVMClusterLifecycleState(input)
	return &out, nil
}

type DiskRedundancy string

const (
	DiskRedundancyHigh   DiskRedundancy = "High"
	DiskRedundancyNormal DiskRedundancy = "Normal"
)

func PossibleValuesForDiskRedundancy() []string {
	return []string{
		string(DiskRedundancyHigh),
		string(DiskRedundancyNormal),
	}
}

func (s *DiskRedundancy) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDiskRedundancy(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDiskRedundancy(input string) (*DiskRedundancy, error) {
	vals := map[string]DiskRedundancy{
		"high":   DiskRedundancyHigh,
		"normal": DiskRedundancyNormal,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DiskRedundancy(input)
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
