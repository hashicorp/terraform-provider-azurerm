package networkwatchers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HTTPConfiguration struct {
	Headers          *[]HTTPHeader `json:"headers,omitempty"`
	Method           *HTTPMethod   `json:"method,omitempty"`
	ValidStatusCodes *[]int64      `json:"validStatusCodes,omitempty"`
}
