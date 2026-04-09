package volumes

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OperationalStatus string

const (
	OperationalStatusHealthy            OperationalStatus = "Healthy"
	OperationalStatusInvalid            OperationalStatus = "Invalid"
	OperationalStatusRunning            OperationalStatus = "Running"
	OperationalStatusStopped            OperationalStatus = "Stopped"
	OperationalStatusStoppedDeallocated OperationalStatus = "Stopped (deallocated)"
	OperationalStatusUnhealthy          OperationalStatus = "Unhealthy"
	OperationalStatusUnknown            OperationalStatus = "Unknown"
	OperationalStatusUpdating           OperationalStatus = "Updating"
)

func PossibleValuesForOperationalStatus() []string {
	return []string{
		string(OperationalStatusHealthy),
		string(OperationalStatusInvalid),
		string(OperationalStatusRunning),
		string(OperationalStatusStopped),
		string(OperationalStatusStoppedDeallocated),
		string(OperationalStatusUnhealthy),
		string(OperationalStatusUnknown),
		string(OperationalStatusUpdating),
	}
}

func (s *OperationalStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOperationalStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOperationalStatus(input string) (*OperationalStatus, error) {
	vals := map[string]OperationalStatus{
		"healthy":               OperationalStatusHealthy,
		"invalid":               OperationalStatusInvalid,
		"running":               OperationalStatusRunning,
		"stopped":               OperationalStatusStopped,
		"stopped (deallocated)": OperationalStatusStoppedDeallocated,
		"unhealthy":             OperationalStatusUnhealthy,
		"unknown":               OperationalStatusUnknown,
		"updating":              OperationalStatusUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OperationalStatus(input)
	return &out, nil
}

type ProvisioningStates string

const (
	ProvisioningStatesCanceled  ProvisioningStates = "Canceled"
	ProvisioningStatesCreating  ProvisioningStates = "Creating"
	ProvisioningStatesDeleting  ProvisioningStates = "Deleting"
	ProvisioningStatesFailed    ProvisioningStates = "Failed"
	ProvisioningStatesInvalid   ProvisioningStates = "Invalid"
	ProvisioningStatesPending   ProvisioningStates = "Pending"
	ProvisioningStatesSucceeded ProvisioningStates = "Succeeded"
	ProvisioningStatesUpdating  ProvisioningStates = "Updating"
)

func PossibleValuesForProvisioningStates() []string {
	return []string{
		string(ProvisioningStatesCanceled),
		string(ProvisioningStatesCreating),
		string(ProvisioningStatesDeleting),
		string(ProvisioningStatesFailed),
		string(ProvisioningStatesInvalid),
		string(ProvisioningStatesPending),
		string(ProvisioningStatesSucceeded),
		string(ProvisioningStatesUpdating),
	}
}

func (s *ProvisioningStates) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProvisioningStates(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProvisioningStates(input string) (*ProvisioningStates, error) {
	vals := map[string]ProvisioningStates{
		"canceled":  ProvisioningStatesCanceled,
		"creating":  ProvisioningStatesCreating,
		"deleting":  ProvisioningStatesDeleting,
		"failed":    ProvisioningStatesFailed,
		"invalid":   ProvisioningStatesInvalid,
		"pending":   ProvisioningStatesPending,
		"succeeded": ProvisioningStatesSucceeded,
		"updating":  ProvisioningStatesUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningStates(input)
	return &out, nil
}

type VolumeCreateOption string

const (
	VolumeCreateOptionDisk             VolumeCreateOption = "Disk"
	VolumeCreateOptionDiskRestorePoint VolumeCreateOption = "DiskRestorePoint"
	VolumeCreateOptionDiskSnapshot     VolumeCreateOption = "DiskSnapshot"
	VolumeCreateOptionNone             VolumeCreateOption = "None"
	VolumeCreateOptionVolumeSnapshot   VolumeCreateOption = "VolumeSnapshot"
)

func PossibleValuesForVolumeCreateOption() []string {
	return []string{
		string(VolumeCreateOptionDisk),
		string(VolumeCreateOptionDiskRestorePoint),
		string(VolumeCreateOptionDiskSnapshot),
		string(VolumeCreateOptionNone),
		string(VolumeCreateOptionVolumeSnapshot),
	}
}

func (s *VolumeCreateOption) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVolumeCreateOption(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVolumeCreateOption(input string) (*VolumeCreateOption, error) {
	vals := map[string]VolumeCreateOption{
		"disk":             VolumeCreateOptionDisk,
		"diskrestorepoint": VolumeCreateOptionDiskRestorePoint,
		"disksnapshot":     VolumeCreateOptionDiskSnapshot,
		"none":             VolumeCreateOptionNone,
		"volumesnapshot":   VolumeCreateOptionVolumeSnapshot,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VolumeCreateOption(input)
	return &out, nil
}

type XMsDeleteSnapshots string

const (
	XMsDeleteSnapshotsFalse XMsDeleteSnapshots = "false"
	XMsDeleteSnapshotsTrue  XMsDeleteSnapshots = "true"
)

func PossibleValuesForXMsDeleteSnapshots() []string {
	return []string{
		string(XMsDeleteSnapshotsFalse),
		string(XMsDeleteSnapshotsTrue),
	}
}

func (s *XMsDeleteSnapshots) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseXMsDeleteSnapshots(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseXMsDeleteSnapshots(input string) (*XMsDeleteSnapshots, error) {
	vals := map[string]XMsDeleteSnapshots{
		"false": XMsDeleteSnapshotsFalse,
		"true":  XMsDeleteSnapshotsTrue,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := XMsDeleteSnapshots(input)
	return &out, nil
}

type XMsForceDelete string

const (
	XMsForceDeleteFalse XMsForceDelete = "false"
	XMsForceDeleteTrue  XMsForceDelete = "true"
)

func PossibleValuesForXMsForceDelete() []string {
	return []string{
		string(XMsForceDeleteFalse),
		string(XMsForceDeleteTrue),
	}
}

func (s *XMsForceDelete) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseXMsForceDelete(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseXMsForceDelete(input string) (*XMsForceDelete, error) {
	vals := map[string]XMsForceDelete{
		"false": XMsForceDeleteFalse,
		"true":  XMsForceDeleteTrue,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := XMsForceDelete(input)
	return &out, nil
}
