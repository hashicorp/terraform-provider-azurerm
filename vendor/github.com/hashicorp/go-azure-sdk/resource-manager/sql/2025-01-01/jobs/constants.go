package jobs

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobScheduleType string

const (
	JobScheduleTypeOnce      JobScheduleType = "Once"
	JobScheduleTypeRecurring JobScheduleType = "Recurring"
)

func PossibleValuesForJobScheduleType() []string {
	return []string{
		string(JobScheduleTypeOnce),
		string(JobScheduleTypeRecurring),
	}
}

func (s *JobScheduleType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseJobScheduleType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseJobScheduleType(input string) (*JobScheduleType, error) {
	vals := map[string]JobScheduleType{
		"once":      JobScheduleTypeOnce,
		"recurring": JobScheduleTypeRecurring,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := JobScheduleType(input)
	return &out, nil
}
