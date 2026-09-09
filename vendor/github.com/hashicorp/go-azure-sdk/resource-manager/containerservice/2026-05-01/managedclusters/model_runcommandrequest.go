package managedclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RunCommandRequest struct {
	ClusterToken *string `json:"clusterToken,omitempty"`
	Command      string  `json:"command"`
	Context      *string `json:"context,omitempty"`
}
