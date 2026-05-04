package servicetags

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServiceTagInformationOperationPredicate struct {
	Id                     *string
	Name                   *string
	ServiceTagChangeNumber *string
}

func (p ServiceTagInformationOperationPredicate) Matches(input ServiceTagInformation) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.ServiceTagChangeNumber != nil && (input.ServiceTagChangeNumber == nil || *p.ServiceTagChangeNumber != *input.ServiceTagChangeNumber) {
		return false
	}

	return true
}
