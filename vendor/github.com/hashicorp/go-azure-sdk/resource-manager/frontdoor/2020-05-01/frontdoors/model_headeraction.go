package frontdoors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HeaderAction struct {
	HeaderActionType HeaderActionType `json:"headerActionType"`
	HeaderName       string           `json:"headerName"`
	Value            *string          `json:"value,omitempty"`
}
