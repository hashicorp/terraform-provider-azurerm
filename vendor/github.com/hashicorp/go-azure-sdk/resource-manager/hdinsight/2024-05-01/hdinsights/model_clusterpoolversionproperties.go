package hdinsights

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterPoolVersionProperties struct {
	AksVersion         *string `json:"aksVersion,omitempty"`
	ClusterPoolVersion *string `json:"clusterPoolVersion,omitempty"`
	IsPreview          *bool   `json:"isPreview,omitempty"`
}
