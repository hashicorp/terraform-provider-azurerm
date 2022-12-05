package streamingpoliciesandstreaminglocators

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListPathsResponse struct {
	DownloadPaths  *[]string        `json:"downloadPaths,omitempty"`
	StreamingPaths *[]StreamingPath `json:"streamingPaths,omitempty"`
}
