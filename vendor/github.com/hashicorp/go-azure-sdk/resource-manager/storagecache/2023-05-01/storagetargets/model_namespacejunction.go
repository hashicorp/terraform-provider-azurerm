package storagetargets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NamespaceJunction struct {
	NamespacePath   *string `json:"namespacePath,omitempty"`
	NfsAccessPolicy *string `json:"nfsAccessPolicy,omitempty"`
	NfsExport       *string `json:"nfsExport,omitempty"`
	TargetPath      *string `json:"targetPath,omitempty"`
}
