package softwareupdateconfigurationrun

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SoftwareUpdateConfigurationRunOperationPredicate struct {
	Id   *string
	Name *string
}

func (p SoftwareUpdateConfigurationRunOperationPredicate) Matches(input SoftwareUpdateConfigurationRun) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	return true
}
