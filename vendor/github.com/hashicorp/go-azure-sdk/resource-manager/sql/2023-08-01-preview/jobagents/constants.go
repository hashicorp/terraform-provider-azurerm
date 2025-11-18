package jobagents

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobAgentState string

const (
	JobAgentStateCreating JobAgentState = "Creating"
	JobAgentStateDeleting JobAgentState = "Deleting"
	JobAgentStateDisabled JobAgentState = "Disabled"
	JobAgentStateReady    JobAgentState = "Ready"
	JobAgentStateUpdating JobAgentState = "Updating"
)

func PossibleValuesForJobAgentState() []string {
	return []string{
		string(JobAgentStateCreating),
		string(JobAgentStateDeleting),
		string(JobAgentStateDisabled),
		string(JobAgentStateReady),
		string(JobAgentStateUpdating),
	}
}

func (s *JobAgentState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseJobAgentState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseJobAgentState(input string) (*JobAgentState, error) {
	vals := map[string]JobAgentState{
		"creating": JobAgentStateCreating,
		"deleting": JobAgentStateDeleting,
		"disabled": JobAgentStateDisabled,
		"ready":    JobAgentStateReady,
		"updating": JobAgentStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := JobAgentState(input)
	return &out, nil
}
