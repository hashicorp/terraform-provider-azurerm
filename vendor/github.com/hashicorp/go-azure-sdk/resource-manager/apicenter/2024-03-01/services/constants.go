package services

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

type MetadataSchemaExportFormat string

const (
	MetadataSchemaExportFormatInline MetadataSchemaExportFormat = "inline"
	MetadataSchemaExportFormatLink   MetadataSchemaExportFormat = "link"
)

func PossibleValuesForMetadataSchemaExportFormat() []string {
	return []string{
		string(MetadataSchemaExportFormatInline),
		string(MetadataSchemaExportFormatLink),
	}
}

func (s *MetadataSchemaExportFormat) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMetadataSchemaExportFormat(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMetadataSchemaExportFormat(input string) (*MetadataSchemaExportFormat, error) {
	vals := map[string]MetadataSchemaExportFormat{
		"inline": MetadataSchemaExportFormatInline,
		"link":   MetadataSchemaExportFormatLink,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MetadataSchemaExportFormat(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateCanceled  ProvisioningState = "Canceled"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateCanceled),
		string(ProvisioningStateFailed),
		string(ProvisioningStateSucceeded),
	}
}

func (s *ProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"canceled":  ProvisioningStateCanceled,
		"failed":    ProvisioningStateFailed,
		"succeeded": ProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}
