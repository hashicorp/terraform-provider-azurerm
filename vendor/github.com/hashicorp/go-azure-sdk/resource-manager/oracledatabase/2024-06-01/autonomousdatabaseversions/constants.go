package autonomousdatabaseversions

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkloadType string

const (
	WorkloadTypeAJD  WorkloadType = "AJD"
	WorkloadTypeAPEX WorkloadType = "APEX"
	WorkloadTypeDW   WorkloadType = "DW"
	WorkloadTypeOLTP WorkloadType = "OLTP"
)

func PossibleValuesForWorkloadType() []string {
	return []string{
		string(WorkloadTypeAJD),
		string(WorkloadTypeAPEX),
		string(WorkloadTypeDW),
		string(WorkloadTypeOLTP),
	}
}

func (s *WorkloadType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseWorkloadType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseWorkloadType(input string) (*WorkloadType, error) {
	vals := map[string]WorkloadType{
		"ajd":  WorkloadTypeAJD,
		"apex": WorkloadTypeAPEX,
		"dw":   WorkloadTypeDW,
		"oltp": WorkloadTypeOLTP,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := WorkloadType(input)
	return &out, nil
}
