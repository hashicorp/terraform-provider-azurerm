package keys

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeletedKeyItemOperationPredicate struct {
	DeletedDate        *int64
	Kid                *string
	Managed            *bool
	RecoveryId         *string
	ScheduledPurgeDate *int64
}

func (p DeletedKeyItemOperationPredicate) Matches(input DeletedKeyItem) bool {

	if p.DeletedDate != nil && (input.DeletedDate == nil || *p.DeletedDate != *input.DeletedDate) {
		return false
	}

	if p.Kid != nil && (input.Kid == nil || *p.Kid != *input.Kid) {
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

type KeyItemOperationPredicate struct {
	Kid     *string
	Managed *bool
}

func (p KeyItemOperationPredicate) Matches(input KeyItem) bool {

	if p.Kid != nil && (input.Kid == nil || *p.Kid != *input.Kid) {
		return false
	}

	if p.Managed != nil && (input.Managed == nil || *p.Managed != *input.Managed) {
		return false
	}

	return true
}
