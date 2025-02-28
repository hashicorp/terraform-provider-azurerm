package productapi

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SubscriptionKeyParameterNamesContract struct {
	Header *string `json:"header,omitempty"`
	Query  *string `json:"query,omitempty"`
}
