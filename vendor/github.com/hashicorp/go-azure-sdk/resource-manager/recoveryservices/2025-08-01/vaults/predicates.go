package vaults

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReplicationUsageOperationPredicate struct {
	ProtectedItemCount               *int64
	RecoveryPlanCount                *int64
	RecoveryServicesProviderAuthType *int64
	RegisteredServersCount           *int64
}

func (p ReplicationUsageOperationPredicate) Matches(input ReplicationUsage) bool {

	if p.ProtectedItemCount != nil && (input.ProtectedItemCount == nil || *p.ProtectedItemCount != *input.ProtectedItemCount) {
		return false
	}

	if p.RecoveryPlanCount != nil && (input.RecoveryPlanCount == nil || *p.RecoveryPlanCount != *input.RecoveryPlanCount) {
		return false
	}

	if p.RecoveryServicesProviderAuthType != nil && (input.RecoveryServicesProviderAuthType == nil || *p.RecoveryServicesProviderAuthType != *input.RecoveryServicesProviderAuthType) {
		return false
	}

	if p.RegisteredServersCount != nil && (input.RegisteredServersCount == nil || *p.RegisteredServersCount != *input.RegisteredServersCount) {
		return false
	}

	return true
}

type VaultOperationPredicate struct {
	Etag     *string
	Id       *string
	Location *string
	Name     *string
	Type     *string
}

func (p VaultOperationPredicate) Matches(input Vault) bool {

	if p.Etag != nil && (input.Etag == nil || *p.Etag != *input.Etag) {
		return false
	}

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

type VaultUsageOperationPredicate struct {
	CurrentValue  *int64
	Limit         *int64
	NextResetTime *string
	QuotaPeriod   *string
}

func (p VaultUsageOperationPredicate) Matches(input VaultUsage) bool {

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
