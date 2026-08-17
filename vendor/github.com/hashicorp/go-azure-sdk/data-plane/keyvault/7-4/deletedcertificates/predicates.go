package deletedcertificates

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeletedCertificateItemOperationPredicate struct {
	DeletedDate        *int64
	Id                 *string
	RecoveryId         *string
	ScheduledPurgeDate *int64
	X5t                *string
}

func (p DeletedCertificateItemOperationPredicate) Matches(input DeletedCertificateItem) bool {

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

	if p.X5t != nil && (input.X5t == nil || *p.X5t != *input.X5t) {
		return false
	}

	return true
}
