package managementlocks

import (
	"encoding/json"
	"fmt"
	"strings"
)

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

func (s *LockLevel) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLockLevel(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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
