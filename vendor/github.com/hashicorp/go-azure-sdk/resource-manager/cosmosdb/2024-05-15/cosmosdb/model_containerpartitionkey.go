package cosmosdb

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContainerPartitionKey struct {
	Kind      *PartitionKind `json:"kind,omitempty"`
	Paths     *[]string      `json:"paths,omitempty"`
	SystemKey *bool          `json:"systemKey,omitempty"`
	Version   *int64         `json:"version,omitempty"`
}
