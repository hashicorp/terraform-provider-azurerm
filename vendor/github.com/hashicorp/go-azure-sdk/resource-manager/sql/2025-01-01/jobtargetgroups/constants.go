package jobtargetgroups

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobTargetGroupMembershipType string

const (
	JobTargetGroupMembershipTypeExclude JobTargetGroupMembershipType = "Exclude"
	JobTargetGroupMembershipTypeInclude JobTargetGroupMembershipType = "Include"
)

func PossibleValuesForJobTargetGroupMembershipType() []string {
	return []string{
		string(JobTargetGroupMembershipTypeExclude),
		string(JobTargetGroupMembershipTypeInclude),
	}
}

func (s *JobTargetGroupMembershipType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseJobTargetGroupMembershipType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseJobTargetGroupMembershipType(input string) (*JobTargetGroupMembershipType, error) {
	vals := map[string]JobTargetGroupMembershipType{
		"exclude": JobTargetGroupMembershipTypeExclude,
		"include": JobTargetGroupMembershipTypeInclude,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := JobTargetGroupMembershipType(input)
	return &out, nil
}

type JobTargetType string

const (
	JobTargetTypeSqlDatabase    JobTargetType = "SqlDatabase"
	JobTargetTypeSqlElasticPool JobTargetType = "SqlElasticPool"
	JobTargetTypeSqlServer      JobTargetType = "SqlServer"
	JobTargetTypeSqlShardMap    JobTargetType = "SqlShardMap"
	JobTargetTypeTargetGroup    JobTargetType = "TargetGroup"
)

func PossibleValuesForJobTargetType() []string {
	return []string{
		string(JobTargetTypeSqlDatabase),
		string(JobTargetTypeSqlElasticPool),
		string(JobTargetTypeSqlServer),
		string(JobTargetTypeSqlShardMap),
		string(JobTargetTypeTargetGroup),
	}
}

func (s *JobTargetType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseJobTargetType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseJobTargetType(input string) (*JobTargetType, error) {
	vals := map[string]JobTargetType{
		"sqldatabase":    JobTargetTypeSqlDatabase,
		"sqlelasticpool": JobTargetTypeSqlElasticPool,
		"sqlserver":      JobTargetTypeSqlServer,
		"sqlshardmap":    JobTargetTypeSqlShardMap,
		"targetgroup":    JobTargetTypeTargetGroup,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := JobTargetType(input)
	return &out, nil
}
