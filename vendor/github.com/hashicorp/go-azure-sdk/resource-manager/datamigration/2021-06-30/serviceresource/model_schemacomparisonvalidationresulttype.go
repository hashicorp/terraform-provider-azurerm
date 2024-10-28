package serviceresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SchemaComparisonValidationResultType struct {
	ObjectName   *string           `json:"objectName,omitempty"`
	ObjectType   *ObjectType       `json:"objectType,omitempty"`
	UpdateAction *UpdateActionType `json:"updateAction,omitempty"`
}
