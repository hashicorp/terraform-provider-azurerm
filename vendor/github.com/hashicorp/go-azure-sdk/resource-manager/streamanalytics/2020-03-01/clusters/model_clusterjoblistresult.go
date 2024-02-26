package clusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterJobListResult struct {
	NextLink *string       `json:"nextLink,omitempty"`
	Value    *[]ClusterJob `json:"value,omitempty"`
}
