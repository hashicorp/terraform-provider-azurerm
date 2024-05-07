package jobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContainerAppProbeHTTPGet struct {
	HTTPHeaders *[]ContainerAppProbeHTTPGetHTTPHeadersInlined `json:"httpHeaders,omitempty"`
	Host        *string                                       `json:"host,omitempty"`
	Path        *string                                       `json:"path,omitempty"`
	Port        int64                                         `json:"port"`
	Scheme      *Scheme                                       `json:"scheme,omitempty"`
}
