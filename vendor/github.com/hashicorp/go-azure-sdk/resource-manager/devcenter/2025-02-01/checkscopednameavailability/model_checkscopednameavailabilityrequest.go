package checkscopednameavailability

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CheckScopedNameAvailabilityRequest struct {
	Name  *string `json:"name,omitempty"`
	Scope *string `json:"scope,omitempty"`
	Type  *string `json:"type,omitempty"`
}
