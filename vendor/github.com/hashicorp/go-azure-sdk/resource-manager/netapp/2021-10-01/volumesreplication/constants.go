package volumesreplication

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

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
	RelationshipStatusIdle         RelationshipStatus = "Idle"
	RelationshipStatusTransferring RelationshipStatus = "Transferring"
)

func PossibleValuesForRelationshipStatus() []string {
	return []string{
		string(RelationshipStatusIdle),
		string(RelationshipStatusTransferring),
	}
}

func parseRelationshipStatus(input string) (*RelationshipStatus, error) {
	vals := map[string]RelationshipStatus{
		"idle":         RelationshipStatusIdle,
		"transferring": RelationshipStatusTransferring,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RelationshipStatus(input)
	return &out, nil
}
