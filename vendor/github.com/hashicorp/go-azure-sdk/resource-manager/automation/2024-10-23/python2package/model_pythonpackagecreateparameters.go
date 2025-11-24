package python2package

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PythonPackageCreateParameters struct {
	Properties PythonPackageCreateProperties `json:"properties"`
	Tags       *map[string]string            `json:"tags,omitempty"`
}
