package jobsteps

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobStepActionSource string

const (
	JobStepActionSourceInline JobStepActionSource = "Inline"
)

func PossibleValuesForJobStepActionSource() []string {
	return []string{
		string(JobStepActionSourceInline),
	}
}

func (s *JobStepActionSource) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseJobStepActionSource(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseJobStepActionSource(input string) (*JobStepActionSource, error) {
	vals := map[string]JobStepActionSource{
		"inline": JobStepActionSourceInline,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := JobStepActionSource(input)
	return &out, nil
}

type JobStepActionType string

const (
	JobStepActionTypeTSql JobStepActionType = "TSql"
)

func PossibleValuesForJobStepActionType() []string {
	return []string{
		string(JobStepActionTypeTSql),
	}
}

func (s *JobStepActionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseJobStepActionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseJobStepActionType(input string) (*JobStepActionType, error) {
	vals := map[string]JobStepActionType{
		"tsql": JobStepActionTypeTSql,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := JobStepActionType(input)
	return &out, nil
}

type JobStepOutputType string

const (
	JobStepOutputTypeSqlDatabase JobStepOutputType = "SqlDatabase"
)

func PossibleValuesForJobStepOutputType() []string {
	return []string{
		string(JobStepOutputTypeSqlDatabase),
	}
}

func (s *JobStepOutputType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseJobStepOutputType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseJobStepOutputType(input string) (*JobStepOutputType, error) {
	vals := map[string]JobStepOutputType{
		"sqldatabase": JobStepOutputTypeSqlDatabase,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := JobStepOutputType(input)
	return &out, nil
}
