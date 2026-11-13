package agents

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SandboxPackageSpec struct {
	Name           string                `json:"name"`
	PackageManager SandboxPackageManager `json:"packageManager"`
	Version        *string               `json:"version,omitempty"`
}
