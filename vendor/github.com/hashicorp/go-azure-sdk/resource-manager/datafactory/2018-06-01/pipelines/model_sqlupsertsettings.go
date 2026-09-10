package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SqlUpsertSettings struct {
	InterimSchemaName *interface{} `json:"interimSchemaName,omitempty"`
	Keys              *interface{} `json:"keys,omitempty"`
	UseTempDB         *interface{} `json:"useTempDB,omitempty"`
}
