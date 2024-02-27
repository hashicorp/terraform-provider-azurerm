package apikey

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DatadogApiKeyOperationPredicate struct {
	Created   *string
	CreatedBy *string
	Key       *string
	Name      *string
}

func (p DatadogApiKeyOperationPredicate) Matches(input DatadogApiKey) bool {

	if p.Created != nil && (input.Created == nil || *p.Created != *input.Created) {
		return false
	}

	if p.CreatedBy != nil && (input.CreatedBy == nil || *p.CreatedBy != *input.CreatedBy) {
		return false
	}

	if p.Key != nil && *p.Key != input.Key {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	return true
}
