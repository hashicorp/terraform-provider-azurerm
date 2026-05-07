package packageresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PackageUpdateParameters struct {
	AllOf      *TrackedResource         `json:"allOf,omitempty"`
	Properties *PackageUpdateProperties `json:"properties,omitempty"`
}
