package jobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CloudJobOperationPredicate struct {
	AllowTaskPreemption         *bool
	CreationTime                *string
	DisplayName                 *string
	ETag                        *string
	Id                          *string
	LastModified                *string
	MaxParallelTasks            *int64
	PreviousStateTransitionTime *string
	Priority                    *int64
	StateTransitionTime         *string
	Url                         *string
	UsesTaskDependencies        *bool
}

func (p CloudJobOperationPredicate) Matches(input CloudJob) bool {

	if p.AllowTaskPreemption != nil && (input.AllowTaskPreemption == nil || *p.AllowTaskPreemption != *input.AllowTaskPreemption) {
		return false
	}

	if p.CreationTime != nil && (input.CreationTime == nil || *p.CreationTime != *input.CreationTime) {
		return false
	}

	if p.DisplayName != nil && (input.DisplayName == nil || *p.DisplayName != *input.DisplayName) {
		return false
	}

	if p.ETag != nil && (input.ETag == nil || *p.ETag != *input.ETag) {
		return false
	}

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.LastModified != nil && (input.LastModified == nil || *p.LastModified != *input.LastModified) {
		return false
	}

	if p.MaxParallelTasks != nil && (input.MaxParallelTasks == nil || *p.MaxParallelTasks != *input.MaxParallelTasks) {
		return false
	}

	if p.PreviousStateTransitionTime != nil && (input.PreviousStateTransitionTime == nil || *p.PreviousStateTransitionTime != *input.PreviousStateTransitionTime) {
		return false
	}

	if p.Priority != nil && (input.Priority == nil || *p.Priority != *input.Priority) {
		return false
	}

	if p.StateTransitionTime != nil && (input.StateTransitionTime == nil || *p.StateTransitionTime != *input.StateTransitionTime) {
		return false
	}

	if p.Url != nil && (input.Url == nil || *p.Url != *input.Url) {
		return false
	}

	if p.UsesTaskDependencies != nil && (input.UsesTaskDependencies == nil || *p.UsesTaskDependencies != *input.UsesTaskDependencies) {
		return false
	}

	return true
}

type JobPreparationAndReleaseTaskExecutionInformationOperationPredicate struct {
	NodeId  *string
	NodeURL *string
	PoolId  *string
}

func (p JobPreparationAndReleaseTaskExecutionInformationOperationPredicate) Matches(input JobPreparationAndReleaseTaskExecutionInformation) bool {

	if p.NodeId != nil && (input.NodeId == nil || *p.NodeId != *input.NodeId) {
		return false
	}

	if p.NodeURL != nil && (input.NodeURL == nil || *p.NodeURL != *input.NodeURL) {
		return false
	}

	if p.PoolId != nil && (input.PoolId == nil || *p.PoolId != *input.PoolId) {
		return false
	}

	return true
}
