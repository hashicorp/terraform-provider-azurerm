package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SSISChildPackage struct {
	PackageContent          interface{} `json:"packageContent"`
	PackageLastModifiedDate *string     `json:"packageLastModifiedDate,omitempty"`
	PackageName             *string     `json:"packageName,omitempty"`
	PackagePath             interface{} `json:"packagePath"`
}
