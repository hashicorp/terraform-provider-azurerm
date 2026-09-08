package webhooks

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Request struct {
	Addr      *string `json:"addr,omitempty"`
	Host      *string `json:"host,omitempty"`
	Id        *string `json:"id,omitempty"`
	Method    *string `json:"method,omitempty"`
	Useragent *string `json:"useragent,omitempty"`
}
