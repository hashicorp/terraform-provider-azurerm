package cosmosdb

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SqlContainerResource struct {
	AnalyticalStorageTtl     *int64                    `json:"analyticalStorageTtl,omitempty"`
	ClientEncryptionPolicy   *ClientEncryptionPolicy   `json:"clientEncryptionPolicy,omitempty"`
	ComputedProperties       *[]ComputedProperty       `json:"computedProperties,omitempty"`
	ConflictResolutionPolicy *ConflictResolutionPolicy `json:"conflictResolutionPolicy,omitempty"`
	CreateMode               *CreateMode               `json:"createMode,omitempty"`
	DefaultTtl               *int64                    `json:"defaultTtl,omitempty"`
	Id                       string                    `json:"id"`
	IndexingPolicy           *IndexingPolicy           `json:"indexingPolicy,omitempty"`
	PartitionKey             *ContainerPartitionKey    `json:"partitionKey,omitempty"`
	RestoreParameters        *RestoreParametersBase    `json:"restoreParameters,omitempty"`
	UniqueKeyPolicy          *UniqueKeyPolicy          `json:"uniqueKeyPolicy,omitempty"`
}
