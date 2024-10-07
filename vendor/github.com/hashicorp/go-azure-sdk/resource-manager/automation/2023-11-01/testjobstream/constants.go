package testjobstream

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobStreamType string

const (
	JobStreamTypeAny      JobStreamType = "Any"
	JobStreamTypeDebug    JobStreamType = "Debug"
	JobStreamTypeError    JobStreamType = "Error"
	JobStreamTypeOutput   JobStreamType = "Output"
	JobStreamTypeProgress JobStreamType = "Progress"
	JobStreamTypeVerbose  JobStreamType = "Verbose"
	JobStreamTypeWarning  JobStreamType = "Warning"
)

func PossibleValuesForJobStreamType() []string {
	return []string{
		string(JobStreamTypeAny),
		string(JobStreamTypeDebug),
		string(JobStreamTypeError),
		string(JobStreamTypeOutput),
		string(JobStreamTypeProgress),
		string(JobStreamTypeVerbose),
		string(JobStreamTypeWarning),
	}
}

func (s *JobStreamType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseJobStreamType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseJobStreamType(input string) (*JobStreamType, error) {
	vals := map[string]JobStreamType{
		"any":      JobStreamTypeAny,
		"debug":    JobStreamTypeDebug,
		"error":    JobStreamTypeError,
		"output":   JobStreamTypeOutput,
		"progress": JobStreamTypeProgress,
		"verbose":  JobStreamTypeVerbose,
		"warning":  JobStreamTypeWarning,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := JobStreamType(input)
	return &out, nil
}
