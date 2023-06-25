package cluster

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EndpointRangeDescription struct {
	EndPort   int64 `json:"endPort"`
	StartPort int64 `json:"startPort"`
}
