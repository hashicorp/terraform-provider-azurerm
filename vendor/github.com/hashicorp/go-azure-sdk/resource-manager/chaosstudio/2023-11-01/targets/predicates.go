package targets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TargetOperationPredicate struct {
	Id         *string
	Location   *string
	Name       *string
	Properties *interface{}
	Type       *string
}

func (p TargetOperationPredicate) Matches(input Target) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Location != nil && (input.Location == nil || *p.Location != *input.Location) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.Properties != nil && *p.Properties != input.Properties {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}
