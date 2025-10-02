package availableservicealiases

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AvailableServiceAliasOperationPredicate struct {
	Id           *string
	Name         *string
	ResourceName *string
	Type         *string
}

func (p AvailableServiceAliasOperationPredicate) Matches(input AvailableServiceAlias) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.ResourceName != nil && (input.ResourceName == nil || *p.ResourceName != *input.ResourceName) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}
