package nodepools

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConditionType string

const (
	ConditionTypeAvailable   ConditionType = "Available"
	ConditionTypeDegraded    ConditionType = "Degraded"
	ConditionTypeProgressing ConditionType = "Progressing"
)

func PossibleValuesForConditionType() []string {
	return []string{
		string(ConditionTypeAvailable),
		string(ConditionTypeDegraded),
		string(ConditionTypeProgressing),
	}
}

func (s *ConditionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseConditionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseConditionType(input string) (*ConditionType, error) {
	vals := map[string]ConditionType{
		"available":   ConditionTypeAvailable,
		"degraded":    ConditionTypeDegraded,
		"progressing": ConditionTypeProgressing,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ConditionType(input)
	return &out, nil
}

type DiskStorageAccountType string

const (
	DiskStorageAccountTypePremiumLRS     DiskStorageAccountType = "Premium_LRS"
	DiskStorageAccountTypeStandardLRS    DiskStorageAccountType = "Standard_LRS"
	DiskStorageAccountTypeStandardSSDLRS DiskStorageAccountType = "StandardSSD_LRS"
)

func PossibleValuesForDiskStorageAccountType() []string {
	return []string{
		string(DiskStorageAccountTypePremiumLRS),
		string(DiskStorageAccountTypeStandardLRS),
		string(DiskStorageAccountTypeStandardSSDLRS),
	}
}

func (s *DiskStorageAccountType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDiskStorageAccountType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDiskStorageAccountType(input string) (*DiskStorageAccountType, error) {
	vals := map[string]DiskStorageAccountType{
		"premium_lrs":     DiskStorageAccountTypePremiumLRS,
		"standard_lrs":    DiskStorageAccountTypeStandardLRS,
		"standardssd_lrs": DiskStorageAccountTypeStandardSSDLRS,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DiskStorageAccountType(input)
	return &out, nil
}

type Effect string

const (
	EffectNoExecute        Effect = "NoExecute"
	EffectNoSchedule       Effect = "NoSchedule"
	EffectPreferNoSchedule Effect = "PreferNoSchedule"
)

func PossibleValuesForEffect() []string {
	return []string{
		string(EffectNoExecute),
		string(EffectNoSchedule),
		string(EffectPreferNoSchedule),
	}
}

func (s *Effect) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEffect(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEffect(input string) (*Effect, error) {
	vals := map[string]Effect{
		"noexecute":        EffectNoExecute,
		"noschedule":       EffectNoSchedule,
		"prefernoschedule": EffectPreferNoSchedule,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Effect(input)
	return &out, nil
}

type ManagedServiceIdentityType string

const (
	ManagedServiceIdentityTypeNone                       ManagedServiceIdentityType = "None"
	ManagedServiceIdentityTypeSystemAssigned             ManagedServiceIdentityType = "SystemAssigned"
	ManagedServiceIdentityTypeSystemAssignedUserAssigned ManagedServiceIdentityType = "SystemAssigned,UserAssigned"
	ManagedServiceIdentityTypeUserAssigned               ManagedServiceIdentityType = "UserAssigned"
)

func PossibleValuesForManagedServiceIdentityType() []string {
	return []string{
		string(ManagedServiceIdentityTypeNone),
		string(ManagedServiceIdentityTypeSystemAssigned),
		string(ManagedServiceIdentityTypeSystemAssignedUserAssigned),
		string(ManagedServiceIdentityTypeUserAssigned),
	}
}

func (s *ManagedServiceIdentityType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseManagedServiceIdentityType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseManagedServiceIdentityType(input string) (*ManagedServiceIdentityType, error) {
	vals := map[string]ManagedServiceIdentityType{
		"none":                        ManagedServiceIdentityTypeNone,
		"systemassigned":              ManagedServiceIdentityTypeSystemAssigned,
		"systemassigned,userassigned": ManagedServiceIdentityTypeSystemAssignedUserAssigned,
		"userassigned":                ManagedServiceIdentityTypeUserAssigned,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ManagedServiceIdentityType(input)
	return &out, nil
}

type OsDiskType string

const (
	OsDiskTypeEphemeral OsDiskType = "Ephemeral"
	OsDiskTypeManaged   OsDiskType = "Managed"
)

func PossibleValuesForOsDiskType() []string {
	return []string{
		string(OsDiskTypeEphemeral),
		string(OsDiskTypeManaged),
	}
}

func (s *OsDiskType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOsDiskType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOsDiskType(input string) (*OsDiskType, error) {
	vals := map[string]OsDiskType{
		"ephemeral": OsDiskTypeEphemeral,
		"managed":   OsDiskTypeManaged,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OsDiskType(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateAccepted     ProvisioningState = "Accepted"
	ProvisioningStateCanceled     ProvisioningState = "Canceled"
	ProvisioningStateDeleting     ProvisioningState = "Deleting"
	ProvisioningStateFailed       ProvisioningState = "Failed"
	ProvisioningStateProvisioning ProvisioningState = "Provisioning"
	ProvisioningStateSucceeded    ProvisioningState = "Succeeded"
	ProvisioningStateUpdating     ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateAccepted),
		string(ProvisioningStateCanceled),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateProvisioning),
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
		"accepted":     ProvisioningStateAccepted,
		"canceled":     ProvisioningStateCanceled,
		"deleting":     ProvisioningStateDeleting,
		"failed":       ProvisioningStateFailed,
		"provisioning": ProvisioningStateProvisioning,
		"succeeded":    ProvisioningStateSucceeded,
		"updating":     ProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}

type StatusType string

const (
	StatusTypeFalse   StatusType = "False"
	StatusTypeTrue    StatusType = "True"
	StatusTypeUnknown StatusType = "Unknown"
)

func PossibleValuesForStatusType() []string {
	return []string{
		string(StatusTypeFalse),
		string(StatusTypeTrue),
		string(StatusTypeUnknown),
	}
}

func (s *StatusType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStatusType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStatusType(input string) (*StatusType, error) {
	vals := map[string]StatusType{
		"false":   StatusTypeFalse,
		"true":    StatusTypeTrue,
		"unknown": StatusTypeUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StatusType(input)
	return &out, nil
}
