package secrets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SecretItemOperationPredicate struct {
	ContentType *string
	Id          *string
	Managed     *bool
}

func (p SecretItemOperationPredicate) Matches(input SecretItem) bool {

	if p.ContentType != nil && (input.ContentType == nil || *p.ContentType != *input.ContentType) {
		return false
	}

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Managed != nil && (input.Managed == nil || *p.Managed != *input.Managed) {
		return false
	}

	return true
}
