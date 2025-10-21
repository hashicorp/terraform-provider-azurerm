package deploymentstacksatresourcegroup

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DenySettingsMode string

const (
	DenySettingsModeDenyDelete         DenySettingsMode = "denyDelete"
	DenySettingsModeDenyWriteAndDelete DenySettingsMode = "denyWriteAndDelete"
	DenySettingsModeNone               DenySettingsMode = "none"
)

func PossibleValuesForDenySettingsMode() []string {
	return []string{
		string(DenySettingsModeDenyDelete),
		string(DenySettingsModeDenyWriteAndDelete),
		string(DenySettingsModeNone),
	}
}

func (s *DenySettingsMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDenySettingsMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDenySettingsMode(input string) (*DenySettingsMode, error) {
	vals := map[string]DenySettingsMode{
		"denydelete":         DenySettingsModeDenyDelete,
		"denywriteanddelete": DenySettingsModeDenyWriteAndDelete,
		"none":               DenySettingsModeNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DenySettingsMode(input)
	return &out, nil
}

type DenyStatusMode string

const (
	DenyStatusModeDenyDelete         DenyStatusMode = "denyDelete"
	DenyStatusModeDenyWriteAndDelete DenyStatusMode = "denyWriteAndDelete"
	DenyStatusModeInapplicable       DenyStatusMode = "inapplicable"
	DenyStatusModeNone               DenyStatusMode = "none"
	DenyStatusModeNotSupported       DenyStatusMode = "notSupported"
	DenyStatusModeRemovedBySystem    DenyStatusMode = "removedBySystem"
)

func PossibleValuesForDenyStatusMode() []string {
	return []string{
		string(DenyStatusModeDenyDelete),
		string(DenyStatusModeDenyWriteAndDelete),
		string(DenyStatusModeInapplicable),
		string(DenyStatusModeNone),
		string(DenyStatusModeNotSupported),
		string(DenyStatusModeRemovedBySystem),
	}
}

func (s *DenyStatusMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDenyStatusMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDenyStatusMode(input string) (*DenyStatusMode, error) {
	vals := map[string]DenyStatusMode{
		"denydelete":         DenyStatusModeDenyDelete,
		"denywriteanddelete": DenyStatusModeDenyWriteAndDelete,
		"inapplicable":       DenyStatusModeInapplicable,
		"none":               DenyStatusModeNone,
		"notsupported":       DenyStatusModeNotSupported,
		"removedbysystem":    DenyStatusModeRemovedBySystem,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DenyStatusMode(input)
	return &out, nil
}

type DeploymentStackProvisioningState string

const (
	DeploymentStackProvisioningStateCanceled                DeploymentStackProvisioningState = "canceled"
	DeploymentStackProvisioningStateCanceling               DeploymentStackProvisioningState = "canceling"
	DeploymentStackProvisioningStateCreating                DeploymentStackProvisioningState = "creating"
	DeploymentStackProvisioningStateDeleting                DeploymentStackProvisioningState = "deleting"
	DeploymentStackProvisioningStateDeletingResources       DeploymentStackProvisioningState = "deletingResources"
	DeploymentStackProvisioningStateDeploying               DeploymentStackProvisioningState = "deploying"
	DeploymentStackProvisioningStateFailed                  DeploymentStackProvisioningState = "failed"
	DeploymentStackProvisioningStateSucceeded               DeploymentStackProvisioningState = "succeeded"
	DeploymentStackProvisioningStateUpdatingDenyAssignments DeploymentStackProvisioningState = "updatingDenyAssignments"
	DeploymentStackProvisioningStateValidating              DeploymentStackProvisioningState = "validating"
	DeploymentStackProvisioningStateWaiting                 DeploymentStackProvisioningState = "waiting"
)

func PossibleValuesForDeploymentStackProvisioningState() []string {
	return []string{
		string(DeploymentStackProvisioningStateCanceled),
		string(DeploymentStackProvisioningStateCanceling),
		string(DeploymentStackProvisioningStateCreating),
		string(DeploymentStackProvisioningStateDeleting),
		string(DeploymentStackProvisioningStateDeletingResources),
		string(DeploymentStackProvisioningStateDeploying),
		string(DeploymentStackProvisioningStateFailed),
		string(DeploymentStackProvisioningStateSucceeded),
		string(DeploymentStackProvisioningStateUpdatingDenyAssignments),
		string(DeploymentStackProvisioningStateValidating),
		string(DeploymentStackProvisioningStateWaiting),
	}
}

func (s *DeploymentStackProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDeploymentStackProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDeploymentStackProvisioningState(input string) (*DeploymentStackProvisioningState, error) {
	vals := map[string]DeploymentStackProvisioningState{
		"canceled":                DeploymentStackProvisioningStateCanceled,
		"canceling":               DeploymentStackProvisioningStateCanceling,
		"creating":                DeploymentStackProvisioningStateCreating,
		"deleting":                DeploymentStackProvisioningStateDeleting,
		"deletingresources":       DeploymentStackProvisioningStateDeletingResources,
		"deploying":               DeploymentStackProvisioningStateDeploying,
		"failed":                  DeploymentStackProvisioningStateFailed,
		"succeeded":               DeploymentStackProvisioningStateSucceeded,
		"updatingdenyassignments": DeploymentStackProvisioningStateUpdatingDenyAssignments,
		"validating":              DeploymentStackProvisioningStateValidating,
		"waiting":                 DeploymentStackProvisioningStateWaiting,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DeploymentStackProvisioningState(input)
	return &out, nil
}

type DeploymentStacksDeleteDetachEnum string

const (
	DeploymentStacksDeleteDetachEnumDelete DeploymentStacksDeleteDetachEnum = "delete"
	DeploymentStacksDeleteDetachEnumDetach DeploymentStacksDeleteDetachEnum = "detach"
)

func PossibleValuesForDeploymentStacksDeleteDetachEnum() []string {
	return []string{
		string(DeploymentStacksDeleteDetachEnumDelete),
		string(DeploymentStacksDeleteDetachEnumDetach),
	}
}

func (s *DeploymentStacksDeleteDetachEnum) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDeploymentStacksDeleteDetachEnum(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDeploymentStacksDeleteDetachEnum(input string) (*DeploymentStacksDeleteDetachEnum, error) {
	vals := map[string]DeploymentStacksDeleteDetachEnum{
		"delete": DeploymentStacksDeleteDetachEnumDelete,
		"detach": DeploymentStacksDeleteDetachEnumDetach,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DeploymentStacksDeleteDetachEnum(input)
	return &out, nil
}

type ResourceStatusMode string

const (
	ResourceStatusModeDeleteFailed     ResourceStatusMode = "deleteFailed"
	ResourceStatusModeManaged          ResourceStatusMode = "managed"
	ResourceStatusModeRemoveDenyFailed ResourceStatusMode = "removeDenyFailed"
)

func PossibleValuesForResourceStatusMode() []string {
	return []string{
		string(ResourceStatusModeDeleteFailed),
		string(ResourceStatusModeManaged),
		string(ResourceStatusModeRemoveDenyFailed),
	}
}

func (s *ResourceStatusMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseResourceStatusMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseResourceStatusMode(input string) (*ResourceStatusMode, error) {
	vals := map[string]ResourceStatusMode{
		"deletefailed":     ResourceStatusModeDeleteFailed,
		"managed":          ResourceStatusModeManaged,
		"removedenyfailed": ResourceStatusModeRemoveDenyFailed,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ResourceStatusMode(input)
	return &out, nil
}
