package reports

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReportRecordContractOperationPredicate struct {
	ApiId            *string
	ApiRegion        *string
	ApiTimeAvg       *float64
	ApiTimeMax       *float64
	ApiTimeMin       *float64
	Bandwidth        *int64
	CacheHitCount    *int64
	CacheMissCount   *int64
	CallCountBlocked *int64
	CallCountFailed  *int64
	CallCountOther   *int64
	CallCountSuccess *int64
	CallCountTotal   *int64
	Country          *string
	Interval         *string
	Name             *string
	OperationId      *string
	ProductId        *string
	Region           *string
	ServiceTimeAvg   *float64
	ServiceTimeMax   *float64
	ServiceTimeMin   *float64
	SubscriptionId   *string
	Timestamp        *string
	UserId           *string
	Zip              *string
}

func (p ReportRecordContractOperationPredicate) Matches(input ReportRecordContract) bool {

	if p.ApiId != nil && (input.ApiId == nil || *p.ApiId != *input.ApiId) {
		return false
	}

	if p.ApiRegion != nil && (input.ApiRegion == nil || *p.ApiRegion != *input.ApiRegion) {
		return false
	}

	if p.ApiTimeAvg != nil && (input.ApiTimeAvg == nil || *p.ApiTimeAvg != *input.ApiTimeAvg) {
		return false
	}

	if p.ApiTimeMax != nil && (input.ApiTimeMax == nil || *p.ApiTimeMax != *input.ApiTimeMax) {
		return false
	}

	if p.ApiTimeMin != nil && (input.ApiTimeMin == nil || *p.ApiTimeMin != *input.ApiTimeMin) {
		return false
	}

	if p.Bandwidth != nil && (input.Bandwidth == nil || *p.Bandwidth != *input.Bandwidth) {
		return false
	}

	if p.CacheHitCount != nil && (input.CacheHitCount == nil || *p.CacheHitCount != *input.CacheHitCount) {
		return false
	}

	if p.CacheMissCount != nil && (input.CacheMissCount == nil || *p.CacheMissCount != *input.CacheMissCount) {
		return false
	}

	if p.CallCountBlocked != nil && (input.CallCountBlocked == nil || *p.CallCountBlocked != *input.CallCountBlocked) {
		return false
	}

	if p.CallCountFailed != nil && (input.CallCountFailed == nil || *p.CallCountFailed != *input.CallCountFailed) {
		return false
	}

	if p.CallCountOther != nil && (input.CallCountOther == nil || *p.CallCountOther != *input.CallCountOther) {
		return false
	}

	if p.CallCountSuccess != nil && (input.CallCountSuccess == nil || *p.CallCountSuccess != *input.CallCountSuccess) {
		return false
	}

	if p.CallCountTotal != nil && (input.CallCountTotal == nil || *p.CallCountTotal != *input.CallCountTotal) {
		return false
	}

	if p.Country != nil && (input.Country == nil || *p.Country != *input.Country) {
		return false
	}

	if p.Interval != nil && (input.Interval == nil || *p.Interval != *input.Interval) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.OperationId != nil && (input.OperationId == nil || *p.OperationId != *input.OperationId) {
		return false
	}

	if p.ProductId != nil && (input.ProductId == nil || *p.ProductId != *input.ProductId) {
		return false
	}

	if p.Region != nil && (input.Region == nil || *p.Region != *input.Region) {
		return false
	}

	if p.ServiceTimeAvg != nil && (input.ServiceTimeAvg == nil || *p.ServiceTimeAvg != *input.ServiceTimeAvg) {
		return false
	}

	if p.ServiceTimeMax != nil && (input.ServiceTimeMax == nil || *p.ServiceTimeMax != *input.ServiceTimeMax) {
		return false
	}

	if p.ServiceTimeMin != nil && (input.ServiceTimeMin == nil || *p.ServiceTimeMin != *input.ServiceTimeMin) {
		return false
	}

	if p.SubscriptionId != nil && (input.SubscriptionId == nil || *p.SubscriptionId != *input.SubscriptionId) {
		return false
	}

	if p.Timestamp != nil && (input.Timestamp == nil || *p.Timestamp != *input.Timestamp) {
		return false
	}

	if p.UserId != nil && (input.UserId == nil || *p.UserId != *input.UserId) {
		return false
	}

	if p.Zip != nil && (input.Zip == nil || *p.Zip != *input.Zip) {
		return false
	}

	return true
}
