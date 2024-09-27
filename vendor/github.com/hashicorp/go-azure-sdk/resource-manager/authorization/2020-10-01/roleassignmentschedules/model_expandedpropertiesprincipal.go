package roleassignmentschedules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpandedPropertiesPrincipal struct {
	DisplayName *string `json:"displayName,omitempty"`
	Email       *string `json:"email,omitempty"`
	Id          *string `json:"id,omitempty"`
	Type        *string `json:"type,omitempty"`
}
