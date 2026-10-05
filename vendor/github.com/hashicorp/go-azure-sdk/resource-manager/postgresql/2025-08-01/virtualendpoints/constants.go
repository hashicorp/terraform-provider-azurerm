package virtualendpoints

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualEndpointType string

const (
	VirtualEndpointTypeReadWrite VirtualEndpointType = "ReadWrite"
)

func PossibleValuesForVirtualEndpointType() []string {
	return []string{
		string(VirtualEndpointTypeReadWrite),
	}
}

func (s *VirtualEndpointType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVirtualEndpointType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVirtualEndpointType(input string) (*VirtualEndpointType, error) {
	vals := map[string]VirtualEndpointType{
		"readwrite": VirtualEndpointTypeReadWrite,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VirtualEndpointType(input)
	return &out, nil
}
