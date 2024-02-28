package streamingjobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeltaSerializationProperties struct {
	DeltaTablePath   string    `json:"deltaTablePath"`
	PartitionColumns *[]string `json:"partitionColumns,omitempty"`
}
