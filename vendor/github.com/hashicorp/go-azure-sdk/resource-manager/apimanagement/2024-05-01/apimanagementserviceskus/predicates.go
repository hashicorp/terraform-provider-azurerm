package apimanagementserviceskus

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceSkuResultOperationPredicate struct {
	ResourceType *string
}

func (p ResourceSkuResultOperationPredicate) Matches(input ResourceSkuResult) bool {

	if p.ResourceType != nil && (input.ResourceType == nil || *p.ResourceType != *input.ResourceType) {
		return false
	}

	return true
}
