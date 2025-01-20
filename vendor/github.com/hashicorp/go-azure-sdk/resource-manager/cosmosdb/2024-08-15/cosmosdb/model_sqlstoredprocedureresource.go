package cosmosdb

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SqlStoredProcedureResource struct {
	Body *string `json:"body,omitempty"`
	Id   string  `json:"id"`
}
