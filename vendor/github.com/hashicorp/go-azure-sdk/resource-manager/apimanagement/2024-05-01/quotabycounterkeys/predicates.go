package quotabycounterkeys

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type QuotaCounterContractOperationPredicate struct {
	CounterKey      *string
	PeriodEndTime   *string
	PeriodKey       *string
	PeriodStartTime *string
}

func (p QuotaCounterContractOperationPredicate) Matches(input QuotaCounterContract) bool {

	if p.CounterKey != nil && *p.CounterKey != input.CounterKey {
		return false
	}

	if p.PeriodEndTime != nil && *p.PeriodEndTime != input.PeriodEndTime {
		return false
	}

	if p.PeriodKey != nil && *p.PeriodKey != input.PeriodKey {
		return false
	}

	if p.PeriodStartTime != nil && *p.PeriodStartTime != input.PeriodStartTime {
		return false
	}

	return true
}
