package agentversions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AgentVersion struct {
	AgentVersion *string `json:"agentVersion,omitempty"`
	DownloadLink *string `json:"downloadLink,omitempty"`
	OsType       *string `json:"osType,omitempty"`
}
