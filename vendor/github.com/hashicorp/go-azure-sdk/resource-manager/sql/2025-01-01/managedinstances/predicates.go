package managedinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedInstanceOperationPredicate struct {
	Id       *string
	Location *string
	Name     *string
	Type     *string
}

func (p ManagedInstanceOperationPredicate) Matches(input ManagedInstance) bool {

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

type OutboundEnvironmentEndpointOperationPredicate struct {
	Category *string
}

func (p OutboundEnvironmentEndpointOperationPredicate) Matches(input OutboundEnvironmentEndpoint) bool {

	if p.Category != nil && (input.Category == nil || *p.Category != *input.Category) {
		return false
	}

	return true
}

type TopQueriesOperationPredicate struct {
	AggregationFunction *string
	EndTime             *string
	NumberOfQueries     *int64
	ObservationMetric   *string
	StartTime           *string
}

func (p TopQueriesOperationPredicate) Matches(input TopQueries) bool {

	if p.AggregationFunction != nil && (input.AggregationFunction == nil || *p.AggregationFunction != *input.AggregationFunction) {
		return false
	}

	if p.EndTime != nil && (input.EndTime == nil || *p.EndTime != *input.EndTime) {
		return false
	}

	if p.NumberOfQueries != nil && (input.NumberOfQueries == nil || *p.NumberOfQueries != *input.NumberOfQueries) {
		return false
	}

	if p.ObservationMetric != nil && (input.ObservationMetric == nil || *p.ObservationMetric != *input.ObservationMetric) {
		return false
	}

	if p.StartTime != nil && (input.StartTime == nil || *p.StartTime != *input.StartTime) {
		return false
	}

	return true
}
