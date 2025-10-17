package organizationresources

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterStatusEntity struct {
	Cku   *int64  `json:"cku,omitempty"`
	Phase *string `json:"phase,omitempty"`
}
