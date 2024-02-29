package elasticsanskus

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SkuInformationOperationPredicate struct {
	ResourceType *string
}

func (p SkuInformationOperationPredicate) Matches(input SkuInformation) bool {

	if p.ResourceType != nil && (input.ResourceType == nil || *p.ResourceType != *input.ResourceType) {
		return false
	}

	return true
}
