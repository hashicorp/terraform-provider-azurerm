package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HostKeys struct {
	FunctionKeys *map[string]string `json:"functionKeys,omitempty"`
	MasterKey    *string            `json:"masterKey,omitempty"`
	SystemKeys   *map[string]string `json:"systemKeys,omitempty"`
}
