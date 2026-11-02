package keys

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

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
