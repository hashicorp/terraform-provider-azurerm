package bastionhosts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BastionActiveSessionOperationPredicate struct {
	ResourceType          *string
	SessionDurationInMins *float64
	SessionId             *string
	StartTime             *interface{}
	TargetHostName        *string
	TargetIPAddress       *string
	TargetResourceGroup   *string
	TargetResourceId      *string
	TargetSubscriptionId  *string
	UserName              *string
}

func (p BastionActiveSessionOperationPredicate) Matches(input BastionActiveSession) bool {

	if p.ResourceType != nil && (input.ResourceType == nil || *p.ResourceType != *input.ResourceType) {
		return false
	}

	if p.SessionDurationInMins != nil && (input.SessionDurationInMins == nil || *p.SessionDurationInMins != *input.SessionDurationInMins) {
		return false
	}

	if p.SessionId != nil && (input.SessionId == nil || *p.SessionId != *input.SessionId) {
		return false
	}

	if p.StartTime != nil && (input.StartTime == nil || *p.StartTime != *input.StartTime) {
		return false
	}

	if p.TargetHostName != nil && (input.TargetHostName == nil || *p.TargetHostName != *input.TargetHostName) {
		return false
	}

	if p.TargetIPAddress != nil && (input.TargetIPAddress == nil || *p.TargetIPAddress != *input.TargetIPAddress) {
		return false
	}

	if p.TargetResourceGroup != nil && (input.TargetResourceGroup == nil || *p.TargetResourceGroup != *input.TargetResourceGroup) {
		return false
	}

	if p.TargetResourceId != nil && (input.TargetResourceId == nil || *p.TargetResourceId != *input.TargetResourceId) {
		return false
	}

	if p.TargetSubscriptionId != nil && (input.TargetSubscriptionId == nil || *p.TargetSubscriptionId != *input.TargetSubscriptionId) {
		return false
	}

	if p.UserName != nil && (input.UserName == nil || *p.UserName != *input.UserName) {
		return false
	}

	return true
}

type BastionHostOperationPredicate struct {
	Etag     *string
	Id       *string
	Location *string
	Name     *string
	Type     *string
}

func (p BastionHostOperationPredicate) Matches(input BastionHost) bool {

	if p.Etag != nil && (input.Etag == nil || *p.Etag != *input.Etag) {
		return false
	}

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Location != nil && (input.Location == nil || *p.Location != *input.Location) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}

type BastionSessionStateOperationPredicate struct {
	Message   *string
	SessionId *string
	State     *string
}

func (p BastionSessionStateOperationPredicate) Matches(input BastionSessionState) bool {

	if p.Message != nil && (input.Message == nil || *p.Message != *input.Message) {
		return false
	}

	if p.SessionId != nil && (input.SessionId == nil || *p.SessionId != *input.SessionId) {
		return false
	}

	if p.State != nil && (input.State == nil || *p.State != *input.State) {
		return false
	}

	return true
}

type BastionShareableLinkOperationPredicate struct {
	Bsl       *string
	CreatedAt *string
	Message   *string
}

func (p BastionShareableLinkOperationPredicate) Matches(input BastionShareableLink) bool {

	if p.Bsl != nil && (input.Bsl == nil || *p.Bsl != *input.Bsl) {
		return false
	}

	if p.CreatedAt != nil && (input.CreatedAt == nil || *p.CreatedAt != *input.CreatedAt) {
		return false
	}

	if p.Message != nil && (input.Message == nil || *p.Message != *input.Message) {
		return false
	}

	return true
}
