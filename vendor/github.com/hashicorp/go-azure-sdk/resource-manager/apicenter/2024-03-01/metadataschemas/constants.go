package metadataschemas

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MetadataAssignmentEntity string

const (
	MetadataAssignmentEntityApi         MetadataAssignmentEntity = "api"
	MetadataAssignmentEntityDeployment  MetadataAssignmentEntity = "deployment"
	MetadataAssignmentEntityEnvironment MetadataAssignmentEntity = "environment"
)

func PossibleValuesForMetadataAssignmentEntity() []string {
	return []string{
		string(MetadataAssignmentEntityApi),
		string(MetadataAssignmentEntityDeployment),
		string(MetadataAssignmentEntityEnvironment),
	}
}

func (s *MetadataAssignmentEntity) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMetadataAssignmentEntity(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMetadataAssignmentEntity(input string) (*MetadataAssignmentEntity, error) {
	vals := map[string]MetadataAssignmentEntity{
		"api":         MetadataAssignmentEntityApi,
		"deployment":  MetadataAssignmentEntityDeployment,
		"environment": MetadataAssignmentEntityEnvironment,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MetadataAssignmentEntity(input)
	return &out, nil
}
