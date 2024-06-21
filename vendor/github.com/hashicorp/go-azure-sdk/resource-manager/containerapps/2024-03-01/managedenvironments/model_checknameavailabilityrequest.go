package managedenvironments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CheckNameAvailabilityRequest struct {
	Name *string `json:"name,omitempty"`
	Type *string `json:"type,omitempty"`
}
