package onlineendpoint

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PartialMinimalTrackedResourceWithIdentity struct {
	Identity *PartialManagedServiceIdentity `json:"identity,omitempty"`
	Tags     *map[string]string             `json:"tags,omitempty"`
}
