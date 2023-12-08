package appplatform

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SupportedRuntimeVersion struct {
	Platform *SupportedRuntimePlatform `json:"platform,omitempty"`
	Value    *SupportedRuntimeValue    `json:"value,omitempty"`
	Version  *string                   `json:"version,omitempty"`
}
