package bastionhosts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BastionShareableLink struct {
	Bsl       *string  `json:"bsl,omitempty"`
	CreatedAt *string  `json:"createdAt,omitempty"`
	Message   *string  `json:"message,omitempty"`
	VM        Resource `json:"vm"`
}
