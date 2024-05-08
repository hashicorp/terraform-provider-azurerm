package hdinsights

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterConfigFile struct {
	Content  *string            `json:"content,omitempty"`
	Encoding *ContentEncoding   `json:"encoding,omitempty"`
	FileName string             `json:"fileName"`
	Path     *string            `json:"path,omitempty"`
	Values   *map[string]string `json:"values,omitempty"`
}
