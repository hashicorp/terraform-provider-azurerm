package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SqlUpsertSettings struct {
	InterimSchemaName *string   `json:"interimSchemaName,omitempty"`
	Keys              *[]string `json:"keys,omitempty"`
	UseTempDB         *bool     `json:"useTempDB,omitempty"`
}
