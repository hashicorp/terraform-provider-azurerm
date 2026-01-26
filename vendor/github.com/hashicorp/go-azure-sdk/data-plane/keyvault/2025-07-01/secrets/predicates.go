package secrets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeletedSecretItemOperationPredicate struct {
	ContentType        *string
	DeletedDate        *int64
	Id                 *string
	Managed            *bool
	RecoveryId         *string
	ScheduledPurgeDate *int64
}

func (p DeletedSecretItemOperationPredicate) Matches(input DeletedSecretItem) bool {

	if p.ContentType != nil && (input.ContentType == nil || *p.ContentType != *input.ContentType) {
		return false
	}

	if p.DeletedDate != nil && (input.DeletedDate == nil || *p.DeletedDate != *input.DeletedDate) {
		return false
	}

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Managed != nil && (input.Managed == nil || *p.Managed != *input.Managed) {
		return false
	}

	if p.RecoveryId != nil && (input.RecoveryId == nil || *p.RecoveryId != *input.RecoveryId) {
		return false
	}

	if p.ScheduledPurgeDate != nil && (input.ScheduledPurgeDate == nil || *p.ScheduledPurgeDate != *input.ScheduledPurgeDate) {
		return false
	}

	return true
}

type SecretItemOperationPredicate struct {
	ContentType *string
	Id          *string
	Managed     *bool
}

func (p SecretItemOperationPredicate) Matches(input SecretItem) bool {

	if p.ContentType != nil && (input.ContentType == nil || *p.ContentType != *input.ContentType) {
		return false
	}

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Managed != nil && (input.Managed == nil || *p.Managed != *input.Managed) {
		return false
	}

	return true
}
