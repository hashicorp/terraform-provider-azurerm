package virtualmachineimages

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineImageResourceOperationPredicate struct {
	Id       *string
	Location *string
	Name     *string
}

func (p VirtualMachineImageResourceOperationPredicate) Matches(input VirtualMachineImageResource) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Location != nil && *p.Location != input.Location {
		return false
	}

	if p.Name != nil && *p.Name != input.Name {
		return false
	}

	return true
}
