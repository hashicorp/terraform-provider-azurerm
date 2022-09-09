package cosmosdb

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Indexes struct {
	DataType  *DataType  `json:"dataType,omitempty"`
	Kind      *IndexKind `json:"kind,omitempty"`
	Precision *int64     `json:"precision,omitempty"`
}
