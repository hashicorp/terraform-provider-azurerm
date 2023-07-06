package deployments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeploymentModel struct {
	CallRateLimit *CallRateLimit `json:"callRateLimit,omitempty"`
	Format        *string        `json:"format,omitempty"`
	Name          *string        `json:"name,omitempty"`
	Source        *string        `json:"source,omitempty"`
	Version       *string        `json:"version,omitempty"`
}
