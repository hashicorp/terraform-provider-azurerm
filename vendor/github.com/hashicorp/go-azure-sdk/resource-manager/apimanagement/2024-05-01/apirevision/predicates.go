package apirevision

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApiRevisionContractOperationPredicate struct {
	ApiId           *string
	ApiRevision     *string
	CreatedDateTime *string
	Description     *string
	IsCurrent       *bool
	IsOnline        *bool
	PrivateURL      *string
	UpdatedDateTime *string
}

func (p ApiRevisionContractOperationPredicate) Matches(input ApiRevisionContract) bool {

	if p.ApiId != nil && (input.ApiId == nil || *p.ApiId != *input.ApiId) {
		return false
	}

	if p.ApiRevision != nil && (input.ApiRevision == nil || *p.ApiRevision != *input.ApiRevision) {
		return false
	}

	if p.CreatedDateTime != nil && (input.CreatedDateTime == nil || *p.CreatedDateTime != *input.CreatedDateTime) {
		return false
	}

	if p.Description != nil && (input.Description == nil || *p.Description != *input.Description) {
		return false
	}

	if p.IsCurrent != nil && (input.IsCurrent == nil || *p.IsCurrent != *input.IsCurrent) {
		return false
	}

	if p.IsOnline != nil && (input.IsOnline == nil || *p.IsOnline != *input.IsOnline) {
		return false
	}

	if p.PrivateURL != nil && (input.PrivateURL == nil || *p.PrivateURL != *input.PrivateURL) {
		return false
	}

	if p.UpdatedDateTime != nil && (input.UpdatedDateTime == nil || *p.UpdatedDateTime != *input.UpdatedDateTime) {
		return false
	}

	return true
}
