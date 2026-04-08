package afdendpoints

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AfdPurgeParameters struct {
	ContentPaths []string  `json:"contentPaths"`
	Domains      *[]string `json:"domains,omitempty"`
}
