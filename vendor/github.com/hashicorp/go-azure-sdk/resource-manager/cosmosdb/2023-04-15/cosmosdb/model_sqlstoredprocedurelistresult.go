package cosmosdb

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SqlStoredProcedureListResult struct {
	Value *[]SqlStoredProcedureGetResults `json:"value,omitempty"`
}
