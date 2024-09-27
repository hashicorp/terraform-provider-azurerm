package cognitiveservicesaccounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccountOperationPredicate struct {
	Etag     *string
	Id       *string
	Kind     *string
	Location *string
	Name     *string
	Type     *string
}

func (p AccountOperationPredicate) Matches(input Account) bool {

	if p.Etag != nil && (input.Etag == nil || *p.Etag != *input.Etag) {
		return false
	}

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Kind != nil && (input.Kind == nil || *p.Kind != *input.Kind) {
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

type AccountModelOperationPredicate struct {
	Format           *string
	IsDefaultVersion *bool
	MaxCapacity      *int64
	Name             *string
	Source           *string
	Version          *string
}

func (p AccountModelOperationPredicate) Matches(input AccountModel) bool {

	if p.Format != nil && (input.Format == nil || *p.Format != *input.Format) {
		return false
	}

	if p.IsDefaultVersion != nil && (input.IsDefaultVersion == nil || *p.IsDefaultVersion != *input.IsDefaultVersion) {
		return false
	}

	if p.MaxCapacity != nil && (input.MaxCapacity == nil || *p.MaxCapacity != *input.MaxCapacity) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.Source != nil && (input.Source == nil || *p.Source != *input.Source) {
		return false
	}

	if p.Version != nil && (input.Version == nil || *p.Version != *input.Version) {
		return false
	}

	return true
}

type ResourceSkuOperationPredicate struct {
	Kind         *string
	Name         *string
	ResourceType *string
	Tier         *string
}

func (p ResourceSkuOperationPredicate) Matches(input ResourceSku) bool {

	if p.Kind != nil && (input.Kind == nil || *p.Kind != *input.Kind) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.ResourceType != nil && (input.ResourceType == nil || *p.ResourceType != *input.ResourceType) {
		return false
	}

	if p.Tier != nil && (input.Tier == nil || *p.Tier != *input.Tier) {
		return false
	}

	return true
}

type UsageOperationPredicate struct {
	CurrentValue  *float64
	Limit         *float64
	NextResetTime *string
	QuotaPeriod   *string
}

func (p UsageOperationPredicate) Matches(input Usage) bool {

	if p.CurrentValue != nil && (input.CurrentValue == nil || *p.CurrentValue != *input.CurrentValue) {
		return false
	}

	if p.Limit != nil && (input.Limit == nil || *p.Limit != *input.Limit) {
		return false
	}

	if p.NextResetTime != nil && (input.NextResetTime == nil || *p.NextResetTime != *input.NextResetTime) {
		return false
	}

	if p.QuotaPeriod != nil && (input.QuotaPeriod == nil || *p.QuotaPeriod != *input.QuotaPeriod) {
		return false
	}

	return true
}
