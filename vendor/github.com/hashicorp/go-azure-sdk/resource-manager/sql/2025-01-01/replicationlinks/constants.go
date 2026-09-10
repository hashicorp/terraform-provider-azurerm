package replicationlinks

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReplicationLinkType string

const (
	ReplicationLinkTypeGEO     ReplicationLinkType = "GEO"
	ReplicationLinkTypeNAMED   ReplicationLinkType = "NAMED"
	ReplicationLinkTypeSTANDBY ReplicationLinkType = "STANDBY"
)

func PossibleValuesForReplicationLinkType() []string {
	return []string{
		string(ReplicationLinkTypeGEO),
		string(ReplicationLinkTypeNAMED),
		string(ReplicationLinkTypeSTANDBY),
	}
}

func (s *ReplicationLinkType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseReplicationLinkType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseReplicationLinkType(input string) (*ReplicationLinkType, error) {
	vals := map[string]ReplicationLinkType{
		"geo":     ReplicationLinkTypeGEO,
		"named":   ReplicationLinkTypeNAMED,
		"standby": ReplicationLinkTypeSTANDBY,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ReplicationLinkType(input)
	return &out, nil
}

type ReplicationRole string

const (
	ReplicationRoleCopy                 ReplicationRole = "Copy"
	ReplicationRoleNonReadableSecondary ReplicationRole = "NonReadableSecondary"
	ReplicationRolePrimary              ReplicationRole = "Primary"
	ReplicationRoleSecondary            ReplicationRole = "Secondary"
	ReplicationRoleSource               ReplicationRole = "Source"
)

func PossibleValuesForReplicationRole() []string {
	return []string{
		string(ReplicationRoleCopy),
		string(ReplicationRoleNonReadableSecondary),
		string(ReplicationRolePrimary),
		string(ReplicationRoleSecondary),
		string(ReplicationRoleSource),
	}
}

func (s *ReplicationRole) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseReplicationRole(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseReplicationRole(input string) (*ReplicationRole, error) {
	vals := map[string]ReplicationRole{
		"copy":                 ReplicationRoleCopy,
		"nonreadablesecondary": ReplicationRoleNonReadableSecondary,
		"primary":              ReplicationRolePrimary,
		"secondary":            ReplicationRoleSecondary,
		"source":               ReplicationRoleSource,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ReplicationRole(input)
	return &out, nil
}

type ReplicationState string

const (
	ReplicationStateCATCHUP   ReplicationState = "CATCH_UP"
	ReplicationStatePENDING   ReplicationState = "PENDING"
	ReplicationStateSEEDING   ReplicationState = "SEEDING"
	ReplicationStateSUSPENDED ReplicationState = "SUSPENDED"
)

func PossibleValuesForReplicationState() []string {
	return []string{
		string(ReplicationStateCATCHUP),
		string(ReplicationStatePENDING),
		string(ReplicationStateSEEDING),
		string(ReplicationStateSUSPENDED),
	}
}

func (s *ReplicationState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseReplicationState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseReplicationState(input string) (*ReplicationState, error) {
	vals := map[string]ReplicationState{
		"catch_up":  ReplicationStateCATCHUP,
		"pending":   ReplicationStatePENDING,
		"seeding":   ReplicationStateSEEDING,
		"suspended": ReplicationStateSUSPENDED,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ReplicationState(input)
	return &out, nil
}
