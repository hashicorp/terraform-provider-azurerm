package redis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RedisCacheAccessPolicyOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p RedisCacheAccessPolicyOperationPredicate) Matches(input RedisCacheAccessPolicy) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
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

type RedisCacheAccessPolicyAssignmentOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p RedisCacheAccessPolicyAssignmentOperationPredicate) Matches(input RedisCacheAccessPolicyAssignment) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
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

type RedisFirewallRuleOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p RedisFirewallRuleOperationPredicate) Matches(input RedisFirewallRule) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
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

type RedisLinkedServerWithPropertiesOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p RedisLinkedServerWithPropertiesOperationPredicate) Matches(input RedisLinkedServerWithProperties) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
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

type RedisPatchScheduleOperationPredicate struct {
	Id       *string
	Location *string
	Name     *string
	Type     *string
}

func (p RedisPatchScheduleOperationPredicate) Matches(input RedisPatchSchedule) bool {

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

type RedisResourceOperationPredicate struct {
	Id       *string
	Location *string
	Name     *string
	Type     *string
}

func (p RedisResourceOperationPredicate) Matches(input RedisResource) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Location != nil && *p.Location != input.Location {
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

type UpgradeNotificationOperationPredicate struct {
	Name      *string
	Timestamp *string
}

func (p UpgradeNotificationOperationPredicate) Matches(input UpgradeNotification) bool {

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.Timestamp != nil && (input.Timestamp == nil || *p.Timestamp != *input.Timestamp) {
		return false
	}

	return true
}
