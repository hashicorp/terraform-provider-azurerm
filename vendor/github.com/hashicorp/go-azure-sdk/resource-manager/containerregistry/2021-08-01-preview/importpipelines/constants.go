package importpipelines

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PipelineOptions string

const (
	PipelineOptionsContinueOnErrors          PipelineOptions = "ContinueOnErrors"
	PipelineOptionsDeleteSourceBlobOnSuccess PipelineOptions = "DeleteSourceBlobOnSuccess"
	PipelineOptionsOverwriteBlobs            PipelineOptions = "OverwriteBlobs"
	PipelineOptionsOverwriteTags             PipelineOptions = "OverwriteTags"
)

func PossibleValuesForPipelineOptions() []string {
	return []string{
		string(PipelineOptionsContinueOnErrors),
		string(PipelineOptionsDeleteSourceBlobOnSuccess),
		string(PipelineOptionsOverwriteBlobs),
		string(PipelineOptionsOverwriteTags),
	}
}

func parsePipelineOptions(input string) (*PipelineOptions, error) {
	vals := map[string]PipelineOptions{
		"continueonerrors":          PipelineOptionsContinueOnErrors,
		"deletesourceblobonsuccess": PipelineOptionsDeleteSourceBlobOnSuccess,
		"overwriteblobs":            PipelineOptionsOverwriteBlobs,
		"overwritetags":             PipelineOptionsOverwriteTags,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PipelineOptions(input)
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

type TriggerStatus string

const (
	TriggerStatusDisabled TriggerStatus = "Disabled"
	TriggerStatusEnabled  TriggerStatus = "Enabled"
)

func PossibleValuesForTriggerStatus() []string {
	return []string{
		string(TriggerStatusDisabled),
		string(TriggerStatusEnabled),
	}
}

func parseTriggerStatus(input string) (*TriggerStatus, error) {
	vals := map[string]TriggerStatus{
		"disabled": TriggerStatusDisabled,
		"enabled":  TriggerStatusEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TriggerStatus(input)
	return &out, nil
}
