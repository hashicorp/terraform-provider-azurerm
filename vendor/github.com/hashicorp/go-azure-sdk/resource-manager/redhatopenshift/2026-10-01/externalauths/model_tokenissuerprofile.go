package externalauths

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TokenIssuerProfile struct {
	Audiences []string `json:"audiences"`
	Ca        *string  `json:"ca,omitempty"`
	Url       string   `json:"url"`
}
