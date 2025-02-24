package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SSISPackageLocation struct {
	PackagePath    *string                            `json:"packagePath,omitempty"`
	Type           *SsisPackageLocationType           `json:"type,omitempty"`
	TypeProperties *SSISPackageLocationTypeProperties `json:"typeProperties,omitempty"`
}
