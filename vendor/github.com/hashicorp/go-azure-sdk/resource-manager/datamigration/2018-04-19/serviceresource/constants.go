package serviceresource

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CommandState string

const (
	CommandStateAccepted  CommandState = "Accepted"
	CommandStateFailed    CommandState = "Failed"
	CommandStateRunning   CommandState = "Running"
	CommandStateSucceeded CommandState = "Succeeded"
	CommandStateUnknown   CommandState = "Unknown"
)

func PossibleValuesForCommandState() []string {
	return []string{
		string(CommandStateAccepted),
		string(CommandStateFailed),
		string(CommandStateRunning),
		string(CommandStateSucceeded),
		string(CommandStateUnknown),
	}
}

func parseCommandState(input string) (*CommandState, error) {
	vals := map[string]CommandState{
		"accepted":  CommandStateAccepted,
		"failed":    CommandStateFailed,
		"running":   CommandStateRunning,
		"succeeded": CommandStateSucceeded,
		"unknown":   CommandStateUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CommandState(input)
	return &out, nil
}

type ServiceProvisioningState string

const (
	ServiceProvisioningStateAccepted      ServiceProvisioningState = "Accepted"
	ServiceProvisioningStateDeleting      ServiceProvisioningState = "Deleting"
	ServiceProvisioningStateDeploying     ServiceProvisioningState = "Deploying"
	ServiceProvisioningStateFailed        ServiceProvisioningState = "Failed"
	ServiceProvisioningStateFailedToStart ServiceProvisioningState = "FailedToStart"
	ServiceProvisioningStateFailedToStop  ServiceProvisioningState = "FailedToStop"
	ServiceProvisioningStateStarting      ServiceProvisioningState = "Starting"
	ServiceProvisioningStateStopped       ServiceProvisioningState = "Stopped"
	ServiceProvisioningStateStopping      ServiceProvisioningState = "Stopping"
	ServiceProvisioningStateSucceeded     ServiceProvisioningState = "Succeeded"
)

func PossibleValuesForServiceProvisioningState() []string {
	return []string{
		string(ServiceProvisioningStateAccepted),
		string(ServiceProvisioningStateDeleting),
		string(ServiceProvisioningStateDeploying),
		string(ServiceProvisioningStateFailed),
		string(ServiceProvisioningStateFailedToStart),
		string(ServiceProvisioningStateFailedToStop),
		string(ServiceProvisioningStateStarting),
		string(ServiceProvisioningStateStopped),
		string(ServiceProvisioningStateStopping),
		string(ServiceProvisioningStateSucceeded),
	}
}

func parseServiceProvisioningState(input string) (*ServiceProvisioningState, error) {
	vals := map[string]ServiceProvisioningState{
		"accepted":      ServiceProvisioningStateAccepted,
		"deleting":      ServiceProvisioningStateDeleting,
		"deploying":     ServiceProvisioningStateDeploying,
		"failed":        ServiceProvisioningStateFailed,
		"failedtostart": ServiceProvisioningStateFailedToStart,
		"failedtostop":  ServiceProvisioningStateFailedToStop,
		"starting":      ServiceProvisioningStateStarting,
		"stopped":       ServiceProvisioningStateStopped,
		"stopping":      ServiceProvisioningStateStopping,
		"succeeded":     ServiceProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ServiceProvisioningState(input)
	return &out, nil
}

type ServiceScalability string

const (
	ServiceScalabilityAutomatic ServiceScalability = "automatic"
	ServiceScalabilityManual    ServiceScalability = "manual"
	ServiceScalabilityNone      ServiceScalability = "none"
)

func PossibleValuesForServiceScalability() []string {
	return []string{
		string(ServiceScalabilityAutomatic),
		string(ServiceScalabilityManual),
		string(ServiceScalabilityNone),
	}
}

func parseServiceScalability(input string) (*ServiceScalability, error) {
	vals := map[string]ServiceScalability{
		"automatic": ServiceScalabilityAutomatic,
		"manual":    ServiceScalabilityManual,
		"none":      ServiceScalabilityNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ServiceScalability(input)
	return &out, nil
}

type TaskState string

const (
	TaskStateCanceled              TaskState = "Canceled"
	TaskStateFailed                TaskState = "Failed"
	TaskStateFailedInputValidation TaskState = "FailedInputValidation"
	TaskStateFaulted               TaskState = "Faulted"
	TaskStateQueued                TaskState = "Queued"
	TaskStateRunning               TaskState = "Running"
	TaskStateSucceeded             TaskState = "Succeeded"
	TaskStateUnknown               TaskState = "Unknown"
)

func PossibleValuesForTaskState() []string {
	return []string{
		string(TaskStateCanceled),
		string(TaskStateFailed),
		string(TaskStateFailedInputValidation),
		string(TaskStateFaulted),
		string(TaskStateQueued),
		string(TaskStateRunning),
		string(TaskStateSucceeded),
		string(TaskStateUnknown),
	}
}

func parseTaskState(input string) (*TaskState, error) {
	vals := map[string]TaskState{
		"canceled":              TaskStateCanceled,
		"failed":                TaskStateFailed,
		"failedinputvalidation": TaskStateFailedInputValidation,
		"faulted":               TaskStateFaulted,
		"queued":                TaskStateQueued,
		"running":               TaskStateRunning,
		"succeeded":             TaskStateSucceeded,
		"unknown":               TaskStateUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TaskState(input)
	return &out, nil
}
