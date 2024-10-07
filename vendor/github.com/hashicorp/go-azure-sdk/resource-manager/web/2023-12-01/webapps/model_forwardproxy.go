package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ForwardProxy struct {
	Convention            *ForwardProxyConvention `json:"convention,omitempty"`
	CustomHostHeaderName  *string                 `json:"customHostHeaderName,omitempty"`
	CustomProtoHeaderName *string                 `json:"customProtoHeaderName,omitempty"`
}
