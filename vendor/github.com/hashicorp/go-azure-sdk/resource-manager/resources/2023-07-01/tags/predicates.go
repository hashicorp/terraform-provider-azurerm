package tags

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TagDetailsOperationPredicate struct {
	Id      *string
	TagName *string
}

func (p TagDetailsOperationPredicate) Matches(input TagDetails) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.TagName != nil && (input.TagName == nil || *p.TagName != *input.TagName) {
		return false
	}

	return true
}
