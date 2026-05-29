package hybridrunbookworker

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkerType string

const (
	WorkerTypeHybridVOne WorkerType = "HybridV1"
	WorkerTypeHybridVTwo WorkerType = "HybridV2"
)

func PossibleValuesForWorkerType() []string {
	return []string{
		string(WorkerTypeHybridVOne),
		string(WorkerTypeHybridVTwo),
	}
}

func (s *WorkerType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseWorkerType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseWorkerType(input string) (*WorkerType, error) {
	vals := map[string]WorkerType{
		"hybridv1": WorkerTypeHybridVOne,
		"hybridv2": WorkerTypeHybridVTwo,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := WorkerType(input)
	return &out, nil
}
