package managementlocks

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LockLevel string

const (
	LockLevelCanNotDelete LockLevel = "CanNotDelete"
	LockLevelNotSpecified LockLevel = "NotSpecified"
	LockLevelReadOnly     LockLevel = "ReadOnly"
)

func PossibleValuesForLockLevel() []string {
	return []string{
		string(LockLevelCanNotDelete),
		string(LockLevelNotSpecified),
		string(LockLevelReadOnly),
	}
}

func parseLockLevel(input string) (*LockLevel, error) {
	vals := map[string]LockLevel{
		"cannotdelete": LockLevelCanNotDelete,
		"notspecified": LockLevelNotSpecified,
		"readonly":     LockLevelReadOnly,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LockLevel(input)
	return &out, nil
}
