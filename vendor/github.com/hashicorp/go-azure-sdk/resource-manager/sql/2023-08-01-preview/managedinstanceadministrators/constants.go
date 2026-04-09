package managedinstanceadministrators

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedInstanceAdministratorType string

const (
	ManagedInstanceAdministratorTypeActiveDirectory ManagedInstanceAdministratorType = "ActiveDirectory"
)

func PossibleValuesForManagedInstanceAdministratorType() []string {
	return []string{
		string(ManagedInstanceAdministratorTypeActiveDirectory),
	}
}

func (s *ManagedInstanceAdministratorType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseManagedInstanceAdministratorType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseManagedInstanceAdministratorType(input string) (*ManagedInstanceAdministratorType, error) {
	vals := map[string]ManagedInstanceAdministratorType{
		"activedirectory": ManagedInstanceAdministratorTypeActiveDirectory,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ManagedInstanceAdministratorType(input)
	return &out, nil
}
