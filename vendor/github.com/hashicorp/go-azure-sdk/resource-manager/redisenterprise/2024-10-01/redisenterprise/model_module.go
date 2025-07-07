package redisenterprise

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Module struct {
	Args    *string `json:"args,omitempty"`
	Name    string  `json:"name"`
	Version *string `json:"version,omitempty"`
}
