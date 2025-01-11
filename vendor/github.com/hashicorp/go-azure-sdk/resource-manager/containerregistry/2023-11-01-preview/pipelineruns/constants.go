package pipelineruns

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PipelineRunSourceType string

const (
	PipelineRunSourceTypeAzureStorageBlob PipelineRunSourceType = "AzureStorageBlob"
)

func PossibleValuesForPipelineRunSourceType() []string {
	return []string{
		string(PipelineRunSourceTypeAzureStorageBlob),
	}
}

func (s *PipelineRunSourceType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePipelineRunSourceType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePipelineRunSourceType(input string) (*PipelineRunSourceType, error) {
	vals := map[string]PipelineRunSourceType{
		"azurestorageblob": PipelineRunSourceTypeAzureStorageBlob,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PipelineRunSourceType(input)
	return &out, nil
}

type PipelineRunTargetType string

const (
	PipelineRunTargetTypeAzureStorageBlob PipelineRunTargetType = "AzureStorageBlob"
)

func PossibleValuesForPipelineRunTargetType() []string {
	return []string{
		string(PipelineRunTargetTypeAzureStorageBlob),
	}
}

func (s *PipelineRunTargetType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePipelineRunTargetType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePipelineRunTargetType(input string) (*PipelineRunTargetType, error) {
	vals := map[string]PipelineRunTargetType{
		"azurestorageblob": PipelineRunTargetTypeAzureStorageBlob,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PipelineRunTargetType(input)
	return &out, nil
}

type PipelineSourceType string

const (
	PipelineSourceTypeAzureStorageBlobContainer PipelineSourceType = "AzureStorageBlobContainer"
)

func PossibleValuesForPipelineSourceType() []string {
	return []string{
		string(PipelineSourceTypeAzureStorageBlobContainer),
	}
}

func (s *PipelineSourceType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePipelineSourceType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePipelineSourceType(input string) (*PipelineSourceType, error) {
	vals := map[string]PipelineSourceType{
		"azurestorageblobcontainer": PipelineSourceTypeAzureStorageBlobContainer,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PipelineSourceType(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateCanceled  ProvisioningState = "Canceled"
	ProvisioningStateCreating  ProvisioningState = "Creating"
	ProvisioningStateDeleting  ProvisioningState = "Deleting"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
	ProvisioningStateUpdating  ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateCanceled),
		string(ProvisioningStateCreating),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateSucceeded),
		string(ProvisioningStateUpdating),
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
		"creating":  ProvisioningStateCreating,
		"deleting":  ProvisioningStateDeleting,
		"failed":    ProvisioningStateFailed,
		"succeeded": ProvisioningStateSucceeded,
		"updating":  ProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}
