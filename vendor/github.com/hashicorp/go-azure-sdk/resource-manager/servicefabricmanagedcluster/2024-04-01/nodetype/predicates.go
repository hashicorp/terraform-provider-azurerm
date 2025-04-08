package nodetype

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NodeTypeOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p NodeTypeOperationPredicate) Matches(input NodeType) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}

type NodeTypeAvailableSkuOperationPredicate struct {
	ResourceType *string
}

func (p NodeTypeAvailableSkuOperationPredicate) Matches(input NodeTypeAvailableSku) bool {

	if p.ResourceType != nil && (input.ResourceType == nil || *p.ResourceType != *input.ResourceType) {
		return false
	}

	return true
}
