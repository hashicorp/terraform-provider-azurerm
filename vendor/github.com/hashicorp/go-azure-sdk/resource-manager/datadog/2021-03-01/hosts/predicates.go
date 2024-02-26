package hosts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DatadogHostListResponseOperationPredicate struct {
	NextLink *string
}

func (p DatadogHostListResponseOperationPredicate) Matches(input DatadogHostListResponse) bool {

	if p.NextLink != nil && (input.NextLink == nil || *p.NextLink != *input.NextLink) {
		return false
	}

	return true
}
