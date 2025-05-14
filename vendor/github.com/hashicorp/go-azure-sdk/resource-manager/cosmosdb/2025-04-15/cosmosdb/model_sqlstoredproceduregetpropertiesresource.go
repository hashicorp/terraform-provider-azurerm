package cosmosdb

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SqlStoredProcedureGetPropertiesResource struct {
	Body *string  `json:"body,omitempty"`
	Etag *string  `json:"_etag,omitempty"`
	Id   *string  `json:"id,omitempty"`
	Rid  *string  `json:"_rid,omitempty"`
	Ts   *float64 `json:"_ts,omitempty"`
}
