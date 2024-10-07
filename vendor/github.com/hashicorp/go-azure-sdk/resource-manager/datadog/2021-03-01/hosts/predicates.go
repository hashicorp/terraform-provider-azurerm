package hosts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DatadogHostOperationPredicate struct {
	Name *string
}

func (p DatadogHostOperationPredicate) Matches(input DatadogHost) bool {

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	return true
}
