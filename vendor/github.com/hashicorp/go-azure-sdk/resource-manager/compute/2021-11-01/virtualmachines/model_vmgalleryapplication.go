package virtualmachines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VMGalleryApplication struct {
	ConfigurationReference *string `json:"configurationReference,omitempty"`
	Order                  *int64  `json:"order,omitempty"`
	PackageReferenceId     string  `json:"packageReferenceId"`
	Tags                   *string `json:"tags,omitempty"`
}
