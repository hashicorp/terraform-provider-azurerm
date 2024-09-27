package appplatform

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SupportedServerVersion struct {
	Server  *string `json:"server,omitempty"`
	Value   *string `json:"value,omitempty"`
	Version *string `json:"version,omitempty"`
}
