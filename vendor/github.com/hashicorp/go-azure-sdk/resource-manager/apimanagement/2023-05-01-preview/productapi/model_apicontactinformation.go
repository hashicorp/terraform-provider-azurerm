package productapi

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApiContactInformation struct {
	Email *string `json:"email,omitempty"`
	Name  *string `json:"name,omitempty"`
	Url   *string `json:"url,omitempty"`
}
