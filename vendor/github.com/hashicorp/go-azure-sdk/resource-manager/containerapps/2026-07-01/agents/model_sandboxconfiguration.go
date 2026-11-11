package agents

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SandboxConfiguration struct {
	Egress   *SandboxEgressConfiguration `json:"egress,omitempty"`
	Packages *[]SandboxPackageSpec       `json:"packages,omitempty"`
}
