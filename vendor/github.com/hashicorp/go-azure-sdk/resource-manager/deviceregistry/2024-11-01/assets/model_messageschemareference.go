package assets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MessageSchemaReference struct {
	SchemaName              string `json:"schemaName"`
	SchemaRegistryNamespace string `json:"schemaRegistryNamespace"`
	SchemaVersion           string `json:"schemaVersion"`
}
