package availabilitygrouplisteners

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Commit string

const (
	CommitAsynchronousCommit Commit = "Asynchronous_Commit"
	CommitSynchronousCommit  Commit = "Synchronous_Commit"
)

func PossibleValuesForCommit() []string {
	return []string{
		string(CommitAsynchronousCommit),
		string(CommitSynchronousCommit),
	}
}

func parseCommit(input string) (*Commit, error) {
	vals := map[string]Commit{
		"asynchronous_commit": CommitAsynchronousCommit,
		"synchronous_commit":  CommitSynchronousCommit,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Commit(input)
	return &out, nil
}

type Failover string

const (
	FailoverAutomatic Failover = "Automatic"
	FailoverManual    Failover = "Manual"
)

func PossibleValuesForFailover() []string {
	return []string{
		string(FailoverAutomatic),
		string(FailoverManual),
	}
}

func parseFailover(input string) (*Failover, error) {
	vals := map[string]Failover{
		"automatic": FailoverAutomatic,
		"manual":    FailoverManual,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Failover(input)
	return &out, nil
}

type ReadableSecondary string

const (
	ReadableSecondaryAll      ReadableSecondary = "All"
	ReadableSecondaryNo       ReadableSecondary = "No"
	ReadableSecondaryReadOnly ReadableSecondary = "Read_Only"
)

func PossibleValuesForReadableSecondary() []string {
	return []string{
		string(ReadableSecondaryAll),
		string(ReadableSecondaryNo),
		string(ReadableSecondaryReadOnly),
	}
}

func parseReadableSecondary(input string) (*ReadableSecondary, error) {
	vals := map[string]ReadableSecondary{
		"all":       ReadableSecondaryAll,
		"no":        ReadableSecondaryNo,
		"read_only": ReadableSecondaryReadOnly,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ReadableSecondary(input)
	return &out, nil
}

type Role string

const (
	RolePrimary   Role = "Primary"
	RoleSecondary Role = "Secondary"
)

func PossibleValuesForRole() []string {
	return []string{
		string(RolePrimary),
		string(RoleSecondary),
	}
}

func parseRole(input string) (*Role, error) {
	vals := map[string]Role{
		"primary":   RolePrimary,
		"secondary": RoleSecondary,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Role(input)
	return &out, nil
}
