package fileservice

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FileServicePropertiesProperties struct {
	Cors                       *CorsRules             `json:"cors,omitempty"`
	ProtocolSettings           *ProtocolSettings      `json:"protocolSettings,omitempty"`
	ShareDeleteRetentionPolicy *DeleteRetentionPolicy `json:"shareDeleteRetentionPolicy,omitempty"`
}
