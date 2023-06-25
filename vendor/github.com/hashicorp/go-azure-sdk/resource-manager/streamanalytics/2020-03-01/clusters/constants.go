package clusters

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterProvisioningState string

const (
	ClusterProvisioningStateCanceled   ClusterProvisioningState = "Canceled"
	ClusterProvisioningStateFailed     ClusterProvisioningState = "Failed"
	ClusterProvisioningStateInProgress ClusterProvisioningState = "InProgress"
	ClusterProvisioningStateSucceeded  ClusterProvisioningState = "Succeeded"
)

func PossibleValuesForClusterProvisioningState() []string {
	return []string{
		string(ClusterProvisioningStateCanceled),
		string(ClusterProvisioningStateFailed),
		string(ClusterProvisioningStateInProgress),
		string(ClusterProvisioningStateSucceeded),
	}
}

func parseClusterProvisioningState(input string) (*ClusterProvisioningState, error) {
	vals := map[string]ClusterProvisioningState{
		"canceled":   ClusterProvisioningStateCanceled,
		"failed":     ClusterProvisioningStateFailed,
		"inprogress": ClusterProvisioningStateInProgress,
		"succeeded":  ClusterProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ClusterProvisioningState(input)
	return &out, nil
}

type ClusterSkuName string

const (
	ClusterSkuNameDefault ClusterSkuName = "Default"
)

func PossibleValuesForClusterSkuName() []string {
	return []string{
		string(ClusterSkuNameDefault),
	}
}

func parseClusterSkuName(input string) (*ClusterSkuName, error) {
	vals := map[string]ClusterSkuName{
		"default": ClusterSkuNameDefault,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ClusterSkuName(input)
	return &out, nil
}

type JobState string

const (
	JobStateCreated    JobState = "Created"
	JobStateDegraded   JobState = "Degraded"
	JobStateDeleting   JobState = "Deleting"
	JobStateFailed     JobState = "Failed"
	JobStateRestarting JobState = "Restarting"
	JobStateRunning    JobState = "Running"
	JobStateScaling    JobState = "Scaling"
	JobStateStarting   JobState = "Starting"
	JobStateStopped    JobState = "Stopped"
	JobStateStopping   JobState = "Stopping"
)

func PossibleValuesForJobState() []string {
	return []string{
		string(JobStateCreated),
		string(JobStateDegraded),
		string(JobStateDeleting),
		string(JobStateFailed),
		string(JobStateRestarting),
		string(JobStateRunning),
		string(JobStateScaling),
		string(JobStateStarting),
		string(JobStateStopped),
		string(JobStateStopping),
	}
}

func parseJobState(input string) (*JobState, error) {
	vals := map[string]JobState{
		"created":    JobStateCreated,
		"degraded":   JobStateDegraded,
		"deleting":   JobStateDeleting,
		"failed":     JobStateFailed,
		"restarting": JobStateRestarting,
		"running":    JobStateRunning,
		"scaling":    JobStateScaling,
		"starting":   JobStateStarting,
		"stopped":    JobStateStopped,
		"stopping":   JobStateStopping,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := JobState(input)
	return &out, nil
}
