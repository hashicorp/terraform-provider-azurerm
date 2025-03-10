package performconnectivitycheck

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectivityCheckRequestProtocolConfigurationHTTPConfiguration struct {
	Headers          *[]HTTPHeader `json:"headers,omitempty"`
	Method           *Method       `json:"method,omitempty"`
	ValidStatusCodes *[]int64      `json:"validStatusCodes,omitempty"`
}
