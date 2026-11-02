package deletedstorage

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeletedSasDefinitionItemOperationPredicate struct {
	DeletedDate        *int64
	Id                 *string
	RecoveryId         *string
	ScheduledPurgeDate *int64
	Sid                *string
}

func (p DeletedSasDefinitionItemOperationPredicate) Matches(input DeletedSasDefinitionItem) bool {

	if p.DeletedDate != nil && (input.DeletedDate == nil || *p.DeletedDate != *input.DeletedDate) {
		return false
	}

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.RecoveryId != nil && (input.RecoveryId == nil || *p.RecoveryId != *input.RecoveryId) {
		return false
	}

	if p.ScheduledPurgeDate != nil && (input.ScheduledPurgeDate == nil || *p.ScheduledPurgeDate != *input.ScheduledPurgeDate) {
		return false
	}

	if p.Sid != nil && (input.Sid == nil || *p.Sid != *input.Sid) {
		return false
	}

	return true
}

type DeletedStorageAccountItemOperationPredicate struct {
	DeletedDate        *int64
	Id                 *string
	RecoveryId         *string
	ResourceId         *string
	ScheduledPurgeDate *int64
}

func (p DeletedStorageAccountItemOperationPredicate) Matches(input DeletedStorageAccountItem) bool {

	if p.DeletedDate != nil && (input.DeletedDate == nil || *p.DeletedDate != *input.DeletedDate) {
		return false
	}

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.RecoveryId != nil && (input.RecoveryId == nil || *p.RecoveryId != *input.RecoveryId) {
		return false
	}

	if p.ResourceId != nil && (input.ResourceId == nil || *p.ResourceId != *input.ResourceId) {
		return false
	}

	if p.ScheduledPurgeDate != nil && (input.ScheduledPurgeDate == nil || *p.ScheduledPurgeDate != *input.ScheduledPurgeDate) {
		return false
	}

	return true
}
