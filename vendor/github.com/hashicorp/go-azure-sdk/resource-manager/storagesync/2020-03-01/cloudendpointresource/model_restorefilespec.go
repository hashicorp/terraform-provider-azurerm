package cloudendpointresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RestoreFileSpec struct {
	Isdir *bool   `json:"isdir,omitempty"`
	Path  *string `json:"path,omitempty"`
}
