package managedcluster

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedClusterUpdateParameters struct {
	Tags *map[string]string `json:"tags,omitempty"`
}
