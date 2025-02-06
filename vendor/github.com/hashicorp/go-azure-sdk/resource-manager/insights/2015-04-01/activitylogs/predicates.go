package activitylogs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EventDataOperationPredicate struct {
	Caller              *string
	CorrelationId       *string
	Description         *string
	EventDataId         *string
	EventTimestamp      *string
	Id                  *string
	OperationId         *string
	ResourceGroupName   *string
	ResourceId          *string
	SubmissionTimestamp *string
	SubscriptionId      *string
	TenantId            *string
}

func (p EventDataOperationPredicate) Matches(input EventData) bool {

	if p.Caller != nil && (input.Caller == nil || *p.Caller != *input.Caller) {
		return false
	}

	if p.CorrelationId != nil && (input.CorrelationId == nil || *p.CorrelationId != *input.CorrelationId) {
		return false
	}

	if p.Description != nil && (input.Description == nil || *p.Description != *input.Description) {
		return false
	}

	if p.EventDataId != nil && (input.EventDataId == nil || *p.EventDataId != *input.EventDataId) {
		return false
	}

	if p.EventTimestamp != nil && (input.EventTimestamp == nil || *p.EventTimestamp != *input.EventTimestamp) {
		return false
	}

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.OperationId != nil && (input.OperationId == nil || *p.OperationId != *input.OperationId) {
		return false
	}

	if p.ResourceGroupName != nil && (input.ResourceGroupName == nil || *p.ResourceGroupName != *input.ResourceGroupName) {
		return false
	}

	if p.ResourceId != nil && (input.ResourceId == nil || *p.ResourceId != *input.ResourceId) {
		return false
	}

	if p.SubmissionTimestamp != nil && (input.SubmissionTimestamp == nil || *p.SubmissionTimestamp != *input.SubmissionTimestamp) {
		return false
	}

	if p.SubscriptionId != nil && (input.SubscriptionId == nil || *p.SubscriptionId != *input.SubscriptionId) {
		return false
	}

	if p.TenantId != nil && (input.TenantId == nil || *p.TenantId != *input.TenantId) {
		return false
	}

	return true
}
