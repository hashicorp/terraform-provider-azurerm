package containerinstance

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContainerHTTPGet struct {
	HTTPHeaders *[]HTTPHeader `json:"httpHeaders,omitempty"`
	Path        *string       `json:"path,omitempty"`
	Port        int64         `json:"port"`
	Scheme      *Scheme       `json:"scheme,omitempty"`
}
