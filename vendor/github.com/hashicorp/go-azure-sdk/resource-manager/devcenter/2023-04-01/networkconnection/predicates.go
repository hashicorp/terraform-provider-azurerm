package networkconnection

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OutboundEnvironmentEndpointOperationPredicate struct {
	Category *string
}

func (p OutboundEnvironmentEndpointOperationPredicate) Matches(input OutboundEnvironmentEndpoint) bool {

	if p.Category != nil && (input.Category == nil || *p.Category != *input.Category) {
		return false
	}

	return true
}
