package backups

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BackupType string

const (
	BackupTypeManual    BackupType = "Manual"
	BackupTypeScheduled BackupType = "Scheduled"
)

func PossibleValuesForBackupType() []string {
	return []string{
		string(BackupTypeManual),
		string(BackupTypeScheduled),
	}
}

func (s *BackupType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBackupType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBackupType(input string) (*BackupType, error) {
	vals := map[string]BackupType{
		"manual":    BackupTypeManual,
		"scheduled": BackupTypeScheduled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BackupType(input)
	return &out, nil
}

type MirrorState string

const (
	MirrorStateBroken        MirrorState = "Broken"
	MirrorStateMirrored      MirrorState = "Mirrored"
	MirrorStateUninitialized MirrorState = "Uninitialized"
)

func PossibleValuesForMirrorState() []string {
	return []string{
		string(MirrorStateBroken),
		string(MirrorStateMirrored),
		string(MirrorStateUninitialized),
	}
}

func (s *MirrorState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMirrorState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMirrorState(input string) (*MirrorState, error) {
	vals := map[string]MirrorState{
		"broken":        MirrorStateBroken,
		"mirrored":      MirrorStateMirrored,
		"uninitialized": MirrorStateUninitialized,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MirrorState(input)
	return &out, nil
}

type RelationshipStatus string

const (
	RelationshipStatusFailed       RelationshipStatus = "Failed"
	RelationshipStatusIdle         RelationshipStatus = "Idle"
	RelationshipStatusTransferring RelationshipStatus = "Transferring"
	RelationshipStatusUnknown      RelationshipStatus = "Unknown"
)

func PossibleValuesForRelationshipStatus() []string {
	return []string{
		string(RelationshipStatusFailed),
		string(RelationshipStatusIdle),
		string(RelationshipStatusTransferring),
		string(RelationshipStatusUnknown),
	}
}

func (s *RelationshipStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRelationshipStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRelationshipStatus(input string) (*RelationshipStatus, error) {
	vals := map[string]RelationshipStatus{
		"failed":       RelationshipStatusFailed,
		"idle":         RelationshipStatusIdle,
		"transferring": RelationshipStatusTransferring,
		"unknown":      RelationshipStatusUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RelationshipStatus(input)
	return &out, nil
}
