package dbsystemshapes

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ComputeModel string

const (
	ComputeModelECPU ComputeModel = "ECPU"
	ComputeModelOCPU ComputeModel = "OCPU"
)

func PossibleValuesForComputeModel() []string {
	return []string{
		string(ComputeModelECPU),
		string(ComputeModelOCPU),
	}
}

func (s *ComputeModel) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseComputeModel(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseComputeModel(input string) (*ComputeModel, error) {
	vals := map[string]ComputeModel{
		"ecpu": ComputeModelECPU,
		"ocpu": ComputeModelOCPU,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ComputeModel(input)
	return &out, nil
}
