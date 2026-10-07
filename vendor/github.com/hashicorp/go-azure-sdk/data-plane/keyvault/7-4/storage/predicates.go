package storage

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SasDefinitionItemOperationPredicate struct {
	Id  *string
	Sid *string
}

func (p SasDefinitionItemOperationPredicate) Matches(input SasDefinitionItem) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Sid != nil && (input.Sid == nil || *p.Sid != *input.Sid) {
		return false
	}

	return true
}

type StorageAccountItemOperationPredicate struct {
	Id         *string
	ResourceId *string
}

func (p StorageAccountItemOperationPredicate) Matches(input StorageAccountItem) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.ResourceId != nil && (input.ResourceId == nil || *p.ResourceId != *input.ResourceId) {
		return false
	}

	return true
}
