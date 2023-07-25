package projectresource

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProjectProvisioningState string

const (
	ProjectProvisioningStateDeleting  ProjectProvisioningState = "Deleting"
	ProjectProvisioningStateSucceeded ProjectProvisioningState = "Succeeded"
)

func PossibleValuesForProjectProvisioningState() []string {
	return []string{
		string(ProjectProvisioningStateDeleting),
		string(ProjectProvisioningStateSucceeded),
	}
}

func (s *ProjectProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProjectProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProjectProvisioningState(input string) (*ProjectProvisioningState, error) {
	vals := map[string]ProjectProvisioningState{
		"deleting":  ProjectProvisioningStateDeleting,
		"succeeded": ProjectProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProjectProvisioningState(input)
	return &out, nil
}

type ProjectSourcePlatform string

const (
	ProjectSourcePlatformSQL     ProjectSourcePlatform = "SQL"
	ProjectSourcePlatformUnknown ProjectSourcePlatform = "Unknown"
)

func PossibleValuesForProjectSourcePlatform() []string {
	return []string{
		string(ProjectSourcePlatformSQL),
		string(ProjectSourcePlatformUnknown),
	}
}

func (s *ProjectSourcePlatform) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProjectSourcePlatform(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProjectSourcePlatform(input string) (*ProjectSourcePlatform, error) {
	vals := map[string]ProjectSourcePlatform{
		"sql":     ProjectSourcePlatformSQL,
		"unknown": ProjectSourcePlatformUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProjectSourcePlatform(input)
	return &out, nil
}

type ProjectTargetPlatform string

const (
	ProjectTargetPlatformSQLDB   ProjectTargetPlatform = "SQLDB"
	ProjectTargetPlatformUnknown ProjectTargetPlatform = "Unknown"
)

func PossibleValuesForProjectTargetPlatform() []string {
	return []string{
		string(ProjectTargetPlatformSQLDB),
		string(ProjectTargetPlatformUnknown),
	}
}

func (s *ProjectTargetPlatform) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProjectTargetPlatform(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProjectTargetPlatform(input string) (*ProjectTargetPlatform, error) {
	vals := map[string]ProjectTargetPlatform{
		"sqldb":   ProjectTargetPlatformSQLDB,
		"unknown": ProjectTargetPlatformUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProjectTargetPlatform(input)
	return &out, nil
}
