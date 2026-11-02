package statistics

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StatisticsOperationPredicate struct {
	CounterProperty *string
	CounterValue    *int64
	EndTime         *string
	Id              *string
	StartTime       *string
}

func (p StatisticsOperationPredicate) Matches(input Statistics) bool {

	if p.CounterProperty != nil && (input.CounterProperty == nil || *p.CounterProperty != *input.CounterProperty) {
		return false
	}

	if p.CounterValue != nil && (input.CounterValue == nil || *p.CounterValue != *input.CounterValue) {
		return false
	}

	if p.EndTime != nil && (input.EndTime == nil || *p.EndTime != *input.EndTime) {
		return false
	}

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.StartTime != nil && (input.StartTime == nil || *p.StartTime != *input.StartTime) {
		return false
	}

	return true
}
