package agentpools

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContainerServiceOSTypes string

const (
	ContainerServiceOSTypesLinux   ContainerServiceOSTypes = "Linux"
	ContainerServiceOSTypesWindows ContainerServiceOSTypes = "Windows"
)

func PossibleValuesForContainerServiceOSTypes() []string {
	return []string{
		string(ContainerServiceOSTypesLinux),
		string(ContainerServiceOSTypesWindows),
	}
}

func (s *ContainerServiceOSTypes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseContainerServiceOSTypes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseContainerServiceOSTypes(input string) (*ContainerServiceOSTypes, error) {
	vals := map[string]ContainerServiceOSTypes{
		"linux":   ContainerServiceOSTypesLinux,
		"windows": ContainerServiceOSTypesWindows,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ContainerServiceOSTypes(input)
	return &out, nil
}
