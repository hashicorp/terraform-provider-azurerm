package connectionmonitors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectionMonitorHTTPConfiguration struct {
	Method                *HTTPConfigurationMethod `json:"method,omitempty"`
	Path                  *string                  `json:"path,omitempty"`
	Port                  *int64                   `json:"port,omitempty"`
	PreferHTTPS           *bool                    `json:"preferHTTPS,omitempty"`
	RequestHeaders        *[]HTTPHeader            `json:"requestHeaders,omitempty"`
	ValidStatusCodeRanges *[]string                `json:"validStatusCodeRanges,omitempty"`
}
