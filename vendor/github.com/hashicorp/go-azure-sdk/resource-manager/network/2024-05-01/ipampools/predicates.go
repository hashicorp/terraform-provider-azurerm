package ipampools

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IPamPoolOperationPredicate struct {
	Id       *string
	Location *string
	Name     *string
	Type     *string
}

func (p IPamPoolOperationPredicate) Matches(input IPamPool) bool {

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

type PoolAssociationOperationPredicate struct {
	CreatedAt                   *string
	Description                 *string
	NumberOfReservedIPAddresses *string
	PoolId                      *string
	ReservationExpiresAt        *string
	ResourceId                  *string
	TotalNumberOfIPAddresses    *string
}

func (p PoolAssociationOperationPredicate) Matches(input PoolAssociation) bool {

	if p.CreatedAt != nil && (input.CreatedAt == nil || *p.CreatedAt != *input.CreatedAt) {
		return false
	}

	if p.Description != nil && (input.Description == nil || *p.Description != *input.Description) {
		return false
	}

	if p.NumberOfReservedIPAddresses != nil && (input.NumberOfReservedIPAddresses == nil || *p.NumberOfReservedIPAddresses != *input.NumberOfReservedIPAddresses) {
		return false
	}

	if p.PoolId != nil && (input.PoolId == nil || *p.PoolId != *input.PoolId) {
		return false
	}

	if p.ReservationExpiresAt != nil && (input.ReservationExpiresAt == nil || *p.ReservationExpiresAt != *input.ReservationExpiresAt) {
		return false
	}

	if p.ResourceId != nil && *p.ResourceId != input.ResourceId {
		return false
	}

	if p.TotalNumberOfIPAddresses != nil && (input.TotalNumberOfIPAddresses == nil || *p.TotalNumberOfIPAddresses != *input.TotalNumberOfIPAddresses) {
		return false
	}

	return true
}
