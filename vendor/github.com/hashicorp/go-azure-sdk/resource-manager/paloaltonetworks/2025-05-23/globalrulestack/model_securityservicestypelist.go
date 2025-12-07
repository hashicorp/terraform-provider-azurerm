package globalrulestack

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SecurityServicesTypeList struct {
	Entry []NameDescriptionObject `json:"entry"`
	Type  *string                 `json:"type,omitempty"`
}
