package managementpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagementPolicyFilter struct {
	BlobIndexMatch *[]TagFilter `json:"blobIndexMatch,omitempty"`
	BlobTypes      []string     `json:"blobTypes"`
	PrefixMatch    *[]string    `json:"prefixMatch,omitempty"`
}
