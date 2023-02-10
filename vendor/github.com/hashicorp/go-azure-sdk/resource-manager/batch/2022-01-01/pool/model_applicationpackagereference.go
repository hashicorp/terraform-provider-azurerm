package pool

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationPackageReference struct {
	Id      string  `json:"id"`
	Version *string `json:"version,omitempty"`
}
