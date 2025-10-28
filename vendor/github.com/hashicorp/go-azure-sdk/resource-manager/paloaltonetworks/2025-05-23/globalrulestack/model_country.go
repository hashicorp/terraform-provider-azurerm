package globalrulestack

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Country struct {
	Code        string  `json:"code"`
	Description *string `json:"description,omitempty"`
}
