package datasources

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataSource struct {
	Etag       *string            `json:"etag,omitempty"`
	Id         *string            `json:"id,omitempty"`
	Kind       DataSourceKind     `json:"kind"`
	Name       *string            `json:"name,omitempty"`
	Properties interface{}        `json:"properties"`
	Tags       *map[string]string `json:"tags,omitempty"`
	Type       *string            `json:"type,omitempty"`
}
